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
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/analytics"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/auth"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
	"github.com/Bereke1t2/Memere/memere-backend/internal/worker"
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

	// Build the fully-wired application (router + background workers).
	app := buildApp(cfg, dbPool, redisClient)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      app.Router,
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

	// Start the background attempt-expiry sweeper in its own goroutine. It shares
	// a cancelable context with shutdown so SIGINT/SIGTERM stops it cleanly.
	sweeperCtx, stopSweeper := context.WithCancel(context.Background())
	defer stopSweeper()
	if cfg.Sweeper.Enabled {
		go app.Sweeper.Run(sweeperCtx)
	} else {
		log.Println("attempt sweeper disabled (SWEEPER_ENABLED=false)")
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	stopSweeper()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")
}

// App bundles the constructed HTTP router and the background workers so main can
// own their lifecycle.
type App struct {
	Router  *gin.Engine
	Sweeper *worker.AttemptSweeper
}

// buildApp performs the explicit constructor wiring for the API: repositories
// over the pool/cache, the JWT manager, every usecase, the HTTP handlers, the
// router with its middleware, and the background sweeper. Plain dependency
// injection — no DI framework (matches the spec's explicit-wiring intent).
func buildApp(cfg *config.Config, pool *pgxpool.Pool, redisClient *redis.Client) *App {
	// Repositories.
	userRepo := postgres.NewUserRepo(pool)
	tokenRepo := postgres.NewRefreshTokenRepo(pool)
	courseRepo := postgres.NewCourseRepo(pool)
	sectionRepo := postgres.NewSectionRepo(pool)
	lessonRepo := postgres.NewLessonRepo(pool)
	quizRepo := postgres.NewQuizRepo(pool)
	questionRepo := postgres.NewQuestionRepo(pool)
	quizAttemptRepo := postgres.NewQuizAttemptRepo(pool)
	examRepo := postgres.NewExamRepo(pool)
	examAttemptRepo := postgres.NewExamAttemptRepo(pool)
	txManager := postgres.NewTxManager(pool)
	sessionRepo := redisrepo.NewSessionRepo(redisClient)
	attemptStateRepo := redisrepo.NewAttemptStateRepo(redisClient)
	scoreRankingRepo := redisrepo.NewScoreRankingRepo(redisClient)

	// JWT manager.
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, cfg.JWT.Issuer)

	// Usecases.
	authSvc := auth.NewService(userRepo, tokenRepo, sessionRepo, jwtManager)
	courseSvc := course.NewService(courseRepo, sectionRepo, lessonRepo, txManager)
	quizSvc := quiz.NewService(quizRepo, questionRepo, quizAttemptRepo, courseRepo, attemptStateRepo, txManager)
	examSvc := exam.NewService(examRepo, examAttemptRepo, courseRepo, attemptStateRepo, scoreRankingRepo, txManager)
	analyticsSvc := analytics.NewService(examRepo, examAttemptRepo, courseRepo, scoreRankingRepo)

	// Handlers + router.
	router := delivery_http.NewRouter(delivery_http.Deps{
		Config:    cfg,
		DB:        pool,
		Cache:     redisClient,
		JWT:       jwtManager,
		Auth:      delivery_http.NewAuthHandler(authSvc, userRepo),
		Courses:   delivery_http.NewCourseHandler(courseSvc),
		Quizzes:   delivery_http.NewQuizHandler(quizSvc),
		Exams:     delivery_http.NewExamHandler(examSvc),
		Analytics: delivery_http.NewAnalyticsHandler(analyticsSvc),
	})

	// Background workers: the expiry sweeper drives both engines off one ticker.
	sweeper := worker.NewAttemptSweeper(cfg.Sweeper.Interval, quizSvc, examSvc)

	return &App{Router: router, Sweeper: sweeper}
}
