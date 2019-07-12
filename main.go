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
	owner       string
	repo        string

	// optional
	branch              string
	changelogTemplate   string
	releaseNoteTemplate string
	labelPrefix         string
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
			envString("GITHUB_BRANCH", "master"),
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

		flLabelPrefix = flagset.String(
			"label-prefix",
			"",
			"Prefix for filtering changelog-related labels; text after prefix is used to determine category",
		)
	)

	if err := flagset.Parse(args); err != nil {
		return nil, nil, err
	}

	if *flGitHubToken == "" {
		return nil, nil, errors.New("GitHub token must be set via -github-token or $GITHUB_TOKEN")
	}

	if *flOwner == "" {
		return nil, nil, errors.New("GitHub repository owner must be set via -owner or $GITHUB_OWNER")
	}

	if *flRepo == "" {
		return nil, nil, errors.New("GitHub repository must be set via -repo or $GITHUB_REPO")
	}

	return flagset.Args(), &options{
		githubToken: *flGitHubToken,
		owner:       *flOwner,
		repo:        *flRepo,

		branch:              *flBranch,
		changelogTemplate:   *flChangelogTemplate,
		releaseNoteTemplate: *flReleaseNoteTemplate,
		labelPrefix:         *flLabelPrefix,
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
			return errors.New("2 arguments are required")
		}

		branch := opts.branch
		if branch == "" {
			branch = "master"
		}

		changelogTemplate, err := loadTemplate(opts.changelogTemplate)
		if err != nil {
			return err
		}

		releaseNoteTemplate, err := loadTemplate(opts.releaseNoteTemplate)
		if err != nil {
			return err
		}

		ctx := context.Background()
		client := githubClient(ctx, opts.githubToken)

		startCommit, startTime, err := parseCommitOrTime(args[0])
		if err != nil {
			return err
		}
		if startCommit != "" {
			startTime, err = changelog.TimeFromCommit(ctx, client, opts.owner, opts.repo, startCommit)
			if err != nil {
				return err
			}
		}

		endCommit, endTime, err := parseCommitOrTime(args[1])
		if err != nil {
			return err
		}
		if endCommit != "" {
			endTime, err = changelog.TimeFromCommit(ctx, client, opts.owner, opts.repo, endCommit)
			if err != nil {
				return err
			}
		}

		cl, err := changelog.BuildChangelog(
			ctx, client, logger,
			changelogTemplate,
			releaseNoteTemplate,
			opts.owner,
			opts.repo,
			branch,
			opts.labelPrefix,
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
