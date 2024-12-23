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

func GetWidgetSonarQube(w http.ResponseWriter, r *http.Request) {

	var sonarApiPath = os.Getenv("SONARQUBE_URL")

	component := utils.GetQueryParam(r, "component")
	metricKeys := utils.GetQueryParam(r, "metricKeys")

	fmt.Println("SonarQube widget", component, metricKeys)

	link := fmt.Sprintf("%s/api/measures/component?component=%s&metricKeys=%s", sonarApiPath, component, metricKeys)
	response, statusCode, err := utils.SendRequestWithAuth("GET", link, "Basic", "SONARQUBE_TOKEN")
	if err != nil {
		log.Printf(`{"error": "failed to get sonarqube project info: %v"}`, err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "failed to get sonarqube project info: %v"}`, err), statusCode)
		return
	}

	var tryDecode bool
	decodedCollection := response

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

	output := map[string]interface{}{
		"widget":     "sonarqube",
		"request":    "get",
		"component":  component,
		"metricKeys": metricKeys,
		"responce":   decodedCollection,
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
	_, writeErr := w.Write(decodedCollection)
	if writeErr != nil {
		log.Printf("Failed to write response: %v", writeErr)
	}
}
