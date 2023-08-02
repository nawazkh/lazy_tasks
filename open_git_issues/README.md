# open_issues 

`open_issues` is a go utility intended to be used to open git issues on public repos.

## Pre-requisites

- Create a github token which has a access to the following scopes.
  - `repo:status` - Grants access to commit status on public and private repositories.
  - `repo_deployment` - Grants access to deployment statuses on public and private repositories.
  - `public_repo` - Grants access to public repositories

- [repos.txt](./repos.txt) is populated with list of public(or private) repositories
  - If the issue is to be created against `https://github.com/my-org/public-repo`, _repos.txt_ should contain `my-org/public-repo`.
  - Add a new repo at a new line in the file.

- decide upon the title of the issue to be opened.

- [issue_body.md](./issue_body.md) should have description of the issue in a markdown format.

## How to run the tool

- build the utility.
  ```
  go build open_issues.go
  ```
- run the utility
  ```
  ./open_issues --token $GITHUB_ISSUE_OPENER_TOKEN --file repos.txt --title "This is a test issue" --body issue_body.md
  ```

## Sample output

```
❯ ./open_issues --token $GITHUB_ISSUE_OPENER_TOKEN --file repos.txt --title "This is a test issue" --body issue_body.md
request will be sent to https://api.github.com/repos/nawazkh/cluster-api-provider-azure/issues
Failed to create issue for repository 'nawazkh/cluster-api-provider-azure'. Status code: 410, Response: Status:410 Gone


request will be sent to https://api.github.com/repos/nawazkh/lazy_tasks/issues
Issue created for repository 'nawazkh/lazy_tasks'. URL: https://github.com/nawazkh/lazy_tasks/issues/5
```
