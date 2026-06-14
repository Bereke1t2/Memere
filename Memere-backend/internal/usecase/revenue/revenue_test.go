package revenue

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ---- fakes -------------------------------------------------------------------

type fakeRevenueRepo struct {
	completed  repository.RevenueTotals
	refunded   repository.RevenueTotals
	teacher    repository.RevenueTotals
	courseSale repository.RevenueTotals
	byProvider []repository.ProviderRevenue
}

var _ repository.RevenueRepository = (*fakeRevenueRepo)(nil)

func (f *fakeRevenueRepo) CompletedTotals(context.Context, time.Time, time.Time) (repository.RevenueTotals, error) {
	return f.completed, nil
}
func (f *fakeRevenueRepo) RefundedTotals(context.Context, time.Time, time.Time) (repository.RevenueTotals, error) {
	return f.refunded, nil
}
func (f *fakeRevenueRepo) TeacherGross(context.Context, uuid.UUID, time.Time, time.Time) (repository.RevenueTotals, error) {
	return f.teacher, nil
}
func (f *fakeRevenueRepo) CourseSales(context.Context, uuid.UUID) (repository.RevenueTotals, error) {
	return f.courseSale, nil
}
func (f *fakeRevenueRepo) CompletedByProvider(context.Context, time.Time, time.Time) ([]repository.ProviderRevenue, error) {
	return f.byProvider, nil
}

// fakeCourseRepo embeds the interface so only FindByID is implemented; any other
// method call would panic, flagging an unexpected dependency.
type fakeCourseRepo struct {
	repository.CourseRepository
	byID map[uuid.UUID]*entity.Course
}

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	if c, ok := f.byID[id]; ok {
		return c, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

func totals(gross string, units int64) repository.RevenueTotals {
	return repository.RevenueTotals{Gross: decimal.RequireFromString(gross), Units: units}
}

func window() (time.Time, time.Time) {
	from := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC)
	return from, to
}

// ---- PlatformRevenue ---------------------------------------------------------

func TestPlatformRevenue_AdminOnly(t *testing.T) {
	rev := &fakeRevenueRepo{
		completed: totals("1000", 10),
		refunded:  totals("150", 2),
		byProvider: []repository.ProviderRevenue{
			{Provider: entity.ProviderChapa, Gross: decimal.NewFromInt(700), Units: 7},
			{Provider: entity.ProviderStripe, Gross: decimal.NewFromInt(300), Units: 3},
		},
	}
	svc := NewService(rev, &fakeCourseRepo{}, Config{})
	from, to := window()

	t.Run("admin", func(t *testing.T) {
		rep, err := svc.PlatformRevenue(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleAdmin}, from, to)
		if err != nil {
			t.Fatalf("admin: %v", err)
		}
		if rep.Gross != "1000.00" || rep.Refunds != "150.00" || rep.Net != "850.00" {
			t.Fatalf("got gross=%s refunds=%s net=%s", rep.Gross, rep.Refunds, rep.Net)
		}
		if rep.Units != 10 || rep.RefundUnits != 2 {
			t.Fatalf("got units=%d refundUnits=%d", rep.Units, rep.RefundUnits)
		}
		if len(rep.ByProvider) != 2 || rep.ByProvider[0].Provider != "chapa" || rep.ByProvider[0].Gross != "700.00" {
			t.Fatalf("provider breakdown wrong: %+v", rep.ByProvider)
		}
	})

	t.Run("teacher forbidden", func(t *testing.T) {
		_, err := svc.PlatformRevenue(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleTeacher}, from, to)
		if !apperror.IsCode(err, "FORBIDDEN") {
			t.Fatalf("want FORBIDDEN, got %v", err)
		}
	})

	t.Run("anonymous unauthorized", func(t *testing.T) {
		_, err := svc.PlatformRevenue(context.Background(), Actor{}, from, to)
		if !apperror.IsCode(err, "UNAUTHORIZED") {
			t.Fatalf("want UNAUTHORIZED, got %v", err)
		}
	})

	t.Run("inverted window", func(t *testing.T) {
		_, err := svc.PlatformRevenue(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleAdmin}, to, from)
		if !apperror.IsCode(err, "INVALID_WINDOW") {
			t.Fatalf("want INVALID_WINDOW, got %v", err)
		}
	})
}

// ---- TeacherEarnings ---------------------------------------------------------

func TestTeacherEarnings_Split(t *testing.T) {
	teacher := uuid.New()
	rev := &fakeRevenueRepo{teacher: totals("1000", 8)}
	svc := NewService(rev, &fakeCourseRepo{}, Config{}) // default 0.70 share
	from, to := window()

	t.Run("teacher reads own (70/30)", func(t *testing.T) {
		rep, err := svc.TeacherEarnings(context.Background(), Actor{UserID: teacher, Role: entity.RoleTeacher}, teacher, from, to)
		if err != nil {
			t.Fatalf("own earnings: %v", err)
		}
		if rep.Gross != "1000.00" || rep.Earnings != "700.00" || rep.PlatformFee != "300.00" {
			t.Fatalf("got gross=%s earnings=%s fee=%s", rep.Gross, rep.Earnings, rep.PlatformFee)
		}
		if rep.Units != 8 {
			t.Fatalf("units: got %d, want 8", rep.Units)
		}
	})

	t.Run("admin reads any", func(t *testing.T) {
		_, err := svc.TeacherEarnings(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleAdmin}, teacher, from, to)
		if err != nil {
			t.Fatalf("admin earnings: %v", err)
		}
	})

	t.Run("other teacher forbidden", func(t *testing.T) {
		_, err := svc.TeacherEarnings(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleTeacher}, teacher, from, to)
		if !apperror.IsCode(err, "FORBIDDEN") {
			t.Fatalf("want FORBIDDEN, got %v", err)
		}
	})
}

func TestTeacherEarnings_ConfigurableShare(t *testing.T) {
	teacher := uuid.New()
	rev := &fakeRevenueRepo{teacher: totals("1000", 4)}
	svc := NewService(rev, &fakeCourseRepo{}, Config{TeacherShare: decimal.NewFromFloat(0.80)})
	from, to := window()

	rep, err := svc.TeacherEarnings(context.Background(), Actor{UserID: teacher, Role: entity.RoleTeacher}, teacher, from, to)
	if err != nil {
		t.Fatalf("earnings: %v", err)
	}
	if rep.Earnings != "800.00" || rep.PlatformFee != "200.00" {
		t.Fatalf("configurable split: got earnings=%s fee=%s, want 800/200", rep.Earnings, rep.PlatformFee)
	}
}

func TestTeacherEarnings_RoundingSumsToGross(t *testing.T) {
	teacher := uuid.New()
	// 99.99 * 0.70 = 69.993 -> rounds to 69.99; fee must be the remainder 30.00.
	rev := &fakeRevenueRepo{teacher: totals("99.99", 1)}
	svc := NewService(rev, &fakeCourseRepo{}, Config{})
	from, to := window()

	rep, err := svc.TeacherEarnings(context.Background(), Actor{UserID: teacher, Role: entity.RoleTeacher}, teacher, from, to)
	if err != nil {
		t.Fatalf("earnings: %v", err)
	}
	earn := decimal.RequireFromString(rep.Earnings)
	fee := decimal.RequireFromString(rep.PlatformFee)
	if !earn.Add(fee).Equal(decimal.RequireFromString("99.99")) {
		t.Fatalf("earnings+fee must equal gross: %s + %s", rep.Earnings, rep.PlatformFee)
	}
}

// ---- CourseSalesStats --------------------------------------------------------

func TestCourseSalesStats_OwnerOrAdmin(t *testing.T) {
	teacher := uuid.New()
	courseID := uuid.New()
	courses := &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{
		courseID: {ID: courseID, TeacherID: teacher},
	}}
	rev := &fakeRevenueRepo{courseSale: totals("2500", 25)}
	svc := NewService(rev, courses, Config{})

	t.Run("owner", func(t *testing.T) {
		rep, err := svc.CourseSalesStats(context.Background(), Actor{UserID: teacher, Role: entity.RoleTeacher}, courseID)
		if err != nil {
			t.Fatalf("owner: %v", err)
		}
		if rep.Gross != "2500.00" || rep.Units != 25 {
			t.Fatalf("got gross=%s units=%d", rep.Gross, rep.Units)
		}
	})

	t.Run("admin", func(t *testing.T) {
		_, err := svc.CourseSalesStats(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleAdmin}, courseID)
		if err != nil {
			t.Fatalf("admin: %v", err)
		}
	})

	t.Run("non-owner forbidden", func(t *testing.T) {
		_, err := svc.CourseSalesStats(context.Background(), Actor{UserID: uuid.New(), Role: entity.RoleTeacher}, courseID)
		if !apperror.IsCode(err, "FORBIDDEN") {
			t.Fatalf("want FORBIDDEN, got %v", err)
		}
	})

	t.Run("missing course", func(t *testing.T) {
		_, err := svc.CourseSalesStats(context.Background(), Actor{UserID: teacher, Role: entity.RoleTeacher}, uuid.New())
		if !apperror.IsNotFound(err) {
			t.Fatalf("want NotFound, got %v", err)
		}
	})
}
