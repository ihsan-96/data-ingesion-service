package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"data-ingestion-service/internal/config"
	"data-ingestion-service/internal/models"
	"google.golang.org/api/iterator"
)

type BigQueryClient struct {
	client   *bigquery.Client
	config   *config.Config
	rawTable *bigquery.Table
}

func NewBigQueryClient(ctx context.Context, cfg *config.Config) (*BigQueryClient, error) {
	client, err := bigquery.NewClient(ctx, cfg.GCPProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %w", err)
	}

	rawTable := client.Dataset(cfg.Dataset).Table(cfg.RawTable)

	return &BigQueryClient{
		client:   client,
		config:   cfg,
		rawTable: rawTable,
	}, nil
}

func (bq *BigQueryClient) InsertRawLogs(ctx context.Context, logs []models.LogEntry, batchID string) error {
	var rawLogs []models.RawLog
	now := time.Now().UTC()

	for _, log := range logs {
		payload, _ := json.Marshal(log)
		rawLogs = append(rawLogs, models.RawLog{
			Payload:    string(payload),
			IngestedAt: now,
			Source:     bq.config.Source,
			BatchID:    batchID,
		})
	}

	inserter := bq.rawTable.Inserter()
	return inserter.Put(ctx, rawLogs)
}

func (bq *BigQueryClient) GetProcessedLogs(ctx context.Context, limit int, offset int) ([]models.ProcessedLog, error) {
	query := fmt.Sprintf(`
		SELECT * FROM `+"`%s.%s.%s`"+` 
		ORDER BY created_at DESC 
		LIMIT %d OFFSET %d`,
		bq.config.GCPProjectID, bq.config.Dataset, bq.config.ProdTable, limit, offset)

	q := bq.client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	var logs []models.ProcessedLog
	for {
		var log models.ProcessedLog
		err := it.Next(&log)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (bq *BigQueryClient) Close() error {
	return bq.client.Close()
}