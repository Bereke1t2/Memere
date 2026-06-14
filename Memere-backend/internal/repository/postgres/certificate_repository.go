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

// CertificateRepo is the sqlc-backed implementation of repository.CertificateRepository.
type CertificateRepo struct {
	q    *sqlcgen.Queries
	pool *pgxpool.Pool
}

var _ repository.CertificateRepository = (*CertificateRepo)(nil)

func NewCertificateRepo(pool *pgxpool.Pool) *CertificateRepo {
	return &CertificateRepo{q: sqlcgen.New(pool), pool: pool}
}

// CreateIfNotExists inserts a certificate for (studentID, courseID) with serial.
// If a row already exists (ON CONFLICT DO NOTHING), it fetches and returns that row
// with created=false. Returns the row and created=true when a new row was inserted.
func (r *CertificateRepo) CreateIfNotExists(ctx context.Context, studentID, courseID uuid.UUID, serial string) (*entity.Certificate, bool, error) {
	row, err := r.q.CreateCertificateIfNotExists(ctx, sqlcgen.CreateCertificateIfNotExistsParams{
		StudentID: toPgUUID(studentID),
		CourseID:  toPgUUID(courseID),
		Serial:    serial,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// ON CONFLICT DO NOTHING produced no row — fetch the existing cert.
			existing, ferr := r.GetByStudentCourse(ctx, studentID, courseID)
			if ferr != nil {
				return nil, false, ferr
			}
			return existing, false, nil
		}
		return nil, false, apperror.Internal(err)
	}
	return certRowToEntity(row), true, nil
}

func (r *CertificateRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Certificate, error) {
	row, err := r.q.GetCertificateByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("certificate not found", nil)
		}
		return nil, apperror.Internal(err)
	}
	return certRowToEntity(row), nil
}

func (r *CertificateRepo) GetByStudentCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Certificate, error) {
	row, err := r.q.GetCertificateByStudentCourse(ctx, toPgUUID(studentID), toPgUUID(courseID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("certificate not found", nil)
		}
		return nil, apperror.Internal(err)
	}
	return certRowToEntity(row), nil
}

func (r *CertificateRepo) GetBySerial(ctx context.Context, serial string) (*entity.Certificate, error) {
	row, err := r.q.GetCertificateBySerial(ctx, serial)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("certificate not found", nil)
		}
		return nil, apperror.Internal(err)
	}
	return certRowToEntity(row), nil
}

func (r *CertificateRepo) SetFileKey(ctx context.Context, id uuid.UUID, fileKey string) error {
	if err := r.q.SetCertificateFileKey(ctx, toPgUUID(id), fileKey); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (r *CertificateRepo) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.Certificate, error) {
	rows, err := r.q.ListCertificatesByStudent(ctx, toPgUUID(studentID))
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]*entity.Certificate, 0, len(rows))
	for _, row := range rows {
		out = append(out, certRowToEntity(row))
	}
	return out, nil
}

func certRowToEntity(r sqlcgen.CoursesCertificate) *entity.Certificate {
	return &entity.Certificate{
		ID:        fromPgUUID(r.ID),
		StudentID: fromPgUUID(r.StudentID),
		CourseID:  fromPgUUID(r.CourseID),
		Serial:    r.Serial,
		FileKey:   r.FileKey,
		IssuedAt:  fromPgTimestamptzValue(r.IssuedAt),
	}
}
