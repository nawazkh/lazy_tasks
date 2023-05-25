# Rebase current branch with remote

- This folder provides a script named `rebase_with_upstream.go` and generates a binary `./bin/rebaseAllDirs`.
- Running `rebaseAllDirs` at the root of all the git repos will rebase the locally checked out branch with `<remote>`/`<active development>`
  - `<remote>` : script will figure out if the remote is `upstream` or `origin`
  - `<active development>`: script will figure out if this is `master` or `main`

## Needed Variables

- Export `UPSTREAM` to the root of the dir containing all the git repositories
- In Makefile
  - update the `#VARIABLES` section for local testing

## Things to do to run this utility

I created an alias for my local pointing to the binary of `rebase_with_upstream.go`.
My alias in `~/.zshrc` looks like below:

```shell
alias rebaseAll=$UPSTREAM/lazy_tasks/rebase_current_branch_with_remote/bin/rebaseAllDirs
```

**Note**: This alias hack assumes that

- _lazy_tasks/rebase_current_branch_with_remote/_ was checked out to your `$UPSTREAM` location
- `make all` was run at least once before creating the alias so that `lazy_tasks/rebase_current_branch_with_remote/bin/rebaseAllDirs` is created.