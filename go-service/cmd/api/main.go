package main

import (
	"context"
	"log"

	"data-ingestion-service/internal/config"
	"data-ingestion-service/internal/handler"
	"data-ingestion-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	ctx := context.Background()
	
	bqClient, err := storage.NewBigQueryClient(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create BigQuery client:", err)
	}
	defer bqClient.Close()

	h := handler.New(bqClient)
	
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	
	r.GET("/health", h.Health)
	r.GET("/logs", h.GetLogs)
	
	log.Printf("Starting API server on port %s...", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}