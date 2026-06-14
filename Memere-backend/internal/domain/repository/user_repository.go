// Package repository defines the data-access contracts (interfaces) that
// usecases depend on. Implementations live in internal/repository/* (Skills
// 3–4). Per Clean Architecture, this package imports only domain entities,
// pkg utilities, and the standard library — never infrastructure.
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// UserRepository persists and retrieves users (spec §4.2.1). Every read filters
// out soft-deleted rows (deleted_at IS NULL).
type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, u *entity.User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	SetLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error
	// List returns users matching filter (admin use). Bounded by limit.
	List(ctx context.Context, filter AdminUserFilter, cursor *pagination.Cursor, limit int) ([]*entity.User, *pagination.Cursor, error)
	// CountByRole returns the number of active users with the given role.
	CountByRole(ctx context.Context, role entity.Role) (int, error)
}
