package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"data-ingestion-service/internal/client"
)

func TestFetchLogs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"userId":1,"id":1,"title":"test","body":"test body"}]`))
	}))
	defer server.Close()

	c := client.New(server.URL)
	logs, err := c.FetchLogs()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(logs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(logs))
	}

	if logs[0].Title != "test" {
		t.Fatalf("Expected title 'test', got '%s'", logs[0].Title)
	}
}

func TestFetchLogs_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := client.New(server.URL)
	_, err := c.FetchLogs()

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
