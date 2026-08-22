package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// ErrObjectNotFound is returned by the Drive store when a key is not present in
// the object index (or Drive reports the file gone). The media proxy maps it to
// a 404 rather than a 500.
var ErrObjectNotFound = errors.New("storage: object not found")

// GoogleDriveConfig carries what GoogleDriveStore needs. main.go maps the typed
// config onto this so infrastructure stays free of the config package.
type GoogleDriveConfig struct {
	ClientID     string
	ClientSecret string
	RootFolderID string
	// RefreshToken is an optional env-provided fallback; the store prefers the
	// CredentialStore (DB) value obtained via the admin connect page.
	RefreshToken string
	// MediaBaseURL is the absolute URL of the media proxy that streams objects,
	// e.g. "http://localhost:8080/api/v1/media". PresignGet returns signed links
	// against it.
	MediaBaseURL string
	// SignSecret keys the HMAC that protects media-proxy URLs (derived from the
	// app JWT secret in wiring). Must be non-empty.
	SignSecret string
}

// GoogleDriveStore implements service.ObjectStore over a single admin-owned
// Google Drive account using the Drive v3 REST API directly (no SDK dependency).
// A DB-backed ObjectIndex maps stable keys to opaque Drive file ids; a
// CredentialStore holds the OAuth refresh token used to mint short-lived access
// tokens. No Google credential is ever exposed to a client: downloads either
// stream through the backend (Get) or via short-lived HMAC-signed proxy URLs
// (PresignGet -> /media), never through Drive's public sharing.
type GoogleDriveStore struct {
	cfg   GoogleDriveConfig
	index ObjectIndex
	creds CredentialStore
	hc    *http.Client
	now   func() time.Time

	mediaKey []byte // HMAC key for media-proxy URLs (derived from SignSecret)

	mu          sync.Mutex // guards the cached access token
	accessToken string
	tokenExp    time.Time
}

var (
	_ service.ObjectStore   = (*GoogleDriveStore)(nil)
	_ service.SizedUploader = (*GoogleDriveStore)(nil)
)

// NewGoogleDriveStore builds the store. It does NOT require a refresh token at
// construction (the admin may connect later); calls fail with a clear error
// until the account is connected.
func NewGoogleDriveStore(cfg GoogleDriveConfig, index ObjectIndex, creds CredentialStore) (*GoogleDriveStore, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, errors.New("storage: gdrive requires GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET")
	}
	if cfg.RootFolderID == "" {
		return nil, errors.New("storage: gdrive requires GOOGLE_DRIVE_ROOT_FOLDER_ID")
	}
	if cfg.SignSecret == "" {
		return nil, errors.New("storage: gdrive requires a media URL signing secret")
	}
	if index == nil {
		return nil, errors.New("storage: gdrive requires an object index")
	}
	mk := hmac.New(sha256.New, []byte(cfg.SignSecret))
	mk.Write([]byte("memere-media-url-v1"))
	return &GoogleDriveStore{
		cfg:   cfg,
		index: index,
		creds: creds,
		// No client-level timeout: per-request contexts govern cancellation, and a
		// fixed deadline would abort legitimate large uploads/streams mid-transfer.
		hc:       &http.Client{},
		now:      time.Now,
		mediaKey: mk.Sum(nil),
	}, nil
}

// ---- OAuth access-token management ----------------------------------------

// refreshToken returns the admin account's OAuth refresh token, preferring the
// server-side credential store (DB) and falling back to the configured env var.
func (s *GoogleDriveStore) refreshToken(ctx context.Context) (string, error) {
	if s.creds != nil {
		tok, ok, err := s.creds.Load(ctx)
		if err != nil {
			return "", err
		}
		if ok && tok != "" {
			return tok, nil
		}
	}
	if s.cfg.RefreshToken != "" {
		return s.cfg.RefreshToken, nil
	}
	return "", errors.New("storage: gdrive not connected — no refresh token (visit /api/v1/admin/google/connect or set GOOGLE_REFRESH_TOKEN)")
}

// accessTokenFor returns a valid OAuth access token, refreshing via the stored
// refresh token when the cached one is missing or within 60s of expiry.
func (s *GoogleDriveStore) accessTokenFor(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.accessToken != "" && s.now().Before(s.tokenExp.Add(-60*time.Second)) {
		return s.accessToken, nil
	}
	refresh, err := s.refreshToken(ctx)
	if err != nil {
		return "", err
	}
	form := url.Values{}
	form.Set("client_id", s.cfg.ClientID)
	form.Set("client_secret", s.cfg.ClientSecret)
	form.Set("refresh_token", refresh)
	form.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("storage: gdrive token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("storage: gdrive token exchange: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// Never log the body: it can echo client_secret hints. Report status only.
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("storage: gdrive token exchange failed: status %d", resp.StatusCode)
	}
	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("storage: gdrive token decode: %w", err)
	}
	if tr.AccessToken == "" {
		return "", errors.New("storage: gdrive token exchange returned no access_token")
	}
	s.accessToken = tr.AccessToken
	s.tokenExp = s.now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return s.accessToken, nil
}

// ---- Upload ----------------------------------------------------------------

type driveFile struct {
	ID string `json:"id"`
}

// store uploads bytes, records the key -> id mapping, and best-effort deletes a
// prior Drive file for the same key so re-uploads don't orphan objects.
func (s *GoogleDriveStore) store(ctx context.Context, key, contentType string, r io.Reader, size int64) error {
	oldID, hadOld, _ := s.index.Lookup(ctx, key)

	newID, err := s.uploadResumable(ctx, key, contentType, r, size)
	if err != nil {
		return err
	}
	idxSize := size
	if idxSize < 0 {
		idxSize = 0
	}
	if err := s.index.Put(ctx, key, newID, contentType, idxSize); err != nil {
		// The upload landed but we couldn't record it: delete the new file so it
		// doesn't leak, then surface the error.
		_ = s.deleteDriveFile(ctx, newID)
		return err
	}
	if hadOld && oldID != "" && oldID != newID {
		_ = s.deleteDriveFile(ctx, oldID) // best-effort; a leftover is harmless
	}
	return nil
}

// uploadResumable creates a file in the root folder via a resumable session and
// returns its Drive file id. size must be the exact byte count r will yield.
func (s *GoogleDriveStore) uploadResumable(ctx context.Context, key, contentType string, r io.Reader, size int64) (string, error) {
	token, err := s.accessTokenFor(ctx)
	if err != nil {
		return "", err
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 1. Initiate the resumable session.
	meta, _ := json.Marshal(map[string]any{
		"name":          key,
		"parents":       []string{s.cfg.RootFolderID},
		"appProperties": map[string]string{"key": key},
	})
	const initURL = "https://www.googleapis.com/upload/drive/v3/files?uploadType=resumable&fields=id&supportsAllDrives=true"
	initReq, err := http.NewRequestWithContext(ctx, http.MethodPost, initURL, bytes.NewReader(meta))
	if err != nil {
		return "", fmt.Errorf("storage: gdrive init upload: %w", err)
	}
	initReq.Header.Set("Authorization", "Bearer "+token)
	initReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	initReq.Header.Set("X-Upload-Content-Type", contentType)
	if size >= 0 {
		initReq.Header.Set("X-Upload-Content-Length", strconv.FormatInt(size, 10))
	}
	initResp, err := s.hc.Do(initReq)
	if err != nil {
		return "", fmt.Errorf("storage: gdrive init upload: %w", err)
	}
	sessionURI := initResp.Header.Get("Location")
	_, _ = io.Copy(io.Discard, io.LimitReader(initResp.Body, 4096))
	initResp.Body.Close()
	if initResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("storage: gdrive init upload: status %d", initResp.StatusCode)
	}
	if sessionURI == "" {
		return "", errors.New("storage: gdrive init upload: no session URI")
	}

	// 2. Upload the bytes in a single PUT (length known).
	putReq, err := http.NewRequestWithContext(ctx, http.MethodPut, sessionURI, r)
	if err != nil {
		return "", fmt.Errorf("storage: gdrive put: %w", err)
	}
	putReq.Header.Set("Content-Type", contentType)
	if size >= 0 {
		putReq.ContentLength = size // ensures a Content-Length header on a generic reader
	}
	putResp, err := s.hc.Do(putReq)
	if err != nil {
		return "", fmt.Errorf("storage: gdrive put: %w", err)
	}
	defer putResp.Body.Close()
	if putResp.StatusCode != http.StatusOK && putResp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("storage: gdrive put: status %d", putResp.StatusCode)
	}
	var f driveFile
	if err := json.NewDecoder(putResp.Body).Decode(&f); err != nil {
		return "", fmt.Errorf("storage: gdrive put decode: %w", err)
	}
	if f.ID == "" {
		return "", errors.New("storage: gdrive put: empty file id")
	}
	return f.ID, nil
}

// Put uploads bytes of unknown length. Callers of Put are small (lesson PDFs
// ≤50MB, thumbnails, certificate PDFs), so buffering to learn the size is
// acceptable; large videos use PutSized and never buffer here.
func (s *GoogleDriveStore) Put(ctx context.Context, key, contentType string, body io.Reader) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("storage: gdrive read body %q: %w", key, err)
	}
	return s.store(ctx, key, contentType, bytes.NewReader(data), int64(len(data)))
}

// PutSized streams an upload of known length straight to Drive without buffering
// (service.SizedUploader). Used by the proxied video upload.
func (s *GoogleDriveStore) PutSized(ctx context.Context, key, contentType string, r io.Reader, size int64) error {
	return s.store(ctx, key, contentType, r, size)
}

// ---- Download / stream ------------------------------------------------------

// DriveContent is a streaming response from Drive, used by the media proxy to
// forward Range semantics to the client. The caller must close Body.
type DriveContent struct {
	Body          io.ReadCloser
	StatusCode    int    // 200 or 206
	ContentType   string
	ContentLength int64  // -1 if unknown
	ContentRange  string // set on 206
	AcceptRanges  string
}

// download fetches an object's bytes (optionally a byte range) from Drive.
func (s *GoogleDriveStore) download(ctx context.Context, key, rangeHeader string) (*DriveContent, error) {
	fileID, ok, err := s.index.Lookup(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrObjectNotFound
	}
	token, err := s.accessTokenFor(ctx)
	if err != nil {
		return nil, err
	}
	getURL := "https://www.googleapis.com/drive/v3/files/" + url.PathEscape(fileID) + "?alt=media&supportsAllDrives=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("storage: gdrive get %q: %w", key, err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	resp, err := s.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("storage: gdrive get %q: %w", key, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, ErrObjectNotFound
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		resp.Body.Close()
		return nil, fmt.Errorf("storage: gdrive get %q: status %d", key, resp.StatusCode)
	}
	return &DriveContent{
		Body:          resp.Body,
		StatusCode:    resp.StatusCode,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
		ContentRange:  resp.Header.Get("Content-Range"),
		AcceptRanges:  resp.Header.Get("Accept-Ranges"),
	}, nil
}

// Get streams an object (whole). The caller closes the returned reader.
func (s *GoogleDriveStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	dc, err := s.download(ctx, key, "")
	if err != nil {
		return nil, err
	}
	return dc.Body, nil
}

// Stream fetches an object honoring an optional HTTP Range header. Used by the
// media proxy so MP4 playback can seek. The caller closes DriveContent.Body.
func (s *GoogleDriveStore) Stream(ctx context.Context, key, rangeHeader string) (*DriveContent, error) {
	return s.download(ctx, key, rangeHeader)
}

// Exists reports whether the key is present in the object index.
func (s *GoogleDriveStore) Exists(ctx context.Context, key string) (bool, error) {
	_, ok, err := s.index.Lookup(ctx, key)
	return ok, err
}

// Delete removes the Drive file and its index entry. A missing key is not an
// error (matches S3Store semantics).
func (s *GoogleDriveStore) Delete(ctx context.Context, key string) error {
	fileID, ok, err := s.index.Lookup(ctx, key)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if err := s.deleteDriveFile(ctx, fileID); err != nil {
		return err
	}
	return s.index.Delete(ctx, key)
}

func (s *GoogleDriveStore) deleteDriveFile(ctx context.Context, fileID string) error {
	token, err := s.accessTokenFor(ctx)
	if err != nil {
		return err
	}
	delURL := "https://www.googleapis.com/drive/v3/files/" + url.PathEscape(fileID) + "?supportsAllDrives=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return fmt.Errorf("storage: gdrive delete: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.hc.Do(req)
	if err != nil {
		return fmt.Errorf("storage: gdrive delete: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 2048))
	switch resp.StatusCode {
	case http.StatusNoContent, http.StatusOK, http.StatusNotFound:
		return nil
	default:
		return fmt.Errorf("storage: gdrive delete: status %d", resp.StatusCode)
	}
}

// ---- Presign (proxy URLs) ---------------------------------------------------

// PresignPut is unsupported for Drive (no direct client upload). The proxied
// upload endpoint is used instead.
func (s *GoogleDriveStore) PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	return "", errors.New("storage: gdrive does not support presigned uploads; use the proxied upload endpoint")
}

// PresignGet returns a short-lived, HMAC-signed media-proxy URL for key. The
// backend streams the bytes via /media after verifying the signature and expiry;
// Drive is never exposed publicly.
func (s *GoogleDriveStore) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	exp := s.now().Add(ttl).Unix()
	q := url.Values{}
	q.Set("key", key)
	q.Set("exp", strconv.FormatInt(exp, 10))
	q.Set("sig", s.signMedia(key, exp))
	return s.cfg.MediaBaseURL + "?" + q.Encode(), nil
}

func (s *GoogleDriveStore) signMedia(key string, exp int64) string {
	mac := hmac.New(sha256.New, s.mediaKey)
	fmt.Fprintf(mac, "%s\n%d", key, exp)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// VerifyMediaURL reports whether sig authenticates key+exp and exp is not past.
func (s *GoogleDriveStore) VerifyMediaURL(key string, exp int64, sig string) bool {
	if exp < s.now().Unix() {
		return false
	}
	want := s.signMedia(key, exp)
	return hmac.Equal([]byte(want), []byte(sig))
}
