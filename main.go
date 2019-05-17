package main // import "github.com/paultyng/changelog-gen"

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/paultyng/changelog-gen/changelog"
)

type options struct {
	// required
	githubToken string

	// optional
	owner               string
	repo                string
	branch              string
	changelogTemplate   string
	releaseNoteTemplate string
}

func envString(key, def string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return def
}

func parseOptions(args []string) ([]string, *options, error) {
	flagset := flag.NewFlagSet("changelog-gen", flag.ExitOnError)

	var (
		flGitHubToken = flagset.String(
			"github-token",
			envString("GITHUB_TOKEN", ""),
			"A personal GitHub access token (required)",
		)

		flOwner = flagset.String(
			"owner",
			envString("GITHUB_OWNER", ""),
			"GitHub repository owner",
		)

		flRepo = flagset.String(
			"repo",
			envString("GITHUB_REPO", ""),
			"Github repository name",
		)

		flBranch = flagset.String(
			"branch",
			envString("GITHUB_BRANCH", ""),
			"Github branch (defaults to master)",
		)

		flChangelogTemplate = flagset.String(
			"changelog",
			"",
			"Changelog template path (leave blank for built-in template)",
		)

		flReleaseNoteTemplate = flagset.String(
			"releasenote",
			"",
			"Release note template path (leave blank for built-in template)",
		)
	)

	if err := flagset.Parse(args); err != nil {
		return nil, nil, err
	}

	if *flGitHubToken == "" {
		return nil, nil, errors.New("GitHub token must be set via -github-token or $GITHUB_TOKEN")
	}

	return flagset.Args(), &options{
		githubToken: *flGitHubToken,

		owner:               *flOwner,
		repo:                *flRepo,
		branch:              *flBranch,
		changelogTemplate:   *flChangelogTemplate,
		releaseNoteTemplate: *flReleaseNoteTemplate,
	}, nil
}

var commitRE = regexp.MustCompile("^[0-9a-f]{5,40}$")

func parseCommitOrTime(v string) (string, time.Time, error) {
	if commitRE.MatchString(v) {
		return v, time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return "", time.Time{}, err
	}
	return "", t, nil
}

func loadTemplate(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		// TODO: load level from somewhere else?
		Level:  hclog.Info,
		Output: os.Stderr,
	})
	err := func() error {

		args, opts, err := parseOptions(os.Args[1:])
		if err != nil {
			return err
		}
		if len(args) != 2 {
			// TODO: use most recent tag if nothing passed
			// ie. git describe --tags --abbrev=0
			// like goreleaser, for current commit
			return errors.New("2 arguments are required")
		}

		wd, err := os.Getwd()
		if err != nil {
			// TODO: ignore this and skip the git stuff?
			return err
		}

		gitOwner, gitRepo, gitBranch, err := parseGitWD(wd, "origin")
		if err != nil {
			// TODO: ignore this?
			return err
		}

		owner := coalesce(opts.owner, gitOwner)
		repo := coalesce(opts.repo, gitRepo)
		branch := coalesce(opts.branch, gitBranch, "master")

		if owner == "" || repo == "" {
			return errors.New("you must provide an owner and repository or run the tool inside a git repo with the proper origin set")
		}

		changelogTemplate, err := loadTemplate(opts.changelogTemplate)
		if err != nil {
			return err
		}

		releaseNoteTemplate, err := loadTemplate(opts.releaseNoteTemplate)
		if err != nil {
			return err
		}

		logger.Info("parsing commits", "owner", owner, "repo", repo, "branch", branch)

		ctx := context.Background()
		client := githubClient(ctx, opts.githubToken)

		startCommit, startTime, err := parseCommitOrTime(args[0])
		if err != nil {
			return err
		}
		if startCommit != "" {
			startTime, err = changelog.TimeFromCommit(ctx, client, owner, repo, startCommit)
			if err != nil {
				return err
			}
		}

		endCommit, endTime, err := parseCommitOrTime(args[1])
		if err != nil {
			return err
		}
		if endCommit != "" {
			endTime, err = changelog.TimeFromCommit(ctx, client, owner, repo, endCommit)
			if err != nil {
				return err
			}
		}

		cl, err := changelog.BuildChangelog(
			ctx, client, logger,
			changelogTemplate, releaseNoteTemplate,
			owner, repo, branch,
			startTime, endTime,
		)
		if err != nil {
			panic(err)
		}

		fmt.Println(cl)
		return nil
	}()
	if err != nil {
		logger.Error("error parsing options", "err", err)
		os.Exit(1)
	}
}

func githubClient(ctx context.Context, token string) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	return client
}
