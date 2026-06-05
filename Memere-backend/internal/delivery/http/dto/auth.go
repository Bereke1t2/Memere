// Package dto holds the HTTP request/response shapes for the delivery layer.
// This is the only layer that carries json tags (the domain entities stay pure)
// and the only place entity→wire mapping happens. Response DTOs deliberately
// omit secret fields (password hash, tokens) so a leak is structurally
// impossible (Non-Negotiable #8, spec §7.3).
package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// RegisterRequest is the body of POST /auth/register.
type RegisterRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     *string `json:"phone"`
	Role      string  `json:"role"`
}

// LoginRequest is the body of POST /auth/login.
type LoginRequest struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	DeviceInfo *string `json:"device_info"`
}

// RefreshRequest is the body of POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest is the body of POST /auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// UserResponse is the public projection of a user. It has no field for the
// password hash or any token, so those values can never be serialized to a
// client even if a caller passes an un-sanitized entity.
type UserResponse struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Phone           *string    `json:"phone,omitempty"`
	Role            string     `json:"role"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	IsActive        bool       `json:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// NewUserResponse maps a domain user to its public projection.
func NewUserResponse(u *entity.User) UserResponse {
	return UserResponse{
		ID:              u.ID.String(),
		Email:           u.Email,
		Phone:           u.Phone,
		Role:            string(u.Role),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		AvatarURL:       u.AvatarURL,
		IsActive:        u.IsActive,
		IsEmailVerified: u.IsEmailVerified,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// AuthResponse is returned by login and refresh: the token pair plus the
// authenticated user.
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
	User         *UserResponse `json:"user,omitempty"`
}
