package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// FCMSender delivers push notifications via the Firebase Cloud Messaging
// legacy HTTP API. Tokens reported as NotRegistered are returned as invalid
// so the caller (NotificationWorker) can prune them from device_tokens.
type FCMSender struct {
	serverKey  string
	httpClient *http.Client
}

var _ service.PushSender = (*FCMSender)(nil)

// NewFCMSender builds an FCMSender. serverKey is the FCM_SERVER_KEY env value.
// Returns a LogPushSender when serverKey is empty (dev fallback).
func NewFCMSender(serverKey string) service.PushSender {
	if serverKey == "" {
		return LogPushSender{}
	}
	return &FCMSender{serverKey: serverKey, httpClient: &http.Client{}}
}

type fcmPayload struct {
	RegistrationIDs []string          `json:"registration_ids"`
	Notification    fcmNotification   `json:"notification"`
	Data            map[string]string `json:"data,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type fcmResponse struct {
	Results []fcmResult `json:"results"`
}

type fcmResult struct {
	Error string `json:"error"`
}

func (s *FCMSender) Send(ctx context.Context, tokens []string, ev service.NotifyEvent) ([]string, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	payload := fcmPayload{
		RegistrationIDs: tokens,
		Notification:    fcmNotification{Title: ev.Title, Body: ev.Body},
		Data:            ev.Data,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("fcm: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://fcm.googleapis.com/fcm/send", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("fcm: build request: %w", err)
	}
	req.Header.Set("Authorization", "key="+s.serverKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fcm: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fcm: HTTP %d", resp.StatusCode)
	}

	var fcmResp fcmResponse
	if err := json.NewDecoder(resp.Body).Decode(&fcmResp); err != nil {
		return nil, fmt.Errorf("fcm: decode response: %w", err)
	}

	var invalid []string
	for i, r := range fcmResp.Results {
		if r.Error == "NotRegistered" || r.Error == "InvalidRegistration" {
			if i < len(tokens) {
				invalid = append(invalid, tokens[i])
			}
		}
	}
	return invalid, nil
}
