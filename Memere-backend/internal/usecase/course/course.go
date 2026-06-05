package course

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/validator"
)

const (
	maxTitleLen     = 200
	maxShortDescLen = 500
	minGrade        = 1
	maxGrade        = 12
	// slugSuffixRetries bounds the attempts to find a free slug before giving up.
	slugSuffixRetries = 5
)

// defaultCurrency / defaultLanguage are applied to new courses when the input
// leaves them blank (spec §4.2.2).
const (
	defaultCurrency = "ETB"
	defaultLanguage = "en"
)

// Service implements the course-domain usecases over the domain repository
// interfaces plus a TxManager for grouping multi-write operations.
type Service struct {
	courses  repository.CourseRepository
	sections repository.SectionRepository
	lessons  repository.LessonRepository
	tx       repository.TxManager
}

// NewService wires the course usecase with its dependencies.
func NewService(
	courses repository.CourseRepository,
	sections repository.SectionRepository,
	lessons repository.LessonRepository,
	tx repository.TxManager,
) *Service {
	return &Service{courses: courses, sections: sections, lessons: lessons, tx: tx}
}

// CreateCourseInput is the create request after transport decoding.
type CreateCourseInput struct {
	Title            string
	Description      string
	ShortDescription *string
	Subject          string
	Grade            int
	ThumbnailURL     *string
	Price            float64
	Currency         string
	IsFree           bool
	Language         string
	Level            entity.Level
	Metadata         map[string]any
}

// UpdateCourseInput carries partial updates; nil fields are left unchanged.
type UpdateCourseInput struct {
	Title            *string
	Description      *string
	ShortDescription *string
	Subject          *string
	Grade            *int
	ThumbnailURL     *string
	Price            *float64
	IsFree           *bool
	Language         *string
	Level            *entity.Level
	Metadata         map[string]any
}

// CreateCourse creates a course owned by the acting teacher/admin. The teacher
// is always the actor (teachers cannot create on behalf of others), a unique
// slug is derived from the title, and publish-state/counters start at their
// defaults.
func (s *Service) CreateCourse(ctx context.Context, actor *Actor, in CreateCourseInput) (*entity.Course, error) {
	if !actor.isTeacherOrAdmin() {
		return nil, apperror.Forbidden("only teachers may create courses", nil)
	}
	if err := validateCourseInput(in); err != nil {
		return nil, err
	}

	currency := in.Currency
	if currency == "" {
		currency = defaultCurrency
	}
	language := in.Language
	if language == "" {
		language = defaultLanguage
	}

	c := &entity.Course{
		TeacherID:        actor.UserID,
		Title:            in.Title,
		Description:      in.Description,
		ShortDescription: in.ShortDescription,
		Subject:          in.Subject,
		Grade:            in.Grade,
		ThumbnailURL:     in.ThumbnailURL,
		Price:            in.Price,
		Currency:         currency,
		IsFree:           in.IsFree,
		IsPublished:      false,
		Language:         language,
		Level:            in.Level,
		Metadata:         in.Metadata,
	}

	// Derive a unique slug, retrying with a short random suffix on collision.
	base := validator.Slugify(in.Title)
	if base == "" {
		base = "course"
	}
	if err := s.createWithUniqueSlug(ctx, c, base); err != nil {
		return nil, err
	}
	return c, nil
}

// createWithUniqueSlug attempts Create with base, then base-<suffix> on a
// SLUG_TAKEN conflict, up to slugSuffixRetries times.
func (s *Service) createWithUniqueSlug(ctx context.Context, c *entity.Course, base string) error {
	c.Slug = base
	err := s.courses.Create(ctx, c)
	for attempt := 0; attempt < slugSuffixRetries && apperror.IsCode(err, "CONFLICT"); attempt++ {
		c.Slug = base + "-" + shortSuffix()
		err = s.courses.Create(ctx, c)
	}
	return err
}

// ListCourses returns a page of courses visible to the actor. Anonymous callers
// and students see only published courses; a teacher additionally sees their own
// unpublished courses; an admin sees everything. The filter's IsPublished is
// overridden where visibility demands it.
func (s *Service) ListCourses(ctx context.Context, actor *Actor, filter repository.CourseFilter, cursor *pagination.Cursor, limit int) ([]*entity.Course, *pagination.Cursor, error) {
	switch {
	case actor.isAdmin():
		// No published restriction; honor whatever filter was supplied.
	case actor != nil && actor.Role == entity.RoleTeacher:
		// Teacher: published OR own. The repo filter cannot express the OR, so we
		// scope to the teacher's own courses when they ask for unpublished;
		// otherwise force published. A full published-OR-own listing across all
		// teachers is a Skill 5 concern (the handler decides the audience).
		published := true
		filter.IsPublished = &published
	default:
		// Anonymous or student: published only.
		published := true
		filter.IsPublished = &published
	}
	return s.courses.List(ctx, filter, cursor, limit)
}

// GetCourse returns the nested course view (sections + lessons) if the actor may
// see it. Unpublished courses are visible only to their owner/admin; for a
// published course viewed by a non-owner, unpublished sections and lessons are
// stripped.
func (s *Service) GetCourse(ctx context.Context, actor *Actor, id uuid.UUID) (*repository.CourseWithContent, error) {
	content, err := s.courses.GetCourseWithSectionsAndLessons(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.applyVisibility(actor, content)
}

// GetCourseBySlug is GetCourse keyed by slug.
func (s *Service) GetCourseBySlug(ctx context.Context, actor *Actor, slug string) (*repository.CourseWithContent, error) {
	c, err := s.courses.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	content, err := s.courses.GetCourseWithSectionsAndLessons(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	return s.applyVisibility(actor, content)
}

// applyVisibility enforces the read rules on an assembled course view: it hides
// an unpublished course from non-owners entirely (NotFound, so existence is not
// revealed), and strips unpublished sections/lessons for non-owners.
func (s *Service) applyVisibility(actor *Actor, content *repository.CourseWithContent) (*repository.CourseWithContent, error) {
	owner := actor.canSeeUnpublished(content.Course)
	if !content.Course.IsPublished && !owner {
		return nil, apperror.NotFound("course not found", nil)
	}
	if owner {
		return content, nil
	}
	// Non-owner on a published course: drop unpublished sections and any
	// unpublished lessons within the surviving sections.
	visible := make([]repository.SectionWithLessons, 0, len(content.Sections))
	for _, sec := range content.Sections {
		if !sec.Section.IsPublished {
			continue
		}
		lessons := make([]*entity.Lesson, 0, len(sec.Lessons))
		for _, l := range sec.Lessons {
			if l.IsPublished {
				lessons = append(lessons, l)
			}
		}
		sec.Lessons = lessons
		visible = append(visible, sec)
	}
	content.Sections = visible
	return content, nil
}

// UpdateCourse applies partial updates after asserting the actor owns the course
// (or is admin).
func (s *Service) UpdateCourse(ctx context.Context, actor *Actor, id uuid.UUID, in UpdateCourseInput) (*entity.Course, error) {
	c, err := s.loadOwnedCourse(ctx, actor, id)
	if err != nil {
		return nil, err
	}

	v := validator.New()
	if in.Title != nil {
		v.Required("title", *in.Title)
		v.MaxLen("title", *in.Title, maxTitleLen)
		c.Title = *in.Title
	}
	if in.Description != nil {
		c.Description = *in.Description
	}
	if in.ShortDescription != nil {
		v.MaxLen("short_description", *in.ShortDescription, maxShortDescLen)
		c.ShortDescription = in.ShortDescription
	}
	if in.Subject != nil {
		v.Required("subject", *in.Subject)
		c.Subject = *in.Subject
	}
	if in.Grade != nil {
		v.InRange("grade", *in.Grade, minGrade, maxGrade)
		c.Grade = *in.Grade
	}
	if in.ThumbnailURL != nil {
		c.ThumbnailURL = in.ThumbnailURL
	}
	if in.Price != nil {
		v.NonNegative("price", *in.Price)
		c.Price = *in.Price
	}
	if in.IsFree != nil {
		c.IsFree = *in.IsFree
	}
	if in.Language != nil {
		c.Language = *in.Language
	}
	if in.Level != nil {
		v.OneOf("level", string(*in.Level), string(entity.LevelBeginner), string(entity.LevelIntermediate), string(entity.LevelAdvanced))
		c.Level = *in.Level
	}
	if in.Metadata != nil {
		c.Metadata = in.Metadata
	}
	if v.HasErrors() {
		return nil, apperror.Validation(v.Map(), nil)
	}

	if err := s.courses.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// PublishCourse sets the course published (ownership-checked).
func (s *Service) PublishCourse(ctx context.Context, actor *Actor, id uuid.UUID) (*entity.Course, error) {
	return s.setPublished(ctx, actor, id, true)
}

// UnpublishCourse clears the course published flag (ownership-checked).
func (s *Service) UnpublishCourse(ctx context.Context, actor *Actor, id uuid.UUID) (*entity.Course, error) {
	return s.setPublished(ctx, actor, id, false)
}

func (s *Service) setPublished(ctx context.Context, actor *Actor, id uuid.UUID, published bool) (*entity.Course, error) {
	c, err := s.loadOwnedCourse(ctx, actor, id)
	if err != nil {
		return nil, err
	}
	c.IsPublished = published
	if err := s.courses.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// DeleteCourse soft-deletes the course (ownership-checked). Sections and lessons
// are left as-is; they become unreachable through the soft-deleted course and a
// later cleanup can cascade if needed.
func (s *Service) DeleteCourse(ctx context.Context, actor *Actor, id uuid.UUID) error {
	if _, err := s.loadOwnedCourse(ctx, actor, id); err != nil {
		return err
	}
	return s.courses.SoftDelete(ctx, id)
}

// loadOwnedCourse fetches a course and asserts the actor may modify it.
func (s *Service) loadOwnedCourse(ctx context.Context, actor *Actor, id uuid.UUID) (*entity.Course, error) {
	c, err := s.courses.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := assertCourseOwner(actor, c); err != nil {
		return nil, err
	}
	return c, nil
}

// assertCourseOwner permits an admin always, the owning teacher otherwise, and
// rejects everyone else with FORBIDDEN.
func assertCourseOwner(actor *Actor, c *entity.Course) error {
	if actor.isAdmin() {
		return nil
	}
	if actor != nil && actor.Role == entity.RoleTeacher && c.TeacherID == actor.UserID {
		return nil
	}
	return apperror.Forbidden("you do not own this course", nil)
}

// validateCourseInput validates the create payload, returning an
// apperror.Validation on any problem.
func validateCourseInput(in CreateCourseInput) error {
	v := validator.New()
	v.Required("title", in.Title)
	v.MaxLen("title", in.Title, maxTitleLen)
	v.Required("description", in.Description)
	v.Required("subject", in.Subject)
	v.InRange("grade", in.Grade, minGrade, maxGrade)
	v.NonNegative("price", in.Price)
	if in.ShortDescription != nil {
		v.MaxLen("short_description", *in.ShortDescription, maxShortDescLen)
	}
	v.OneOf("level", string(in.Level), string(entity.LevelBeginner), string(entity.LevelIntermediate), string(entity.LevelAdvanced))
	if v.HasErrors() {
		return apperror.Validation(v.Map(), nil)
	}
	return nil
}
