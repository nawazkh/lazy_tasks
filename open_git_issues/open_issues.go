package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const baseURL = "https://api.github.com"

type Issue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type IssueResponse struct {
	HTMLURL string `json:"html_url"`
}

func main() {
	token := flag.String("token", "", "GitHub personal access token")
	repoFile := flag.String("file", "", "Path to the text file containing the list of repositories")
	issueTitle := flag.String("title", "", "Issue title")
	issueFile := flag.String("body", "", "Path to the markdown file containing the issue body")
	dryRun := flag.String("dryRun", "", "Dry run")
	flag.Parse()

	if *token == "" {
		fmt.Println("GitHub personal access token is required.")
		os.Exit(1)
	}

	if *repoFile == "" {
		fmt.Println("Path to the text file containing the list of repositories is required.")
		os.Exit(1)
	}

	if *issueTitle == "" {
		fmt.Println("Issue title is required.")
		os.Exit(1)
	}

	issueBody := ""
	if *issueFile != "" {
		fileContent, err := os.ReadFile(*issueFile)
		if err != nil {
			fmt.Printf("Failed to read file: %s\n", err)
			os.Exit(1)
		}
		issueBody = string(fileContent)
	}

	repoList, err := readLines(*repoFile)
	if err != nil {
		fmt.Printf("Failed to read repository list: %s\n", err)
		os.Exit(1)
	}

	if *dryRun == "true" {
		fmt.Printf("\n\n")
		fmt.Println("###############################################")
		fmt.Println("Dry run is enabled. No issues will be created.")
		fmt.Println("###############################################")
		fmt.Printf("\n\n")
	}

	fmt.Println("\n\nCreating issues for the following repositories:")
	fmt.Println(strings.Join(repoList, "\n"))
	fmt.Printf("\n")

	if *dryRun != "true" {
		fmt.Println("Shall I continue creating issues? (y/n)")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			fmt.Println("Aborting...")
			os.Exit(0)
		}
	}
	fmt.Printf("\n")

	for _, repo := range repoList {
		issue := Issue{
			Title: *issueTitle,
			Body:  issueBody,
		}

		issueJSON, err := json.Marshal(issue)
		if err != nil {
			fmt.Printf("Failed to marshal issue: %s\n", err)
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/repos/%s/issues", baseURL, repo)
		if *dryRun == "true" {
			continue
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(issueJSON))
		// fmt.Println("Request will be sent at:", req.URL)
		if err != nil {
			fmt.Printf("Failed to create request: %s\n", err)
			os.Exit(1)
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
		req.Header.Set("X-Github-Api-Version", "2022-11-28")
		req.Header.Set("User-Agent", "open_git_issues")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request: %s\n", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Failed to create issue for repository '%s'\nStatus code: %d\nStatus:%s\n\n\n", repo, resp.StatusCode, resp.Status)
		} else {
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to read response body: %s\n", err)
				os.Exit(1)
			}

			var issueResponse IssueResponse
			err = json.Unmarshal(responseBody, &issueResponse)
			if err != nil {
				fmt.Printf("Failed to unmarshal issue response: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("\nIssue created for repository '%s'\nURL: %s\n", repo, issueResponse.HTMLURL)
		}
	}
}

// readLines reads a text file and returns a slice of lines.
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	return lines, scanner.Err()
}
