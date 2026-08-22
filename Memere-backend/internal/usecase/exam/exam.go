package exam

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/grading"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

const stateGrace = 60 * time.Second

// CourseAccess is the access gate for graded exam work. *access.Service
// satisfies it; an exam tied to a course requires FullAccess (preview is for
// watching sample lessons, never for a graded sitting).
type CourseAccess interface {
	RequireFullAccess(ctx context.Context, actor access.Actor, courseID uuid.UUID) error
}

// Notifier fires async notifications after an exam is graded. Implementations
// must never block or fail the API request path.
type Notifier interface {
	ExamGraded(ctx context.Context, studentID, attemptID uuid.UUID, score float64)
}

// noopExamNotifier is the default when no notifier is wired.
type noopExamNotifier struct{}

func (noopExamNotifier) ExamGraded(context.Context, uuid.UUID, uuid.UUID, float64) {}

// Service implements the exam-engine usecases over the domain repositories plus
// the shared Redis attempt-state store. now is injectable so tests drive the
// server-side timer with a fake clock.
type Service struct {
	exams    repository.ExamRepository
	attempts repository.ExamAttemptRepository
	courses  repository.CourseRepository
	state    repository.AttemptStateStore
	ranking  repository.ScoreRanking
	tx       repository.TxManager
	access   CourseAccess
	notify   Notifier
	now      func() time.Time
}

// NewService wires the exam usecase with its dependencies. ranking and notify
// may be nil (e.g. in tests that don't exercise percentile/notifications).
func NewService(
	exams repository.ExamRepository,
	attempts repository.ExamAttemptRepository,
	courses repository.CourseRepository,
	state repository.AttemptStateStore,
	ranking repository.ScoreRanking,
	tx repository.TxManager,
	accessSvc CourseAccess,
	notify Notifier,
) *Service {
	if notify == nil {
		notify = noopExamNotifier{}
	}
	return &Service{
		exams:    exams,
		attempts: attempts,
		courses:  courses,
		state:    state,
		ranking:  ranking,
		tx:       tx,
		access:   accessSvc,
		notify:   notify,
		now:      time.Now,
	}
}

// ListByCourse returns all exams for a course. Used by the teacher authoring UI.
func (s *Service) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Exam, error) {
	return s.exams.ListByCourse(ctx, courseID)
}

// GuestUserID is the system identifier used for anonymous / unregistered exam attempts.
var GuestUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

// StartExam begins (or resumes) a student's exam sitting. It is idempotent: an
// in-progress attempt still within its window is resumed with remaining time. An
// in-progress attempt found past its deadline is finalized (expired→graded) per
// the lazy-expiry rule before a fresh attempt is created.
//
// Unregistered/anonymous guests (actor == nil) are allowed to take published mock exams.
func (s *Service) StartExam(ctx context.Context, actor *Actor, examID uuid.UUID) (*ExamAttemptClientView, error) {
	exam, err := s.exams.FindByID(ctx, examID)
	if err != nil {
		return nil, err
	}
	if err := s.assertExamAccess(ctx, actor, exam); err != nil {
		return nil, err
	}

	if actor != nil {
		if existing, err := s.attempts.GetActive(ctx, actor.UserID, examID); err == nil {
			if s.expired(exam, existing) {
				_, _ = s.finalizeExpired(ctx, exam, existing)
			} else {
				return s.attemptView(ctx, exam, existing)
			}
		} else if !apperror.IsNotFound(err) {
			return nil, err
		}
	}

	studentID := GuestUserID
	if actor != nil {
		studentID = actor.UserID
	}

	now := s.now()
	expiresAt := now.Add(time.Duration(exam.DurationMinutes) * time.Minute)
	attempt := &entity.ExamAttempt{
		ExamID:    examID,
		StudentID: studentID,
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
	studentID := GuestUserID
	if actor != nil {
		studentID = actor.UserID
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, studentID)
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
	studentID := GuestUserID
	if actor != nil {
		studentID = actor.UserID
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, studentID)
	if err != nil {
		return nil, err
	}
	exam, err := s.exams.FindByID(ctx, attempt.ExamID)
	if err != nil {
		return nil, err
	}
	// Already finalized (e.g. the sweeper graded it first): return the existing
	// result rather than erroring, so a late client submit is benign.
	if attempt.Status == entity.AttemptGraded {
		return s.gradeView(ctx, exam, attempt)
	}
	if attempt.Status != entity.AttemptInProgress {
		return nil, invalidState("attempt is not in progress")
	}

	saved, _ := s.state.GetAnswers(ctx, attemptID)
	merged := mergeAnswers(saved, answers)

	target := entity.AttemptSubmitted
	if s.expired(exam, attempt) {
		target = entity.AttemptExpired
	}
	res, _, err := s.finalize(ctx, exam, attempt, merged, target)
	return res, err
}

// GetExamResult returns the graded result for one of the actor's own attempts.
// It is available only in a terminal state; an attempt found past its deadline
// while still in progress is finalized first (lazy expiry).
func (s *Service) GetExamResult(ctx context.Context, actor *Actor, attemptID uuid.UUID) (*ExamResult, error) {
	studentID := GuestUserID
	if actor != nil {
		studentID = actor.UserID
	}
	attempt, err := s.attempts.FindByID(ctx, attemptID, studentID)
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
	studentID := GuestUserID
	if actor != nil {
		studentID = actor.UserID
	}
	return s.ListAttemptsByStudent(ctx, studentID, examID)
}

// ListAttemptsByStudent returns the student's attempt history, optionally filtered by examID.
func (s *Service) ListAttemptsByStudent(ctx context.Context, studentID uuid.UUID, examID uuid.UUID) ([]*entity.ExamAttempt, error) {
	all, err := s.attempts.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if examID == uuid.Nil {
		return all, nil
	}
	filtered := make([]*entity.ExamAttempt, 0, len(all))
	for _, a := range all {
		if a.ExamID == examID {
			filtered = append(filtered, a)
		}
	}
	return filtered, nil
}

// finalizeExpired drives an in-progress, past-deadline attempt through
// expired→graded using whatever answers were auto-saved.
func (s *Service) finalizeExpired(ctx context.Context, exam *entity.Exam, attempt *entity.ExamAttempt) (*ExamResult, error) {
	saved, _ := s.state.GetAnswers(ctx, attempt.ID)
	res, _, err := s.finalize(ctx, exam, attempt, mergeAnswers(saved, nil), entity.AttemptExpired)
	return res, err
}

// finalize performs the in_progress → {submitted|expired} → graded sequence.
// Each transition is validated against the §9.2 state machine, and the move out
// of in_progress goes through the guarded ClaimForGrading so a race between a
// client submit and the background sweeper grades exactly once: the first writer
// to flip the row wins (claimed=true); the loser returns the winner's already-
// graded result with claimed=false. Grading is server-side; Redis state is
// cleared at the end.
func (s *Service) finalize(ctx context.Context, exam *entity.Exam, attempt *entity.ExamAttempt, answers map[string]any, target entity.AttemptStatus) (*ExamResult, bool, error) {
	if err := guard(attempt.Status, target); err != nil {
		return nil, false, err
	}

	gradables, err := s.loadGradables(ctx, exam.ID)
	if err != nil {
		return nil, false, err
	}
	sub := parseSubmitted(gradables, answers)
	core := grading.Grade(gradables, sub)

	now := s.now()
	attempt.Status = target
	attempt.SubmittedAt = &now
	attempt.AnswersSnapshot = snapshotFor(answers, core)

	claimed, err := s.attempts.ClaimForGrading(ctx, attempt)
	if err != nil {
		return nil, false, err
	}
	if !claimed {
		// Lost the race — return the winner's graded result.
		fresh, ferr := s.attempts.FindByID(ctx, attempt.ID, attempt.StudentID)
		if ferr != nil {
			return nil, false, ferr
		}
		res, gerr := s.gradeView(ctx, exam, fresh)
		return res, false, gerr
	}

	if err := guard(attempt.Status, entity.AttemptGraded); err != nil {
		return nil, false, err
	}
	score := core.Score
	pct := core.Percentage
	attempt.Score = &score
	attempt.Percentage = &pct
	if err := s.attempts.Grade(ctx, attempt); err != nil {
		return nil, false, err
	}
	attempt.Status = entity.AttemptGraded

	if s.ranking != nil {
		_ = s.ranking.RecordExamScore(ctx, exam.ID, attempt.StudentID, core.Percentage)
	}
	_ = s.state.DeleteAttemptState(ctx, attempt.ID)

	// Fire the exam-graded notification asynchronously via the hook (must not
	// block or fail this request).
	s.notify.ExamGraded(ctx, attempt.StudentID, attempt.ID, core.Score)

	res := s.result(exam, attempt, core)
	res.SubmittedAt = &now
	return res, true, nil
}

// SweepExpired finalizes abandoned (in-progress, past-deadline) exam attempts so
// the §9.2 server timer fires even without a client request. Called by the
// background worker (Skill 4). Each is graded with its last Redis-saved answers;
// the guarded claim makes a race with a late submit grade exactly once. Returns
// the number of attempts graded this pass.
func (s *Service) SweepExpired(ctx context.Context, now time.Time, limit int) (int, error) {
	expired, err := s.attempts.ListExpired(ctx, now, limit)
	if err != nil {
		return 0, err
	}
	graded := 0
	for _, attempt := range expired {
		exam, err := s.exams.FindByID(ctx, attempt.ExamID)
		if err != nil {
			return graded, err
		}
		saved, _ := s.state.GetAnswers(ctx, attempt.ID)
		_, claimed, err := s.finalize(ctx, exam, attempt, saved, entity.AttemptExpired)
		if err != nil {
			return graded, err
		}
		if claimed {
			graded++
		}
	}
	return graded, nil
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

// assertExamAccess gates an exam sitting. Any published mock exam is open to all
// participants (including unregistered guests and non-enrolled students).
// Unpublished/draft exams remain restricted to course owners and admins.
func (s *Service) assertExamAccess(ctx context.Context, actor *Actor, exam *entity.Exam) error {
	if exam.IsPublished {
		return nil
	}
	if actor != nil && actor.Role == entity.RoleAdmin {
		return nil
	}
	if exam.CourseID != nil && actor != nil {
		return s.access.RequireFullAccess(ctx, access.Actor{UserID: actor.UserID, Role: actor.Role}, *exam.CourseID)
	}
	return apperror.Forbidden("you do not have access to this exam", nil)
}
