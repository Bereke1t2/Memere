// Package certificate implements course completion certificate issuance,
// short-lived download URLs, and public serial verification.
package certificate

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// Renderer abstracts PDF generation so the implementation is swappable.
// The infrastructure/pdf package provides the production implementation.
type Renderer interface {
	Certificate(ctx context.Context, data CertData) ([]byte, error)
}

// CertData carries the display values rendered into the certificate PDF.
type CertData struct {
	StudentName string
	CourseTitle string
	Serial      string
	IssuedDate  string
}

// Notifier is the narrow notification port for certificate events.
type Notifier interface {
	CertificateReady(ctx context.Context, studentID, courseID uuid.UUID)
}

// noopNotifier is the default when no notifier is wired (tests that only care
// about certificate logic, not notification side-effects).
type noopNotifier struct{}

func (noopNotifier) CertificateReady(context.Context, uuid.UUID, uuid.UUID) {}

// PublicCertView is the minimal, privacy-respecting response for serial
// verification. Only student name and course title are exposed — no email,
// user ID, or other PII (privacy choice documented here per spec §11.2).
type PublicCertView struct {
	StudentName string
	CourseTitle string
	IssuedAt    time.Time
	Valid       bool
}

// Actor is the authenticated caller's identity.
type Actor struct {
	UserID uuid.UUID
}

// Config carries runtime settings for the certificate service.
type Config struct {
	DownloadURLTTL time.Duration // how long download signed URLs stay valid (default 1h)
}

// Service orchestrates certificate issuance, download, and public verification.
type Service struct {
	repo     repository.CertificateRepository
	progress repository.ProgressRepository
	courses  repository.CourseRepository
	users    repository.UserRepository
	store    service.ObjectStore
	signer   service.URLSigner
	render   Renderer
	notify   Notifier
	cfg      Config
}

// NewService constructs a certificate Service. notify may be nil (defaults to
// noopNotifier). DownloadURLTTL defaults to 1 hour when zero.
func NewService(
	repo repository.CertificateRepository,
	progress repository.ProgressRepository,
	courses repository.CourseRepository,
	users repository.UserRepository,
	store service.ObjectStore,
	signer service.URLSigner,
	render Renderer,
	notify Notifier,
	cfg Config,
) *Service {
	if notify == nil {
		notify = noopNotifier{}
	}
	if cfg.DownloadURLTTL == 0 {
		cfg.DownloadURLTTL = time.Hour
	}
	return &Service{
		repo:     repo,
		progress: progress,
		courses:  courses,
		users:    users,
		store:    store,
		signer:   signer,
		render:   render,
		notify:   notify,
		cfg:      cfg,
	}
}

// Issue grants a certificate when the student has completed 100% of the course.
// Idempotent: returns the existing certificate when called again without
// re-rendering the PDF. Fires CertificateReady asynchronously on first issue.
func (s *Service) Issue(ctx context.Context, actor Actor, courseID uuid.UUID) (*entity.Certificate, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	cp, err := s.progress.GetCourseProgress(ctx, actor.UserID, courseID)
	if err != nil {
		return nil, err
	}
	if cp.PercentComplete.LessThan(decimal.NewFromInt(100)) {
		return nil, apperror.Conflict("COURSE_NOT_COMPLETE", nil)
	}
	serial := s.newSerial()
	cert, created, err := s.repo.CreateIfNotExists(ctx, actor.UserID, courseID, serial)
	if err != nil {
		return nil, err
	}
	if created {
		// Render the PDF in a background goroutine so Issue returns immediately.
		go s.renderAndStore(context.Background(), cert)
	}
	return cert, nil
}

// GetDownloadURL returns a short-lived signed URL for the certificate PDF.
// Only the certificate's student may download it (ownership check prevents IDOR).
func (s *Service) GetDownloadURL(ctx context.Context, actor Actor, certID uuid.UUID) (string, error) {
	if actor.UserID == uuid.Nil {
		return "", apperror.Unauthorized("authentication required", nil)
	}
	cert, err := s.repo.GetByID(ctx, certID)
	if err != nil {
		return "", err
	}
	if cert.StudentID != actor.UserID {
		return "", apperror.Forbidden("certificate belongs to another student", nil)
	}
	if cert.FileKey == nil || *cert.FileKey == "" {
		// PDF render hasn't completed yet (or failed); try synchronously.
		if syncErr := s.renderAndStoreSync(ctx, cert); syncErr != nil {
			return "", syncErr
		}
		// Re-fetch to get the updated file_key.
		cert, err = s.repo.GetByID(ctx, certID)
		if err != nil {
			return "", err
		}
	}
	if cert.FileKey == nil || *cert.FileKey == "" {
		return "", apperror.Internal(fmt.Errorf("certificate PDF not available"))
	}
	return s.signer.SignGet(ctx, *cert.FileKey, s.cfg.DownloadURLTTL)
}

// VerifyCertificate returns minimal public data about a certificate given its
// serial code. This is a public endpoint — no authentication required. Only
// student name and course title are returned; no email, UUID, or other PII.
func (s *Service) VerifyCertificate(ctx context.Context, serial string) (*PublicCertView, error) {
	cert, err := s.repo.GetBySerial(ctx, serial)
	if err != nil {
		return nil, err
	}
	student, err := s.users.FindByID(ctx, cert.StudentID)
	if err != nil {
		return nil, err
	}
	course, err := s.courses.FindByID(ctx, cert.CourseID)
	if err != nil {
		return nil, err
	}
	return &PublicCertView{
		StudentName: student.FirstName + " " + student.LastName,
		CourseTitle: course.Title,
		IssuedAt:    cert.IssuedAt,
		Valid:        true,
	}, nil
}

// ListMyCertificates returns all certificates for the authenticated student.
func (s *Service) ListMyCertificates(ctx context.Context, actor Actor) ([]*entity.Certificate, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("authentication required", nil)
	}
	return s.repo.ListByStudent(ctx, actor.UserID)
}

// renderAndStore renders and persists the PDF in a background goroutine.
func (s *Service) renderAndStore(ctx context.Context, cert *entity.Certificate) {
	_ = s.renderAndStoreSync(ctx, cert)
}

// renderAndStoreSync renders the PDF, uploads to object store, records the
// file_key, and fires the CertificateReady notification.
func (s *Service) renderAndStoreSync(ctx context.Context, cert *entity.Certificate) error {
	student, err := s.users.FindByID(ctx, cert.StudentID)
	if err != nil {
		return err
	}
	course, err := s.courses.FindByID(ctx, cert.CourseID)
	if err != nil {
		return err
	}
	data := CertData{
		StudentName: student.FirstName + " " + student.LastName,
		CourseTitle: course.Title,
		Serial:      cert.Serial,
		IssuedDate:  cert.IssuedAt.Format("January 2, 2006"),
	}
	pdfBytes, err := s.render.Certificate(ctx, data)
	if err != nil {
		return fmt.Errorf("certificate: render PDF: %w", err)
	}
	key := fmt.Sprintf("certificates/%s.pdf", cert.ID)
	if err := s.store.Put(ctx, key, "application/pdf", bytes.NewReader(pdfBytes)); err != nil {
		return fmt.Errorf("certificate: upload PDF: %w", err)
	}
	if err := s.repo.SetFileKey(ctx, cert.ID, key); err != nil {
		return err
	}
	s.notify.CertificateReady(ctx, cert.StudentID, cert.CourseID)
	return nil
}

// newSerial generates a non-guessable public verification code.
// Format: MEM-{year}-{base32(8 random bytes)} — e.g. MEM-2026-ABCDEFGHIJKLMNO
func (s *Service) newSerial() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	return fmt.Sprintf("MEM-%d-%s", time.Now().Year(), code)
}
