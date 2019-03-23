package main // import "github.com/paultyng/changelog-gen"

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/paultyng/changelog-gen/changelog"
)

type options struct {
	githubToken string
	owner       string
	repo        string
	branch      string
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
		branch:      *flBranch,
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
			// TODO: load override templates from files somewhere
			"", "",
			opts.owner, opts.repo, opts.branch,
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
