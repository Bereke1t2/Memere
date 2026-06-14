package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/enrollment"
)

// EnrollmentHandler adapts the enrollment usecase to HTTP. Eligibility (free
// course only) and ownership are enforced in the usecase, not middleware.
type EnrollmentHandler struct {
	svc *enrollment.Service
}

// NewEnrollmentHandler builds an EnrollmentHandler.
func NewEnrollmentHandler(svc *enrollment.Service) *EnrollmentHandler {
	return &EnrollmentHandler{svc: svc}
}

// enrollmentActor adapts the context actor to the enrollment usecase's Actor.
func enrollmentActor(c *gin.Context) enrollment.Actor {
	a, ok := middleware.ActorFromContext(c)
	if !ok || a == nil {
		return enrollment.Actor{}
	}
	return enrollment.Actor{UserID: a.UserID, Role: a.Role}
}

// EnrollFree handles POST /courses/:id/enroll-free (bearer): grants access to a
// free course directly. The usecase rejects a paid course (which must go through
// the payment flow).
func (h *EnrollmentHandler) EnrollFree(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.EnrollFree(c.Request.Context(), enrollmentActor(c), courseID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListMine handles GET /me/enrollments (bearer): the caller's own enrollments.
func (h *EnrollmentHandler) ListMine(c *gin.Context) {
	limit := atoiDefault(c.Query("limit"), 0)
	es, err := h.svc.ListMyEnrollments(c.Request.Context(), enrollmentActor(c), limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"data": dto.NewEnrollmentList(es)})
}
