// Package metrics registers all Prometheus metrics for the Memere API.
// Import this package in main.go via a blank import or by calling Register().
// All metrics use consistent label sets to avoid cardinality explosions.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RED metrics (Rate / Errors / Duration) per HTTP endpoint.
var (
	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by method, route template, and status class.",
	}, []string{"method", "route", "status"})

	HTTPDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency (seconds). Buckets tuned for p95 < 200ms SLO.",
		Buckets: []float64{.005, .01, .025, .05, .1, .2, .5, 1, 2, 5},
	}, []string{"method", "route"})
)

// Business event counters — incremented by usecases / workers via the
// exported Observe* helpers below. Labels are intentionally coarse to keep
// series count low.
var (
	PaymentsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_payments_total",
		Help: "Payment events (completed / failed).",
	}, []string{"outcome"}) // labels: "completed", "failed"

	EnrollmentsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "memere_enrollments_total",
		Help: "Total enrollment grants (free + paid).",
	})

	VideoTranscodeTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_video_transcode_total",
		Help: "Video transcode completions by outcome.",
	}, []string{"outcome"}) // labels: "success", "failure"

	AttemptsGradedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_attempts_graded_total",
		Help: "Quiz/exam attempts graded by type.",
	}, []string{"type"}) // labels: "quiz", "exam"

	NotificationsSentTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_notifications_sent_total",
		Help: "Notification dispatches by channel.",
	}, []string{"channel"}) // labels: "push", "email", "in_app"

	CertificatesIssuedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "memere_certificates_issued_total",
		Help: "Completion certificates issued.",
	})
)

// Cache hit/miss counters (Phase 6 Skill 3).
var (
	CacheHitTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_cache_hit_total",
		Help: "Redis cache hits by key class.",
	}, []string{"class"})

	CacheMissTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_cache_miss_total",
		Help: "Redis cache misses (loaded from DB) by key class.",
	}, []string{"class"})
)

func ObserveCacheHit(class string)  { CacheHitTotal.WithLabelValues(class).Inc() }
func ObserveCacheMiss(class string) { CacheMissTotal.WithLabelValues(class).Inc() }

// Security event counters (Phase 6 Skill 2).
var (
	RateLimitBlockTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_rate_limit_block_total",
		Help: "Requests blocked by the rate limiter, by limiter type.",
	}, []string{"limiter"}) // labels: "global", "auth", "user", "provider"

	SuspiciousLoginTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "memere_suspicious_login_total",
		Help: "Account lockout events triggered by too many failed logins.",
	}, []string{"reason"}) // labels: "lockout"
)

// Worker / queue depth gauges.
var (
	WorkerQueueDepth = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memere_worker_queue_depth",
		Help: "Approximate number of jobs pending in worker queues.",
	}, []string{"queue"}) // labels: "transcode", "notification"
)

// DB pool gauges — set by a periodic scrape in main.go.
var (
	DBPoolAcquired = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_acquired_connections",
		Help: "Number of connections currently acquired from the pgx pool.",
	})
	DBPoolTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_total_connections",
		Help: "Total connections in the pgx pool (acquired + idle).",
	})
)

// ----- Observer helpers (called by application code) ------------------------

func ObserveRateLimitBlock(limiter string)       { RateLimitBlockTotal.WithLabelValues(limiter).Inc() }
func ObserveSuspiciousLogin(reason string)       { SuspiciousLoginTotal.WithLabelValues(reason).Inc() }

func ObservePayment(outcome string)              { PaymentsTotal.WithLabelValues(outcome).Inc() }
func ObserveEnrollment()                        { EnrollmentsTotal.Inc() }
func ObserveTranscode(outcome string)           { VideoTranscodeTotal.WithLabelValues(outcome).Inc() }
func ObserveAttemptGraded(kind string)          { AttemptsGradedTotal.WithLabelValues(kind).Inc() }
func ObserveNotification(channel string)        { NotificationsSentTotal.WithLabelValues(channel).Inc() }
func ObserveCertificateIssued()                 { CertificatesIssuedTotal.Inc() }
func SetQueueDepth(queue string, depth float64) { WorkerQueueDepth.WithLabelValues(queue).Set(depth) }
