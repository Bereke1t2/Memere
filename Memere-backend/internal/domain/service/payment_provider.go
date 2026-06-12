package service

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// CheckoutRequest is the provider-agnostic input to start a hosted checkout. The
// usecase fills it from a pending payment row; the provider turns it into a
// redirect URL the client opens. Amount is decimal money (never float).
type CheckoutRequest struct {
	PaymentID      string
	Amount         decimal.Decimal
	Currency       string
	CustomerEmail  string
	Description    string
	ReturnURL      string // where the provider redirects the user after pay
	CallbackURL    string // our webhook endpoint the provider posts to
	IdempotencyKey string
}

// CheckoutResult is the provider's response to CreateCheckout: a reference we
// store as provider_checkout_id and a URL the client opens (Stripe may return a
// client_secret here instead).
type CheckoutResult struct {
	CheckoutID  string
	RedirectURL string
}

// WebhookEvent is a normalized provider webhook after signature verification.
// PaymentRef maps back to our PaymentID or the provider transaction id; Status
// is normalized to entity.PaymentStatus values ("completed" | "failed").
type WebhookEvent struct {
	ProviderEventID string
	Type            string
	PaymentRef      string
	Status          string
	ProviderTxnID   string
	RawPayload      []byte
}

// PaymentProvider is the domain PORT for a payment backend. Chapa (primary),
// Telebirr and Stripe implement it in the infrastructure ring (spec §10.1); the
// usecase depends only on this interface so provider churn never reaches the
// business logic (a named risk in §14.4).
type PaymentProvider interface {
	// Name returns the provider's entity.PaymentProvider value as a string.
	Name() string
	// CreateCheckout starts a hosted payment and returns the redirect URL.
	CreateCheckout(ctx context.Context, req CheckoutRequest) (*CheckoutResult, error)
	// VerifyWebhook validates the signature on a raw webhook body and parses it
	// into a normalized event. It MUST reject an unsigned or tampered payload.
	VerifyWebhook(ctx context.Context, headers map[string]string, body []byte) (*WebhookEvent, error)
	// VerifyPayment is a server-to-server confirmation used as a polling fallback
	// when a webhook is delayed; it returns the normalized status string.
	VerifyPayment(ctx context.Context, providerRef string) (status string, err error)
}

// PaymentProviderRegistry resolves a configured provider by its entity name. The
// usecase asks the registry for the provider named on the payment row, so adding
// a provider is a wiring change, not a usecase change.
type PaymentProviderRegistry interface {
	Get(p entity.PaymentProvider) (PaymentProvider, error)
}
