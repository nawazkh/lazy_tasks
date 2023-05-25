package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

var (
	infoLogger         *log.Logger
	errorLogger        *log.Logger
	reposRoot          = flag.String("reposRoot", "", "The root of all the repos you want to rebase with UPSTREAM.")
	customRemote       = flag.String("customRemote", "", "The custom remote. It is generally \"upstream\" or \"origin\"")
	customRemoteBranch = flag.String("customRemoteBranch", "", "The custom remote branch. It is generally \"master\" or \"main\"")
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO ", log.LstdFlags) // 5 -> 2023/05/24 14:52:51.507584, 4 -> 14:52:57.482177
	errorLogger = log.New(os.Stdout, "ERROR ", log.LstdFlags)
	flag.Parse()
}

func main() {
	os.Exit(run())
}

func run() int {
	// Setup

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

	// get all the dirs in root
	infoLogger.Printf("Root Dir: \"%s\"\n", *reposRoot)
	entries, err := os.ReadDir(*reposRoot)
	dirPath, _ := filepath.Abs(*reposRoot)
	if err != nil {
		errorLogger.Printf("error in fetching directories from %s dir. Error: %s\n", *reposRoot, err.Error())
		return 1
	}

	var wg sync.WaitGroup
	returnCode := 0
	resultChan := make(chan int, len(entries))

	// fetch UPSTREAM and attempt rebase
	for i, e := range entries {
		wg.Add(1)
		go performRebasing(i, &wg, resultChan, dirPath, e.Name())
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		returnCode += result
	}
	return returnCode
}

func performRebasing(id int, wg *sync.WaitGroup, resultChan chan<- int, dirPath, dirName string) {
	defer wg.Done()

	// get full directory dir
	dir := filepath.Join(dirPath, dirName)
	fileInfo, err := os.Stat(dir)
	if err != nil {
		errorLogger.Printf("Error in fetching the file info of the %s. Error: %s\n", dir, err.Error())
		resultChan <- 1
		return
	}

	if !fileInfo.IsDir() {
		resultChan <- 0
		return
	}

	// Is input directory a git repository ?
	isGitRepoCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	isGitRepoCmd.Dir = dir
	isGitRepo, err := isGitRepoCmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr.ExitCode() == 128 {
			resultChan <- 0
			return
			// infoLogger.Println(dir, "dir is not a git repo")
		}
		errorLogger.Printf("Failed to execute git command on the dir %s. Error: %s\n", dir, err.Error())
		resultChan <- 1
		return
	}
	if strings.TrimSpace(string(isGitRepo)) != "true" { // ignore other non git directories
		resultChan <- 0
		return
	}
	// infoLogger.Printf("Git repository:%s\n", dir)

	remote := ""
	if *customRemote == "" {
		// get remotes of the repository
		gitRemoteCmd := exec.Command("git", "remote", "show")
		gitRemoteCmd.Dir = dir
		gitRemote, err := gitRemoteCmd.CombinedOutput()
		if err != nil {
			errorLogger.Printf("Filed to execute git remote show on the dir %s. Error: %s\n", dir, err.Error())
			resultChan <- 1
			return
		}
		remotes := strings.Fields(strings.TrimSpace(string(gitRemote)))
		// infoLogger.Println(remotes)

		// check if upstream exists else select origin
		if slices.Contains(remotes, "upstream") {
			remote = "upstream"
		} else {
			remote = "origin"
		}
	} else {
		remote = *customRemote
	}
	// infoLogger.Printf("Remote: %s\n", remote)

	// fetch tracking remote branch
	gitFetchCmd := exec.Command("git", "fetch", remote)
	gitFetchCmd.Dir = dir
	_, err = gitFetchCmd.CombinedOutput()
	if err != nil {
		errorLogger.Printf("Filed to execute git fetch remote on the dir %s. Error: %s\n", dir, err.Error())
		resultChan <- 1
		return
	}
	// infoLogger.Println(string(gitFetch))

	remoteBranch := ""
	if *customRemoteBranch == "" {
		// remote master or main branch?
		gitRemoteBCmd := exec.Command("git", "ls-remote", "--exit-code", "--heads", remote, "master", remote, "main")
		gitRemoteBCmd.Dir = dir
		gitRemoteBranch, err := gitRemoteBCmd.CombinedOutput()
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok && exitErr.ExitCode() == 2 {
				errorLogger.Printf("no matching refs are found in %s master or for %s main. Error: %s \n", remote, remote, err.Error())
				resultChan <- 1
				return
			} else {
				errorLogger.Printf("Unable to run git ls-remote on dir:%s. Error:%s\n", dir, err.Error())
				resultChan <- 1
				return
			}
		}
		remoteBranchList := strings.Fields(strings.TrimSpace(string(gitRemoteBranch))) // assuming remote will not have master and main set together
		remoteBranch = strings.Split(remoteBranchList[1], "/")[2]
		// infoLogger.Printf("Remote Branch:%s\n", remoteBranch)
	} else {
		remoteBranch = *customRemoteBranch
	}

	// find default branch a.k.a current branch name
	gitCurrentBranchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	gitCurrentBranchCmd.Dir = dir
	currentBranchOp, err := gitCurrentBranchCmd.CombinedOutput()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr.ExitCode() == 128 {
			infoLogger.Printf("Failed to execute git symbolic-ref on the dir %s. Check dir manually.\n", dir)
			resultChan <- 0
			return
		}
		errorLogger.Printf("Failed to execute git symbolic-ref on the dir %s. Error: %s\n", dir, err)
		resultChan <- 1
		return
	}
	currentBranch := strings.TrimSpace(string(currentBranchOp))
	// infoLogger.Printf("Current Branch: %s", currentBranch)

	// run rebase on current branch using remote/remoteBranch
	gitRebaseCmd := exec.Command("git", "rebase", remote+"/"+remoteBranch)
	gitRebaseCmd.Dir = dir
	_, err = gitRebaseCmd.CombinedOutput()
	if err != nil {
		errorLogger.Printf("Error in performing rebase for dir %s. Error %s\n", dir, err.Error())
		resultChan <- 1
		return
	}
	infoLogger.Printf("Successfully rebased and updated %s's local/%s with %s/%s\n", dirName, currentBranch, remote, remoteBranch)
	resultChan <- 0
	return
}
