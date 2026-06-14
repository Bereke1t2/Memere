package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	adminuc "github.com/Bereke1t2/Memere/memere-backend/internal/usecase/admin"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// AdminHandler adapts the admin usecase to HTTP.
type AdminHandler struct {
	svc *adminuc.Service
}

func NewAdminHandler(svc *adminuc.Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// adminActor converts the context actor to an admin.Actor.
func adminActor(c *gin.Context) adminuc.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return adminuc.Actor{}
	}
	return adminuc.Actor{UserID: a.UserID, Role: a.Role}
}

// ---- User management --------------------------------------------------------

// ListUsers handles GET /admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))
	filter := repository.AdminUserFilter{}
	if r := c.Query("role"); r != "" {
		filter.Role = &r
	}

	users, next, err := h.svc.ListUsers(c.Request.Context(), adminActor(c), filter, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		r := dto.NewUserResponse(u)
		items = append(items, r)
	}
	nextCursor := ""
	if next != nil {
		nextCursor = next.Encode()
	}
	respondJSON(c, http.StatusOK, gin.H{
		"users": items,
		"next":  nextCursor,
	})
}

// GetUser handles GET /admin/users/:id
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid user id", err))
		return
	}
	u, err := h.svc.GetUser(c.Request.Context(), adminActor(c), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	r := dto.NewUserResponse(u)
	respondJSON(c, http.StatusOK, &r)
}

// SuspendUser handles POST /admin/users/:id/suspend → 204
func (h *AdminHandler) SuspendUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid user id", err))
		return
	}
	var req dto.SuspendUserRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.SuspendUser(c.Request.Context(), adminActor(c), userID, req.Reason); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ReactivateUser handles POST /admin/users/:id/reactivate → 204
func (h *AdminHandler) ReactivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid user id", err))
		return
	}
	if err := h.svc.ReactivateUser(c.Request.Context(), adminActor(c), userID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ChangeRole handles POST /admin/users/:id/role → 204
func (h *AdminHandler) ChangeRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid user id", err))
		return
	}
	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.ChangeRole(c.Request.Context(), adminActor(c), userID, entity.Role(req.Role)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ---- Content moderation -----------------------------------------------------

// ListCourses handles GET /admin/courses
func (h *AdminHandler) ListCourses(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))

	courses, next, err := h.svc.ListAllCourses(c.Request.Context(), adminActor(c), repository.CourseFilter{}, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.CourseResponse, 0, len(courses))
	for _, course := range courses {
		items = append(items, dto.NewCourseResponse(course))
	}
	nextCursor2 := ""
	if next != nil {
		nextCursor2 = next.Encode()
	}
	respondJSON(c, http.StatusOK, gin.H{
		"courses": items,
		"next":    nextCursor2,
	})
}

// UnpublishCourse handles POST /admin/courses/:id/unpublish → 204
func (h *AdminHandler) UnpublishCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid course id", err))
		return
	}
	var req dto.UnpublishCourseRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.UnpublishCourse(c.Request.Context(), adminActor(c), courseID, req.Reason); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ---- Payment reconciliation -------------------------------------------------

// ListPayments handles GET /admin/payments
func (h *AdminHandler) ListPayments(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))
	filter := repository.AdminPaymentFilter{}
	if s := c.Query("status"); s != "" {
		filter.Status = &s
	}

	payments, next, err := h.svc.ListPayments(c.Request.Context(), adminActor(c), filter, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AdminPaymentResponse, 0, len(payments))
	for _, p := range payments {
		items = append(items, dto.NewAdminPaymentResponse(p))
	}
	nextCursor3 := ""
	if next != nil {
		nextCursor3 = next.Encode()
	}
	respondJSON(c, http.StatusOK, gin.H{
		"payments": items,
		"next":     nextCursor3,
	})
}

// GetPaymentDetail handles GET /admin/payments/:id
func (h *AdminHandler) GetPaymentDetail(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid payment id", err))
		return
	}
	p, err := h.svc.GetPaymentDetail(c.Request.Context(), adminActor(c), paymentID)
	if err != nil {
		respondError(c, err)
		return
	}
	r := dto.NewAdminPaymentResponse(p)
	respondJSON(c, http.StatusOK, &r)
}

// ReconcilePending handles POST /admin/payments/reconcile → 200 with count
func (h *AdminHandler) ReconcilePending(c *gin.Context) {
	n, err := h.svc.ReconcilePending(c.Request.Context(), adminActor(c), 10*time.Minute)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.ReconcileResponse{Reconciled: n})
}

// ---- Announcements ----------------------------------------------------------

// Broadcast handles POST /admin/announcements → 204
func (h *AdminHandler) Broadcast(c *gin.Context) {
	var req dto.BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.Broadcast(c.Request.Context(), adminActor(c), adminuc.BroadcastInput{
		Title:   req.Title,
		Body:    req.Body,
		Segment: adminuc.BroadcastSegment(req.Segment),
		Data:    req.Data,
	}); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ---- Analytics --------------------------------------------------------------

// Overview handles GET /admin/analytics/overview
func (h *AdminHandler) Overview(c *gin.Context) {
	from := parseTimeQuery(c, "from", time.Now().AddDate(0, -1, 0))
	to := parseTimeQuery(c, "to", time.Now())

	ov, err := h.svc.Overview(c.Request.Context(), adminActor(c), from, to)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewOverviewResponse(ov))
}

// RevenueBreakdown handles GET /admin/analytics/revenue
func (h *AdminHandler) RevenueBreakdown(c *gin.Context) {
	from := parseTimeQuery(c, "from", time.Now().AddDate(0, -1, 0))
	to := parseTimeQuery(c, "to", time.Now())

	rows, err := h.svc.RevenueBreakdown(c.Request.Context(), adminActor(c), from, to)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"breakdown": dto.NewRevenueBreakdown(rows)})
}

// GetEngagementStats handles GET /admin/analytics/engagement
func (h *AdminHandler) GetEngagementStats(c *gin.Context) {
	stats, err := h.svc.GetEngagementStats(c.Request.Context(), adminActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewEngagementStatsResponse(stats))
}

// parseTimeQuery parses a ?name= query param as RFC3339; returns def on failure.
func parseTimeQuery(c *gin.Context, name string, def time.Time) time.Time {
	s := c.Query(name)
	if s == "" {
		return def
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return def
	}
	return t
}
