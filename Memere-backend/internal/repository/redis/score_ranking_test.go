package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

// testClient connects to a local Redis for integration tests, skipping the test
// when none is reachable (so unit runs in CI without Redis still pass). Set
// REDIS_TEST_ADDR to override the address.
func testClient(t *testing.T) *goredis.Client {
	t.Helper()
	addr := os.Getenv("REDIS_TEST_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	c := goredis.NewClient(&goredis.Options{Addr: addr, Password: os.Getenv("REDIS_PASSWORD")})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := c.Ping(ctx).Err(); err != nil {
		c.Close()
		t.Skipf("redis not reachable at %s: %v", addr, err)
	}
	return c
}

func TestScoreRanking_PercentileFromSortedSet(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	repo := NewScoreRankingRepo(c)
	ctx := context.Background()

	examID := uuid.New() // unique key per run; no cleanup collisions
	defer c.Del(ctx, examScoresKey(examID))

	// Five students with distinct percentages.
	students := make([]uuid.UUID, 5)
	scores := []float64{20, 40, 60, 80, 100}
	for i := range students {
		students[i] = uuid.New()
		if err := repo.RecordExamScore(ctx, examID, students[i], scores[i]); err != nil {
			t.Fatalf("record: %v", err)
		}
	}

	// The top scorer (100) ties-or-beats all 5 → 100th percentile.
	pct, ok, err := repo.PercentileRank(ctx, examID, students[4])
	if err != nil || !ok {
		t.Fatalf("percentile top: ok=%v err=%v", ok, err)
	}
	if pct != 100 {
		t.Errorf("top scorer percentile = %v, want 100", pct)
	}

	// The 60 scorer ties-or-beats {20,40,60} = 3 of 5 → 60th percentile.
	pct, ok, _ = repo.PercentileRank(ctx, examID, students[2])
	if !ok || pct != 60 {
		t.Errorf("middle scorer percentile = %v, want 60", pct)
	}

	// An unknown student has no recorded score.
	if _, ok, _ := repo.PercentileRank(ctx, examID, uuid.New()); ok {
		t.Error("unknown student should report ok=false")
	}
}

func TestScoreRanking_KeepsBestScore(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	repo := NewScoreRankingRepo(c)
	ctx := context.Background()

	examID := uuid.New()
	student := uuid.New()
	defer c.Del(ctx, examScoresKey(examID))

	if err := repo.RecordExamScore(ctx, examID, student, 90); err != nil {
		t.Fatalf("record high: %v", err)
	}
	// A later, lower score must not regress the recorded best (ZADD GT).
	if err := repo.RecordExamScore(ctx, examID, student, 50); err != nil {
		t.Fatalf("record low: %v", err)
	}
	got, err := c.ZScore(ctx, examScoresKey(examID), student.String()).Result()
	if err != nil {
		t.Fatalf("zscore: %v", err)
	}
	if got != 90 {
		t.Errorf("stored score = %v, want 90 (best kept)", got)
	}
}
