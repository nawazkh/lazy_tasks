# Rebase current branch with remote

- This folder provides a script named `rebase_with_upstream.go` and generates a binary `./bin/rebaseAllDirs`.
- Running `rebaseAllDirs` at the root of all the git repos will rebase the locally checked out branch with `<remote>`/`<active development>`
  - `<remote>` : script will figure out if the remote is `upstream` or `origin`
  - `<active development>`: script will figure out if this is `master` or `main`

## Sample Output

```shell
❯ rebaseAll
INFO 2023/05/25 10:31:32 Root Dir: "/Users/nawazhussain/msftcode"
INFO 2023/05/25 10:31:33 Successfully rebased and updated goss's local/master with upstream/master
INFO 2023/05/25 10:31:33 Successfully rebased and updated image-builder's local/error_out_cluster_create with upstream/master
INFO 2023/05/25 10:31:33 Successfully rebased and updated org's local/main with upstream/main
INFO 2023/05/25 10:31:33 Successfully rebased and updated k8s.io's local/add_nawazkh_to_capz_team with upstream/main
INFO 2023/05/25 10:31:33 Successfully rebased and updated windows-testing's local/master with upstream/master
INFO 2023/05/25 10:31:33 Successfully rebased and updated test-infra's local/add_area_labels_capi with upstream/master
INFO 2023/05/25 10:31:33 Successfully rebased and updated ip-masq-agent-v2's local/master with upstream/master
INFO 2023/05/25 10:31:33 Failed to execute git symbolic-ref on the dir /Users/nawazhussain/msftcode/azure-container-networking. Check dir manually.
INFO 2023/05/25 10:31:33 Successfully rebased and updated aks-engine's local/test with upstream/master
INFO 2023/05/25 10:31:33 Successfully rebased and updated cluster-api-provider-azure's local/main with upstream/main
INFO 2023/05/25 10:31:33 Successfully rebased and updated cluster-api's local/add_area_labels with upstream/main
INFO 2023/05/25 10:31:34 Successfully rebased and updated lazy_tasks's local/main with origin/main
INFO 2023/05/25 10:31:35 Successfully rebased and updated blinds-manager's local/main with origin/main
```

## Needed Variables

- Export `UPSTREAM` to the root of the dir containing all the git repositories
- In Makefile
  - update the `#VARIABLES` section for local testing
- Set `parent` alias in your `~/.gitconfig`. Script will use this command to get the parent of the branch your current-branch has branched off of (wow, so many branches)

```shell
[alias]
    parent = "!git show-branch 2>&1 | grep '*' | grep -v \"$(git rev-parse --abbrev-ref HEAD)\" | head -n1 | sed 's/.*\\[\\(.*\\)\\].*/\\1/' | sed 's/[\\^~].*//' #"
```

## Things to do to run this utility

I created an alias for my local pointing to the binary of `rebase_with_upstream.go`.
My alias in `~/.zshrc` looks like below:

```shell
alias rebaseAll=$UPSTREAM/lazy_tasks/rebase_current_branch_with_remote/bin/rebaseAllDirs
```

**Note**: This alias hack assumes that

- _lazy_tasks/rebase_current_branch_with_remote/_ was checked out to your `$UPSTREAM` location
- `make all` was run at least once before creating the alias so that `lazy_tasks/rebase_current_branch_with_remote/bin/rebaseAllDirs` is created.