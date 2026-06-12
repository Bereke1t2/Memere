package entity

import (
	"time"

	"github.com/google/uuid"
)

// WebhookEvent is a row in the dedup ledger (spec §7.3 "webhook deduplication").
// The (Provider, ProviderEventID) pair is unique; InsertIfNew uses it to process
// a re-delivered provider event at most once. ProcessedAt is nil until the
// fulfillment for the event commits.
type WebhookEvent struct {
	ID              uuid.UUID
	Provider        PaymentProvider
	ProviderEventID string
	EventType       string
	Payload         []byte
	ProcessedAt     *time.Time
	CreatedAt       time.Time
}
