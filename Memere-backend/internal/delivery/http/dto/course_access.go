package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/admin"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/courseaccess"
)

// CourseAccessStatusResponse is returned by GET /courses/:id/access-status.
type CourseAccessStatusResponse struct {
	HasAccess       bool    `json:"has_access"`
	AccessLevel     string  `json:"access_level"`
	Status          string  `json:"status"` // granted | account_pending | request_pending | not_enrolled | rejected | anonymous
	RequestID       *string `json:"request_id,omitempty"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

// NewCourseAccessStatusResponse maps courseaccess.AccessStatusResult to wire DTO.
func NewCourseAccessStatusResponse(r *courseaccess.AccessStatusResult) CourseAccessStatusResponse {
	var reqID *string
	if r.RequestID != nil {
		s := r.RequestID.String()
		reqID = &s
	}
	return CourseAccessStatusResponse{
		HasAccess:       r.HasAccess,
		AccessLevel:     r.AccessLevel,
		Status:          r.Status,
		RequestID:       reqID,
		RejectionReason: r.RejectionReason,
	}
}

// CourseAccessRequestResponse is the wire projection of a CourseAccessRequest.
type CourseAccessRequestResponse struct {
	ID              string     `json:"id"`
	StudentID       string     `json:"student_id"`
	StudentName     string     `json:"student_name"`
	StudentEmail    string     `json:"student_email"`
	CourseID        string     `json:"course_id"`
	CourseTitle     string     `json:"course_title"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
}

// NewCourseAccessRequestResponse maps domain request entity to wire DTO.
func NewCourseAccessRequestResponse(r *entity.CourseAccessRequest) CourseAccessRequestResponse {
	return CourseAccessRequestResponse{
		ID:              r.ID.String(),
		StudentID:       r.StudentID.String(),
		StudentName:     r.StudentName,
		StudentEmail:    r.StudentEmail,
		CourseID:        r.CourseID.String(),
		CourseTitle:     r.CourseTitle,
		Status:          string(r.Status),
		RejectionReason: r.RejectionReason,
		CreatedAt:       r.CreatedAt,
		ReviewedAt:      r.ReviewedAt,
	}
}

// RejectRequestInput is the optional body for request rejection.
type RejectRequestInput struct {
	Reason string `json:"reason,omitempty"`
}

// GrantCourseAccessRequest is the body for admin course access granting.
type GrantCourseAccessRequest struct {
	CourseIDs []string `json:"course_ids,omitempty"`
	GrantAll  bool     `json:"grant_all,omitempty"`
}

// RevokeCourseAccessRequest is the body for admin course access revocation.
type RevokeCourseAccessRequest struct {
	CourseID string `json:"course_id"`
}

// UserCourseAccessResponse item in the admin user course access list.
type UserCourseAccessResponse struct {
	CourseID      string  `json:"course_id"`
	Title         string  `json:"title"`
	Subject       string  `json:"subject"`
	Grade         int     `json:"grade"`
	Price         string  `json:"price"`
	Currency      string  `json:"currency"`
	IsFree        bool    `json:"is_free"`
	HasAccess     bool    `json:"has_access"`
	RequestStatus *string `json:"request_status,omitempty"`
}

// NewUserCourseAccessList maps domain items to wire DTOs.
func NewUserCourseAccessList(items []admin.UserCourseAccessItem) []UserCourseAccessResponse {
	out := make([]UserCourseAccessResponse, 0, len(items))
	for _, item := range items {
		out = append(out, UserCourseAccessResponse{
			CourseID:      item.CourseID.String(),
			Title:         item.Title,
			Subject:       item.Subject,
			Grade:         item.Grade,
			Price:         item.Price,
			Currency:      item.Currency,
			IsFree:        item.IsFree,
			HasAccess:     item.HasAccess,
			RequestStatus: item.RequestStatus,
		})
	}
	return out
}
