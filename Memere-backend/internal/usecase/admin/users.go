package admin

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ListUsers returns a paginated list of users. Admin only.
func (s *Service) ListUsers(ctx context.Context, actor Actor, filter repository.AdminUserFilter, cursor *pagination.Cursor, limit int) ([]*entity.User, *pagination.Cursor, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.users.List(ctx, filter, cursor, limit)
}

// GetUser returns one user by ID. Admin only.
func (s *Service) GetUser(ctx context.Context, actor Actor, userID uuid.UUID) (*entity.User, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.users.FindByID(ctx, userID)
}

// SuspendUser sets is_active=false so the user cannot log in. Admin only.
// Audited.
func (s *Service) SuspendUser(ctx context.Context, actor Actor, userID uuid.UUID, reason string) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !u.IsActive {
		return nil // idempotent
	}
	u.IsActive = false
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "user.suspend", "user", &userID, map[string]any{"reason": reason})
	return nil
}

// ReactivateUser sets is_active=true. Admin only. Audited.
func (s *Service) ReactivateUser(ctx context.Context, actor Actor, userID uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.IsActive {
		return nil // idempotent
	}
	u.IsActive = true
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "user.reactivate", "user", &userID, nil)
	return nil
}

// ChangeRole promotes or demotes a user's role. Admin only. Guards against
// self-demotion of the last admin. Audited.
func (s *Service) ChangeRole(ctx context.Context, actor Actor, userID uuid.UUID, newRole entity.Role) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if !newRole.Valid() {
		return apperror.BadRequest("invalid role", nil)
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.Role == newRole {
		return nil // idempotent
	}
	// Guard: do not demote the last active admin.
	if u.Role == entity.RoleAdmin && newRole != entity.RoleAdmin {
		count, err := s.users.CountByRole(ctx, entity.RoleAdmin)
		if err != nil {
			return err
		}
		if count <= 1 {
			return apperror.BadRequest("cannot demote the last admin", nil)
		}
	}
	oldRole := string(u.Role)
	u.Role = newRole
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "user.change_role", "user", &userID, map[string]any{
		"old_role": oldRole,
		"new_role": string(newRole),
	})
	return nil
}

// SoftDeleteUser tombstones a user. Admin only. Audited.
func (s *Service) SoftDeleteUser(ctx context.Context, actor Actor, userID uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if err := s.users.SoftDelete(ctx, userID); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "user.delete", "user", &userID, nil)
	return nil
}
