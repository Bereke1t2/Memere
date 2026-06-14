package dto

import (
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/certificate"
)

// CertificateResponse is the wire projection of entity.Certificate.
type CertificateResponse struct {
	ID        string  `json:"id"`
	StudentID string  `json:"student_id"`
	CourseID  string  `json:"course_id"`
	Serial    string  `json:"serial"`
	FileReady bool    `json:"file_ready"`
	IssuedAt  time.Time `json:"issued_at"`
}

// CertDownloadURLResponse is returned by GET /certificates/:id/download.
type CertDownloadURLResponse struct {
	URL string `json:"url"`
}

// PublicCertResponse is returned by GET /verify/certificates/:serial.
type PublicCertResponse struct {
	StudentName string    `json:"student_name"`
	CourseTitle string    `json:"course_title"`
	IssuedAt    time.Time `json:"issued_at"`
	Valid       bool      `json:"valid"`
}

// NewCertificateResponse maps a domain certificate to the wire DTO.
func NewCertificateResponse(c *entity.Certificate) CertificateResponse {
	return CertificateResponse{
		ID:        c.ID.String(),
		StudentID: c.StudentID.String(),
		CourseID:  c.CourseID.String(),
		Serial:    c.Serial,
		FileReady: c.FileKey != nil && *c.FileKey != "",
		IssuedAt:  c.IssuedAt,
	}
}

// NewPublicCertResponse maps a public cert view to the wire DTO.
func NewPublicCertResponse(v *certificate.PublicCertView) PublicCertResponse {
	return PublicCertResponse{
		StudentName: v.StudentName,
		CourseTitle: v.CourseTitle,
		IssuedAt:    v.IssuedAt,
		Valid:       v.Valid,
	}
}
