package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// BrevoSender delivers transactional emails via Brevo (Sendinblue) v3 REST API.
// Returns a LogEmailSender when the API key is absent (dev fallback).
type BrevoSender struct {
	apiKey     string
	fromEmail  string
	fromName   string
	apiURL     string
	httpClient *http.Client
}

var _ service.EmailSender = (*BrevoSender)(nil)

// NewBrevoSender builds a BrevoSender. apiKey is BREVO_API_KEY.
// fromEmail and fromName are the sender identity. Returns LogEmailSender when
// apiKey is empty.
func NewBrevoSender(apiKey, fromEmail, fromName string) service.EmailSender {
	if apiKey == "" {
		return LogEmailSender{}
	}
	if fromName == "" {
		fromName = "Memere"
	}
	if fromEmail == "" {
		fromEmail = "noreply@memere.app"
	}
	return &BrevoSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		apiURL:    "https://api.brevo.com/v3/smtp/email",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type brevoContact struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoPayload struct {
	Sender      brevoContact   `json:"sender"`
	To          []brevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent"`
}

func (b *BrevoSender) Send(ctx context.Context, to, subject, html string) error {
	payload := brevoPayload{
		Sender: brevoContact{
			Email: b.fromEmail,
			Name:  b.fromName,
		},
		To: []brevoContact{
			{Email: to},
		},
		Subject:     subject,
		HTMLContent: html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("brevo: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("brevo: build request: %w", err)
	}
	req.Header.Set("api-key", b.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("brevo: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
