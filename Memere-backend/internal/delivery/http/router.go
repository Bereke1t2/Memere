package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
)

// Deps bundles everything NewRouter needs: the constructed handlers, the JWT
// manager (for per-route auth), the Redis client (for rate limiting), the
// HTTP config, and the DB/cache handles for the health probe.
type Deps struct {
	Config    *config.Config
	DB        *pgxpool.Pool
	Cache     *redis.Client
	JWT       *jwt.Manager
	// Sessions is the Redis SessionRepository used for the JTI denylist check in
	// RequireAuth. Optional: when nil the denylist check is skipped.
	Sessions  repository.SessionRepository
	Auth      *AuthHandler
	Courses   *CourseHandler
	Quizzes   *QuizHandler
	Exams     *ExamHandler
	Analytics *AnalyticsHandler
	Video     *VideoHandler

	// Phase 4 (payments) handlers. Nil when payments are not configured, in which
	// case their routes are not registered.
	Payments      *PaymentHandler
	Enrollments   *EnrollmentHandler
	Subscriptions *SubscriptionHandler
	Revenue       *RevenueHandler

	// Phase 5 handlers. Non-nil when the respective feature is wired.
	Progress      *ProgressHandler
	Notifications *NotificationHandler
	Certificates  *CertificateHandler
	Admin         *AdminHandler
}

// NewRouter assembles the Gin engine: the global middleware stack (in order),
// the health probe, and the versioned /api/v1 routes with their per-route
// auth/role middleware.
func NewRouter(deps Deps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware — order matters (spec §3.2, Phase 6 §12):
	//  1. SecurityHeaders — hardening headers on every response
	//  2. RequestID       — correlation id available to every subsequent layer
	//  3. Recovery        — catches panics before they can abort tracing spans
	//  4. Tracing         — opens OTel span; uses request_id attribute
	//  5. Logger          — slog access log; enriches ctx with request-scoped logger
	//  6. Metrics         — records Prometheus RED metrics after the handler returns
	//  7. BodyLimit       — rejects oversized JSON bodies (1 MiB default)
	//  8. CORS            — sets Access-Control-* headers (fail-closed in prod)
	//  9. RateLimit       — Redis fixed-window; runs after logging so blocks appear
	// 10. Compress        — gzip response bodies > 1 KiB when client supports it
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery())
	r.Use(middleware.Tracing())
	r.Use(middleware.Logger())
	r.Use(middleware.Metrics())
	r.Use(middleware.BodyLimit(deps.Config.Security.BodyLimitBytes))
	r.Use(middleware.CORS(deps.Config.HTTP.CORSAllowedOrigins, deps.Config.App.IsProduction()))
	r.Use(middleware.RateLimit(deps.Cache, deps.Config.HTTP.RateLimitRPM))
	r.Use(middleware.Compress())

	r.GET("/healthz", healthHandler(deps.DB, deps.Cache))
	r.GET("/readyz", readyzHandler(deps.DB, deps.Cache))
	r.GET("/health", healthHandler(deps.DB, deps.Cache)) // legacy alias kept for backwards compat
	r.GET("/version", versionHandler())

	requireAuth := middleware.RequireAuth(deps.JWT, deps.Sessions)
	optionalAuth := middleware.OptionalAuth(deps.JWT)
	teacherOrAdmin := middleware.RequireRole(entity.RoleTeacher, entity.RoleAdmin)
	adminOnly := middleware.RequireRole(entity.RoleAdmin)
	loginLimit := middleware.LoginRateLimit(deps.Cache, deps.Config.HTTP.LoginRateLimit, deps.Config.HTTP.LoginRateWindow)
	userLimit := middleware.UserRateLimit(deps.Cache, 30)
	webhookLimit := middleware.ProviderWebhookRateLimit(deps.Cache, 120)

	v1 := r.Group("/api/v1")

	// Auth routes.
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", deps.Auth.Register)
		authGroup.POST("/login", loginLimit, deps.Auth.Login)
		authGroup.POST("/refresh", deps.Auth.Refresh)
		authGroup.POST("/logout", requireAuth, deps.Auth.Logout)
		authGroup.GET("/me", requireAuth, deps.Auth.Me)
	}

	// Course routes. Reads use OptionalAuth so visibility adapts to the viewer;
	// writes require auth and the teacher/admin role (ownership is enforced in
	// the usecase).
	courses := v1.Group("/courses")
	{
		courses.GET("", optionalAuth, deps.Courses.List)
		courses.GET("/:id", optionalAuth, deps.Courses.Get)
		courses.GET("/:id/sections", optionalAuth, deps.Courses.ListSections)

		courses.POST("", requireAuth, teacherOrAdmin, deps.Courses.Create)
		courses.PUT("/:id", requireAuth, teacherOrAdmin, deps.Courses.Update)
		courses.DELETE("/:id", requireAuth, teacherOrAdmin, deps.Courses.Delete)
		courses.POST("/:id/publish", requireAuth, teacherOrAdmin, deps.Courses.Publish)
		courses.POST("/:id/sections", requireAuth, teacherOrAdmin, deps.Courses.AddSection)
	}

	// Section-scoped lesson routes.
	sections := v1.Group("/sections")
	{
		sections.GET("/:id/lessons", optionalAuth, deps.Courses.ListLessons)
		sections.POST("/:id/lessons", requireAuth, teacherOrAdmin, deps.Courses.AddLesson)
	}

	// Lesson update/delete routes.
	lessons := v1.Group("/lessons")
	{
		lessons.PUT("/:id", requireAuth, teacherOrAdmin, deps.Courses.UpdateLesson)
		lessons.DELETE("/:id", requireAuth, teacherOrAdmin, deps.Courses.DeleteLesson)
	}

	// Quiz/exam authoring nested under a course (teacher/admin).
	courses.GET("/:id/quizzes", requireAuth, teacherOrAdmin, deps.Quizzes.ListByCourse)
	courses.POST("/:id/quizzes", requireAuth, teacherOrAdmin, deps.Quizzes.CreateQuiz)
	courses.GET("/:id/exams", requireAuth, teacherOrAdmin, deps.Exams.ListByCourse)
	courses.POST("/:id/exams", requireAuth, teacherOrAdmin, deps.Exams.CreateExam)

	// Quiz authoring + taking. Authoring is teacher/admin; taking requires auth,
	// with ownership enforced in the usecase (not middleware).
	quizzes := v1.Group("/quizzes")
	{
		quizzes.POST("/:id/questions", requireAuth, teacherOrAdmin, deps.Quizzes.AddQuestion)
		quizzes.PUT("/:id", requireAuth, teacherOrAdmin, deps.Quizzes.UpdateQuiz)
		quizzes.GET("/:id", requireAuth, deps.Quizzes.GetQuiz)
		quizzes.POST("/:id/attempts", requireAuth, deps.Quizzes.StartAttempt)
	}
	quizAttempts := v1.Group("/quiz-attempts")
	{
		quizAttempts.PATCH("/:id", requireAuth, deps.Quizzes.SaveProgress)
		quizAttempts.POST("/:id/submit", requireAuth, deps.Quizzes.Submit)
		quizAttempts.GET("/:id/result", requireAuth, deps.Quizzes.GetResult)
	}

	// Exam authoring + analytics stats (teacher/admin) and the mock-exam catalog.
	exams := v1.Group("/exams")
	{
		exams.POST("/:id/questions", requireAuth, teacherOrAdmin, deps.Exams.AddQuestion)
		exams.POST("/:id/publish", requireAuth, teacherOrAdmin, deps.Exams.Publish)
		exams.GET("/:id/stats", requireAuth, teacherOrAdmin, deps.Analytics.ExamStats)
		// Leaderboard: top-N + caller's own rank. Any authenticated student may view.
		exams.GET("/:id/leaderboard", requireAuth, deps.Analytics.Leaderboard)
	}
	mockExams := v1.Group("/mock-exams")
	{
		mockExams.GET("", optionalAuth, deps.Exams.ListMockExams)
		mockExams.POST("/:id/start", requireAuth, deps.Exams.Start)
	}
	examAttempts := v1.Group("/exam-attempts")
	{
		examAttempts.PATCH("/:id", requireAuth, deps.Exams.SaveProgress)
		examAttempts.POST("/:id/submit", requireAuth, deps.Exams.Submit)
		examAttempts.GET("/:id/results", requireAuth, deps.Exams.GetResult)
		examAttempts.GET("/:id/analytics", requireAuth, deps.Analytics.AttemptAnalytics)
	}

	// Student self-analytics.
	v1.GET("/me/trend", requireAuth, deps.Analytics.Trend)

	// Video pipeline (Phase 3). Reads (status/stream/download) require an
	// authenticated caller; access control (owner/admin or free/preview student)
	// is enforced in the usecase. Mutations are gated to teacher/admin by role,
	// with course ownership checked in the usecase.
	if deps.Video != nil {
		videos := v1.Group("/videos")
		{
			videos.GET("/:id/status", requireAuth, deps.Video.Status)
			videos.GET("/:id/stream", requireAuth, deps.Video.Stream)
			videos.GET("/:id/download-url", requireAuth, deps.Video.DownloadURL)
			videos.GET("/download/:token", requireAuth, deps.Video.ConsumeDownload)

			videos.POST("/:id/confirm", requireAuth, teacherOrAdmin, deps.Video.Confirm)
			videos.POST("/:id/retry", requireAuth, teacherOrAdmin, deps.Video.Retry)
		}
		v1.POST("/lessons/:id/videos/upload-url", requireAuth, teacherOrAdmin, deps.Video.RequestUpload)
	}

	// Phase 4 — payments, enrollments, subscriptions, revenue. Registered as a
	// block so an unconfigured payment stack leaves the routes absent rather than
	// nil-panicking.
	if deps.Payments != nil {
		// Provider webhook: PUBLIC and raw-body. It is registered OUTSIDE requireAuth
		// (unauthenticated providers must reach it) and does no JSON binding — the
		// handler reads the raw body itself for signature verification. The global
		// logger redacts bodies, so no payload/secret is logged.
		v1.POST("/webhooks/payments/:provider", webhookLimit, deps.Payments.Webhook)

		payments := v1.Group("/payments")
		{
			payments.POST("/initiate", requireAuth, userLimit, deps.Payments.Initiate)
			payments.GET("", requireAuth, deps.Payments.ListMine)
			payments.GET("/:id/status", requireAuth, deps.Payments.Status)
			payments.POST("/:id/refund", requireAuth, adminOnly, deps.Payments.Refund)
		}

		// Free-course enrollment (paid access goes through the payment flow).
		courses.POST("/:id/enroll-free", requireAuth, deps.Enrollments.EnrollFree)
		v1.GET("/me/enrollments", requireAuth, deps.Enrollments.ListMine)

		// Subscriptions. The plan catalogue is public; everything else is bearer.
		v1.GET("/subscription-plans", deps.Subscriptions.ListPlans)
		subs := v1.Group("/subscriptions")
		{
			subs.POST("", requireAuth, deps.Subscriptions.Subscribe)
			subs.POST("/:id/cancel", requireAuth, deps.Subscriptions.Cancel)
		}
		v1.GET("/me/subscription", requireAuth, deps.Subscriptions.GetMine)

		// Revenue reporting. Platform totals are admin-only; earnings are
		// teacher/admin; per-course sales ownership is checked in the usecase.
		v1.GET("/admin/revenue", requireAuth, adminOnly, deps.Revenue.PlatformRevenue)
		v1.GET("/me/earnings", requireAuth, teacherOrAdmin, deps.Revenue.MyEarnings)
		courses.GET("/:id/sales", requireAuth, deps.Revenue.CourseSales)
	}

	// Phase 5 — progress, notifications, certificates, admin.
	if deps.Progress != nil {
		v1.POST("/lessons/:id/complete", requireAuth, deps.Progress.Complete)
		v1.PUT("/lessons/:id/video-progress", requireAuth, deps.Progress.VideoProgress)
		v1.GET("/courses/:id/progress", requireAuth, deps.Progress.GetCourseProgress)
		v1.GET("/me/dashboard", requireAuth, deps.Progress.GetDashboard)
		v1.GET("/me/streak", requireAuth, deps.Progress.GetStreak)
	}

	if deps.Notifications != nil {
		v1.GET("/me/notifications", requireAuth, deps.Notifications.List)
		v1.GET("/me/notifications/unread-count", requireAuth, deps.Notifications.UnreadCount)
		v1.POST("/me/notifications/:id/read", requireAuth, deps.Notifications.MarkRead)
		v1.POST("/me/notifications/read-all", requireAuth, deps.Notifications.MarkAllRead)
		v1.POST("/me/devices", requireAuth, deps.Notifications.RegisterDevice)
		v1.DELETE("/me/devices/:token", requireAuth, deps.Notifications.UnregisterDevice)
		v1.GET("/me/notification-preferences", requireAuth, deps.Notifications.GetPreferences)
		v1.PUT("/me/notification-preferences", requireAuth, deps.Notifications.UpdatePreferences)
	}

	if deps.Certificates != nil {
		// Public serial verification — no auth required.
		v1.GET("/verify/certificates/:serial", deps.Certificates.Verify)
		v1.POST("/courses/:id/certificate", requireAuth, deps.Certificates.Issue)
		v1.GET("/me/certificates", requireAuth, deps.Certificates.ListMine)
		v1.GET("/certificates/:id/download", requireAuth, deps.Certificates.GetDownloadURL)
	}

	if deps.Admin != nil {
		adminGroup := v1.Group("/admin", requireAuth, adminOnly)
		{
			adminGroup.GET("/users", deps.Admin.ListUsers)
			adminGroup.GET("/users/:id", deps.Admin.GetUser)
			adminGroup.POST("/users/:id/suspend", deps.Admin.SuspendUser)
			adminGroup.POST("/users/:id/reactivate", deps.Admin.ReactivateUser)
			adminGroup.POST("/users/:id/role", deps.Admin.ChangeRole)

			adminGroup.GET("/courses", deps.Admin.ListCourses)
			adminGroup.POST("/courses/:id/unpublish", deps.Admin.UnpublishCourse)

			adminGroup.GET("/payments", deps.Admin.ListPayments)
			adminGroup.GET("/payments/:id", deps.Admin.GetPaymentDetail)
			adminGroup.POST("/payments/reconcile", deps.Admin.ReconcilePending)

			adminGroup.POST("/announcements", deps.Admin.Broadcast)

			adminGroup.GET("/analytics/overview", deps.Admin.Overview)
			adminGroup.GET("/analytics/revenue", deps.Admin.RevenueBreakdown)
			adminGroup.GET("/analytics/engagement", deps.Admin.GetEngagementStats)
		}
	}

	return r
}
