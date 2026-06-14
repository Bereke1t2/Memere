package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// NotificationRepo implements repository.NotificationRepository.
type NotificationRepo struct {
	q *sqlcgen.Queries
}

var _ repository.NotificationRepository = (*NotificationRepo)(nil)

func NewNotificationRepo(pool *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{q: sqlcgen.New(pool)}
}

func (r *NotificationRepo) Create(ctx context.Context, n *entity.Notification) (*entity.Notification, error) {
	data, err := json.Marshal(n.Data)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	row, err := queriesFor(ctx, r.q).CreateNotification(ctx, sqlcgen.CreateNotificationParams{
		UserID: toPgUUID(n.UserID),
		Type:   n.Type,
		Title:  n.Title,
		Body:   n.Body,
		Data:   data,
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return notificationRowToEntity(row), nil
}

func (r *NotificationRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.Notification, error) {
	rows, err := queriesFor(ctx, r.q).ListNotificationsByUser(ctx, toPgUUID(userID), int32(limit))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.Notification, 0, len(rows))
	for _, row := range rows {
		out = append(out, notificationRowToEntity(row))
	}
	return out, nil
}

func (r *NotificationRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	row, err := queriesFor(ctx, r.q).GetNotificationByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("notification not found", nil)
		}
		return nil, apperror.Internal(err)
	}
	return notificationRowToEntity(row), nil
}

func (r *NotificationRepo) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	if err := queriesFor(ctx, r.q).MarkNotificationRead(ctx, toPgUUID(id), toPgUUID(userID)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if err := queriesFor(ctx, r.q).MarkAllNotificationsRead(ctx, toPgUUID(userID)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *NotificationRepo) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	n, err := queriesFor(ctx, r.q).UnreadNotificationCount(ctx, toPgUUID(userID))
	if err != nil {
		return 0, apperror.Internal(err)
	}
	return int(n), nil
}

// --- DeviceTokenRepo ---------------------------------------------------------

type DeviceTokenRepo struct {
	q *sqlcgen.Queries
}

var _ repository.DeviceTokenRepository = (*DeviceTokenRepo)(nil)

func NewDeviceTokenRepo(pool *pgxpool.Pool) *DeviceTokenRepo {
	return &DeviceTokenRepo{q: sqlcgen.New(pool)}
}

func (r *DeviceTokenRepo) Upsert(ctx context.Context, t *entity.DeviceToken) error {
	if err := queriesFor(ctx, r.q).UpsertDeviceToken(ctx, sqlcgen.UpsertDeviceTokenParams{
		UserID:   toPgUUID(t.UserID),
		FCMToken: t.FCMToken,
		Platform: t.Platform,
	}); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *DeviceTokenRepo) Delete(ctx context.Context, fcmToken string) error {
	if err := queriesFor(ctx, r.q).DeleteDeviceToken(ctx, fcmToken); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *DeviceTokenRepo) DeleteInvalid(ctx context.Context, tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	if err := queriesFor(ctx, r.q).DeleteInvalidDeviceTokens(ctx, tokens); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *DeviceTokenRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.DeviceToken, error) {
	rows, err := queriesFor(ctx, r.q).ListDeviceTokensByUser(ctx, toPgUUID(userID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.DeviceToken, 0, len(rows))
	for _, row := range rows {
		out = append(out, deviceTokenRowToEntity(row))
	}
	return out, nil
}

// --- PreferenceRepo ----------------------------------------------------------

type PreferenceRepo struct {
	q *sqlcgen.Queries
}

var _ repository.PreferenceRepository = (*PreferenceRepo)(nil)

func NewPreferenceRepo(pool *pgxpool.Pool) *PreferenceRepo {
	return &PreferenceRepo{q: sqlcgen.New(pool)}
}

func (r *PreferenceRepo) Get(ctx context.Context, userID uuid.UUID) (*entity.NotificationPreference, error) {
	row, err := queriesFor(ctx, r.q).GetNotificationPreferences(ctx, toPgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return all-enabled default when no explicit preference exists.
			return &entity.NotificationPreference{
				UserID:       userID,
				PushEnabled:  true,
				EmailEnabled: true,
			}, nil
		}
		return nil, apperror.Internal(err)
	}
	return preferenceRowToEntity(row), nil
}

func (r *PreferenceRepo) Upsert(ctx context.Context, p *entity.NotificationPreference) error {
	if err := queriesFor(ctx, r.q).UpsertNotificationPreferences(ctx, sqlcgen.UpsertNotificationPreferencesParams{
		UserID:       toPgUUID(p.UserID),
		PushEnabled:  p.PushEnabled,
		EmailEnabled: p.EmailEnabled,
	}); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// --- mappers -----------------------------------------------------------------

func notificationRowToEntity(r sqlcgen.NotificationsNotification) *entity.Notification {
	n := &entity.Notification{
		ID:        fromPgUUID(r.ID),
		UserID:    fromPgUUID(r.UserID),
		Type:      r.Type,
		Title:     r.Title,
		Body:      r.Body,
		ReadAt:    fromPgTimestamptz(r.ReadAt),
		CreatedAt: fromPgTimestamptzValue(r.CreatedAt),
	}
	if len(r.Data) > 0 {
		_ = json.Unmarshal(r.Data, &n.Data)
	}
	return n
}

func deviceTokenRowToEntity(r sqlcgen.NotificationsDeviceToken) *entity.DeviceToken {
	return &entity.DeviceToken{
		ID:        fromPgUUID(r.ID),
		UserID:    fromPgUUID(r.UserID),
		FCMToken:  r.FCMToken,
		Platform:  r.Platform,
		CreatedAt: fromPgTimestamptzValue(r.CreatedAt),
	}
}

func preferenceRowToEntity(r sqlcgen.NotificationsPreference) *entity.NotificationPreference {
	return &entity.NotificationPreference{
		UserID:       fromPgUUID(r.UserID),
		PushEnabled:  r.PushEnabled,
		EmailEnabled: r.EmailEnabled,
		UpdatedAt:    fromPgTimestamptzValue(r.UpdatedAt),
	}
}
