package config_test

import (
	"testing"

	"github.com/Bereke1t2/Memere/memere-backend/config"
)

func prodConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{Env: "production"},
		JWT: config.JWTConfig{Secret: "this-is-a-strong-secret-for-testing-32ch"},
		HTTP: config.HTTPConfig{
			CORSAllowedOrigins: []string{"https://app.memere.et"},
		},
		Security: config.SecurityConfig{
			MinJWTSecretLen: 32,
		},
	}
}

func TestValidateProduction_OK(t *testing.T) {
	cfg := prodConfig()
	if err := cfg.ValidateProduction(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateProduction_WeakSecret(t *testing.T) {
	cfg := prodConfig()
	cfg.JWT.Secret = "short"
	if err := cfg.ValidateProduction(); err == nil {
		t.Error("expected error for weak JWT secret, got nil")
	}
}

func TestValidateProduction_WildcardCORS(t *testing.T) {
	cfg := prodConfig()
	cfg.HTTP.CORSAllowedOrigins = []string{"*"}
	if err := cfg.ValidateProduction(); err == nil {
		t.Error("expected error for wildcard CORS in production, got nil")
	}
}

func TestValidateProduction_MissingChapaWebhookSecret(t *testing.T) {
	cfg := prodConfig()
	cfg.Chapa.SecretKey = "sk_live_xxx"
	cfg.Chapa.WebhookSecret = ""
	if err := cfg.ValidateProduction(); err == nil {
		t.Error("expected error for missing Chapa webhook secret, got nil")
	}
}

func TestValidateProduction_DevMode_NoCheck(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{Env: "development"},
		JWT: config.JWTConfig{Secret: "short"},
		HTTP: config.HTTPConfig{CORSAllowedOrigins: []string{"*"}},
	}
	if err := cfg.ValidateProduction(); err != nil {
		t.Errorf("dev mode must not validate: %v", err)
	}
}
