package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"data-ingestion-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	h := &handler.Handler{}
	
	r := gin.Default()
	r.GET("/health", h.Health)
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}