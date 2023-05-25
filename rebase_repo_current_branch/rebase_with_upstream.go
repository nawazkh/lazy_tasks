package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	infoLogger      *log.Logger
	errorLogger     *log.Logger
	reposRoot       = flag.String("reposRoot", "", "The root of all the repos you want to rebase with UPSTREAM.")
	defaultUpstream = flag.String("defaultUpstream", "main", "The default source of truth at UPSTREAM. It could be `main` or `master`.")
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO ", 5) // 5 -> 2023/05/24 14:52:51.507584, 4 -> 14:52:57.482177
	errorLogger = log.New(os.Stdout, "ERROR ", 5)
	flag.Parse()
}

func main() {
	os.Exit(run())
}

func run() int {
	// Setup
	// find the right root
	if reposRoot == nil {
		errorLogger.Println("reposRoot is nil")
		return 1
	}

	if defaultUpstream == nil {
		errorLogger.Println("defaultUpstream is nil")
		return 1
	}

	if *reposRoot == "" {
		value, isDefined := os.LookupEnv("UPSTREAM")
		if isDefined {
			*reposRoot = value
		}

		if !isDefined || value == "" {
			errorLogger.Println("neither reposRoot nor $UPSTREAM was defined. Exiting")
			return 1
		}
	}

	// TODO: uncomment the below lines if you want to enforce rebasing with `main` or `master` branch only
	// if !(*defaultUpstream == "main" || *defaultUpstream == "master") {
	// 	errorLogger.Println("default upstream bran")
	// }

	// get all the dirs in root
	infoLogger.Printf("Root Dir: \"%s\"\n\n", *reposRoot)
	// infoLogger.Printf("\"%s\" <------ Upstream's branch to rebase existing branch with\n\n\n", *defaultUpstream)

	entries, err := os.ReadDir(*reposRoot)
	dirPath, _ := filepath.Abs(*reposRoot)
	if err != nil {
		errorLogger.Printf("error in fetching directories from %s dir\n", *reposRoot)
		errorLogger.Println(err)
		return 1
	}

	// fetch UPSTREAM and attempt rebase
	// TODO: launch parallel go routines for each fetch and rebase
	for _, e := range entries {
		// get full directory dir
		dir := filepath.Join(dirPath, e.Name())
		fileInfo, err := os.Stat(dir)
		if err != nil {
			errorLogger.Println("Error in fetching the file info of the", dir, "dir", err)
			continue
		}

		if !fileInfo.IsDir() {
			continue
		}

		// Is input directory a git repository ?
		isGitRepoCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		isGitRepoCmd.Dir = dir
		isGitRepo, err := isGitRepoCmd.CombinedOutput()
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok && exitErr.ExitCode() == 128 {
				infoLogger.Println(dir, "dir is not a git repo")
			} else {
				errorLogger.Println("Filed to execute git command on the dir", dir, ". Error: ", err)
			}
			continue
		}
		if strings.TrimSpace(string(isGitRepo)) != "true" { // ignore other non git directories
			continue
		}
		infoLogger.Printf("Git repository:%s\n", dir)

		// get remotes of the repository
		gitRemoteCmd := exec.Command("git", "remote", "show")
		gitRemoteCmd.Dir = dir
		gitRemote, err := gitRemoteCmd.CombinedOutput()
		if err != nil {
			errorLogger.Println("Filed to execute git remote show on the dir", dir, ". Error: ", err)
			return 1
		}
		remotes := strings.Fields(strings.TrimSpace(string(gitRemote)))
		// infoLogger.Println(remotes)

		// check if upstream exists else select origin
		remote := ""
		if slices.Contains(remotes, "upstream") {
			remote = "upstream"
		} else {
			remote = "origin"
		}

		// fetch tracking remote branch
		infoLogger.Printf("Remote: %s\n", remote)
		gitFetchCmd := exec.Command("git", "fetch", remote)
		gitFetchCmd.Dir = dir
		_, err = gitFetchCmd.CombinedOutput()
		if err != nil {
			errorLogger.Println("Filed to execute git fetch remote on the dir", dir, ". Error: ", err)
			return 1
		}
		// infoLogger.Println(string(gitFetch))

		// remote master or main branch?
		gitRemoteBCmd := exec.Command("git", "ls-remote", "--exit-code", "--heads", remote, "master", remote, "main")
		gitRemoteBCmd.Dir = dir
		gitRemoteBranch, err := gitRemoteBCmd.CombinedOutput()
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok && exitErr.ExitCode() == 2 {
				errorLogger.Println("no matching refs are found in", remote, "master and for", remote, "main", err)
			} else {
				errorLogger.Print("Unable to run git ls-remote on dir", dir, ". Error:", err)
				return 1
			}
		}
		remoteBranchList := strings.Fields(strings.TrimSpace(string(gitRemoteBranch))) // assuming remote will not have master and main set together
		remoteBranch := strings.Split(remoteBranchList[1], "/")[2]
		infoLogger.Printf("Remote Branch:%s\n", remoteBranch)

		// find default branch a.k.a current branch name
		gitCurrentBranchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
		gitCurrentBranchCmd.Dir = dir
		currentBranchOp, err := gitCurrentBranchCmd.CombinedOutput()
		if err != nil {
			errorLogger.Println("Filed to execute git symbolic-ref on the dir", dir, ". Error: ", err)
			return 1
		}
		currentBranch := strings.TrimSpace(string(currentBranchOp))
		infoLogger.Printf("Current Branch: %s", currentBranch)

		// run rebase on current branch using remote/remoteBranch
		gitRebaseCmd := exec.Command("git", "rebase", remote+"/"+remoteBranch)
		gitRebaseCmd.Dir = dir
		_, err = gitRebaseCmd.CombinedOutput()
		if err != nil {
			errorLogger.Println("Error in performing rebase for dir", dir, "check manually")
			return 1
		}
		infoLogger.Printf("Successfully rebased and updated local/%s with %s/%s\n", currentBranch, remote, remoteBranch)
		break
	}

	// return 0 on all successful rebases
	// return any error code if any one of them fails
	return 0
}