package course

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/validator"
)

// LessonInput is the create/update payload for a lesson.
type LessonInput struct {
	Title           string
	Type            entity.LessonType
	IsFreePreview   bool
	DurationSeconds int
	IsPublished     bool
	// OrderIndex is optional on create; when nil the lesson is appended after the
	// section's current last lesson.
	OrderIndex *int
}

// AddLesson adds a lesson to a section the actor owns. Creating the lesson and
// recomputing the parent course's counters run in one transaction so
// total_lessons / total_duration_seconds never drift from the lessons.
func (s *Service) AddLesson(ctx context.Context, actor *Actor, sectionID uuid.UUID, in LessonInput) (*entity.Lesson, error) {
	sec, err := s.loadOwnedSection(ctx, actor, sectionID)
	if err != nil {
		return nil, err
	}
	if err := validateLessonInput(in); err != nil {
		return nil, err
	}

	orderIndex := 0
	if in.OrderIndex != nil {
		orderIndex = *in.OrderIndex
	} else {
		max, err := s.lessons.MaxOrderIndex(ctx, sectionID)
		if err != nil {
			return nil, err
		}
		orderIndex = max + 1
	}

	l := &entity.Lesson{
		SectionID:       sectionID,
		CourseID:        sec.CourseID,
		Title:           in.Title,
		Type:            in.Type,
		OrderIndex:      orderIndex,
		IsFreePreview:   in.IsFreePreview,
		DurationSeconds: in.DurationSeconds,
		IsPublished:     in.IsPublished,
	}
	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.lessons.Create(ctx, l); err != nil {
			return err
		}
		return s.courses.RecomputeCounters(ctx, sec.CourseID)
	})
	if err != nil {
		return nil, err
	}
	return l, nil
}

// ListLessons returns a section's lessons. A non-owner viewing a published
// course sees only published lessons; unpublished course/section access is gated
// by loading the section's course first.
func (s *Service) ListLessons(ctx context.Context, actor *Actor, sectionID uuid.UUID) ([]*entity.Lesson, error) {
	sec, err := s.sections.FindByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	c, err := s.courses.FindByID(ctx, sec.CourseID)
	if err != nil {
		return nil, err
	}
	owner := actor.canSeeUnpublished(c)
	if (!c.IsPublished || !sec.IsPublished) && !owner {
		return nil, apperror.NotFound("section not found", nil)
	}

	lessons, err := s.lessons.ListBySection(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if owner {
		return lessons, nil
	}
	visible := make([]*entity.Lesson, 0, len(lessons))
	for _, l := range lessons {
		if l.IsPublished {
			visible = append(visible, l)
		}
	}
	return visible, nil
}

// UpdateLesson applies updates after asserting ownership via the lesson's
// course. When the duration changes the course counters are recomputed in the
// same transaction as the update.
func (s *Service) UpdateLesson(ctx context.Context, actor *Actor, lessonID uuid.UUID, in LessonInput) (*entity.Lesson, error) {
	l, err := s.loadOwnedLesson(ctx, actor, lessonID)
	if err != nil {
		return nil, err
	}
	if err := validateLessonInput(in); err != nil {
		return nil, err
	}

	durationChanged := l.DurationSeconds != in.DurationSeconds
	l.Title = in.Title
	l.Type = in.Type
	l.IsFreePreview = in.IsFreePreview
	l.DurationSeconds = in.DurationSeconds
	l.IsPublished = in.IsPublished
	if in.OrderIndex != nil {
		l.OrderIndex = *in.OrderIndex
	}

	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.lessons.Update(ctx, l); err != nil {
			return err
		}
		if durationChanged {
			return s.courses.RecomputeCounters(ctx, l.CourseID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLesson soft-deletes a lesson and recomputes the course counters in one
// transaction.
func (s *Service) DeleteLesson(ctx context.Context, actor *Actor, lessonID uuid.UUID) error {
	l, err := s.loadOwnedLesson(ctx, actor, lessonID)
	if err != nil {
		return err
	}
	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.lessons.SoftDelete(ctx, lessonID); err != nil {
			return err
		}
		return s.courses.RecomputeCounters(ctx, l.CourseID)
	})
}

// loadOwnedLesson fetches a lesson and asserts the actor owns its course.
func (s *Service) loadOwnedLesson(ctx context.Context, actor *Actor, lessonID uuid.UUID) (*entity.Lesson, error) {
	l, err := s.lessons.FindByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	c, err := s.courses.FindByID(ctx, l.CourseID)
	if err != nil {
		return nil, err
	}
	if err := assertCourseOwner(actor, c); err != nil {
		return nil, err
	}
	return l, nil
}

// validateLessonInput validates a lesson payload, including the type enum.
func validateLessonInput(in LessonInput) error {
	v := validator.New()
	v.Required("title", in.Title)
	v.MaxLen("title", in.Title, maxTitleLen)
	v.OneOf("type", string(in.Type),
		string(entity.LessonTypeVideo), string(entity.LessonTypeNote),
		string(entity.LessonTypeQuiz), string(entity.LessonTypeMixed))
	v.NonNegative("duration_seconds", float64(in.DurationSeconds))
	if in.OrderIndex != nil {
		v.InRange("order_index", *in.OrderIndex, 0, 1_000_000)
	}
	if v.HasErrors() {
		return apperror.Validation(v.Map(), nil)
	}
	return nil
}
