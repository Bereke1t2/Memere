package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/payment"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// maxWebhookBody caps the bytes we read from a provider callback (spec §7.3): the
// raw body is needed for signature verification, but it must be bounded so a
// malicious or buggy provider cannot exhaust memory.
const maxWebhookBody = 1 << 20 // 1 MiB

// PaymentHandler adapts the payment-lifecycle usecase to HTTP. Role gating
// happens in middleware; ownership/access control stays in the usecase (IDOR
// rule), so handlers only bind/validate, resolve the actor, and map results.
type PaymentHandler struct {
	svc *payment.Service
}

// NewPaymentHandler builds a PaymentHandler.
func NewPaymentHandler(svc *payment.Service) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// paymentActor adapts the context actor to the payment usecase's Actor. Returns a
// zero Actor (UserID == Nil) for an anonymous request, which the usecase rejects.
func paymentActor(c *gin.Context) payment.Actor {
	a, ok := middleware.ActorFromContext(c)
	if !ok || a == nil {
		return payment.Actor{}
	}
	return payment.Actor{UserID: a.UserID, Role: a.Role}
}

// Initiate handles POST /payments/initiate (bearer): starts an idempotent
// checkout. The Idempotency-Key header is mandatory (no double-charge) and is
// checked before the body is even parsed.
func (h *PaymentHandler) Initiate(c *gin.Context) {
	idem := c.GetHeader("Idempotency-Key")
	if idem == "" {
		respondError(c, apperror.New(http.StatusBadRequest, "IDEMPOTENCY_KEY_REQUIRED",
			"an Idempotency-Key header is required to start a payment", nil))
		return
	}
	var req dto.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	courseID, err := parseUUIDPtr(req.CourseID)
	if err != nil {
		respondError(c, err)
		return
	}
	out, err := h.svc.Initiate(c.Request.Context(), payment.InitiateInput{
		Actor:          paymentActor(c),
		CourseID:       courseID,
		Plan:           req.Plan,
		Provider:       entity.PaymentProvider(req.Provider),
		CouponCode:     req.CouponCode,
		IdempotencyKey: idem,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusCreated, out)
}

// Status handles GET /payments/:id/status (bearer, owner/admin): the current
// payment state, with a server-to-server verify fallback for a lost webhook.
func (h *PaymentHandler) Status(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.GetPaymentStatus(c.Request.Context(), paymentActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, view)
}

// ListMine handles GET /payments (bearer): the caller's own payment history.
func (h *PaymentHandler) ListMine(c *gin.Context) {
	limit := atoiDefault(c.Query("limit"), 0)
	views, err := h.svc.ListMyPayments(c.Request.Context(), paymentActor(c), limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"data": views})
}

// Refund handles POST /payments/:id/refund (bearer + admin): marks a completed
// payment refunded (admin enforced in middleware and re-checked in the usecase).
func (h *PaymentHandler) Refund(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	view, err := h.svc.RefundPayment(c.Request.Context(), paymentActor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, view)
}

// Webhook handles POST /webhooks/payments/:provider (public, raw body). It reads
// the bounded raw body (needed for signature verification), forwards every header
// to the usecase for provider-specific verification, and returns 200 on any
// successful (idempotent) processing — a re-delivered event is a no-op success.
// Errors map through the usecase's apperror (bad signature → 401, unknown
// payment → 404). It is registered OUTSIDE the auth + JSON middleware so the body
// is untouched and unauthenticated providers can reach it.
func (h *PaymentHandler) Webhook(c *gin.Context) {
	provider := entity.PaymentProvider(c.Param("provider"))
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxWebhookBody))
	if err != nil {
		respondError(c, apperror.BadRequest("could not read webhook body", err))
		return
	}
	headers := make(map[string]string, len(c.Request.Header))
	for k := range c.Request.Header {
		headers[k] = c.GetHeader(k)
	}
	if err := h.svc.HandleWebhook(c.Request.Context(), provider, headers, body); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// parseUUIDPtr parses an optional string id into a *uuid.UUID. A nil or empty
// pointer yields (nil, nil) — a subscription purchase carries no course id. A
// malformed id is a clean 400.
func parseUUIDPtr(s *string) (*uuid.UUID, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, apperror.New(http.StatusBadRequest, "INVALID_ID", "invalid course_id", err)
	}
	return &id, nil
}
