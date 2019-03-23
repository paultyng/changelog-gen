package changelog

import (
	"context"
	"errors"
	"sort"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
)

func TimeFromCommit(
	ctx context.Context,
	client *githubv4.Client,
	owner, repo, commit string,
) (time.Time, error) {
	var q struct {
		Repository struct {
			Object *struct {
				Commit struct {
					CommittedDate time.Time
				} `graphql:"... on Commit"`
			} `graphql:"object(expression: $commit)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}

	err := client.Query(ctx, &q, map[string]interface{}{
		"repoOwner": githubv4.String(owner),
		"repoName":  githubv4.String(repo),
		"commit":    githubv4.String(commit),
	})
	if err != nil {
		return time.Time{}, err
	}
	if q.Repository.Object == nil {
		return time.Time{}, errors.New("unable to find commit")
	}
	return q.Repository.Object.Commit.CommittedDate, nil
}

func BuildChangelog(
	ctx context.Context,
	client *githubv4.Client,
	logger hclog.Logger,
	changelogTemplate,
	releaseNoteTemplate,
	owner, repo, branch string,
	start, end time.Time,
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
