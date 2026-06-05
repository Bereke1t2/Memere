package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/auth"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// AuthHandler adapts the auth usecase to HTTP. It also holds the user repository
// so /auth/me can load the current user from the authenticated actor.
type AuthHandler struct {
	svc   *auth.Service
	users repository.UserRepository
}

// NewAuthHandler builds an AuthHandler.
func NewAuthHandler(svc *auth.Service, users repository.UserRepository) *AuthHandler {
	return &AuthHandler{svc: svc, users: users}
}

// Register handles POST /auth/register → 201 with the sanitized user.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}

	user, err := h.svc.Register(c.Request.Context(), auth.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      entity.Role(req.Role),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewUserResponse(user)
	respondJSON(c, http.StatusCreated, &resp)
}

// Login handles POST /auth/login → 200 with the token pair + user.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}

	tokens, user, err := h.svc.Login(c.Request.Context(), auth.LoginInput{
		Email:      req.Email,
		Password:   req.Password,
		DeviceInfo: req.DeviceInfo,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, authResponse(tokens, user))
}

// Refresh handles POST /auth/refresh → 200 with a fresh token pair (the
// presented refresh token is rotated/revoked by the usecase).
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	tokens, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, authResponse(tokens, nil))
}

// Logout handles POST /auth/logout → 200. It revokes the presented refresh
// token and clears the session; it is idempotent.
func (h *AuthHandler) Logout(c *gin.Context) {
	actor, ok := middleware.ActorFromContext(c)
	if !ok || actor == nil {
		respondError(c, apperror.Unauthorized("authentication required", nil))
		return
	}
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.Logout(c.Request.Context(), actor.UserID, req.RefreshToken); err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"status": "ok"})
}

// Me handles GET /auth/me → 200 with the current user.
func (h *AuthHandler) Me(c *gin.Context) {
	actor, ok := middleware.ActorFromContext(c)
	if !ok || actor == nil {
		respondError(c, apperror.Unauthorized("authentication required", nil))
		return
	}
	user, err := h.users.FindByID(c.Request.Context(), actor.UserID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewUserResponse(user)
	respondJSON(c, http.StatusOK, &resp)
}

// authResponse builds the AuthResponse envelope from tokens and an optional
// user (refresh omits the user).
func authResponse(tokens *auth.AuthTokens, user *entity.User) *dto.AuthResponse {
	resp := &dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}
	if user != nil {
		u := dto.NewUserResponse(user)
		resp.User = &u
	}
	return resp
}
