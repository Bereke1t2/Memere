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

// QuestionRepo is the sqlc-backed implementation of repository.QuestionRepository.
type QuestionRepo struct {
	q *sqlcgen.Queries
}

var _ repository.QuestionRepository = (*QuestionRepo)(nil)

// NewQuestionRepo builds a QuestionRepo over a pgx pool.
func NewQuestionRepo(pool *pgxpool.Pool) *QuestionRepo {
	return &QuestionRepo{q: sqlcgen.New(pool)}
}

// Create inserts the question and its answer options through the same queries
// handle, so when run inside a WithinTx the question and its answer key commit
// atomically (the answer key is never left half-written).
func (r *QuestionRepo) Create(ctx context.Context, q *entity.Question, answers []*entity.Answer) error {
	queries := queriesFor(ctx, r.q)

	row, err := queries.CreateQuestion(ctx, sqlcgen.CreateQuestionParams{
		QuizID:      toPgUUID(q.QuizID),
		Text:        q.Text,
		Type:        string(q.Type),
		Points:      int32(q.Points),
		Explanation: q.Explanation,
		OrderIndex:  int32(q.OrderIndex),
		Subject:     q.Subject,
		Topic:       q.Topic,
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*q = *questionFromRow(row)

	for _, a := range answers {
		a.QuestionID = q.ID
		ar, err := queries.CreateAnswer(ctx, sqlcgen.CreateAnswerParams{
			QuestionID: toPgUUID(a.QuestionID),
			Text:       a.Text,
			IsCorrect:  a.IsCorrect,
			OrderIndex: int32(a.OrderIndex),
		})
		if err != nil {
			return apperror.Internal(err)
		}
		*a = *answerFromRow(ar)
	}
	return nil
}

func (r *QuestionRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	row, err := queriesFor(ctx, r.q).GetQuestionByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("question not found", err)
		}
		return nil, apperror.Internal(err)
	}
	return questionFromRow(row), nil
}

func (r *QuestionRepo) ListByQuiz(ctx context.Context, quizID uuid.UUID) ([]*entity.Question, error) {
	rows, err := queriesFor(ctx, r.q).ListQuestionsByQuiz(ctx, toPgUUID(quizID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	questions := make([]*entity.Question, len(rows))
	for i, row := range rows {
		questions[i] = questionFromRow(row)
	}
	return questions, nil
}

// ListAnswers returns a question's options WITH the answer key — server-internal
// (grading / teacher) use only.
func (r *QuestionRepo) ListAnswers(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error) {
	rows, err := queriesFor(ctx, r.q).ListAnswersByQuestion(ctx, toPgUUID(questionID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	answers := make([]*entity.Answer, len(rows))
	for i, row := range rows {
		answers[i] = answerFromRow(row)
	}
	return answers, nil
}

func (r *QuestionRepo) Update(ctx context.Context, q *entity.Question) error {
	row, err := queriesFor(ctx, r.q).UpdateQuestion(ctx, sqlcgen.UpdateQuestionParams{
		ID:          toPgUUID(q.ID),
		Text:        q.Text,
		Type:        string(q.Type),
		Points:      int32(q.Points),
		Explanation: q.Explanation,
		OrderIndex:  int32(q.OrderIndex),
		Subject:     q.Subject,
		Topic:       q.Topic,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.NotFound("question not found", err)
		}
		return apperror.Internal(err)
	}
	*q = *questionFromRow(row)
	return nil
}

// Delete removes the question; its answers cascade via the FK. Questions are not
// user-facing content with a tombstone — they are removed with their quiz tree.
func (r *QuestionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).DeleteQuestion(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
