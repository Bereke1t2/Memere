package repository

import (
	"context"

	"github.com/google/uuid"
)

// LeaderboardEntry is one row in a top-N leaderboard view.
type LeaderboardEntry struct {
	StudentID uuid.UUID
	Score     float64 // best percentage (0–100)
	Rank      int64   // 1-based position (highest score = 1)
}

// ScoreRanking is a derived, Redis-backed ranking of exam scores used for the
// §9.3 percentile metric ("student score vs all attempts for the same exam").
// It is a cache, not a source of truth: courses.exam_attempts remains
// authoritative and the ranking is rebuildable from it.
type ScoreRanking interface {
	// RecordExamScore records a student's best percentage for an exam. Only the
	// student's highest score is kept (a later, lower score does not regress it).
	RecordExamScore(ctx context.Context, examID, studentID uuid.UUID, percentage float64) error
	// PercentileRank returns the student's percentile (0–100) among all recorded
	// scores for the exam: the share of students scoring at or below them. Returns
	// (0, false) when the student has no recorded score.
	PercentileRank(ctx context.Context, examID, studentID uuid.UUID) (percentile float64, ok bool, err error)

	// --- Phase 6 Skill 3: leaderboard (§FR-11, §9.3) ---

	// GetTopN returns the top limit students ordered by score descending.
	GetTopN(ctx context.Context, examID uuid.UUID, limit int) ([]LeaderboardEntry, error)
	// GetRank returns the caller's 1-based rank (highest = 1), score, and total
	// entries. ok=false when the student has no recorded score.
	GetRank(ctx context.Context, examID, studentID uuid.UUID) (rank, total int64, score float64, ok bool, err error)
	// RebuildFromScores repopulates the sorted set from a map of studentID →
	// best percentage. Used to recover from Redis loss.
	RebuildFromScores(ctx context.Context, examID uuid.UUID, scores map[uuid.UUID]float64) error
}
