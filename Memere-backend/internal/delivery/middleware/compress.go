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
// exceeded it flushes through gzip. Responses whose Content-Type is already
// compressed (images, video, audio, PDF, archives) are passed through untouched
// so byte-range streaming and Content-Length survive.
type gzipResponseWriter struct {
	gin.ResponseWriter
	gz      *gzip.Writer
	buf     []byte
	active  bool
	decided bool // whether the compress/passthrough choice has been made
	skip    bool // true => never compress, write straight through
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.skip {
		return w.ResponseWriter.Write(b)
	}
	if w.active {
		return w.gz.Write(b)
	}
	// Decide once, on the first write, using the Content-Type the handler set.
	// Precompressed/binary payloads (notably the /media video proxy) must not be
	// gzipped: activate() deletes Content-Length, which would break HTTP Range.
	if !w.decided {
		w.decided = true
		if isIncompressibleType(w.ResponseWriter.Header().Get("Content-Type")) {
			w.skip = true
			return w.ResponseWriter.Write(b)
		}
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

// isIncompressibleType reports whether a Content-Type is already compressed (or
// otherwise must not be gzipped). Gzipping these wastes CPU and — because the
// gzip path deletes Content-Length — would break HTTP Range/seeking for the
// media proxy. Parameters (e.g. "; charset=utf-8") are ignored.
func isIncompressibleType(ct string) bool {
	if ct == "" {
		return false
	}
	ct = strings.ToLower(strings.TrimSpace(ct))
	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = strings.TrimSpace(ct[:i])
	}
	if strings.HasPrefix(ct, "image/") ||
		strings.HasPrefix(ct, "video/") ||
		strings.HasPrefix(ct, "audio/") {
		return true
	}
	switch ct {
	case "application/pdf",
		"application/zip",
		"application/gzip",
		"application/x-gzip",
		"application/octet-stream",
		"application/x-7z-compressed",
		"application/x-rar-compressed":
		return true
	}
	return false
}

// noopWriter satisfies the http.ResponseWriter interface for pool initialization.
type noopWriter struct{ http.ResponseWriter }

func (n *noopWriter) Write(b []byte) (int, error) { return len(b), nil }
