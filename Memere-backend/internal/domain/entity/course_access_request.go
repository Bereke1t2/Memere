package entity

import (
	"time"

	"github.com/google/uuid"
)

// RequestStatus represents the state of a course access request.
type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"
	RequestStatusApproved RequestStatus = "approved"
	RequestStatusRejected RequestStatus = "rejected"
)

// Valid reports whether s is a valid request status.
func (s RequestStatus) Valid() bool {
	switch s {
	case RequestStatusPending, RequestStatusApproved, RequestStatusRejected:
		return true
	default:
		return false
	}
}

// CourseAccessRequest represents a student's application to access a course.
type CourseAccessRequest struct {
	ID              uuid.UUID
	StudentID       uuid.UUID
	CourseID        uuid.UUID
	Status          RequestStatus
	ReviewedBy      *uuid.UUID
	ReviewedAt      *time.Time
	RejectionReason *string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Display metadata populated in list/get queries
	StudentName  string
	StudentEmail string
	CourseTitle  string
}

// IsPending reports whether the request is waiting for review.
func (r *CourseAccessRequest) IsPending() bool {
	return r.Status == RequestStatusPending
}

// IsApproved reports whether the request was approved.
func (r *CourseAccessRequest) IsApproved() bool {
	return r.Status == RequestStatusApproved
}

// IsRejected reports whether the request was rejected.
func (r *CourseAccessRequest) IsRejected() bool {
	return r.Status == RequestStatusRejected
}
