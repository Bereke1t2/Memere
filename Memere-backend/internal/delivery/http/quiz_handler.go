package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
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
		CourseID:           courseID,
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
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewQuizResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
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
