package handlers

import (
	"log"
	"net/http"
	"giproxy/internal/services/gitlab"
	"giproxy/internal/utils"
)


func HandleGitLabRequest(w http.ResponseWriter, r *http.Request) {
	var RouteMapGitlab = map[string]func(http.ResponseWriter, map[string]interface{}){
		"/api/git/branches": services.GetBranchesFromGitlab,
	}
	
	body, err := utils.ParseRequestBody(r)

	if err != nil {
		log.Println("Failed to parse body:", err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if handlerFunc, exists := RouteMapGitlab[r.URL.Path]; exists {
		handlerFunc(w, body)
	} else {
		log.Printf("No handler found for path: %s\n", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}
}