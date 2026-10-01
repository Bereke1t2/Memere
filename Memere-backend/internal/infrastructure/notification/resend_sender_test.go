package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResendSender_EmptyKeyReturnsLogSender(t *testing.T) {
	sender := NewResendSender("", "onboarding@resend.dev", "Mirkuz")
	if _, ok := sender.(LogEmailSender); !ok {
		t.Fatalf("expected LogEmailSender when apiKey is empty, got %T", sender)
	}
}

func TestResendSender_SendSuccess(t *testing.T) {
	var receivedAuthHeader string
	var receivedContentType string
	var receivedPayload resendPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		receivedContentType = r.Header.Get("Content-Type")

		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"49633864-637a-4dec-a4a0-781d61b05e63"}`))
	}))
	defer srv.Close()

	sender := &ResendSender{
		apiKey:     "re_123456789_abcdefg",
		fromEmail:  "updates@mirkuz.et",
		fromName:   "Mirkuz Platform",
		apiURL:     srv.URL,
		httpClient: srv.Client(),
	}

	err := sender.Send(context.Background(), "student@example.com", "Verify Email", "<p>Code: 123456</p>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedAuthHeader != "Bearer re_123456789_abcdefg" {
		t.Errorf("expected Authorization 'Bearer re_123456789_abcdefg', got %q", receivedAuthHeader)
	}
	if receivedContentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", receivedContentType)
	}
	if receivedPayload.From != "Mirkuz Platform <updates@mirkuz.et>" {
		t.Errorf("expected from 'Mirkuz Platform <updates@mirkuz.et>', got %q", receivedPayload.From)
	}
	if len(receivedPayload.To) != 1 || receivedPayload.To[0] != "student@example.com" {
		t.Errorf("unexpected to recipient: %+v", receivedPayload.To)
	}
	if receivedPayload.Subject != "Verify Email" {
		t.Errorf("expected subject 'Verify Email', got %q", receivedPayload.Subject)
	}
	if receivedPayload.HTML != "<p>Code: 123456</p>" {
		t.Errorf("expected HTML '<p>Code: 123456</p>', got %q", receivedPayload.HTML)
	}
}

func TestResendSender_SendErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"statusCode":422,"message":"Domain not verified","name":"validation_error"}`))
	}))
	defer srv.Close()

	sender := &ResendSender{
		apiKey:     "invalid-key",
		fromEmail:  "noreply@mirkuz.et",
		fromName:   "Mirkuz",
		apiURL:     srv.URL,
		httpClient: srv.Client(),
	}

	err := sender.Send(context.Background(), "student@example.com", "Test", "body")
	if err == nil {
		t.Fatalf("expected error from 422 response, got nil")
	}
	if expected := "resend: HTTP 422"; !containsSubstring(err.Error(), expected) {
		t.Errorf("expected error containing %q, got %q", expected, err.Error())
	}
}
