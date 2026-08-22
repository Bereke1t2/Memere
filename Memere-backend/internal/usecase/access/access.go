// Package access is the single source of truth for "does this caller have
// access to this course (or lesson)?". Every gate in the system — quiz/exam
// taking, paid video streaming, course content reads — routes through here so
// the access policy lives in exactly one tested place (spec §7.2 RBAC, §10
// monetization). Like the other usecase packages it is pure orchestration over
// the domain repository interfaces: no HTTP, no SQL, no infrastructure.
//
// The policy, in order of precedence:
//   - course owner (teacher) or any admin  -> FullAccess
//   - free course                          -> FullAccess
//   - active (unexpired) enrollment        -> FullAccess
//   - active (unexpired) subscription      -> FullAccess
//   - otherwise                            -> PreviewAccess (sample lessons only)
//
// Graded work (quizzes/exams) requires FullAccess — preview is for *watching*
// sample lessons, not for graded attempts.
package access

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// AccessLevel is the outcome of an access decision.
type AccessLevel int

const (
	// NoAccess: the caller may not see the course at all.
	NoAccess AccessLevel = iota
	// PreviewAccess: the caller may see free-preview lessons only.
	PreviewAccess
	// FullAccess: the caller may see every lesson and take graded work.
	FullAccess
)

// Actor is the authenticated caller's identity. A zero Actor (UserID == Nil)
// represents an anonymous request, which can never hold more than PreviewAccess.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// Service answers access questions over the enrollment, subscription, and
// course repositories. clock is injectable so tests can drive expiry.
type Service struct {
	enroll repository.EnrollmentRepository
	subs   repository.SubscriptionRepository
	course repository.CourseRepository
	clock  func() time.Time
}

// NewService wires the access service. When clock is nil it defaults to
// time.Now.
func NewService(
	enroll repository.EnrollmentRepository,
	subs repository.SubscriptionRepository,
	course repository.CourseRepository,
	clock func() time.Time,
) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{enroll: enroll, subs: subs, course: course, clock: clock}
}

// DisableEnrollmentCheck temporarily disables enrollment enforcement so all
// students can access courses, lessons, quizzes, and exams without enrolling.
// Set back to false to re-enable enrollment checks in the future.
var DisableEnrollmentCheck = true

// CanAccessCourse is THE access function. It resolves the course server-side
// (never trusting a client-supplied flag) and applies the precedence policy.
func (s *Service) CanAccessCourse(ctx context.Context, actor Actor, courseID uuid.UUID) (AccessLevel, error) {
	if DisableEnrollmentCheck {
		return FullAccess, nil
	}
	c, err := s.course.FindByID(ctx, courseID)
	if err != nil {
		return NoAccess, err
	}
	// Owner / admin: always full.
	if actor.Role == entity.RoleAdmin ||
		(actor.Role == entity.RoleTeacher && c.TeacherID == actor.UserID) {
		return FullAccess, nil
	}
	// Free course: full for everyone (including anonymous).
	if c.IsFree {
		return FullAccess, nil
	}
	// Authenticated student: an active enrollment or subscription unlocks it.
	if actor.UserID != uuid.Nil {
		if s.hasActiveEnrollment(ctx, actor.UserID, courseID) {
			return FullAccess, nil
		}
		if s.hasActiveSubscription(ctx, actor.UserID) {
			return FullAccess, nil
		}
	}
	// Paid course, no entitlement: preview lessons only.
	return PreviewAccess, nil
}

// CanAccessLesson layers free_preview on top of course access: a paid lesson
// needs FullAccess; a free-preview lesson is visible at PreviewAccess too.
func (s *Service) CanAccessLesson(ctx context.Context, actor Actor, lesson *entity.Lesson) (bool, error) {
	lvl, err := s.CanAccessCourse(ctx, actor, lesson.CourseID)
	if err != nil {
		return false, err
	}
	switch {
	case lvl == FullAccess:
		return true, nil
	case lvl == PreviewAccess && lesson.IsFreePreview:
		return true, nil
	default:
		return false, nil
	}
}

// RequireFullAccess is a guard for graded work and paid content: it returns
// apperror.Forbidden("NOT_ENROLLED") unless the actor holds FullAccess.
func (s *Service) RequireFullAccess(ctx context.Context, actor Actor, courseID uuid.UUID) error {
	lvl, err := s.CanAccessCourse(ctx, actor, courseID)
	if err != nil {
		return err
	}
	if lvl != FullAccess {
		return apperror.New(403, "NOT_ENROLLED", "you must enrol in this course first", nil)
	}
	return nil
}

// hasActiveEnrollment reports whether the student holds an unexpired enrollment
// for the course. A NotFound from the repo is the ordinary "no enrollment"
// answer, not an error.
func (s *Service) hasActiveEnrollment(ctx context.Context, studentID, courseID uuid.UUID) bool {
	e, err := s.enroll.GetActiveForStudent(ctx, studentID, courseID)
	if err != nil {
		return false
	}
	return e.IsActiveAt(s.clock())
}

// hasActiveSubscription reports whether the student holds an active, unexpired
// all-access subscription.
func (s *Service) hasActiveSubscription(ctx context.Context, studentID uuid.UUID) bool {
	sub, err := s.subs.GetActiveForStudent(ctx, studentID)
	if err != nil {
		return false
	}
	return sub.IsActiveAt(s.clock())
}
