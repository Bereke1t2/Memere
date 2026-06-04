# Phase 6 · Skill 1 — Observability (Logging, Metrics, Tracing, Sentry)

> **Prerequisite:** Phases 1–5 complete and green (feature-complete monolith).
> Read [`docs/skill.md`](../../skill.md) §2. Phase 6 adds **no product features** —
> it makes the system production-ready and scalable.
>
> **Spec references:** `memere_Design_Specification.md` §13.3 (monitoring &
> alerting — Prometheus, Grafana, ELK, Sentry), §1.5 (NFRs: p95 < 200ms, 99.9%
> uptime), §3.2 (gateway).

---

## Goal

Make the system **observable**: structured leveled logging with request
correlation, Prometheus metrics (RED: rate/errors/duration + business gauges),
distributed tracing (OpenTelemetry), and error tracking (Sentry). You can't
operate or scale what you can't see — this is the foundation for the rest of
Phase 6.

---

## Tasks

### 1.1 — Structured logging (`pkg/logger`)

Replace the ad-hoc logger from Phase 1 with `log/slog` (stdlib, Go 1.21+).

```go
package logger

import ( "log/slog"; "os" )

func New(env string) *slog.Logger {
    var h slog.Handler
    if env == "production" {
        h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    } else {
        h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    }
    return slog.New(h)
}
```

- A middleware that injects a child logger with `request_id`, `method`, `path`,
  `user_id` (if authed) into the request context; handlers/usecases log via
  `logger.FromContext(ctx)`.
- **Redaction is mandatory** (Non-Negotiable #8): a helper that strips
  Authorization headers, tokens, passwords, card data, and provider secrets before
  they can be logged. Add a test that feeds known-sensitive fields and asserts
  they're absent from output.
- Replace `fmt.Print*`/ad-hoc logging across the codebase with `slog`.

### 1.2 — Prometheus metrics (`internal/infrastructure/metrics` + `/metrics`)

Use `prometheus/client_golang`.

```go
var (
    HTTPRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "http_requests_total", Help: "HTTP requests",
    }, []string{"method", "route", "status"})
    HTTPDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name: "http_request_duration_seconds",
        Buckets: prometheus.DefBuckets, // tune for p95<200ms SLO
    }, []string{"method", "route"})
)
```

- A metrics middleware records count + duration per **route template** (not raw
  path — avoid cardinality explosion; use the Gin route, not `:id` values).
- Business gauges/counters: payments completed/failed, enrollments, video
  transcode success/failure, attempts graded, notifications sent, worker queue
  depth. Expose worker health.
- DB pool stats (pgx) and Redis stats as gauges.
- Expose `GET /metrics` (no auth, but bind to an internal port or protect at the
  gateway/k8s level — document).

### 1.3 — Distributed tracing (OpenTelemetry)

- `go.opentelemetry.io/otel` with an OTLP exporter (configurable endpoint;
  no-op/stdout exporter in dev). A tracing middleware starts a span per request
  (carrying `request_id`/trace correlation), and the pgx + redis + outbound HTTP
  (provider) calls are instrumented so a slow query/provider call is visible in a
  trace.
- Propagate context through usecases → repos (already `ctx`-threaded since Phase
  1, so this is mostly wiring instrumentation wrappers).

### 1.4 — Sentry error tracking

- `getsentry/sentry-go`; init from `SENTRY_DSN` (disabled if empty).
- A recovery integration: panics and `500`-class `apperror.Internal` are reported
  to Sentry **with** the `request_id`/trace id and **without** PII/secrets (reuse
  the redactor). 4xx client errors are **not** reported (they're not bugs).
- Capture worker failures (transcode, payment fulfillment, notification send) too.

### 1.5 — Config + wiring

Add `SENTRY_DSN`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `METRICS_PORT`,
`LOG_LEVEL` to config + `.env.example`. Wire logger/metrics/tracing/Sentry in
`main.go` early (before building the router) so everything downstream is
instrumented. Order the global middleware:
`recover(+sentry) → request-id → tracing → logger → metrics → cors → ratelimit`.

### 1.6 — Health & readiness split

Split Phase 1's `/health`:
- `GET /healthz` (liveness) — process is up; no dependency checks (cheap).
- `GET /readyz` (readiness) — pings DB + Redis (+ object store HEAD optionally);
  `503` until dependencies are reachable. k8s uses these (Skill 5).

### 1.7 — Tests

Redactor strips all sensitive fields; metrics middleware uses route templates (no
high-cardinality labels); `/metrics` exposes the registered collectors; tracing
no-op exporter works without a collector; Sentry disabled when DSN empty; readiness
returns 503 when DB down.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] All logging goes through `slog`; `grep -rn 'fmt.Print' internal/ cmd/` is
      empty.
- [ ] Sensitive fields never appear in logs or Sentry (redaction test green).
- [ ] `GET /metrics` exposes RED metrics with **route-template** labels + business
      metrics; DB/Redis/worker stats present.
- [ ] Requests produce traces with correlated `request_id`; pgx/redis/provider
      calls are spans.
- [ ] Sentry reports 5xx/panics (with correlation, no PII) and ignores 4xx; no-ops
      cleanly when `SENTRY_DSN` empty.
- [ ] `/healthz` (liveness) and `/readyz` (readiness, dep-checked) both work.

## Verification commands

```bash
go build ./... && golangci-lint run && go test ./...
grep -rn 'fmt.Print' internal/ cmd/ && echo "FAIL: use slog" || echo OK
make run & sleep 2
curl -s localhost:${METRICS_PORT:-8080}/metrics | grep -E 'http_requests_total|http_request_duration_seconds'
curl -i localhost:8080/readyz
```

## Hand-off to Skill 2

The system is observable. Skill 2 hardens **security** (gateway concerns: security
headers, tightened rate limits, input/size limits, secret management, dependency
scanning) now that you can measure the impact of every change.
