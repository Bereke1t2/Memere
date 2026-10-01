package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// ResendSender delivers transactional emails via the Resend REST API (https://resend.com).
// Returns a LogEmailSender when the API key is absent (dev fallback).
type ResendSender struct {
	apiKey     string
	fromEmail  string
	fromName   string
	apiURL     string
	httpClient *http.Client
}

var _ service.EmailSender = (*ResendSender)(nil)

// NewResendSender builds a ResendSender. apiKey is RESEND_API_KEY.
// fromEmail and fromName define the sender identity. Returns LogEmailSender when
// apiKey is empty.
func NewResendSender(apiKey, fromEmail, fromName string) service.EmailSender {
	if apiKey == "" {
		return LogEmailSender{}
	}
	if fromName == "" {
		fromName = "Memere"
	}
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}
	return &ResendSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		apiURL:    "https://api.resend.com/emails",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (r *ResendSender) formatFrom() string {
	if r.fromName != "" {
		return fmt.Sprintf("%s <%s>", r.fromName, r.fromEmail)
	}
	return r.fromEmail
}

func (r *ResendSender) Send(ctx context.Context, to, subject, html string) error {
	payload := resendPayload{
		From:    r.formatFrom(),
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("resend: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
