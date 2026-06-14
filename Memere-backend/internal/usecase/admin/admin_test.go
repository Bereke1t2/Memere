package admin

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
)

// adminActor returns an authenticated admin Actor.
func adminActor() Actor {
	return Actor{UserID: uuid.New(), Role: entity.RoleAdmin}
}

// studentActor returns a non-admin Actor (student role).
func studentActor() Actor {
	return Actor{UserID: uuid.New(), Role: entity.RoleStudent}
}

// ---- Role gate tests --------------------------------------------------------

func TestSuspendUser_NonAdminForbidden(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)
	err := svc.SuspendUser(ctx, studentActor(), uuid.New(), "test")
	if err == nil {
		t.Error("expected Forbidden for non-admin actor")
	}
}

func TestListUsers_NonAdminForbidden(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)
	_, _, err := svc.ListUsers(ctx, studentActor(), repository.AdminUserFilter{}, nil, 10)
	if err == nil {
		t.Error("expected Forbidden for non-admin actor")
	}
}

func TestUnpublishCourse_NonAdminForbidden(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)
	err := svc.UnpublishCourse(ctx, studentActor(), uuid.New(), "spam")
	if err == nil {
		t.Error("expected Forbidden for non-admin actor")
	}
}

func TestGetEngagementStats_StudentForbidden(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)
	_, err := svc.GetEngagementStats(ctx, studentActor())
	if err == nil {
		t.Error("expected Forbidden for student actor")
	}
}

// ---- SuspendUser tests ------------------------------------------------------

func TestSuspendUser_SetsInactive(t *testing.T) {
	ctx := context.Background()
	users := newFakeUserRepo()
	audit := &fakeAuditRepo{}
	userID := uuid.New()
	users.add(&entity.User{ID: userID, IsActive: true, Role: entity.RoleStudent})

	svc := newSvc(users, newFakeCourseRepo(), newFakePaymentRepo(), audit, nil)
	if err := svc.SuspendUser(ctx, adminActor(), userID, "tos violation"); err != nil {
		t.Fatalf("SuspendUser: %v", err)
	}

	u, _ := users.FindByID(ctx, userID)
	if u.IsActive {
		t.Error("expected IsActive=false after suspension")
	}
	if audit.countAction("user.suspend") != 1 {
		t.Error("expected 1 audit row for user.suspend")
	}
}

func TestSuspendUser_IdempotentAlreadySuspended(t *testing.T) {
	ctx := context.Background()
	users := newFakeUserRepo()
	audit := &fakeAuditRepo{}
	userID := uuid.New()
	users.add(&entity.User{ID: userID, IsActive: false, Role: entity.RoleStudent})

	svc := newSvc(users, newFakeCourseRepo(), newFakePaymentRepo(), audit, nil)
	if err := svc.SuspendUser(ctx, adminActor(), userID, "again"); err != nil {
		t.Fatalf("SuspendUser idempotent: %v", err)
	}
	// No audit row — already suspended.
	if audit.countAction("user.suspend") != 0 {
		t.Error("expected no audit row when already suspended")
	}
}

// ---- ChangeRole tests -------------------------------------------------------

func TestChangeRole_LastAdminDemotionBlocked(t *testing.T) {
	ctx := context.Background()
	users := newFakeUserRepo()
	adminID := uuid.New()
	users.add(&entity.User{ID: adminID, IsActive: true, Role: entity.RoleAdmin})

	svc := newSvc(users, newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)
	err := svc.ChangeRole(ctx, adminActor(), adminID, entity.RoleStudent)
	if err == nil {
		t.Error("expected error: cannot demote the last admin")
	}
}

func TestChangeRole_AllowedWhenMultipleAdmins(t *testing.T) {
	ctx := context.Background()
	users := newFakeUserRepo()
	admin1 := uuid.New()
	admin2 := uuid.New()
	users.add(&entity.User{ID: admin1, IsActive: true, Role: entity.RoleAdmin})
	users.add(&entity.User{ID: admin2, IsActive: true, Role: entity.RoleAdmin})

	audit := &fakeAuditRepo{}
	svc := newSvc(users, newFakeCourseRepo(), newFakePaymentRepo(), audit, nil)
	if err := svc.ChangeRole(ctx, adminActor(), admin1, entity.RoleTeacher); err != nil {
		t.Fatalf("ChangeRole: %v", err)
	}
	u, _ := users.FindByID(ctx, admin1)
	if u.Role != entity.RoleTeacher {
		t.Errorf("expected role=teacher, got %s", u.Role)
	}
	if audit.countAction("user.change_role") != 1 {
		t.Error("expected 1 audit row for user.change_role")
	}
}

// ---- UnpublishCourse tests --------------------------------------------------

func TestUnpublishCourse_SetsUnpublishedAndAudits(t *testing.T) {
	ctx := context.Background()
	courses := newFakeCourseRepo()
	audit := &fakeAuditRepo{}
	notify := &fakeNotifier{}
	courseID := uuid.New()
	teacherID := uuid.New()
	courses.add(&entity.Course{ID: courseID, TeacherID: teacherID, IsPublished: true})

	svc := newSvc(newFakeUserRepo(), courses, newFakePaymentRepo(), audit, notify)
	if err := svc.UnpublishCourse(ctx, adminActor(), courseID, "low quality"); err != nil {
		t.Fatalf("UnpublishCourse: %v", err)
	}

	c, _ := courses.FindByID(ctx, courseID)
	if c.IsPublished {
		t.Error("expected IsPublished=false after unpublish")
	}
	if audit.countAction("course.unpublish") != 1 {
		t.Error("expected 1 audit row for course.unpublish")
	}
	if notify.count() == 0 {
		t.Error("expected teacher notification after unpublish")
	}
}

// ---- ManualRefund tests -----------------------------------------------------

func TestManualRefund_AuditsOnSuccess(t *testing.T) {
	ctx := context.Background()
	payments := newFakePaymentRepo()
	audit := &fakeAuditRepo{}
	paymentID := uuid.New()
	key := "key-1"
	payments.add(&entity.Payment{
		ID:             paymentID,
		Status:         entity.PayCompleted,
		IdempotencyKey: &key,
	})

	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), payments, audit, nil)
	if err := svc.ManualRefund(ctx, adminActor(), paymentID, "customer request"); err != nil {
		t.Fatalf("ManualRefund: %v", err)
	}
	if audit.countAction("payment.refund") != 1 {
		t.Error("expected 1 audit row for payment.refund")
	}
	p, _ := payments.GetByID(ctx, paymentID)
	if p.Status != entity.PayRefunded {
		t.Errorf("expected status=refunded, got %s", p.Status)
	}
}

func TestManualRefund_RejectsPendingPayment(t *testing.T) {
	ctx := context.Background()
	payments := newFakePaymentRepo()
	paymentID := uuid.New()
	key := "key-2"
	payments.add(&entity.Payment{
		ID:             paymentID,
		Status:         entity.PayPending,
		IdempotencyKey: &key,
	})

	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), payments, &fakeAuditRepo{}, nil)
	err := svc.ManualRefund(ctx, adminActor(), paymentID, "test")
	if err == nil {
		t.Error("expected error: cannot refund a pending payment")
	}
}

// ---- Overview tests ---------------------------------------------------------

func TestOverview_ReturnsRevenueTotals(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)

	from := time.Now().Add(-30 * 24 * time.Hour)
	to := time.Now()
	ov, err := svc.Overview(ctx, adminActor(), from, to)
	if err != nil {
		t.Fatalf("Overview: %v", err)
	}
	gross, _ := ov.GrossRevenue.Float64()
	if gross != 1000 {
		t.Errorf("expected gross=1000, got %.2f", gross)
	}
}

// ---- Broadcast tests --------------------------------------------------------

func TestBroadcast_NonAdminForbidden(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeUserRepo(), newFakeCourseRepo(), newFakePaymentRepo(), &fakeAuditRepo{}, nil)
	err := svc.Broadcast(ctx, studentActor(), BroadcastInput{Title: "hi", Segment: SegmentAll})
	if err == nil {
		t.Error("expected Forbidden for student actor")
	}
}

func TestBroadcast_AuditsOnSuccess(t *testing.T) {
	ctx := context.Background()
	audit := &fakeAuditRepo{}
	notify := &fakeNotifier{}
	users := newFakeUserRepo()
	users.add(&entity.User{ID: uuid.New(), IsActive: true, Role: entity.RoleStudent})

	svc := newSvc(users, newFakeCourseRepo(), newFakePaymentRepo(), audit, notify)
	if err := svc.Broadcast(ctx, adminActor(), BroadcastInput{
		Title:   "Platform update",
		Body:    "We've added new features.",
		Segment: SegmentStudents,
	}); err != nil {
		t.Fatalf("Broadcast: %v", err)
	}
	if audit.countAction("broadcast.send") != 1 {
		t.Error("expected 1 audit row for broadcast.send")
	}
	if notify.count() == 0 {
		t.Error("expected at least 1 notification sent to students")
	}
}

