package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// GoogleCredentialStore is the subset of the Drive credential store this handler
// needs: it persists the single admin account's OAuth refresh token (encrypted
// at rest, server-side only) and reports whether an account is connected. The
// refresh token itself is written by Callback and NEVER read back into a
// response — Load is used only for its ok flag in Status.
type GoogleCredentialStore interface {
	Save(ctx context.Context, refreshToken, googleEmail string) error
	Load(ctx context.Context) (refreshToken string, ok bool, err error)
}

// GoogleOAuthConfig carries the client credentials + redirect for the one-time
// admin connect flow. main.go maps config.GoogleDriveConfig onto it.
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// GoogleOAuthHandler implements the one-time admin flow that connects the single
// admin-owned Google Drive account: Connect returns a Google consent URL, and
// Callback exchanges the returned code for a refresh token and stores it
// encrypted. Students and teachers NEVER touch this flow. The refresh token,
// client secret and access tokens are never returned to any client, never logged,
// and never rendered in the result page.
type GoogleOAuthHandler struct {
	cfg   GoogleOAuthConfig
	creds GoogleCredentialStore
	redis *redis.Client
	hc    *http.Client
}

// NewGoogleOAuthHandler wires the handler. redis backs the one-time CSRF state.
func NewGoogleOAuthHandler(cfg GoogleOAuthConfig, creds GoogleCredentialStore, rc *redis.Client) *GoogleOAuthHandler {
	return &GoogleOAuthHandler{
		cfg:   cfg,
		creds: creds,
		redis: rc,
		hc:    &http.Client{Timeout: 15 * time.Second},
	}
}

const (
	oauthStatePrefix = "gdrive_oauth_state:"
	oauthStateTTL    = 10 * time.Minute
	// Drive scope: full drive access to the admin's own dedicated account so the
	// backend can create files inside, read, and delete them. openid+email let us
	// record which account was connected (best-effort, for the admin UI). No scope
	// is ever requested from students/teachers.
	googleDriveScope = "https://www.googleapis.com/auth/drive openid email"
)

// Connect handles GET /api/v1/admin/google/connect (admin only, bearer). It mints
// a one-time CSRF state, stores it in Redis, and returns the Google consent URL
// for the admin to open. It does NOT redirect: the admin panel calls this with
// its bearer token and opens auth_url in a browser tab, so the public Callback
// receives the code without needing the token.
func (h *GoogleOAuthHandler) Connect(c *gin.Context) {
	if h.cfg.ClientID == "" || h.cfg.ClientSecret == "" {
		respondError(c, apperror.New(http.StatusServiceUnavailable, "GDRIVE_NOT_CONFIGURED",
			"Google Drive OAuth is not configured (set GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET)", nil))
		return
	}
	state, err := randomState()
	if err != nil {
		respondError(c, apperror.Internal(err))
		return
	}
	if err := h.redis.Set(c.Request.Context(), oauthStatePrefix+state, "1", oauthStateTTL).Err(); err != nil {
		respondError(c, apperror.Internal(err))
		return
	}

	q := url.Values{}
	q.Set("client_id", h.cfg.ClientID)
	q.Set("redirect_uri", h.cfg.RedirectURI)
	q.Set("response_type", "code")
	q.Set("scope", googleDriveScope)
	q.Set("access_type", "offline") // required to receive a refresh token
	q.Set("prompt", "consent")      // force a refresh token even on re-consent
	q.Set("include_granted_scopes", "true")
	q.Set("state", state)
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()

	respondJSON(c, http.StatusOK, gin.H{"auth_url": authURL})
}

// Status handles GET /api/v1/admin/google/status (admin only): reports whether a
// Drive account is currently connected. It never returns the refresh token.
func (h *GoogleOAuthHandler) Status(c *gin.Context) {
	_, ok, err := h.creds.Load(c.Request.Context())
	if err != nil {
		respondError(c, apperror.Internal(err))
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"connected": ok})
}

// Callback handles GET /api/v1/admin/google/callback — PUBLIC, because Google
// redirects the admin's browser here with no bearer token. CSRF is enforced by
// the one-time state param (consumed from Redis). It exchanges the code for a
// refresh token, stores it encrypted, and renders a plain success/failure page.
// The token is never displayed, returned, or logged.
func (h *GoogleOAuthHandler) Callback(c *gin.Context) {
	if errParam := c.Query("error"); errParam != "" {
		h.renderResult(c, http.StatusBadRequest, "Google authorization was cancelled or denied.", false)
		return
	}
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		h.renderResult(c, http.StatusBadRequest, "Missing authorization code or state.", false)
		return
	}

	// Validate + consume the one-time state (CSRF). GetDel returns redis.Nil when
	// the key is absent/expired — any error means "not a state we issued".
	if _, err := h.redis.GetDel(c.Request.Context(), oauthStatePrefix+state).Result(); err != nil {
		h.renderResult(c, http.StatusForbidden, "Invalid or expired authorization state. Please start again.", false)
		return
	}

	refresh, accessToken, err := h.exchangeCode(c.Request.Context(), code)
	if err != nil || refresh == "" {
		// No refresh token usually means the account previously consented without
		// prompt=consent; our URL forces it, so this is rare.
		h.renderResult(c, http.StatusBadGateway,
			"Could not complete Google authorization. Ensure you granted offline access, then retry.", false)
		return
	}

	email := h.fetchEmail(c.Request.Context(), accessToken) // best-effort, may be ""
	if err := h.creds.Save(c.Request.Context(), refresh, email); err != nil {
		respondError(c, apperror.Internal(err)) // JSON — this is a server fault, not a user step
		return
	}

	h.renderResult(c, http.StatusOK, "Google Drive connected successfully. You can close this window.", true)
}

// exchangeCode swaps an authorization code for tokens. It returns the refresh
// token and a (short-lived) access token; it never logs the response body.
func (h *GoogleOAuthHandler) exchangeCode(ctx context.Context, code string) (refresh, access string, err error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", h.cfg.ClientID)
	form.Set("client_secret", h.cfg.ClientSecret)
	form.Set("redirect_uri", h.cfg.RedirectURI)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := h.hc.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 2048)) // never log the body
		return "", "", fmt.Errorf("google token exchange: status %d", resp.StatusCode)
	}
	var tr struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", "", err
	}
	return tr.RefreshToken, tr.AccessToken, nil
}

// fetchEmail best-effort resolves which Google account was connected, for display
// in the admin UI. Any failure returns "" — it never blocks the connect flow.
func (h *GoogleOAuthHandler) fetchEmail(ctx context.Context, accessToken string) string {
	if accessToken == "" {
		return ""
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := h.hc.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 2048))
		return ""
	}
	var u struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return ""
	}
	return u.Email
}

// renderResult writes a minimal self-contained HTML page. It contains no secrets.
func (h *GoogleOAuthHandler) renderResult(c *gin.Context, status int, message string, ok bool) {
	title := "Connection failed"
	color := "#b00020"
	if ok {
		title = "Connected"
		color = "#0b8043"
	}
	page := fmt.Sprintf(`<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title></head>
<body style="font-family:system-ui,-apple-system,Segoe UI,Roboto,sans-serif;background:#f7f7f8;margin:0">
<div style="max-width:520px;margin:12vh auto;padding:32px;background:#fff;border-radius:12px;box-shadow:0 1px 4px rgba(0,0,0,.08);text-align:center">
<h1 style="color:%s;font-size:20px;margin:0 0 12px">%s</h1>
<p style="color:#333;font-size:15px;line-height:1.5;margin:0">%s</p>
</div></body></html>`, htmlEscape(title), color, htmlEscape(title), htmlEscape(message))
	c.Data(status, "text/html; charset=utf-8", []byte(page))
}

// randomState mints an unguessable 256-bit URL-safe state token.
func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// htmlEscape minimally escapes text interpolated into the result page.
func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}
