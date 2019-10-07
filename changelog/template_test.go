package changelog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestRender_defaultChangelogTemplate(t *testing.T) {
	expected := `BREAKING CHANGES

* this is a breaking bug ([0]() by []())
* this is a breaking feature ([0]() by []())

FEATURES

* this is a new data-source ([0]() by []())
* this is a new resource ([0]() by []())

IMPROVEMENTS

* this is an improvement & 'stuff' ([0]() by []())

BUGS

* this is a bug ([0]() by []())
`

	actual, err := renderChangelog(defaultChangelogTemplate, defaultReleaseNoteTemplate, []ReleaseNote{
		{
			BreakingChange: true,
			Text:           "this is a breaking feature",
		},
		{
			BreakingChange: true,
			Bug:            true,
			Text:           "this is a breaking bug",
		},
		{
			Labels: []string{"new-resource"},
			Text:   "this is a new resource",
		},
		{
			Labels: []string{"new-data-source"},
			Text:   "this is a new data-source",
		},
		{
			Text: "this is an improvement & 'stuff'",
		},
		{
			Bug:  true,
			Text: "this is a bug",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestRender_defaultBlockTypeChangelogTemplate(t *testing.T) {
	expected := `BREAKING CHANGES

* this is a breaking change ([0]() by []())

FEATURES

* this is a new data-source ([0]() by []())
* this is a new resource ([0]() by []())

IMPROVEMENTS

* this is an improvement & 'stuff' ([0]() by []())

BUGS

* this is a bug ([0]() by []())
`

	actual, err := renderChangelog(defaultBlockTypeChangelogTemplate, defaultReleaseNoteTemplate, []ReleaseNote{
		{
			Type: "breaking-change",
			Text: "this is a breaking change",
		},
		{
			Type: "new-resource",
			Text: "this is a new resource",
		},
		{
			Type: "new-data-source",
			Text: "this is a new data-source",
		},
		{
			Type: "improvement",
			Text: "this is an improvement & 'stuff'",
		},
		{
			Type: "bug",
			Text: "this is a bug",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestRender_defaultReleaseNoteTemplate(t *testing.T) {
	for i, c := range []struct {
		expected string
		rn       ReleaseNote
	}{
		{"bar ([2](baz) by [qux](quux))", ReleaseNote{
			Text:      "bar",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
		{"bar 'apos' ([2](baz) by [qux](quux))", ReleaseNote{
			Text:      "bar 'apos'",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
		{"bar ([2](baz) by [qux](quux))", ReleaseNote{
			Text:      "bar",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
		{"**foo:** bar ([2](baz) by [qux](quux))", ReleaseNote{
			Labels:    []string{"service/foo"},
			Text:      "bar",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
		{"**bar, foo:** bar ([2](baz) by [qux](quux))", ReleaseNote{
			Labels:    []string{"service/bar", "service/foo"},
			Text:      "bar",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
		{"**bar, foo:** bar ([2](baz) by [qux](quux))", ReleaseNote{
			Labels:    []string{"service/foo", "service/bar"},
			Text:      "bar",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
		{"bar ([2](baz) by [qux](quux))", ReleaseNote{
			Labels:    []string{"foo", "bar"},
			Text:      "bar",
			PRNumber:  2,
			PRURL:     "baz",
			Author:    "qux",
			AuthorURL: "quux",
		}},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := renderReleaseNoteFunc(defaultReleaseNoteTemplate)(c.rn)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}
