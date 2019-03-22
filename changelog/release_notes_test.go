package changelog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextFromPR(t *testing.T) {
	for i, c := range []struct {
		expected string
		title    string
		body     string
	}{
		{"", "", ""},
		{"foo", "foo", ""},
		{"foo", "bar", "```release-note\nfoo\n```"},
		{"foo", "bar", "```releasenote\nfoo\n```"},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			actual := textFromPR(c.title, c.body)
			assert.Equal(t, c.expected, actual)
		})
	}
}
