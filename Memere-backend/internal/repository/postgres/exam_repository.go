package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ExamRepo is the sqlc-backed implementation of repository.ExamRepository.
type ExamRepo struct {
	q    *sqlcgen.Queries
	pool *pgxpool.Pool
}

var _ repository.ExamRepository = (*ExamRepo)(nil)

// NewExamRepo builds an ExamRepo over a pgx pool.
func NewExamRepo(pool *pgxpool.Pool) *ExamRepo {
	return &ExamRepo{q: sqlcgen.New(pool), pool: pool}
}

func (r *ExamRepo) Create(ctx context.Context, e *entity.Exam) error {
	row, err := queriesFor(ctx, r.q).CreateExam(ctx, sqlcgen.CreateExamParams{
		CourseID:        toPgUUIDPtr(e.CourseID),
		Title:           e.Title,
		Subject:         e.Subject,
		Grade:           int32(e.Grade),
		DurationMinutes: int32(e.DurationMinutes),
		TotalMarks:      int32(e.TotalMarks),
		PassMarks:       int32(e.PassMarks),
		Instructions:    e.Instructions,
		IsPublished:     e.IsPublished,
	})
	if err != nil {
		return apperror.Internal(err)
	}
	*e = *examFromRow(row)
	return nil
}

func (r *ExamRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Exam, error) {
	row, err := queriesFor(ctx, r.q).GetExamByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapExamErr(err)
	}
	return examFromRow(row), nil
}

func (r *ExamRepo) List(ctx context.Context, filter repository.ExamFilter, cursor *pagination.Cursor, limit int) ([]*entity.Exam, *pagination.Cursor, error) {
	limit = pagination.NormalizeLimit(limit)

	params := sqlcgen.ListExamsParams{
		Subject:     filter.Subject,
		IsPublished: filter.IsPublished,
		RowLimit:    int32(limit + 1),
	}
	if filter.Grade != nil {
		g := int32(*filter.Grade)
		params.Grade = &g
	}
	if cursor != nil {
		params.AfterCreatedAt = pgTimestamptzValue(cursor.CreatedAt)
		params.AfterID = toPgUUID(cursor.ID)
	}

	rows, err := queriesFor(ctx, r.q).ListExams(ctx, params)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}

	var next *pagination.Cursor
	if len(rows) > limit {
		last := rows[limit-1]
		next = &pagination.Cursor{
			CreatedAt: fromPgTimestamptzValue(last.CreatedAt),
			ID:        fromPgUUID(last.ID),
		}
		rows = rows[:limit]
	}

	exams := make([]*entity.Exam, len(rows))
	for i, row := range rows {
		exams[i] = examFromRow(row)
	}
	return exams, next, nil
}

func (r *ExamRepo) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Exam, error) {
	const q = `SELECT id, course_id, title, subject, grade, duration_minutes, total_marks,
		pass_marks, instructions, is_published, created_at, updated_at, deleted_at
		FROM courses.exams
		WHERE course_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC, id DESC`

	conn := r.pool
	rows, err := conn.Query(ctx, q, toPgUUID(courseID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	defer rows.Close()

	var exams []*entity.Exam
	for rows.Next() {
		var row sqlcgen.CoursesExam
		if err := rows.Scan(
			&row.ID, &row.CourseID, &row.Title, &row.Subject, &row.Grade,
			&row.DurationMinutes, &row.TotalMarks, &row.PassMarks, &row.Instructions,
			&row.IsPublished, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt,
		); err != nil {
			return nil, apperror.Internal(err)
		}
		exams = append(exams, examFromRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.Internal(err)
	}
	return exams, nil
}

func (r *ExamRepo) Update(ctx context.Context, e *entity.Exam) error {
	row, err := queriesFor(ctx, r.q).UpdateExam(ctx, sqlcgen.UpdateExamParams{
		ID:              toPgUUID(e.ID),
		Title:           e.Title,
		Subject:         e.Subject,
		Grade:           int32(e.Grade),
		DurationMinutes: int32(e.DurationMinutes),
		TotalMarks:      int32(e.TotalMarks),
		PassMarks:       int32(e.PassMarks),
		Instructions:    e.Instructions,
		IsPublished:     e.IsPublished,
	})
	if err != nil {
		return mapExamErr(err)
	}
	*e = *examFromRow(row)
	return nil
}

// Delete soft-deletes the exam (Non-Negotiable #5).
func (r *ExamRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteExam(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *ExamRepo) AddQuestion(ctx context.Context, eq *entity.ExamQuestion) error {
	row, err := queriesFor(ctx, r.q).AddExamQuestion(ctx, sqlcgen.AddExamQuestionParams{
		ExamID:     toPgUUID(eq.ExamID),
		QuestionID: toPgUUID(eq.QuestionID),
		OrderIndex: int32(eq.OrderIndex),
		Marks:      int32(eq.Marks),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("QUESTION_ALREADY_IN_EXAM", err)
		}
		return apperror.Internal(err)
	}
	*eq = *examQuestionFromRow(row)
	return nil
}

func (r *ExamRepo) ListQuestions(ctx context.Context, examID uuid.UUID) ([]*entity.ExamQuestion, error) {
	rows, err := queriesFor(ctx, r.q).ListExamQuestions(ctx, toPgUUID(examID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.ExamQuestion, len(rows))
	for i, row := range rows {
		out[i] = examQuestionFromRow(row)
	}
	return out, nil
}

// GetExamWithQuestions loads the exam and, per linked question, its full answer
// set INCLUDING the answer key — for server-side grading only.
func (r *ExamRepo) GetExamWithQuestions(ctx context.Context, examID uuid.UUID) (*repository.ExamWithQuestions, error) {
	q := queriesFor(ctx, r.q)

	examRow, err := q.GetExamByID(ctx, toPgUUID(examID))
	if err != nil {
		return nil, mapExamErr(err)
	}
	linkRows, err := q.ListExamQuestions(ctx, toPgUUID(examID))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	questions := make([]repository.QuestionWithAnswers, len(linkRows))
	for i, link := range linkRows {
		qRow, err := q.GetQuestionByID(ctx, link.QuestionID)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		answerRows, err := q.ListAnswersByQuestion(ctx, link.QuestionID)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		answers := make([]*entity.Answer, len(answerRows))
		for j, ar := range answerRows {
			answers[j] = answerFromRow(ar)
		}
		questions[i] = repository.QuestionWithAnswers{
			Question: questionFromRow(qRow),
			Answers:  answers,
		}
	}

	return &repository.ExamWithQuestions{
		Exam:      examFromRow(examRow),
		Questions: questions,
	}, nil
}

// GetQuestionsForClient loads the exam's questions and options with the answer
// key structurally excluded.
func (r *ExamRepo) GetQuestionsForClient(ctx context.Context, examID uuid.UUID) ([]repository.ClientExamQuestion, error) {
	q := queriesFor(ctx, r.q)

	linkRows, err := q.ListExamQuestionsForClient(ctx, toPgUUID(examID))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	out := make([]repository.ClientExamQuestion, len(linkRows))
	for i, link := range linkRows {
		answerRows, err := q.ListAnswersForClient(ctx, link.QuestionID)
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
		out[i] = repository.ClientExamQuestion{
			QuestionID: fromPgUUID(link.QuestionID),
			OrderIndex: int(link.OrderIndex),
			Marks:      int(link.Marks),
			Text:       link.Text,
			Type:       entity.QuestionType(link.Type),
			Subject:    link.Subject,
			Topic:      link.Topic,
			Answers:    answers,
		}
	}
	return out, nil
}

func mapExamErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("exam not found", err)
	}
	return apperror.Internal(err)
}

func examFromRow(row sqlcgen.CoursesExam) *entity.Exam {
	return &entity.Exam{
		ID:              fromPgUUID(row.ID),
		CourseID:        fromPgUUIDPtr(row.CourseID),
		Title:           row.Title,
		Subject:         row.Subject,
		Grade:           int(row.Grade),
		DurationMinutes: int(row.DurationMinutes),
		TotalMarks:      int(row.TotalMarks),
		PassMarks:       int(row.PassMarks),
		Instructions:    row.Instructions,
		IsPublished:     row.IsPublished,
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}

func examQuestionFromRow(row sqlcgen.CoursesExamQuestion) *entity.ExamQuestion {
	return &entity.ExamQuestion{
		ID:         fromPgUUID(row.ID),
		ExamID:     fromPgUUID(row.ExamID),
		QuestionID: fromPgUUID(row.QuestionID),
		OrderIndex: int(row.OrderIndex),
		Marks:      int(row.Marks),
	}
}
