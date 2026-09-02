#!/usr/bin/env bash
#
# deploy.sh — build the Memere API image (with FFmpeg) and deploy it to Cloud Run
# as a SINGLE always-on service.
#
# Why a single always-on service (not scale-to-zero): the transcode and
# push/email pipelines use an in-process, non-durable queue — a job is enqueued
# and consumed in the SAME process, so those workers must run inside the API
# service, and the instance needs CPU allocated at all times (background
# goroutines don't run while a scale-to-zero instance is throttled between
# requests). Hence --min-instances=1 --no-cpu-throttling. See README
# "Google Cloud Run + Upstash".
#
# Secrets are NEVER passed on the command line or baked into the image. Create
# them once in Secret Manager (see README §3) and this script wires them in with
# --set-secrets. Non-secret settings come from the environment below.
#
# Usage:
#   export PROJECT_ID=... REGION=... [SERVICE=...] [REPO=...] [IMAGE_TAG=...]
#   export DB_HOST=... DB_USER=... DB_NAME=... CORS_ALLOWED_ORIGINS=... APP_PUBLIC_URL=...
#   ./deploy/cloudrun/deploy.sh
#
# Prerequisites: `gcloud auth login`, and the APIs enabled (README §2).

set -euo pipefail

# ---- Required configuration (fail fast if unset) --------------------------
: "${PROJECT_ID:?set PROJECT_ID to your GCP project id}"
: "${REGION:?set REGION, e.g. europe-west1 (match Neon/Upstash region)}"
: "${DB_HOST:?set DB_HOST to the Neon host (prefer the -pooler host with DB_POOLED=true)}"
: "${DB_USER:?set DB_USER}"
: "${DB_NAME:?set DB_NAME}"
: "${CORS_ALLOWED_ORIGINS:?set CORS_ALLOWED_ORIGINS to an explicit https origin allowlist (no '*')}"
: "${APP_PUBLIC_URL:?set APP_PUBLIC_URL to this API's public https base URL}"

# ---- Optional configuration (sensible defaults) ---------------------------
SERVICE="${SERVICE:-memere-api}"
REPO="${REPO:-memere}"
IMAGE_TAG="${IMAGE_TAG:-$(git rev-parse --short HEAD 2>/dev/null || echo latest)}"
IMAGE="${IMAGE:-$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/memere-api:$IMAGE_TAG}"

DB_PORT="${DB_PORT:-5432}"
DB_SSL_MODE="${DB_SSL_MODE:-require}"   # Neon requires TLS
DB_POOLED="${DB_POOLED:-true}"          # true when DB_HOST is the -pooler endpoint
DB_MAX_CONNS="${DB_MAX_CONNS:-10}"      # DB_MAX_CONNS × MAX_INSTANCES must stay < Neon's limit

MIN_INSTANCES="${MIN_INSTANCES:-1}"     # keep >=1: in-proc workers must stay warm
MAX_INSTANCES="${MAX_INSTANCES:-4}"
CPU="${CPU:-1}"
MEMORY="${MEMORY:-512Mi}"
CONCURRENCY="${CONCURRENCY:-80}"
TIMEOUT="${TIMEOUT:-300}"

# Secret Manager secret names (values live only in Secret Manager, never here).
JWT_SECRET_NAME="${JWT_SECRET_NAME:-JWT_SECRET}"
INTERNAL_JOB_TOKEN_NAME="${INTERNAL_JOB_TOKEN_NAME:-INTERNAL_JOB_TOKEN}"
REDIS_URL_NAME="${REDIS_URL_NAME:-REDIS_URL}"
DB_PASSWORD_NAME="${DB_PASSWORD_NAME:-DB_PASSWORD}"

echo ">> project=$PROJECT_ID region=$REGION service=$SERVICE"
echo ">> image=$IMAGE"

gcloud config set project "$PROJECT_ID" >/dev/null

# ---- 1. Ensure the Artifact Registry repo exists --------------------------
if ! gcloud artifacts repositories describe "$REPO" --location="$REGION" >/dev/null 2>&1; then
  echo ">> creating Artifact Registry repo '$REPO' in $REGION"
  gcloud artifacts repositories create "$REPO" \
    --repository-format=docker --location="$REGION"
fi

# ---- 2. Build & push the image FROM Dockerfile.worker (it has FFmpeg) ------
# The distroless Dockerfile has no FFmpeg and would fail transcoding; the single
# always-on service must use the worker image.
echo ">> building $IMAGE from Dockerfile.worker"
gcloud builds submit --tag "$IMAGE" --config=- <<EOF
steps:
  - name: gcr.io/cloud-builders/docker
    args: ['build', '-f', 'Dockerfile.worker', '-t', '$IMAGE', '.']
images: ['$IMAGE']
EOF

# ---- 3. Deploy ------------------------------------------------------------
# Non-secret env (--set-env-vars) and secret env (--set-secrets) are kept
# strictly separate. METRICS_PORT is emptied so the server does not start the
# second /metrics listener Cloud Run cannot route to. The sweep tickers are
# disabled here and driven by Cloud Scheduler (see scheduler.sh); the in-proc
# workers are enabled because their jobs are produced in this same process.
echo ">> deploying $SERVICE"
gcloud run deploy "$SERVICE" \
  --image="$IMAGE" \
  --region="$REGION" \
  --platform=managed \
  --allow-unauthenticated \
  --min-instances="$MIN_INSTANCES" \
  --max-instances="$MAX_INSTANCES" \
  --no-cpu-throttling \
  --cpu="$CPU" \
  --memory="$MEMORY" \
  --concurrency="$CONCURRENCY" \
  --timeout="$TIMEOUT" \
  --set-env-vars="APP_ENV=production" \
  --set-env-vars="APP_PUBLIC_URL=$APP_PUBLIC_URL" \
  --set-env-vars="DB_HOST=$DB_HOST,DB_PORT=$DB_PORT,DB_USER=$DB_USER,DB_NAME=$DB_NAME,DB_SSL_MODE=$DB_SSL_MODE,DB_POOLED=$DB_POOLED,DB_MAX_CONNS=$DB_MAX_CONNS,DB_MIN_CONNS=0" \
  --set-env-vars="CORS_ALLOWED_ORIGINS=$CORS_ALLOWED_ORIGINS" \
  --set-env-vars="METRICS_PORT=" \
  --set-env-vars="SWEEPER_ENABLED=false,SUBSCRIPTION_SWEEP_ENABLED=false,ENGAGEMENT_SWEEP_ENABLED=false" \
  --set-env-vars="TRANSCODE_WORKER_ENABLED=true,NOTIFICATION_WORKER_ENABLED=true" \
  --set-secrets="JWT_SECRET=$JWT_SECRET_NAME:latest,INTERNAL_JOB_TOKEN=$INTERNAL_JOB_TOKEN_NAME:latest,REDIS_URL=$REDIS_URL_NAME:latest,DB_PASSWORD=$DB_PASSWORD_NAME:latest"

echo ">> done. Service URL:"
gcloud run services describe "$SERVICE" --region="$REGION" --format='value(status.url)'

cat <<'NOTE'

Next steps:
  1. Run migrations against Neon (once per schema change):
       DB_HOST=... DB_USER=... DB_PASSWORD=... DB_NAME=... DB_SSL_MODE=require DB_POOLED=true make migrate-up
  2. Create the Cloud Scheduler sweep jobs:
       ./deploy/cloudrun/scheduler.sh
  3. Verify: curl -fsS "$SERVICE_URL/readyz"
NOTE
