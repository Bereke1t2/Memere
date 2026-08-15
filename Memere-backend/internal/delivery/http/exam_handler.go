package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ExamHandler adapts the exam usecase to HTTP.
type ExamHandler struct {
	svc *exam.Service
}

// NewExamHandler builds an ExamHandler.
func NewExamHandler(svc *exam.Service) *ExamHandler {
	return &ExamHandler{svc: svc}
}

// examActor converts the context actor into the exam usecase's actor (nil for
// anonymous).
func examActor(c *gin.Context) *exam.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return nil
	}
	return &exam.Actor{UserID: a.UserID, Role: a.Role}
}

// ListByCourse handles GET /courses/:id/exams → 200.
func (h *ExamHandler) ListByCourse(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	exams, err := h.svc.ListByCourse(c.Request.Context(), courseID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"data": dto.NewExamListResponse(exams)})
}

// CreateExam handles POST /courses/:id/exams → 201.
func (h *ExamHandler) CreateExam(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	cid := courseID
	created, err := h.svc.CreateExam(c.Request.Context(), examActor(c), exam.CreateExamInput{
		CourseID:        &cid,
		Title:           req.Title,
		Subject:         req.Subject,
		Grade:           req.Grade,
		DurationMinutes: req.DurationMinutes,
		PassMarks:       req.PassMarks,
		Instructions:    req.Instructions,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResponse(created)
	respondJSON(c, http.StatusCreated, &resp)
}

// AddQuestion handles POST /exams/:id/questions → 201 (references an existing
// question by id; returns the updated exam with recomputed total_marks).
func (h *ExamHandler) AddQuestion(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.AddExamQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	updated, err := h.svc.AddExamQuestion(c.Request.Context(), examActor(c), examID, req.QuestionID, req.Marks, req.OrderIndex)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResponse(updated)
	respondJSON(c, http.StatusCreated, &resp)
}

// Publish handles POST /exams/:id/publish → 200. Required for an exam to appear
// in the /mock-exams catalog and to be startable by students.
func (h *ExamHandler) Publish(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	updated, err := h.svc.PublishExam(c.Request.Context(), examActor(c), examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// ListMockExams handles GET /mock-exams → 200, the published exam catalog.
func (h *ExamHandler) ListMockExams(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))

	published := true
	filter := repository.ExamFilter{IsPublished: &published}
	if s := c.Query("subject"); s != "" {
		filter.Subject = &s
	}
	if g := c.Query("grade"); g != "" {
		grade, convErr := strconv.Atoi(g)
		if convErr != nil {
			respondError(c, apperror.BadRequest("invalid grade", convErr))
			return
		}
		filter.Grade = &grade
	}

	exams, next, err := h.svc.ListExams(c.Request.Context(), examActor(c), filter, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	nextCursor := ""
	if next != nil {
		nextCursor = next.Encode()
	}
	respondJSON(c, http.StatusOK, dto.Paginated{
		Data:       dto.NewExamListResponse(exams),
		NextCursor: nextCursor,
		Limit:      limit,
	})
}

// Start handles POST /mock-exams/:id/start → 201 (questions + server timer).
func (h *ExamHandler) Start(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.StartExam(c.Request.Context(), examActor(c), examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewStartExamResponse(view)
	respondJSON(c, http.StatusCreated, &resp)
}

// SaveProgress handles PATCH /exam-attempts/:id → 204 (auto-save).
func (h *ExamHandler) SaveProgress(c *gin.Context) {
	attemptID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.SaveProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.SaveExamProgress(c.Request.Context(), examActor(c), attemptID, req.Answers); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Submit handles POST /exam-attempts/:id/submit → 200 (graded result).
func (h *ExamHandler) Submit(c *gin.Context) {
	attemptID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.SubmitExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	result, err := h.svc.SubmitExam(c.Request.Context(), examActor(c), attemptID, req.Answers)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResultResponse(result)
	respondJSON(c, http.StatusOK, &resp)
}

// GetResult handles GET /exam-attempts/:id/results → 200 (post-submit only).
func (h *ExamHandler) GetResult(c *gin.Context) {
	attemptID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	result, err := h.svc.GetExamResult(c.Request.Context(), examActor(c), attemptID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResultResponse(result)
	respondJSON(c, http.StatusOK, &resp)
}

// ListMyAttempts handles GET /mock-exams/:id/attempts or GET /exam-attempts/my → 200.
func (h *ExamHandler) ListMyAttempts(c *gin.Context) {
	examIDStr := c.Param("id")
	var examID uuid.UUID
	var err error
	if examIDStr != "" {
		examID, err = uuid.Parse(examIDStr)
		if err != nil {
			respondError(c, apperror.BadRequest("invalid exam id", err))
			return
		}
	}

	actor := examActor(c)
	var studentID uuid.UUID
	if actor != nil {
		studentID = actor.UserID
	} else {
		studentID = exam.GuestUserID
	}

	attempts, err := h.svc.ListAttemptsByStudent(c.Request.Context(), studentID, examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamAttemptHistoryListResponse(attempts)
	respondJSON(c, http.StatusOK, resp)
}
