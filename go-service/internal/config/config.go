package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIEndpoint    string
	GCPProjectID   string
	Dataset        string
	RawTable       string
	ProdTable      string
	Source         string
	BatchSize      int
	FetchInterval  time.Duration
	Port           string
}

func Load() *Config {
	godotenv.Load()
	
	batchSize, _ := strconv.Atoi(getEnv("BATCH_SIZE", "100"))
	interval, _ := time.ParseDuration(getEnv("FETCH_INTERVAL", "5m"))
	
	return &Config{
		APIEndpoint:   getEnv("API_ENDPOINT", ""),
		GCPProjectID:  getEnv("GCP_PROJECT_ID", ""),
		Dataset:       getEnv("DATASET", "data_ingestion"),
		RawTable:      getEnv("RAW_TABLE", "raw_logs"),
		ProdTable:     getEnv("PROD_TABLE", "processed_logs"),
		Source:        getEnv("SOURCE", "placeholder_api"),
		BatchSize:     batchSize,
		FetchInterval: interval,
		Port:          getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}