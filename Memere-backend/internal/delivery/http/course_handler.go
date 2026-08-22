package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http/dto"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// CourseHandler adapts the course usecase to HTTP.
type CourseHandler struct {
	svc   *course.Service
	store service.ObjectStore
	// publicURL is the API's externally reachable base URL (config APP_PUBLIC_URL).
	// It builds the absolute, backend-served thumbnail URL persisted on the course
	// so a Drive/S3-backed image is fetched via GET /courses/:id/thumbnail.
	publicURL string
}

// NewCourseHandler builds a CourseHandler. publicURL is the API base URL used to
// construct thumbnail URLs; store is optional (nil disables file endpoints).
func NewCourseHandler(svc *course.Service, publicURL string, store ...service.ObjectStore) *CourseHandler {
	var s service.ObjectStore
	if len(store) > 0 {
		s = store[0]
	}
	return &CourseHandler{svc: svc, store: s, publicURL: publicURL}
}

// actor pulls the authenticated caller from context (nil for anonymous).
func actor(c *gin.Context) *course.Actor {
	a, _ := middleware.ActorFromContext(c)
	return a
}

// List handles GET /courses → paginated, visibility-filtered course list.
func (h *CourseHandler) List(c *gin.Context) {
	cursor, err := pagination.Decode(c.Query("after"))
	if err != nil {
		respondError(c, apperror.BadRequest("invalid cursor", err))
		return
	}
	limit := pagination.NormalizeLimit(atoiDefault(c.Query("limit"), 0))

	filter := repository.CourseFilter{}
	if s := c.Query("subject"); s != "" {
		filter.Subject = &s
	}
	if g := c.Query("grade"); g != "" {
		grade, convErr := strconv.Atoi(g)
		if convErr != nil {
			respondError(c, apperror.BadRequest("invalid grade", convErr))
			return
		}
		filter.Grade = &grade
	}

	courses, next, err := h.svc.ListCourses(c.Request.Context(), actor(c), filter, cursor, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	nextCursor := ""
	if next != nil {
		nextCursor = next.Encode()
	}
	respondJSON(c, http.StatusOK, dto.Paginated{
		Data:       dto.NewCourseListResponse(courses),
		NextCursor: nextCursor,
		Limit:      limit,
	})
}

// Get handles GET /courses/:id → nested course detail.
func (h *CourseHandler) Get(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	content, err := h.svc.GetCourse(c.Request.Context(), actor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, dto.NewCourseDetailResponse(content))
}

// Create handles POST /courses → 201.
func (h *CourseHandler) Create(c *gin.Context) {
	var req dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	created, err := h.svc.CreateCourse(c.Request.Context(), actor(c), course.CreateCourseInput{
		Title:            req.Title,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Subject:          req.Subject,
		Grade:            req.Grade,
		ThumbnailURL:     req.ThumbnailURL,
		Price:            req.Price,
		Currency:         req.Currency,
		IsFree:           req.IsFree,
		Language:         req.Language,
		Level:            entity.Level(req.Level),
		Metadata:         req.Metadata,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewCourseResponse(created)
	respondJSON(c, http.StatusCreated, &resp)
}

// Update handles PUT /courses/:id → 200.
func (h *CourseHandler) Update(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	in := course.UpdateCourseInput{
		Title:            req.Title,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Subject:          req.Subject,
		Grade:            req.Grade,
		ThumbnailURL:     req.ThumbnailURL,
		Price:            req.Price,
		IsFree:           req.IsFree,
		Language:         req.Language,
		Metadata:         req.Metadata,
	}
	if req.Level != nil {
		lvl := entity.Level(*req.Level)
		in.Level = &lvl
	}
	updated, err := h.svc.UpdateCourse(c.Request.Context(), actor(c), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewCourseResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// Delete handles DELETE /courses/:id → 204 (soft delete).
func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.DeleteCourse(c.Request.Context(), actor(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Publish handles POST /courses/:id/publish → 200.
func (h *CourseHandler) Publish(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	updated, err := h.svc.PublishCourse(c.Request.Context(), actor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewCourseResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// AddSection handles POST /courses/:id/sections → 201.
func (h *CourseHandler) AddSection(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	sec, err := h.svc.AddSection(c.Request.Context(), actor(c), courseID, course.SectionInput{
		Title:       req.Title,
		Description: req.Description,
		IsPublished: req.IsPublished,
		OrderIndex:  req.OrderIndex,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewSectionResponse(sec)
	respondJSON(c, http.StatusCreated, &resp)
}

// ListSections handles GET /courses/:id/sections → 200.
func (h *CourseHandler) ListSections(c *gin.Context) {
	courseID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	sections, err := h.svc.ListSections(c.Request.Context(), actor(c), courseID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"data": dto.NewSectionListResponse(sections)})
}

// AddLesson handles POST /sections/:id/lessons → 201.
func (h *CourseHandler) AddLesson(c *gin.Context) {
	sectionID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	lesson, err := h.svc.AddLesson(c.Request.Context(), actor(c), sectionID, course.LessonInput{
		Title:           req.Title,
		Type:            entity.LessonType(req.Type),
		IsFreePreview:   req.IsFreePreview,
		DurationSeconds: req.DurationSeconds,
		IsPublished:     req.IsPublished,
		OrderIndex:      req.OrderIndex,
		Content:         req.Content,
		PdfURL:          req.PdfURL,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewLessonResponse(lesson)
	respondJSON(c, http.StatusCreated, &resp)
}

// UpdateLesson handles PUT /lessons/:id → 200.
func (h *CourseHandler) UpdateLesson(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}
	updated, err := h.svc.UpdateLesson(c.Request.Context(), actor(c), id, course.LessonInput{
		Title:           req.Title,
		Type:            entity.LessonType(req.Type),
		IsFreePreview:   req.IsFreePreview,
		DurationSeconds: req.DurationSeconds,
		IsPublished:     req.IsPublished,
		OrderIndex:      req.OrderIndex,
		Content:         req.Content,
		PdfURL:          req.PdfURL,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	resp := dto.NewLessonResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// DeleteLesson handles DELETE /lessons/:id → 200.
func (h *CourseHandler) DeleteLesson(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.svc.DeleteLesson(c.Request.Context(), actor(c), id); err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"message": "lesson deleted"})
}

// ListLessons handles GET /sections/:id/lessons → 200.
func (h *CourseHandler) ListLessons(c *gin.Context) {
	sectionID, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	lessons, err := h.svc.ListLessons(c.Request.Context(), actor(c), sectionID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"data": dto.NewLessonListResponse(lessons)})
}

// parseUUIDParam parses a path parameter as a UUID, returning a clean 400
// INVALID_ID on failure.
func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, apperror.New(http.StatusBadRequest, "INVALID_ID", "invalid id parameter", err)
	}
	return id, nil
}

// atoiDefault parses s as an int, returning def when s is empty or invalid.
func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// UploadLessonPDF handles POST /lessons/:id/pdf (teacher/admin).
func (h *CourseHandler) UploadLessonPDF(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, apperror.BadRequest("missing 'file' form field", err))
		return
	}
	defer file.Close()

	if header.Size > 50*1024*1024 {
		respondError(c, apperror.BadRequest("file size exceeds maximum 50MB limit", nil))
		return
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		respondError(c, apperror.Internal(err))
		return
	}

	rawBytes := buf.Bytes()
	limit := len(rawBytes)
	if limit > 1024 {
		limit = 1024
	}
	if len(rawBytes) < 5 || !bytes.Contains(rawBytes[:limit], []byte("%PDF-")) {
		respondError(c, apperror.BadRequest("uploaded file is not a valid PDF document", nil))
		return
	}

	key := fmt.Sprintf("lessons/%s/notes.pdf", id.String())
	if h.store != nil {
		if err := h.store.Put(c.Request.Context(), key, "application/pdf", bytes.NewReader(rawBytes)); err != nil {
			respondError(c, apperror.Internal(err))
			return
		}
	}

	updated, err := h.svc.SetLessonPdfURL(c.Request.Context(), actor(c), id, key)
	if err != nil {
		respondError(c, err)
		return
	}

	resp := dto.NewLessonResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// DownloadLessonPDF handles GET /lessons/:id/pdf.
func (h *CourseHandler) DownloadLessonPDF(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	l, err := h.svc.GetLessonByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	if l.PdfURL == nil || strings.TrimSpace(*l.PdfURL) == "" {
		respondError(c, apperror.NotFound("no PDF attached to this lesson", nil))
		return
	}

	pdfPath := strings.TrimSpace(*l.PdfURL)

	// If it's already an absolute URL
	if strings.HasPrefix(pdfPath, "http://") || strings.HasPrefix(pdfPath, "https://") {
		c.Redirect(http.StatusFound, pdfPath)
		return
	}

	if h.store != nil {
		// Stream the PDF bytes with an explicit Content-Length header via c.Data.
		// io.Copy uses chunked transfer encoding (no Content-Length), causing mobile Dio client to hang at ~100%.
		rc, err := h.store.Get(c.Request.Context(), pdfPath)
		if err == nil {
			defer rc.Close()
			data, rerr := io.ReadAll(rc)
			if rerr == nil && len(data) > 0 {
				c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filepath.Base(pdfPath)))
				c.Data(http.StatusOK, "application/pdf", data)
				return
			}
		}

		// Fallback: hand back a presigned URL (works in prod where the CDN/S3
		// host is publicly reachable).
		if presigned, perr := h.store.PresignGet(c.Request.Context(), pdfPath, 1*time.Hour); perr == nil && presigned != "" {
			c.Redirect(http.StatusFound, presigned)
			return
		}
	}

	respondError(c, apperror.NotFound("PDF document not found in storage", nil))
}

// UploadCourseThumbnail handles POST /courses/:id/thumbnail (teacher/admin). The
// image is sent as multipart/form-data under "file". Ownership is enforced by
// UpdateCourse (called first, so an object is only written for a course the actor
// owns); the bytes are then stored under a fixed, extension-less key so re-uploads
// overwrite cleanly and the public serve endpoint needs no key bookkeeping. Works
// with the S3/MinIO and Google Drive backends alike. The global body limit is
// waived for this route (see middleware.BodyLimit); this handler caps at 5 MiB.
func (h *CourseHandler) UploadCourseThumbnail(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if h.store == nil {
		respondError(c, apperror.New(http.StatusServiceUnavailable, "STORAGE_UNAVAILABLE", "object storage is not configured", nil))
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, apperror.BadRequest("missing 'file' form field", err))
		return
	}
	defer file.Close()

	const maxThumb = 5 * 1024 * 1024
	if header.Size > maxThumb {
		respondError(c, apperror.BadRequest("thumbnail exceeds the maximum 5MB limit", nil))
		return
	}
	// Read with a +1 guard so a lying Content-Length can't slip past the cap.
	data, err := io.ReadAll(io.LimitReader(file, maxThumb+1))
	if err != nil {
		respondError(c, apperror.Internal(err))
		return
	}
	if len(data) > maxThumb {
		respondError(c, apperror.BadRequest("thumbnail exceeds the maximum 5MB limit", nil))
		return
	}
	// Sniff the real type from the bytes (never trust the client's declared type).
	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "image/") {
		respondError(c, apperror.BadRequest("uploaded file is not an image", nil))
		return
	}

	// Authorize + persist the URL FIRST: UpdateCourse asserts owner/admin, so the
	// object write below only happens for a course the caller may edit (prevents
	// overwriting another teacher's thumbnail object at the shared key).
	thumbURL := strings.TrimRight(h.publicURL, "/") + "/api/v1/courses/" + id.String() + "/thumbnail"
	updated, err := h.svc.UpdateCourse(c.Request.Context(), actor(c), id, course.UpdateCourseInput{ThumbnailURL: &thumbURL})
	if err != nil {
		respondError(c, err)
		return
	}

	key := fmt.Sprintf("courses/%s/thumbnail", id.String())
	if err := h.store.Put(c.Request.Context(), key, contentType, bytes.NewReader(data)); err != nil {
		respondError(c, apperror.Internal(err))
		return
	}

	resp := dto.NewCourseResponse(updated)
	respondJSON(c, http.StatusOK, &resp)
}

// ServeCourseThumbnail handles GET /courses/:id/thumbnail — PUBLIC. A course
// thumbnail is marketing content, so it is served cacheable and unsigned. The
// bytes are streamed through the backend from the object store (so the Google
// Drive backend is never exposed publicly); the type is sniffed from the bytes.
func (h *CourseHandler) ServeCourseThumbnail(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}
	if h.store == nil {
		respondError(c, apperror.NotFound("thumbnail not found", nil))
		return
	}
	key := fmt.Sprintf("courses/%s/thumbnail", id.String())
	rc, err := h.store.Get(c.Request.Context(), key)
	if err != nil {
		// Missing object (or an upstream fault) — a placeholder client shows its own
		// fallback, so report a clean 404 rather than leaking backend detail.
		respondError(c, apperror.NotFound("thumbnail not found", nil))
		return
	}
	defer rc.Close()

	const maxThumb = 5 * 1024 * 1024
	data, err := io.ReadAll(io.LimitReader(rc, maxThumb+1))
	if err != nil {
		respondError(c, apperror.Internal(err))
		return
	}
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, http.DetectContentType(data), data)
}

