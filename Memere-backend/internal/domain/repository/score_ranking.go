package repository

import (
	"context"

	"github.com/google/uuid"
)

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
}
