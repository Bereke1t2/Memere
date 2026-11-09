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
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// QuizHandler adapts the quiz usecase to HTTP.
type QuizHandler struct {
	svc *quiz.Service
}

// NewQuizHandler builds a QuizHandler.
func NewQuizHandler(svc *quiz.Service) *QuizHandler {
	return &QuizHandler{svc: svc}
}

// quizActor converts the context actor into the quiz usecase's actor (nil for
// anonymous). Ownership and visibility are enforced inside the usecase.
func quizActor(c *gin.Context) *quiz.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return nil
	}
	return &quiz.Actor{UserID: a.UserID, Role: a.Role}
}

// ListByCourse handles GET /courses/:id/quizzes → 200.
func (h *QuizHandler) ListByCourse(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	quizzes, err := h.svc.ListByCourse(c.Request.Context(), courseID)
	if err != nil {
		respondError(c, err)
		return
	}
	out := make([]dto.QuizResponse, len(quizzes))
	for i, q := range quizzes {
		out[i] = dto.NewQuizResponse(q)
	}
	respondJSON(c, http.StatusOK, gin.H{"data": out})
}

// CreateQuiz handles POST /courses/:id/quizzes → 201.
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	created, err := h.svc.CreateQuiz(c.Request.Context(), quizActor(c), quiz.CreateQuizInput{
		CourseID:           &courseID,
		LessonID:           req.LessonID,
		Title:              req.Title,
		TimeLimitSeconds:   req.TimeLimitSeconds,
		PassPercentage:     req.PassPercentage,
		RandomizeQuestions: req.RandomizeQuestions,
		MaxAttempts:        req.MaxAttempts,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuizResponse(created)
	respondJSON(c, http.StatusCreated, &resp)
}

// AddQuestion handles POST /quizzes/:id/questions → 201. Answers are inline.
func (h *QuizHandler) AddQuestion(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.AddQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	answers := make([]quiz.AnswerInput, 0, len(req.Answers))
	for _, a := range req.Answers {
		answers = append(answers, quiz.AnswerInput{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: a.OrderIndex})
	}
	created, err := h.svc.AddQuestion(c.Request.Context(), quizActor(c), quizID, quiz.QuestionInput{
		Text:        req.Text,
		Type:        entity.QuestionType(req.Type),
		Points:      req.Points,
		Explanation: req.Explanation,
		OrderIndex:  req.OrderIndex,
		Subject:     req.Subject,
		Topic:       req.Topic,
		Answers:     answers,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuestionResponse(created)
	respondJSON(c, http.StatusCreated, &resp)
}

// UpdateQuiz handles PUT /quizzes/:id → 200.
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	updated, err := h.svc.UpdateQuiz(c.Request.Context(), quizActor(c), quizID, quiz.UpdateQuizInput{
		Title:              req.Title,
		TimeLimitSeconds:   req.TimeLimitSeconds,
		PassPercentage:     req.PassPercentage,
		RandomizeQuestions: req.RandomizeQuestions,
		MaxAttempts:        req.MaxAttempts,
		LessonID:           req.LessonID,
		IsPublished:        req.IsPublished,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuizResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// DeleteQuiz handles DELETE /quizzes/:id → 204.
func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.DeleteQuiz(c.Request.Context(), quizActor(c), quizID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetQuiz handles GET /quizzes/:id → 200 (student metadata, no answer keys).
func (h *QuizHandler) GetQuiz(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.GetQuizForStudent(c.Request.Context(), quizActor(c), quizID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuizClientResponse(view)
	respondJSON(c, http.StatusOK, &resp)
}

// StartAttempt handles POST /quizzes/:id/attempts → 201 (questions + timer).
func (h *QuizHandler) StartAttempt(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.StartAttempt(c.Request.Context(), quizActor(c), quizID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewStartAttemptResponse(view)
	respondJSON(c, http.StatusCreated, &resp)
}

// SaveProgress handles PATCH /quiz-attempts/:id → 204 (auto-save).
func (h *QuizHandler) SaveProgress(c *gin.Context) {
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
	if err := h.svc.SaveProgress(c.Request.Context(), quizActor(c), attemptID, req.Answers); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Submit handles POST /quiz-attempts/:id/submit → 200 (graded result).
func (h *QuizHandler) Submit(c *gin.Context) {
	attemptID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.SubmitAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	result, err := h.svc.SubmitAttempt(c.Request.Context(), quizActor(c), attemptID, req.Answers)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAttemptResultResponse(result)
	respondJSON(c, http.StatusOK, &resp)
}

// GetResult handles GET /quiz-attempts/:id/result → 200 (post-submit only).
func (h *QuizHandler) GetResult(c *gin.Context) {
	attemptID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	result, err := h.svc.GetAttemptResult(c.Request.Context(), quizActor(c), attemptID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAttemptResultResponse(result)
	respondJSON(c, http.StatusOK, &resp)
}

// GetForDownload handles GET /quizzes/:id/download → 200. It returns the full
// quiz WITH answer keys so a downloaded quiz can be graded on-device offline.
// This is the sanctioned Non-Negotiable #1 exception (see dto/download.go): it is
// gated exactly like taking the quiz (GetQuizForDownload → assertCourseAccess),
// and the normal GET /quizzes/:id path stays answer-free.
func (h *QuizHandler) GetForDownload(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	d, err := h.svc.GetQuizForDownload(c.Request.Context(), quizActor(c), quizID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuizDownloadResponse(d)
	respondJSON(c, http.StatusOK, &resp)
}

// CreateStandaloneQuiz handles POST /quizzes → 201. Unlike the nested
// POST /courses/:id/quizzes route the course comes from the body and may be
// omitted entirely, which creates a standalone (course-less) quiz — admin only.
func (h *QuizHandler) CreateStandaloneQuiz(c *gin.Context) {
	var req dto.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	created, err := h.svc.CreateQuiz(c.Request.Context(), quizActor(c), quiz.CreateQuizInput{
		CourseID:           req.CourseID,
		LessonID:           req.LessonID,
		Title:              req.Title,
		TimeLimitSeconds:   req.TimeLimitSeconds,
		PassPercentage:     req.PassPercentage,
		RandomizeQuestions: req.RandomizeQuestions,
		MaxAttempts:        req.MaxAttempts,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuizResponse(created)
	respondJSON(c, http.StatusCreated, &resp)
}

// ListQuizzes handles GET /quizzes → 200: the cross-course authoring list with
// cursor pagination. Supported filters: course_id, standalone, published.
func (h *QuizHandler) ListQuizzes(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))

	var filter repository.QuizFilter
	if raw := c.Query("course_id"); raw != "" {
		courseID, parseErr := uuid.Parse(raw)
		if parseErr != nil {
			respondError(c, apperror.BadRequest("invalid course_id", parseErr))
			return
		}
		filter.CourseID = &courseID
	}
	if raw := c.Query("standalone"); raw != "" {
		standalone, parseErr := strconv.ParseBool(raw)
		if parseErr != nil {
			respondError(c, apperror.BadRequest("invalid standalone", parseErr))
			return
		}
		filter.Standalone = &standalone
	}
	if raw := c.Query("published"); raw != "" {
		published, parseErr := strconv.ParseBool(raw)
		if parseErr != nil {
			respondError(c, apperror.BadRequest("invalid published", parseErr))
			return
		}
		filter.IsPublished = &published
	}

	quizzes, counts, next, err := h.svc.ListQuizzes(c.Request.Context(), quizActor(c), filter, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	nextCursor := ""
	if next != nil {
		nextCursor = next.Encode()
	}
	respondJSON(c, http.StatusOK, gin.H{
		"data":        dto.NewQuizListResponse(quizzes, counts),
		"next_cursor": nextCursor,
	})
}

// GetQuizForAuthoring handles GET /quizzes/:id/questions → 200. The response
// carries answer keys, so this route MUST stay behind requireAuth +
// teacherOrAdmin; the usecase additionally proves course ownership.
func (h *QuizHandler) GetQuizForAuthoring(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.GetQuizForAuthoring(c.Request.Context(), quizActor(c), quizID)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAuthoredQuizResponse(view)
	respondJSON(c, http.StatusOK, &resp)
}

// UpdateQuestion handles PUT /questions/:id → 200.
func (h *QuizHandler) UpdateQuestion(c *gin.Context) {
	questionID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}

	in := quiz.UpdateQuestionInput{
		Text:        req.Text,
		Points:      req.Points,
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
		answers := make([]quiz.AnswerInput, 0, len(*req.Answers))
		for i, a := range *req.Answers {
			order := a.OrderIndex
			if order == 0 {
				// A client that omits order_index gets the array order, so a
				// replacement set never collapses to all-zero indexes.
				order = i
			}
			answers = append(answers, quiz.AnswerInput{Text: a.Text, IsCorrect: a.IsCorrect, OrderIndex: order})
		}
		in.Answers = answers
	}

	updated, err := h.svc.UpdateQuestion(c.Request.Context(), quizActor(c), questionID, in)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewAuthoredQuizQuestionResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// DeleteQuestion handles DELETE /questions/:id → 204.
func (h *QuizHandler) DeleteQuestion(c *gin.Context) {
	questionID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.DeleteQuestion(c.Request.Context(), quizActor(c), questionID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ReorderQuestions handles PUT /quizzes/:id/questions/order → 204.
func (h *QuizHandler) ReorderQuestions(c *gin.Context) {
	quizID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.ReorderQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	if err := h.svc.ReorderQuestions(c.Request.Context(), quizActor(c), quizID, req.QuestionIDs); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ImportQuiz handles POST /quizzes/import → 201: a whole quiz and its question
// set from one uploaded document, created atomically.
func (h *QuizHandler) ImportQuiz(c *gin.Context) {
	var req dto.ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid import document", err))
		return
	}
	in, err := req.ToQuizImport()
	if err != nil {
		respondError(c, err)
		return
	}
	created, count, err := h.svc.ImportQuiz(c.Request.Context(), quizActor(c), in)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusCreated, &dto.ImportResponse{
		Kind:          "quiz",
		ID:            created.ID,
		Title:         created.Title,
		QuestionCount: count,
	})
}
