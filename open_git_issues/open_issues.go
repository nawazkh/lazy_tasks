package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
		fileContent, err := ioutil.ReadFile(*issueFile)
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
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(issueJSON))
		fmt.Println("request will be sent to", req.URL)
		if err != nil {
			fmt.Printf("Failed to create request: %s\n", err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request: %s\n", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Failed to create issue for repository '%s'. Status code: %d, Response: Status:%s\n\n\n", repo, resp.StatusCode, resp.Status)
		} else {
			responseBody, err := ioutil.ReadAll(resp.Body)
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

			fmt.Printf("Issue created for repository '%s'. URL: %s\n\n\n", repo, issueResponse.HTMLURL)
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
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	return lines, scanner.Err()
}
