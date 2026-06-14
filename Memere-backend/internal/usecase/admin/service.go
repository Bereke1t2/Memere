// Package admin implements privileged admin/operations usecases (spec §1.3,
// §1.4, §3.2). All state-changing operations are (a) gated by role and (b)
// write an admin_audit_log row. Reads are also role-gated but not audited
// (audit log is for mutations only).
package admin

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Actor is the authenticated caller. A zero Actor (UserID == Nil) or one with
// a non-admin role is rejected by every state-changing admin usecase.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// isAdmin reports whether the actor holds the admin role.
func (a Actor) isAdmin() bool { return a.Role == entity.RoleAdmin }

// isTeacherOrAdmin reports whether the actor may read teacher-scoped reports.
func (a Actor) isTeacherOrAdmin() bool {
	return a.Role == entity.RoleAdmin || a.Role == entity.RoleTeacher
}

// Service is the stateless admin usecase orchestrator.
type Service struct {
	users    repository.UserRepository
	courses  repository.CourseRepository
	payments repository.PaymentRepository
	enrolls  repository.EnrollmentRepository
	subs     repository.SubscriptionRepository
	revenue  repository.RevenueRepository
	audit    repository.AdminAuditRepository
	notify   service.Notifier
	providers service.PaymentProviderRegistry
}

func NewService(
	users repository.UserRepository,
	courses repository.CourseRepository,
	payments repository.PaymentRepository,
	enrolls repository.EnrollmentRepository,
	subs repository.SubscriptionRepository,
	revenue repository.RevenueRepository,
	audit repository.AdminAuditRepository,
	notify service.Notifier,
	providers service.PaymentProviderRegistry,
) *Service {
	return &Service{
		users:     users,
		courses:   courses,
		payments:  payments,
		enrolls:   enrolls,
		subs:      subs,
		revenue:   revenue,
		audit:     audit,
		notify:    notify,
		providers: providers,
	}
}

// writeAudit records an audit entry best-effort: errors are logged but never
// surfaced to the caller (the audit write must not roll back a business action).
func (s *Service) writeAudit(ctx context.Context, actor Actor, action, targetType string, targetID *uuid.UUID, details map[string]any) {
	if err := s.audit.Insert(ctx, &entity.AdminAuditLog{
		ActorID:    actor.UserID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Details:    details,
	}); err != nil {
		log.Printf("admin_audit: failed to write audit log: %v", err)
	}
}

// requireAdmin returns Forbidden when the actor is not an admin.
func requireAdmin(actor Actor) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	if !actor.isAdmin() {
		return apperror.Forbidden("admin role required", nil)
	}
	return nil
}

// requireTeacherOrAdmin returns Forbidden when the actor is neither teacher nor admin.
func requireTeacherOrAdmin(actor Actor) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	if !actor.isTeacherOrAdmin() {
		return apperror.Forbidden("teacher or admin role required", nil)
	}
	return nil
}
