package quiz

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// untimedStateTTL is the Redis TTL applied to attempt state for quizzes with no
// time limit. Timed quizzes use time_limit + stateGrace instead.
const (
	untimedStateTTL = 2 * time.Hour
	stateGrace      = 60 * time.Second
	// shuffle salts keep the question-order and answer-order permutations
	// independent within a single attempt.
	questionSalt = 1
	answerSalt   = 2
)

// CourseAccess is the access gate for graded work. *access.Service satisfies
// it; the quiz engine requires FullAccess (preview is for watching sample
// lessons, never for graded attempts).
type CourseAccess interface {
	RequireFullAccess(ctx context.Context, actor access.Actor, courseID uuid.UUID) error
}

// Service implements the quiz-engine usecases over the domain repositories plus
// the Redis attempt-state store. now is injectable so tests can drive the
// server-side timer with a fake clock.
type Service struct {
	quizzes   repository.QuizRepository
	questions repository.QuestionRepository
	attempts  repository.QuizAttemptRepository
	courses   repository.CourseRepository
	state     repository.AttemptStateStore
	tx        repository.TxManager
	access    CourseAccess
	now       func() time.Time
}

// NewService wires the quiz usecase with its dependencies.
func NewService(
	quizzes repository.QuizRepository,
	questions repository.QuestionRepository,
	attempts repository.QuizAttemptRepository,
	courses repository.CourseRepository,
	state repository.AttemptStateStore,
	tx repository.TxManager,
	accessSvc CourseAccess,
) *Service {
	return &Service{
		quizzes:   quizzes,
		questions: questions,
		attempts:  attempts,
		courses:   courses,
		state:     state,
		tx:        tx,
		access:    accessSvc,
		now:       time.Now,
	}
}

// ListByCourse returns all quizzes for a course. Used by the teacher authoring UI.
func (s *Service) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Quiz, error) {
	return s.quizzes.ListByCourse(ctx, courseID)
}

// GetQuizForStudent returns quiz metadata and the question count (not the
// questions). Access is gated on FullAccess to the parent course (owner/admin,
// free course, active enrollment, or active subscription) via access.Service.
func (s *Service) GetQuizForStudent(ctx context.Context, actor *Actor, quizID uuid.UUID) (*QuizClientView, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCourseAccess(ctx, actor, quiz.CourseID); err != nil {
		return nil, err
	}

	questions, err := s.quizzes.GetQuestionsForClient(ctx, quizID)
	if err != nil {
		return nil, err
	}
	used, err := s.attempts.CountByStudentAndQuiz(ctx, actor.UserID, quizID)
	if err != nil {
		return nil, err
	}

	return &QuizClientView{
		ID:                 quiz.ID,
		CourseID:           quiz.CourseID,
		Title:              quiz.Title,
		TimeLimitSeconds:   quiz.TimeLimitSeconds,
		PassPercentage:     quiz.PassPercentage,
		RandomizeQuestions: quiz.RandomizeQuestions,
		MaxAttempts:        quiz.MaxAttempts,
		QuestionCount:      len(questions),
		AttemptsUsed:       used,
	}, nil
}

// StartAttempt begins (or resumes) a student's attempt. It is idempotent: if an
// in-progress attempt already exists it is returned with the remaining time,
// rather than consuming another of the max_attempts budget. Otherwise it
// enforces max_attempts, creates the Postgres attempt, builds and persists the
// randomized snapshot (Redis + the durable question_order column), and returns
// the questions via the answer-key-free client view.
func (s *Service) StartAttempt(ctx context.Context, actor *Actor, quizID uuid.UUID) (*AttemptClientView, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	quiz, err := s.quizzes.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if err := s.assertCourseAccess(ctx, actor, quiz.CourseID); err != nil {
		return nil, err
	}

	// Idempotent start: resume an existing in-progress attempt.
	if existing, err := s.attempts.GetActive(ctx, actor.UserID, quizID); err == nil {
		clientQs, err := s.quizzes.GetQuestionsForClient(ctx, quizID)
		if err != nil {
			return nil, err
		}
		ordered := s.applyOrder(existing, quiz, clientQs)
		return s.attemptView(existing, quizID, ordered), nil
	} else if !apperror.IsNotFound(err) {
		return nil, err
	}

	// Enforce attempt limits (§9.1).
	used, err := s.attempts.CountByStudentAndQuiz(ctx, actor.UserID, quizID)
	if err != nil {
		return nil, err
	}
	if quiz.MaxAttempts != nil && used >= *quiz.MaxAttempts {
		return nil, apperror.Forbidden("maximum attempts reached for this quiz", nil)
	}

	clientQs, err := s.quizzes.GetQuestionsForClient(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if len(clientQs) == 0 {
		return nil, apperror.BadRequest("quiz has no questions", nil)
	}

	now := s.now()
	var expiresAt *time.Time
	if quiz.TimeLimitSeconds != nil {
		t := now.Add(time.Duration(*quiz.TimeLimitSeconds) * time.Second)
		expiresAt = &t
	}

	attempt := &entity.QuizAttempt{
		QuizID:        quizID,
		StudentID:     actor.UserID,
		AttemptNumber: used + 1,
		StartedAt:     now,
		ExpiresAt:     expiresAt,
		Status:        entity.AttemptInProgress,
	}
	if err := s.attempts.Create(ctx, attempt); err != nil {
		return nil, err
	}

	// Build the randomized order now that we have the attempt ID, persist it to
	// the durable column, and mirror it into Redis for fast live reads.
	order := buildOrder(attempt.ID, quiz, clientQs)
	attempt.QuestionOrder = order
	if err := s.attempts.Update(ctx, attempt); err != nil {
		return nil, err
	}
	_ = s.state.SetSnapshot(ctx, attempt.ID, order, s.stateTTL(quiz))

	ordered := orderQuestions(clientQs, order)
	return s.attemptView(attempt, quizID, ordered), nil
}

// SaveProgress records the student's in-progress answers to Redis (auto-save,
// called ~every 30s). It asserts ownership, that the attempt is in progress, and
// that the server-side deadline has not passed.
func (s *Service) SaveProgress(ctx context.Context, actor *Actor, attemptID uuid.UUID, answers map[string]any) error {
	if actor == nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return err
	}
	if attempt.Status != entity.AttemptInProgress {
		return apperror.Conflict("attempt is not in progress", nil)
	}
	if s.expired(attempt) {
		return apperror.Conflict("attempt has expired", nil)
	}
	quiz, err := s.quizzes.FindByID(ctx, attempt.QuizID)
	if err != nil {
		return err
	}
	return s.state.SaveAnswers(ctx, attemptID, answers, s.stateTTL(quiz))
}

// SubmitAttempt grades the attempt server-side and returns the result with
// feedback. It merges the final payload over any Redis-saved answers, grades
// using the server-only grading tree (never the client-reported score), persists
// the snapshot + score/percentage/passed, transitions status to graded, and
// clears the Redis state. If the deadline has already passed the attempt is
// marked expired and graded with whatever answers exist (§9.2). If the background
// sweeper finalized it first, the already-graded result is returned.
func (s *Service) SubmitAttempt(ctx context.Context, actor *Actor, attemptID uuid.UUID, answers map[string]any) (*AttemptResult, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return nil, err
	}
	// Already finalized (e.g. the sweeper graded an expired attempt first): return
	// the existing result rather than erroring, so a late client submit is benign.
	if attempt.Status != entity.AttemptInProgress {
		return s.GetAttemptResult(ctx, actor, attemptID)
	}

	saved, _ := s.state.GetAnswers(ctx, attemptID)
	merged := mergeAnswers(saved, answers)

	result, claimed, err := s.finalize(ctx, attempt, merged)
	if err != nil {
		return nil, err
	}
	if !claimed {
		// The sweeper (or a concurrent submit) graded it first — return that result.
		return s.GetAttemptResult(ctx, actor, attemptID)
	}
	return result, nil
}

// finalize grades an in-progress attempt and transitions it to graded, guarded so
// that only one caller (this submit OR the sweeper) ever grades it. The first
// writer to flip the row out of in_progress (ClaimForGrading) wins; a loser
// returns claimed=false and grades nothing. Returns the result on success.
func (s *Service) finalize(ctx context.Context, attempt *entity.QuizAttempt, answers map[string]any) (*AttemptResult, bool, error) {
	quiz, err := s.quizzes.FindByID(ctx, attempt.QuizID)
	if err != nil {
		return nil, false, err
	}
	tree, err := s.quizzes.GetQuizWithQuestions(ctx, attempt.QuizID)
	if err != nil {
		return nil, false, err
	}

	sub := parseSubmitted(tree.Questions, answers)
	result := gradeQuiz(quiz, tree.Questions, sub)
	result.AttemptID = attempt.ID
	result.AttemptNumber = attempt.AttemptNumber

	target := entity.AttemptSubmitted
	if s.expired(attempt) {
		target = entity.AttemptExpired
	}

	now := s.now()
	attempt.AnswersSnapshot = snapshotFor(answers, result)
	attempt.SubmittedAt = &now
	attempt.Status = target

	claimed, err := s.attempts.ClaimForGrading(ctx, attempt)
	if err != nil {
		return nil, false, err
	}
	if !claimed {
		return nil, false, nil
	}

	score := result.Score
	pct := result.Percentage
	passed := result.Passed
	attempt.Score = &score
	attempt.Percentage = &pct
	attempt.Passed = &passed
	if err := s.attempts.Grade(ctx, attempt); err != nil {
		return nil, false, err
	}
	result.Status = entity.AttemptGraded
	result.SubmittedAt = attempt.SubmittedAt

	_ = s.state.DeleteAttemptState(ctx, attempt.ID)
	return &result, true, nil
}

// SweepExpired finalizes abandoned (in-progress, past-deadline) attempts so the
// §9.2 "server timer fires" rule holds even when no client request arrives. It
// is called by the background worker (Skill 4). Each attempt is graded with its
// last Redis-saved answers (the PG snapshot is empty until finalize). Already-
// finalized rows are skipped via the guarded claim, so a race with a late client
// submit grades exactly once. Returns the number of attempts graded this pass.
func (s *Service) SweepExpired(ctx context.Context, now time.Time, limit int) (int, error) {
	expired, err := s.attempts.ListExpired(ctx, now, limit)
	if err != nil {
		return 0, err
	}
	graded := 0
	for _, attempt := range expired {
		saved, _ := s.state.GetAnswers(ctx, attempt.ID)
		_, claimed, err := s.finalize(ctx, attempt, saved)
		if err != nil {
			return graded, err
		}
		if claimed {
			graded++
		}
	}
	return graded, nil
}

// GetAttemptResult returns the graded result for one of the actor's own
// attempts. It is available only once the attempt has left in_progress.
func (s *Service) GetAttemptResult(ctx context.Context, actor *Actor, attemptID uuid.UUID) (*AttemptResult, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return nil, err
	}
	if attempt.Status == entity.AttemptInProgress {
		return nil, apperror.Conflict("attempt is still in progress", nil)
	}

	quiz, err := s.quizzes.FindByID(ctx, attempt.QuizID)
	if err != nil {
		return nil, err
	}
	tree, err := s.quizzes.GetQuizWithQuestions(ctx, attempt.QuizID)
	if err != nil {
		return nil, err
	}
	// The persisted snapshot wraps the answers under "answers"; grade those, not
	// the wrapper.
	sub := parseSubmitted(tree.Questions, answersFromSnapshot(attempt.AnswersSnapshot))
	result := gradeQuiz(quiz, tree.Questions, sub)
	result.AttemptID = attempt.ID
	result.AttemptNumber = attempt.AttemptNumber
	result.Status = attempt.Status
	result.SubmittedAt = attempt.SubmittedAt
	return &result, nil
}

// ListMyAttempts returns the actor's own attempt history for a quiz.
func (s *Service) ListMyAttempts(ctx context.Context, actor *Actor, quizID uuid.UUID) ([]*entity.QuizAttempt, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.attempts.ListByStudentAndQuiz(ctx, actor.UserID, quizID)
}

// assertCourseAccess gates graded quiz work on FullAccess to the parent course
// (owner/admin, free course, active enrollment, or active subscription) via the
// shared access.Service. PreviewAccess is deliberately insufficient: previews
// are for watching sample lessons, never for graded attempts.
func (s *Service) assertCourseAccess(ctx context.Context, actor *Actor, courseID uuid.UUID) error {
	return s.access.RequireFullAccess(ctx, access.Actor{UserID: actor.UserID, Role: actor.Role}, courseID)
}

// expired reports whether a timed attempt is past its server-side deadline.
func (s *Service) expired(a *entity.QuizAttempt) bool {
	return a.ExpiresAt != nil && s.now().After(*a.ExpiresAt)
}

// stateTTL is the Redis TTL for an attempt's live state.
func (s *Service) stateTTL(quiz *entity.Quiz) time.Duration {
	if quiz.TimeLimitSeconds == nil {
		return untimedStateTTL
	}
	return time.Duration(*quiz.TimeLimitSeconds)*time.Second + stateGrace
}

// applyOrder reconstructs the persisted order for a resumed attempt, falling
// back to a freshly-derived deterministic order if neither Redis nor the durable
// column has it.
func (s *Service) applyOrder(attempt *entity.QuizAttempt, quiz *entity.Quiz, qs []repository.ClientQuestion) []QuestionClientView {
	order := attempt.QuestionOrder
	if order == nil {
		if snap, err := s.state.GetSnapshot(context.Background(), attempt.ID); err == nil && snap != nil {
			order = snap
		}
	}
	if order == nil {
		order = buildOrder(attempt.ID, quiz, qs)
	}
	return orderQuestions(qs, order)
}

// attemptView assembles the client view for an attempt, computing remaining time
// against the server clock.
func (s *Service) attemptView(attempt *entity.QuizAttempt, quizID uuid.UUID, questions []QuestionClientView) *AttemptClientView {
	var remaining *int
	if attempt.ExpiresAt != nil {
		secs := int(attempt.ExpiresAt.Sub(s.now()).Seconds())
		if secs < 0 {
			secs = 0
		}
		remaining = &secs
	}
	return &AttemptClientView{
		AttemptID:        attempt.ID,
		QuizID:           quizID,
		AttemptNumber:    attempt.AttemptNumber,
		Status:           attempt.Status,
		StartedAt:        attempt.StartedAt,
		ExpiresAt:        attempt.ExpiresAt,
		RemainingSeconds: remaining,
		Questions:        questions,
	}
}
