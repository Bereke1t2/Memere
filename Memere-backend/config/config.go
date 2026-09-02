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
	GoogleDrive GoogleDriveConfig
	Video      VideoConfig
	Payment    PaymentConfig
	Observability ObservabilityConfig
	Security      SecurityConfig
	Notifications NotificationConfig
	Jobs          JobsConfig

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
	// CloudRunPort mirrors the platform-injected PORT variable (Cloud Run,
	// Heroku, etc.). It is empty in local dev. ListenPort prefers it over Port so
	// the same binary binds the platform's port without any APP_PORT juggling.
	CloudRunPort string `envconfig:"PORT" default:""`
	// PublicURL is the externally reachable base URL of this API (scheme+host+
	// port, no trailing slash). It is used only to build absolute media-proxy
	// URLs when STORAGE_PROVIDER=gdrive (e.g. <PublicURL>/api/v1/media?...), so
	// the mobile client can stream Drive-backed objects through the backend.
	PublicURL string `envconfig:"APP_PUBLIC_URL" default:"http://localhost:8080"`
}

func (a AppConfig) IsProduction() bool { return a.Env == "production" }

// ListenPort is the TCP port the HTTP server binds. It prefers the platform-
// injected PORT (Cloud Run and most PaaS route only to that port) and falls back
// to APP_PORT for local dev / Docker Compose. The server binds ":"+ListenPort,
// i.e. all interfaces (0.0.0.0), as Cloud Run requires.
func (a AppConfig) ListenPort() string {
	if a.CloudRunPort != "" {
		return a.CloudRunPort
	}
	return a.Port
}

type DBConfig struct {
	// URL, when set, is a full connection string (postgres:// or postgresql://) that
	// overrides Host/Port/User/Password/Name/SSLMode — convenient for managed
	// providers like Neon that provide a single connection string.
	URL      string `envconfig:"DATABASE_URL" default:""`
	Host     string `envconfig:"DB_HOST" default:""`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:""`
	Password string `envconfig:"DB_PASSWORD" default:""`
	Name     string `envconfig:"DB_NAME" default:""`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`

	// Pool sizing. On Cloud Run each instance owns its own pool, so the ceiling is
	// MaxConns × max-instances against the database's global connection limit
	// (Neon free tier ~100). Keep MaxConns modest and lean on a small MinConns so
	// scaled-out instances don't exhaust it. MaxConnLifetime recycles connections
	// so a scaled-in instance's conns are released promptly.
	MaxConns        int32         `envconfig:"DB_MAX_CONNS" default:"10"`
	MinConns        int32         `envconfig:"DB_MIN_CONNS" default:"0"`
	MaxConnIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"5m"`
	MaxConnLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"30m"`

	// Pooled is set true when connecting through a transaction-mode pooler
	// (Neon's -pooler endpoint / PgBouncer). That mode does not support the
	// implicit prepared statements pgx caches by default, so the driver must use
	// the simple query protocol. Leave false for a direct endpoint.
	Pooled bool `envconfig:"DB_POOLED" default:"false"`
}

// DSN builds a connection URL for pgx, using URL directly when set.
func (d DBConfig) DSN() string {
	if d.URL != "" {
		return d.URL
	}
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
	// URL, when set, is a full connection string (redis:// or rediss://) that
	// overrides Host/Port/Password/DB/TLS — convenient for managed providers like
	// Upstash that hand you one URL. A rediss:// scheme enables TLS automatically.
	URL      string `envconfig:"REDIS_URL" default:""`
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	// TLS enables a TLS connection (required by Upstash and most managed Redis).
	// Leave false for a local/Docker Redis. Ignored when URL is set (the URL's
	// scheme decides). Upstash uses a valid public cert, so no extra CA config.
	TLS bool `envconfig:"REDIS_TLS" default:"false"`

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

	// TrustedProxies lists the CIDRs/IPs Gin trusts to set X-Forwarded-For, so
	// client IPs (used by the rate limiter and logs) are read from the real edge
	// rather than spoofable headers. Behind Cloud Run set this to the container
	// network range that fronts the app (see README); empty means "trust none"
	// (Gin uses the direct RemoteAddr — the safe default on an untrusted network).
	TrustedProxies []string `envconfig:"TRUSTED_PROXIES" default:""`
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
// (spec §3.3, §8). The same shape serves AWS S3 in production, MinIO in local
// dev (Endpoint + UsePathStyle set), and Backblaze B2 (STORAGE_PROVIDER=b2) —
// B2 speaks the S3 API, so it reuses the S3 store with its own endpoint and
// application-key credentials (the B2* fields below). Bucket is intentionally
// NOT required so the API still boots when video is unconfigured; the S3 store
// validates it at construction and the upload usecase fails clearly if empty.
type StorageConfig struct {
	Provider        string `envconfig:"STORAGE_PROVIDER" default:"s3"` // s3 | minio | b2 | gdrive
	Endpoint        string `envconfig:"S3_ENDPOINT"`                   // empty for AWS; set for MinIO
	Region          string `envconfig:"AWS_REGION" default:"af-south-1"`
	Bucket          string `envconfig:"AWS_S3_BUCKET" default:"memere-media"`
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	UsePathStyle    bool   `envconfig:"S3_USE_PATH_STYLE" default:"false"` // true for MinIO

	// Backblaze B2 (S3-compatible API), used when STORAGE_PROVIDER=b2. B2's S3
	// endpoint takes the Application Key ID as the access key and the Application
	// Key as the secret. Endpoint + region come from the bucket's details in the
	// B2 console, e.g. Endpoint "s3.us-west-004.backblazeb2.com" -> region
	// "us-west-004". The scheme is optional; the wiring prepends https:// if
	// omitted. Only the backend ever holds these values (never sent to a client).
	B2KeyID    string `envconfig:"B2_APPLICATION_KEY_ID"`
	B2AppKey   string `envconfig:"B2_APPLICATION_KEY"`
	B2Endpoint string `envconfig:"B2_ENDPOINT"`
	B2Region   string `envconfig:"B2_REGION"`
	B2Bucket   string `envconfig:"B2_BUCKET"`

	UploadURLTTL   time.Duration `envconfig:"UPLOAD_URL_TTL" default:"15m"`
	StreamURLTTL   time.Duration `envconfig:"STREAM_URL_TTL" default:"2h"` // spec §8.3: 2h
	DownloadURLTTL time.Duration `envconfig:"DOWNLOAD_URL_TTL" default:"2h"`

	CDNDomain        string `envconfig:"CDN_DOMAIN"`          // e.g. dxxxx.cloudfront.net
	CDNKeyPairID     string `envconfig:"CDN_KEY_PAIR_ID"`     // CloudFront signing
	CDNPrivateKeyPEM string `envconfig:"CDN_PRIVATE_KEY_PEM"` // PEM contents (or path)
}

// B2EndpointURL returns the Backblaze B2 S3 endpoint as an absolute URL. B2's
// console shows the endpoint without a scheme (e.g. s3.us-west-004.backblazeb2.com);
// the AWS SDK needs a full URL, so we default a missing scheme to https. An empty
// B2_ENDPOINT is returned unchanged (the store then errors clearly).
func (s StorageConfig) B2EndpointURL() string {
	if s.B2Endpoint == "" || strings.Contains(s.B2Endpoint, "://") {
		return s.B2Endpoint
	}
	return "https://" + s.B2Endpoint
}

// GoogleDriveConfig groups the settings for the single admin-owned Google Drive
// account that backs object storage when STORAGE_PROVIDER=gdrive. ONLY the
// backend ever holds these values — they are never sent to a client, never
// logged, and never committed (spec: students/teachers must not authenticate
// with Google; one admin account owns the Drive).
//
// The refresh token is normally obtained ONCE through the admin connect page
// (GET /api/v1/admin/google/connect) and stored server-side (encrypted in the
// storage.google_credentials table). GOOGLE_REFRESH_TOKEN is an alternative for
// deployments that inject it from a secret manager; the store prefers the DB
// value and falls back to this env var.
type GoogleDriveConfig struct {
	ClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
	ClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`
	RedirectURI  string `envconfig:"GOOGLE_REDIRECT_URI" default:"http://localhost:8080/api/v1/admin/google/callback"`
	RefreshToken string `envconfig:"GOOGLE_REFRESH_TOKEN"`
	RootFolderID string `envconfig:"GOOGLE_DRIVE_ROOT_FOLDER_ID"`
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
	// WorkerEnabled runs the in-process transcode worker (needs FFmpeg on PATH).
	// The scale-to-zero Cloud Run API image is FFmpeg-free and CPU-throttled
	// between requests, so it sets this false and lets the dedicated worker image
	// (Dockerfile.worker) own transcoding. Defaults true for local dev / Compose.
	WorkerEnabled bool `envconfig:"TRANSCODE_WORKER_ENABLED" default:"true"`
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

// NotificationConfig gates the in-process notification worker (Phase 5). On a
// scale-to-zero Cloud Run instance whose CPU is throttled between requests, an
// always-on delivery loop can't run reliably, so production disables it here and
// drives delivery from the dedicated worker image / a scheduled sweep instead.
// Defaults true so local dev and Docker Compose keep their current behavior.
type NotificationConfig struct {
	WorkerEnabled bool `envconfig:"NOTIFICATION_WORKER_ENABLED" default:"true"`
}

// JobsConfig secures the internal /internal/jobs/* endpoints that Cloud
// Scheduler calls to drive periodic sweeps (attempt expiry, subscription expiry,
// engagement) when the in-process tickers are disabled on a scale-to-zero
// service. InternalToken is a shared secret compared in constant time; when it
// is empty the internal routes are not registered at all (local dev keeps using
// the in-process tickers, so it needs no token). Never hardcode it — inject via
// a secret and give the same value to the Scheduler job's header.
type JobsConfig struct {
	InternalToken string `envconfig:"INTERNAL_JOB_TOKEN" default:""`
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
	if cfg.DB.URL == "" && (cfg.DB.Host == "" || cfg.DB.User == "" || cfg.DB.Password == "" || cfg.DB.Name == "") {
		return nil, fmt.Errorf("load config: either DATABASE_URL or DB_HOST/DB_USER/DB_PASSWORD/DB_NAME must be set")
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
	if c.Storage.Provider == "gdrive" {
		if c.GoogleDrive.ClientID == "" || c.GoogleDrive.ClientSecret == "" {
			errs = append(errs, "GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET are required when STORAGE_PROVIDER=gdrive")
		}
		if c.GoogleDrive.RootFolderID == "" {
			errs = append(errs, "GOOGLE_DRIVE_ROOT_FOLDER_ID is required when STORAGE_PROVIDER=gdrive")
		}
	}
	if c.Storage.Provider == "b2" {
		if c.Storage.B2KeyID == "" || c.Storage.B2AppKey == "" {
			errs = append(errs, "B2_APPLICATION_KEY_ID and B2_APPLICATION_KEY are required when STORAGE_PROVIDER=b2")
		}
		if c.Storage.B2Bucket == "" {
			errs = append(errs, "B2_BUCKET is required when STORAGE_PROVIDER=b2")
		}
		if c.Storage.B2Endpoint == "" {
			errs = append(errs, "B2_ENDPOINT is required when STORAGE_PROVIDER=b2 (e.g. s3.us-west-004.backblazeb2.com)")
		}
		if c.Storage.B2Region == "" {
			errs = append(errs, "B2_REGION is required when STORAGE_PROVIDER=b2 (e.g. us-west-004)")
		}
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
