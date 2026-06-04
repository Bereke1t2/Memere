# Phase 6 · Skill 5 — Kubernetes, Autoscaling & the Microservices Endgame

> **Prerequisite:** Phase 6 Skills 1–4 done (observable, hardened, performant,
> containerized + CI/CD). Read [`docs/skill.md`](../../skill.md) §2. This is the
> **final skill** of the build.
>
> **Spec references:** `memere_Design_Specification.md` §13.2 (k8s structure),
> §12.1–12.2 (scaling tiers & triggers), §3.1–3.2 (microservices target &
> service responsibilities), §13.3 (alerting), README "Kubernetes" & "Scaling
> Tiers".

---

## Goal

Deploy to Kubernetes and lay out the path from the modular monolith to the §3.2
microservices — **without a risky big-bang rewrite**. Deliver the k8s manifests,
HPA, ingress/TLS, alerting rules, and a concrete **strangler-fig extraction plan**
that uses the clean-architecture boundaries you've maintained since Phase 1.

The reward for the dependency rule we've enforced for six phases: each
`usecase/<domain>` package is already a near-self-contained service. Extraction is
mechanical, not architectural.

---

## Tasks

### 5.1 — Kubernetes manifests (`k8s/`) per §13.2

```
k8s/
├── namespaces/        production.yaml, staging.yaml
├── deployments/       api.yaml, worker.yaml            # api + transcode/notif worker
├── services/          api-svc.yaml (ClusterIP)
├── ingress/           api-ingress.yaml                 # nginx-ingress + cert-manager TLS
├── configmaps/        app-config.yaml                  # non-secret config
├── secrets/           app-secrets.yaml                 # SealedSecrets/ExternalSecrets ref
├── hpa/               api-hpa.yaml, worker-hpa.yaml
└── jobs/              migrate-job.yaml                  # run before rollout
```

Key manifest properties:

```yaml
# deployments/api.yaml (essentials)
spec:
  replicas: 2
  template:
    spec:
      securityContext: { runAsNonRoot: true, readOnlyRootFilesystem: true }
      containers:
        - name: api
          image: <registry>/memere-api:<tag>
          envFrom: [{ configMapRef: { name: app-config } }, { secretRef: { name: app-secrets } }]
          ports: [{ containerPort: 8080 }]
          readinessProbe: { httpGet: { path: /readyz, port: 8080 }, initialDelaySeconds: 5 }
          livenessProbe:  { httpGet: { path: /healthz, port: 8080 } }
          resources:
            requests: { cpu: "100m", memory: "128Mi" }
            limits:   { cpu: "500m", memory: "256Mi" }
```

- **Migrations as a pre-rollout `Job`/init container** (`/migrate up`) — never in
  the app's normal start path. Tie to the deploy (Skill 4).
- The **worker** deployment is separate from `api` (different image w/ ffmpeg,
  different scaling profile, no ingress). This is the first real "service split".
- Secrets via SealedSecrets or External Secrets (never plain manifests in git).
- Graceful shutdown: `terminationGracePeriodSeconds` long enough for in-flight
  requests + worker drain (the app already handles `SIGTERM` since Phase 1).

### 5.2 — Autoscaling (§12.1) + PodDisruptionBudgets

```yaml
# hpa/api-hpa.yaml
spec:
  scaleTargetRef: { apiVersion: apps/v1, kind: Deployment, name: api }
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource: { name: cpu, target: { type: Utilization, averageUtilization: 65 } }
```

- HPA on CPU (and, if metrics-adapter available, on the Skill-1 custom metrics:
  request rate / queue depth).
- A `PodDisruptionBudget` (minAvailable) so rollouts/node drains don't drop below
  capacity.
- Worker HPA can scale on **queue depth** (transcode backlog) — note the metric
  source.

### 5.3 — Ingress, TLS, edge

- nginx-ingress with cert-manager (Let's Encrypt) for TLS — satisfies the
  HTTPS-only Non-Negotiable at the edge (HTTP→HTTPS redirect here).
- The ingress is the real **API Gateway** (§3.2): rate limiting can move here at
  scale, but keep the app-level limiter as defense-in-depth.
- `/metrics` and `/healthz`/`/readyz` not exposed publicly (internal only / scrape
  via ServiceMonitor).

### 5.4 — Monitoring stack hookup (§13.3)

- `ServiceMonitor` (Prometheus Operator) scraping `/metrics`.
- Grafana dashboards: RED per route, DB/Redis saturation, worker queue depth,
  business KPIs (payments, enrollments, MRR, transcode backlog).
- **Alert rules** matching §13.3: 5xx rate > 1%, p99 latency > 500ms, DB
  replication lag > 30s (when replicas exist), transcode failure spike, payment
  webhook failures, worker stalled. Route to PagerDuty/OpsGenie/Slack.

### 5.5 — Scaling-tier playbook (`docs/scaling.md`) per §12.1–12.2

Document the concrete triggers and the action at each tier (1K → 10K → 100K → 1M):
- when to add read replicas (split read-only repos to a replica pool — the repo
  interface already allows a reader/writer split),
- when to move rate-limiting/caching to the edge,
- when the in-proc `JobQueue` must become **SQS/RabbitMQ** (the port is already
  there — swap the impl, run the worker as its own deployment consuming the queue),
- when to extract the first microservice.

### 5.6 — Microservices extraction plan (`docs/microservices-plan.md`) per §3.2

The strangler-fig plan, leveraging six phases of clean boundaries:

1. **Pre-work (already done):** every domain is a `usecase/<x>` package depending
   only on `domain` interfaces; the `JobQueue`, `ObjectStore`, `Notifier`,
   `PaymentProvider` are ports. No cross-domain DB joins that can't become an API
   call (audit and note any that exist).
2. **Order of extraction** (lowest-coupling first), matching §3.2 ports:
   - **Notification Service** (8086) — already async via queue; near-zero
     coupling. Extract first as the template.
   - **Media/Video worker** (transcode) — already a separate deployment; promote
     to its own service + DB schema ownership.
   - **Payment Service** (8085) — clear boundary, webhook-driven.
   - **Quiz/Exam Service** (8083/8084), **Course** (8082), **Progress** (8087),
     **Auth** (8081) last (everyone depends on it — keep it stable, expose JWKS).
3. **Mechanics per extraction:** carve the `usecase` + its repo + its schema into a
   new module; replace in-process calls from the monolith with an HTTP/gRPC client
   implementing the **same domain interface** (so the monolith doesn't know the
   difference); move the schema to the service's own DB when ready (DB-per-service).
4. **Cross-cutting:** shared `pkg/` (jwt, apperror, pagination) becomes a shared
   library or is duplicated; auth via JWT validation (JWKS from Auth service);
   service-to-service authn (mTLS / signed internal tokens); the API gateway routes
   per §3.2 port map.
5. **Data:** introduce the **saga** pattern for cross-service flows that were single
   transactions (e.g. payment→enrollment once Payment and Course are separate) —
   note where the current in-transaction fulfillment (Phase 4) becomes an
   event/saga.

This is a **plan**, not an implementation — extraction happens when §12.2 triggers
fire (a service needs independent scaling). Do not split prematurely; the monolith
is the correct architecture until the metrics say otherwise.

### 5.7 — Validation

Deploy to a local cluster (kind/minikube) or staging: manifests apply; migrate Job
runs then api+worker roll out; readiness gates traffic; HPA object is valid;
ingress serves over TLS; Prometheus scrapes; an alert rule fires in a forced-error
test. Document the deploy runbook.

---

## Definition of Done

- [ ] `kubectl apply -k k8s/` (or per-dir) brings up api + worker + migrate Job on
      a test cluster; pods reach Ready via `/readyz`.
- [ ] Migrations run as a pre-rollout Job/init container, not in app start.
- [ ] HPA (CPU + custom metric where available) and a PodDisruptionBudget exist;
      rolling update keeps the service available.
- [ ] Ingress terminates TLS (cert-manager); HTTP redirects to HTTPS;
      `/metrics`/health endpoints are not public.
- [ ] Prometheus scrapes via ServiceMonitor; Grafana dashboards + §13.3 alert
      rules are defined and a forced error triggers an alert.
- [ ] Secrets are managed (SealedSecrets/ExternalSecrets), never plain in git.
- [ ] `docs/scaling.md` (tier playbook) and `docs/microservices-plan.md`
      (strangler-fig extraction order + mechanics) are written and reviewed.

## Verification commands

```bash
kubectl apply -k k8s/ --dry-run=server     # validate manifests
kind create cluster && kubectl apply -k k8s/   # local end-to-end
kubectl get pods -n staging -w             # api+worker Ready, migrate Job Completed
kubectl rollout status deploy/api -n staging
```

---

## 🎉 Phase 6 complete — the build is DONE

All six phases are delivered:

1. **Foundation + Auth + Courses** — clean-architecture monolith, JWT/RBAC.
2. **Quiz & Exam engines** — server-side grading, server-enforced timers,
   analytics.
3. **Video pipeline** — pre-signed upload, FFmpeg→HLS, signed delivery.
4. **Payments & Enrollments** — idempotent, webhook-deduped, transactional
   fulfillment; subscriptions; revenue.
5. **Progress, Notifications, Admin, Certificates** — feature-complete product.
6. **Production-ready** — observable, hardened, performant (p95<200ms),
   containerized, CI/CD, on Kubernetes with a clear microservices roadmap.

The result honors every Non-Negotiable from `docs/skill.md` §2 and follows the
design spec end to end. Scale by following `docs/scaling.md`; split services by
following `docs/microservices-plan.md` **when the §12.2 triggers fire — not
before**.

There are no further phases to author. Future work is product iteration
(Phase 2/3 roadmap features from the spec: AI tutor, adaptive learning, live
sessions, Amharic localization, parent portal, B2B licensing) — each can be added
as a new `docs/skills/phaseN/` set using this same skill format.
