package main

import (
	"net/url"
	"path"
	"regexp"
	"strings"
)

// from https://github.com/src-d/go-git/blob/2ab6d5cd72b59cfd36b08078ddeebd1efb0d2254/internal/url/url.go
// but it is an internal package, so cannot reference it directly, so copy and paste
var (
	isSchemeRegExp   = regexp.MustCompile(`^[^:]+://`)
	scpLikeURLRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})/)?(?P<path>[^\\].*)$`)
)

func matchesScpLike(url string) bool {
	return !isSchemeRegExp.MatchString(url) && scpLikeURLRegExp.MatchString(url)
}

func findScpLikeComponents(url string) (user, host, port, path string) {
	m := scpLikeURLRegExp.FindStringSubmatch(url)
	return m[1], m[2], m[3], m[4]
}

func parseGitWD(gitPath, remoteName string) (owner string, repo string, branch string, err error) {
	gitRepo, err := git.PlainOpen(gitPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	head, err := gitRepo.Head()
	if err != nil {
		if err != plumbing.ErrReferenceNotFound {
			return "", "", "", err
		}
	} else {
		if n := head.Name(); n.IsBranch() {
			branch = n.Short()
		}
	}

	remote, err := gitRepo.Remote(remoteName)
	if err != nil {
		if err == git.ErrRemoteNotFound {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	for _, raw := range remote.Config().URLs {
		if matchesScpLike(raw) {
			_, host, _, remotePath := findScpLikeComponents(raw)
			if host != "github.com" {
				continue
			}

			owner, repo = path.Split(remotePath)
			break
		}

		u, err := url.Parse(raw)
		if err != nil {
			return "", "", "", err
		}

		if u.Hostname() != "github.com" {
			continue
		}

		owner, repo = path.Split(u.Path)
		break
	}

	owner = strings.Trim(owner, "/")
	repo = strings.TrimRight(repo, ".git")

	return
}
