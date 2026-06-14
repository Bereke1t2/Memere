package entity

import (
	"time"

	"github.com/google/uuid"
)

// Notification is an in-app notification row (spec §11.1).
type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      string
	Title     string
	Body      string
	Data      map[string]string // deep-link payload
	ReadAt    *time.Time
	CreatedAt time.Time
}

// IsRead reports whether the notification has been read.
func (n *Notification) IsRead() bool { return n.ReadAt != nil }

// DeviceToken holds an FCM registration token for a user's device.
type DeviceToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FCMToken  string
	Platform  string // android | ios
	CreatedAt time.Time
}

// NotificationPreference holds per-user channel opt-in flags.
type NotificationPreference struct {
	UserID       uuid.UUID
	PushEnabled  bool
	EmailEnabled bool
	UpdatedAt    time.Time
}
