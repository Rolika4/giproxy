package services

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "giproxy/internal/utils"
)

func GetBranchesFromBitbucket(w http.ResponseWriter, body map[string]interface{}) {
    var bitbucketApiPath = os.Getenv("BITBUCKET_URL")

    var values struct {
        Values json.RawMessage `json:"values"`
    }

    type BranchInfo struct {
        Name  string `json:"name"`
        Commit struct {
            SHA string `json:"hash"`
        } `json:"target"`
    }

    owner, err := utils.GetBodyValue(body, "owner")
    if err != nil {
        log.Printf("Error fetching owner: %v", err)
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Missing or invalid 'owner' field in request body"}`, http.StatusBadRequest)
        return
    }

    repo, err := utils.GetBodyValue(body, "repo")
    if err != nil {
        log.Printf("Error fetching repo: %v", err)
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Missing or invalid 'repo' field in request body"}`, http.StatusBadRequest)
        return
    }

    ownerStr, ok := owner.(string)
    if !ok {
        log.Println("Owner is not a valid string")
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Invalid 'owner' field format"}`, http.StatusBadRequest)
        return
    }

    repoStr, ok := repo.(string)
    if !ok {
        log.Println("Repo is not a valid string")
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Invalid 'repo' field format"}`, http.StatusBadRequest)
        return
    }

    link := fmt.Sprintf("%s/2.0/repositories/%s/%s/refs/branches", bitbucketApiPath, ownerStr, repoStr)

	response, statusCode, err := utils.SendRequestWithAuth("GET", link, "Basic", "BITBUCKET_TOKEN")
	if err != nil {
        log.Printf("Failed to fetch branches: %v", err)
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch branches: %v"}`, err), statusCode)
        return
    }

    var branches []BranchInfo
    err = json.Unmarshal(response, &values)
    if err != nil {
        log.Printf("Failed to parse response val: %v", err)
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Failed to parse response from Bitbucket"}`, http.StatusInternalServerError)
        return
    }

    err = json.Unmarshal(values.Values, &branches)
    if err != nil {
        log.Printf("Failed to parse 'values' block: %v", err)
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Failed to parse response values"}`, http.StatusInternalServerError)
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
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Failed to process branches data"}`, http.StatusInternalServerError)
        return
    }

    output := map[string]interface{}{
        "git_provider": "bitbucket",
        "request":      "branches",
        "workspace":    ownerStr,
        "repository":   repo,
        "response":     branches,
    }

    prettyJSON, err := json.Marshal(output)
    if err != nil {
        log.Printf("Failed to serialize output: %v", err)
        w.Header().Set("Content-Type", "application/json")
        http.Error(w, `{"error": "Failed to serialize output"}`, http.StatusInternalServerError)
        return
    }

    log.Printf("%s", string(prettyJSON))

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, writeErr := w.Write(result)
    if writeErr != nil {
        log.Printf("Failed to write response: %v", writeErr)
    }
}