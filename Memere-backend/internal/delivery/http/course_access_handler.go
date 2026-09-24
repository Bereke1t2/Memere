package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/courseaccess"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// CourseAccessHandler adapts courseaccess.Service to HTTP.
type CourseAccessHandler struct {
	svc *courseaccess.Service
}

// NewCourseAccessHandler constructs the handler.
func NewCourseAccessHandler(svc *courseaccess.Service) *CourseAccessHandler {
	return &CourseAccessHandler{svc: svc}
}

func courseAccessActor(c *gin.Context) courseaccess.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return courseaccess.Actor{UserID: uuid.Nil, Role: entity.RoleStudent}
	}
	return courseaccess.Actor{UserID: a.UserID, Role: a.Role}
}

// GetAccessStatus handles GET /courses/:id/access-status
func (h *CourseAccessHandler) GetAccessStatus(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid course id", err))
		return
	}

	res, err := h.svc.GetAccessStatus(c.Request.Context(), courseAccessActor(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusOK, dto.NewCourseAccessStatusResponse(res))
}

// RequestAccess handles POST /courses/:id/request-access
func (h *CourseAccessHandler) RequestAccess(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid course id", err))
		return
	}

	req, err := h.svc.RequestAccess(c.Request.Context(), courseAccessActor(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondJSON(c, http.StatusCreated, dto.NewCourseAccessRequestResponse(req))
}

// ListMyRequests handles GET /me/course-requests
func (h *CourseAccessHandler) ListMyRequests(c *gin.Context) {
	requests, err := h.svc.ListStudentRequests(c.Request.Context(), courseAccessActor(c))
	if err != nil {
		respondError(c, err)
		return
	}

	items := make([]dto.CourseAccessRequestResponse, 0, len(requests))
	for _, req := range requests {
		items = append(items, dto.NewCourseAccessRequestResponse(req))
	}

	respondJSON(c, http.StatusOK, gin.H{"requests": items})
}

// ListTeacherRequests handles GET /teacher/course-requests
func (h *CourseAccessHandler) ListTeacherRequests(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	requests, next, err := h.svc.ListTeacherRequests(c.Request.Context(), courseAccessActor(c), status, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	items := make([]dto.CourseAccessRequestResponse, 0, len(requests))
	for _, req := range requests {
		items = append(items, dto.NewCourseAccessRequestResponse(req))
	}

	nextCursor := ""
	if next != nil {
		nextCursor = next.Encode()
	}

	respondJSON(c, http.StatusOK, gin.H{
		"requests": items,
		"next":     nextCursor,
	})
}

// ApproveTeacherRequest handles POST /teacher/course-requests/:id/approve
func (h *CourseAccessHandler) ApproveTeacherRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid request id", err))
		return
	}

	if err := h.svc.ApproveTeacherRequest(c.Request.Context(), courseAccessActor(c), requestID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RejectTeacherRequest handles POST /teacher/course-requests/:id/reject
func (h *CourseAccessHandler) RejectTeacherRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid request id", err))
		return
	}

	var in dto.RejectRequestInput
	_ = c.ShouldBindJSON(&in)

	if err := h.svc.RejectTeacherRequest(c.Request.Context(), courseAccessActor(c), requestID, in.Reason); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
