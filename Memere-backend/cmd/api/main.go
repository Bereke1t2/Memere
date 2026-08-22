package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	redis "github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	_ "go.uber.org/automaxprocs" // auto-sets GOMAXPROCS to cgroup CPU quota in containers

	"github.com/Bereke1t2/Memere/memere-backend/config"
	infmetrics "github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/metrics"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/tracing"
	repocache "github.com/Bereke1t2/Memere/memere-backend/internal/repository/cache"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/logger"
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
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/notification"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/pdf"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/access"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/admin"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/analytics"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/auth"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/certificate"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/coupon"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/enrollment"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/exam"
	notificationuc "github.com/Bereke1t2/Memere/memere-backend/internal/usecase/notification"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/payment"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/progress"
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

	// Phase 6 — structured logger. Must come first so all subsequent log calls
	// use the configured level/format and the redacting handler.
	appLogger := logger.New(cfg.App.Env, cfg.Observability.LogLevel)
	_ = appLogger // slog.SetDefault is called inside logger.New

	// Phase 6 Skill 2 — fail fast on weak/missing secrets in production.
	if err := cfg.ValidateProduction(); err != nil {
		slog.Error("production config validation failed", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Phase 6 — OpenTelemetry tracing.
	shutdownTracing, err := tracing.Setup(ctx, "memere-api", cfg.Observability.OTELEndpoint)
	if err != nil {
		slog.Error("failed to init tracing", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdownTracing(context.Background()); err != nil {
			slog.Error("tracing shutdown error", "err", err)
		}
	}()

	// Phase 6 — Prometheus metrics server (separate port so /metrics is never
	// exposed on the public API port). An empty MetricsPort disables it.
	if cfg.Observability.MetricsPort != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		metricsSrv := &http.Server{
			Addr:         ":" + cfg.Observability.MetricsPort,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}
		go func() {
			slog.Info("metrics server listening", "port", cfg.Observability.MetricsPort)
			if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("metrics server error", "err", err)
			}
		}()
	}

	// Connect to PostgreSQL
	dbPool, err := database.Connect(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Connect to Redis
	redisClient, err := cache.Connect(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to cache", "err", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Phase 6 — scrape pgx pool stats into Prometheus gauges every 15 s.
	go func() {
		t := time.NewTicker(15 * time.Second)
		defer t.Stop()
		for range t.C {
			s := dbPool.Stat()
			infmetrics.DBPoolAcquired.Set(float64(s.AcquiredConns()))
			infmetrics.DBPoolTotal.Set(float64(s.TotalConns()))
		}
	}()

	// Build the fully-wired application (router + background workers).
	app, err := buildApp(ctx, cfg, dbPool, redisClient)
	if err != nil {
		slog.Error("failed to build application", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		slog.Info("API server listening", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve error", "err", err)
			os.Exit(1)
		}
	}()

	// Background workers share a cancelable context with shutdown so SIGINT/
	// SIGTERM stops them cleanly.
	bgCtx, stopBackground := context.WithCancel(context.Background())
	defer stopBackground()

	if cfg.Sweeper.Enabled {
		go app.Sweeper.Run(bgCtx)
	} else {
		slog.Info("attempt sweeper disabled", "reason", "SWEEPER_ENABLED=false")
	}

	if cfg.SubSweeper.Enabled {
		go app.SubSweeper.Run(bgCtx)
	} else {
		slog.Info("subscription sweeper disabled", "reason", "SUBSCRIPTION_SWEEP_ENABLED=false")
	}

	// Transcode worker: re-enqueue videos stranded in 'processing' by a previous
	// crash (the in-proc queue is non-durable), then start consuming jobs.
	if n, err := app.VideoUC.RequeueStuck(bgCtx); err != nil {
		slog.Error("requeue stuck videos failed", "err", err)
	} else if n > 0 {
		slog.Info("re-enqueued stuck videos", "count", n)
	}
	go app.Worker.Run(bgCtx)

	// Phase 5 workers: notification fan-out + engagement (streak) sweeper.
	go app.NotifWorker.Run(bgCtx)

	if cfg.EngagementSweep.Enabled {
		go app.EngSweeper.Run(bgCtx)
	} else {
		slog.Info("engagement sweeper disabled", "reason", "ENGAGEMENT_SWEEP_ENABLED=false")
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	stopBackground()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		slog.Error("server forced to shutdown", "err", err)
		os.Exit(1)
	}

	slog.Info("server exited gracefully")
}

// App bundles the constructed HTTP router and the background workers so main can
// own their lifecycle.
type App struct {
	Router        *gin.Engine
	Sweeper       *worker.AttemptSweeper
	SubSweeper    *worker.SubscriptionSweeper
	Worker        *worker.TranscodeWorker
	NotifWorker   *worker.NotificationWorker
	EngSweeper    *worker.EngagementSweeper
	VideoUC       *video.Service
}

// buildApp performs the explicit constructor wiring for the API: repositories
// over the pool/cache, the JWT manager, every usecase, the HTTP handlers, the
// router with its middleware, and the background workers. Plain dependency
// injection — no DI framework (matches the spec's explicit-wiring intent).
func buildApp(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, redisClient *redis.Client) (*App, error) {
	// Repositories (Phase 6 Skill 3: hot reads wrapped with Redis cache-aside).
	userRepo := repocache.NewCachedUserRepo(postgres.NewUserRepo(pool), redisClient)
	tokenRepo := postgres.NewRefreshTokenRepo(pool)
	courseRepo := repocache.NewCachedCourseRepo(postgres.NewCourseRepo(pool), redisClient)
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
	progressRepo := postgres.NewProgressRepo(pool)
	notifRepo := postgres.NewNotificationRepo(pool)
	deviceRepo := postgres.NewDeviceTokenRepo(pool)
	prefRepo := postgres.NewPreferenceRepo(pool)
	auditRepo := postgres.NewAdminAuditRepo(pool)
	certRepo := postgres.NewCertificateRepo(pool)
	sessionRepo := redisrepo.NewSessionRepo(redisClient)
	attemptStateRepo := redisrepo.NewAttemptStateRepo(redisClient)
	scoreRankingRepo := redisrepo.NewScoreRankingRepo(redisClient)
	downloadTokens := redisrepo.NewDownloadTokenStore(redisClient)

	// JWT manager.
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, cfg.JWT.Issuer)

	// Object storage backend. STORAGE_PROVIDER selects S3/MinIO (default), the
	// S3-compatible Backblaze B2 (b2), or the single admin-owned Google Drive
	// account (gdrive). The CDN signer sits in front of whichever store is chosen;
	// in dev (no CDN_DOMAIN) it delegates to the store's PresignGet — which for
	// Drive returns a short-lived HMAC-signed /media proxy URL — so the video
	// delivery path is reused unchanged.
	var store service.ObjectStore
	var driveStore *storage.GoogleDriveStore  // non-nil only in gdrive mode
	var driveCreds *storage.PgCredentialStore // shared by the store and OAuth flow
	if cfg.Storage.Provider == "gdrive" {
		var err error
		driveCreds, err = storage.NewPgCredentialStore(pool, cfg.JWT.Secret)
		if err != nil {
			return nil, fmt.Errorf("gdrive credential store: %w", err)
		}
		driveStore, err = storage.NewGoogleDriveStore(storage.GoogleDriveConfig{
			ClientID:     cfg.GoogleDrive.ClientID,
			ClientSecret: cfg.GoogleDrive.ClientSecret,
			RootFolderID: cfg.GoogleDrive.RootFolderID,
			RefreshToken: cfg.GoogleDrive.RefreshToken,
			MediaBaseURL: strings.TrimRight(cfg.App.PublicURL, "/") + "/api/v1/media",
			SignSecret:   cfg.JWT.Secret,
		}, storage.NewPgObjectIndex(pool), driveCreds)
		if err != nil {
			return nil, fmt.Errorf("gdrive store: %w", err)
		}
		store = driveStore
	} else {
		// S3-compatible backends. AWS S3 and MinIO use the generic AWS_*/S3_* keys;
		// Backblaze B2 (STORAGE_PROVIDER=b2) speaks the same S3 API but takes its
		// own endpoint + application-key credentials, virtual-hosted addressing, and
		// needs the SDK's default upload checksums turned off (B2 rejects them).
		scfg := storage.Config{
			Endpoint:        cfg.Storage.Endpoint,
			Region:          cfg.Storage.Region,
			Bucket:          cfg.Storage.Bucket,
			AccessKeyID:     cfg.Storage.AccessKeyID,
			SecretAccessKey: cfg.Storage.SecretAccessKey,
			UsePathStyle:    cfg.Storage.UsePathStyle,
		}
		if cfg.Storage.Provider == "b2" {
			scfg = storage.Config{
				Endpoint:         cfg.Storage.B2EndpointURL(),
				Region:           cfg.Storage.B2Region,
				Bucket:           cfg.Storage.B2Bucket,
				AccessKeyID:      cfg.Storage.B2KeyID,
				SecretAccessKey:  cfg.Storage.B2AppKey,
				UsePathStyle:     false, // B2 uses virtual-hosted-style addressing
				DisableChecksums: true,  // B2 rejects the SDK's default CRC32 upload trailers
			}
		}
		s3store, err := storage.NewS3Store(ctx, scfg)
		if err != nil {
			return nil, fmt.Errorf("object store: %w", err)
		}
		store = s3store
	}
	// Log the selected backend at boot (no secrets) so it is obvious where
	// uploaded objects actually land — the usual cause of "I switched providers
	// but nothing shows up" is the process still running with the old provider.
	switch cfg.Storage.Provider {
	case "gdrive":
		slog.Info("object storage configured", "provider", "gdrive", "root_folder_id", cfg.GoogleDrive.RootFolderID)
	case "b2":
		slog.Info("object storage configured", "provider", "b2",
			"endpoint", cfg.Storage.B2EndpointURL(), "region", cfg.Storage.B2Region, "bucket", cfg.Storage.B2Bucket)
	default:
		slog.Info("object storage configured", "provider", cfg.Storage.Provider,
			"endpoint", cfg.Storage.Endpoint, "region", cfg.Storage.Region,
			"bucket", cfg.Storage.Bucket, "path_style", cfg.Storage.UsePathStyle)
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

	// Phase 5 — notification senders (real if keys present, LogSender in dev).
	// Created before usecases so hooks can be passed to all services from the start.
	pushSender := notification.NewFCMSender(cfg.FCM.ServerKey)
	emailSender := notification.NewSendGridSender(cfg.SendGrid.APIKey, cfg.SendGrid.FromEmail)
	notifDispatcher := notificationuc.NewDispatcher(notifRepo, queue)
	hooks := notificationuc.NewHooks(notifDispatcher)
	notifWorker := worker.NewNotificationWorker(queue, pushSender, emailSender, deviceRepo, prefRepo)

	// Usecases.
	authSvc := auth.NewService(userRepo, tokenRepo, sessionRepo, jwtManager).
		WithLockout(auth.LockoutConfig{
			MaxFailures: cfg.Security.LoginMaxFailures,
			LockoutTTL:  cfg.Security.LoginLockoutTTL,
		})
	courseSvc := course.NewService(courseRepo, sectionRepo, lessonRepo, videoRepo, store, txManager)
	// access.Service is the single authority on "can this caller reach this
	// course/lesson?" — quiz/exam taking and paid video all route through it.
	accessSvc := access.NewService(enrollmentRepo, subscriptionRepo, courseRepo, nil)
	quizSvc := quiz.NewService(quizRepo, questionRepo, quizAttemptRepo, courseRepo, attemptStateRepo, txManager, accessSvc)
	examSvc := exam.NewService(examRepo, examAttemptRepo, courseRepo, attemptStateRepo, scoreRankingRepo, txManager, accessSvc, hooks)
	analyticsSvc := analytics.NewService(examRepo, examAttemptRepo, courseRepo, scoreRankingRepo, quizAttemptRepo)
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
		txManager, hooks,
		payment.Config{
			CallbackURL:     cfg.Payment.CallbackURL,
			ReturnURL:       cfg.Payment.ReturnURL,
			DefaultCurrency: cfg.Payment.DefaultCurrency,
		},
		nil,
	)

	// Engagement sweeper warns students whose streak is at risk.
	engSweeper := worker.NewEngagementSweeper(progressRepo, hooks, nil, cfg.EngagementSweep.Interval)

	// Phase 5 usecases.
	progressSvc := progress.NewService(progressRepo, courseRepo, enrollmentRepo, accessSvc, lessonRepo, hooks,
		progress.Config{StreakTZ: "Africa/Addis_Ababa"}, nil)
	notifSvc := notificationuc.NewService(notifRepo, deviceRepo, prefRepo)
	adminSvc := admin.NewService(userRepo, courseRepo, paymentRepo, enrollmentRepo, subscriptionRepo,
		revenueRepo, auditRepo, notifDispatcher, registry)
	pdfRenderer := pdf.NewFPDFRenderer()
	certSvc := certificate.NewService(certRepo, progressRepo, courseRepo, userRepo, store, signer,
		pdfRenderer, hooks, certificate.Config{})

	// Google Drive-only handlers: the signed /media proxy (fronted by an adapter
	// that maps the store's types onto the delivery package's, keeping the HTTP
	// layer free of the storage ring) and the one-time admin connect flow. Both
	// stay nil under the S3/MinIO backend so their routes are never registered.
	var mediaHandler *delivery_http.MediaHandler
	var googleOAuthHandler *delivery_http.GoogleOAuthHandler
	if driveStore != nil {
		mediaHandler = delivery_http.NewMediaHandler(driveMediaAdapter{s: driveStore})
		googleOAuthHandler = delivery_http.NewGoogleOAuthHandler(delivery_http.GoogleOAuthConfig{
			ClientID:     cfg.GoogleDrive.ClientID,
			ClientSecret: cfg.GoogleDrive.ClientSecret,
			RedirectURI:  cfg.GoogleDrive.RedirectURI,
		}, driveCreds, redisClient)
	}

	// Handlers + router.
	router := delivery_http.NewRouter(delivery_http.Deps{
		Config:    cfg,
		DB:        pool,
		Cache:     redisClient,
		JWT:       jwtManager,
		Sessions:  sessionRepo,
		Auth:      delivery_http.NewAuthHandler(authSvc, userRepo),
		Courses:   delivery_http.NewCourseHandler(courseSvc, cfg.App.PublicURL, store),
		Quizzes:   delivery_http.NewQuizHandler(quizSvc),
		Exams:     delivery_http.NewExamHandler(examSvc),
		Analytics: delivery_http.NewAnalyticsHandler(analyticsSvc),
		Video:     delivery_http.NewVideoHandler(videoSvc),

		Payments:      delivery_http.NewPaymentHandler(paymentSvc),
		Enrollments:   delivery_http.NewEnrollmentHandler(enrollmentSvc),
		Subscriptions: delivery_http.NewSubscriptionHandler(subscriptionSvc, paymentSvc),
		Revenue:       delivery_http.NewRevenueHandler(revenueSvc, nil),

		Progress:      delivery_http.NewProgressHandler(progressSvc),
		Notifications: delivery_http.NewNotificationHandler(notifSvc),
		Certificates:  delivery_http.NewCertificateHandler(certSvc),
		Admin:         delivery_http.NewAdminHandler(adminSvc),

		Media:       mediaHandler,
		GoogleOAuth: googleOAuthHandler,
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
		Router:      router,
		Sweeper:     sweeper,
		SubSweeper:  subSweeper,
		Worker:      transcodeWorker,
		NotifWorker: notifWorker,
		EngSweeper:  engSweeper,
		VideoUC:     videoSvc,
	}, nil
}

// driveMediaAdapter adapts *storage.GoogleDriveStore to delivery_http.MediaStore.
// It maps the storage package's DriveContent/ErrObjectNotFound onto the delivery
// package's MediaContent/ErrMediaNotFound, so the HTTP layer serves Drive objects
// without importing the infrastructure ring (clean-architecture dependency rule).
type driveMediaAdapter struct{ s *storage.GoogleDriveStore }

func (a driveMediaAdapter) VerifyMediaURL(key string, exp int64, sig string) bool {
	return a.s.VerifyMediaURL(key, exp, sig)
}

func (a driveMediaAdapter) Stream(ctx context.Context, key, rangeHeader string) (*delivery_http.MediaContent, error) {
	dc, err := a.s.Stream(ctx, key, rangeHeader)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			return nil, delivery_http.ErrMediaNotFound
		}
		return nil, err
	}
	return &delivery_http.MediaContent{
		Body:          dc.Body,
		StatusCode:    dc.StatusCode,
		ContentType:   dc.ContentType,
		ContentLength: dc.ContentLength,
		ContentRange:  dc.ContentRange,
		AcceptRanges:  dc.AcceptRanges,
	}, nil
}
