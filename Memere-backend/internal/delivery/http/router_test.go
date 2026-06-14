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

// TestNewRouter_RegistersPhase4PaymentRoutes verifies the Phase 4 payment,
// enrollment, subscription, and revenue routes assemble without a panic (the new
// /courses/:id/{enroll-free,sales} sub-routes coexist with the existing :id
// routes) and are registered. In particular the public, unauthenticated webhook
// route must be present.
func TestNewRouter_RegistersPhase4PaymentRoutes(t *testing.T) {
	deps := Deps{
		Config:        &config.Config{},
		JWT:           jwt.NewManager("test-secret", time.Minute, time.Hour, "memere-test"),
		Auth:          &AuthHandler{},
		Courses:       &CourseHandler{},
		Quizzes:       &QuizHandler{},
		Exams:         &ExamHandler{},
		Analytics:     &AnalyticsHandler{},
		Payments:      &PaymentHandler{},
		Enrollments:   &EnrollmentHandler{},
		Subscriptions: &SubscriptionHandler{},
		Revenue:       &RevenueHandler{},
	}

	r := NewRouter(deps)

	want := map[string]bool{
		"POST /api/v1/payments/initiate":           false,
		"GET /api/v1/payments":                     false,
		"GET /api/v1/payments/:id/status":          false,
		"POST /api/v1/payments/:id/refund":         false,
		"POST /api/v1/webhooks/payments/:provider": false,
		"POST /api/v1/courses/:id/enroll-free":     false,
		"GET /api/v1/me/enrollments":               false,
		"GET /api/v1/subscription-plans":           false,
		"POST /api/v1/subscriptions":               false,
		"POST /api/v1/subscriptions/:id/cancel":    false,
		"GET /api/v1/me/subscription":              false,
		"GET /api/v1/admin/revenue":                false,
		"GET /api/v1/me/earnings":                  false,
		"GET /api/v1/courses/:id/sales":            false,
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
