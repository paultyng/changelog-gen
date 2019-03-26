package changelog

import (
	"html/template"
	"strings"

	"github.com/Masterminds/sprig"
)

const defaultReleaseNoteTemplate = `{{with .Labels | filterPrefix "service/" true | sortAlpha }}{{if len . | lt 0 }}**{{. | join ", " }}:** {{end}}{{end}}{{.Text | raw}} ([{{.PRNumber}}]({{.PRURL}}) by [{{.Author}}]({{.AuthorURL}}))`
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
* {{. | raw}}
{{end -}}
{{- end -}}
{{- if gt (len $features) 0}}
FEATURES

{{range $features | sortAlpha -}}
* {{. | raw}}
{{end -}}
{{- end -}}
{{- if gt (len $improvements) 0}}
IMPROVEMENTS

{{range $improvements | sortAlpha -}}
* {{. | raw}}
{{end -}}
{{- end -}}
{{- if gt (len $bugs) 0}}
BUGS

{{range $bugs | sortAlpha -}}
* {{. | raw}}
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

func renderReleaseNoteFunc(templateText string) func(ReleaseNote) (template.HTML, error) {
	return func(note ReleaseNote) (template.HTML, error) {
		raw, err := render(templateText, note, nil)
		if err != nil {
			return "", err
		}
		return template.HTML(raw), nil
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
	funcs := sprig.FuncMap()
	for n, f := range additionalFuncs {
		funcs[n] = f
	}
	funcs["filterPrefix"] = filterPrefix
	funcs["newStringList"] = newStringList
	funcs["raw"] = func(s string) template.HTML {
		return template.HTML(s)
	}

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
