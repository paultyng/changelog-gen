package changelog

import (
	"context"
	"sort"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
)

func BuildChangelog(
	ctx context.Context,
	client *githubv4.Client,
	logger hclog.Logger,
	changelogTemplate,
	releaseNoteTemplate,
	owner,
	repo,
	branch,
	start,
	end string,
) (string, error) {
	prIDs, err := listPullRequestIDs(ctx, client, logger, owner, repo, branch, start, end)
	if err != nil {
		return "", err
	}

	logger.Info("found PRs", "count", len(prIDs))

	notes, err := pullRequestsToReleaseNotes(ctx, client, logger, prIDs)
	if err != nil {
		return "", err
	}

	sort.Slice(notes, func(i int, j int) bool {
		return notes[i].PRDate.After(notes[j].PRDate)
	})

	if changelogTemplate == "" {
		changelogTemplate = defaultChangelogTemplate
	}

	if releaseNoteTemplate == "" {
		releaseNoteTemplate = defaultReleaseNoteTemplate
	}

	return renderChangelog(changelogTemplate, releaseNoteTemplate, notes)
}

func filterReleaseNotes(notes []ReleaseNote, filter func(ReleaseNote) bool) []ReleaseNote {
	if filter == nil {
		filtered := notes
		notes = []ReleaseNote{}
		return filtered
	}

	filtered := []ReleaseNote{}
	unfiltered := []ReleaseNote{}
	for _, n := range notes {
		if filter(n) {
			filtered = append(filtered, n)
			continue
		}
		unfiltered = append(unfiltered, n)
	}

	notes = unfiltered

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PRDate.After(filtered[j].PRDate)
	})
	return filtered
}
