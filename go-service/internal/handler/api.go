package handler

import (
	"net/http"
	"strconv"

	"data-ingestion-service/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage *storage.BigQueryClient
}

func New(storage *storage.BigQueryClient) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 1000 {
		limit = 1000
	}

	logs, err := h.storage.GetProcessedLogs(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"limit":  limit,
		"offset": offset,
		"count":  len(logs),
	})
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}