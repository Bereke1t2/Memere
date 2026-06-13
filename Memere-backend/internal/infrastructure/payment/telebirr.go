package payment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// TelebirrProvider is a sandbox implementation of the payment port (spec §10.1).
// Real Telebirr integration (encrypted RSA payloads) lands with production keys;
// for Phase 4 it returns a sandbox checkout URL and verifies webhooks with the
// same HMAC scheme so the dedup/fulfillment path is exercised end to end.
type TelebirrProvider struct {
	webhookSecret string
	returnBase    string
}

var _ service.PaymentProvider = (*TelebirrProvider)(nil)

// TelebirrConfig configures the sandbox Telebirr client.
type TelebirrConfig struct {
	WebhookSecret string
	// CheckoutBaseURL is the sandbox hosted-payment page; the PaymentID is
	// appended as a query param.
	CheckoutBaseURL string
}

// NewTelebirrProvider builds a sandbox TelebirrProvider.
func NewTelebirrProvider(cfg TelebirrConfig) *TelebirrProvider {
	base := cfg.CheckoutBaseURL
	if base == "" {
		base = "https://sandbox.telebirr.local/checkout"
	}
	return &TelebirrProvider{webhookSecret: cfg.WebhookSecret, returnBase: base}
}

// Name returns the provider's entity name.
func (t *TelebirrProvider) Name() string { return string(entity.ProviderTelebirr) }

// CreateCheckout returns a sandbox hosted-payment URL keyed on the PaymentID.
func (t *TelebirrProvider) CreateCheckout(ctx context.Context, req service.CheckoutRequest) (*service.CheckoutResult, error) {
	return &service.CheckoutResult{
		CheckoutID:  req.PaymentID,
		RedirectURL: t.returnBase + "?tx_ref=" + req.PaymentID,
	}, nil
}

// telebirrWebhook is the sandbox webhook body shape.
type telebirrWebhook struct {
	EventID string `json:"event_id"`
	Type    string `json:"type"`
	TxRef   string `json:"tx_ref"`
	Status  string `json:"status"`
	TxnID   string `json:"transaction_id"`
}

// VerifyWebhook validates the HMAC signature and parses the sandbox event.
func (t *TelebirrProvider) VerifyWebhook(ctx context.Context, headers map[string]string, body []byte) (*service.WebhookEvent, error) {
	sig := headerValue(headers, "Telebirr-Signature", "x-telebirr-signature")
	if !verifyHexSig(t.webhookSecret, sig, body) {
		return nil, apperror.New(http.StatusUnauthorized, "BAD_SIGNATURE", "webhook signature verification failed", nil)
	}
	var wh telebirrWebhook
	if err := json.Unmarshal(body, &wh); err != nil {
		return nil, apperror.New(http.StatusBadRequest, "BAD_WEBHOOK", "could not parse webhook body", err)
	}
	eventID := wh.EventID
	if eventID == "" {
		eventID = wh.TxRef + ":" + wh.Status
	}
	return &service.WebhookEvent{
		ProviderEventID: eventID,
		Type:            wh.Type,
		PaymentRef:      wh.TxRef,
		Status:          normalizeSandboxStatus(wh.Status),
		ProviderTxnID:   wh.TxnID,
		RawPayload:      body,
	}, nil
}

// VerifyPayment is unsupported in the sandbox; callers treat a pending result as
// "no information yet".
func (t *TelebirrProvider) VerifyPayment(ctx context.Context, providerRef string) (string, error) {
	return string(entity.PayPending), nil
}

// normalizeSandboxStatus maps sandbox status strings onto entity values; shared
// by the sandbox providers.
func normalizeSandboxStatus(s string) string {
	switch s {
	case "success", "successful", "completed", "paid":
		return string(entity.PayCompleted)
	case "failed", "cancelled", "canceled", "declined":
		return string(entity.PayFailed)
	default:
		return s
	}
}
