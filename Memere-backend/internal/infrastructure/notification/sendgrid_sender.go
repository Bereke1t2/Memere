package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// SendGridSender delivers emails via the SendGrid v3 Mail Send API.
// Returns a LogEmailSender when the API key is absent (dev fallback).
type SendGridSender struct {
	apiKey     string
	from       string
	httpClient *http.Client
}

var _ service.EmailSender = (*SendGridSender)(nil)

// NewSendGridSender builds a SendGridSender. apiKey is SENDGRID_API_KEY.
// fromEmail is the verified sender address. Returns LogEmailSender when apiKey
// is empty.
func NewSendGridSender(apiKey, fromEmail string) service.EmailSender {
	if apiKey == "" {
		return LogEmailSender{}
	}
	return &SendGridSender{
		apiKey:     apiKey,
		from:       fromEmail,
		httpClient: &http.Client{},
	}
}

type sgPayload struct {
	Personalizations []sgPersonalization `json:"personalizations"`
	From             sgEmail             `json:"from"`
	Subject          string              `json:"subject"`
	Content          []sgContent         `json:"content"`
}

type sgPersonalization struct {
	To []sgEmail `json:"to"`
}

type sgEmail struct {
	Email string `json:"email"`
}

type sgContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (s *SendGridSender) Send(ctx context.Context, to, subject, html string) error {
	payload := sgPayload{
		Personalizations: []sgPersonalization{{To: []sgEmail{{Email: to}}}},
		From:             sgEmail{Email: s.from},
		Subject:          subject,
		Content:          []sgContent{{Type: "text/html", Value: html}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sendgrid: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sendgrid: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sendgrid: send: %w", err)
	}
	defer resp.Body.Close()

	// 202 Accepted is the success response for SendGrid.
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("sendgrid: HTTP %d", resp.StatusCode)
	}
	return nil
}
