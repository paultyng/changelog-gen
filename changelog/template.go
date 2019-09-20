package changelog

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

const defaultReleaseNoteTemplate = `{{with .Labels | filterPrefix "service/" true | sortAlpha }}{{if len . | lt 0 }}**{{. | join ", " }}:** {{end}}{{end}}{{.Text }} ([{{.PRNumber}}]({{.PRURL}}) by [{{.Author}}]({{.AuthorURL}}))`
const defaultChangelogTemplate = `
{{- $breaking := newStringList -}}
{{- $features := newStringList -}}
{{- $improvements := newStringList -}}
{{- $bugs := newStringList -}}
{{- range . -}}
  {{if .BreakingChange -}}
	{{$breaking = append $breaking (renderReleaseNote .) -}}
  {{else if or (has "new-resource" .Labels) (has "new-data-source" .Labels) -}}
	{{$features = append $features (renderReleaseNote .) -}}
  {{else if not .Bug -}}
	{{$improvements = append $improvements (renderReleaseNote .) -}}
  {{else -}}
	{{$bugs = append $bugs (renderReleaseNote .) -}}
  {{end -}}
{{- end -}}
{{- if gt (len $breaking) 0 -}}
BREAKING CHANGES

{{range $breaking | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
{{- if gt (len $features) 0}}
FEATURES

{{range $features | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
{{- if gt (len $improvements) 0}}
IMPROVEMENTS

{{range $improvements | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
{{- if gt (len $bugs) 0}}
BUGS

{{range $bugs | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
`

const defaultBlockTypeChangelogTemplate = `
{{- $breaking := newStringList -}}
{{- $features := newStringList -}}
{{- $improvements := newStringList -}}
{{- $bugs := newStringList -}}
{{- range . -}}
  {{if eq "breaking-change" .Type -}}
	{{$breaking = append $breaking (renderReleaseNote .) -}}
  {{else if or (eq "new-resource" .Type) (eq "new-data-source" .Type) (eq "feature" .Type) -}}
	{{$features = append $features (renderReleaseNote .) -}}
  {{else if eq "improvement" .Type -}}
	{{$improvements = append $improvements (renderReleaseNote .) -}}
  {{else if eq "bug" .Type -}}
	{{$bugs = append $bugs (renderReleaseNote .) -}}
  {{end -}}
{{- end -}}
{{- if gt (len $breaking) 0 -}}
BREAKING CHANGES

{{range $breaking | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
{{- if gt (len $features) 0}}
FEATURES

{{range $features | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
{{- if gt (len $improvements) 0}}
IMPROVEMENTS

{{range $improvements | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
{{- if gt (len $bugs) 0}}
BUGS

{{range $bugs | sortAlpha -}}
* {{. }}
{{end -}}
{{- end -}}
`

func filterPrefix(prefix string, trim bool, data []string) []string {
	result := []string{}
	for _, s := range data {
		if strings.HasPrefix(s, prefix) {
			if trim {
				s = strings.TrimPrefix(s, prefix)
			}
			result = append(result, s)
		}
	}

	return result
}

func newStringList() []string {
	return make([]string, 0)
}

func renderReleaseNoteFunc(templateText string) func(ReleaseNote) (string, error) {
	return func(note ReleaseNote) (string, error) {
		return render(templateText, note, nil)
	}
}

func renderChangelog(changelogTemplateText, releaseNoteTemplateText string, notes []ReleaseNote) (string, error) {
	funcs := template.FuncMap{
		"renderReleaseNote": renderReleaseNoteFunc(releaseNoteTemplateText),
	}

	return render(changelogTemplateText, notes, funcs)
}

func render(templateText string, data interface{}, additionalFuncs template.FuncMap) (string, error) {
	const renderTemplateName = "render"

	tmpl := template.New(renderTemplateName)
	funcs := sprig.TxtFuncMap()
	for n, f := range additionalFuncs {
		funcs[n] = f
	}
	funcs["filterPrefix"] = filterPrefix
	funcs["newStringList"] = newStringList
	tmpl = tmpl.Funcs(funcs)

	var err error
	tmpl, err = tmpl.Parse(templateText)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	err = tmpl.Execute(builder, data)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
