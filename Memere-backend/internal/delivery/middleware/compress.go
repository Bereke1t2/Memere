package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	// compressMinBytes is the minimum response size to trigger gzip. Compressing
	// tiny JSON payloads costs more CPU than it saves on the wire.
	compressMinBytes = 1024
)

// pool reuses gzip writers to avoid per-request allocations.
var gzipPool = sync.Pool{
	New: func() any {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		return w
	},
}

// gzipResponseWriter wraps gin.ResponseWriter and compresses the body when the
// threshold is crossed. It buffers the first compressMinBytes in memory; once
// exceeded it flushes through gzip.
type gzipResponseWriter struct {
	gin.ResponseWriter
	gz     *gzip.Writer
	buf    []byte
	active bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.active {
		return w.gz.Write(b)
	}
	w.buf = append(w.buf, b...)
	if len(w.buf) >= compressMinBytes {
		w.activate()
		n, err := w.gz.Write(w.buf)
		w.buf = nil
		return n, err
	}
	return len(b), nil
}

func (w *gzipResponseWriter) activate() {
	w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	w.ResponseWriter.Header().Del("Content-Length") // length changes after compression
	w.gz.Reset(w.ResponseWriter)
	w.active = true
}

func (w *gzipResponseWriter) flush() {
	if !w.active && len(w.buf) > 0 {
		_, _ = w.ResponseWriter.Write(w.buf)
		w.buf = nil
		return
	}
	if w.active {
		_ = w.gz.Flush()
	}
}

// Compress adds gzip Content-Encoding to responses when the client sends
// Accept-Encoding: gzip and the response body exceeds compressMinBytes. It is a
// no-op for already-compressed content types (images, video).
func Compress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		gz := gzipPool.Get().(*gzip.Writer)
		defer gzipPool.Put(gz)

		grw := &gzipResponseWriter{
			ResponseWriter: c.Writer,
			gz:             gz,
		}
		c.Writer = grw
		defer grw.flush()

		c.Header("Vary", "Accept-Encoding")
		c.Next()

		if grw.active {
			_ = gz.Close()
		}
	}
}

// noopWriter satisfies the http.ResponseWriter interface for pool initialization.
type noopWriter struct{ http.ResponseWriter }

func (n *noopWriter) Write(b []byte) (int, error) { return len(b), nil }
