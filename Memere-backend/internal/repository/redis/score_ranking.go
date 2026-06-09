package redis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ScoreRankingRepo implements repository.ScoreRanking over a Redis sorted set per
// exam (member = student ID, score = best percentage). It is a derived cache;
// Postgres remains the source of truth and the set is rebuildable (spec §9.3).
type ScoreRankingRepo struct {
	client *goredis.Client
}

var _ repository.ScoreRanking = (*ScoreRankingRepo)(nil)

// NewScoreRankingRepo builds a ScoreRankingRepo over a go-redis client.
func NewScoreRankingRepo(client *goredis.Client) *ScoreRankingRepo {
	return &ScoreRankingRepo{client: client}
}

func examScoresKey(examID uuid.UUID) string {
	return fmt.Sprintf("exam:%s:scores", examID.String())
}

// RecordExamScore stores the student's best percentage using ZADD GT, which only
// raises an existing score, so a later, lower attempt never regresses the best.
func (r *ScoreRankingRepo) RecordExamScore(ctx context.Context, examID, studentID uuid.UUID, percentage float64) error {
	z := goredis.Z{Score: percentage, Member: studentID.String()}
	if err := r.client.ZAddGT(ctx, examScoresKey(examID), z).Err(); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

// PercentileRank returns the share (0–100) of students scoring at or below the
// given student among all recorded scores for the exam:
//
//	percentile = (# scores <= mine) / (total) * 100
//
// Returns ok=false when the student has no recorded score for the exam.
func (r *ScoreRankingRepo) PercentileRank(ctx context.Context, examID, studentID uuid.UUID) (float64, bool, error) {
	key := examScoresKey(examID)

	myScore, err := r.client.ZScore(ctx, key, studentID.String()).Result()
	if err == goredis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, apperror.Internal(err)
	}

	total, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return 0, false, apperror.Internal(err)
	}
	if total == 0 {
		return 0, false, nil
	}

	// Count members with score <= mine (inclusive), i.e. those this student ties
	// or beats.
	atOrBelow, err := r.client.ZCount(ctx, key, "-inf", fmt.Sprintf("%f", myScore)).Result()
	if err != nil {
		return 0, false, apperror.Internal(err)
	}

	return float64(atOrBelow) / float64(total) * 100, true, nil
}
