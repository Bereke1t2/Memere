package course

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// harness bundles a Service with its fakes for a test.
type harness struct {
	svc      *Service
	courses  *fakeCourseRepo
	sections *fakeSectionRepo
	lessons  *fakeLessonRepo
	videos   *fakeVideoRepo
	store    *fakeStore
}

func newHarness() *harness {
	sections := newFakeSectionRepo()
	lessons := newFakeLessonRepo()
	courses := newFakeCourseRepo(sections, lessons)
	videos := newFakeVideoRepo()
	store := newFakeStore()
	svc := NewService(courses, sections, lessons, videos, store, fakeTxManager{})
	return &harness{svc: svc, courses: courses, sections: sections, lessons: lessons, videos: videos, store: store}
}

func teacher() *Actor { return &Actor{UserID: uuid.New(), Role: entity.RoleTeacher} }
func student() *Actor { return &Actor{UserID: uuid.New(), Role: entity.RoleStudent} }
func admin() *Actor   { return &Actor{UserID: uuid.New(), Role: entity.RoleAdmin} }

func validCreateInput() CreateCourseInput {
	return CreateCourseInput{
		Title:       "Intro to Calculus",
		Description: "Limits, derivatives, integrals.",
		Subject:     "Mathematics",
		Grade:       12,
		Price:       100,
		Level:       entity.LevelBeginner,
	}
}

func TestCreateCourse_AsTeacherSetsTeacherID(t *testing.T) {
	h := newHarness()
	act := teacher()

	c, err := h.svc.CreateCourse(context.Background(), act, validCreateInput())
	if err != nil {
		t.Fatalf("CreateCourse: %v", err)
	}
	if c.TeacherID != act.UserID {
		t.Errorf("TeacherID = %v, want actor %v", c.TeacherID, act.UserID)
	}
	if c.Slug != "intro-to-calculus" {
		t.Errorf("Slug = %q, want slugified title", c.Slug)
	}
	if c.IsPublished {
		t.Error("new course should be unpublished")
	}
	if c.Currency != "ETB" || c.Language != "en" {
		t.Errorf("defaults not applied: currency=%q language=%q", c.Currency, c.Language)
	}
}

func TestCreateCourse_AsStudentForbidden(t *testing.T) {
	h := newHarness()
	_, err := h.svc.CreateCourse(context.Background(), student(), validCreateInput())
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("err = %v, want FORBIDDEN", err)
	}
}

func TestCreateCourse_DuplicateSlugGetsSuffix(t *testing.T) {
	h := newHarness()
	act := teacher()
	c1, err := h.svc.CreateCourse(context.Background(), act, validCreateInput())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	c2, err := h.svc.CreateCourse(context.Background(), act, validCreateInput())
	if err != nil {
		t.Fatalf("second create: %v", err)
	}
	if c1.Slug == c2.Slug {
		t.Errorf("expected distinct slugs, both %q", c1.Slug)
	}
}

func TestCreateCourse_ValidationError(t *testing.T) {
	h := newHarness()
	in := validCreateInput()
	in.Title = ""
	in.Grade = 99
	_, err := h.svc.CreateCourse(context.Background(), teacher(), in)
	if !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("err = %v, want VALIDATION_ERROR", err)
	}
}

func TestUpdateCourse_NotOwnerForbidden(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())

	other := teacher()
	newTitle := "Hijacked"
	_, err := h.svc.UpdateCourse(context.Background(), other, c.ID, UpdateCourseInput{Title: &newTitle})
	if !apperror.IsCode(err, "FORBIDDEN") {
		t.Fatalf("err = %v, want FORBIDDEN", err)
	}
}

func TestUpdateCourse_AdminCanEditAnyCourse(t *testing.T) {
	h := newHarness()
	c, _ := h.svc.CreateCourse(context.Background(), teacher(), validCreateInput())
	newTitle := "Admin Edited"
	updated, err := h.svc.UpdateCourse(context.Background(), admin(), c.ID, UpdateCourseInput{Title: &newTitle})
	if err != nil {
		t.Fatalf("admin update: %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("Title = %q, want %q", updated.Title, newTitle)
	}
}

func TestDeleteCourse_IsSoftDelete(t *testing.T) {
	h := newHarness()
	act := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), act, validCreateInput())

	if err := h.svc.DeleteCourse(context.Background(), act, c.ID); err != nil {
		t.Fatalf("DeleteCourse: %v", err)
	}
	// The row must still exist with DeletedAt set, not be physically removed.
	stored, ok := h.courses.byID[c.ID]
	if !ok {
		t.Fatal("course row was physically removed; soft-delete expected")
	}
	if stored.DeletedAt == nil {
		t.Error("DeletedAt not set after delete")
	}
	// And it must no longer be findable.
	if _, err := h.svc.GetCourse(context.Background(), act, c.ID); !apperror.IsNotFound(err) {
		t.Errorf("deleted course still readable: %v", err)
	}
}

func TestListCourses_AnonymousHidesUnpublished(t *testing.T) {
	h := newHarness()
	act := teacher()
	pub, _ := h.svc.CreateCourse(context.Background(), act, validCreateInput())
	if _, err := h.svc.PublishCourse(context.Background(), act, pub.ID); err != nil {
		t.Fatalf("publish: %v", err)
	}
	in := validCreateInput()
	in.Title = "Draft Course"
	if _, err := h.svc.CreateCourse(context.Background(), act, in); err != nil {
		t.Fatalf("create draft: %v", err)
	}

	// Anonymous (nil actor) must see only the published course.
	got, _, err := h.svc.ListCourses(context.Background(), nil, repository.CourseFilter{}, nil, 20)
	if err != nil {
		t.Fatalf("ListCourses: %v", err)
	}
	if len(got) != 1 || !got[0].IsPublished {
		t.Fatalf("anonymous listing = %d courses, want 1 published", len(got))
	}
}

func TestGetCourse_NonOwnerSeesAllPublishedCourseContent(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	if _, err := h.svc.PublishCourse(context.Background(), owner, c.ID); err != nil {
		t.Fatalf("publish: %v", err)
	}

	// A published course exposes ALL of its content to students, regardless of the
	// per-section/per-lesson publish flags (there is no student-facing gate on
	// those — only the course publish flag gates student visibility).
	secA, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "A", IsPublished: true})
	_, _ = h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "B", IsPublished: false})
	_, _ = h.svc.AddLesson(context.Background(), owner, secA.ID, LessonInput{Title: "L1", Type: entity.LessonTypeNote, IsPublished: true})
	_, _ = h.svc.AddLesson(context.Background(), owner, secA.ID, LessonInput{Title: "L2", Type: entity.LessonTypeNote, IsPublished: false})

	// A student (non-owner) sees both sections and both lessons of section A.
	view, err := h.svc.GetCourse(context.Background(), student(), c.ID)
	if err != nil {
		t.Fatalf("GetCourse: %v", err)
	}
	if len(view.Sections) != 2 {
		t.Fatalf("sections = %d, want 2 (published course shows all)", len(view.Sections))
	}
	totalLessons := 0
	for _, sec := range view.Sections {
		totalLessons += len(sec.Lessons)
	}
	if totalLessons != 2 {
		t.Fatalf("lessons = %d, want 2 (published course shows all)", totalLessons)
	}

	// The owner also sees everything.
	ownerView, err := h.svc.GetCourse(context.Background(), owner, c.ID)
	if err != nil {
		t.Fatalf("owner GetCourse: %v", err)
	}
	if len(ownerView.Sections) != 2 {
		t.Errorf("owner sections = %d, want 2", len(ownerView.Sections))
	}
}

func TestGetCourse_UnpublishedHiddenFromNonOwner(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput()) // unpublished

	if _, err := h.svc.GetCourse(context.Background(), student(), c.ID); !apperror.IsNotFound(err) {
		t.Errorf("student saw unpublished course: %v", err)
	}
	if _, err := h.svc.GetCourse(context.Background(), owner, c.ID); err != nil {
		t.Errorf("owner cannot see own unpublished course: %v", err)
	}
}

func TestAddLesson_BumpsCourseCounters(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	sec, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "S1", IsPublished: true})

	if _, err := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "A", Type: entity.LessonTypeVideo, DurationSeconds: 60, IsPublished: true}); err != nil {
		t.Fatalf("AddLesson A: %v", err)
	}
	if _, err := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "B", Type: entity.LessonTypeVideo, DurationSeconds: 90, IsPublished: true}); err != nil {
		t.Fatalf("AddLesson B: %v", err)
	}

	got, _ := h.courses.FindByID(context.Background(), c.ID)
	if got.TotalLessons != 2 {
		t.Errorf("TotalLessons = %d, want 2", got.TotalLessons)
	}
	if got.TotalDurationSeconds != 150 {
		t.Errorf("TotalDurationSeconds = %d, want 150", got.TotalDurationSeconds)
	}
}

func TestAddLesson_AutoOrderIndex(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	sec, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "S1", IsPublished: true})

	l0, _ := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "A", Type: entity.LessonTypeNote, IsPublished: true})
	l1, _ := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "B", Type: entity.LessonTypeNote, IsPublished: true})
	if l0.OrderIndex != 0 || l1.OrderIndex != 1 {
		t.Errorf("order indices = %d,%d, want 0,1", l0.OrderIndex, l1.OrderIndex)
	}
}

func TestDeleteSection_CascadesAndRecomputes(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	sec, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "S1", IsPublished: true})
	l, _ := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "A", Type: entity.LessonTypeVideo, DurationSeconds: 120, IsPublished: true})

	if err := h.svc.DeleteSection(context.Background(), owner, sec.ID); err != nil {
		t.Fatalf("DeleteSection: %v", err)
	}
	// Lesson cascade soft-deleted.
	if h.lessons.byID[l.ID].DeletedAt == nil {
		t.Error("lesson not cascade soft-deleted")
	}
	// Counters back to zero.
	got, _ := h.courses.FindByID(context.Background(), c.ID)
	if got.TotalLessons != 0 || got.TotalDurationSeconds != 0 {
		t.Errorf("counters = %d/%d, want 0/0 after section delete", got.TotalLessons, got.TotalDurationSeconds)
	}
}

func TestDeleteLesson_IsSoftAndRecomputes(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	sec, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "S1", IsPublished: true})
	l, _ := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "A", Type: entity.LessonTypeVideo, DurationSeconds: 30, IsPublished: true})

	if err := h.svc.DeleteLesson(context.Background(), owner, l.ID); err != nil {
		t.Fatalf("DeleteLesson: %v", err)
	}
	if h.lessons.byID[l.ID].DeletedAt == nil {
		t.Error("lesson physically removed; soft-delete expected")
	}
	got, _ := h.courses.FindByID(context.Background(), c.ID)
	if got.TotalLessons != 0 {
		t.Errorf("TotalLessons = %d, want 0 after delete", got.TotalLessons)
	}
}

// TestDeleteLesson_PurgesVideoStorageAndPdf verifies that deleting a lesson
// tombstones both the lesson and its video row (soft-delete) while permanently
// purging the video's storage artifacts (scalar keys + hls/<id>/ and
// thumbnails/<id>/ prefixes) and the notes PDF at its fixed key.
func TestDeleteLesson_PurgesVideoStorageAndPdf(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	sec, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "S1", IsPublished: true})
	l, _ := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "A", Type: entity.LessonTypeVideo, DurationSeconds: 30, IsPublished: true})

	vidID := uuid.New()
	orig := "originals/" + c.ID.String() + "/" + l.ID.String() + "/" + vidID.String() + "/source.mp4"
	master := "hls/" + vidID.String() + "/master.m3u8"
	thumb := "thumbnails/" + vidID.String() + "/thumb.jpg"
	h.videos.seed(&entity.Video{
		ID: vidID, LessonID: l.ID, CourseID: c.ID, Status: entity.VideoReady,
		OriginalFileKey: &orig, HLSMasterKey: &master, ThumbnailKey: &thumb,
	})

	if err := h.svc.DeleteLesson(context.Background(), owner, l.ID); err != nil {
		t.Fatalf("DeleteLesson: %v", err)
	}

	if h.lessons.byID[l.ID].DeletedAt == nil {
		t.Error("lesson row not soft-deleted")
	}
	if h.videos.byID[vidID].DeletedAt == nil {
		t.Error("video row not soft-deleted")
	}
	for _, k := range []string{orig, master, thumb} {
		if !h.store.wasDeleted(k) {
			t.Errorf("storage key %q was not deleted", k)
		}
	}
	if !h.store.prefixDeleted("hls/" + vidID.String() + "/") {
		t.Errorf("hls/ prefix for video %s was not purged", vidID)
	}
	if !h.store.prefixDeleted("thumbnails/" + vidID.String() + "/") {
		t.Errorf("thumbnails/ prefix for video %s was not purged", vidID)
	}
	if !h.store.wasDeleted("lessons/" + l.ID.String() + "/notes.pdf") {
		t.Error("notes PDF at the fixed key was not deleted")
	}
}

func TestAddLesson_InvalidTypeRejected(t *testing.T) {
	h := newHarness()
	owner := teacher()
	c, _ := h.svc.CreateCourse(context.Background(), owner, validCreateInput())
	sec, _ := h.svc.AddSection(context.Background(), owner, c.ID, SectionInput{Title: "S1", IsPublished: true})

	_, err := h.svc.AddLesson(context.Background(), owner, sec.ID, LessonInput{Title: "A", Type: entity.LessonType("podcast"), IsPublished: true})
	if !apperror.IsCode(err, "VALIDATION_ERROR") {
		t.Fatalf("err = %v, want VALIDATION_ERROR", err)
	}
}

func TestListCourses_PaginationCursor(t *testing.T) {
	h := newHarness()
	act := admin()
	// Create 3 courses; admin sees all regardless of publish state.
	for i := 0; i < 3; i++ {
		in := validCreateInput()
		in.Title = "Course " + string(rune('A'+i))
		if _, err := h.svc.CreateCourse(context.Background(), act, in); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}

	// Page size 2 → first page returns 2 with a nextCursor.
	page1, next, err := h.svc.ListCourses(context.Background(), act, repository.CourseFilter{}, nil, 2)
	if err != nil {
		t.Fatalf("page1: %v", err)
	}
	if len(page1) != 2 || next == nil {
		t.Fatalf("page1 len=%d next=%v, want 2 and a cursor", len(page1), next)
	}
	// Second page returns the remaining 1 with no further cursor.
	page2, next2, err := h.svc.ListCourses(context.Background(), act, repository.CourseFilter{}, next, 2)
	if err != nil {
		t.Fatalf("page2: %v", err)
	}
	if len(page2) != 1 {
		t.Fatalf("page2 len=%d, want 1", len(page2))
	}
	if next2 != nil {
		t.Errorf("nextCursor on last page = %v, want nil", next2)
	}
}
