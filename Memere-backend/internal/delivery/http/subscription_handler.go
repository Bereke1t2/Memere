package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/payment"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/subscription"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// SubscriptionHandler adapts the subscription usecase to HTTP. Subscribing is
// sugar over the payment flow (a checkout with no course id), so the handler also
// holds the payment service to initiate the purchase. Ownership is enforced in
// the usecase.
type SubscriptionHandler struct {
	svc *subscription.Service
	pay *payment.Service
}

// NewSubscriptionHandler builds a SubscriptionHandler.
func NewSubscriptionHandler(svc *subscription.Service, pay *payment.Service) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, pay: pay}
}

// subscriptionActor adapts the context actor to the subscription usecase's Actor.
func subscriptionActor(c *gin.Context) subscription.Actor {
	a, ok := middleware.ActorFromContext(c)
	if !ok || a == nil {
		return subscription.Actor{}
	}
	return subscription.Actor{UserID: a.UserID, Role: a.Role}
}

// ListPlans handles GET /subscription-plans (public): the configured plan
// catalogue with prices and billing periods.
func (h *SubscriptionHandler) ListPlans(c *gin.Context) {
	respondJSON(c, http.StatusOK, gin.H{"data": dto.NewPlanList(h.svc.ListPlans())})
}

// Subscribe handles POST /subscriptions (bearer): starts a subscription purchase
// by initiating a payment with no course id (CourseID == nil ⇒ subscription).
// The Idempotency-Key header is mandatory, exactly as for a course purchase.
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	idem := c.GetHeader("Idempotency-Key")
	if idem == "" {
		respondError(c, apperror.New(http.StatusBadRequest, "IDEMPOTENCY_KEY_REQUIRED",
			"an Idempotency-Key header is required to start a payment", nil))
		return
	}
	var req dto.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	out, err := h.pay.Initiate(c.Request.Context(), payment.InitiateInput{
		Actor:          paymentActor(c),
		CourseID:       nil,
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

// GetMine handles GET /me/subscription (bearer): the caller's active subscription.
func (h *SubscriptionHandler) GetMine(c *gin.Context) {
	sub, err := h.svc.GetMySubscription(c.Request.Context(), subscriptionActor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewSubscriptionResponse(sub))
}

// Cancel handles POST /subscriptions/:id/cancel (bearer, owner/admin): flags the
// subscription to not renew while keeping access until the period ends.
func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.CancelSubscription(c.Request.Context(), subscriptionActor(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
