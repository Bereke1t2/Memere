# Phase 6 · Skill 3 — Performance, Caching & Scalability

> **Prerequisite:** Phase 6 Skills 1–2 done (you can now measure latency and have
> hardened the surface). Read [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §12 (scaling tiers &
> decision triggers), §1.5 (NFRs: p95 < 200ms, 1K→10K concurrent), §9.3
> (leaderboard via Redis sorted sets), §3.3 (Redis cache-aside).

---

## Goal

Hit the §1.5 performance SLOs and prepare for the §12 scaling path: add
**cache-aside** caching for hot read paths, tune queries/indexes, finish the
**leaderboard**, size connection pools, and **load-test** to prove p95 < 200ms.
Measurement (Skill 1) drives every change here — optimize what the metrics/traces
say is slow, not what you guess.

---

## Tasks

### 3.1 — Cache-aside layer (`internal/repository/cache`)

Wrap hot, read-heavy, infrequently-changing reads with Redis cache-aside. Good
candidates (per §12.2 "introduce Redis when DB CPU > 60% / p95 > 100ms"):

- published course list + course detail (high read volume),
- subscription-plan list, coupon lookups,
- user-by-id (auth hot path),
- leaderboard/percentile (already Redis).

```go
// generic cache-aside helper
func GetOrSet[T any](ctx context.Context, rdb *redis.Client, key string, ttl time.Duration,
    load func(ctx context.Context) (T, error)) (T, error) {
    if b, err := rdb.Get(ctx, key).Bytes(); err == nil {
        var v T; if json.Unmarshal(b, &v) == nil { return v, nil }
    }
    v, err := load(ctx)
    if err != nil { return v, err }
    if b, e := json.Marshal(v); e == nil { _ = rdb.Set(ctx, key, b, ttl).Err() }
    return v, nil
}
```

- **Invalidation:** on course update/publish/delete, `DEL` the relevant keys
  (write-through invalidation). Prefer short TTLs (e.g. 60s) + explicit
  invalidation over long TTLs. Document each cached key's TTL + invalidation
  trigger.
- **Never cache** authz decisions, attempt state, payment status, or anything
  per-request-sensitive. Cache data, not access decisions.
- Add cache hit/miss metrics (Skill 1) per key class.

### 3.2 — Query & index tuning

- Use the traces/`pg_stat_statements` to find the slow queries; `EXPLAIN ANALYZE`
  them. Add missing indexes (the spec listed many in §4 — verify they exist and
  are used). Watch for N+1 in the nested course/section/lesson and dashboard reads
  — batch with a single query or `IN (...)`.
- Add composite/partial indexes for the hottest filters (e.g.
  `courses (is_published, subject, grade)` for the catalog; `payments (status,
  paid_at)` for reconciliation/revenue).
- Consider **materialized views** for the heaviest admin analytics (overview/MRR)
  with periodic refresh (a worker) — only if metrics show the live query is too
  slow. Document the trade-off (staleness vs cost).

### 3.3 — Leaderboard (finish §9.3 / §FR-11)

Build the student-facing leaderboard on Redis sorted sets (the percentile
machinery exists from Phase 2 Skill 4):

- `leaderboard:course:{course_id}` and/or `leaderboard:exam:{exam_id}` ZSETs,
  score = best percentage. Updated on attempt grade.
- `GetLeaderboard(ctx, scope, limit)` → top N + the caller's own rank (`ZREVRANK`).
- Rebuildable from PG (a `RebuildLeaderboard` admin/worker task) so Redis loss is
  recoverable. Add a usecase + (Skill 5 already exposes analytics — add the route
  there or here-noted for the API skill).

### 3.4 — Connection pool & resource tuning

- Size the pgx pool from config based on expected concurrency and DB max
  connections (§12: t3.micro → small pool; scale with tier). Expose pool
  saturation metric; alert when near max.
- Redis pool sizing; sensible dial/read/write timeouts.
- Worker concurrency caps (transcode is CPU-heavy — keep it low and separate from
  request handling; note that at scale the transcode worker should be its **own
  deployment**, foreshadowing Skill 4/5).
- `GOMAXPROCS` correctness in containers (use `automaxprocs`).

### 3.5 — Pagination & payload hygiene

- Confirm all list endpoints are cursor-paginated with clamped limits (Phase 1
  pattern) — no unbounded scans. Add any that slipped.
- Minimize payloads (§14.4 risk "minimal API payloads" for poor connectivity):
  ensure list DTOs are lean (no heavy nested blobs); offer field selection only if
  metrics justify it.
- Enable gzip/deflate response compression middleware for JSON.

### 3.6 — Load testing (`scripts/load/`)

- A `k6` (or `vegeta`) script exercising the hot paths: login, course list/detail,
  start+submit quiz, stream URL, dashboard. Parameterize concurrency to simulate
  the §1.5 targets (1K concurrent now, headroom toward 10K).
- Run against a realistic seed (extend `scripts/seed.go`). Capture p50/p95/p99 from
  the metrics; **assert p95 < 200ms** on read paths under target load.
- Document results + any tier recommendation (§12.1) in a short `docs/perf.md`.

### 3.7 — Tests

Cache-aside returns cached value + invalidation `DEL`s on write; cache miss falls
back to loader; leaderboard top-N + own-rank correct and rebuildable; pagination
limits clamped; no list endpoint returns unbounded results.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] Hot read paths are cache-aside with documented TTL + invalidation; cache
      hit/miss metrics exposed; no sensitive data cached.
- [ ] Slow queries identified via traces and fixed; required indexes verified
      present and used (EXPLAIN); no N+1 on course-detail/dashboard.
- [ ] Leaderboard (top-N + own rank) works on Redis ZSETs and is rebuildable.
- [ ] pgx/Redis pools + worker concurrency sized from config; saturation metrics
      exist; `automaxprocs` in place.
- [ ] Load test runs; **p95 < 200ms** on read paths at the 1K-concurrent target;
      results recorded in `docs/perf.md`.

## Verification commands

```bash
go build ./... && golangci-lint run && go test ./...
make run & sleep 2
k6 run scripts/load/hot_paths.js            # or the vegeta equivalent
# inspect p95:
curl -s localhost:${METRICS_PORT:-8080}/metrics | grep http_request_duration_seconds
```

## Hand-off to Skill 4

Performance meets SLO and caching is in place. Skill 4 builds the **CI/CD pipeline
and containerization** (multi-stage images, GitHub Actions per §13.1, migrations
in the pipeline) so deploys are automated and repeatable.
