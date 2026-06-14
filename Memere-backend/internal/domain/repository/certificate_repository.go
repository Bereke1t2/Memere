package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// CertificateRepository is the persistence port for course completion certificates.
type CertificateRepository interface {
	// CreateIfNotExists inserts a certificate row with the given serial if no
	// certificate already exists for (studentID, courseID). Returns the cert
	// (existing or new) and created=true when a new row was inserted.
	CreateIfNotExists(ctx context.Context, studentID, courseID uuid.UUID, serial string) (cert *entity.Certificate, created bool, err error)

	// GetByID returns a certificate by primary key.
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Certificate, error)

	// GetByStudentCourse returns the certificate for a student+course pair.
	GetByStudentCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Certificate, error)

	// GetBySerial looks up a certificate by its public verification serial.
	GetBySerial(ctx context.Context, serial string) (*entity.Certificate, error)

	// SetFileKey records the S3 key of the rendered PDF.
	SetFileKey(ctx context.Context, id uuid.UUID, fileKey string) error

	// ListByStudent returns all certificates for a student, newest first.
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.Certificate, error)
}
