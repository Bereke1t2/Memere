package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// NotificationResponse is the wire projection of entity.Notification.
type NotificationResponse struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Data      map[string]string `json:"data,omitempty"`
	ReadAt    *time.Time        `json:"read_at,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// UnreadCountResponse is returned by GET /me/notifications/unread-count.
type UnreadCountResponse struct {
	Count int `json:"count"`
}

// RegisterDeviceRequest is the body of POST /me/devices.
type RegisterDeviceRequest struct {
	FCMToken string `json:"fcm_token"`
	Platform string `json:"platform"` // android | ios
}

// UpdatePreferencesRequest is the body of PUT /me/notification-preferences.
type UpdatePreferencesRequest struct {
	PushEnabled  bool `json:"push_enabled"`
	EmailEnabled bool `json:"email_enabled"`
}

// PreferencesResponse is the wire projection of entity.NotificationPreference.
type PreferencesResponse struct {
	PushEnabled  bool `json:"push_enabled"`
	EmailEnabled bool `json:"email_enabled"`
}

// NewNotificationResponse maps a domain notification to the wire DTO.
func NewNotificationResponse(n *entity.Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID.String(),
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		Data:      n.Data,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
	}
}

// NewPreferencesResponse maps preferences to the wire DTO.
func NewPreferencesResponse(p *entity.NotificationPreference) PreferencesResponse {
	return PreferencesResponse{
		PushEnabled:  p.PushEnabled,
		EmailEnabled: p.EmailEnabled,
	}
}
