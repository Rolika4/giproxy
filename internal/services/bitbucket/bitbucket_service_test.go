package services

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

func TestGetBranchesFromBitbucket_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, err := w.Write([]byte(`{"values": [{"name": "main", "target": {"hash": "abc123"}}]}`))
        if err != nil {
            t.Fatalf("could not write response: %v", err)
        }
    }))
    defer server.Close()

	os.Setenv("BITBUCKET_URL", server.URL)
	os.Setenv("BITBUCKET_TOKEN", "mockToken")

    body := map[string]interface{}{
        "owner": "testOwner",
        "repo":  "testRepo",
    }

    w := httptest.NewRecorder()

    GetBranchesFromBitbucket(w, body)

    resp := w.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %v", resp.StatusCode)
    }

    var response []map[string]string
    err := json.NewDecoder(resp.Body).Decode(&response)
    if err != nil {
        t.Fatalf("could not decode response: %v", err)
    }

    expected := []map[string]string{
        {"name": "main", "sha": "abc123"},
    }
    if len(response) != len(expected) || response[0]["name"] != expected[0]["name"] || response[0]["sha"] != expected[0]["sha"] {
        t.Errorf("expected %v, got %v", expected, response)
    }
}

func TestGetBranchesFromBitbucket_MissingOwner(t *testing.T) {
    body := map[string]interface{}{
        "repo": "testRepo",
    }

    w := httptest.NewRecorder()

    GetBranchesFromBitbucket(w, body)

    resp := w.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("expected status 400, got %v", resp.StatusCode)
    }
}

func TestGetBranchesFromBitbucket_MissingRepo(t *testing.T) {
    body := map[string]interface{}{
        "owner": "testOwner",
    }

    w := httptest.NewRecorder()

    GetBranchesFromBitbucket(w, body)

    resp := w.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("expected status 400, got %v", resp.StatusCode)
    }
}

func TestGetBranchesFromBitbucket_InvalidOwnerFormat(t *testing.T) {
    body := map[string]interface{}{
        "owner": 123,
        "repo":  "testRepo",
    }

    w := httptest.NewRecorder()

    GetBranchesFromBitbucket(w, body)

    resp := w.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("expected status %v; got %v", http.StatusBadRequest, resp.StatusCode)
    }

    var responseBody map[string]interface{}
    err := json.NewDecoder(resp.Body).Decode(&responseBody)
    if err != nil {
        t.Fatalf("could not decode response: %v", err)
    }

    expectedError := "Invalid 'owner' field format"
    if responseBody["error"] != expectedError {
        t.Errorf("expected error message %v; got %v", expectedError, responseBody["error"])
    }
}