package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// AdminAuditRepository persists and reads admin audit log entries.
type AdminAuditRepository interface {
	Insert(ctx context.Context, log *entity.AdminAuditLog) error
	List(ctx context.Context, filter AdminAuditFilter, cursor *pagination.Cursor, limit int) ([]*entity.AdminAuditLog, *pagination.Cursor, error)
}

// AdminAuditFilter narrows audit log queries. Nil fields are ignored.
type AdminAuditFilter struct {
	ActorID    *uuid.UUID
	TargetType *string
	TargetID   *uuid.UUID
	From       *time.Time
	To         *time.Time
}

// AdminUserFilter narrows user listing for admin views.
type AdminUserFilter struct {
	Role     *string
	IsActive *bool
	Email    *string
}

// AdminPaymentFilter narrows payment listing for reconciliation views.
type AdminPaymentFilter struct {
	Status   *string
	Provider *string
	From     *time.Time
	To       *time.Time
	StudentID *uuid.UUID
}
