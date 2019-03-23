# CHANGELOG.md Generator

Heavily inspired by Kubernetes' [release note generator](https://github.com/kubernetes/release/tree/master/cmd/release-notes).

## Usage

```shell
$ export GITHUB_TOKEN=<your token>
$ changelog-gen \
  -owner terraform-providers \
  -repo terraform-provider-aws \
  441ec74e66706cc0a75d4d207724cd6460f5f6a4 \
  f3bdfeaaa7ddd9522c549e9f134948e0d698ab02
```
([example output](https://gist.github.com/paultyng/373566df174f8febfbff1f7cf21a238c#file-changelog-md))

The following flags are supported:

* **-github-token** GitHub token, environment variable: `GITHUB_TOKEN`
* **-owner** repository owner, environment variable: `GITHUB_OWNER`
* **-repo** repository name, environment variable: `GITHUB_NAME`
* **-branch** branch, defaults to `master`, environment variable: `GITHUB_BRANCH`

In addition to flags you must also supply either 2 commit shas or 2 RFC3339 timestamps indicating the portion of the commit log to pull PRs for.

## How Entries are Created

Each commit within the supplied range is has its associated PRs queried. Those PRs are check to find any who were merged with the base ref targetting the branch supplied in flags (in case PRs have been opened and closed on the same commit, or the commit was also part of a PR on a fork). PRs with the label `no-release-note` are also excluded.

The unique list of PRs is then converted to release notes. For each PR, the body is checked for a section with the release note copy:

    ```release-note
    This is an example release note!
    ```

If no release note is found in the body, the text is taken from the PR title.

Additionally, the PR body is checked for an override author, this can used when a bot creates PRs to indicate the original author:

    Original Author: @paultyng

If no author information is found, it defaults to the PR author.

## Templating

Only the internal template is supported currently.
