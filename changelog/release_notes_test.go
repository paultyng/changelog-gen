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
		// zero case
		{"", "", ""},

		// text in title
		{"foo", "foo", ""},

		// text in body
		{"foo", "bar", "```release-note\nfoo\n```"},
		{"foo", "bar", "```releasenote\nfoo\n```"},
		{"foo", "bar", "\n```releasenote\nfoo\n```\n"},

		// text in title (malformed body)
		{"bar", "bar", "\n ```releasenote\nfoo\n```"},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			actual := textFromPR(c.title, c.body)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestAuthorFromPR(t *testing.T) {
	for i, c := range []struct {
		expected string
		body     string
	}{
		// zero case
		{"", ""},

		{"foo", "original author: @foo"},
		{"foo", "original author:@foo"},
		{"foo", "original author:     @foo"},
		{"foo", "Original Author: @foo"},
		{"foo", "**Original Author:** @foo"},
		{"foo", "\n**Original Author:** @foo\n"},

		{"", "\n **Original Author:** @foo\n"},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			actual, actualURL, ok := authorFromPR(c.body)
			assert.Equal(t, c.expected != "", ok)
			if ok {
				assert.Equal(t, c.expected, actual)
				// TODO: confirm URL encoding appropriately?
				assert.Equal(t, fmt.Sprintf("https://github.com/%s", c.expected), actualURL)
			}
		})
	}
}
