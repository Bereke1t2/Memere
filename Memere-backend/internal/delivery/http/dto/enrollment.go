package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// EnrollmentResponse is the client-safe view of one enrollment. expires_at is
// omitted for a permanent grant (purchase/coupon/free) and set for
// subscription-derived access.
type EnrollmentResponse struct {
	ID         string  `json:"id"`
	CourseID   string  `json:"course_id"`
	Source     string  `json:"source"`
	EnrolledAt string  `json:"enrolled_at"`
	ExpiresAt  *string `json:"expires_at,omitempty"`
}

// NewEnrollmentResponse maps an entity enrollment to the wire shape.
func NewEnrollmentResponse(e *entity.Enrollment) EnrollmentResponse {
	resp := EnrollmentResponse{
		ID:         e.ID.String(),
		CourseID:   e.CourseID.String(),
		Source:     string(e.Source),
		EnrolledAt: e.EnrolledAt.UTC().Format(time.RFC3339),
	}
	if e.ExpiresAt != nil {
		s := e.ExpiresAt.UTC().Format(time.RFC3339)
		resp.ExpiresAt = &s
	}
	return resp
}

// NewEnrollmentList maps a slice of enrollments to wire shapes.
func NewEnrollmentList(es []*entity.Enrollment) []EnrollmentResponse {
	out := make([]EnrollmentResponse, 0, len(es))
	for _, e := range es {
		out = append(out, NewEnrollmentResponse(e))
	}
	return out
}
