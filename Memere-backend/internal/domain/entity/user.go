package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role is the RBAC role of a user (spec §7.2).
type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

// Valid reports whether r is a recognized role.
func (r Role) Valid() bool {
	switch r {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return true
	default:
		return false
	}
}

// User is the core account entity (spec §4.2.1). Pure domain type: no db/json
// tags. Nullable columns are represented as pointers.
type User struct {
	ID                     uuid.UUID
	Email                  string
	Phone                  *string
	PasswordHash           string
	Role                   Role
	FirstName              string
	LastName               string
	AvatarURL              *string
	IsActive               bool
	IsEmailVerified        bool
	EmailVerificationToken *string
	PasswordResetToken     *string
	PasswordResetExpiresAt *time.Time
	LastLoginAt            *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              *time.Time
}
