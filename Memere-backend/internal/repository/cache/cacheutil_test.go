package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/cache"
)

func newTestRedis(t *testing.T) *goredis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	return goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
}

func TestGetOrSet_CacheHit(t *testing.T) {
	rdb := newTestRedis(t)
	ctx := context.Background()

	calls := 0
	load := func(_ context.Context) (string, error) {
		calls++
		return "from-db", nil
	}

	// First call: miss → loader invoked.
	v1, err := cache.GetOrSet(ctx, rdb, "test:key", "test", time.Minute, load)
	if err != nil || v1 != "from-db" || calls != 1 {
		t.Fatalf("first call: got %q err=%v calls=%d", v1, err, calls)
	}

	// Second call: hit → loader NOT invoked again.
	v2, err := cache.GetOrSet(ctx, rdb, "test:key", "test", time.Minute, load)
	if err != nil || v2 != "from-db" || calls != 1 {
		t.Fatalf("second call (should hit cache): got %q err=%v calls=%d", v2, err, calls)
	}
}

func TestGetOrSet_CacheMiss_FallsBackToLoader(t *testing.T) {
	rdb := newTestRedis(t)
	ctx := context.Background()

	calls := 0
	load := func(_ context.Context) (int, error) {
		calls++
		return 42, nil
	}

	v, err := cache.GetOrSet(ctx, rdb, "miss:key", "test", time.Minute, load)
	if err != nil || v != 42 || calls != 1 {
		t.Fatalf("miss: got %d err=%v calls=%d", v, err, calls)
	}
}

func TestGetOrSet_LoaderError_NotCached(t *testing.T) {
	rdb := newTestRedis(t)
	ctx := context.Background()

	boom := errors.New("db exploded")
	calls := 0
	load := func(_ context.Context) (string, error) {
		calls++
		return "", boom
	}

	_, err := cache.GetOrSet(ctx, rdb, "err:key", "test", time.Minute, load)
	if !errors.Is(err, boom) {
		t.Fatalf("expected loader error, got %v", err)
	}

	// A subsequent call must re-invoke the loader (error must not be cached).
	_, _ = cache.GetOrSet(ctx, rdb, "err:key", "test", time.Minute, load)
	if calls != 2 {
		t.Errorf("loader should be called again after error, calls=%d", calls)
	}
}

func TestDel_RemovesCachedValue(t *testing.T) {
	rdb := newTestRedis(t)
	ctx := context.Background()

	calls := 0
	load := func(_ context.Context) (string, error) { calls++; return "v", nil }

	_, _ = cache.GetOrSet(ctx, rdb, "del:key", "test", time.Minute, load)
	if calls != 1 {
		t.Fatal("precondition: loader not called on first get")
	}

	cache.Del(ctx, rdb, "del:key")

	_, _ = cache.GetOrSet(ctx, rdb, "del:key", "test", time.Minute, load)
	if calls != 2 {
		t.Errorf("after Del, loader must be called again; calls=%d", calls)
	}
}
