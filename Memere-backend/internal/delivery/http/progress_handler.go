package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/progress"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ProgressHandler adapts the progress usecase to HTTP.
type ProgressHandler struct {
	svc *progress.Service
}

func NewProgressHandler(svc *progress.Service) *ProgressHandler {
	return &ProgressHandler{svc: svc}
}

// progressActor converts the context actor to a progress.Actor.
func progressActor(c *gin.Context) progress.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return progress.Actor{}
	}
	return progress.Actor{UserID: a.UserID, Role: a.Role}
}

// Complete handles POST /lessons/:id/complete → 204
func (h *ProgressHandler) Complete(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid lesson id", err))
		return
	}
	if _, err := h.svc.MarkLessonComplete(c.Request.Context(), progressActor(c), lessonID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// VideoProgress handles PUT /lessons/:id/video-progress → 200 with updated progress
func (h *ProgressHandler) VideoProgress(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid lesson id", err))
		return
	}
	var req dto.VideoProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	p, err := h.svc.UpdateVideoProgress(c.Request.Context(), progressActor(c), lessonID, req.PositionSeconds)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewLessonProgressResponse(p))
}

// GetCourseProgress handles GET /courses/:id/progress → 200 with rollup + lesson list
func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid course id", err))
		return
	}
	cp, lessons, err := h.svc.GetCourseProgress(c.Request.Context(), progressActor(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewCourseProgressResponse(cp, lessons))
}

// GetDashboard handles GET /me/dashboard → 200 with enrolled courses + streak
func (h *ProgressHandler) GetDashboard(c *gin.Context) {
	dash, err := h.svc.GetMyDashboard(c.Request.Context(), progressActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewDashboardResponse(dash))
}

// GetStreak handles GET /me/streak → 200 with streak data
func (h *ProgressHandler) GetStreak(c *gin.Context) {
	streak, err := h.svc.GetMyStreak(c.Request.Context(), progressActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewStreakResponse(streak))
}
