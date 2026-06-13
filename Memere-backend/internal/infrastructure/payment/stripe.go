package payment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// StripeProvider is a sandbox implementation of the payment port (spec §10.1).
// The real Stripe SDK (Checkout Sessions + Stripe-Signature with a timestamped
// scheme) lands with production keys; for Phase 4 it returns a sandbox checkout
// URL and verifies webhooks with the shared HMAC scheme so the dedup/fulfillment
// path is exercised end to end.
type StripeProvider struct {
	webhookSecret string
	checkoutBase  string
}

var _ service.PaymentProvider = (*StripeProvider)(nil)

// StripeConfig configures the sandbox Stripe client.
type StripeConfig struct {
	WebhookSecret string
	CheckoutBase  string
}

// NewStripeProvider builds a sandbox StripeProvider.
func NewStripeProvider(cfg StripeConfig) *StripeProvider {
	base := cfg.CheckoutBase
	if base == "" {
		base = "https://sandbox.stripe.local/checkout"
	}
	return &StripeProvider{webhookSecret: cfg.WebhookSecret, checkoutBase: base}
}

// Name returns the provider's entity name.
func (s *StripeProvider) Name() string { return string(entity.ProviderStripe) }

// CreateCheckout returns a sandbox hosted-checkout URL keyed on the PaymentID.
func (s *StripeProvider) CreateCheckout(ctx context.Context, req service.CheckoutRequest) (*service.CheckoutResult, error) {
	return &service.CheckoutResult{
		CheckoutID:  req.PaymentID,
		RedirectURL: s.checkoutBase + "?client_reference_id=" + req.PaymentID,
	}, nil
}

// stripeWebhook is the sandbox webhook body shape.
type stripeWebhook struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	TxRef  string `json:"client_reference_id"`
	Status string `json:"status"`
	TxnID  string `json:"payment_intent"`
}

// VerifyWebhook validates the HMAC signature and parses the sandbox event.
func (s *StripeProvider) VerifyWebhook(ctx context.Context, headers map[string]string, body []byte) (*service.WebhookEvent, error) {
	sig := headerValue(headers, "Stripe-Signature", "x-stripe-signature")
	if !verifyHexSig(s.webhookSecret, sig, body) {
		return nil, apperror.New(http.StatusUnauthorized, "BAD_SIGNATURE", "webhook signature verification failed", nil)
	}
	var wh stripeWebhook
	if err := json.Unmarshal(body, &wh); err != nil {
		return nil, apperror.New(http.StatusBadRequest, "BAD_WEBHOOK", "could not parse webhook body", err)
	}
	eventID := wh.ID
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

// VerifyPayment is unsupported in the sandbox.
func (s *StripeProvider) VerifyPayment(ctx context.Context, providerRef string) (string, error) {
	return string(entity.PayPending), nil
}
