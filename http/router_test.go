package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouterInit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	InitRouter(v1)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/async", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
