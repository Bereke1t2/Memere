package course

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/validator"
)

// SectionInput is the create/update payload for a section.
type SectionInput struct {
	Title       string
	Description *string
	IsPublished bool
	// OrderIndex is optional on create; when nil the section is appended after
	// the course's current last section.
	OrderIndex *int
}

// AddSection adds a section to a course the actor owns. When OrderIndex is
// omitted it is auto-assigned as (current max + 1) so sections append in order.
func (s *Service) AddSection(ctx context.Context, actor *Actor, courseID uuid.UUID, in SectionInput) (*entity.CourseSection, error) {
	if _, err := s.loadOwnedCourse(ctx, actor, courseID); err != nil {
		return nil, err
	}
	if err := validateSectionInput(in); err != nil {
		return nil, err
	}

	orderIndex := 0
	if in.OrderIndex != nil {
		orderIndex = *in.OrderIndex
	} else {
		max, err := s.sections.MaxOrderIndex(ctx, courseID)
		if err != nil {
			return nil, err
		}
		orderIndex = max + 1
	}

	sec := &entity.CourseSection{
		CourseID:    courseID,
		Title:       in.Title,
		Description: in.Description,
		OrderIndex:  orderIndex,
		IsPublished: in.IsPublished,
	}
	if err := s.sections.Create(ctx, sec); err != nil {
		return nil, err
	}
	return sec, nil
}

// ListSections returns a course's sections. Visibility mirrors the course: a
// non-owner viewing a published course sees only published sections; the
// unpublished-course case is gated by loading the course first.
func (s *Service) ListSections(ctx context.Context, actor *Actor, courseID uuid.UUID) ([]*entity.CourseSection, error) {
	c, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	owner := actor.canSeeUnpublished(c)
	if !c.IsPublished && !owner {
		return nil, apperror.NotFound("course not found", nil)
	}

	sections, err := s.sections.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if owner {
		return sections, nil
	}
	visible := make([]*entity.CourseSection, 0, len(sections))
	for _, sec := range sections {
		if sec.IsPublished {
			visible = append(visible, sec)
		}
	}
	return visible, nil
}

// UpdateSection applies updates after asserting ownership via the parent course.
func (s *Service) UpdateSection(ctx context.Context, actor *Actor, sectionID uuid.UUID, in SectionInput) (*entity.CourseSection, error) {
	sec, err := s.loadOwnedSection(ctx, actor, sectionID)
	if err != nil {
		return nil, err
	}
	if err := validateSectionInput(in); err != nil {
		return nil, err
	}
	sec.Title = in.Title
	sec.Description = in.Description
	sec.IsPublished = in.IsPublished
	if in.OrderIndex != nil {
		sec.OrderIndex = *in.OrderIndex
	}
	if err := s.sections.Update(ctx, sec); err != nil {
		return nil, err
	}
	return sec, nil
}

// DeleteSection soft-deletes a section and cascades the soft-delete to its
// lessons, then recomputes the course counters — all in one transaction so the
// course's total_lessons / total_duration_seconds stay consistent.
func (s *Service) DeleteSection(ctx context.Context, actor *Actor, sectionID uuid.UUID) error {
	sec, err := s.loadOwnedSection(ctx, actor, sectionID)
	if err != nil {
		return err
	}
	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.lessons.SoftDeleteBySection(ctx, sectionID); err != nil {
			return err
		}
		if err := s.sections.SoftDelete(ctx, sectionID); err != nil {
			return err
		}
		return s.courses.RecomputeCounters(ctx, sec.CourseID)
	})
}

// loadOwnedSection fetches a section and asserts the actor owns its course.
func (s *Service) loadOwnedSection(ctx context.Context, actor *Actor, sectionID uuid.UUID) (*entity.CourseSection, error) {
	sec, err := s.sections.FindByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	c, err := s.courses.FindByID(ctx, sec.CourseID)
	if err != nil {
		return nil, err
	}
	if err := assertCourseOwner(actor, c); err != nil {
		return nil, err
	}
	return sec, nil
}

// validateSectionInput validates a section payload.
func validateSectionInput(in SectionInput) error {
	v := validator.New()
	v.Required("title", in.Title)
	v.MaxLen("title", in.Title, maxTitleLen)
	if in.OrderIndex != nil {
		v.InRange("order_index", *in.OrderIndex, 0, 1_000_000)
	}
	if v.HasErrors() {
		return apperror.Validation(v.Map(), nil)
	}
	return nil
}
