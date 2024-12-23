package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"giproxy/internal/utils"
	"log"
	"net/http"
	"os"
)

func GetWidgetDepTrack(w http.ResponseWriter, r *http.Request) {

	var depTrackApiPath = os.Getenv("DEPTRACK_URL")

	name := utils.GetQueryParam(r, "name")

	link := fmt.Sprintf("%s/api/v1/project?name=%s", depTrackApiPath, name)
	response, statusCode, err := utils.SendRequestWithAuth("GET", link, "X-Api-Key", "DEPTRACK_TOKEN")
	if err != nil {
		log.Printf("Failed to get dependency-track project info: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to get dependency-track project info: %v"}`, err), statusCode)
		return
	}

	decodedCollection := response
	var tryDecode bool

	_, err = base64.StdEncoding.DecodeString(string(response))
	if err == nil {
		tryDecode = true
	}

	if tryDecode {
		decodedCollection, err = base64.StdEncoding.DecodeString(string(response))
		if err != nil {
			log.Printf("Failed to decode base64 collection: %v", err)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Failed to decode base64 collection"}`, http.StatusInternalServerError)
			return
		}
	}

	var decodedJSON interface{}
	err = json.Unmarshal(decodedCollection, &decodedJSON)
	if err != nil {
		log.Printf("Failed to unmarshal decoded response: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Failed to parse decoded collection"}`, http.StatusInternalServerError)
		return
	}

	wrappedResponse := map[string]interface{}{
		"collection": decodedJSON,
	}

	prettyResponce, err := json.Marshal(wrappedResponse)
	if err != nil {
		log.Printf("Failed to serialize output: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Failed to serialize output"}`, http.StatusInternalServerError)
		return
	}

	output := map[string]interface{}{
		"widget":       "dependency-track",
		"request":      "get",
		"project-name": name,
		"response":     wrappedResponse,
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
	_, writeErr := w.Write(prettyResponce)
	if writeErr != nil {
		log.Printf("Failed to write response: %v", writeErr)
	}
}
