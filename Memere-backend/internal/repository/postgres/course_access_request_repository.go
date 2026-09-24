package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// CourseAccessRequestRepo implements repository.CourseAccessRequestRepository.
type CourseAccessRequestRepo struct {
	pool *pgxpool.Pool
}

var _ repository.CourseAccessRequestRepository = (*CourseAccessRequestRepo)(nil)

// NewCourseAccessRequestRepo builds a CourseAccessRequestRepo over a pgx pool.
func NewCourseAccessRequestRepo(pool *pgxpool.Pool) *CourseAccessRequestRepo {
	return &CourseAccessRequestRepo{pool: pool}
}

// Create inserts or upserts a course access request.
func (r *CourseAccessRequestRepo) Create(ctx context.Context, req *entity.CourseAccessRequest) error {
	id := req.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	query := `
INSERT INTO courses.course_access_requests (
    id, student_id, course_id, status, reviewed_by, reviewed_at, rejection_reason
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (student_id, course_id) DO UPDATE
SET status = EXCLUDED.status,
    reviewed_by = EXCLUDED.reviewed_by,
    reviewed_at = EXCLUDED.reviewed_at,
    rejection_reason = EXCLUDED.rejection_reason,
    updated_at = now()
RETURNING id, student_id, course_id, status, reviewed_by, reviewed_at, rejection_reason, created_at, updated_at;`

	var (
		statusStr string
	)
	err := r.pool.QueryRow(ctx, query,
		id,
		req.StudentID,
		req.CourseID,
		string(req.Status),
		req.ReviewedBy,
		req.ReviewedAt,
		req.RejectionReason,
	).Scan(
		&req.ID,
		&req.StudentID,
		&req.CourseID,
		&statusStr,
		&req.ReviewedBy,
		&req.ReviewedAt,
		&req.RejectionReason,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err != nil {
		return apperror.Internal(err)
	}
	req.Status = entity.RequestStatus(statusStr)
	return nil
}

// GetByID returns the request by ID.
func (r *CourseAccessRequestRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.CourseAccessRequest, error) {
	query := `
SELECT r.id, r.student_id, r.course_id, r.status, r.reviewed_by, r.reviewed_at, r.rejection_reason,
       r.created_at, r.updated_at,
       u.first_name || ' ' || u.last_name AS student_name, u.email AS student_email,
       c.title AS course_title
FROM courses.course_access_requests r
JOIN auth.users u ON u.id = r.student_id
JOIN courses.courses c ON c.id = r.course_id
WHERE r.id = $1;`

	return r.scanOne(r.pool.QueryRow(ctx, query, id))
}

// GetByStudentAndCourse returns the request for student + course.
func (r *CourseAccessRequestRepo) GetByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.CourseAccessRequest, error) {
	query := `
SELECT r.id, r.student_id, r.course_id, r.status, r.reviewed_by, r.reviewed_at, r.rejection_reason,
       r.created_at, r.updated_at,
       u.first_name || ' ' || u.last_name AS student_name, u.email AS student_email,
       c.title AS course_title
FROM courses.course_access_requests r
JOIN auth.users u ON u.id = r.student_id
JOIN courses.courses c ON c.id = r.course_id
WHERE r.student_id = $1 AND r.course_id = $2;`

	return r.scanOne(r.pool.QueryRow(ctx, query, studentID, courseID))
}

// Update persists status changes.
func (r *CourseAccessRequestRepo) Update(ctx context.Context, req *entity.CourseAccessRequest) error {
	query := `
UPDATE courses.course_access_requests
SET status = $2,
    reviewed_by = $3,
    reviewed_at = $4,
    rejection_reason = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, student_id, course_id, status, reviewed_by, reviewed_at, rejection_reason, created_at, updated_at;`

	var statusStr string
	err := r.pool.QueryRow(ctx, query,
		req.ID,
		string(req.Status),
		req.ReviewedBy,
		req.ReviewedAt,
		req.RejectionReason,
	).Scan(
		&req.ID,
		&req.StudentID,
		&req.CourseID,
		&statusStr,
		&req.ReviewedBy,
		&req.ReviewedAt,
		&req.RejectionReason,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.NotFound("course access request not found", err)
		}
		return apperror.Internal(err)
	}
	req.Status = entity.RequestStatus(statusStr)
	return nil
}

// ListByTeacher returns requests for courses owned by the teacher.
func (r *CourseAccessRequestRepo) ListByTeacher(ctx context.Context, teacherID uuid.UUID, status *string, cursor *pagination.Cursor, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	limit = pagination.NormalizeLimit(limit)

	var (
		where []string
		args  []any
	)

	args = append(args, teacherID)
	where = append(where, "c.teacher_id = $"+itoa(len(args)))

	if status != nil && *status != "" {
		args = append(args, *status)
		where = append(where, "r.status = $"+itoa(len(args)))
	}

	if cursor != nil {
		args = append(args, cursor.CreatedAt, cursor.ID)
		where = append(where, "(r.created_at, r.id) < ($"+itoa(len(args)-1)+", $"+itoa(len(args))+")")
	}

	args = append(args, int32(limit+1))
	query := `
SELECT r.id, r.student_id, r.course_id, r.status, r.reviewed_by, r.reviewed_at, r.rejection_reason,
       r.created_at, r.updated_at,
       u.first_name || ' ' || u.last_name AS student_name, u.email AS student_email,
       c.title AS course_title
FROM courses.course_access_requests r
JOIN auth.users u ON u.id = r.student_id
JOIN courses.courses c ON c.id = r.course_id
WHERE ` + strings.Join(where, " AND ") + `
ORDER BY r.created_at DESC, r.id DESC
LIMIT $` + itoa(len(args))

	return r.queryList(ctx, query, args, limit)
}

// ListAll returns all requests across courses (admin).
func (r *CourseAccessRequestRepo) ListAll(ctx context.Context, status *string, cursor *pagination.Cursor, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	limit = pagination.NormalizeLimit(limit)

	var (
		where []string
		args  []any
	)

	if status != nil && *status != "" {
		args = append(args, *status)
		where = append(where, "r.status = $"+itoa(len(args)))
	}

	if cursor != nil {
		args = append(args, cursor.CreatedAt, cursor.ID)
		where = append(where, "(r.created_at, r.id) < ($"+itoa(len(args)-1)+", $"+itoa(len(args))+")")
	}

	args = append(args, int32(limit+1))
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	query := `
SELECT r.id, r.student_id, r.course_id, r.status, r.reviewed_by, r.reviewed_at, r.rejection_reason,
       r.created_at, r.updated_at,
       u.first_name || ' ' || u.last_name AS student_name, u.email AS student_email,
       c.title AS course_title
FROM courses.course_access_requests r
JOIN auth.users u ON u.id = r.student_id
JOIN courses.courses c ON c.id = r.course_id
` + whereClause + `
ORDER BY r.created_at DESC, r.id DESC
LIMIT $` + itoa(len(args))

	return r.queryList(ctx, query, args, limit)
}

// ListByStudent returns requests made by a specific student.
func (r *CourseAccessRequestRepo) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.CourseAccessRequest, error) {
	query := `
SELECT r.id, r.student_id, r.course_id, r.status, r.reviewed_by, r.reviewed_at, r.rejection_reason,
       r.created_at, r.updated_at,
       u.first_name || ' ' || u.last_name AS student_name, u.email AS student_email,
       c.title AS course_title
FROM courses.course_access_requests r
JOIN auth.users u ON u.id = r.student_id
JOIN courses.courses c ON c.id = r.course_id
WHERE r.student_id = $1
ORDER BY r.created_at DESC;`

	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	defer rows.Close()

	var items []*entity.CourseAccessRequest
	for rows.Next() {
		var req entity.CourseAccessRequest
		var statusStr string
		if err := rows.Scan(
			&req.ID,
			&req.StudentID,
			&req.CourseID,
			&statusStr,
			&req.ReviewedBy,
			&req.ReviewedAt,
			&req.RejectionReason,
			&req.CreatedAt,
			&req.UpdatedAt,
			&req.StudentName,
			&req.StudentEmail,
			&req.CourseTitle,
		); err != nil {
			return nil, apperror.Internal(err)
		}
		req.Status = entity.RequestStatus(statusStr)
		items = append(items, &req)
	}
	return items, rows.Err()
}

func (r *CourseAccessRequestRepo) scanOne(row pgx.Row) (*entity.CourseAccessRequest, error) {
	var req entity.CourseAccessRequest
	var statusStr string
	err := row.Scan(
		&req.ID,
		&req.StudentID,
		&req.CourseID,
		&statusStr,
		&req.ReviewedBy,
		&req.ReviewedAt,
		&req.RejectionReason,
		&req.CreatedAt,
		&req.UpdatedAt,
		&req.StudentName,
		&req.StudentEmail,
		&req.CourseTitle,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("course access request not found", err)
		}
		return nil, apperror.Internal(err)
	}
	req.Status = entity.RequestStatus(statusStr)
	return &req, nil
}

func (r *CourseAccessRequestRepo) queryList(ctx context.Context, query string, args []any, limit int) ([]*entity.CourseAccessRequest, *pagination.Cursor, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}
	defer rows.Close()

	var items []*entity.CourseAccessRequest
	for rows.Next() {
		var req entity.CourseAccessRequest
		var statusStr string
		if err := rows.Scan(
			&req.ID,
			&req.StudentID,
			&req.CourseID,
			&statusStr,
			&req.ReviewedBy,
			&req.ReviewedAt,
			&req.RejectionReason,
			&req.CreatedAt,
			&req.UpdatedAt,
			&req.StudentName,
			&req.StudentEmail,
			&req.CourseTitle,
		); err != nil {
			return nil, nil, apperror.Internal(err)
		}
		req.Status = entity.RequestStatus(statusStr)
		items = append(items, &req)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, apperror.Internal(err)
	}

	var next *pagination.Cursor
	if len(items) > limit {
		last := items[limit-1]
		next = &pagination.Cursor{CreatedAt: last.CreatedAt, ID: last.ID}
		items = items[:limit]
	}
	return items, next, nil
}
