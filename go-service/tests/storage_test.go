package tests

import (
	"testing"

	"data-ingestion-service/internal/models"
)

func TestRawLogModel(t *testing.T) {
	log := models.RawLog{
		Payload: `{"userId":1,"id":1,"title":"test","body":"test body"}`,
		Source:  "test_api",
		BatchID: "test-batch-123",
	}

	if log.Payload == "" {
		t.Fatal("Payload should not be empty")
	}

	if log.Source != "test_api" {
		t.Fatalf("Expected source 'test_api', got '%s'", log.Source)
	}
}

func TestLogEntryModel(t *testing.T) {
	entry := models.LogEntry{
		UserID: 1,
		ID:     1,
		Title:  "test",
		Body:   "test body",
	}

	if entry.UserID != 1 {
		t.Fatalf("Expected UserID 1, got %d", entry.UserID)
	}

	if entry.Title != "test" {
		t.Fatalf("Expected title 'test', got '%s'", entry.Title)
	}
}