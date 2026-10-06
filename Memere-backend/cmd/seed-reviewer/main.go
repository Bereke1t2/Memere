// Command seed-reviewer creates a verified student user for Google Play Store review.
// It is idempotent: running it twice updates the password instead of failing.
//
// Usage:
//
//	make seed-reviewer
//	REVIEWER_EMAIL=reviewer@mirkuz.app REVIEWER_PASSWORD=secret make seed-reviewer
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
	email := envOr("REVIEWER_EMAIL", "reviewer@mirkuz.app")
	plain := envOr("REVIEWER_PASSWORD", "MirkuzReview2026!")
	firstName := envOr("REVIEWER_FIRST_NAME", "Google")
	lastName := envOr("REVIEWER_LAST_NAME", "Reviewer")

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

	var id uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO auth.users
		  (id, email, password_hash, role, first_name, last_name,
		   is_active, is_email_verified, created_at, updated_at)
		VALUES
		  (gen_random_uuid(), $1, $2, 'student', $3, $4,
		   true, true, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE
		  SET password_hash     = EXCLUDED.password_hash,
		      role              = 'student',
		      is_active         = true,
		      is_email_verified = true,
		      updated_at        = NOW()
		RETURNING id`,
		email, hash, firstName, lastName,
	).Scan(&id)
	// Grant full access to all published courses so reviewer can access all screens
	tag, err := pool.Exec(ctx, `
		INSERT INTO payments.enrollments (id, student_id, course_id, source, expires_at)
		SELECT gen_random_uuid(), $1, id, 'admin_grant', NOW() + INTERVAL '10 years'
		FROM courses.courses
		ON CONFLICT (student_id, course_id) DO UPDATE
		  SET expires_at = NOW() + INTERVAL '10 years'`,
		id,
	)
	if err != nil {
		log.Printf("warning: grant courses: %v", err)
	} else {
		fmt.Printf("  Access:   %d courses unlocked\n", tag.RowsAffected())
	}

	fmt.Printf("✅ Google Play Reviewer student account ready\n")
	fmt.Printf("  ID:       %s\n", id)
	fmt.Printf("  Email:    %s\n", email)
	fmt.Printf("  Password: %s\n", plain)
	fmt.Printf("  Role:     %s\n", entity.RoleStudent)
	fmt.Printf("  Verified: true\n")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
