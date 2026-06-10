package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
)

// TestVideoFromRow pins the sqlc-model -> entity mapping, including the
// nullable-key and nullable-timestamp paths (no DB required).
func TestVideoFromRow(t *testing.T) {
	id, lesson, course := uuid.New(), uuid.New(), uuid.New()
	now := time.Now().UTC().Truncate(time.Second)
	key := "originals/c/l/v/source.mp4"

	v := videoFromRow(sqlcgen.CoursesVideo{
		ID:               toPgUUID(id),
		LessonID:         toPgUUID(lesson),
		CourseID:         toPgUUID(course),
		ProcessingStatus: string(entity.VideoProcessing),
		OriginalFileKey:  &key,
		HlsMasterKey:     nil,
		DurationSeconds:  42,
		FileSizeBytes:    1024,
		CreatedAt:        pgTimestamptzValue(now),
		UpdatedAt:        pgTimestamptzValue(now),
		ProcessedAt:      pgtype.Timestamptz{}, // NULL
		DeletedAt:        pgtype.Timestamptz{}, // NULL
	})

	if v.ID != id || v.LessonID != lesson || v.CourseID != course {
		t.Fatalf("uuid mapping wrong: %+v", v)
	}
	if v.Status != entity.VideoProcessing {
		t.Fatalf("status = %q, want processing", v.Status)
	}
	if v.OriginalFileKey == nil || *v.OriginalFileKey != key {
		t.Fatalf("original key not mapped: %v", v.OriginalFileKey)
	}
	if v.HLSMasterKey != nil {
		t.Fatalf("NULL key should map to nil, got %v", *v.HLSMasterKey)
	}
	if v.DurationSeconds != 42 || v.FileSizeBytes != 1024 {
		t.Fatalf("numeric mapping wrong: dur=%d size=%d", v.DurationSeconds, v.FileSizeBytes)
	}
	if v.ProcessedAt != nil || v.DeletedAt != nil {
		t.Fatalf("NULL timestamps should map to nil")
	}
	if !v.CreatedAt.Equal(now) {
		t.Fatalf("created_at = %v, want %v", v.CreatedAt, now)
	}
}

// TestVideoUpdateStatusGuardedNoRow proves the guarded transition reports
// claimed=false (no error) when no row matches the from-status — here via a
// random id that matches nothing. Skips when no database is configured.
func TestVideoUpdateStatusGuardedNoRow(t *testing.T) {
	dsn := dsnFromEnv()
	if dsn == "" {
		t.Skip("no database env (DB_HOST/TEST_DATABASE_URL); skipping integration check")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("cannot create pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("database unreachable: %v", err)
	}

	repo := NewVideoRepo(pool)
	ok, err := repo.UpdateStatus(ctx, uuid.New(), entity.VideoPending, entity.VideoProcessing)
	if err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
	if ok {
		t.Fatalf("guarded update on a nonexistent id must be false")
	}
}

func dsnFromEnv() string {
	if dsn := os.Getenv("TEST_DATABASE_URL"); dsn != "" {
		return dsn
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		envOr("DB_USER", "postgres"), os.Getenv("DB_PASSWORD"), host,
		envOr("DB_PORT", "5432"), envOr("DB_NAME", "memere"),
		envOr("DB_SSL_MODE", "disable"))
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
