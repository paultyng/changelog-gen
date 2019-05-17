package main

import (
	"fmt"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	
}

func TestParseGitWD(t *testing.T) {
	for i, c := range []struct {
		expectedOwner  string
		expectedRepo   string
		expectedBranch string
		path           string
		remote         string
	}{
		{"paultyng", "changelog-gen", "", "goodoriginhttps", "origin"},
		{"paultyng", "changelog-gen", "master", "goodoriginssh", "origin"},
		{"", "", "", "goodoriginhttps", "notfound"},
		{"", "", "", "badorigin", "origin"},
	} {
		t.Run(fmt.Sprintf("%d %s/%s", i, c.expectedOwner, c.expectedRepo), func(t *testing.T) {
			actualOwner, actualRepo, actualBranch, err := parseGitWD(path.Join("./testdata", c.path), c.remote)
			assert.NoError(t, err)
			assert.Equal(t, c.expectedOwner, actualOwner)
			assert.Equal(t, c.expectedRepo, actualRepo)
			assert.Equal(t, c.expectedBranch, actualBranch)
		})
	}
}
