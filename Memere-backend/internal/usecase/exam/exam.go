package exam

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/grading"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

const stateGrace = 60 * time.Second

// Service implements the exam-engine usecases over the domain repositories plus
// the shared Redis attempt-state store. now is injectable so tests drive the
// server-side timer with a fake clock.
type Service struct {
	exams    repository.ExamRepository
	attempts repository.ExamAttemptRepository
	courses  repository.CourseRepository
	state    repository.AttemptStateStore
	tx       repository.TxManager
	now      func() time.Time
}

// NewService wires the exam usecase with its dependencies.
func NewService(
	exams repository.ExamRepository,
	attempts repository.ExamAttemptRepository,
	courses repository.CourseRepository,
	state repository.AttemptStateStore,
	tx repository.TxManager,
) *Service {
	return &Service{
		exams:    exams,
		attempts: attempts,
		courses:  courses,
		state:    state,
		tx:       tx,
		now:      time.Now,
	}
}

// StartExam begins (or resumes) a student's exam sitting. It is idempotent: an
// in-progress attempt still within its window is resumed with remaining time. An
// in-progress attempt found past its deadline is finalized (expired→graded) per
// the lazy-expiry rule before a fresh attempt is created.
func (s *Service) StartExam(ctx context.Context, actor *Actor, examID uuid.UUID) (*ExamAttemptClientView, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	exam, err := s.exams.FindByID(ctx, examID)
	if err != nil {
		return nil, err
	}
	if err := s.assertExamAccess(ctx, actor, exam); err != nil {
		return nil, err
	}

	if existing, err := s.attempts.GetActive(ctx, actor.UserID, examID); err == nil {
		if s.expired(exam, existing) {
			if _, err := s.finalizeExpired(ctx, exam, existing); err != nil {
				return nil, err
			}
		} else {
			return s.attemptView(ctx, exam, existing)
		}
	} else if !apperror.IsNotFound(err) {
		return nil, err
	}

	now := s.now()
	expiresAt := now.Add(time.Duration(exam.DurationMinutes) * time.Minute)
	attempt := &entity.ExamAttempt{
		ExamID:    examID,
		StudentID: actor.UserID,
		StartedAt: now,
		Status:    entity.AttemptInProgress,
	}
	if err := s.attempts.Create(ctx, attempt); err != nil {
		return nil, err
	}

	// Persist the presentation order (teacher-fixed order_index for exams) into
	// Redis for fast resume; the order is reproducible from the DB regardless.
	clientQs, err := s.exams.GetQuestionsForClient(ctx, examID)
	if err != nil {
		return nil, err
	}
	_ = s.state.SetSnapshot(ctx, attempt.ID, orderSnapshot(clientQs), s.stateTTL(exam))

	// attempt.ExpiresAt is derived (exam_attempts has no expires_at column); set it
	// on the in-memory entity so the view can report remaining time.
	attempt.SubmittedAt = nil
	view := s.buildView(exam, attempt, &expiresAt, clientQs)
	return view, nil
}

// SaveExamProgress records in-progress answers to Redis (auto-save). If the
// deadline has passed it lazily finalizes the attempt (expired→graded) and
// reports the expiry rather than accepting more answers.
func (s *Service) SaveExamProgress(ctx context.Context, actor *Actor, attemptID uuid.UUID, answers map[string]any) error {
	if actor == nil {
		return apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return err
	}
	if attempt.Status != entity.AttemptInProgress {
		return invalidState("attempt is not in progress")
	}
	exam, err := s.exams.FindByID(ctx, attempt.ExamID)
	if err != nil {
		return err
	}
	if s.expired(exam, attempt) {
		if _, err := s.finalizeExpired(ctx, exam, attempt); err != nil {
			return err
		}
		return invalidState("attempt has expired")
	}
	return s.state.SaveAnswers(ctx, attemptID, answers, s.stateTTL(exam))
}

// SubmitExam grades an attempt server-side and returns the result. Past the
// deadline the attempt goes expired→graded; otherwise submitted→graded. Every
// transition is validated through the §9.2 state machine.
func (s *Service) SubmitExam(ctx context.Context, actor *Actor, attemptID uuid.UUID, answers map[string]any) (*ExamResult, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return nil, err
	}
	if attempt.Status != entity.AttemptInProgress {
		return nil, invalidState("attempt is not in progress")
	}
	exam, err := s.exams.FindByID(ctx, attempt.ExamID)
	if err != nil {
		return nil, err
	}

	saved, _ := s.state.GetAnswers(ctx, attemptID)
	merged := mergeAnswers(saved, answers)

	target := entity.AttemptSubmitted
	if s.expired(exam, attempt) {
		target = entity.AttemptExpired
	}
	return s.finalize(ctx, exam, attempt, merged, target)
}

// GetExamResult returns the graded result for one of the actor's own attempts.
// It is available only in a terminal state; an attempt found past its deadline
// while still in progress is finalized first (lazy expiry).
func (s *Service) GetExamResult(ctx context.Context, actor *Actor, attemptID uuid.UUID) (*ExamResult, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, actor.UserID)
	if err != nil {
		return nil, err
	}
	exam, err := s.exams.FindByID(ctx, attempt.ExamID)
	if err != nil {
		return nil, err
	}

	if attempt.Status == entity.AttemptInProgress {
		if !s.expired(exam, attempt) {
			return nil, invalidState("attempt is still in progress")
		}
		return s.finalizeExpired(ctx, exam, attempt)
	}
	return s.gradeView(ctx, exam, attempt)
}

// ListMyExamAttempts returns the actor's own attempt history for an exam.
func (s *Service) ListMyExamAttempts(ctx context.Context, actor *Actor, examID uuid.UUID) ([]*entity.ExamAttempt, error) {
	if actor == nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.attempts.ListByStudent(ctx, actor.UserID)
}

// finalizeExpired drives an in-progress, past-deadline attempt through
// expired→graded using whatever answers were auto-saved.
func (s *Service) finalizeExpired(ctx context.Context, exam *entity.Exam, attempt *entity.ExamAttempt) (*ExamResult, error) {
	saved, _ := s.state.GetAnswers(ctx, attempt.ID)
	return s.finalize(ctx, exam, attempt, mergeAnswers(saved, nil), entity.AttemptExpired)
}

// finalize performs the in_progress → {submitted|expired} → graded sequence,
// validating each step against the state machine, grading server-side, and
// persisting the snapshot and score. Redis state is cleared at the end.
func (s *Service) finalize(ctx context.Context, exam *entity.Exam, attempt *entity.ExamAttempt, answers map[string]any, target entity.AttemptStatus) (*ExamResult, error) {
	if err := guard(attempt.Status, target); err != nil {
		return nil, err
	}

	gradables, err := s.loadGradables(ctx, exam.ID)
	if err != nil {
		return nil, err
	}
	sub := parseSubmitted(gradables, answers)
	core := grading.Grade(gradables, sub)

	now := s.now()
	attempt.Status = target
	attempt.SubmittedAt = &now
	attempt.AnswersSnapshot = snapshotFor(answers, core)
	if err := s.attempts.Update(ctx, attempt); err != nil {
		return nil, err
	}

	if err := guard(attempt.Status, entity.AttemptGraded); err != nil {
		return nil, err
	}
	score := core.Score
	pct := core.Percentage
	attempt.Score = &score
	attempt.Percentage = &pct
	if err := s.attempts.Grade(ctx, attempt); err != nil {
		return nil, err
	}
	attempt.Status = entity.AttemptGraded

	_ = s.state.DeleteAttemptState(ctx, attempt.ID)

	res := s.result(exam, attempt, core)
	res.SubmittedAt = &now
	return res, nil
}

// guard validates a state transition, returning INVALID_ATTEMPT_STATE on an
// illegal move (§9.2 via entity.AttemptStatus.CanTransitionTo).
func guard(from, to entity.AttemptStatus) error {
	if !from.CanTransitionTo(to) {
		return invalidState("illegal attempt state transition")
	}
	return nil
}

func invalidState(msg string) error {
	return apperror.New(http.StatusConflict, "INVALID_ATTEMPT_STATE", msg, nil)
}

// expired reports whether an attempt is past its server-side deadline
// (started_at + duration_minutes). The client timer is display-only
// (Non-Negotiable #2); this server clock is the authority.
func (s *Service) expired(exam *entity.Exam, a *entity.ExamAttempt) bool {
	return s.now().After(a.StartedAt.Add(time.Duration(exam.DurationMinutes) * time.Minute))
}

func (s *Service) stateTTL(exam *entity.Exam) time.Duration {
	return time.Duration(exam.DurationMinutes)*time.Minute + stateGrace
}

// assertExamAccess enforces the published-or-owner gate.
// TODO(phase4): replace with an enrollment check.
func (s *Service) assertExamAccess(ctx context.Context, actor *Actor, exam *entity.Exam) error {
	if exam.IsPublished {
		return nil
	}
	if exam.CourseID != nil {
		course, err := s.courses.FindByID(ctx, *exam.CourseID)
		if err != nil {
			return err
		}
		if actor.ownsCourse(course) {
			return nil
		}
	} else if actor.Role == entity.RoleAdmin {
		return nil
	}
	return apperror.Forbidden("you do not have access to this exam", nil)
}
