package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/certificate"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// CertificateHandler adapts the certificate usecase to HTTP.
type CertificateHandler struct {
	svc *certificate.Service
}

func NewCertificateHandler(svc *certificate.Service) *CertificateHandler {
	return &CertificateHandler{svc: svc}
}

// certActor converts the context actor to a certificate.Actor.
func certActor(c *gin.Context) certificate.Actor {
	a, _ := middleware.ActorFromContext(c)
	if a == nil {
		return certificate.Actor{}
	}
	return certificate.Actor{UserID: a.UserID}
}

// Issue handles POST /courses/:id/certificate → 201 with the certificate
func (h *CertificateHandler) Issue(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid course id", err))
		return
	}
	cert, err := h.svc.Issue(c.Request.Context(), certActor(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusCreated, dto.NewCertificateResponse(cert))
}

// ListMine handles GET /me/certificates → 200 with certificate list
func (h *CertificateHandler) ListMine(c *gin.Context) {
	certs, err := h.svc.ListMyCertificates(c.Request.Context(), certActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.CertificateResponse, 0, len(certs))
	for _, cert := range certs {
		items = append(items, dto.NewCertificateResponse(cert))
	}
	respondJSON(c, http.StatusOK, gin.H{"certificates": items})
}

// GetDownloadURL handles GET /certificates/:id/download → 200 with signed URL
func (h *CertificateHandler) GetDownloadURL(c *gin.Context) {
	certID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid certificate id", err))
		return
	}
	url, err := h.svc.GetDownloadURL(c.Request.Context(), certActor(c), certID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.CertDownloadURLResponse{URL: url})
}

// Verify handles GET /verify/certificates/:serial → 200 (public, no auth)
func (h *CertificateHandler) Verify(c *gin.Context) {
	serial := c.Param("serial")
	view, err := h.svc.VerifyCertificate(c.Request.Context(), serial)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewPublicCertResponse(view))
}
