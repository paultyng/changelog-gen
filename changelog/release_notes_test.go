package changelog

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextFromPR(t *testing.T) {
	for i, c := range []struct {
		expected []releaseNoteEntry
		title    string
		body     string
	}{
		// zero case
		{nil, "", ""},

		// text in title
		{[]releaseNoteEntry{{Text: "foo"}}, "foo", ""},

		// text in body, type in labels
		{[]releaseNoteEntry{{Text: "foo"}}, "bar", "```release-note\nfoo\n```"},
		{[]releaseNoteEntry{{Text: "foo"}}, "bar", "```releasenote\nfoo\n```"},
		{[]releaseNoteEntry{{Text: "foo"}}, "bar", "\n```releasenote\nfoo\n```\n"},

		// text in title (malformed body)
		{[]releaseNoteEntry{{Text: "bar"}}, "bar", "\n ```releasenote\nfoo\n```"},

		// text in body, type in body
		{[]releaseNoteEntry{{Type: "bug", Text: "foo"}}, "", "```release-note:bug\nfoo\n```"},
		{[]releaseNoteEntry{{Type: "enhancement", Text: "bar"}}, "", "```releasenote:enhancement\nbar\n```"},

		// text in body, type in body, multiple blocks
		{[]releaseNoteEntry{{Type: "bug", Text: "foo"}, {Type: "enhancement", Text: "bar"}},
			"", "\n```releasenote:bug\nfoo\n```\n\n```release-note:enhancement\nbar\n```\n"},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			res := c.expected
			sort.Slice(res, func(i, j int) bool {
				if res[i].Type < res[j].Type {
					return true
				} else if res[j].Type < res[i].Type {
					return false
				} else if res[i].Text < res[j].Text {
					return true
				} else if res[j].Text < res[i].Text {
					return false
				}
				return false
			})
			actual := releaseNoteBlocks(c.title, c.body)
			assert.Equal(t, res, actual)
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
		{"foo", "original author:      @foo"},
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
