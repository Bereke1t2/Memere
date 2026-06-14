package notification

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Actor is the authenticated caller. A zero Actor (UserID == Nil) is unauthenticated.
type Actor struct {
	UserID uuid.UUID
}

// Service provides in-app notification, device-token, and preference management.
// All operations are scoped to the authenticated caller (IDOR prevention).
type Service struct {
	notifications repository.NotificationRepository
	devices       repository.DeviceTokenRepository
	prefs         repository.PreferenceRepository
}

func NewService(
	notifications repository.NotificationRepository,
	devices repository.DeviceTokenRepository,
	prefs repository.PreferenceRepository,
) *Service {
	return &Service{
		notifications: notifications,
		devices:       devices,
		prefs:         prefs,
	}
}

// ListNotifications returns up to limit notifications for the calling user,
// unread first.
func (s *Service) ListNotifications(ctx context.Context, actor Actor, limit int) ([]*entity.Notification, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.notifications.ListByUser(ctx, actor.UserID, limit)
}

// MarkRead marks a single notification as read. Returns NotFound when the
// notification doesn't exist or belongs to a different user (IDOR-safe).
func (s *Service) MarkRead(ctx context.Context, actor Actor, notificationID uuid.UUID) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	// Verify ownership before marking (the SQL also filters by user_id, so this
	// is defence-in-depth).
	n, err := s.notifications.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if n.UserID != actor.UserID {
		return apperror.NotFound("notification not found", nil)
	}
	return s.notifications.MarkRead(ctx, notificationID, actor.UserID)
}

// MarkAllRead marks every unread notification as read for the calling user.
func (s *Service) MarkAllRead(ctx context.Context, actor Actor) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	return s.notifications.MarkAllRead(ctx, actor.UserID)
}

// UnreadCount returns the count of unread notifications for the calling user.
func (s *Service) UnreadCount(ctx context.Context, actor Actor) (int, error) {
	if actor.UserID == uuid.Nil {
		return 0, apperror.Unauthorized("authentication required", nil)
	}
	return s.notifications.UnreadCount(ctx, actor.UserID)
}

// RegisterDevice registers (or refreshes) an FCM device token for the caller.
func (s *Service) RegisterDevice(ctx context.Context, actor Actor, fcmToken, platform string) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	if fcmToken == "" {
		return apperror.BadRequest("fcm_token is required", nil)
	}
	if platform != "android" && platform != "ios" {
		return apperror.BadRequest("platform must be android or ios", nil)
	}
	return s.devices.Upsert(ctx, &entity.DeviceToken{
		UserID:   actor.UserID,
		FCMToken: fcmToken,
		Platform: platform,
	})
}

// UnregisterDevice removes an FCM device token (on logout or token rotation).
func (s *Service) UnregisterDevice(ctx context.Context, actor Actor, fcmToken string) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	if fcmToken == "" {
		return apperror.BadRequest("fcm_token is required", nil)
	}
	return s.devices.Delete(ctx, fcmToken)
}

// GetPreferences returns the caller's notification preferences (defaults to
// all-enabled when no explicit preference has been saved).
func (s *Service) GetPreferences(ctx context.Context, actor Actor) (*entity.NotificationPreference, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.prefs.Get(ctx, actor.UserID)
}

// UpdatePreferences persists the caller's channel opt-in flags.
func (s *Service) UpdatePreferences(ctx context.Context, actor Actor, pushEnabled, emailEnabled bool) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	return s.prefs.Upsert(ctx, &entity.NotificationPreference{
		UserID:       actor.UserID,
		PushEnabled:  pushEnabled,
		EmailEnabled: emailEnabled,
	})
}
