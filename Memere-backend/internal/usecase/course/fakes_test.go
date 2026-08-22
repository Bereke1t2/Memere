package course

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// fakeCourseRepo is an in-memory CourseRepository keyed by id (and slug). It
// records soft-deletes by setting DeletedAt — never removing the row — so tests
// can assert the soft-delete rule. RecomputeCounters recomputes from the shared
// lesson repo to exercise the counter-maintenance path.
type fakeCourseRepo struct {
	mu       sync.Mutex
	byID     map[uuid.UUID]*entity.Course
	sections *fakeSectionRepo // shared, for nested view assembly
	lessons  *fakeLessonRepo  // shared, for RecomputeCounters and nested view
}

func newFakeCourseRepo(sections *fakeSectionRepo, lessons *fakeLessonRepo) *fakeCourseRepo {
	return &fakeCourseRepo{byID: map[uuid.UUID]*entity.Course{}, sections: sections, lessons: lessons}
}

func (f *fakeCourseRepo) Create(_ context.Context, c *entity.Course) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ex := range f.byID {
		if ex.DeletedAt == nil && ex.Slug == c.Slug {
			return apperror.Conflict("SLUG_TAKEN", nil)
		}
	}
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt
	cp := *c
	f.byID[c.ID] = &cp
	return nil
}

func (f *fakeCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok && c.DeletedAt == nil {
		cp := *c
		return &cp, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

func (f *fakeCourseRepo) FindBySlug(_ context.Context, slug string) (*entity.Course, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.byID {
		if c.DeletedAt == nil && c.Slug == slug {
			cp := *c
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("course not found", nil)
}

func (f *fakeCourseRepo) List(_ context.Context, filter repository.CourseFilter, cursor *pagination.Cursor, limit int) ([]*entity.Course, *pagination.Cursor, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	limit = pagination.NormalizeLimit(limit)

	var rows []*entity.Course
	for _, c := range f.byID {
		if c.DeletedAt != nil {
			continue
		}
		if filter.IsPublished != nil && c.IsPublished != *filter.IsPublished {
			continue
		}
		if filter.TeacherID != nil && c.TeacherID != *filter.TeacherID {
			continue
		}
		if filter.Subject != nil && c.Subject != *filter.Subject {
			continue
		}
		if filter.Grade != nil && c.Grade != *filter.Grade {
			continue
		}
		cp := *c
		rows = append(rows, &cp)
	}
	// Order by (created_at, id) DESC for a stable keyset.
	sort.Slice(rows, func(i, j int) bool {
		if !rows[i].CreatedAt.Equal(rows[j].CreatedAt) {
			return rows[i].CreatedAt.After(rows[j].CreatedAt)
		}
		return rows[i].ID.String() > rows[j].ID.String()
	})
	// Apply cursor (skip everything up to and including the cursor row).
	if cursor != nil {
		filtered := rows[:0]
		for _, c := range rows {
			if c.CreatedAt.Before(cursor.CreatedAt) ||
				(c.CreatedAt.Equal(cursor.CreatedAt) && c.ID.String() < cursor.ID.String()) {
				filtered = append(filtered, c)
			}
		}
		rows = filtered
	}

	var next *pagination.Cursor
	if len(rows) > limit {
		last := rows[limit-1]
		next = &pagination.Cursor{CreatedAt: last.CreatedAt, ID: last.ID}
		rows = rows[:limit]
	}
	return rows, next, nil
}

func (f *fakeCourseRepo) Update(_ context.Context, c *entity.Course) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	ex, ok := f.byID[c.ID]
	if !ok || ex.DeletedAt != nil {
		return apperror.NotFound("course not found", nil)
	}
	c.UpdatedAt = time.Now()
	cp := *c
	f.byID[c.ID] = &cp
	return nil
}

func (f *fakeCourseRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if c, ok := f.byID[id]; ok {
		now := time.Now()
		c.DeletedAt = &now
	}
	return nil
}

func (f *fakeCourseRepo) RecomputeCounters(_ context.Context, courseID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.byID[courseID]
	if !ok {
		return apperror.NotFound("course not found", nil)
	}
	count, dur := f.lessons.statsForCourse(courseID)
	c.TotalLessons = count
	c.TotalDurationSeconds = dur
	return nil
}

func (f *fakeCourseRepo) GetCourseWithSectionsAndLessons(ctx context.Context, courseID uuid.UUID) (*repository.CourseWithContent, error) {
	c, err := f.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	sections, _ := f.sections.ListByCourse(ctx, courseID)
	out := make([]repository.SectionWithLessons, len(sections))
	for i, sec := range sections {
		lessons, _ := f.lessons.ListBySection(ctx, sec.ID)
		out[i] = repository.SectionWithLessons{Section: sec, Lessons: lessons}
	}
	return &repository.CourseWithContent{Course: c, Sections: out}, nil
}

// fakeSectionRepo is an in-memory SectionRepository.
type fakeSectionRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.CourseSection
}

func newFakeSectionRepo() *fakeSectionRepo {
	return &fakeSectionRepo{byID: map[uuid.UUID]*entity.CourseSection{}}
}

func (f *fakeSectionRepo) Create(_ context.Context, s *entity.CourseSection) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt
	cp := *s
	f.byID[s.ID] = &cp
	return nil
}

func (f *fakeSectionRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.CourseSection, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if s, ok := f.byID[id]; ok && s.DeletedAt == nil {
		cp := *s
		return &cp, nil
	}
	return nil, apperror.NotFound("section not found", nil)
}

func (f *fakeSectionRepo) ListByCourse(_ context.Context, courseID uuid.UUID) ([]*entity.CourseSection, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.CourseSection
	for _, s := range f.byID {
		if s.DeletedAt == nil && s.CourseID == courseID {
			cp := *s
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderIndex < out[j].OrderIndex })
	return out, nil
}

func (f *fakeSectionRepo) Update(_ context.Context, s *entity.CourseSection) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if ex, ok := f.byID[s.ID]; !ok || ex.DeletedAt != nil {
		return apperror.NotFound("section not found", nil)
	}
	s.UpdatedAt = time.Now()
	cp := *s
	f.byID[s.ID] = &cp
	return nil
}

func (f *fakeSectionRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if s, ok := f.byID[id]; ok {
		now := time.Now()
		s.DeletedAt = &now
	}
	return nil
}

func (f *fakeSectionRepo) MaxOrderIndex(_ context.Context, courseID uuid.UUID) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	max := -1
	for _, s := range f.byID {
		if s.DeletedAt == nil && s.CourseID == courseID && s.OrderIndex > max {
			max = s.OrderIndex
		}
	}
	return max, nil
}

// fakeLessonRepo is an in-memory LessonRepository.
type fakeLessonRepo struct {
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.Lesson
}

func newFakeLessonRepo() *fakeLessonRepo {
	return &fakeLessonRepo{byID: map[uuid.UUID]*entity.Lesson{}}
}

func (f *fakeLessonRepo) Create(_ context.Context, l *entity.Lesson) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	l.CreatedAt = time.Now()
	l.UpdatedAt = l.CreatedAt
	cp := *l
	f.byID[l.ID] = &cp
	return nil
}

func (f *fakeLessonRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Lesson, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if l, ok := f.byID[id]; ok && l.DeletedAt == nil {
		cp := *l
		return &cp, nil
	}
	return nil, apperror.NotFound("lesson not found", nil)
}

func (f *fakeLessonRepo) ListBySection(_ context.Context, sectionID uuid.UUID) ([]*entity.Lesson, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Lesson
	for _, l := range f.byID {
		if l.DeletedAt == nil && l.SectionID == sectionID {
			cp := *l
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderIndex < out[j].OrderIndex })
	return out, nil
}

func (f *fakeLessonRepo) ListByCourse(_ context.Context, courseID uuid.UUID) ([]*entity.Lesson, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Lesson
	for _, l := range f.byID {
		if l.DeletedAt == nil && l.CourseID == courseID {
			cp := *l
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderIndex < out[j].OrderIndex })
	return out, nil
}

func (f *fakeLessonRepo) Update(_ context.Context, l *entity.Lesson) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if ex, ok := f.byID[l.ID]; !ok || ex.DeletedAt != nil {
		return apperror.NotFound("lesson not found", nil)
	}
	l.UpdatedAt = time.Now()
	cp := *l
	f.byID[l.ID] = &cp
	return nil
}

func (f *fakeLessonRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if l, ok := f.byID[id]; ok {
		now := time.Now()
		l.DeletedAt = &now
	}
	return nil
}

func (f *fakeLessonRepo) SoftDeleteBySection(_ context.Context, sectionID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	for _, l := range f.byID {
		if l.DeletedAt == nil && l.SectionID == sectionID {
			l.DeletedAt = &now
		}
	}
	return nil
}

func (f *fakeLessonRepo) MaxOrderIndex(_ context.Context, sectionID uuid.UUID) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	max := -1
	for _, l := range f.byID {
		if l.DeletedAt == nil && l.SectionID == sectionID && l.OrderIndex > max {
			max = l.OrderIndex
		}
	}
	return max, nil
}

// statsForCourse returns the count and total duration of a course's live
// lessons, mirroring the RecomputeCourseCounters SQL.
func (f *fakeLessonRepo) statsForCourse(courseID uuid.UUID) (count, dur int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, l := range f.byID {
		if l.DeletedAt == nil && l.CourseID == courseID {
			count++
			dur += l.DurationSeconds
		}
	}
	return count, dur
}

// fakeTxManager runs fn directly — the in-memory fakes give no real
// transactional isolation, but this lets the usecase's WithinTx grouping be
// exercised end-to-end.
type fakeTxManager struct{}

func (fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// fakeVideoRepo is a minimal in-memory VideoRepository. The course usecase only
// calls GetByLessonID (to find a lesson's video on delete) and SoftDelete; the
// remaining methods are unused and would panic via the embedded nil interface,
// surfacing an unexpected dependency.
type fakeVideoRepo struct {
	repository.VideoRepository
	mu   sync.Mutex
	byID map[uuid.UUID]*entity.Video
}

func newFakeVideoRepo() *fakeVideoRepo {
	return &fakeVideoRepo{byID: map[uuid.UUID]*entity.Video{}}
}

func (f *fakeVideoRepo) seed(v *entity.Video) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	cp := *v
	f.byID[v.ID] = &cp
}

func (f *fakeVideoRepo) GetByLessonID(_ context.Context, lessonID uuid.UUID) (*entity.Video, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, v := range f.byID {
		if v.DeletedAt == nil && v.LessonID == lessonID {
			cp := *v
			return &cp, nil
		}
	}
	return nil, apperror.NotFound("video not found", nil)
}

func (f *fakeVideoRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.byID[id]; ok {
		now := time.Now()
		v.DeletedAt = &now
	}
	return nil
}

// fakeStore is a minimal ObjectStore that records every Delete/DeletePrefix so a
// test can assert which objects were purged. It implements PrefixDeleter so the
// segment-sweep path is exercised.
type fakeStore struct {
	service.ObjectStore
	mu       sync.Mutex
	deleted  map[string]bool
	prefixes map[string]bool
}

func newFakeStore() *fakeStore {
	return &fakeStore{deleted: map[string]bool{}, prefixes: map[string]bool{}}
}

func (f *fakeStore) Delete(_ context.Context, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deleted[key] = true
	return nil
}

func (f *fakeStore) DeletePrefix(_ context.Context, prefix string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.prefixes[prefix] = true
	return nil
}

func (f *fakeStore) wasDeleted(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.deleted[key]
}

func (f *fakeStore) prefixDeleted(prefix string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.prefixes[prefix]
}

var _ service.PrefixDeleter = (*fakeStore)(nil)
