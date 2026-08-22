package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ErrMediaNotFound is returned by a MediaStore when the requested key is absent.
// The handler maps it to 404 (rather than a 500) so a missing object reads as
// "gone", not "server error".
var ErrMediaNotFound = errors.New("media: object not found")

// MediaContent is one streaming response the media proxy forwards to the client.
// It preserves HTTP Range semantics (StatusCode 206 + Content-Range) so MP4
// playback can seek. The caller must close Body.
type MediaContent struct {
	Body          io.ReadCloser
	StatusCode    int    // 200 (whole) or 206 (range)
	ContentType   string
	ContentLength int64  // -1 when unknown
	ContentRange  string // set on 206
	AcceptRanges  string
}

// MediaStore is the subset of the object store the media proxy needs. Keeping it
// a local interface (satisfied in wiring by an adapter over the Google Drive
// store) means this delivery package never imports the storage/infrastructure
// ring — the clean-architecture dependency rule holds.
type MediaStore interface {
	// VerifyMediaURL reports whether sig authenticates key+exp and exp is fresh.
	VerifyMediaURL(key string, exp int64, sig string) bool
	// Stream fetches key from the backing store, honoring an optional HTTP Range
	// header. It returns ErrMediaNotFound when the key is unknown.
	Stream(ctx context.Context, key, rangeHeader string) (*MediaContent, error)
}

// MediaHandler serves object bytes over short-lived, HMAC-signed proxy URLs used
// when STORAGE_PROVIDER=gdrive. The signature (minted by the store's PresignGet
// and carried in ?exp=&sig=) is the ONLY credential the client presents: Google
// Drive is never exposed publicly and no Google OAuth token ever reaches the
// client. App-level authorization is enforced upstream — a caller only ever
// receives a signed /media URL after the video usecase's access check passes, so
// the short expiry bounds replay of an already-authorized link.
type MediaHandler struct {
	store MediaStore
}

// NewMediaHandler builds a MediaHandler over the signed-URL-aware store.
func NewMediaHandler(store MediaStore) *MediaHandler {
	return &MediaHandler{store: store}
}

// Serve handles GET /media?key=&exp=&sig= — it verifies the signature and expiry
// then streams the object, forwarding Range/Content-Range/Accept-Ranges so the
// mobile player can seek within an MP4. On a bad/expired signature it returns 403
// without touching the store; a missing object is 404; an upstream Drive fault is
// 502. No error body leaks Drive detail.
func (h *MediaHandler) Serve(c *gin.Context) {
	key := c.Query("key")
	expStr := c.Query("exp")
	sig := c.Query("sig")
	if key == "" || expStr == "" || sig == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if !h.store.VerifyMediaURL(key, exp, sig) {
		c.Status(http.StatusForbidden)
		return
	}

	content, err := h.store.Stream(c.Request.Context(), key, c.GetHeader("Range"))
	if err != nil {
		if errors.Is(err, ErrMediaNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		// Upstream (Drive) problem — let the client retry, without leaking detail.
		c.Status(http.StatusBadGateway)
		return
	}
	defer content.Body.Close()

	hdr := c.Writer.Header()
	if content.ContentType != "" {
		hdr.Set("Content-Type", content.ContentType)
	}
	if content.AcceptRanges != "" {
		hdr.Set("Accept-Ranges", content.AcceptRanges)
	} else {
		hdr.Set("Accept-Ranges", "bytes")
	}
	if content.ContentRange != "" {
		hdr.Set("Content-Range", content.ContentRange)
	}
	if content.ContentLength >= 0 {
		hdr.Set("Content-Length", strconv.FormatInt(content.ContentLength, 10))
	}

	status := content.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	c.Status(status)

	// Stream the body straight through. The Compress middleware skips binary
	// content types (see middleware/compress.go), so Content-Length/Range are not
	// rewritten out from under us here.
	_, _ = io.Copy(c.Writer, content.Body)
}
