package handlers

import (
	"log"
	"net/http"
	"giproxy/internal/utils"
)


func CommonHandler(w http.ResponseWriter, r *http.Request) {

	gitserver := utils.GetQueryParam(r, "gitserver")

	switch gitserver {
	case "github":
		HandleGitHubRequest(w, r)
	case "bitbucket":
		HandleBitbucketRequest(w, r)
	case "gitlab":
		HandleGitLabRequest(w, r)
	default:
		http.Error(w, "Unsupported git server", http.StatusBadRequest)
		log.Printf("Unsupported git server: %s", gitserver)
	}
}