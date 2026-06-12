// Package enrollment holds the enrollment-engine usecases: granting a student
// access to a course (the fulfillment path used by payments and the free path)
// and listing a student's enrollments. Like the other usecase packages it is
// pure orchestration over the domain repository interfaces — no HTTP, no SQL.
//
// Granting is idempotent: the (student_id, course_id) UNIQUE constraint is the
// hard guarantee, and an Exists check is the friendly fast path, so a
// re-processed payment grants access exactly once (spec §4.2.7, Non-Negotiable
// #4 idempotency).
package enrollment

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Actor is the authenticated caller's identity. A zero Actor (UserID == Nil) is
// an anonymous request, which may not enrol.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// Service implements the enrollment usecases over the enrollment and course
// repositories. clock is injectable so tests can pin enrollment timestamps.
type Service struct {
	enroll repository.EnrollmentRepository
	course repository.CourseRepository
	clock  func() time.Time
}

// NewService wires the enrollment service. When clock is nil it defaults to
// time.Now.
func NewService(
	enroll repository.EnrollmentRepository,
	course repository.CourseRepository,
	clock func() time.Time,
) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{enroll: enroll, course: course, clock: clock}
}

// GrantEnrollment grants a student access to a course from the given source.
// It is idempotent: if the student already has any enrollment for the course it
// no-ops and returns nil. expiresAt is nil for a permanent grant
// (purchase/coupon/free) and set for subscription-derived access.
//
// This is the fulfillment primitive the payment path (Skill 3) calls after a
// payment settles; it never decides *whether* the student paid — only records
// that they may now access the course.
func (s *Service) GrantEnrollment(
	ctx context.Context,
	studentID, courseID uuid.UUID,
	source entity.EnrollmentSource,
	expiresAt *time.Time,
) error {
	if studentID == uuid.Nil || courseID == uuid.Nil {
		return apperror.New(400, "INVALID_ARGUMENT", "student and course are required", nil)
	}
	if !source.Valid() {
		return apperror.New(400, "INVALID_ARGUMENT", "unknown enrollment source", nil)
	}
	// Friendly fast path; Create is also idempotent at the DB (UNIQUE), so a
	// race that slips past this check still grants exactly once.
	exists, err := s.enroll.Exists(ctx, studentID, courseID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	now := s.clock()
	return s.enroll.Create(ctx, &entity.Enrollment{
		ID:         uuid.New(),
		StudentID:  studentID,
		CourseID:   courseID,
		Source:     source,
		EnrolledAt: now,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
}

// EnrollFree enrolls the actor in a course that is free. It refuses a paid
// course (those must go through the payment flow), preventing a student from
// self-granting access to paid content.
func (s *Service) EnrollFree(ctx context.Context, actor Actor, courseID uuid.UUID) error {
	if actor.UserID == uuid.Nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	c, err := s.course.FindByID(ctx, courseID)
	if err != nil {
		return err
	}
	if !c.IsFree {
		return apperror.New(409, "COURSE_NOT_FREE", "this course requires purchase", nil)
	}
	return s.GrantEnrollment(ctx, actor.UserID, courseID, entity.SourceFree, nil)
}

// ListMyEnrollments returns the actor's own enrollments (IDOR-safe: scoped to
// the authenticated student id, never a client-supplied id).
func (s *Service) ListMyEnrollments(ctx context.Context, actor Actor, limit int) ([]*entity.Enrollment, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.enroll.ListByStudent(ctx, actor.UserID, limit)
}

// IsEnrolled reports whether the actor holds any enrollment for the course. It
// is a thin, IDOR-safe wrapper over the repository Exists check.
func (s *Service) IsEnrolled(ctx context.Context, actor Actor, courseID uuid.UUID) (bool, error) {
	if actor.UserID == uuid.Nil {
		return false, nil
	}
	return s.enroll.Exists(ctx, actor.UserID, courseID)
}
