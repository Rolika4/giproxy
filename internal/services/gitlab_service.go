package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"giproxy/internal/utils"
	"os"
)


func GetBranchesFromGitlab(w http.ResponseWriter, body map[string]interface{}) {

	var gitlabApiPath = os.Getenv("GITLAB_URL")

	type BranchInfo struct {
		Name   string `json:"name"`
		Commit struct {
			SHA string `json:"id"`
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

	ownerRepoPath := fmt.Sprintf("%s/%s", ownerStr, repoStr)
	link := fmt.Sprintf("%s/api/v4/projects/%s/repository/branches", gitlabApiPath, url.PathEscape(ownerRepoPath))

	response, err := utils.SendRequestWithAuth("GET", link, "Bearer", "GITLAB_TOKEN")
	if err != nil {
		log.Printf("Failed to fetch branches: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch branches: %v", err), http.StatusInternalServerError)
		return
	}

	var branches []BranchInfo
	err = json.Unmarshal(response, &branches)
	if err != nil {
		log.Printf("Failed to parse response: %v", err)
		http.Error(w, "Failed to parse response from Gitlab", http.StatusInternalServerError)
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
		"git_provider": "gitlab",
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