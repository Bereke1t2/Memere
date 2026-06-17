// Command seed-admin creates an admin user for local development and testing.
// It is idempotent: running it twice updates the password instead of failing.
//
// Usage:
//
//	make seed-admin                                 # default creds
//	ADMIN_EMAIL=me@example.com ADMIN_PASSWORD=secret make seed-admin
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/database"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/password"
)

func main() {
	email := envOr("ADMIN_EMAIL", "admin@memere.app")
	plain := envOr("ADMIN_PASSWORD", "Admin1234!")
	firstName := envOr("ADMIN_FIRST_NAME", "Platform")
	lastName := envOr("ADMIN_LAST_NAME", "Admin")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	hash, err := password.Hash(plain)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	email = strings.ToLower(strings.TrimSpace(email))

	// Upsert: insert new admin or update the password + role if the email exists.
	var id uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO auth.users
		  (id, email, password_hash, role, first_name, last_name,
		   is_active, is_email_verified, created_at, updated_at)
		VALUES
		  (gen_random_uuid(), $1, $2, 'admin', $3, $4,
		   true, true, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE
		  SET password_hash    = EXCLUDED.password_hash,
		      role             = 'admin',
		      is_active        = true,
		      is_email_verified = true,
		      updated_at       = NOW()
		RETURNING id`,
		email, hash, firstName, lastName,
	).Scan(&id)
	if err != nil {
		log.Fatalf("upsert admin: %v", err)
	}

	fmt.Printf("Admin user ready\n")
	fmt.Printf("  ID:       %s\n", id)
	fmt.Printf("  Email:    %s\n", email)
	fmt.Printf("  Password: %s\n", plain)
	fmt.Printf("  Role:     %s\n", entity.RoleAdmin)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
