package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	delivery_http "github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/cache"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/database"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres"
	redisrepo "github.com/Bereke1t2/Memere/memere-backend/internal/repository/redis"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/auth"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	// Connect to PostgreSQL
	dbPool, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Connect to Redis
	redisClient, err := cache.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to cache: %v", err)
	}
	defer redisClient.Close()

	// Build the fully-wired HTTP router.
	router := buildRouter(cfg, dbPool, redisClient)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen and serve error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")
}

// buildRouter performs the explicit constructor wiring for the API: repositories
// over the pool/cache, the JWT manager, the auth and course usecases, the HTTP
// handlers, and finally the router with its middleware. Plain dependency
// injection — no DI framework (matches the spec's explicit-wiring intent).
func buildRouter(cfg *config.Config, pool *pgxpool.Pool, redisClient *redis.Client) *gin.Engine {
	// Repositories.
	userRepo := postgres.NewUserRepo(pool)
	tokenRepo := postgres.NewRefreshTokenRepo(pool)
	courseRepo := postgres.NewCourseRepo(pool)
	sectionRepo := postgres.NewSectionRepo(pool)
	lessonRepo := postgres.NewLessonRepo(pool)
	txManager := postgres.NewTxManager(pool)
	sessionRepo := redisrepo.NewSessionRepo(redisClient)

	// JWT manager.
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, cfg.JWT.Issuer)

	// Usecases.
	authSvc := auth.NewService(userRepo, tokenRepo, sessionRepo, jwtManager)
	courseSvc := course.NewService(courseRepo, sectionRepo, lessonRepo, txManager)

	// Handlers + router.
	return delivery_http.NewRouter(delivery_http.Deps{
		Config:  cfg,
		DB:      pool,
		Cache:   redisClient,
		JWT:     jwtManager,
		Auth:    delivery_http.NewAuthHandler(authSvc, userRepo),
		Courses: delivery_http.NewCourseHandler(courseSvc),
	})
}
