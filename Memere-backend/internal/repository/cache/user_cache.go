package cache

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

const (
	userTTL   = 5 * time.Minute
	userClass = "user_id"
)

// CachedUserRepo is a cache-aside decorator for UserRepository. It caches
// FindByID — the hot path hit on every authenticated request — and delegates
// all other calls to the underlying Postgres repository.
//
// Key: users:id:{uuid}  TTL: 5 minutes
// Invalidation: Update and SoftDelete DEL the key.
type CachedUserRepo struct {
	next repository.UserRepository
	rdb  *goredis.Client
}

var _ repository.UserRepository = (*CachedUserRepo)(nil)

// NewCachedUserRepo wraps next with a 5-minute FindByID cache.
func NewCachedUserRepo(next repository.UserRepository, rdb *goredis.Client) *CachedUserRepo {
	return &CachedUserRepo{next: next, rdb: rdb}
}

func userIDKey(id uuid.UUID) string { return "users:id:" + id.String() }

// FindByID returns the user from cache if available; falls back to Postgres on
// a miss.
func (r *CachedUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return GetOrSet(ctx, r.rdb, userIDKey(id), userClass, userTTL,
		func(ctx context.Context) (*entity.User, error) {
			return r.next.FindByID(ctx, id)
		})
}

// --- Writes invalidate the cache ---

func (r *CachedUserRepo) Update(ctx context.Context, u *entity.User) error {
	if err := r.next.Update(ctx, u); err != nil {
		return err
	}
	Del(ctx, r.rdb, userIDKey(u.ID))
	return nil
}

func (r *CachedUserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.next.SoftDelete(ctx, id); err != nil {
		return err
	}
	Del(ctx, r.rdb, userIDKey(id))
	return nil
}

// --- Uncached pass-throughs ---

func (r *CachedUserRepo) Create(ctx context.Context, u *entity.User) error {
	return r.next.Create(ctx, u)
}

func (r *CachedUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.next.FindByEmail(ctx, email)
}

func (r *CachedUserRepo) SetLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error {
	return r.next.SetLastLogin(ctx, id, t)
}

func (r *CachedUserRepo) List(ctx context.Context, filter repository.AdminUserFilter, cursor *pagination.Cursor, limit int) ([]*entity.User, *pagination.Cursor, error) {
	return r.next.List(ctx, filter, cursor, limit)
}

func (r *CachedUserRepo) CountByRole(ctx context.Context, role entity.Role) (int, error) {
	return r.next.CountByRole(ctx, role)
}
