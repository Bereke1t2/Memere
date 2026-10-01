package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewBrevoSender_EmptyKeyReturnsLogSender(t *testing.T) {
	sender := NewBrevoSender("", "noreply@example.com", "Mirkuz")
	if _, ok := sender.(LogEmailSender); !ok {
		t.Fatalf("expected LogEmailSender when apiKey is empty, got %T", sender)
	}
}

func TestBrevoSender_SendSuccess(t *testing.T) {
	var receivedHeaderKey string
	var receivedContentType string
	var receivedPayload brevoPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaderKey = r.Header.Get("api-key")
		receivedContentType = r.Header.Get("Content-Type")

		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"messageId":"<test-msg-id>"}`))
	}))
	defer srv.Close()

	sender := &BrevoSender{
		apiKey:     "xkeysib-testkey-12345",
		fromEmail:  "verified@mirkuz.et",
		fromName:   "Mirkuz Platform",
		apiURL:     srv.URL,
		httpClient: srv.Client(),
	}

	err := sender.Send(context.Background(), "student@example.com", "Verify Email", "<p>Code: 123456</p>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedHeaderKey != "xkeysib-testkey-12345" {
		t.Errorf("expected api-key 'xkeysib-testkey-12345', got %q", receivedHeaderKey)
	}
	if receivedContentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", receivedContentType)
	}
	if receivedPayload.Sender.Email != "verified@mirkuz.et" {
		t.Errorf("expected sender email 'verified@mirkuz.et', got %q", receivedPayload.Sender.Email)
	}
	if receivedPayload.Sender.Name != "Mirkuz Platform" {
		t.Errorf("expected sender name 'Mirkuz Platform', got %q", receivedPayload.Sender.Name)
	}
	if len(receivedPayload.To) != 1 || receivedPayload.To[0].Email != "student@example.com" {
		t.Errorf("unexpected to recipient: %+v", receivedPayload.To)
	}
	if receivedPayload.Subject != "Verify Email" {
		t.Errorf("expected subject 'Verify Email', got %q", receivedPayload.Subject)
	}
	if receivedPayload.HTMLContent != "<p>Code: 123456</p>" {
		t.Errorf("expected HTML content '<p>Code: 123456</p>', got %q", receivedPayload.HTMLContent)
	}
}

func TestBrevoSender_SendErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":"unauthorized","message":"Key not found"}`))
	}))
	defer srv.Close()

	sender := &BrevoSender{
		apiKey:     "invalid-key",
		fromEmail:  "noreply@mirkuz.et",
		fromName:   "Mirkuz",
		apiURL:     srv.URL,
		httpClient: srv.Client(),
	}

	err := sender.Send(context.Background(), "student@example.com", "Test", "body")
	if err == nil {
		t.Fatalf("expected error from 401 response, got nil")
	}
	if expected := "brevo: HTTP 401"; !containsSubstring(err.Error(), expected) {
		t.Errorf("expected error containing %q, got %q", expected, err.Error())
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && searchSubstr(s, substr)))
}

func searchSubstr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
