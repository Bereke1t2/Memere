package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// NotificationRepository persists in-app notification rows (spec §11.1).
// All reads are scoped to the authenticated user — callers supply the correct
// userID; there is no server-side IDOR check at this layer.
type NotificationRepository interface {
	// Create inserts a new notification row and returns it with a generated ID.
	Create(ctx context.Context, n *entity.Notification) (*entity.Notification, error)
	// ListByUser returns unread-first notifications for a user, limited to limit rows.
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.Notification, error)
	// GetByID returns a single notification or apperror.NotFound.
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	// MarkRead sets read_at = now() on one notification.
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	// MarkAllRead sets read_at = now() on every unread notification for a user.
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	// UnreadCount returns the count of unread notifications for a user.
	UnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

// DeviceTokenRepository persists FCM device tokens (spec §11.1 push channel).
type DeviceTokenRepository interface {
	// Upsert registers (or refreshes) a device token for a user. Duplicate
	// fcm_token rows from the same device are collapsed by the UNIQUE constraint.
	Upsert(ctx context.Context, t *entity.DeviceToken) error
	// Delete removes a device token (logout / token rotation).
	Delete(ctx context.Context, fcmToken string) error
	// DeleteInvalid prunes tokens reported invalid by FCM (NotRegistered error).
	DeleteInvalid(ctx context.Context, tokens []string) error
	// ListByUser returns all live FCM tokens for a user.
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.DeviceToken, error)
}

// PreferenceRepository persists per-user notification opt-in flags.
type PreferenceRepository interface {
	// Get returns the user's preferences, or a default (all-enabled) when none
	// exist yet.
	Get(ctx context.Context, userID uuid.UUID) (*entity.NotificationPreference, error)
	// Upsert creates or updates the user's preferences.
	Upsert(ctx context.Context, p *entity.NotificationPreference) error
}
