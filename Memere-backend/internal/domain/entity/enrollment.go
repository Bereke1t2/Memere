package entity

import (
	"time"

	"github.com/google/uuid"
)

// EnrollmentSource records how a student gained access to a course (spec §4.2.7:
// enrollments.source ENUM purchase / subscription / free / coupon).
type EnrollmentSource string

const (
	SourcePurchase     EnrollmentSource = "purchase"
	SourceSubscription EnrollmentSource = "subscription"
	SourceFree         EnrollmentSource = "free"
	SourceCoupon       EnrollmentSource = "coupon"
)

// Valid reports whether s is a known enrollment source.
func (s EnrollmentSource) Valid() bool {
	switch s {
	case SourcePurchase, SourceSubscription, SourceFree, SourceCoupon:
		return true
	}
	return false
}

// Enrollment grants one student access to one course. ExpiresAt is nil for a
// permanent (purchase/coupon/free) grant and set for subscription-derived access
// (spec §4.2.7: enrollments.expires_at nullable, subscription-based expiry).
type Enrollment struct {
	ID         uuid.UUID
	StudentID  uuid.UUID
	CourseID   uuid.UUID
	Source     EnrollmentSource
	EnrolledAt time.Time
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// IsActiveAt reports whether the enrollment grants access at instant t: it is
// active when it has no expiry or has not yet expired. Access checks call this
// rather than inspecting ExpiresAt directly.
func (e *Enrollment) IsActiveAt(t time.Time) bool {
	return e.ExpiresAt == nil || e.ExpiresAt.After(t)
}
