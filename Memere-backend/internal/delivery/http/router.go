package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/Bereke1t2/Memere/memere-backend/internal/delivery/middleware"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
)

// Deps bundles everything NewRouter needs: the constructed handlers, the JWT
// manager (for per-route auth), the Redis client (for rate limiting), the
// HTTP config, and the DB/cache handles for the health probe.
type Deps struct {
	Config  *config.Config
	DB      *pgxpool.Pool
	Cache   *redis.Client
	JWT     *jwt.Manager
	Auth    *AuthHandler
	Courses *CourseHandler
}

// NewRouter assembles the Gin engine: the global middleware stack (in order),
// the health probe, and the versioned /api/v1 routes with their per-route
// auth/role middleware.
func NewRouter(deps Deps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global middleware — order matters (spec §3.2): request-id first so every
	// later log/recovery line can reference it, then recovery, logging, CORS,
	// and the global rate limit.
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(deps.Config.HTTP.CORSAllowedOrigins))
	r.Use(middleware.RateLimit(deps.Cache, deps.Config.HTTP.RateLimitRPM))

	r.GET("/health", healthHandler(deps.DB, deps.Cache))

	requireAuth := middleware.RequireAuth(deps.JWT)
	optionalAuth := middleware.OptionalAuth(deps.JWT)
	teacherOrAdmin := middleware.RequireRole(entity.RoleTeacher, entity.RoleAdmin)
	loginLimit := middleware.LoginRateLimit(deps.Cache, deps.Config.HTTP.LoginRateLimit, deps.Config.HTTP.LoginRateWindow)

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

	return r
}
