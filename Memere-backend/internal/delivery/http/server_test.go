package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSystemStatsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health/system", systemStatsHandler(nil, nil))

	req := httptest.NewRequest(http.MethodGet, "/health/system", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	if res["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", res["status"])
	}

	runtimeInfo, ok := res["runtime"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected runtime map in response")
	}

	if _, ok := runtimeInfo["goroutines"]; !ok {
		t.Errorf("expected goroutines in runtime info")
	}
	if _, ok := runtimeInfo["alloc_mb"]; !ok {
		t.Errorf("expected alloc_mb in runtime info")
	}
	if _, ok := runtimeInfo["sys_mb"]; !ok {
		t.Errorf("expected sys_mb in runtime info")
	}
}
