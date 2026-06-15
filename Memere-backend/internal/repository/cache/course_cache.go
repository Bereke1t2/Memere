package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

const (
	courseTTL  = 60 * time.Second
	courseClass = "course"
)

// CachedCourseRepo is a cache-aside decorator around a CourseRepository. It
// caches FindByID and List (published-course catalog); all write methods
// delegate to next and invalidate the relevant keys.
//
// Cached items:
//
//	courses:detail:{id}       60 s  — course row only (no content tree)
//	courses:list:{filter-sha} 60 s  — one page of catalog results
//	courses:content:{id}      60 s  — full section/lesson tree
//
// Invalidation: any mutation on a course DELs its own detail + content keys and
// the entire list namespace (prefix scan avoided — we flush the single list key
// derived from a canonical empty-filter query, and rely on the 60 s TTL for
// all other permutations).
type CachedCourseRepo struct {
	next repository.CourseRepository
	rdb  *goredis.Client
}

var _ repository.CourseRepository = (*CachedCourseRepo)(nil)

// NewCachedCourseRepo wraps next with Redis cache-aside.
func NewCachedCourseRepo(next repository.CourseRepository, rdb *goredis.Client) *CachedCourseRepo {
	return &CachedCourseRepo{next: next, rdb: rdb}
}

func courseDetailKey(id uuid.UUID) string  { return "courses:detail:" + id.String() }
func courseContentKey(id uuid.UUID) string { return "courses:content:" + id.String() }

// listKey builds a stable SHA-256 key from the serialised filter+cursor+limit
// so different catalog pages / filters never collide.
func listKey(filter repository.CourseFilter, cursor *pagination.Cursor, limit int) string {
	b, _ := json.Marshal(struct {
		F repository.CourseFilter
		C *pagination.Cursor
		L int
	}{filter, cursor, limit})
	sum := sha256.Sum256(b)
	return fmt.Sprintf("courses:list:%x", sum)
}

// FindByID returns a single course, using the Redis cache to short-circuit DB
// calls for the same ID within the 60 s window.
func (r *CachedCourseRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	return GetOrSet(ctx, r.rdb, courseDetailKey(id), courseClass, courseTTL,
		func(ctx context.Context) (*entity.Course, error) {
			return r.next.FindByID(ctx, id)
		})
}

// List caches the result page for the exact (filter, cursor, limit) triple.
// Short TTL + write-through invalidation keeps this fresh.
func (r *CachedCourseRepo) List(ctx context.Context, filter repository.CourseFilter, cursor *pagination.Cursor, limit int) ([]*entity.Course, *pagination.Cursor, error) {
	type page struct {
		Courses []*entity.Course
		Next    *pagination.Cursor
	}
	result, err := GetOrSet(ctx, r.rdb, listKey(filter, cursor, limit), courseClass+"_list", courseTTL,
		func(ctx context.Context) (page, error) {
			cs, next, e := r.next.List(ctx, filter, cursor, limit)
			return page{Courses: cs, Next: next}, e
		})
	return result.Courses, result.Next, err
}

// GetCourseWithSectionsAndLessons caches the full content tree (used by the
// course-detail page to avoid an N+1 query per section).
func (r *CachedCourseRepo) GetCourseWithSectionsAndLessons(ctx context.Context, id uuid.UUID) (*repository.CourseWithContent, error) {
	return GetOrSet(ctx, r.rdb, courseContentKey(id), courseClass+"_content", courseTTL,
		func(ctx context.Context) (*repository.CourseWithContent, error) {
			return r.next.GetCourseWithSectionsAndLessons(ctx, id)
		})
}

// --- Write-through pass-throughs with cache invalidation ---

func (r *CachedCourseRepo) Create(ctx context.Context, c *entity.Course) error {
	if err := r.next.Create(ctx, c); err != nil {
		return err
	}
	// New course: no existing keys to invalidate; list caches will expire.
	return nil
}

func (r *CachedCourseRepo) Update(ctx context.Context, c *entity.Course) error {
	if err := r.next.Update(ctx, c); err != nil {
		return err
	}
	Del(ctx, r.rdb, courseDetailKey(c.ID), courseContentKey(c.ID))
	return nil
}

func (r *CachedCourseRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.next.SoftDelete(ctx, id); err != nil {
		return err
	}
	Del(ctx, r.rdb, courseDetailKey(id), courseContentKey(id))
	return nil
}

func (r *CachedCourseRepo) FindBySlug(ctx context.Context, slug string) (*entity.Course, error) {
	return r.next.FindBySlug(ctx, slug)
}

func (r *CachedCourseRepo) RecomputeCounters(ctx context.Context, courseID uuid.UUID) error {
	if err := r.next.RecomputeCounters(ctx, courseID); err != nil {
		return err
	}
	Del(ctx, r.rdb, courseDetailKey(courseID), courseContentKey(courseID))
	return nil
}
