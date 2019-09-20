package changelog

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
)

const (
	labelBreakingChange = "breaking-change"
)

var (
	labelsBug           = []string{"bug", "kind/bug"}
	labelsNoReleaseNote = []string{"no-release-note", "release-note-none"}
)

// ReleaseNote is the type that represents the total sum of all the information
// we've gathered about a single release note.
type ReleaseNote struct {
	// Text is the actual content of the release note
	Text string `json:"text"`

	// Author is the GitHub username of the commit author
	Author string `json:"author"`

	// AuthorURL is the GitHub URL of the commit author
	AuthorURL string `json:"author_url"`

	//PRDate is the Date the PR was merged
	PRDate time.Time `json:"pr_date"`

	// PRUrl is a URL to the PR
	PRURL string `json:"pr_url"`

	// PRNumber is the number of the PR
	PRNumber int `json:"pr_number"`

	// Labels is a list of all the labels on the PR.
	Labels []string `labels:"areas,omitempty"`

	// Indicates whether or not a note will appear as a bug (`bug` label).
	Bug bool `json:"bug,omitempty"`

	// BreakingChange indicates if this change was breaking (the
	// `breaking-change` label was applied to the PR).
	BreakingChange bool `json:"breaking_change,omitempty"`

	// Type is the type of entry the ReleaseNote is
	Type string
}

type releaseNoteEntry struct {
	Type string
	Text string
}

func listPullRequestIDs(
	ctx context.Context,
	client *githubv4.Client,
	logger hclog.Logger,
	owner, repo, branch string,
	start, end time.Time,
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
										Labels struct {
											Nodes []struct {
												Name string
											}
										} `graphql:"labels(first: 100)"`
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

	logger = logger.With("since", start, "until", end)

	logger.Info("checking commits for associated PRs")
	err := client.Query(ctx, &q, map[string]interface{}{
		"repoOwner": githubv4.String(owner),
		"repoName":  githubv4.String(repo),
		"ref":       githubv4.String(fmt.Sprintf("refs/heads/%s", branch)),
		"since":     githubv4.GitTimestamp{Time: start},
		"until":     githubv4.GitTimestamp{Time: end},
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

			noChangelog := false
			for _, ln := range prn.Labels.Nodes {
				for _, nrn := range labelsNoReleaseNote {
					if ln.Name == nrn {
						noChangelog = true
						break
					}
				}
			}
			if noChangelog {
				logger.Debug("no-changelog label applied, skipping")
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

func stringInSlice(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
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

		author, authorURL, found := authorFromPR(n.PullRequest.Body)
		if !found {
			author = n.PullRequest.Author.Login
			authorURL = n.PullRequest.Author.URL
		}

		note := ReleaseNote{
			PRDate:    n.PullRequest.MergedAt,
			PRNumber:  n.PullRequest.Number,
			PRURL:     n.PullRequest.URL,
			Author:    author,
			AuthorURL: authorURL,
		}

		for _, ln := range n.PullRequest.Labels.Nodes {
			switch {
			case stringInSlice(labelsBug, ln.Name):
				note.Bug = true
			case ln.Name == labelBreakingChange:
				note.BreakingChange = true
			default:
				note.Labels = append(note.Labels, ln.Name)
			}
		}

		for _, entry := range releaseNoteBlocks(n.PullRequest.Title, n.PullRequest.Body) {
			n := note
			n.Text = entry.Text
			n.Type = entry.Type
			notes = append(notes, n)
		}
	}

	return notes, nil
}

var textInBodyREs = []*regexp.Regexp{
	regexp.MustCompile("(?m)^```release-note\n(?P<note>.+)\n```"),
	regexp.MustCompile("(?m)^```releasenote\n(?P<note>.+)\n```"),
	regexp.MustCompile("(?m)^```release-note:(?P<type>[^\n]*)\n(?P<note>.+)\n```"),
	regexp.MustCompile("(?m)^```releasenote:(?P<type>[^\n]*)\n(?P<note>.+)\n```"),
}

func releaseNoteBlocks(title, body string) []releaseNoteEntry {
	var res []releaseNoteEntry
	for _, re := range textInBodyREs {
		matches := re.FindAllStringSubmatch(body, -1)
		if len(matches) == 0 {
			continue
		}

		for _, match := range matches {
			note := ""
			typ := ""
			for i, name := range re.SubexpNames() {
				switch name {
				case "note":
					note = match[i]
				case "type":
					typ = match[i]
				}
				if note != "" && typ != "" {
					break
				}
			}

			note = strings.TrimRight(note, "\r")
			note = stripMarkdownBullet(note)

			note = strings.TrimSpace(note)
			typ = strings.TrimSpace(typ)

			if note == "" {
				continue
			}

			res = append(res, releaseNoteEntry{
				Type: typ,
				Text: note,
			})
		}
	}
	if len(res) < 1 && title != "" {
		res = append(res, releaseNoteEntry{
			Text: title,
		})
	}
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
	return res
}

func stripMarkdownBullet(note string) string {
	re := regexp.MustCompile(`(?i)\*\s`)
	return re.ReplaceAllString(note, "")
}

var authorInBodyREs = []*regexp.Regexp{
	// /cc syntax is too ambiguous probably
	// regexp.MustCompile("(?m)^/[Cc][Cc] *@(?P<login>.+)"),
	regexp.MustCompile("(?m)^(\\*\\*)?[Oo]riginal [Aa]uthor:(\\*\\*)? *@(?P<login>.+)"),
}

func authorFromPR(body string) (string, string, bool) {
	for _, re := range authorInBodyREs {
		match := re.FindStringSubmatch(body)
		if len(match) == 0 {
			continue
		}

		author := ""
		for i, name := range re.SubexpNames() {
			if name == "login" {
				author = match[i]
				break
			}
		}

		author = strings.TrimLeft(author, "@")

		if author != "" {
			authorURL := fmt.Sprintf("https://github.com/%s", author)
			return author, authorURL, true
		}
	}

	return "", "", false
}
