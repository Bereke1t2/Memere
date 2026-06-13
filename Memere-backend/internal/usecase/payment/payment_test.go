package payment

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

type harness struct {
	svc      *Service
	payments *fakePaymentRepo
	enroll   *fakeEnrollRepo
	coupons  *fakeCouponRepo
	webhooks *fakeWebhookRepo
	courses  *fakeCourseRepo
	provider *mockProvider
	clock    time.Time

	studentID uuid.UUID
	courseID  uuid.UUID
}

func newHarness() *harness {
	h := &harness{
		payments:  newFakePaymentRepo(),
		enroll:    newFakeEnrollRepo(),
		coupons:   newFakeCouponRepo(),
		webhooks:  newFakeWebhookRepo(),
		courses:   newFakeCourseRepo(),
		provider:  &mockProvider{name: string(entity.ProviderChapa)},
		clock:     time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC),
		studentID: uuid.New(),
		courseID:  uuid.New(),
	}
	h.courses.add(&entity.Course{ID: h.courseID, TeacherID: uuid.New(), Price: 100, Currency: "ETB", IsPublished: true})

	// A passthrough quoter unless a test installs its own; defaults to no discount.
	quoter := quoterFunc(func(_ context.Context, _ string, _ uuid.UUID, base decimal.Decimal) (decimal.Decimal, uuid.UUID, error) {
		return base, uuid.Nil, nil
	})

	h.svc = NewService(
		h.payments, h.enroll, h.coupons, h.webhooks, h.courses,
		mockRegistry{p: h.provider}, quoter, fakeTxManager{}, nil,
		Config{CallbackURL: "https://api.memere/webhook", ReturnURL: "https://app.memere/done", DefaultCurrency: "ETB"},
		func() time.Time { return h.clock },
	)
	return h
}

func (h *harness) actor() Actor { return Actor{UserID: h.studentID, Role: entity.RoleStudent} }

func (h *harness) initiate(key string) (*InitiateResult, error) {
	cid := h.courseID
	return h.svc.Initiate(context.Background(), InitiateInput{
		Actor:          h.actor(),
		CourseID:       &cid,
		Provider:       entity.ProviderChapa,
		IdempotencyKey: key,
	})
}

// completedEvent builds a normalized "completed" webhook for a payment id.
func completedEvent(paymentID, eventID string) *service.WebhookEvent {
	return &service.WebhookEvent{
		ProviderEventID: eventID,
		Type:            "charge.success",
		PaymentRef:      paymentID,
		Status:          string(entity.PayCompleted),
		ProviderTxnID:   "txn-" + eventID,
		RawPayload:      []byte(`{"ok":true}`),
	}
}

// --- initiate -----------------------------------------------------------------

func TestInitiate_RequiresIdempotencyKey(t *testing.T) {
	h := newHarness()
	_, err := h.initiate("")
	if !apperror.IsCode(err, "IDEMPOTENCY_KEY_REQUIRED") {
		t.Fatalf("missing key should be rejected, got %v", err)
	}
}

func TestInitiate_ReplaySameKeyReturnsSamePayment(t *testing.T) {
	h := newHarness()
	first, err := h.initiate("key-1")
	if err != nil {
		t.Fatalf("first initiate: %v", err)
	}
	second, err := h.initiate("key-1")
	if err != nil {
		t.Fatalf("replay initiate: %v", err)
	}
	if first.PaymentID != second.PaymentID {
		t.Fatalf("replay returned a different payment: %s vs %s", first.PaymentID, second.PaymentID)
	}
	if h.payments.creates != 1 {
		t.Fatalf("replay created %d payments, want exactly 1", h.payments.creates)
	}
	if second.RedirectURL == "" {
		t.Fatalf("replay should return the stored redirect URL")
	}
}

func TestInitiate_AlreadyEnrolledConflicts(t *testing.T) {
	h := newHarness()
	// Pre-enrol the student.
	_ = h.enroll.Create(context.Background(), &entity.Enrollment{
		ID: uuid.New(), StudentID: h.studentID, CourseID: h.courseID, Source: entity.SourcePurchase,
	})
	_, err := h.initiate("key-1")
	if !apperror.IsCode(err, "ALREADY_ENROLLED") {
		t.Fatalf("already-enrolled should conflict, got %v", err)
	}
	if h.payments.creates != 0 {
		t.Fatalf("no payment should be created for an enrolled student")
	}
}

func TestInitiate_FreeCourseRejected(t *testing.T) {
	h := newHarness()
	h.courses.byID[h.courseID].IsFree = true
	_, err := h.initiate("key-1")
	if !apperror.IsCode(err, "COURSE_IS_FREE") {
		t.Fatalf("free course should be rejected, got %v", err)
	}
}

func TestInitiate_CouponDiscountsAmount(t *testing.T) {
	h := newHarness()
	couponID := uuid.New()
	h.svc.quoter = quoterFunc(func(_ context.Context, code string, _ uuid.UUID, base decimal.Decimal) (decimal.Decimal, uuid.UUID, error) {
		return base.Mul(decimal.NewFromFloat(0.5)), couponID, nil // 50% off
	})
	code := "HALF"
	cid := h.courseID
	res, err := h.svc.Initiate(context.Background(), InitiateInput{
		Actor: h.actor(), CourseID: &cid, Provider: entity.ProviderChapa,
		CouponCode: &code, IdempotencyKey: "key-1",
	})
	if err != nil {
		t.Fatalf("initiate with coupon: %v", err)
	}
	if res.Amount != "50" {
		t.Fatalf("coupon amount = %s, want 50", res.Amount)
	}
	// The coupon id must be recorded on the payment (burned only at fulfillment).
	pid, _ := uuid.Parse(res.PaymentID)
	p, _ := h.payments.GetByID(context.Background(), pid)
	if p.CouponID == nil || *p.CouponID != couponID {
		t.Fatalf("payment should record the coupon id")
	}
	if h.coupons.increments != 0 {
		t.Fatalf("coupon use must NOT be incremented at initiate")
	}
}

// --- webhook + fulfillment ----------------------------------------------------

func TestWebhook_BadSignatureRejected(t *testing.T) {
	h := newHarness()
	h.provider.verifyErr = apperror.New(401, "BAD_SIGNATURE", "bad sig", nil)
	err := h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`))
	if !apperror.IsCode(err, "BAD_SIGNATURE") {
		t.Fatalf("bad signature should propagate, got %v", err)
	}
}

func TestWebhook_SuccessCompletesAndGrants(t *testing.T) {
	h := newHarness()
	res, err := h.initiate("key-1")
	if err != nil {
		t.Fatalf("initiate: %v", err)
	}
	h.provider.verifyEvent = completedEvent(res.PaymentID, "evt-1")
	if err := h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`)); err != nil {
		t.Fatalf("webhook: %v", err)
	}
	pid, _ := uuid.Parse(res.PaymentID)
	p, _ := h.payments.GetByID(context.Background(), pid)
	if p.Status != entity.PayCompleted {
		t.Fatalf("payment status = %s, want completed", p.Status)
	}
	if ok, _ := h.enroll.Exists(context.Background(), h.studentID, h.courseID); !ok {
		t.Fatalf("enrollment should be granted on success")
	}
}

func TestWebhook_DuplicateEventFulfillsOnce(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	h.provider.verifyEvent = completedEvent(res.PaymentID, "evt-1")

	for i := 0; i < 3; i++ {
		if err := h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`)); err != nil {
			t.Fatalf("webhook %d: %v", i, err)
		}
	}
	if h.enroll.creates != 1 {
		t.Fatalf("re-delivered webhook granted %d enrollments, want exactly 1", h.enroll.creates)
	}
}

func TestWebhook_SuccessWithCouponBurnsOnce(t *testing.T) {
	h := newHarness()
	couponID := uuid.New()
	max := 5
	h.coupons.add(&entity.Coupon{ID: couponID, Code: "C", DiscountType: entity.DiscountPercentage, DiscountValue: decimal.NewFromInt(10), MaxUses: &max})
	h.svc.quoter = quoterFunc(func(_ context.Context, _ string, _ uuid.UUID, base decimal.Decimal) (decimal.Decimal, uuid.UUID, error) {
		return base, couponID, nil
	})
	code := "C"
	cid := h.courseID
	res, err := h.svc.Initiate(context.Background(), InitiateInput{
		Actor: h.actor(), CourseID: &cid, Provider: entity.ProviderChapa, CouponCode: &code, IdempotencyKey: "k",
	})
	if err != nil {
		t.Fatalf("initiate: %v", err)
	}
	h.provider.verifyEvent = completedEvent(res.PaymentID, "evt-1")
	// Deliver twice; the coupon must burn exactly once.
	_ = h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`))
	h.provider.verifyEvent = completedEvent(res.PaymentID, "evt-1") // same event id -> dedup
	_ = h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`))
	if h.coupons.increments != 1 {
		t.Fatalf("coupon burned %d times, want exactly 1", h.coupons.increments)
	}
}

func TestWebhook_FailureMarksFailedNoGrant(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	ev := completedEvent(res.PaymentID, "evt-1")
	ev.Status = string(entity.PayFailed)
	h.provider.verifyEvent = ev
	if err := h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`)); err != nil {
		t.Fatalf("webhook: %v", err)
	}
	pid, _ := uuid.Parse(res.PaymentID)
	p, _ := h.payments.GetByID(context.Background(), pid)
	if p.Status != entity.PayFailed {
		t.Fatalf("status = %s, want failed", p.Status)
	}
	if h.enroll.creates != 0 {
		t.Fatalf("failed payment must not grant enrollment")
	}
}

// Webhook and polling both observe success -> exactly one enrollment.
func TestWebhookAndPollRace_GrantsOnce(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)

	// Webhook settles first.
	h.provider.verifyEvent = completedEvent(res.PaymentID, "evt-1")
	if err := h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`)); err != nil {
		t.Fatalf("webhook: %v", err)
	}
	// Then a poll observes success too; the guarded transition no-ops.
	h.provider.verifyPay = string(entity.PayCompleted)
	if _, err := h.svc.GetPaymentStatus(context.Background(), h.actor(), pid); err != nil {
		t.Fatalf("poll: %v", err)
	}
	if h.enroll.creates != 1 {
		t.Fatalf("webhook+poll granted %d enrollments, want exactly 1", h.enroll.creates)
	}
}

// A lost webhook is recovered by polling: status verify completes the payment.
func TestGetPaymentStatus_PollCompletesPending(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)
	h.provider.verifyPay = string(entity.PayCompleted)

	view, err := h.svc.GetPaymentStatus(context.Background(), h.actor(), pid)
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	if view.Status != string(entity.PayCompleted) {
		t.Fatalf("status = %s, want completed", view.Status)
	}
	if h.enroll.creates != 1 {
		t.Fatalf("poll fulfillment should grant exactly one enrollment")
	}
}

func TestGetPaymentStatus_OtherStudentForbidden(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)
	intruder := Actor{UserID: uuid.New(), Role: entity.RoleStudent}
	_, err := h.svc.GetPaymentStatus(context.Background(), intruder, pid)
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("a student must not read another's payment, got %v", err)
	}
}

// --- refund -------------------------------------------------------------------

func TestRefund_AdminOnly(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)
	_, err := h.svc.RefundPayment(context.Background(), h.actor(), pid)
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("non-admin refund should be forbidden, got %v", err)
	}
}

func TestRefund_CompletedToRefunded(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)
	h.provider.verifyEvent = completedEvent(res.PaymentID, "evt-1")
	_ = h.svc.HandleWebhook(context.Background(), entity.ProviderChapa, nil, []byte(`{}`))

	admin := Actor{UserID: uuid.New(), Role: entity.RoleAdmin}
	view, err := h.svc.RefundPayment(context.Background(), admin, pid)
	if err != nil {
		t.Fatalf("refund: %v", err)
	}
	if view.Status != string(entity.PayRefunded) {
		t.Fatalf("status = %s, want refunded", view.Status)
	}
}

func TestRefund_PendingNotRefundable(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)
	admin := Actor{UserID: uuid.New(), Role: entity.RoleAdmin}
	_, err := h.svc.RefundPayment(context.Background(), admin, pid)
	if !apperror.IsCode(err, "NOT_REFUNDABLE") {
		t.Fatalf("pending payment is not refundable, got %v", err)
	}
}

// guard: ensure metadata round-trips the redirect URL (used by replay).
func TestInitiate_StoresRedirectInMetadata(t *testing.T) {
	h := newHarness()
	res, _ := h.initiate("key-1")
	pid, _ := uuid.Parse(res.PaymentID)
	p, _ := h.payments.GetByID(context.Background(), pid)
	if redirectURLOf(p) != res.RedirectURL {
		t.Fatalf("stored redirect %q != returned %q", redirectURLOf(p), res.RedirectURL)
	}
	// sanity: metadata is JSON-serializable (delivery will marshal it)
	if _, err := json.Marshal(p.Metadata); err != nil {
		t.Fatalf("metadata not serializable: %v", err)
	}
}
