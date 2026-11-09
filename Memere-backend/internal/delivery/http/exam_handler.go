package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
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
// question by id or creates an inline question; returns updated exam).
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

	var updated *entity.Exam
	if req.QuestionID != nil && *req.QuestionID != uuid.Nil {
		updated, err = h.svc.AddExamQuestion(c.Request.Context(), examActor(c), examID, *req.QuestionID, req.Marks, req.OrderIndex)
	} else if req.Text != nil && *req.Text != "" {
		answers := make([]exam.InlineAnswerInput, 0, len(req.Answers))
		for _, a := range req.Answers {
			answers = append(answers, exam.InlineAnswerInput{
				Text:       a.Text,
				IsCorrect:  a.IsCorrect,
				OrderIndex: a.OrderIndex,
			})
		}
		qType := entity.QuestionMultipleChoice
		if req.Type != nil && *req.Type != "" {
			qType = entity.QuestionType(*req.Type)
		}
		updated, err = h.svc.AddInlineExamQuestion(c.Request.Context(), examActor(c), examID, exam.AddInlineQuestionInput{
			Text:        *req.Text,
			Type:        qType,
			Marks:       req.Marks,
			Explanation: req.Explanation,
			OrderIndex:  req.OrderIndex,
			Subject:     req.Subject,
			Topic:       req.Topic,
			Answers:     answers,
		})
	} else {
		respondError(c, apperror.BadRequest("either question_id or text is required", nil))
		return
	}

	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResponse(updated)
	respondJSON(c, http.StatusCreated, &resp)
}

// UpdateExam handles PUT /exams/:id → 200.
func (h *ExamHandler) UpdateExam(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.UpdateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	updated, err := h.svc.UpdateExam(c.Request.Context(), examActor(c), examID, exam.UpdateExamInput{
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
	resp := dto.NewExamResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// DeleteExam handles DELETE /exams/:id → 204.
func (h *ExamHandler) DeleteExam(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.DeleteExam(c.Request.Context(), examActor(c), examID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
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

// GetForDownload handles GET /mock-exams/:id/download → 200. It returns the full
// exam WITH answer keys so a downloaded exam can be taken and graded on-device
// offline. This is the sanctioned Non-Negotiable #1 exception (see
// dto/download.go): it is gated exactly like starting the exam
// (GetExamForDownload → assertExamAccess), and the normal sitting path stays
// answer-free and server-graded.
func (h *ExamHandler) GetForDownload(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	d, err := h.svc.GetExamForDownload(c.Request.Context(), examActor(c), examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamDownloadResponse(d)
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

// CreateStandaloneExam handles POST /exams → 201. Unlike the nested
// POST /courses/:id/exams route the course comes from the body and may be
// omitted entirely, which creates a standalone (course-less) exam — admin only.
func (h *ExamHandler) CreateStandaloneExam(c *gin.Context) {
	var req dto.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	created, err := h.svc.CreateExam(c.Request.Context(), examActor(c), exam.CreateExamInput{
		CourseID:        req.CourseID,
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

// ListExamsForAuthoring handles GET /exams → 200: the authoring catalog with
// cursor pagination, including the caller's unpublished drafts. The public
// student-facing catalog stays on GET /mock-exams, which forces is_published.
func (h *ExamHandler) ListExamsForAuthoring(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))

	var filter repository.ExamFilter
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
	if raw := c.Query("published"); raw != "" {
		published, convErr := strconv.ParseBool(raw)
		if convErr != nil {
			respondError(c, apperror.BadRequest("invalid published", convErr))
			return
		}
		filter.IsPublished = &published
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
	respondJSON(c, http.StatusOK, gin.H{
		"data":        dto.NewExamListResponse(exams),
		"next_cursor": nextCursor,
	})
}

// GetExamForAuthoring handles GET /exams/:id/questions → 200. The response
// carries answer keys, so this route MUST stay behind requireAuth +
// teacherOrAdmin; the usecase additionally proves exam ownership.
func (h *ExamHandler) GetExamForAuthoring(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.GetExamForAuthoring(c.Request.Context(), examActor(c), examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAuthoredExamResponse(view)
	respondJSON(c, http.StatusOK, &resp)
}

// UpdateExamQuestion handles PUT /exams/:id/questions/:questionId → 200.
func (h *ExamHandler) UpdateExamQuestion(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	questionID, err := parseUUIDParam(c, "questionId")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.UpdateExamQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}

	in := exam.UpdateExamQuestionInput{
		Text:        req.Text,
		Marks:       req.Marks,
		Explanation: req.Explanation,
		OrderIndex:  req.OrderIndex,
		Subject:     req.Subject,
		Topic:       req.Topic,
	}
	if req.Type != nil {
		qType := entity.QuestionType(*req.Type)
		in.Type = &qType
	}
	if req.Answers != nil {
		answers := make([]exam.InlineAnswerInput, 0, len(*req.Answers))
		for i, a := range *req.Answers {
			order := a.OrderIndex
			if order == 0 {
				// A client that omits order_index gets the array order, so a
				// replacement set never collapses to all-zero indexes.
				order = i
			}
			answers = append(answers, exam.InlineAnswerInput{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: order})
		}
		in.Answers = answers
	}

	updated, err := h.svc.UpdateExamQuestion(c.Request.Context(), examActor(c), examID, questionID, in)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAuthoredExamQuestionResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// RemoveExamQuestion handles DELETE /exams/:id/questions/:questionId → 200,
// returning the exam with its re-derived total_marks.
func (h *ExamHandler) RemoveExamQuestion(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	questionID, err := parseUUIDParam(c, "questionId")
	if err != nil {
		respondError(c, err)
		return
	}
	updated, err := h.svc.RemoveExamQuestion(c.Request.Context(), examActor(c), examID, questionID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// ReorderExamQuestions handles PUT /exams/:id/questions/order → 204.
func (h *ExamHandler) ReorderExamQuestions(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.ReorderQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.ReorderExamQuestions(c.Request.Context(), examActor(c), examID, req.QuestionIDs); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Unpublish handles POST /exams/:id/unpublish → 200, hiding an exam from
// students again so a published exam can be corrected.
func (h *ExamHandler) Unpublish(c *gin.Context) {
	examID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	updated, err := h.svc.UnpublishExam(c.Request.Context(), examActor(c), examID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewExamResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// ImportExam handles POST /exams/import → 201: a whole exam and its question
// set from one uploaded document, created atomically with a derived total_marks.
func (h *ExamHandler) ImportExam(c *gin.Context) {
	var req dto.ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid import document", err))
		return
	}
	in, err := req.ToExamImport()
	if err != nil {
		respondError(c, err)
		return
	}
	created, count, err := h.svc.ImportExam(c.Request.Context(), examActor(c), in)
	if err != nil {
		respondError(c, err)
		return
	}
	total := created.TotalMarks
	respondJSON(c, http.StatusCreated, &dto.ImportResponse{
		Kind:          "exam",
		ID:            created.ID,
		Title:         created.Title,
		QuestionCount: count,
		TotalMarks:    &total,
	})
}
