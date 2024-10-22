#!/bin/bash

# Organization name and output file
ORG="Azure"
OUTPUT_FILE="projects_list.txt"

# Clear the output file if it exists, or create a new one
> "$OUTPUT_FILE"

# Fetch the list of projects using gh and paginate through the results
PAGE=1

while true; do
  # Use `gh api` to fetch projects from the organization
  RESPONSE=$(gh api --paginate "/orgs/$ORG/projects?per_page=100&page=$PAGE" --jq '.')

  # Check if the response contains any projects
  PROJECT_COUNT=$(echo "$RESPONSE" | jq 'length')

  # If no projects are found, stop fetching
  if [[ $PROJECT_COUNT -eq 0 ]]; then
    echo "No more projects found. Exiting."
    break
  fi

  # Save project details to the output file
  echo "$RESPONSE" | jq -r '.[] | "Project Name: \(.name)\nDescription: \(.body // "No description")\nURL: \(.html_url)\nID: \(.id)\n---"' >> "$OUTPUT_FILE"

  # Print a message to indicate progress
  echo "Page $PAGE processed and written to $OUTPUT_FILE with $PROJECT_COUNT projects."

  # Increment page counter
  PAGE=$((PAGE + 1))
done

echo "All projects have been saved to $OUTPUT_FILE."
