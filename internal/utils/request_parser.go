package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRequestBody(r *http.Request) (map[string]interface{}, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("request body is empty")
	}

	defer r.Body.Close()
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func GetBodyValue(body map[string]interface{}, key string) (interface{}, error) {
	value, exists := body[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in body", key)
	}
	return value, nil
}