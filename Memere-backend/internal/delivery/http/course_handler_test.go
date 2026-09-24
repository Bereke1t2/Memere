package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ---- Fakes for testing CourseHandler ----------------------------------------

type mockCourseRepo struct {
	course  *entity.Course
	content *repository.CourseWithContent
}

func (m *mockCourseRepo) Create(context.Context, *entity.Course) error { return nil }
func (m *mockCourseRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Course, error) {
	if m.course != nil && m.course.ID == id {
		return m.course, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}
func (m *mockCourseRepo) FindBySlug(context.Context, string) (*entity.Course, error) {
	return nil, apperror.NotFound("course not found", nil)
}
func (m *mockCourseRepo) List(context.Context, repository.CourseFilter, *pagination.Cursor, int) ([]*entity.Course, *pagination.Cursor, error) {
	return nil, nil, nil
}
func (m *mockCourseRepo) Update(context.Context, *entity.Course) error       { return nil }
func (m *mockCourseRepo) SoftDelete(context.Context, uuid.UUID) error        { return nil }
func (m *mockCourseRepo) RecomputeCounters(context.Context, uuid.UUID) error { return nil }
func (m *mockCourseRepo) GetCourseWithSectionsAndLessons(_ context.Context, id uuid.UUID) (*repository.CourseWithContent, error) {
	if m.content != nil && m.content.Course.ID == id {
		return m.content, nil
	}
	return nil, apperror.NotFound("course not found", nil)
}

type mockLessonRepo struct {
	lessons map[uuid.UUID]*entity.Lesson
}

func (m *mockLessonRepo) Create(context.Context, *entity.Lesson) error { return nil }
func (m *mockLessonRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Lesson, error) {
	if l, ok := m.lessons[id]; ok {
		return l, nil
	}
	return nil, apperror.NotFound("lesson not found", nil)
}
func (m *mockLessonRepo) ListBySection(context.Context, uuid.UUID) ([]*entity.Lesson, error) {
	var list []*entity.Lesson
	for _, l := range m.lessons {
		list = append(list, l)
	}
	return list, nil
}
func (m *mockLessonRepo) ListByCourse(context.Context, uuid.UUID) ([]*entity.Lesson, error) {
	var list []*entity.Lesson
	for _, l := range m.lessons {
		list = append(list, l)
	}
	return list, nil
}
func (m *mockLessonRepo) Update(context.Context, *entity.Lesson) error       { return nil }
func (m *mockLessonRepo) SoftDelete(context.Context, uuid.UUID) error          { return nil }
func (m *mockLessonRepo) SoftDeleteBySection(context.Context, uuid.UUID) error { return nil }
func (m *mockLessonRepo) MaxOrderIndex(context.Context, uuid.UUID) (int, error)   { return 0, nil }

type mockSectionRepo struct{}

func (m *mockSectionRepo) Create(context.Context, *entity.CourseSection) error { return nil }
func (m *mockSectionRepo) FindByID(context.Context, uuid.UUID) (*entity.CourseSection, error) {
	return nil, nil
}
func (m *mockSectionRepo) ListByCourse(context.Context, uuid.UUID) ([]*entity.CourseSection, error) {
	return nil, nil
}
func (m *mockSectionRepo) Update(context.Context, *entity.CourseSection) error { return nil }
func (m *mockSectionRepo) SoftDelete(context.Context, uuid.UUID) error         { return nil }
func (m *mockSectionRepo) MaxOrderIndex(context.Context, uuid.UUID) (int, error) {
	return 0, nil
}

type mockEnrollRepo struct {
	enrolled map[[2]uuid.UUID]bool
}

func (m *mockEnrollRepo) Create(context.Context, *entity.Enrollment) error { return nil }
func (m *mockEnrollRepo) Exists(_ context.Context, studentID, courseID uuid.UUID) (bool, error) {
	return m.enrolled[[2]uuid.UUID{studentID, courseID}], nil
}
func (m *mockEnrollRepo) GetActiveForStudent(_ context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error) {
	if m.enrolled[[2]uuid.UUID{studentID, courseID}] {
		return &entity.Enrollment{
			ID:        uuid.New(),
			StudentID: studentID,
			CourseID:  courseID,
			Source:    entity.SourcePurchase,
		}, nil
	}
	return nil, apperror.NotFound("enrollment not found", nil)
}
func (m *mockEnrollRepo) ListByStudent(context.Context, uuid.UUID, int) ([]*entity.Enrollment, error) {
	return nil, nil
}
func (m *mockEnrollRepo) ListAllByStudent(context.Context, uuid.UUID) ([]*entity.Enrollment, error) {
	return nil, nil
}
func (m *mockEnrollRepo) Delete(context.Context, uuid.UUID, uuid.UUID) error { return nil }

type mockSubRepo struct{}

func (m *mockSubRepo) Create(context.Context, *entity.Subscription) error { return nil }
func (m *mockSubRepo) GetByID(context.Context, uuid.UUID) (*entity.Subscription, error) {
	return nil, nil
}
func (m *mockSubRepo) GetActiveForStudent(context.Context, uuid.UUID) (*entity.Subscription, error) {
	return nil, apperror.NotFound("no subscription", nil)
}
func (m *mockSubRepo) UpdateStatus(context.Context, uuid.UUID, entity.SubscriptionStatus) error {
	return nil
}
func (m *mockSubRepo) ExtendPeriod(context.Context, uuid.UUID, time.Time) error { return nil }
func (m *mockSubRepo) CancelAtPeriodEnd(context.Context, uuid.UUID) (bool, error) {
	return false, nil
}
func (m *mockSubRepo) ListExpiring(context.Context, int) ([]*entity.Subscription, error) {
	return nil, nil
}
func (m *mockSubRepo) ExpireLapsed(context.Context, uuid.UUID) (bool, error) {
	return false, nil
}

type mockStore struct {
	files map[string][]byte
}

func (m *mockStore) Put(_ context.Context, key, _ string, r io.Reader) error {
	data, _ := io.ReadAll(r)
	m.files[key] = data
	return nil
}
func (m *mockStore) Get(_ context.Context, key string) (io.ReadCloser, error) {
	if data, ok := m.files[key]; ok {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return nil, apperror.NotFound("not found", nil)
}
func (m *mockStore) Exists(_ context.Context, key string) (bool, error) {
	_, ok := m.files[key]
	return ok, nil
}
func (m *mockStore) Delete(_ context.Context, key string) error {
	delete(m.files, key)
	return nil
}
func (m *mockStore) PresignGet(context.Context, string, time.Duration) (string, error) {
	return "", nil
}
func (m *mockStore) PresignPut(context.Context, string, string, time.Duration) (string, error) {
	return "", nil
}

type mockTxManager struct{}

func (m *mockTxManager) WithinTx(_ context.Context, fn func(context.Context) error) error {
	return fn(context.Background())
}

// ---- Tests ------------------------------------------------------------------

func TestCourseHandler_Get_ContentRedactionForGuests(t *testing.T) {
	teacherID := uuid.New()
	courseID := uuid.New()
	sectionID := uuid.New()
	previewLessonID := uuid.New()
	paidLessonID := uuid.New()

	previewContent := "Free preview notes"
	previewPdf := "lessons/preview/notes.pdf"
	paidContent := "Secret paid notes"
	paidPdf := "lessons/paid/notes.pdf"

	courseEntity := &entity.Course{
		ID:          courseID,
		TeacherID:   teacherID,
		Title:       "Paid Chemistry Course",
		Slug:        "paid-chemistry",
		Price:       500,
		IsFree:      false,
		IsPublished: true,
	}

	previewLesson := &entity.Lesson{
		ID:            previewLessonID,
		SectionID:     sectionID,
		CourseID:      courseID,
		Title:         "Introduction (Free Preview)",
		Type:          entity.LessonTypeNote,
		IsFreePreview: true,
		IsPublished:   true,
		Content:       &previewContent,
		PdfURL:        &previewPdf,
	}

	paidLesson := &entity.Lesson{
		ID:            paidLessonID,
		SectionID:     sectionID,
		CourseID:      courseID,
		Title:         "Advanced Thermodynamics (Paid)",
		Type:          entity.LessonTypeNote,
		IsFreePreview: false,
		IsPublished:   true,
		Content:       &paidContent,
		PdfURL:        &paidPdf,
	}

	courseWithContent := &repository.CourseWithContent{
		Course: courseEntity,
		Sections: []repository.SectionWithLessons{
			{
				Section: &entity.CourseSection{
					ID:          sectionID,
					CourseID:    courseID,
					Title:       "Section 1",
					IsPublished: true,
				},
				Lessons: []*entity.Lesson{previewLesson, paidLesson},
			},
		},
	}

	cRepo := &mockCourseRepo{course: courseEntity, content: courseWithContent}
	lRepo := &mockLessonRepo{lessons: map[uuid.UUID]*entity.Lesson{previewLessonID: previewLesson, paidLessonID: paidLesson}}
	sRepo := &mockSectionRepo{}
	eRepo := &mockEnrollRepo{enrolled: map[[2]uuid.UUID]bool{}}
	subRepo := &mockSubRepo{}
	txMgr := &mockTxManager{}
	store := &mockStore{files: map[string][]byte{
		previewPdf: []byte("%PDF-preview"),
		paidPdf:    []byte("%PDF-paid"),
	}}

	accessSvc := access.NewService(eRepo, subRepo, cRepo, nil)
	courseSvc := course.NewService(cRepo, sRepo, lRepo, nil, store, txMgr)
	handler := NewCourseHandler(courseSvc, accessSvc, "http://localhost:8080", store)

	r := gin.New()
	r.GET("/api/v1/courses/:id", handler.Get)
	r.GET("/api/v1/lessons/:id/pdf", handler.DownloadLessonPDF)

	// 1. Unauthenticated Guest queries paid course detail
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/courses/"+courseID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp dto.CourseDetailResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(resp.Sections) != 1 || len(resp.Sections[0].Lessons) != 2 {
		t.Fatalf("expected 1 section with 2 lessons, got %d sections", len(resp.Sections))
	}

	lessons := resp.Sections[0].Lessons
	// Lesson 0 is free preview -> content & pdf_url MUST be present
	if lessons[0].Content == nil || *lessons[0].Content != previewContent {
		t.Errorf("preview lesson content should be present, got: %v", lessons[0].Content)
	}
	if lessons[0].PdfURL == nil || *lessons[0].PdfURL != previewPdf {
		t.Errorf("preview lesson pdf_url should be present, got: %v", lessons[0].PdfURL)
	}

	// Lesson 1 is paid non-preview -> content & pdf_url MUST be redacted (nil)
	if lessons[1].Content != nil {
		t.Errorf("paid non-preview lesson content MUST be redacted (nil), got: %v", *lessons[1].Content)
	}
	if lessons[1].PdfURL != nil {
		t.Errorf("paid non-preview lesson pdf_url MUST be redacted (nil), got: %v", *lessons[1].PdfURL)
	}

	// 2. Unauthenticated Guest attempts to download paid lesson PDF -> MUST return 403 NOT_ENROLLED
	pdfReq, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/"+paidLessonID.String()+"/pdf", nil)
	pdfW := httptest.NewRecorder()
	r.ServeHTTP(pdfW, pdfReq)

	if pdfW.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for paid lesson PDF, got %d: %s", pdfW.Code, pdfW.Body.String())
	}

	// 3. Unauthenticated Guest downloads free preview PDF -> MUST return 200 OK
	prevReq, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/"+previewLessonID.String()+"/pdf", nil)
	prevW := httptest.NewRecorder()
	r.ServeHTTP(prevW, prevReq)

	if prevW.Code != http.StatusOK {
		t.Errorf("expected 200 OK for free preview lesson PDF, got %d: %s", prevW.Code, prevW.Body.String())
	}
}

func TestCourseHandler_Get_EnrolledStudentHasFullAccess(t *testing.T) {
	studentID := uuid.New()
	courseID := uuid.New()
	sectionID := uuid.New()
	paidLessonID := uuid.New()

	paidContent := "Secret paid notes"
	paidPdf := "lessons/paid/notes.pdf"

	courseEntity := &entity.Course{
		ID:          courseID,
		TeacherID:   uuid.New(),
		Title:       "Paid Chemistry Course",
		Slug:        "paid-chemistry",
		Price:       500,
		IsFree:      false,
		IsPublished: true,
	}

	paidLesson := &entity.Lesson{
		ID:            paidLessonID,
		SectionID:     sectionID,
		CourseID:      courseID,
		Title:         "Advanced Thermodynamics (Paid)",
		Type:          entity.LessonTypeNote,
		IsFreePreview: false,
		IsPublished:   true,
		Content:       &paidContent,
		PdfURL:        &paidPdf,
	}

	courseWithContent := &repository.CourseWithContent{
		Course: courseEntity,
		Sections: []repository.SectionWithLessons{
			{
				Section: &entity.CourseSection{
					ID:          sectionID,
					CourseID:    courseID,
					Title:       "Section 1",
					IsPublished: true,
				},
				Lessons: []*entity.Lesson{paidLesson},
			},
		},
	}

	cRepo := &mockCourseRepo{course: courseEntity, content: courseWithContent}
	lRepo := &mockLessonRepo{lessons: map[uuid.UUID]*entity.Lesson{paidLessonID: paidLesson}}
	sRepo := &mockSectionRepo{}
	eRepo := &mockEnrollRepo{enrolled: map[[2]uuid.UUID]bool{
		{studentID, courseID}: true,
	}}
	subRepo := &mockSubRepo{}
	txMgr := &mockTxManager{}
	store := &mockStore{files: map[string][]byte{
		paidPdf: []byte("%PDF-paid-content"),
	}}

	accessSvc := access.NewService(eRepo, subRepo, cRepo, nil)
	courseSvc := course.NewService(cRepo, sRepo, lRepo, nil, store, txMgr)
	handler := NewCourseHandler(courseSvc, accessSvc, "http://localhost:8080", store)

	r := gin.New()
	// Middleware simulating authenticated student
	r.Use(func(c *gin.Context) {
		c.Set("actor", &course.Actor{
			UserID: studentID,
			Role:   entity.RoleStudent,
		})
		c.Next()
	})
	r.GET("/api/v1/courses/:id", handler.Get)
	r.GET("/api/v1/lessons/:id/pdf", handler.DownloadLessonPDF)

	// Enrolled student views course detail -> MUST have full content & pdf_url
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/courses/"+courseID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp dto.CourseDetailResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	lesson := resp.Sections[0].Lessons[0]
	if lesson.Content == nil || *lesson.Content != paidContent {
		t.Errorf("enrolled student must receive paid lesson content, got: %v", lesson.Content)
	}
	if lesson.PdfURL == nil || *lesson.PdfURL != paidPdf {
		t.Errorf("enrolled student must receive paid lesson pdf_url, got: %v", lesson.PdfURL)
	}

	// Enrolled student downloads PDF -> MUST return 200 OK
	pdfReq, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/"+paidLessonID.String()+"/pdf", nil)
	pdfW := httptest.NewRecorder()
	r.ServeHTTP(pdfW, pdfReq)

	if pdfW.Code != http.StatusOK {
		t.Errorf("expected 200 OK for enrolled student downloading PDF, got %d: %s", pdfW.Code, pdfW.Body.String())
	}
}

func TestCourseHandler_FreeCourse_GuestAccess(t *testing.T) {
	courseID := uuid.New()
	sectionID := uuid.New()
	lessonID := uuid.New()

	lessonContent := "Free course lesson content"
	lessonPdf := "lessons/free/notes.pdf"

	courseEntity := &entity.Course{
		ID:          courseID,
		TeacherID:   uuid.New(),
		Title:       "Free Biology Course",
		Slug:        "free-biology",
		Price:       0,
		IsFree:      true,
		IsPublished: true,
	}

	lesson := &entity.Lesson{
		ID:            lessonID,
		SectionID:     sectionID,
		CourseID:      courseID,
		Title:         "Cell Structure",
		Type:          entity.LessonTypeNote,
		IsFreePreview: false,
		IsPublished:   true,
		Content:       &lessonContent,
		PdfURL:        &lessonPdf,
	}

	courseWithContent := &repository.CourseWithContent{
		Course: courseEntity,
		Sections: []repository.SectionWithLessons{
			{
				Section: &entity.CourseSection{
					ID:          sectionID,
					CourseID:    courseID,
					Title:       "Section 1",
					IsPublished: true,
				},
				Lessons: []*entity.Lesson{lesson},
			},
		},
	}

	cRepo := &mockCourseRepo{course: courseEntity, content: courseWithContent}
	lRepo := &mockLessonRepo{lessons: map[uuid.UUID]*entity.Lesson{lessonID: lesson}}
	sRepo := &mockSectionRepo{}
	eRepo := &mockEnrollRepo{enrolled: map[[2]uuid.UUID]bool{}}
	subRepo := &mockSubRepo{}
	txMgr := &mockTxManager{}
	store := &mockStore{files: map[string][]byte{
		lessonPdf: []byte("%PDF-free-course-content"),
	}}

	accessSvc := access.NewService(eRepo, subRepo, cRepo, nil)
	courseSvc := course.NewService(cRepo, sRepo, lRepo, nil, store, txMgr)
	handler := NewCourseHandler(courseSvc, accessSvc, "http://localhost:8080", store)

	r := gin.New()
	r.GET("/api/v1/courses/:id", handler.Get)
	r.GET("/api/v1/lessons/:id/pdf", handler.DownloadLessonPDF)

	// Guest queries free course detail -> MUST have full content & pdf_url
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/courses/"+courseID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp dto.CourseDetailResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	resLesson := resp.Sections[0].Lessons[0]
	if resLesson.Content == nil || *resLesson.Content != lessonContent {
		t.Errorf("guest must receive free course lesson content, got: %v", resLesson.Content)
	}
	if resLesson.PdfURL == nil || *resLesson.PdfURL != lessonPdf {
		t.Errorf("guest must receive free course lesson pdf_url, got: %v", resLesson.PdfURL)
	}

	// Guest downloads PDF for free course -> MUST return 200 OK
	pdfReq, _ := http.NewRequest(http.MethodGet, "/api/v1/lessons/"+lessonID.String()+"/pdf", nil)
	pdfW := httptest.NewRecorder()
	r.ServeHTTP(pdfW, pdfReq)

	if pdfW.Code != http.StatusOK {
		t.Errorf("expected 200 OK for guest downloading free course PDF, got %d: %s", pdfW.Code, pdfW.Body.String())
	}
}
