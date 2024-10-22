# Get all contributions and update project board

This Project is aimed at fetching all the PRs I am working on and updating the project board with the PRs/issue.

## Installation

### To run update_project_board.go
`update_project_board.go` is the main script to live in this repo. It is aimed to fetch all the PRs and issues I am working on and update the project board with the PRs/issue.

To run the project:
- Create an file with the name `.env` in this folder.
  - Add `GITHUB_TOKEN=<Your_github_token>`
  - Copy paste this string as well in a new line `PROJECT_COLUMN_ID=your_project_column_id`. You will replace it later on as the code starts working.

- Then run below command to run the script: 
```bash
go run update_project_board.go
```

### find_projects.sh
`find_projects.sh` is a script to find all the projects in the repo. It will list all the projects in the repo.

To run the find_projects.sh:
```bash
./find_projects.sh
```
- Note: make sure you have given executable permissions to the script. If not, run `chmod +x find_projects.sh` to give executable permissions.
