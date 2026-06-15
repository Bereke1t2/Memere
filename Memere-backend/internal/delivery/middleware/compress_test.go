package middleware_test

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
)

func TestCompress_GzipsLargeResponse(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Compress())
	// Return a body large enough to cross the 1 KiB threshold.
	body := strings.Repeat("hello world ", 200) // ~2.4 KiB
	r.GET("/big", func(c *gin.Context) {
		c.String(http.StatusOK, body)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/big", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if enc := w.Header().Get("Content-Encoding"); enc != "gzip" {
		t.Fatalf("expected Content-Encoding: gzip, got %q", enc)
	}

	gz, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatalf("gzip.NewReader: %v", err)
	}
	defer gz.Close()
	decoded, err := io.ReadAll(gz)
	if err != nil {
		t.Fatalf("read gzip: %v", err)
	}
	if string(decoded) != body {
		t.Errorf("decoded body mismatch: got len %d want len %d", len(decoded), len(body))
	}
}

func TestCompress_SkipsSmallResponse(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Compress())
	r.GET("/small", func(c *gin.Context) {
		c.String(http.StatusOK, "hi")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/small", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	r.ServeHTTP(w, req)

	if enc := w.Header().Get("Content-Encoding"); enc == "gzip" {
		t.Error("small response should not be gzip-encoded")
	}
	if w.Body.String() != "hi" {
		t.Errorf("body = %q, want %q", w.Body.String(), "hi")
	}
}

func TestCompress_NoGzipWithoutHeader(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Compress())
	body := strings.Repeat("x", 2000)
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, body) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil) // no Accept-Encoding
	r.ServeHTTP(w, req)

	if enc := w.Header().Get("Content-Encoding"); enc == "gzip" {
		t.Error("should not compress when Accept-Encoding: gzip is absent")
	}
}
