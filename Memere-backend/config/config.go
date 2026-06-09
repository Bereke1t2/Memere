package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config is the typed application configuration loaded from the environment.
// Groups map to the README "Environment Variables" table. Later-phase groups
// (AWS, Chapa, Stripe, FCM, SendGrid) are declared now so the shape is stable;
// they are not required until the phase that uses them.
type Config struct {
	App     AppConfig
	DB      DBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	HTTP    HTTPConfig
	Sweeper SweeperConfig

	// Reserved for later phases (not required to boot in Phase 1).
	AWS      AWSConfig
	Chapa    ChapaConfig
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

type AWSConfig struct {
	S3Bucket string `envconfig:"AWS_S3_BUCKET" default:""`
	Region   string `envconfig:"AWS_REGION" default:"af-south-1"`
}

type ChapaConfig struct {
	SecretKey string `envconfig:"CHAPA_SECRET_KEY" default:""`
}

type StripeConfig struct {
	SecretKey string `envconfig:"STRIPE_SECRET_KEY" default:""`
}

type FCMConfig struct {
	ServerKey string `envconfig:"FCM_SERVER_KEY" default:""`
}

type SendGridConfig struct {
	APIKey string `envconfig:"SENDGRID_API_KEY" default:""`
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
