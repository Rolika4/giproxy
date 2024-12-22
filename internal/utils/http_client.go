package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendRequestWithAuth(method, url, authType, tokenName string) ([]byte, int, error) {
	authToken := os.Getenv(tokenName)
	if authToken == "" {
		return nil, http.StatusUnauthorized, fmt.Errorf("%v is not set in the environment", tokenName)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %v", err)
	}

	authHeader := fmt.Sprintf("%s %s", authType, authToken)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // Используем io.ReadAll вместо ioutil.ReadAll
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, resp.StatusCode, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}