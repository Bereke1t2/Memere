package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
)

// gatherFamily returns the named MetricFamily from the default Prometheus
// registry, or nil if not found.
func gatherFamily(t *testing.T, name string) bool {
	t.Helper()
	all, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather error: %v", err)
	}
	for _, mf := range all {
		if mf.GetName() == name && len(mf.GetMetric()) > 0 {
			return true
		}
	}
	return false
}

func TestMetricsMiddleware_CountsRequest(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Metrics())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ping", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !gatherFamily(t, "http_requests_total") {
		t.Error("http_requests_total not found in registry after request")
	}
}

func TestMetricsMiddleware_RecordsDuration(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Metrics())
	r.GET("/latency", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/latency", nil))

	if !gatherFamily(t, "http_request_duration_seconds") {
		t.Error("http_request_duration_seconds not found in registry after request")
	}
}

func TestMetricsMiddleware_Status5xx(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Metrics())
	r.GET("/fail", func(c *gin.Context) { c.Status(http.StatusInternalServerError) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/fail", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if !gatherFamily(t, "http_requests_total") {
		t.Error("http_requests_total not found after 5xx response")
	}
}
