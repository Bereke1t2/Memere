package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

const chapaDefaultBaseURL = "https://api.chapa.co/v1"

// ChapaProvider is the primary payment backend (spec §10.1). It initializes a
// hosted checkout over the Chapa REST API and verifies the HMAC signature on
// webhook deliveries. secretKey and webhookSecret are credentials — never logged.
type ChapaProvider struct {
	secretKey     string
	webhookSecret string
	baseURL       string
	httpc         *http.Client
}

var _ service.PaymentProvider = (*ChapaProvider)(nil)

// ChapaConfig configures the Chapa client. BaseURL defaults to the production
// API; WebhookSecret falls back to SecretKey when unset (Chapa signs with the
// account secret unless a dedicated webhook secret is configured).
type ChapaConfig struct {
	SecretKey     string
	WebhookSecret string
	BaseURL       string
	HTTPClient    *http.Client
}

// NewChapaProvider builds a ChapaProvider from config.
func NewChapaProvider(cfg ChapaConfig) *ChapaProvider {
	base := cfg.BaseURL
	if base == "" {
		base = chapaDefaultBaseURL
	}
	httpc := cfg.HTTPClient
	if httpc == nil {
		httpc = &http.Client{Timeout: 15 * time.Second}
	}
	whSecret := cfg.WebhookSecret
	if whSecret == "" {
		whSecret = cfg.SecretKey
	}
	return &ChapaProvider{
		secretKey:     cfg.SecretKey,
		webhookSecret: whSecret,
		baseURL:       base,
		httpc:         httpc,
	}
}

// Name returns the provider's entity name.
func (c *ChapaProvider) Name() string { return string(entity.ProviderChapa) }

// CreateCheckout POSTs to /transaction/initialize with tx_ref == PaymentID and
// returns the hosted checkout URL plus the reference we store as
// provider_checkout_id.
func (c *ChapaProvider) CreateCheckout(ctx context.Context, req service.CheckoutRequest) (*service.CheckoutResult, error) {
	body, _ := json.Marshal(map[string]any{
		"amount":       req.Amount.String(),
		"currency":     req.Currency,
		"tx_ref":       req.PaymentID,
		"callback_url": req.CallbackURL,
		"return_url":   req.ReturnURL,
		"description":  req.Description,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/transaction/initialize", bytes.NewReader(body))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpc.Do(httpReq)
	if err != nil {
		return nil, apperror.New(http.StatusBadGateway, "PROVIDER_UNAVAILABLE", "could not reach the payment provider", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			CheckoutURL string `json:"checkout_url"`
			Reference   string `json:"reference"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, apperror.New(http.StatusBadGateway, "PROVIDER_BAD_RESPONSE", "invalid response from the payment provider", err)
	}
	if resp.StatusCode != http.StatusOK || parsed.Data.CheckoutURL == "" {
		return nil, apperror.New(http.StatusBadGateway, "PROVIDER_CHECKOUT_FAILED",
			fmt.Sprintf("provider rejected the checkout (status %d)", resp.StatusCode), nil)
	}
	checkoutID := parsed.Data.Reference
	if checkoutID == "" {
		// Chapa keys the verify call on tx_ref; fall back to it as the reference.
		checkoutID = req.PaymentID
	}
	return &service.CheckoutResult{
		CheckoutID:  checkoutID,
		RedirectURL: parsed.Data.CheckoutURL,
	}, nil
}

// chapaWebhook is the subset of Chapa's webhook body we consume.
type chapaWebhook struct {
	Event  string `json:"event"`
	Status string `json:"status"`
	TxRef  string `json:"tx_ref"`
	Ref    string `json:"reference"`
}

// VerifyWebhook validates the Chapa-Signature HMAC over the raw body and parses
// the event. A missing or mismatched signature is rejected with BAD_SIGNATURE so
// the usecase never settles on an unverified payload (§7.3).
func (c *ChapaProvider) VerifyWebhook(ctx context.Context, headers map[string]string, body []byte) (*service.WebhookEvent, error) {
	sig := headerValue(headers, "Chapa-Signature", "x-chapa-signature")
	if !verifyHexSig(c.webhookSecret, sig, body) {
		return nil, apperror.New(http.StatusUnauthorized, "BAD_SIGNATURE", "webhook signature verification failed", nil)
	}
	var wh chapaWebhook
	if err := json.Unmarshal(body, &wh); err != nil {
		return nil, apperror.New(http.StatusBadRequest, "BAD_WEBHOOK", "could not parse webhook body", err)
	}
	return &service.WebhookEvent{
		ProviderEventID: chapaEventID(wh, body),
		Type:            wh.Event,
		PaymentRef:      wh.TxRef,
		Status:          normalizeChapaStatus(wh.Status),
		ProviderTxnID:   wh.Ref,
		RawPayload:      body,
	}, nil
}

// VerifyPayment confirms a transaction server-to-server (polling fallback) via
// /transaction/verify/{ref}, returning the normalized status.
func (c *ChapaProvider) VerifyPayment(ctx context.Context, providerRef string) (string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/transaction/verify/"+providerRef, nil)
	if err != nil {
		return "", apperror.Internal(err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.secretKey)
	resp, err := c.httpc.Do(httpReq)
	if err != nil {
		return "", apperror.New(http.StatusBadGateway, "PROVIDER_UNAVAILABLE", "could not reach the payment provider", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var parsed struct {
		Status string `json:"status"`
		Data   struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", apperror.New(http.StatusBadGateway, "PROVIDER_BAD_RESPONSE", "invalid response from the payment provider", err)
	}
	return normalizeChapaStatus(parsed.Data.Status), nil
}

// chapaEventID derives a stable dedup id. Chapa does not always send an explicit
// event id, so we key on tx_ref+status (a given payment settles once); when a
// reference is present we prefer it.
func chapaEventID(wh chapaWebhook, body []byte) string {
	switch {
	case wh.Ref != "":
		return wh.Ref
	case wh.TxRef != "":
		return wh.TxRef + ":" + wh.Status
	default:
		return signHex("chapa-event", body)
	}
}

// normalizeChapaStatus maps Chapa status strings onto entity.PaymentStatus
// values the usecase understands.
func normalizeChapaStatus(s string) string {
	switch s {
	case "success", "successful", "completed":
		return string(entity.PayCompleted)
	case "failed", "cancelled", "canceled":
		return string(entity.PayFailed)
	default:
		return s
	}
}
