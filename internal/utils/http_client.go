package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SendRequestWithAuth(method, url, authType string) ([]byte, error) {
	authToken := os.Getenv("GIT_TOKEN")
	if authToken == "" {
		return nil, errors.New("GIT_TOKEN is not set in the environment")
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	authHeader := fmt.Sprintf("%s %s", authType, authToken)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return body, nil
}