package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// UserCourseAccessItem represents a course and the user's entitlement state.
type UserCourseAccessItem struct {
	CourseID      uuid.UUID `json:"course_id"`
	Title         string    `json:"title"`
	Subject       string    `json:"subject"`
	Grade         int       `json:"grade"`
	Price         string    `json:"price"`
	Currency      string    `json:"currency"`
	IsFree        bool      `json:"is_free"`
	HasAccess     bool      `json:"has_access"`
	RequestStatus *string   `json:"request_status,omitempty"`
}

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

// ApproveUser marks an account as approved. Admin only. Audited.
func (s *Service) ApproveUser(ctx context.Context, actor Actor, userID uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.ApprovalStatus == entity.ApprovalStatusApproved {
		return nil // idempotent
	}
	u.ApprovalStatus = entity.ApprovalStatusApproved
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "user.approve", "user", &userID, nil)

	if s.notify != nil {
		_ = s.notify.Notify(ctx, service.NotifyEvent{
			UserID: userID.String(),
			Type:   "account_approved",
			Title:  "Account Approved",
			Body:   "Your account has been approved. You can now explore courses and request access.",
			Channels: []service.Channel{
				service.ChannelPush,
				service.ChannelInApp,
				service.ChannelEmail,
			},
		})
	}
	return nil
}

// RejectUser marks an account as rejected. Admin only. Audited.
func (s *Service) RejectUser(ctx context.Context, actor Actor, userID uuid.UUID, reason string) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	u.ApprovalStatus = entity.ApprovalStatusRejected
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "user.reject", "user", &userID, map[string]any{"reason": reason})

	if s.notify != nil {
		body := "Your account application was reviewed and not approved."
		if reason != "" {
			body = "Your account application was rejected: " + reason
		}
		_ = s.notify.Notify(ctx, service.NotifyEvent{
			UserID: userID.String(),
			Type:   "account_rejected",
			Title:  "Account Update",
			Body:   body,
			Channels: []service.Channel{
				service.ChannelInApp,
				service.ChannelEmail,
			},
		})
	}
	return nil
}

// GrantCourseAccess directly enrolls a user in one or more courses and resolves pending requests.
func (s *Service) GrantCourseAccess(ctx context.Context, actor Actor, userID uuid.UUID, courseIDs []uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, cid := range courseIDs {
		// 1. Create enrollment
		_ = s.enrolls.Create(ctx, &entity.Enrollment{
			ID:         uuid.New(),
			StudentID:  userID,
			CourseID:   cid,
			Source:     entity.SourceAdminGrant,
			EnrolledAt: now,
			CreatedAt:  now,
			UpdatedAt:  now,
		})

		// 2. Resolve any pending request
		if s.requests != nil {
			req, err := s.requests.GetByStudentAndCourse(ctx, userID, cid)
			if err == nil && req != nil && req.IsPending() {
				req.Status = entity.RequestStatusApproved
				req.ReviewedBy = &actor.UserID
				req.ReviewedAt = &now
				_ = s.requests.Update(ctx, req)
			}
		}
	}

	s.writeAudit(ctx, actor, "course.grant_access", "user", &userID, map[string]any{
		"course_ids": courseIDs,
	})

	if s.notify != nil && len(courseIDs) > 0 {
		_ = s.notify.Notify(ctx, service.NotifyEvent{
			UserID: userID.String(),
			Type:   "course_access_granted",
			Title:  "Course Access Granted",
			Body:   "You have been granted access to selected course(s).",
			Channels: []service.Channel{
				service.ChannelPush,
				service.ChannelInApp,
				service.ChannelEmail,
			},
		})
	}
	_ = u
	return nil
}

// GrantAllCoursesAccess gives the user access to all published courses.
func (s *Service) GrantAllCoursesAccess(ctx context.Context, actor Actor, userID uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	isPub := true
	courses, _, err := s.courses.List(ctx, repository.CourseFilter{IsPublished: &isPub}, nil, 1000)
	if err != nil {
		return err
	}

	courseIDs := make([]uuid.UUID, 0, len(courses))
	for _, c := range courses {
		courseIDs = append(courseIDs, c.ID)
	}

	err = s.GrantCourseAccess(ctx, actor, userID, courseIDs)
	if err != nil {
		return err
	}

	s.writeAudit(ctx, actor, "course.grant_all_access", "user", &userID, nil)
	_ = u
	return nil
}

// RevokeCourseAccess removes a user's enrollment for a course.
func (s *Service) RevokeCourseAccess(ctx context.Context, actor Actor, userID, courseID uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if err := s.enrolls.Delete(ctx, userID, courseID); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "course.revoke_access", "user", &userID, map[string]any{
		"course_id": courseID.String(),
	})
	return nil
}

// GetUserCourseAccessList returns all courses with access/request status for a user.
func (s *Service) GetUserCourseAccessList(ctx context.Context, actor Actor, userID uuid.UUID) ([]UserCourseAccessItem, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}

	isPub := true
	courses, _, err := s.courses.List(ctx, repository.CourseFilter{IsPublished: &isPub}, nil, 1000)
	if err != nil {
		return nil, err
	}

	enrollments, err := s.enrolls.ListAllByStudent(ctx, userID)
	if err != nil {
		enrollments = nil
	}
	enrolledMap := make(map[uuid.UUID]bool, len(enrollments))
	for _, e := range enrollments {
		enrolledMap[e.CourseID] = true
	}

	var requests []*entity.CourseAccessRequest
	if s.requests != nil {
		requests, _ = s.requests.ListByStudent(ctx, userID)
	}
	requestMap := make(map[uuid.UUID]string, len(requests))
	for _, r := range requests {
		requestMap[r.CourseID] = string(r.Status)
	}

	items := make([]UserCourseAccessItem, 0, len(courses))
	for _, c := range courses {
		hasAccess := c.IsFree || enrolledMap[c.ID]
		var reqStatus *string
		if status, ok := requestMap[c.ID]; ok {
			reqStatus = &status
		}
		items = append(items, UserCourseAccessItem{
			CourseID:      c.ID,
			Title:         c.Title,
			Subject:       c.Subject,
			Grade:         c.Grade,
			Price:         fmt.Sprintf("%.2f", c.Price),
			Currency:      c.Currency,
			IsFree:        c.IsFree,
			HasAccess:     hasAccess,
			RequestStatus: reqStatus,
		})
	}
	return items, nil
}

// SuspendUser sets is_active=false so the user cannot log in. Admin only. Audited.
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
