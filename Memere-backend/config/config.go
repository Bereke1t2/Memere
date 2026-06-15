package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config is the typed application configuration loaded from the environment.
// Groups map to the README "Environment Variables" table. Later-phase groups
// (AWS, Chapa, Stripe, FCM, SendGrid) are declared now so the shape is stable;
// they are not required until the phase that uses them.
type Config struct {
	App        AppConfig
	DB         DBConfig
	Redis      RedisConfig
	JWT        JWTConfig
	HTTP       HTTPConfig
	Sweeper          SweeperConfig
	SubSweeper       SubscriptionSweeperConfig
	EngagementSweep  EngagementSweeperConfig
	Storage    StorageConfig
	Video      VideoConfig
	Payment    PaymentConfig
	Observability ObservabilityConfig
	Security      SecurityConfig

	// Reserved for later phases (not required to boot in Phase 1).
	AWS      AWSConfig
	Chapa    ChapaConfig
	Telebirr TelebirrConfig
	Stripe   StripeConfig
	FCM      FCMConfig
	SendGrid SendGridConfig
}

type AppConfig struct {
	Env  string `envconfig:"APP_ENV" default:"development"`
	Port string `envconfig:"APP_PORT" default:"8080"`
}

func (a AppConfig) IsProduction() bool { return a.Env == "production" }

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`

	MaxConns        int32         `envconfig:"DB_MAX_CONNS" default:"20"`
	MaxConnIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"5m"`
}

// DSN builds a libpq-style connection URL for pgx.
func (d DBConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%s", d.Host, d.Port),
		Path:   d.Name,
	}
	q := url.Values{}
	q.Set("sslmode", d.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" required:"true"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`

	// Pool sizing and timeouts (Phase 6 Skill 3 §12.2).
	PoolSize    int           `envconfig:"REDIS_POOL_SIZE"    default:"20"`
	DialTimeout time.Duration `envconfig:"REDIS_DIAL_TIMEOUT" default:"3s"`
	ReadTimeout time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"2s"`
	WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"2s"`
}

func (r RedisConfig) Addr() string { return fmt.Sprintf("%s:%s", r.Host, r.Port) }

type JWTConfig struct {
	Secret     string        `envconfig:"JWT_SECRET" required:"true"`
	AccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
	RefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"720h"`
	Issuer     string        `envconfig:"JWT_ISSUER" default:"memere"`
}

// HTTPConfig tunes the delivery layer: CORS and the Redis-backed rate limiters
// (spec §5.4, §7.3). A permissive CORS default suits development; production
// should pin CORS_ALLOWED_ORIGINS to the app's domains.
type HTTPConfig struct {
	// CORSAllowedOrigins is the list of allowed Origin values; the single value
	// "*" allows any origin.
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	// RateLimitRPM is the global per-IP request budget per minute.
	RateLimitRPM int `envconfig:"RATE_LIMIT_RPM" default:"120"`
	// LoginRateLimit / LoginRateWindow are the stricter limiter on /auth/login
	// (default 5 attempts per 15 minutes per IP, spec §7.3).
	LoginRateLimit  int           `envconfig:"LOGIN_RATE_LIMIT" default:"5"`
	LoginRateWindow time.Duration `envconfig:"LOGIN_RATE_WINDOW" default:"15m"`
}

// SweeperConfig tunes the background expiry sweeper (Phase 2 §9.2). Interval is
// how often it scans for abandoned attempts past their deadline; Enabled lets a
// deployment (or test) turn it off.
type SweeperConfig struct {
	Enabled  bool          `envconfig:"SWEEPER_ENABLED" default:"true"`
	Interval time.Duration `envconfig:"SWEEPER_INTERVAL" default:"60s"`
}

// SubscriptionSweeperConfig tunes the Phase 4 subscription expiry sweeper
// (spec §10.2): how often it scans for lapsed subscriptions to expire, and a
// switch to disable it in a deployment that handles renewals elsewhere.
type SubscriptionSweeperConfig struct {
	Enabled  bool          `envconfig:"SUBSCRIPTION_SWEEP_ENABLED" default:"true"`
	Interval time.Duration `envconfig:"SUBSCRIPTION_SWEEP_INTERVAL" default:"1h"`
}

// EngagementSweeperConfig tunes the Phase 5 streak-warning sweeper (spec §11.2).
type EngagementSweeperConfig struct {
	Enabled  bool          `envconfig:"ENGAGEMENT_SWEEP_ENABLED" default:"true"`
	Interval time.Duration `envconfig:"ENGAGEMENT_SWEEP_INTERVAL" default:"24h"`

}

// PaymentConfig groups the Phase 4 payment-flow settings (spec §10): the URLs we
// hand the provider at checkout, the default currency, the configurable
// teacher/platform revenue split (§1.6), config-driven subscription plan pricing
// (§1.6), and the test-only mock provider switch used by the smoke test.
type PaymentConfig struct {
	// CallbackURL is the base webhook URL we give providers; the provider name is
	// appended (…/webhooks/payments/<provider>).
	CallbackURL     string `envconfig:"PAYMENT_CALLBACK_URL" default:"http://localhost:8080/api/v1/webhooks/payments"`
	ReturnURL       string `envconfig:"PAYMENT_RETURN_URL" default:"http://localhost:8080/payments/done"`
	DefaultCurrency string `envconfig:"PAYMENT_DEFAULT_CURRENCY" default:"ETB"`

	// TeacherShare is the fraction of gross a teacher keeps (spec §1.6, 70/30).
	TeacherShare float64 `envconfig:"TEACHER_REVENUE_SHARE" default:"0.70"`

	// MockEnabled registers the test-only mock provider; MockWebhookSecret is the
	// HMAC secret it verifies simulated webhooks against.
	MockEnabled       bool   `envconfig:"PAYMENT_MOCK_ENABLED" default:"false"`
	MockWebhookSecret string `envconfig:"PAYMENT_MOCK_WEBHOOK_SECRET" default:"mock-secret"`

	// Subscription plan pricing (config-driven, never hardcoded in the usecase).
	MonthlyPrice    float64       `envconfig:"SUB_MONTHLY_PRICE" default:"299"`
	MonthlyCurrency string        `envconfig:"SUB_MONTHLY_CURRENCY" default:"ETB"`
	MonthlyPeriod   time.Duration `envconfig:"SUB_MONTHLY_PERIOD" default:"720h"` // 30 days
	AnnualPrice     float64       `envconfig:"SUB_ANNUAL_PRICE" default:"2999"`
	AnnualCurrency  string        `envconfig:"SUB_ANNUAL_CURRENCY" default:"ETB"`
	AnnualPeriod    time.Duration `envconfig:"SUB_ANNUAL_PERIOD" default:"8760h"` // 365 days
}

// StorageConfig groups object-storage settings for the Phase 3 video pipeline
// (spec §3.3, §8). The same shape serves AWS S3 in production and MinIO in local
// dev (Endpoint + UsePathStyle set). Bucket is intentionally NOT required so the
// API still boots when video is unconfigured; the S3 store validates it at
// construction and the upload usecase fails clearly if it is empty.
type StorageConfig struct {
	Provider        string `envconfig:"STORAGE_PROVIDER" default:"s3"` // s3 | minio
	Endpoint        string `envconfig:"S3_ENDPOINT"`                   // empty for AWS; set for MinIO
	Region          string `envconfig:"AWS_REGION" default:"af-south-1"`
	Bucket          string `envconfig:"AWS_S3_BUCKET" default:"memere-media"`
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	UsePathStyle    bool   `envconfig:"S3_USE_PATH_STYLE" default:"false"` // true for MinIO

	UploadURLTTL   time.Duration `envconfig:"UPLOAD_URL_TTL" default:"15m"`
	StreamURLTTL   time.Duration `envconfig:"STREAM_URL_TTL" default:"2h"` // spec §8.3: 2h
	DownloadURLTTL time.Duration `envconfig:"DOWNLOAD_URL_TTL" default:"2h"`

	CDNDomain        string `envconfig:"CDN_DOMAIN"`          // e.g. dxxxx.cloudfront.net
	CDNKeyPairID     string `envconfig:"CDN_KEY_PAIR_ID"`     // CloudFront signing
	CDNPrivateKeyPEM string `envconfig:"CDN_PRIVATE_KEY_PEM"` // PEM contents (or path)
}

// VideoConfig tunes the upload + transcode pipeline (Phase 3 §8.1). MaxUpload
// bounds the size a client may request a pre-signed URL for (rejected before any
// DB write); QueueBuffer sizes the in-process transcode channel; MaxAttempts
// caps how many times a failed transcode is retried before it stays failed.
type VideoConfig struct {
	MaxUploadBytes int64 `envconfig:"MAX_UPLOAD_BYTES" default:"2147483648"` // 2 GiB
	QueueBuffer    int   `envconfig:"TRANSCODE_QUEUE_BUFFER" default:"64"`
	MaxAttempts    int   `envconfig:"TRANSCODE_MAX_ATTEMPTS" default:"3"`
	// Concurrency bounds parallel ffmpeg runs; 0 means the worker picks
	// min(2, NumCPU). WorkDir is the scratch space for downloads/outputs; empty
	// means os.TempDir()/memere-transcode.
	Concurrency int    `envconfig:"TRANSCODE_CONCURRENCY" default:"0"`
	WorkDir     string `envconfig:"TRANSCODE_WORKDIR"`
}

type AWSConfig struct {
	S3Bucket string `envconfig:"AWS_S3_BUCKET" default:""`
	Region   string `envconfig:"AWS_REGION" default:"af-south-1"`
}

type ChapaConfig struct {
	SecretKey     string `envconfig:"CHAPA_SECRET_KEY" default:""`
	WebhookSecret string `envconfig:"CHAPA_WEBHOOK_SECRET" default:""`
	BaseURL       string `envconfig:"CHAPA_BASE_URL" default:""`
}

type StripeConfig struct {
	SecretKey     string `envconfig:"STRIPE_SECRET_KEY" default:""`
	WebhookSecret string `envconfig:"STRIPE_WEBHOOK_SECRET" default:""`
}

type TelebirrConfig struct {
	AppKey        string `envconfig:"TELEBIRR_APP_KEY" default:""`
	WebhookSecret string `envconfig:"TELEBIRR_WEBHOOK_SECRET" default:""`
}

type FCMConfig struct {
	ServerKey string `envconfig:"FCM_SERVER_KEY" default:""`
}

type SendGridConfig struct {
	APIKey    string `envconfig:"SENDGRID_API_KEY" default:""`
	FromEmail string `envconfig:"SENDGRID_FROM_EMAIL" default:"noreply@memere.app"`
}

// SecurityConfig groups Phase 6 / Skill 2 hardening settings.
//
//   - LoginMaxFailures:   Failed-login attempts before the account is locked.
//   - LoginLockoutTTL:    How long a locked account stays locked.
//   - BodyLimitBytes:     Maximum JSON request body size (all non-upload routes).
//   - MinJWTSecretLen:    Minimum length of JWT_SECRET in production.
type SecurityConfig struct {
	LoginMaxFailures int           `envconfig:"LOGIN_MAX_FAILURES"  default:"10"`
	LoginLockoutTTL  time.Duration `envconfig:"LOGIN_LOCKOUT_TTL"   default:"30m"`
	BodyLimitBytes   int64         `envconfig:"BODY_LIMIT_BYTES"    default:"1048576"` // 1 MiB
	MinJWTSecretLen  int           `envconfig:"MIN_JWT_SECRET_LEN"  default:"32"`
}

// ObservabilityConfig groups Phase 6 telemetry settings.
//
//   - LogLevel:      stdlib slog level ("debug" | "info" | "warn" | "error").
//   - OTELEndpoint:  OTLP/HTTP collector URL, "stdout" for dev, "" for no-op.
//   - MetricsPort:   Port for the internal Prometheus /metrics HTTP server.
//     Defaults to 9090; set to "" to disable.
type ObservabilityConfig struct {
	LogLevel     string `envconfig:"LOG_LEVEL"      default:"info"`
	OTELEndpoint string `envconfig:"OTEL_ENDPOINT"  default:""`
	MetricsPort  string `envconfig:"METRICS_PORT"   default:"9090"`
}

// Load reads configuration from the environment, applying defaults and failing
// fast if any required variable is missing.
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &cfg, nil
}

// ValidateProduction fails fast when APP_ENV=production and any security-
// critical setting is absent or dangerously weak (§7.3 secret management):
//
//   - JWT_SECRET must be at least MinJWTSecretLen characters.
//   - CORS_ALLOWED_ORIGINS must not be "*" (fail-closed §2.5).
//   - Each enabled payment provider must have its webhook secret set.
func (c *Config) ValidateProduction() error {
	if !c.App.IsProduction() {
		return nil
	}
	var errs []string

	if len(c.JWT.Secret) < c.Security.MinJWTSecretLen {
		errs = append(errs, fmt.Sprintf("JWT_SECRET must be at least %d characters in production", c.Security.MinJWTSecretLen))
	}
	if len(c.HTTP.CORSAllowedOrigins) == 1 && c.HTTP.CORSAllowedOrigins[0] == "*" {
		errs = append(errs, "CORS_ALLOWED_ORIGINS must not be '*' in production — set explicit allowed origins")
	}
	if c.Chapa.SecretKey != "" && c.Chapa.WebhookSecret == "" {
		errs = append(errs, "CHAPA_WEBHOOK_SECRET required when CHAPA_SECRET_KEY is set")
	}
	if c.Stripe.SecretKey != "" && c.Stripe.WebhookSecret == "" {
		errs = append(errs, "STRIPE_WEBHOOK_SECRET required when STRIPE_SECRET_KEY is set")
	}
	if len(errs) > 0 {
		return fmt.Errorf("production config invalid:\n  - %s", joinErrs(errs))
	}
	return nil
}

func joinErrs(errs []string) string {
	var b strings.Builder
	for i, e := range errs {
		if i > 0 {
			b.WriteString("\n  - ")
		}
		b.WriteString(e)
	}
	return b.String()
}
