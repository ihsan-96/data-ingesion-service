package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"data-ingestion-service/internal/models"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
}

func New(endpoint string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		endpoint:   endpoint,
	}
}

func (c *Client) FetchLogs() ([]models.LogEntry, error) {
	resp, err := c.httpClient.Get(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var logs []models.LogEntry
	if err := json.Unmarshal(body, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return logs, nil
}