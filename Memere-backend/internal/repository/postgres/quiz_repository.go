package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// QuizRepo is the sqlc-backed implementation of repository.QuizRepository.
type QuizRepo struct {
	q *sqlcgen.Queries
}

var _ repository.QuizRepository = (*QuizRepo)(nil)

// NewQuizRepo builds a QuizRepo over a pgx pool.
func NewQuizRepo(pool *pgxpool.Pool) *QuizRepo {
	return &QuizRepo{q: sqlcgen.New(pool)}
}

func (r *QuizRepo) Create(ctx context.Context, q *entity.Quiz) error {
	row, err := queriesFor(ctx, r.q).CreateQuiz(ctx, sqlcgen.CreateQuizParams{
		LessonID:           toPgUUIDPtr(q.LessonID),
		CourseID:           toPgUUID(q.CourseID),
		Title:              q.Title,
		TimeLimitSeconds:   toPgInt4Ptr(q.TimeLimitSeconds),
		PassPercentage:     toPgNumeric(q.PassPercentage),
		RandomizeQuestions: q.RandomizeQuestions,
		MaxAttempts:        toPgInt4Ptr(q.MaxAttempts),
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*q = *quizFromRow(row)
	return nil
}

func (r *QuizRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Quiz, error) {
	row, err := queriesFor(ctx, r.q).GetQuizByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapQuizErr(err)
	}
	return quizFromRow(row), nil
}

func (r *QuizRepo) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Quiz, error) {
	rows, err := queriesFor(ctx, r.q).ListQuizzesByCourse(ctx, toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	quizzes := make([]*entity.Quiz, len(rows))
	for i, row := range rows {
		quizzes[i] = quizFromRow(row)
	}
	return quizzes, nil
}

func (r *QuizRepo) Update(ctx context.Context, q *entity.Quiz) error {
	row, err := queriesFor(ctx, r.q).UpdateQuiz(ctx, sqlcgen.UpdateQuizParams{
		ID:                 toPgUUID(q.ID),
		Title:              q.Title,
		TimeLimitSeconds:   toPgInt4Ptr(q.TimeLimitSeconds),
		PassPercentage:     toPgNumeric(q.PassPercentage),
		RandomizeQuestions: q.RandomizeQuestions,
		MaxAttempts:        toPgInt4Ptr(q.MaxAttempts),
	})
	if err != nil {
		return mapQuizErr(err)
	}
	*q = *quizFromRow(row)
	return nil
}

// Delete soft-deletes the quiz (Non-Negotiable #5).
func (r *QuizRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteQuiz(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// GetQuizWithQuestions loads the quiz and, for each question, its full answer
// set INCLUDING the answer key — for server-side grading only.
func (r *QuizRepo) GetQuizWithQuestions(ctx context.Context, quizID uuid.UUID) (*repository.QuizWithQuestions, error) {
	q := queriesFor(ctx, r.q)

	quizRow, err := q.GetQuizByID(ctx, toPgUUID(quizID))
	if err != nil {
		return nil, mapQuizErr(err)
	}
	questionRows, err := q.ListQuestionsByQuiz(ctx, toPgUUID(quizID))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	questions := make([]repository.QuestionWithAnswers, len(questionRows))
	for i, qr := range questionRows {
		answerRows, err := q.ListAnswersByQuestion(ctx, qr.ID)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		answers := make([]*entity.Answer, len(answerRows))
		for j, ar := range answerRows {
			answers[j] = answerFromRow(ar)
		}
		questions[i] = repository.QuestionWithAnswers{
			Question: questionFromRow(qr),
			Answers:  answers,
		}
	}

	return &repository.QuizWithQuestions{
		Quiz:      quizFromRow(quizRow),
		Questions: questions,
	}, nil
}

// GetQuestionsForClient loads the quiz's questions and options with the answer
// key structurally excluded (separate is_correct-free queries).
func (r *QuizRepo) GetQuestionsForClient(ctx context.Context, quizID uuid.UUID) ([]repository.ClientQuestion, error) {
	q := queriesFor(ctx, r.q)

	questionRows, err := q.GetQuestionsForClient(ctx, toPgUUID(quizID))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	out := make([]repository.ClientQuestion, len(questionRows))
	for i, qr := range questionRows {
		answerRows, err := q.ListAnswersForClient(ctx, qr.ID)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		answers := make([]repository.ClientAnswer, len(answerRows))
		for j, ar := range answerRows {
			answers[j] = repository.ClientAnswer{
				ID:         fromPgUUID(ar.ID),
				QuestionID: fromPgUUID(ar.QuestionID),
				Text:       ar.Text,
				OrderIndex: int(ar.OrderIndex),
			}
		}
		out[i] = repository.ClientQuestion{
			ID:         fromPgUUID(qr.ID),
			QuizID:     fromPgUUID(qr.QuizID),
			Text:       qr.Text,
			Type:       entity.QuestionType(qr.Type),
			Points:     int(qr.Points),
			OrderIndex: int(qr.OrderIndex),
			Subject:    qr.Subject,
			Topic:      qr.Topic,
			Answers:    answers,
		}
	}
	return out, nil
}

func mapQuizErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("quiz not found", err)
	}
	return apperror.Internal(err)
}

func quizFromRow(row sqlcgen.CoursesQuiz) *entity.Quiz {
	return &entity.Quiz{
		ID:                 fromPgUUID(row.ID),
		LessonID:           fromPgUUIDPtr(row.LessonID),
		CourseID:           fromPgUUID(row.CourseID),
		Title:              row.Title,
		TimeLimitSeconds:   fromPgInt4Ptr(row.TimeLimitSeconds),
		PassPercentage:     fromPgNumeric(row.PassPercentage),
		RandomizeQuestions: row.RandomizeQuestions,
		MaxAttempts:        fromPgInt4Ptr(row.MaxAttempts),
		CreatedAt:          fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:          fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:          fromPgTimestamptz(row.DeletedAt),
	}
}

func questionFromRow(row sqlcgen.CoursesQuestion) *entity.Question {
	return &entity.Question{
		ID:          fromPgUUID(row.ID),
		QuizID:      fromPgUUID(row.QuizID),
		Text:        row.Text,
		Type:        entity.QuestionType(row.Type),
		Points:      int(row.Points),
		Explanation: row.Explanation,
		OrderIndex:  int(row.OrderIndex),
		Subject:     row.Subject,
		Topic:       row.Topic,
		CreatedAt:   fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:   fromPgTimestamptzValue(row.UpdatedAt),
	}
}

func answerFromRow(row sqlcgen.CoursesAnswer) *entity.Answer {
	return &entity.Answer{
		ID:         fromPgUUID(row.ID),
		QuestionID: fromPgUUID(row.QuestionID),
		Text:       row.Text,
		IsCorrect:  row.IsCorrect,
		OrderIndex: int(row.OrderIndex),
		CreatedAt:  fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:  fromPgTimestamptzValue(row.UpdatedAt),
	}
}
