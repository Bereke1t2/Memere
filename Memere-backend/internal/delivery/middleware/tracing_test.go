package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
)

// TestTracingMiddleware_PropagatesContext verifies that the tracing middleware
// runs without error and that the handler still receives a valid context. A
// no-op tracer is configured by default (OTEL_ENDPOINT="") so no real exporter
// is needed.
func TestTracingMiddleware_PropagatesContext(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Tracing())

	handlerCalled := false
	r.GET("/trace-test", func(c *gin.Context) {
		handlerCalled = true
		if c.Request.Context() == nil {
			t.Error("handler context must not be nil")
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trace-test", nil)
	r.ServeHTTP(w, req)

	if !handlerCalled {
		t.Fatal("handler was not called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// TestTracingMiddleware_UnknownRoute uses an empty FullPath (404) and confirms
// the middleware falls back to the raw URL path without panicking.
func TestTracingMiddleware_UnknownRoute(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Tracing())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
