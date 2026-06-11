package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// VideoRepo is the sqlc-backed implementation of repository.VideoRepository.
type VideoRepo struct {
	q *sqlcgen.Queries
}

var _ repository.VideoRepository = (*VideoRepo)(nil)

// NewVideoRepo builds a VideoRepo over a pgx pool.
func NewVideoRepo(pool *pgxpool.Pool) *VideoRepo {
	return &VideoRepo{q: sqlcgen.New(pool)}
}

func (r *VideoRepo) Create(ctx context.Context, v *entity.Video) error {
	status := v.Status
	if status == "" {
		status = entity.VideoPending
	}
	row, err := queriesFor(ctx, r.q).CreateVideo(ctx, sqlcgen.CreateVideoParams{
		LessonID:         toPgUUID(v.LessonID),
		CourseID:         toPgUUID(v.CourseID),
		ProcessingStatus: string(status),
		OriginalFileKey:  v.OriginalFileKey,
		FileSizeBytes:    v.FileSizeBytes,
	})
	if err != nil {
		// The partial unique index videos_lesson_id_uniq enforces one live
		// video per lesson; a violation means the lesson already has one.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return apperror.Conflict("VIDEO_EXISTS", err)
		}
		return apperror.Internal(err)
	}
	*v = *videoFromRow(row)
	return nil
}

func (r *VideoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Video, error) {
	row, err := queriesFor(ctx, r.q).GetVideoByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, mapVideoErr(err)
	}
	return videoFromRow(row), nil
}

func (r *VideoRepo) GetByLessonID(ctx context.Context, lessonID uuid.UUID) (*entity.Video, error) {
	row, err := queriesFor(ctx, r.q).GetVideoByLessonID(ctx, toPgUUID(lessonID))
	if err != nil {
		return nil, mapVideoErr(err)
	}
	return videoFromRow(row), nil
}

func (r *VideoRepo) SetReady(ctx context.Context, v *entity.Video) error {
	row, err := queriesFor(ctx, r.q).SetVideoReady(ctx, sqlcgen.SetVideoReadyParams{
		ID:                 toPgUUID(v.ID),
		HlsMasterKey:       v.HLSMasterKey,
		Resolution480pKey:  v.Res480pKey,
		Resolution720pKey:  v.Res720pKey,
		Resolution1080pKey: v.Res1080pKey,
		ThumbnailKey:       v.ThumbnailKey,
		DurationSeconds:    int32(v.DurationSeconds),
	})
	if err != nil {
		return mapVideoErr(err)
	}
	*v = *videoFromRow(row)
	return nil
}

// UpdateStatus runs the guarded transition. A 0-row result means another writer
// already moved the row — reported as ok=false, not an error.
func (r *VideoRepo) UpdateStatus(ctx context.Context, id uuid.UUID, from, to entity.VideoStatus) (bool, error) {
	n, err := queriesFor(ctx, r.q).UpdateVideoStatusGuarded(ctx, sqlcgen.UpdateVideoStatusGuardedParams{
		ID:         toPgUUID(id),
		FromStatus: string(from),
		ToStatus:   string(to),
	})
	if err != nil {
		return false, apperror.Internal(err)
	}
	return n > 0, nil
}

func (r *VideoRepo) ListByStatus(ctx context.Context, status entity.VideoStatus, limit int) ([]*entity.Video, error) {
	rows, err := queriesFor(ctx, r.q).ListVideosByStatus(ctx, sqlcgen.ListVideosByStatusParams{
		ProcessingStatus: string(status),
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	videos := make([]*entity.Video, len(rows))
	for i, row := range rows {
		videos[i] = videoFromRow(row)
	}
	return videos, nil
}

func (r *VideoRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := queriesFor(ctx, r.q).SoftDeleteVideo(ctx, toPgUUID(id)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func mapVideoErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("video not found", err)
	}
	return apperror.Internal(err)
}

func videoFromRow(row sqlcgen.CoursesVideo) *entity.Video {
	return &entity.Video{
		ID:              fromPgUUID(row.ID),
		LessonID:        fromPgUUID(row.LessonID),
		CourseID:        fromPgUUID(row.CourseID),
		Status:          entity.VideoStatus(row.ProcessingStatus),
		OriginalFileKey: row.OriginalFileKey,
		HLSMasterKey:    row.HlsMasterKey,
		Res480pKey:      row.Resolution480pKey,
		Res720pKey:      row.Resolution720pKey,
		Res1080pKey:     row.Resolution1080pKey,
		ThumbnailKey:    row.ThumbnailKey,
		DurationSeconds: int(row.DurationSeconds),
		FileSizeBytes:   row.FileSizeBytes,
		ProcessingError: row.ProcessingError,
		ProcessedAt:     fromPgTimestamptz(row.ProcessedAt),
		CreatedAt:       fromPgTimestamptzValue(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptzValue(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}
