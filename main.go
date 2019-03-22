package main // import "github.com/paultyng/changelog-gen"

import (
	"context"
	"fmt"
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/paultyng/changelog-gen/changelog"
)

func main() {
	const (
		// to get dates in proper format:
		// git show -q --pretty='format:%cI' 441ec74e66706cc0a75d4d207724cd6460f5f6a4 f3bdfeaaa7ddd9522c549e9f134948e0d698ab02
		releaseStart1_59_0 = "2019-02-08T07:42:58+00:00"
		releaseEnd1_59_0   = "2019-02-14T22:23:11+00:00"
	)

	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Info,
		Output: os.Stderr,
	})
	client := githubClient(ctx)
	org := "terraform-providers"
	repo := "terraform-provider-aws"
	branch := "master"

	cl, err := changelog.BuildChangelog(ctx, client, logger, "", "", org, repo, branch, releaseStart1_59_0, releaseEnd1_59_0)
	if err != nil {
		panic(err)
	}

	fmt.Println(cl)
}

func githubClient(ctx context.Context) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	return client
}
