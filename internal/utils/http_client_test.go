package utils

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestSendRequestWithAuth_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, err := w.Write([]byte(`{"status": "ok"}`))
        if err != nil {
            t.Fatalf("could not write response: %v", err)
        }
    }))
    defer server.Close()

    t.Setenv("BITBUCKET_TOKEN", "mockToken")

    body, statusCode, err := SendRequestWithAuth("GET", server.URL, "Bearer", "BITBUCKET_TOKEN")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if statusCode != http.StatusOK {
        t.Errorf("expected status 200, got %v", statusCode)
    }

    if string(body) != `{"status": "ok"}` {
        t.Errorf("unexpected response body: %s", string(body))
    }
}