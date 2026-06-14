package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	delivery_http "github.com/Bereke1t2/Memere/memere-backend/internal/delivery/http"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/cache"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/database"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/messaging"
	paymentinfra "github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/payment"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/storage"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/transcode"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres"
	redisrepo "github.com/Bereke1t2/Memere/memere-backend/internal/repository/redis"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/analytics"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/auth"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/coupon"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/enrollment"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/payment"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/quiz"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/revenue"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/subscription"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/video"
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
	app, err := buildApp(ctx, cfg, dbPool, redisClient)
	if err != nil {
		log.Fatalf("Failed to build application: %v", err)
	}

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

	// Background workers share a cancelable context with shutdown so SIGINT/
	// SIGTERM stops them cleanly.
	bgCtx, stopBackground := context.WithCancel(context.Background())
	defer stopBackground()

	if cfg.Sweeper.Enabled {
		go app.Sweeper.Run(bgCtx)
	} else {
		log.Println("attempt sweeper disabled (SWEEPER_ENABLED=false)")
	}

	if cfg.SubSweeper.Enabled {
		go app.SubSweeper.Run(bgCtx)
	} else {
		log.Println("subscription sweeper disabled (SUBSCRIPTION_SWEEP_ENABLED=false)")
	}

	// Transcode worker: re-enqueue videos stranded in 'processing' by a previous
	// crash (the in-proc queue is non-durable), then start consuming jobs.
	if n, err := app.VideoUC.RequeueStuck(bgCtx); err != nil {
		log.Printf("requeue stuck videos failed: %v", err)
	} else if n > 0 {
		log.Printf("re-enqueued %d stuck video(s) for transcoding", n)
	}
	go app.Worker.Run(bgCtx)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	stopBackground()

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
	Router     *gin.Engine
	Sweeper    *worker.AttemptSweeper
	SubSweeper *worker.SubscriptionSweeper
	Worker     *worker.TranscodeWorker
	VideoUC    *video.Service
}

// buildApp performs the explicit constructor wiring for the API: repositories
// over the pool/cache, the JWT manager, every usecase, the HTTP handlers, the
// router with its middleware, and the background workers. Plain dependency
// injection — no DI framework (matches the spec's explicit-wiring intent).
func buildApp(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, redisClient *redis.Client) (*App, error) {
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
	videoRepo := postgres.NewVideoRepo(pool)
	enrollmentRepo := postgres.NewEnrollmentRepo(pool)
	subscriptionRepo := postgres.NewSubscriptionRepo(pool)
	paymentRepo := postgres.NewPaymentRepo(pool)
	couponRepo := postgres.NewCouponRepo(pool)
	webhookRepo := postgres.NewWebhookEventRepo(pool)
	revenueRepo := postgres.NewRevenueRepo(pool)
	txManager := postgres.NewTxManager(pool)
	sessionRepo := redisrepo.NewSessionRepo(redisClient)
	attemptStateRepo := redisrepo.NewAttemptStateRepo(redisClient)
	scoreRankingRepo := redisrepo.NewScoreRankingRepo(redisClient)
	downloadTokens := redisrepo.NewDownloadTokenStore(redisClient)

	// JWT manager.
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, cfg.JWT.Issuer)

	// Object storage + CDN signer for the video pipeline. The store talks to
	// S3/MinIO; the signer mints CloudFront-signed URLs in production and falls
	// back to S3/MinIO pre-signed GET in dev.
	store, err := storage.NewS3Store(ctx, storage.Config{
		Endpoint:        cfg.Storage.Endpoint,
		Region:          cfg.Storage.Region,
		Bucket:          cfg.Storage.Bucket,
		AccessKeyID:     cfg.Storage.AccessKeyID,
		SecretAccessKey: cfg.Storage.SecretAccessKey,
		UsePathStyle:    cfg.Storage.UsePathStyle,
	})
	if err != nil {
		return nil, fmt.Errorf("object store: %w", err)
	}
	signer, err := storage.NewCDNSigner(storage.CDNConfig{
		Domain:        cfg.Storage.CDNDomain,
		KeyPairID:     cfg.Storage.CDNKeyPairID,
		PrivateKeyPEM: cfg.Storage.CDNPrivateKeyPEM,
	}, store, nil)
	if err != nil {
		return nil, fmt.Errorf("cdn signer: %w", err)
	}

	// In-process transcode queue (non-durable; future swap to SQS).
	queue := messaging.NewInProcQueue(cfg.Video.QueueBuffer)

	// Usecases.
	authSvc := auth.NewService(userRepo, tokenRepo, sessionRepo, jwtManager)
	courseSvc := course.NewService(courseRepo, sectionRepo, lessonRepo, txManager)
	// access.Service is the single authority on "can this caller reach this
	// course/lesson?" — quiz/exam taking and paid video all route through it.
	accessSvc := access.NewService(enrollmentRepo, subscriptionRepo, courseRepo, nil)
	quizSvc := quiz.NewService(quizRepo, questionRepo, quizAttemptRepo, courseRepo, attemptStateRepo, txManager, accessSvc)
	examSvc := exam.NewService(examRepo, examAttemptRepo, courseRepo, attemptStateRepo, scoreRankingRepo, txManager, accessSvc)
	analyticsSvc := analytics.NewService(examRepo, examAttemptRepo, courseRepo, scoreRankingRepo)
	videoSvc := video.NewService(videoRepo, lessonRepo, courseRepo, store, queue, signer, downloadTokens, accessSvc, video.Config{
		UploadURLTTL:   cfg.Storage.UploadURLTTL,
		MaxUploadBytes: cfg.Video.MaxUploadBytes,
		MaxAttempts:    cfg.Video.MaxAttempts,
		StreamURLTTL:   cfg.Storage.StreamURLTTL,
		DownloadURLTTL: cfg.Storage.DownloadURLTTL,
	})

	// Phase 4 payment stack. The provider registry is built from config; the
	// test-only mock provider is appended only when explicitly enabled
	// (PAYMENT_MOCK_ENABLED) so it can never settle a real payment in production.
	providers := []service.PaymentProvider{
		paymentinfra.NewChapaProvider(paymentinfra.ChapaConfig{
			SecretKey:     cfg.Chapa.SecretKey,
			WebhookSecret: cfg.Chapa.WebhookSecret,
			BaseURL:       cfg.Chapa.BaseURL,
		}),
		paymentinfra.NewTelebirrProvider(paymentinfra.TelebirrConfig{
			WebhookSecret: cfg.Telebirr.WebhookSecret,
		}),
		paymentinfra.NewStripeProvider(paymentinfra.StripeConfig{
			WebhookSecret: cfg.Stripe.WebhookSecret,
		}),
	}
	if cfg.Payment.MockEnabled {
		log.Println("payment: mock provider ENABLED (PAYMENT_MOCK_ENABLED=true) — do not use in production")
		providers = append(providers, paymentinfra.NewMockProvider(paymentinfra.MockConfig{
			WebhookSecret: cfg.Payment.MockWebhookSecret,
		}))
	}
	registry := paymentinfra.NewRegistry(providers...)

	// Subscription plans are config-driven (spec §1.6) — never hardcoded here.
	subCfg := subscription.Config{
		Monthly: subscription.PlanSpec{
			Plan:     entity.PlanMonthly,
			Price:    decimal.NewFromFloat(cfg.Payment.MonthlyPrice),
			Currency: cfg.Payment.MonthlyCurrency,
			Period:   cfg.Payment.MonthlyPeriod,
		},
		Annual: subscription.PlanSpec{
			Plan:     entity.PlanAnnual,
			Price:    decimal.NewFromFloat(cfg.Payment.AnnualPrice),
			Currency: cfg.Payment.AnnualCurrency,
			Period:   cfg.Payment.AnnualPeriod,
		},
	}

	couponSvc := coupon.NewService(couponRepo, nil)
	subscriptionSvc := subscription.NewService(subscriptionRepo, subCfg, nil)
	enrollmentSvc := enrollment.NewService(enrollmentRepo, courseRepo, nil)
	revenueSvc := revenue.NewService(revenueRepo, courseRepo, revenue.Config{
		TeacherShare: decimal.NewFromFloat(cfg.Payment.TeacherShare),
	})
	// The subscription service satisfies both the PlanPricer and the
	// SubscriptionActivator ports the payment flow needs.
	paymentSvc := payment.NewService(
		paymentRepo, enrollmentRepo, couponRepo, webhookRepo, courseRepo,
		registry, couponSvc, subscriptionSvc, subscriptionSvc,
		txManager, nil,
		payment.Config{
			CallbackURL:     cfg.Payment.CallbackURL,
			ReturnURL:       cfg.Payment.ReturnURL,
			DefaultCurrency: cfg.Payment.DefaultCurrency,
		},
		nil,
	)

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
		Video:     delivery_http.NewVideoHandler(videoSvc),

		Payments:      delivery_http.NewPaymentHandler(paymentSvc),
		Enrollments:   delivery_http.NewEnrollmentHandler(enrollmentSvc),
		Subscriptions: delivery_http.NewSubscriptionHandler(subscriptionSvc, paymentSvc),
		Revenue:       delivery_http.NewRevenueHandler(revenueSvc, nil),
	})

	// Background workers: the expiry sweeper drives both attempt engines off one
	// ticker; the transcode worker consumes the in-proc queue and runs FFmpeg.
	sweeper := worker.NewAttemptSweeper(cfg.Sweeper.Interval, quizSvc, examSvc)
	subSweeper := worker.NewSubscriptionSweeper(cfg.SubSweeper.Interval, subscriptionSvc, nil)
	coder := transcode.NewFFmpeg()
	transcodeWorker := worker.NewTranscodeWorker(queue.Transcode(), store, videoRepo, coder, queue, nil, worker.WorkerCfg{
		Concurrency: cfg.Video.Concurrency,
		WorkDir:     cfg.Video.WorkDir,
		MaxAttempts: cfg.Video.MaxAttempts,
	})

	return &App{
		Router:     router,
		Sweeper:    sweeper,
		SubSweeper: subSweeper,
		Worker:     transcodeWorker,
		VideoUC:    videoSvc,
	}, nil
}
