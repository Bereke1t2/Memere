package payment

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// signedBody returns a webhook body plus the matching hex HMAC under secret.
func signedBody(secret string, payload map[string]any) ([]byte, string) {
	body, _ := json.Marshal(payload)
	return body, signHex(secret, body)
}

func TestChapa_VerifyWebhook_ValidSignature(t *testing.T) {
	secret := "wh_secret"
	c := NewChapaProvider(ChapaConfig{SecretKey: "sk", WebhookSecret: secret})

	body, sig := signedBody(secret, map[string]any{
		"event": "charge.success", "status": "success",
		"tx_ref": "pay-123", "reference": "chapa-ref-9",
	})
	ev, err := c.VerifyWebhook(context.Background(), map[string]string{"Chapa-Signature": sig}, body)
	if err != nil {
		t.Fatalf("valid signature should verify: %v", err)
	}
	if ev.PaymentRef != "pay-123" {
		t.Fatalf("PaymentRef = %q, want pay-123", ev.PaymentRef)
	}
	if ev.Status != string(entity.PayCompleted) {
		t.Fatalf("Status = %q, want completed", ev.Status)
	}
	if ev.ProviderTxnID != "chapa-ref-9" {
		t.Fatalf("ProviderTxnID = %q, want chapa-ref-9", ev.ProviderTxnID)
	}
}

func TestChapa_VerifyWebhook_BadSignatureRejected(t *testing.T) {
	c := NewChapaProvider(ChapaConfig{SecretKey: "sk", WebhookSecret: "wh_secret"})
	body, _ := signedBody("wrong_secret", map[string]any{"event": "charge.success", "status": "success", "tx_ref": "p"})
	_, err := c.VerifyWebhook(context.Background(), map[string]string{"Chapa-Signature": "deadbeef"}, body)
	if !apperror.IsCode(err, "BAD_SIGNATURE") {
		t.Fatalf("tampered signature must be rejected, got %v", err)
	}
}

func TestChapa_VerifyWebhook_MissingSignatureRejected(t *testing.T) {
	c := NewChapaProvider(ChapaConfig{SecretKey: "sk", WebhookSecret: "wh_secret"})
	body, _ := signedBody("wh_secret", map[string]any{"event": "charge.success", "status": "success", "tx_ref": "p"})
	_, err := c.VerifyWebhook(context.Background(), map[string]string{}, body)
	if !apperror.IsCode(err, "BAD_SIGNATURE") {
		t.Fatalf("missing signature must be rejected, got %v", err)
	}
}

// The sandbox providers must also enforce their signatures and resolve by name.
func TestSandboxProviders_SignatureAndRegistry(t *testing.T) {
	tb := NewTelebirrProvider(TelebirrConfig{WebhookSecret: "tb"})
	st := NewStripeProvider(StripeConfig{WebhookSecret: "st"})
	reg := NewRegistry(NewChapaProvider(ChapaConfig{SecretKey: "sk"}), tb, st)

	for _, name := range []entity.PaymentProvider{entity.ProviderChapa, entity.ProviderTelebirr, entity.ProviderStripe} {
		if _, err := reg.Get(name); err != nil {
			t.Fatalf("registry should resolve %s: %v", name, err)
		}
	}
	if _, err := reg.Get(entity.PaymentProvider("paypal")); !apperror.IsCode(err, "PROVIDER_NOT_CONFIGURED") {
		t.Fatalf("unknown provider should be PROVIDER_NOT_CONFIGURED, got %v", err)
	}

	// Telebirr valid + bad signature.
	body, sig := signedBody("tb", map[string]any{"event_id": "e1", "type": "paid", "tx_ref": "p1", "status": "success"})
	if _, err := tb.VerifyWebhook(context.Background(), map[string]string{"Telebirr-Signature": sig}, body); err != nil {
		t.Fatalf("telebirr valid sig: %v", err)
	}
	if _, err := tb.VerifyWebhook(context.Background(), map[string]string{"Telebirr-Signature": "nope"}, body); !apperror.IsCode(err, "BAD_SIGNATURE") {
		t.Fatalf("telebirr bad sig must reject, got %v", err)
	}

	// Stripe valid signature.
	sbody, ssig := signedBody("st", map[string]any{"id": "s1", "type": "checkout.session.completed", "client_reference_id": "p2", "status": "paid"})
	ev, err := st.VerifyWebhook(context.Background(), map[string]string{"Stripe-Signature": ssig}, sbody)
	if err != nil {
		t.Fatalf("stripe valid sig: %v", err)
	}
	if ev.PaymentRef != "p2" || ev.Status != string(entity.PayCompleted) {
		t.Fatalf("stripe event = %+v", ev)
	}
}
