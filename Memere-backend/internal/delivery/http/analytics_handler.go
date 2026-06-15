package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/analytics"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// AnalyticsHandler adapts the analytics usecase to HTTP.
type AnalyticsHandler struct {
	svc *analytics.Service
}

// NewAnalyticsHandler builds an AnalyticsHandler.
func NewAnalyticsHandler(svc *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// analyticsActor converts the context actor into the analytics usecase's actor
// (nil for anonymous).
func analyticsActor(c *gin.Context) *analytics.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return nil
	}
	return &analytics.Actor{UserID: a.UserID, Role: a.Role}
}

// AttemptAnalytics handles GET /exam-attempts/:id/analytics → 200.
func (h *AnalyticsHandler) AttemptAnalytics(c *gin.Context) {
	attemptID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	a, err := h.svc.GetAttemptAnalytics(c.Request.Context(), analyticsActor(c), attemptID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAttemptAnalyticsResponse(a)
	respondJSON(c, http.StatusOK, &resp)
}

// Trend handles GET /me/trend?subject= → 200.
func (h *AnalyticsHandler) Trend(c *gin.Context) {
	subject := c.Query("subject")
	if subject == "" {
		respondError(c, apperror.BadRequest("subject query parameter is required", nil))
		return
	}
	points, err := h.svc.GetStudentTrend(c.Request.Context(), analyticsActor(c), subject)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewTrendResponse(subject, points)
	respondJSON(c, http.StatusOK, &resp)
}

// ExamStats handles GET /exams/:id/stats → 200 (teacher/admin only).
func (h *AnalyticsHandler) ExamStats(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	stats, err := h.svc.GetExamStats(c.Request.Context(), analyticsActor(c), examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamStatsResponse(stats)
	respondJSON(c, http.StatusOK, &resp)
}

// Leaderboard handles GET /exams/:id/leaderboard?limit=10 → 200.
// Returns the top-N students for the exam plus the caller's own rank.
func (h *AnalyticsHandler) Leaderboard(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	limit := clampInt(parseInt(c.Query("limit"), 10), 1, 100)
	result, err := h.svc.GetLeaderboard(c.Request.Context(), analyticsActor(c), examID, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, result)
}
