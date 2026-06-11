package http

import (
	"testing"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
)

// TestNewRouter_RegistersPhase2Routes verifies the router assembles without a
// panic (no conflicting static/param routes) and that the Phase 2 quiz, exam, and
// analytics routes are registered. Handlers are nil here: route registration only
// captures method values, so nil receivers are never dereferenced.
func TestNewRouter_RegistersPhase2Routes(t *testing.T) {
	deps := Deps{
		Config:    &config.Config{},
		JWT:       jwt.NewManager("test-secret", time.Minute, time.Hour, "memere-test"),
		Auth:      &AuthHandler{},
		Courses:   &CourseHandler{},
		Quizzes:   &QuizHandler{},
		Exams:     &ExamHandler{},
		Analytics: &AnalyticsHandler{},
	}

	r := NewRouter(deps)

	want := map[string]bool{
		"POST /api/v1/courses/:id/quizzes":        false,
		"POST /api/v1/quizzes/:id/questions":      false,
		"PUT /api/v1/quizzes/:id":                 false,
		"GET /api/v1/quizzes/:id":                 false,
		"POST /api/v1/quizzes/:id/attempts":       false,
		"PATCH /api/v1/quiz-attempts/:id":         false,
		"POST /api/v1/quiz-attempts/:id/submit":   false,
		"GET /api/v1/quiz-attempts/:id/result":    false,
		"POST /api/v1/courses/:id/exams":          false,
		"POST /api/v1/exams/:id/questions":        false,
		"POST /api/v1/exams/:id/publish":          false,
		"GET /api/v1/exams/:id/stats":             false,
		"GET /api/v1/mock-exams":                  false,
		"POST /api/v1/mock-exams/:id/start":       false,
		"PATCH /api/v1/exam-attempts/:id":         false,
		"POST /api/v1/exam-attempts/:id/submit":   false,
		"GET /api/v1/exam-attempts/:id/results":   false,
		"GET /api/v1/exam-attempts/:id/analytics": false,
		"GET /api/v1/me/trend":                    false,
	}

	for _, route := range r.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}

	for key, found := range want {
		if !found {
			t.Errorf("route not registered: %s", key)
		}
	}
}

// TestNewRouter_RegistersPhase3VideoRoutes verifies the video pipeline routes
// assemble without a panic — in particular that the static segment "download"
// and the ":id" param coexist at the same /videos level — and are registered.
func TestNewRouter_RegistersPhase3VideoRoutes(t *testing.T) {
	deps := Deps{
		Config:    &config.Config{},
		JWT:       jwt.NewManager("test-secret", time.Minute, time.Hour, "memere-test"),
		Auth:      &AuthHandler{},
		Courses:   &CourseHandler{},
		Quizzes:   &QuizHandler{},
		Exams:     &ExamHandler{},
		Analytics: &AnalyticsHandler{},
		Video:     &VideoHandler{},
	}

	r := NewRouter(deps)

	want := map[string]bool{
		"POST /api/v1/lessons/:id/videos/upload-url": false,
		"POST /api/v1/videos/:id/confirm":            false,
		"GET /api/v1/videos/:id/status":              false,
		"POST /api/v1/videos/:id/retry":              false,
		"GET /api/v1/videos/:id/stream":              false,
		"GET /api/v1/videos/:id/download-url":        false,
		"GET /api/v1/videos/download/:token":         false,
	}
	for _, route := range r.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for key, found := range want {
		if !found {
			t.Errorf("route not registered: %s", key)
		}
	}
}
