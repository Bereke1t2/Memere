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

// GetTopN returns the top limit students ordered by score descending (highest
// first), with 1-based ranks. Ties share the same score but are assigned
// consecutive ranks per Redis ordering.
func (r *ScoreRankingRepo) GetTopN(ctx context.Context, examID uuid.UUID, limit int) ([]repository.LeaderboardEntry, error) {
	key := examScoresKey(examID)
	members, err := r.client.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, apperror.Internal(err)
	}
	out := make([]repository.LeaderboardEntry, 0, len(members))
	for i, m := range members {
		id, err := uuid.Parse(m.Member.(string))
		if err != nil {
			continue
		}
		out = append(out, repository.LeaderboardEntry{
			StudentID: id,
			Score:     m.Score,
			Rank:      int64(i + 1),
		})
	}
	return out, nil
}

// GetRank returns the 1-based rank of the student (highest score = rank 1), the
// student's best score, and the total number of entries in the sorted set.
// ok=false when the student has no recorded score.
func (r *ScoreRankingRepo) GetRank(ctx context.Context, examID, studentID uuid.UUID) (rank, total int64, score float64, ok bool, err error) {
	key := examScoresKey(examID)
	member := studentID.String()

	// ZREVRANK returns 0-based rank (best = 0) in descending order.
	rank0, zErr := r.client.ZRevRank(ctx, key, member).Result()
	if zErr == goredis.Nil {
		return 0, 0, 0, false, nil
	}
	if zErr != nil {
		return 0, 0, 0, false, apperror.Internal(zErr)
	}

	score, zErr = r.client.ZScore(ctx, key, member).Result()
	if zErr != nil {
		return 0, 0, 0, false, apperror.Internal(zErr)
	}

	total, zErr = r.client.ZCard(ctx, key).Result()
	if zErr != nil {
		return 0, 0, 0, false, apperror.Internal(zErr)
	}

	return rank0 + 1, total, score, true, nil
}

// RebuildFromScores clears the sorted set and repopulates it from the provided
// studentID → best-percentage map. Intended for admin-triggered recovery after
// Redis data loss.
func (r *ScoreRankingRepo) RebuildFromScores(ctx context.Context, examID uuid.UUID, scores map[uuid.UUID]float64) error {
	key := examScoresKey(examID)

	pipe := r.client.Pipeline()
	pipe.Del(ctx, key)
	for sid, pct := range scores {
		pipe.ZAdd(ctx, key, goredis.Z{Score: pct, Member: sid.String()})
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
