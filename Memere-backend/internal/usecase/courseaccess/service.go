package courseaccess

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/notification"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// Actor represents the authenticated caller.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// AccessStatusResult describes the caller's entitlement state for a course.
type AccessStatusResult struct {
	HasAccess       bool                 `json:"has_access"`
	AccessLevel     string               `json:"access_level"`
	Status          string               `json:"status"` // granted | account_pending | request_pending | not_enrolled | rejected | anonymous
	RequestID       *uuid.UUID           `json:"request_id,omitempty"`
	RejectionReason *string              `json:"rejection_reason,omitempty"`
}

// Service orchestrates course access requests and moderation.
type Service struct {
	requests repository.CourseAccessRequestRepository
	enroll   repository.EnrollmentRepository
	courses  repository.CourseRepository
	users    repository.UserRepository
	notifier *notification.Hooks
	clock    func() time.Time
}

// NewService constructs the course access service.
func NewService(
	requests repository.CourseAccessRequestRepository,
	enroll repository.EnrollmentRepository,
	courses repository.CourseRepository,
	users repository.UserRepository,
	notifier *notification.Hooks,
	clock func() time.Time,
) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{
		requests: requests,
		enroll:   enroll,
		courses:  courses,
		users:    users,
		notifier: notifier,
		clock:    clock,
	}
}

// RequestAccess allows an approved student to apply for a paid course.
func (s *Service) RequestAccess(ctx context.Context, actor Actor, courseID uuid.UUID) (*entity.CourseAccessRequest, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}

	u, err := s.users.FindByID(ctx, actor.UserID)
	if err != nil {
		return nil, err
	}
	if !u.IsApproved() {
		return nil, apperror.Forbidden("account is awaiting admin approval", nil)
	}

	c, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// If course is free, enroll directly
	if c.IsFree {
		now := s.clock()
		_ = s.enroll.Create(ctx, &entity.Enrollment{
			ID:         uuid.New(),
			StudentID:  actor.UserID,
			CourseID:   courseID,
			Source:     entity.SourceFree,
			EnrolledAt: now,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		return &entity.CourseAccessRequest{
			ID:        uuid.New(),
			StudentID: actor.UserID,
			CourseID:  courseID,
			Status:    entity.RequestStatusApproved,
		}, nil
	}

	// Guard: check if student already enrolled
	alreadyEnrolled, err := s.enroll.Exists(ctx, actor.UserID, courseID)
	if err != nil {
		return nil, err
	}
	if alreadyEnrolled {
		return nil, apperror.Conflict("ALREADY_ENROLLED", nil)
	}

	// Guard: check existing request
	existing, err := s.requests.GetByStudentAndCourse(ctx, actor.UserID, courseID)
	if err == nil && existing != nil {
		if existing.IsPending() {
			return nil, apperror.Conflict("REQUEST_ALREADY_PENDING", nil)
		}
		if existing.IsApproved() {
			return nil, apperror.Conflict("ALREADY_APPROVED", nil)
		}
	}

	req := &entity.CourseAccessRequest{
		ID:        uuid.New(),
		StudentID: actor.UserID,
		CourseID:  courseID,
		Status:    entity.RequestStatusPending,
		CreatedAt: s.clock(),
		UpdatedAt: s.clock(),
	}

	if err := s.requests.Create(ctx, req); err != nil {
		return nil, err
	}

	if s.notifier != nil {
		s.notifier.CourseAccessRequested(ctx, c.TeacherID, u.FirstName+" "+u.LastName, c.Title, req.ID)
	}

	return req, nil
}

// GetAccessStatus computes the granular access and application state for a course.
func (s *Service) GetAccessStatus(ctx context.Context, actor Actor, courseID uuid.UUID) (*AccessStatusResult, error) {
	c, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// Free courses always granted
	if c.IsFree {
		return &AccessStatusResult{
			HasAccess:   true,
			AccessLevel: "full",
			Status:      "granted",
		}, nil
	}

	// Owner or Admin always granted
	if actor.Role == entity.RoleAdmin || (actor.Role == entity.RoleTeacher && c.TeacherID == actor.UserID) {
		return &AccessStatusResult{
			HasAccess:   true,
			AccessLevel: "full",
			Status:      "granted",
		}, nil
	}

	if actor.UserID == uuid.Nil {
		return &AccessStatusResult{
			HasAccess:   false,
			AccessLevel: "preview",
			Status:      "anonymous",
		}, nil
	}

	// Check if student holds active enrollment
	enrolled, err := s.enroll.Exists(ctx, actor.UserID, courseID)
	if err == nil && enrolled {
		return &AccessStatusResult{
			HasAccess:   true,
			AccessLevel: "full",
			Status:      "granted",
		}, nil
	}

	// Check if user account is approved
	u, err := s.users.FindByID(ctx, actor.UserID)
	if err != nil {
		return nil, err
	}
	if !u.IsApproved() {
		return &AccessStatusResult{
			HasAccess:   false,
			AccessLevel: "preview",
			Status:      "account_pending",
		}, nil
	}

	// Check if a course request exists
	req, err := s.requests.GetByStudentAndCourse(ctx, actor.UserID, courseID)
	if err == nil && req != nil {
		if req.IsPending() {
			return &AccessStatusResult{
				HasAccess:   false,
				AccessLevel: "preview",
				Status:      "request_pending",
				RequestID:   &req.ID,
			}, nil
		}
		if req.IsRejected() {
			return &AccessStatusResult{
				HasAccess:       false,
				AccessLevel:     "preview",
				Status:          "rejected",
				RequestID:       &req.ID,
				RejectionReason: req.RejectionReason,
			}, nil
		}
	}

	return &AccessStatusResult{
		HasAccess:   false,
		AccessLevel: "preview",
		Status:      "not_enrolled",
	}, nil
}

// ListTeacherRequests returns requests for courses belonging to the authenticated teacher.
func (s *Service) ListTeacherRequests(ctx context.Context, actor Actor, status *string, cursor *pagination.Cursor, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	if actor.Role != entity.RoleTeacher && actor.Role != entity.RoleAdmin {
		return nil, nil, apperror.Forbidden("teacher or admin role required", nil)
	}
	return s.requests.ListByTeacher(ctx, actor.UserID, status, cursor, limit)
}

// ApproveTeacherRequest approves a course access request for an owned course.
func (s *Service) ApproveTeacherRequest(ctx context.Context, actor Actor, requestID uuid.UUID) error {
	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	c, err := s.courses.FindByID(ctx, req.CourseID)
	if err != nil {
		return err
	}

	if actor.Role != entity.RoleAdmin && (actor.Role != entity.RoleTeacher || c.TeacherID != actor.UserID) {
		return apperror.Forbidden("cannot moderate requests for courses you do not manage", nil)
	}

	now := s.clock()
	req.Status = entity.RequestStatusApproved
	req.ReviewedBy = &actor.UserID
	req.ReviewedAt = &now

	if err := s.requests.Update(ctx, req); err != nil {
		return err
	}

	// Grant enrollment
	_ = s.enroll.Create(ctx, &entity.Enrollment{
		ID:         uuid.New(),
		StudentID:  req.StudentID,
		CourseID:   req.CourseID,
		Source:     entity.SourceAdminGrant,
		EnrolledAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	if s.notifier != nil {
		s.notifier.CourseAccessApproved(ctx, req.StudentID, c.Title, c.ID)
	}

	return nil
}

// RejectTeacherRequest rejects a course access request for an owned course.
func (s *Service) RejectTeacherRequest(ctx context.Context, actor Actor, requestID uuid.UUID, reason string) error {
	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	c, err := s.courses.FindByID(ctx, req.CourseID)
	if err != nil {
		return err
	}

	if actor.Role != entity.RoleAdmin && (actor.Role != entity.RoleTeacher || c.TeacherID != actor.UserID) {
		return apperror.Forbidden("cannot moderate requests for courses you do not manage", nil)
	}

	now := s.clock()
	req.Status = entity.RequestStatusRejected
	req.ReviewedBy = &actor.UserID
	req.ReviewedAt = &now
	req.RejectionReason = &reason

	if err := s.requests.Update(ctx, req); err != nil {
		return err
	}

	if s.notifier != nil {
		s.notifier.CourseAccessRejected(ctx, req.StudentID, c.Title, reason)
	}

	return nil
}

// ListAdminRequests returns all requests on the platform.
func (s *Service) ListAdminRequests(ctx context.Context, actor Actor, status *string, cursor *pagination.Cursor, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	if actor.Role != entity.RoleAdmin {
		return nil, nil, apperror.Forbidden("admin role required", nil)
	}
	return s.requests.ListAll(ctx, status, cursor, limit)
}

// ApproveAdminRequest approves any course access request.
func (s *Service) ApproveAdminRequest(ctx context.Context, actor Actor, requestID uuid.UUID) error {
	if actor.Role != entity.RoleAdmin {
		return apperror.Forbidden("admin role required", nil)
	}

	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	c, err := s.courses.FindByID(ctx, req.CourseID)
	if err != nil {
		return err
	}

	now := s.clock()
	req.Status = entity.RequestStatusApproved
	req.ReviewedBy = &actor.UserID
	req.ReviewedAt = &now

	if err := s.requests.Update(ctx, req); err != nil {
		return err
	}

	_ = s.enroll.Create(ctx, &entity.Enrollment{
		ID:         uuid.New(),
		StudentID:  req.StudentID,
		CourseID:   req.CourseID,
		Source:     entity.SourceAdminGrant,
		EnrolledAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	if s.notifier != nil {
		s.notifier.CourseAccessApproved(ctx, req.StudentID, c.Title, c.ID)
	}

	return nil
}

// RejectAdminRequest rejects any course access request.
func (s *Service) RejectAdminRequest(ctx context.Context, actor Actor, requestID uuid.UUID, reason string) error {
	if actor.Role != entity.RoleAdmin {
		return apperror.Forbidden("admin role required", nil)
	}

	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	c, err := s.courses.FindByID(ctx, req.CourseID)
	if err != nil {
		return err
	}

	now := s.clock()
	req.Status = entity.RequestStatusRejected
	req.ReviewedBy = &actor.UserID
	req.ReviewedAt = &now
	req.RejectionReason = &reason

	if err := s.requests.Update(ctx, req); err != nil {
		return err
	}

	if s.notifier != nil {
		s.notifier.CourseAccessRejected(ctx, req.StudentID, c.Title, reason)
	}

	return nil
}

// ListStudentRequests returns all requests submitted by the calling student.
func (s *Service) ListStudentRequests(ctx context.Context, actor Actor) ([]*entity.CourseAccessRequest, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.requests.ListByStudent(ctx, actor.UserID)
}
