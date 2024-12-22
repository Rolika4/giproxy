package utils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestParseRequestBody_Success(t *testing.T) {
	body := map[string]interface{}{
		"key": "value",
	}
	mockBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(mockBody))

	parsedBody, err := ParseRequestBody(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parsedBody["key"] != "value" {
		t.Errorf("expected key 'key' to have value 'value', got %v", parsedBody["key"])
	}
}

func TestParseRequestBody_EmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)

	_, err := ParseRequestBody(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Проверка на ошибку EOF, которая возникает при пустом теле
	if err.Error() != "EOF" {
		t.Errorf("expected error 'EOF', got '%v'", err.Error())
	}
}

func TestParseRequestBody_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{invalid json`))

	_, err := ParseRequestBody(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetQueryParam_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/?key=value", nil)
	param := GetQueryParam(req, "key")
	if param != "value" {
		t.Errorf("expected 'value', got '%v'", param)
	}
}

func TestGetQueryParam_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "/?key=value", nil)
	param := GetQueryParam(req, "nonexistent")
	if param != "" {
		t.Errorf("expected empty string, got '%v'", param)
	}
}

func TestGetBodyValue_Success(t *testing.T) {
	body := map[string]interface{}{
		"key": "value",
	}

	value, err := GetBodyValue(body, "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if value != "value" {
		t.Errorf("expected 'value', got %v", value)
	}
}

func TestGetBodyValue_MissingKey(t *testing.T) {
	body := map[string]interface{}{
		"key": "value",
	}

	_, err := GetBodyValue(body, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := "key missing not found in body"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%v', got '%v'", expectedErr, err.Error())
	}
}
