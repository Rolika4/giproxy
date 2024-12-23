package handlers

import (
	deptrack "giproxy/internal/services/deptrack"
	sonarqube "giproxy/internal/services/sonarqube"
	"log"
	"net/http"
	"strings"
)

func WidgetHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if strings.HasPrefix(path, "/widgets/deptrack") {
		deptrack.GetWidgetDepTrack(w, r)
	} else if strings.HasPrefix(path, "/widgets/sonarqube") {
		sonarqube.GetWidgetSonarQube(w, r)
	} else {
		log.Println(`{"error": "Unsupported widget type"}`)
		http.Error(w, `{"error": "Unsupported widget type"}`, http.StatusBadRequest)
	}
}
