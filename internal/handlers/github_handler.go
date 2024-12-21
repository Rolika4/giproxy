package handlers

import (
	"log"
	"net/http"
	"giproxy/internal/services"
	"giproxy/internal/utils"
)



func HandleGitHubRequest(w http.ResponseWriter, r *http.Request) {
	var RouteMapGitHub = map[string]func(http.ResponseWriter, map[string]interface{}){
		"/api/git/branches": services.GetBranchesFromGitHub,
	}
	
	body, err := utils.ParseRequestBody(r)
	
	if err != nil {
		log.Println("Failed to parse body:", err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if handlerFunc, exists := RouteMapGitHub[r.URL.Path]; exists {
		handlerFunc(w, body)
	} else {
		log.Printf("No handler found for path: %s\n", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}
}