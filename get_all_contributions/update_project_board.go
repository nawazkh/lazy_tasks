package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// PullRequest represents a GitHub PR object
type PullRequest struct {
	ID         int        `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	URL        string     `json:"html_url"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// Card represents the payload for adding a PR to a project column
type Card struct {
	ContentID   int    `json:"content_id"`
	ContentType string `json:"content_type"`
}

var githubToken string
var projectColumnID string

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	githubToken = os.Getenv("GITHUB_TOKEN")
	projectColumnID = os.Getenv("PROJECT_COLUMN_ID")
}

func main() {
	// Get user input for date range
	fmt.Print("Enter the start date (YYYY-MM-DD) or press Enter to skip: ")
	var startDateInput string
	fmt.Scanln(&startDateInput)

	fmt.Print("Enter the end date (YYYY-MM-DD) or press Enter to skip: ")
	var endDateInput string
	fmt.Scanln(&endDateInput)

	var startDate, endDate *time.Time
	if startDateInput != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDateInput)
		if err != nil {
			log.Fatalf("Invalid start date format: %v", err)
		}
		startDate = &parsedStartDate
	}

	if endDateInput != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDateInput)
		if err != nil {
			log.Fatalf("Invalid end date format: %v", err)
		}
		endDate = &parsedEndDate
	}

	// Fetch PRs authored by the user
	prs, err := fetchUserPullRequests(startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to fetch PRs: %v", err)
	}

	fmt.Printf("Found %d PRs.\n", len(prs))

	// Add each PR to the project board
	for _, pr := range prs {
		fmt.Println("Adding PR to project board...", pr.Title, pr.URL)
		//err := addPRToProject(pr.ID)
		//if err != nil {
		//	log.Printf("Error adding PR #%d (%s) to project: %v\n", pr.Number, pr.URL, err)
		//}
	}
}

func fetchUserPullRequests(startDate, endDate *time.Time) ([]PullRequest, error) {
	var prs []PullRequest
	page := 1

	// GitHub API Search endpoint to find all PRs authored by the user
	for {
		url := fmt.Sprintf("https://api.github.com/search/issues?q=author:@nojnhuh+type:pr&per_page=100&page=%d", page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+githubToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GitHub API request failed: %s", string(body))
		}

		// Response contains a 'total_count' field and a list of 'items'
		var searchResult struct {
			TotalCount int           `json:"total_count"`
			Items      []PullRequest `json:"items"`
		}
		if err := json.Unmarshal(body, &searchResult); err != nil {
			return nil, err
		}

		// Filter PRs based on the date range
		for _, pr := range searchResult.Items {
			if startDate != nil && pr.CreatedAt.Before(*startDate) {
				continue
			}
			if endDate != nil && pr.ClosedAt != nil && pr.ClosedAt.After(*endDate) {
				continue
			}
			prs = append(prs, pr)
		}

		// If there are no more pages of results, stop
		if len(searchResult.Items) < 100 {
			break
		}

		page++
	}
	return prs, nil
}

func addPRToProject(prID int) error {
	card := Card{
		ContentID:   prID,
		ContentType: "PullRequest",
	}

	jsonBody, err := json.Marshal(card)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/projects/columns/%s/cards", projectColumnID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("PR %d added to project board successfully.\n", prID)
		return nil
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		fmt.Printf("PR %d is already on the project board.\n", prID)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("Failed to add PR %d to project board: %s", prID, string(body))
}
