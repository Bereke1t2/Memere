package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
)

// --- SecurityHeaders ---

func TestSecurityHeaders_PresentOnEveryResponse(t *testing.T) {
	r := gin.New()
	r.Use(middleware.SecurityHeaders())
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))

	want := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Referrer-Policy":           "no-referrer",
		"Content-Security-Policy":   "default-src 'none'; frame-ancestors 'none'",
		"Strict-Transport-Security": "max-age=63072000; includeSubDomains; preload",
	}
	for h, v := range want {
		if got := w.Header().Get(h); got != v {
			t.Errorf("header %s: want %q, got %q", h, v, got)
		}
	}
}

// --- BodyLimit ---

func TestBodyLimit_AcceptsSmallBody(t *testing.T) {
	r := gin.New()
	r.Use(middleware.BodyLimit(1024))
	r.POST("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bytes.Repeat([]byte("x"), 100)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBodyLimit_RejectsOversizedContentLength(t *testing.T) {
	r := gin.New()
	r.Use(middleware.BodyLimit(512))
	r.POST("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bytes.Repeat([]byte("x"), 1000)))
	req.ContentLength = 1000
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Code)
	}
}

// --- CORS ---

func TestCORS_AllowedOrigin(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS([]string{"https://app.memere.et"}, false))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://app.memere.et")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.memere.et" {
		t.Errorf("ACAO header: got %q", got)
	}
}

func TestCORS_UnknownOriginOptions_Rejected(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS([]string{"https://app.memere.et"}, true))
	r.OPTIONS("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("ACAO header must not be set for disallowed origin")
	}
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestCORS_WildcardInProduction_FailsClosed(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS([]string{"*"}, true))
	r.OPTIONS("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://attacker.com")
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("wildcard in production must not emit ACAO header")
	}
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for wildcard+production preflight, got %d", w.Code)
	}
}

func TestCORS_WildcardInDev_AllowsAny(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS([]string{"*"}, false))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected '*', got %q", got)
	}
}
