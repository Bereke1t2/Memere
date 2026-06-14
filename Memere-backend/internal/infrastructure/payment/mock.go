package payment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// MockProvider is a TEST-ONLY payment backend (spec §5.6 smoke test). It settles
// entirely offline — CreateCheckout makes no external HTTP call — and verifies a
// simulated webhook with the same HMAC scheme as the real providers, so the
// idempotency, dedup and fulfillment paths are exercised without a live provider.
// It is registered only when PAYMENT_MOCK_ENABLED is set, never in production.
type MockProvider struct {
	webhookSecret string
}

var _ service.PaymentProvider = (*MockProvider)(nil)

// MockConfig configures the mock provider. WebhookSecret is the HMAC secret the
// simulated webhook is signed with.
type MockConfig struct {
	WebhookSecret string
}

// NewMockProvider builds a MockProvider.
func NewMockProvider(cfg MockConfig) *MockProvider {
	return &MockProvider{webhookSecret: cfg.WebhookSecret}
}

// Name returns the provider's entity name.
func (m *MockProvider) Name() string { return string(entity.ProviderMock) }

// CreateCheckout returns an offline checkout: the redirect URL is a local stub
// keyed on the PaymentID (never visited by the smoke test), and the checkout id
// is the PaymentID so the simulated webhook's tx_ref maps straight back.
func (m *MockProvider) CreateCheckout(ctx context.Context, req service.CheckoutRequest) (*service.CheckoutResult, error) {
	return &service.CheckoutResult{
		CheckoutID:  req.PaymentID,
		RedirectURL: "http://localhost/mock-checkout?tx_ref=" + req.PaymentID,
	}, nil
}

// mockWebhook is the simulated webhook body shape.
type mockWebhook struct {
	EventID string `json:"event_id"`
	Type    string `json:"type"`
	TxRef   string `json:"tx_ref"`
	Status  string `json:"status"`
	TxnID   string `json:"transaction_id"`
}

// VerifyWebhook validates the X-Mock-Signature HMAC over the raw body (so the
// signature path is still enforced) and parses the simulated event.
func (m *MockProvider) VerifyWebhook(ctx context.Context, headers map[string]string, body []byte) (*service.WebhookEvent, error) {
	sig := headerValue(headers, "X-Mock-Signature", "x-mock-signature")
	if !verifyHexSig(m.webhookSecret, sig, body) {
		return nil, apperror.New(http.StatusUnauthorized, "BAD_SIGNATURE", "webhook signature verification failed", nil)
	}
	var wh mockWebhook
	if err := json.Unmarshal(body, &wh); err != nil {
		return nil, apperror.New(http.StatusBadRequest, "BAD_WEBHOOK", "could not parse webhook body", err)
	}
	eventID := wh.EventID
	if eventID == "" {
		eventID = wh.TxRef + ":" + wh.Status
	}
	txnID := wh.TxnID
	if txnID == "" {
		// A stable txn id derived from the ref keeps the provider-txn dedup index
		// effective even when the simulated payload omits one.
		txnID = "mock-" + wh.TxRef
	}
	return &service.WebhookEvent{
		ProviderEventID: eventID,
		Type:            wh.Type,
		PaymentRef:      wh.TxRef,
		Status:          normalizeSandboxStatus(wh.Status),
		ProviderTxnID:   txnID,
		RawPayload:      body,
	}, nil
}

// VerifyPayment reports no information (pending), so a settlement only ever
// happens via the simulated webhook — keeping the dedup assertions deterministic.
func (m *MockProvider) VerifyPayment(ctx context.Context, providerRef string) (string, error) {
	return string(entity.PayPending), nil
}
