# Cloud Run deployment

Reference scripts for deploying the Memere backend to **Google Cloud Run** with
**Neon** (PostgreSQL) and **Upstash** (Redis). Full step-by-step instructions —
including provisioning the data stores and creating secrets — are in the root
[`README.md`](../../README.md#google-cloud-run--upstash-production).

## Topology (important)

The app uses an **in-process, non-durable job queue** for video transcoding and
push/email delivery. A job is produced and consumed in the *same* process, so it
does not cross Cloud Run instances. This drives the deployment shape:

- **One always-on service** built from `Dockerfile.worker` (it bundles FFmpeg),
  run with `--min-instances=1 --no-cpu-throttling` so the in-process transcode
  and notification workers stay alive. Pure scale-to-zero would silently drop
  those jobs.
- **Periodic sweeps** run via **Cloud Scheduler** hitting the authenticated
  `/internal/jobs/*` endpoints, with the in-process tickers disabled — so each
  sweep fires once per tick no matter how many instances are up.
- **In-app** notifications persist to Postgres synchronously and are unaffected.

## Scripts

| Script         | What it does                                                                 |
|----------------|------------------------------------------------------------------------------|
| `deploy.sh`    | Builds the image from `Dockerfile.worker`, pushes to Artifact Registry, and `gcloud run deploy`s the always-on service with the right flags, env vars, and Secret Manager wiring. |
| `scheduler.sh` | Creates/updates the three Cloud Scheduler jobs that POST to `/internal/jobs/*` with the `X-Internal-Job-Token` header. |

Both are parameterized by environment variables and contain **no secrets** —
secrets live only in Secret Manager and are referenced by name. Read the header
comment in each script for required/optional variables.

```bash
export PROJECT_ID=your-gcp-project REGION=europe-west1
export DB_HOST=ep-xxx-pooler.REGION.aws.neon.tech DB_USER=... DB_NAME=...
export CORS_ALLOWED_ORIGINS=https://app.memere.et APP_PUBLIC_URL=https://api.memere.et

./deploy/cloudrun/deploy.sh      # build + deploy
./deploy/cloudrun/scheduler.sh   # wire up the periodic sweeps
```

## Migrations

`Dockerfile.worker` does not include the `migrate` binary. Run migrations
against Neon from your machine or CI (they read the same `DB_*` env):

```bash
DB_HOST=... DB_USER=... DB_PASSWORD=... DB_NAME=... DB_SSL_MODE=require DB_POOLED=true \
  make migrate-up
```
