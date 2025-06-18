package main

import (
	"context"
	"log"
	"time"

	"data-ingestion-service/internal/client"
	"data-ingestion-service/internal/config"
	"data-ingestion-service/internal/storage"
	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	
	ctx := context.Background()
	
	apiClient := client.New(cfg.APIEndpoint)
	
	bqClient, err := storage.NewBigQueryClient(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create BigQuery client:", err)
	}
	defer bqClient.Close()

	log.Println("Starting data ingestion service...")
	
	ticker := time.NewTicker(cfg.FetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ingestData(ctx, apiClient, bqClient); err != nil {
				log.Printf("Ingestion failed: %v", err)
			}
		}
	}
}

func ingestData(ctx context.Context, apiClient *client.Client, bqClient *storage.BigQueryClient) error {
	logs, err := apiClient.FetchLogs()
	if err != nil {
		return err
	}

	log.Println(logs)


	if len(logs) == 0 {
		log.Println("No logs to ingest")
		return nil
	}

	batchID := uuid.New().String()
	if err := bqClient.InsertRawLogs(ctx, logs, batchID); err != nil {
		return err
	}

	log.Printf("Successfully ingested %d logs with batch ID: %s", len(logs), batchID)
	return nil
}