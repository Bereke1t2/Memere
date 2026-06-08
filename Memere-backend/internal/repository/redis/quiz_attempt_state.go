package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// AttemptStateRepo implements repository.AttemptStateStore over Redis. It backs
// both quiz and exam attempts (spec §9.1/§9.2); the attempt ID namespaces the
// keys, so quiz and exam attempts (distinct UUIDs) never collide.
type AttemptStateRepo struct {
	client *goredis.Client
}

var _ repository.AttemptStateStore = (*AttemptStateRepo)(nil)

// NewAttemptStateRepo builds an AttemptStateRepo over a go-redis client.
func NewAttemptStateRepo(client *goredis.Client) *AttemptStateRepo {
	return &AttemptStateRepo{client: client}
}

func snapshotKey(attemptID uuid.UUID) string {
	return fmt.Sprintf("quiz:attempt:%s:order", attemptID.String())
}

func answersKey(attemptID uuid.UUID) string {
	return fmt.Sprintf("quiz:attempt:%s:answers", attemptID.String())
}

func (r *AttemptStateRepo) SetSnapshot(ctx context.Context, attemptID uuid.UUID, snapshot map[string]any, ttl time.Duration) error {
	return r.setJSON(ctx, snapshotKey(attemptID), snapshot, ttl)
}

func (r *AttemptStateRepo) GetSnapshot(ctx context.Context, attemptID uuid.UUID) (map[string]any, error) {
	return r.getJSON(ctx, snapshotKey(attemptID))
}

func (r *AttemptStateRepo) SaveAnswers(ctx context.Context, attemptID uuid.UUID, answers map[string]any, ttl time.Duration) error {
	return r.setJSON(ctx, answersKey(attemptID), answers, ttl)
}

func (r *AttemptStateRepo) GetAnswers(ctx context.Context, attemptID uuid.UUID) (map[string]any, error) {
	return r.getJSON(ctx, answersKey(attemptID))
}

// DeleteAttemptState clears both keys in one round trip. Deleting missing keys
// is a no-op.
func (r *AttemptStateRepo) DeleteAttemptState(ctx context.Context, attemptID uuid.UUID) error {
	if err := r.client.Del(ctx, snapshotKey(attemptID), answersKey(attemptID)).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *AttemptStateRepo) setJSON(ctx context.Context, key string, v map[string]any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return apperror.Internal(err)
	}
	if err := r.client.Set(ctx, key, b, ttl).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// getJSON returns (nil, nil) on a Redis miss — a normal outcome the caller
// handles by falling back to the durable Postgres copy.
func (r *AttemptStateRepo) getJSON(ctx context.Context, key string) (map[string]any, error) {
	b, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, apperror.Internal(err)
	}
	return m, nil
}
