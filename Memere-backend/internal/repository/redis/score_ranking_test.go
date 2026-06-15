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

func TestLeaderboard_GetTopN(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	repo := NewScoreRankingRepo(c)
	ctx := context.Background()

	examID := uuid.New()
	defer c.Del(ctx, examScoresKey(examID))

	students := make([]uuid.UUID, 5)
	scores := []float64{10, 30, 50, 80, 95}
	for i := range students {
		students[i] = uuid.New()
		if err := repo.RecordExamScore(ctx, examID, students[i], scores[i]); err != nil {
			t.Fatalf("record: %v", err)
		}
	}

	top3, err := repo.GetTopN(ctx, examID, 3)
	if err != nil {
		t.Fatalf("GetTopN: %v", err)
	}
	if len(top3) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top3))
	}
	// First entry must have the highest score (95) and rank 1.
	if top3[0].Score != 95 || top3[0].Rank != 1 {
		t.Errorf("top entry: score=%v rank=%v, want score=95 rank=1", top3[0].Score, top3[0].Rank)
	}
	if top3[1].Rank != 2 || top3[2].Rank != 3 {
		t.Errorf("ranks not sequential: got %v %v", top3[1].Rank, top3[2].Rank)
	}
}

func TestLeaderboard_GetRank(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	repo := NewScoreRankingRepo(c)
	ctx := context.Background()

	examID := uuid.New()
	defer c.Del(ctx, examScoresKey(examID))

	s1, s2, s3 := uuid.New(), uuid.New(), uuid.New()
	for _, pair := range []struct {
		id    uuid.UUID
		score float64
	}{{s1, 70}, {s2, 85}, {s3, 55}} {
		if err := repo.RecordExamScore(ctx, examID, pair.id, pair.score); err != nil {
			t.Fatalf("record: %v", err)
		}
	}

	// s2 has the highest score → rank 1.
	rank, total, score, ok, err := repo.GetRank(ctx, examID, s2)
	if err != nil || !ok {
		t.Fatalf("GetRank s2: ok=%v err=%v", ok, err)
	}
	if rank != 1 || total != 3 || score != 85 {
		t.Errorf("s2: rank=%d total=%d score=%v, want 1/3/85", rank, total, score)
	}

	// Unknown student → ok=false.
	_, _, _, ok, err = repo.GetRank(ctx, examID, uuid.New())
	if err != nil || ok {
		t.Errorf("unknown student: ok=%v err=%v, want false/nil", ok, err)
	}
}

func TestLeaderboard_RebuildFromScores(t *testing.T) {
	c := testClient(t)
	defer c.Close()
	repo := NewScoreRankingRepo(c)
	ctx := context.Background()

	examID := uuid.New()
	defer c.Del(ctx, examScoresKey(examID))

	// Seed initial scores.
	s1, s2 := uuid.New(), uuid.New()
	_ = repo.RecordExamScore(ctx, examID, s1, 60)
	_ = repo.RecordExamScore(ctx, examID, s2, 90)

	// Rebuild with different data.
	s3 := uuid.New()
	if err := repo.RebuildFromScores(ctx, examID, map[uuid.UUID]float64{
		s3: 100,
		s1: 75,
	}); err != nil {
		t.Fatalf("RebuildFromScores: %v", err)
	}

	// s2 should no longer be present.
	_, _, _, ok, _ := repo.GetRank(ctx, examID, s2)
	if ok {
		t.Error("s2 should be gone after rebuild")
	}
	// s3 should be rank 1.
	rank, total, _, ok, err := repo.GetRank(ctx, examID, s3)
	if err != nil || !ok || rank != 1 || total != 2 {
		t.Errorf("s3 after rebuild: rank=%d total=%d ok=%v err=%v", rank, total, ok, err)
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
