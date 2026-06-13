// Package payment holds the infrastructure-ring implementations of the
// service.PaymentProvider port: Chapa (primary), Telebirr and Stripe. Each turns
// a provider-agnostic CheckoutRequest into a hosted checkout and verifies the
// signature on inbound webhooks before the usecase trusts a settlement (spec
// §10.1, §7.3). Signing secrets arrive from config and are never logged.
package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/textproto"
)

// signHex returns the lowercase hex HMAC-SHA256 of body under secret. Providers
// use it to recompute the expected webhook signature.
func signHex(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// verifyHexSig constant-time compares the provided signature against the HMAC of
// body under secret. An empty secret or signature is always a failure (we never
// accept an unsigned webhook).
func verifyHexSig(secret, provided string, body []byte) bool {
	if secret == "" || provided == "" {
		return false
	}
	expected := signHex(secret, body)
	return hmac.Equal([]byte(expected), []byte(provided))
}

// headerValue looks a header up case-insensitively across the candidate names,
// tolerating both raw and canonical (MIME) forms so verification does not depend
// on how the delivery layer happened to key the map.
func headerValue(h map[string]string, names ...string) string {
	for _, name := range names {
		if v, ok := h[name]; ok && v != "" {
			return v
		}
		if v, ok := h[textproto.CanonicalMIMEHeaderKey(name)]; ok && v != "" {
			return v
		}
	}
	return ""
}
