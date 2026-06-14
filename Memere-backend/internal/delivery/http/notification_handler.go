package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	notificationuc "github.com/Bereke1t2/Memere/memere-backend/internal/usecase/notification"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// NotificationHandler adapts the notification usecase to HTTP.
type NotificationHandler struct {
	svc *notificationuc.Service
}

func NewNotificationHandler(svc *notificationuc.Service) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// notifActor converts the context actor to a notification.Actor.
func notifActor(c *gin.Context) notificationuc.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return notificationuc.Actor{}
	}
	return notificationuc.Actor{UserID: a.UserID}
}

// List handles GET /me/notifications → 200 with notification list
func (h *NotificationHandler) List(c *gin.Context) {
	limit := atoiDefault(c.Query("limit"), 50)
	notifs, err := h.svc.ListNotifications(c.Request.Context(), notifActor(c), limit)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.NotificationResponse, 0, len(notifs))
	for _, n := range notifs {
		items = append(items, dto.NewNotificationResponse(n))
	}
	respondJSON(c, http.StatusOK, gin.H{"notifications": items})
}

// UnreadCount handles GET /me/notifications/unread-count → 200
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	count, err := h.svc.UnreadCount(c.Request.Context(), notifActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.UnreadCountResponse{Count: count})
}

// MarkRead handles POST /me/notifications/:id/read → 204
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	notifID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid notification id", err))
		return
	}
	if err := h.svc.MarkRead(c.Request.Context(), notifActor(c), notifID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// MarkAllRead handles POST /me/notifications/read-all → 204
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	if err := h.svc.MarkAllRead(c.Request.Context(), notifActor(c)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// RegisterDevice handles POST /me/devices → 201
func (h *NotificationHandler) RegisterDevice(c *gin.Context) {
	var req dto.RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.RegisterDevice(c.Request.Context(), notifActor(c), req.FCMToken, req.Platform); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

// UnregisterDevice handles DELETE /me/devices/:token → 204
func (h *NotificationHandler) UnregisterDevice(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		respondError(c, apperror.BadRequest("token is required", nil))
		return
	}
	if err := h.svc.UnregisterDevice(c.Request.Context(), notifActor(c), token); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetPreferences handles GET /me/notification-preferences → 200
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	prefs, err := h.svc.GetPreferences(c.Request.Context(), notifActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewPreferencesResponse(prefs))
}

// UpdatePreferences handles PUT /me/notification-preferences → 200
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	var req dto.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.UpdatePreferences(c.Request.Context(), notifActor(c), req.PushEnabled, req.EmailEnabled); err != nil {
		respondError(c, err)
		return
	}
	prefs, err := h.svc.GetPreferences(c.Request.Context(), notifActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewPreferencesResponse(prefs))
}
