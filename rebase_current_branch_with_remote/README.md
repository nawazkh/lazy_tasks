# Rebase current branch with remote

- This folder provides a script named `rebaseAllDirs`.
- Running `rebaseAllDirs` at the root of all the git repos will rebase the current checked out branch with `<remote>`/`<active development>`
  - `<remote>` : script will figure out if the remote is `upstream` or `origin`
  - `<active development>`: script will figure out if this is `master` or `main`

## Needed Variables

- Export `UPSTREAM` to the root of the dir containing all the git repositories
- In Makefile
  - update the `#VARIABLES` section for local testing
