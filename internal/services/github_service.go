package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"giproxy/internal/utils"
)




func GetBranchesFromGitHub(w http.ResponseWriter, body map[string]interface{}) {
	var githubApiPath = "https://api.github.com"

	type BranchInfo struct {
		Name   string `json:"name"`
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}
	
	owner, err := utils.GetBodyValue(body, "owner")
	if err != nil {
		log.Printf("Error fetching owner: %v", err)
		http.Error(w, "Missing or invalid 'owner' field in request body", http.StatusBadRequest)
		return
	}

	repo, err := utils.GetBodyValue(body, "repo")
	if err != nil {
		log.Printf("Error fetching repo: %v", err)
		http.Error(w, "Missing or invalid 'repo' field in request body", http.StatusBadRequest)
		return
	}

	ownerStr, ok := owner.(string)
	if !ok {
		log.Println("Owner is not a valid string")
		http.Error(w, "Invalid 'owner' field format", http.StatusBadRequest)
		return
	}

	repoStr, ok := repo.(string)
	if !ok {
		log.Println("Repo is not a valid string")
		http.Error(w, "Invalid 'repo' field format", http.StatusBadRequest)
		return
	}

	link := fmt.Sprintf("%s/repos/%s/%s/branches", githubApiPath, ownerStr, repoStr)

	response, statusCode, err :=  utils.SendRequestWithAuth("GET", link, "Bearer", "GITHUB_TOKEN")
	if err != nil {
		log.Printf("Failed to fetch branches: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch branches: %v", err), statusCode)
		return
	}

	var branches []BranchInfo
	err = json.Unmarshal(response, &branches)
	if err != nil {
		log.Printf("Failed to parse response: %v", err)
		http.Error(w, "Failed to parse response from GitHub", http.StatusInternalServerError)
		return
	}

	var filteredBranches []map[string]string
	for _, branch := range branches {
		filteredBranches = append(filteredBranches, map[string]string{
			"name": branch.Name,
			"sha":  branch.Commit.SHA,
		})
	}

	result, err := json.Marshal(filteredBranches)
	if err != nil {
		log.Printf("Failed to serialize filtered branches: %v", err)
		http.Error(w, "Failed to process branches data", http.StatusInternalServerError)
		return
	}

	output := map[string]interface{}{
		"git_provider": "github",
		"request":      "branches",
		"workspace":    ownerStr,
		"repository":   repo,
		"response":     branches,
	}

	prettyJSON, err := json.Marshal(output)

	log.Printf("%s", string(prettyJSON))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(result)
	if writeErr != nil {
		log.Printf("Failed to write response: %v", writeErr)
	}
}