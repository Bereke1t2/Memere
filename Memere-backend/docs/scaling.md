# Memere — Scaling Tier Playbook

> **Reference:** Design spec §12.1–12.2. Read alongside `docs/microservices-plan.md`.
> **Rule:** Do not act on a tier until the metrics in that tier's trigger column fire
> consistently for ≥ 24 hours. Premature splitting trades simplicity for complexity
> with no throughput gain.

---

## Tier overview

| Tier | MAU | RPS (peak) | Architecture |
|------|-----|-----------|--------------|
| 0 – Launch | < 1 K | < 5 | Monolith, single DB, single Redis node |
| 1 – Growth | 1 K – 10 K | 5 – 50 | Monolith + read replica, Redis Sentinel, HPA |
| 2 – Scale | 10 K – 100 K | 50 – 500 | First service extractions, connection pooler, edge cache |
| 3 – Hyper-scale | 100 K – 1 M | 500 – 5 K | Full microservices, DB-per-service, async sagas |

---

## Tier 0 — Launch (< 1 K MAU)

**Current state after Phase 6.** No action required.

- Single PostgreSQL instance (docker-compose / managed RDS `db.t3.medium`)
- Single Redis node (Elasticache `cache.t3.micro`)
- 2 API pods, 1 worker pod (HPA floor)
- CI/CD pipeline ships to production via `kubectl apply -k k8s/`

**Watch these metrics:** p99 latency, DB CPU, Redis memory. They are your early-warning system.

---

## Tier 1 — Growth (1 K – 10 K MAU)

### Triggers (any one fires → act)

- DB CPU sustained > 70%
- p99 API latency > 300ms for > 30 min
- DB connection pool acquired > 80% (`memere_db_pool_acquired / memere_db_pool_total`)
- Redis memory > 60% of instance size

### Actions

#### 1. Add a PostgreSQL read replica

The `domain/repository` interfaces already have a reader/writer split at the port level (both accept a `*pgxpool.Pool` today — split them to `RWPool` + `ROPool`). Wire course catalog, quiz question, exam, and progress reads to the replica pool:

```go
// internal/repository/postgres/pool.go  (add)
type Pools struct {
    RW *pgxpool.Pool
    RO *pgxpool.Pool  // nil falls back to RW
}
```

Target read-heavy repos first: `CourseRepo.List`, `QuizRepo.Get`, `ExamRepo.ListMockExams`, `ProgressRepo.GetDashboard`.

**When to add replica:** DB CPU > 70% OR p99 > 300ms AND query profiling shows reads dominating.

#### 2. Add PgBouncer (connection pooler)

At ~500 concurrent users, pgx's 25-connection pool cap becomes a bottleneck. Deploy PgBouncer in transaction mode in front of PostgreSQL. No application changes required — swap `DB_URL` to point at PgBouncer.

#### 3. Redis Sentinel → Redis Cluster

Switch Elasticache from single-node to a 3-node Sentinel setup for HA. No application code change — the `go-redis` client handles failover transparently via the Sentinel URL format.

#### 4. HPA already covers burst

API HPA scales to 10 pods on CPU/memory pressure. No manifest changes needed at Tier 1. Monitor HPA events (`kubectl describe hpa api -n production`) to ensure it is firing.

#### 5. Raise DB connection limits

```yaml
# k8s/configmaps/app-config.yaml
DB_MAX_CONNS: "50"   # was 25
```

Redeploy api pods (rolling update).

---

## Tier 2 — Scale (10 K – 100 K MAU)

### Triggers

- p99 latency > 200ms sustained (after Tier 1 tuning)
- Worker queue depth sustained > 20 jobs
- Single service accounts for > 30% of DB query time
- Notification fan-out latency > 5s

### Actions

#### 1. Move rate limiting and caching to the edge

- Enable nginx-ingress rate limiting (`nginx.ingress.kubernetes.io/limit-rps`) as the primary limiter. Keep the app-level Redis limiter as defense-in-depth only.
- Add a CDN (CloudFront) for course catalog and static lesson metadata (`GET /courses`, `GET /courses/:id`). Cache-Control headers already set by the API; the CDN caches transparently.
- Course catalog cache TTL: 60 s. Invalidate on publish/update via CloudFront cache invalidation API call from the `course` usecase.

#### 2. Extract the Notification Service (port 8086)

This is the first and safest extraction because:
- Already fully async via the in-proc `JobQueue`
- Zero synchronous callers (fire-and-forget hooks)
- Has its own DB schema (`notifications`, `device_tokens`, `notification_preferences`)

See `docs/microservices-plan.md` §Extraction-1 for mechanics.

**Trigger:** Notification fan-out takes > 5 s for large broadcasts, OR notification load measurably affects API p99.

#### 3. Swap in-proc queue for SQS / RabbitMQ

The `JobQueue` port (`internal/domain/service/queue.go`) is already a swappable interface. The `messaging.InProcQueue` is non-durable — crashes lose queued transcodes. At Tier 2:

1. Deploy RabbitMQ (or use AWS SQS).
2. Write `messaging.SQSQueue` / `messaging.AMQPQueue` implementing the same `JobQueue` interface.
3. Swap the wire in `cmd/api/main.go` — no other code changes.
4. The worker `Deployment` then scales independently on queue depth (HPA custom metric already defined in `k8s/hpa/worker-hpa.yaml` comments).

**Trigger:** > 2 transcode job losses per week due to pod restarts OR queue depth consistently > 20.

#### 4. Read replica for heavy analytics queries

Admin revenue and exam analytics queries (`/admin/analytics/overview`, `/exam-attempts/:id/analytics`) are potentially expensive. Route them to a dedicated analytics read replica with a longer connection pool.

---

## Tier 3 — Hyper-scale (100 K – 1 M MAU)

### Triggers

- Any single service's DB schema is a bottleneck for independent team deploys
- A microservice needs a different scaling profile than the monolith
- Zero-downtime deploys take > 5 min due to migration scope

### Actions

#### 1. Full strangler-fig extraction per `docs/microservices-plan.md`

Follow the extraction order: Notification → Video/Worker → Payment → Quiz/Exam → Course → Auth.
Each extraction reduces the monolith's blast radius and allows independent scaling.

#### 2. DB-per-service

Move each extracted service's schema to its own RDS instance. The schema-per-domain pattern already enforced (§4.3) makes this a `pg_dump --schema=notifications | psql new_db` operation with a cut-over migration window.

#### 3. Cross-service transactions → Saga pattern

Flows that were single DB transactions become compensating sagas. The first saga to implement is **payment→enrollment**:
- Payment Service emits `payment.confirmed` event to SQS
- Course Service consumes event and creates enrollment
- If Course Service fails, Payment Service receives `enrollment.failed` and issues a refund

The `hooks` pattern established in Phase 5 is the local prototype of this; sagas are the distributed equivalent.

#### 4. API Gateway (dedicated)

Replace the nginx-ingress with a proper API gateway (Kong / AWS API Gateway) that handles:
- JWT validation centrally (JWKS endpoint from Auth Service)
- Per-service route mapping matching the §3.2 port table
- Service-to-service mTLS

The app-level auth middleware stays in each service as defense-in-depth.

#### 5. CQRS for leaderboards and analytics

The exam leaderboard (`scoreRankingRepo` → Redis sorted set) is already a CQRS read model. At Tier 3, extend this pattern to user dashboards and revenue reports: write-side publishes events; read-side maintains pre-aggregated views in Redis or a separate read DB.

---

## Metric thresholds summary

| Metric | Tier-1 action | Tier-2 action | Tier-3 action |
|--------|--------------|--------------|--------------|
| DB CPU % | > 70% → read replica | > 80% → extract service | > 90% → DB-per-service |
| p99 latency ms | > 300 → pool tune | > 200 → edge cache | > 150 → CQRS |
| Queue depth (jobs) | > 5 → monitor | > 20 → durable queue | > 100 → dedicated service |
| API pod count (HPA max) | ≤ 5 normal | 5–8 → review bottleneck | hitting 10 → extract |
| DB conn pool % | > 80% → PgBouncer | > 90% → shard | n/a |

---

## Runbook: emergency scale-up

```bash
# Manual scale API pods immediately (before HPA catches up)
kubectl scale deployment api --replicas=6 -n production

# Tail live logs from all API pods
kubectl logs -l app.kubernetes.io/component=api -n production -f --max-log-requests=10

# Check HPA status
kubectl describe hpa api -n production

# Check worker queue depth (exposed via /metrics)
kubectl port-forward svc/api 9090:9090 -n production &
curl -s localhost:9090/metrics | grep memere_transcode_queue_depth
```
