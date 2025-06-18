package models

import "time"

type LogEntry struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type RawLog struct {
	Payload     string    `bigquery:"payload"`
	IngestedAt  time.Time `bigquery:"ingested_at"`
	Source      string    `bigquery:"source"`
	BatchID     string    `bigquery:"batch_id"`
}

type ProcessedLog struct {
	LogID       int       `json:"log_id" bigquery:"log_id"`
	UserID      int       `json:"user_id" bigquery:"user_id"`
	Title       string    `json:"title" bigquery:"title"`
	Body        string    `json:"body" bigquery:"body"`
	WordCount   int       `json:"word_count" bigquery:"word_count"`
	IngestedAt  time.Time `json:"ingested_at" bigquery:"ingested_at"`
	Source      string    `json:"source" bigquery:"source"`
	ProcessDate string    `json:"process_date" bigquery:"process_date"`
	CreatedAt   time.Time `json:"created_at" bigquery:"created_at"`
}