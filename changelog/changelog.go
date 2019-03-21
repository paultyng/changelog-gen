package changelog

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
)

const (
	labelPrefixService = "service/"

	labelBug = "bug"

	labelNewResource   = "new-resource"
	labelNewDataSource = "new-data-source"

	labelBreakingChange = "breaking-change"
)

// ReleaseNote is the type that represents the total sum of all the information
// we've gathered about a single release note.
type ReleaseNote struct {
	// Text is the actual content of the release note
	Text string `json:"text"`

	// Markdown is the markdown formatted note for the changelog
	Markdown string `json:"markdown"`

	// Author is the GitHub username of the commit author
	Author string `json:"author"`

	// AuthorUrl is the GitHub URL of the commit author
	AuthorURL string `json:"author_url"`

	//PRDate is the Date the PR was merged
	PRDate time.Time `json:"pr_date"`

	// PrUrl is a URL to the PR
	PRURL string `json:"pr_url"`

	// PrNumber is the number of the PR
	PRNumber int `json:"pr_number"`

	// Services is a list of the labels beginning with `service/*`.
	Services []string `json:"areas,omitempty"`

	// Indicates whether or not a note will appear as a bug (`bug` label).
	Bug bool `json:"bug,omitempty"`

	// NewResource and NewDataSource indicate if the PR includes new resources
	// and/or data sources.
	NewResource   bool `json:"new_resource,omitempty"`
	NewDataSource bool `json:"new_data_source,omitempty"`

	// BreakingChange indicates if this change was breaking (the
	// `breaking-change` label was applied to the PR).
	BreakingChange bool `json:"breaking_change,omitempty"`
}

func listPullRequestIDs(
	ctx context.Context,
	client *githubv4.Client,
	logger hclog.Logger,
	owner,
	repo,
	branch,
	start,
	end string,
) ([]string, error) {
	var q struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commit struct {
						History struct {
							Nodes []struct {
								OID string

								AssociatedPullRequests struct {
									Nodes []struct {
										BaseRef struct {
											Repository struct {
												Owner struct {
													Login string
												}
												Name string
											}

											Name string
										}
										State  githubv4.PullRequestState
										ID     string
										Number int
									}
								} `graphql:"associatedPullRequests(first: 100)"`
							}
						} `graphql:"history(since: $since, until: $until)"`
					} `graphql:"... on Commit"`
				}
			} `graphql:"ref(qualifiedName: $ref)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}

	prNodeIDs := map[string]bool{}

	since, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, err
	}

	until, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, err
	}

	logger = logger.With("since", since, "until", until)

	logger.Info("checking commits for associated PRs")
	err = client.Query(ctx, &q, map[string]interface{}{
		"repoOwner": githubv4.String(owner),
		"repoName":  githubv4.String(repo),
		"ref":       githubv4.String(fmt.Sprintf("refs/heads/%s", branch)),
		"since":     githubv4.GitTimestamp{Time: since},
		"until":     githubv4.GitTimestamp{Time: until},
	})
	if err != nil {
		return nil, err
	}

	for _, hn := range q.Repository.Ref.Target.Commit.History.Nodes {
		logger := logger.With("commit", hn.OID)
		logger.Debug("checking commit PRs")

		if len(hn.AssociatedPullRequests.Nodes) == 100 {
			// this is weird, should probably log it
		}
		for _, prn := range hn.AssociatedPullRequests.Nodes {
			logger := logger.With("pr", prn.Number)

			if prn.BaseRef.Name != branch ||
				prn.BaseRef.Repository.Name != repo ||
				prn.BaseRef.Repository.Owner.Login != owner {
				logger.Debug("external PR, skipping")
				continue
			}

			if prn.State != githubv4.PullRequestStateMerged {
				logger.Debug("unmerged PR, skipping")
				continue
			}
			// TODO: check base ref on PR to make sure its master?
			prNodeIDs[prn.ID] = true
		}
	}

	prIDs := make([]string, 0, len(prNodeIDs))
	for id := range prNodeIDs {
		prIDs = append(prIDs, id)
	}

	return prIDs, nil
}

func pullRequestsToReleaseNotes(
	ctx context.Context,
	client *githubv4.Client,
	logger hclog.Logger,
	prIDs []string,
) ([]ReleaseNote, error) {
	var q struct {
		Nodes []struct {
			PullRequest struct {
				MergedAt time.Time
				ID       string
				Number   int
				Title    string
				Body     string
				URL      string
				Author   struct {
					Login string
					URL   string
				}
				Labels struct {
					Nodes []struct {
						Name string
					}
				} `graphql:"labels(first: 100)"`
			} `graphql:"... on PullRequest"`
		} `graphql:"nodes(ids: $ids)"`
	}

	logger.Info("retrieving PRs to build release notes")
	err := client.Query(ctx, &q, map[string]interface{}{
		"ids": prIDs,
	})
	if err != nil {
		return nil, err
	}

	notes := make([]ReleaseNote, 0, len(q.Nodes))
	for _, n := range q.Nodes {
		logger := logger.With("pr", n.PullRequest.Number, "prid", n.PullRequest.ID)

		logger.Info("building release note")

		note := ReleaseNote{
			PRDate:    n.PullRequest.MergedAt,
			PRNumber:  n.PullRequest.Number,
			PRURL:     n.PullRequest.URL,
			Author:    n.PullRequest.Author.Login,
			AuthorURL: n.PullRequest.Author.URL,
			Text:      textFromPR(n.PullRequest.Title, n.PullRequest.Body),
		}

		for _, ln := range n.PullRequest.Labels.Nodes {
			switch {
			case strings.HasPrefix(ln.Name, labelPrefixService):
				note.Services = append(note.Services, strings.TrimPrefix(ln.Name, labelPrefixService))
			case ln.Name == labelBug:
				note.Bug = true
			case ln.Name == labelBreakingChange:
				note.BreakingChange = true
			case ln.Name == labelNewResource:
				note.NewResource = true
			case ln.Name == labelNewDataSource:
				note.NewDataSource = true
			}
		}

		servicesMarkdown := ""
		if len(note.Services) > 0 {
			servicesMarkdown = "**" + strings.Join(note.Services, ", ") + ":** "
		}
		note.Markdown = fmt.Sprintf("%s%s ([%d](%s) by [%s](%s))", servicesMarkdown, note.Text, note.PRNumber, note.PRURL, note.Author, note.AuthorURL)

		notes = append(notes, note)
	}

	return notes, nil
}

func textFromPR(title, body string) string {
	// TODO: add body parsing
	return title
}

func BuildChangelog(
	ctx context.Context,
	client *githubv4.Client,
	logger hclog.Logger,
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

	markdown := strings.Builder{}

	for _, group := range []struct {
		heading string
		filter  func(ReleaseNote) bool
	}{
		{"BREAKING CHANGES", func(rn ReleaseNote) bool {
			return rn.BreakingChange
		}},
		{"FEATURES", func(rn ReleaseNote) bool {
			return rn.NewDataSource || rn.NewResource
		}},
		{"ENHANCEMENTS", func(rn ReleaseNote) bool {
			return !rn.Bug
		}},
		{"BUGS", nil},
	} {
		filtered := filterReleaseNotes(notes, group.filter)
		if len(filtered) == 0 {
			continue
		}
		markdown.WriteString(group.heading)
		markdown.WriteString("\n\n")
		for _, note := range filtered {
			markdown.WriteString("* ")
			markdown.WriteString(note.Markdown)
			markdown.WriteString("\n")
		}
		markdown.WriteString("\n")
	}

	return markdown.String(), nil
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
