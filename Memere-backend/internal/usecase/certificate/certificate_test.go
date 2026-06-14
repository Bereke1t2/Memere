package certificate

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// helpers

func completedActor() Actor { return Actor{UserID: uuid.New()} }

func setupFullCourse(progress *fakeProgressRepo, courses *fakeCourseRepo, users *fakeUserRepo, actor Actor, courseID uuid.UUID) {
	progress.setCourseProgress(actor.UserID, courseID, 100)
	courses.add(&entity.Course{ID: courseID, Title: "Grade 12 Mathematics"})
	users.add(&entity.User{ID: actor.UserID, FirstName: "Abel", LastName: "Tesfaye"})
}

// ---- Issue tests ------------------------------------------------------------

func TestIssue_BlockedUnder100Percent(t *testing.T) {
	ctx := context.Background()
	actor := completedActor()
	courseID := uuid.New()
	progress := newFakeProgressRepo()
	progress.setCourseProgress(actor.UserID, courseID, 75)

	svc := newSvc(newFakeCertRepo(), progress, newFakeCourseRepo(), newFakeUserRepo(), newFakeObjectStore(), &fakeRenderer{}, nil)
	_, err := svc.Issue(ctx, actor, courseID)
	if err == nil {
		t.Fatal("expected COURSE_NOT_COMPLETE error for 75% progress")
	}
	if !apperror.IsCode(err, "CONFLICT") {
		t.Errorf("expected code CONFLICT for incomplete course, got error: %v", err)
	}
}

func TestIssue_SucceedsAt100Percent(t *testing.T) {
	ctx := context.Background()
	actor := completedActor()
	courseID := uuid.New()
	progress := newFakeProgressRepo()
	courses := newFakeCourseRepo()
	users := newFakeUserRepo()
	setupFullCourse(progress, courses, users, actor, courseID)

	svc := newSvc(newFakeCertRepo(), progress, courses, users, newFakeObjectStore(), &fakeRenderer{}, nil)
	cert, err := svc.Issue(ctx, actor, courseID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if cert.StudentID != actor.UserID {
		t.Error("cert studentID mismatch")
	}
	if cert.CourseID != courseID {
		t.Error("cert courseID mismatch")
	}
	if cert.Serial == "" {
		t.Error("expected non-empty serial")
	}
}

func TestIssue_Idempotent(t *testing.T) {
	ctx := context.Background()
	actor := completedActor()
	courseID := uuid.New()
	progress := newFakeProgressRepo()
	courses := newFakeCourseRepo()
	users := newFakeUserRepo()
	setupFullCourse(progress, courses, users, actor, courseID)
	certRepo := newFakeCertRepo()

	svc := newSvc(certRepo, progress, courses, users, newFakeObjectStore(), &fakeRenderer{}, nil)
	cert1, err := svc.Issue(ctx, actor, courseID)
	if err != nil {
		t.Fatalf("Issue 1: %v", err)
	}
	cert2, err := svc.Issue(ctx, actor, courseID)
	if err != nil {
		t.Fatalf("Issue 2: %v", err)
	}
	if cert1.ID != cert2.ID {
		t.Error("expected same certificate on second issue (idempotent)")
	}
}

func TestIssue_UnauthenticatedReturnsUnauthorized(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeCertRepo(), newFakeProgressRepo(), newFakeCourseRepo(), newFakeUserRepo(), newFakeObjectStore(), &fakeRenderer{}, nil)
	_, err := svc.Issue(ctx, Actor{UserID: uuid.Nil}, uuid.New())
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
}

// ---- GetDownloadURL tests ---------------------------------------------------

func TestGetDownloadURL_OwnershipCheck(t *testing.T) {
	ctx := context.Background()
	owner := Actor{UserID: uuid.New()}
	thief := Actor{UserID: uuid.New()}
	courseID := uuid.New()
	progress := newFakeProgressRepo()
	courses := newFakeCourseRepo()
	users := newFakeUserRepo()
	setupFullCourse(progress, courses, users, owner, courseID)
	// Also add thief to users so lookup works.
	users.add(&entity.User{ID: thief.UserID, FirstName: "X", LastName: "X"})
	store := newFakeObjectStore()
	renderer := &fakeRenderer{}
	certRepo := newFakeCertRepo()

	svc := newSvc(certRepo, progress, courses, users, store, renderer, nil)

	// Issue a cert for owner.
	cert, err := svc.Issue(ctx, owner, courseID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	// Wait briefly for the background render.
	time.Sleep(50 * time.Millisecond)

	// Thief cannot download.
	_, err = svc.GetDownloadURL(ctx, thief, cert.ID)
	if err == nil {
		t.Fatal("expected forbidden error for non-owner")
	}
}

func TestGetDownloadURL_ReturnSignedURL(t *testing.T) {
	ctx := context.Background()
	actor := completedActor()
	courseID := uuid.New()
	progress := newFakeProgressRepo()
	courses := newFakeCourseRepo()
	users := newFakeUserRepo()
	setupFullCourse(progress, courses, users, actor, courseID)
	store := newFakeObjectStore()
	renderer := &fakeRenderer{}
	certRepo := newFakeCertRepo()

	svc := newSvc(certRepo, progress, courses, users, store, renderer, nil)

	cert, err := svc.Issue(ctx, actor, courseID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	// Render synchronously so file_key is set before GetDownloadURL.
	if syncErr := svc.renderAndStoreSync(ctx, cert); syncErr != nil {
		t.Fatalf("renderAndStoreSync: %v", syncErr)
	}

	url, err := svc.GetDownloadURL(ctx, actor, cert.ID)
	if err != nil {
		t.Fatalf("GetDownloadURL: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty signed URL")
	}
}

// ---- VerifyCertificate tests ------------------------------------------------

func TestVerifyCertificate_BySerial(t *testing.T) {
	ctx := context.Background()
	actor := completedActor()
	courseID := uuid.New()
	progress := newFakeProgressRepo()
	courses := newFakeCourseRepo()
	users := newFakeUserRepo()
	setupFullCourse(progress, courses, users, actor, courseID)
	certRepo := newFakeCertRepo()

	svc := newSvc(certRepo, progress, courses, users, newFakeObjectStore(), &fakeRenderer{}, nil)
	cert, err := svc.Issue(ctx, actor, courseID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	view, err := svc.VerifyCertificate(ctx, cert.Serial)
	if err != nil {
		t.Fatalf("VerifyCertificate: %v", err)
	}
	if !view.Valid {
		t.Error("expected valid=true")
	}
	if view.StudentName != "Abel Tesfaye" {
		t.Errorf("expected student name 'Abel Tesfaye', got %q", view.StudentName)
	}
	if view.CourseTitle != "Grade 12 Mathematics" {
		t.Errorf("expected course title, got %q", view.CourseTitle)
	}
}

func TestVerifyCertificate_UnknownSerialReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	svc := newSvc(newFakeCertRepo(), newFakeProgressRepo(), newFakeCourseRepo(), newFakeUserRepo(), newFakeObjectStore(), &fakeRenderer{}, nil)
	_, err := svc.VerifyCertificate(ctx, "MEM-2026-DOESNOTEXIST")
	if err == nil {
		t.Fatal("expected not-found error for unknown serial")
	}
	if !apperror.IsNotFound(err) {
		t.Errorf("expected not-found, got: %v", err)
	}
}
