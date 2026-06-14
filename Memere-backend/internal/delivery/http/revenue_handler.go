package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/revenue"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// defaultRevenueWindow is the look-back applied when a report request omits the
// from/to query params (last 30 days).
const defaultRevenueWindow = 30 * 24 * time.Hour

// RevenueHandler adapts the read-only revenue-reporting usecase to HTTP. Role and
// ownership scoping live in the usecase; the admin/teacher route middleware is a
// first cheap gate.
type RevenueHandler struct {
	svc *revenue.Service
	now func() time.Time
}

// NewRevenueHandler builds a RevenueHandler. clock may be nil (defaults to
// time.Now) and is used only to default an omitted reporting window.
func NewRevenueHandler(svc *revenue.Service, clock func() time.Time) *RevenueHandler {
	if clock == nil {
		clock = time.Now
	}
	return &RevenueHandler{svc: svc, now: clock}
}

// revenueActor adapts the context actor to the revenue usecase's Actor.
func revenueActor(c *gin.Context) revenue.Actor {
	a, ok := middleware.ActorFromContext(c)
	if !ok || a == nil {
		return revenue.Actor{}
	}
	return revenue.Actor{UserID: a.UserID, Role: a.Role}
}

// PlatformRevenue handles GET /admin/revenue (bearer + admin): the platform-wide
// revenue report over an optional [from,to] window (defaults to the last 30 days).
func (h *RevenueHandler) PlatformRevenue(c *gin.Context) {
	from, to, err := h.window(c)
	if err != nil {
		respondError(c, err)
		return
	}
	rep, err := h.svc.PlatformRevenue(c.Request.Context(), revenueActor(c), from, to)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, rep)
}

// MyEarnings handles GET /me/earnings (bearer + teacher/admin): the caller's own
// earnings under the configured teacher share.
func (h *RevenueHandler) MyEarnings(c *gin.Context) {
	from, to, err := h.window(c)
	if err != nil {
		respondError(c, err)
		return
	}
	actor := revenueActor(c)
	rep, err := h.svc.TeacherEarnings(c.Request.Context(), actor, actor.UserID, from, to)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, rep)
}

// CourseSales handles GET /courses/:id/sales (bearer, owner/admin): lifetime
// sales for one course.
func (h *RevenueHandler) CourseSales(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	rep, err := h.svc.CourseSalesStats(c.Request.Context(), revenueActor(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, rep)
}

// window parses the optional from/to RFC3339 query params, defaulting to the last
// 30 days. A malformed timestamp is a clean 400.
func (h *RevenueHandler) window(c *gin.Context) (time.Time, time.Time, error) {
	to := h.now().UTC()
	from := to.Add(-defaultRevenueWindow)
	if s := c.Query("from"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Time{}, time.Time{}, apperror.BadRequest("invalid 'from' timestamp (want RFC3339)", err)
		}
		from = t
	}
	if s := c.Query("to"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Time{}, time.Time{}, apperror.BadRequest("invalid 'to' timestamp (want RFC3339)", err)
		}
		to = t
	}
	return from, to, nil
}
