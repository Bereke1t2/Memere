#!/usr/bin/env bash
#
# scheduler.sh — create (or update) the Cloud Scheduler jobs that drive the
# periodic sweeps on the always-on Cloud Run service.
#
# The sweeps (attempt expiry, subscription expiry, streak warnings) are
# DB-driven and idempotent. Their in-process tickers are disabled on Cloud Run
# (SWEEPER_ENABLED=false, etc.) so they fire exactly once per tick regardless of
# how many instances are running; Cloud Scheduler is the single trigger. Each
# job authenticates with the shared INTERNAL_JOB_TOKEN via the
# X-Internal-Job-Token header (constant-time checked by the InternalToken
# middleware).
#
# The token is read from Secret Manager at run time and passed to Scheduler as a
# header — it is not printed or committed. Re-running this script updates the
# existing jobs in place (create-or-update), so it is safe to run repeatedly.
#
# Usage:
#   export PROJECT_ID=... REGION=... [SERVICE=...]
#   ./deploy/cloudrun/scheduler.sh

set -euo pipefail

: "${PROJECT_ID:?set PROJECT_ID to your GCP project id}"
: "${REGION:?set REGION (must match the Cloud Run + Scheduler region)}"

SERVICE="${SERVICE:-memere-api}"
INTERNAL_JOB_TOKEN_NAME="${INTERNAL_JOB_TOKEN_NAME:-INTERNAL_JOB_TOKEN}"

gcloud config set project "$PROJECT_ID" >/dev/null

SERVICE_URL="$(gcloud run services describe "$SERVICE" --region="$REGION" --format='value(status.url)')"
if [[ -z "$SERVICE_URL" ]]; then
  echo "could not resolve URL for service '$SERVICE' in $REGION — deploy it first" >&2
  exit 1
fi

TOKEN="$(gcloud secrets versions access latest --secret="$INTERNAL_JOB_TOKEN_NAME")"
if [[ -z "$TOKEN" ]]; then
  echo "secret '$INTERNAL_JOB_TOKEN_NAME' is empty — create it first (README §3)" >&2
  exit 1
fi

# create-or-update one HTTP scheduler job.
#   $1 job name   $2 cron schedule   $3 endpoint path
upsert_job() {
  local name="$1" schedule="$2" path="$3"
  local uri="$SERVICE_URL$path"
  local verb=create
  if gcloud scheduler jobs describe "$name" --location="$REGION" >/dev/null 2>&1; then
    verb=update
  fi
  echo ">> $verb $name  ($schedule)  -> $path"
  gcloud scheduler jobs "$verb" http "$name" \
    --location="$REGION" \
    --schedule="$schedule" \
    --time-zone="${SCHEDULER_TZ:-Etc/UTC}" \
    --uri="$uri" \
    --http-method=POST \
    --headers="X-Internal-Job-Token=$TOKEN" \
    --attempt-deadline="${ATTEMPT_DEADLINE:-320s}"
}

# Schedules use off-round minutes to avoid the top-of-hour thundering herd.
upsert_job sweep-attempts      "*/2 * * * *"  /internal/jobs/sweep-attempts
upsert_job sweep-subscriptions "17 * * * *"   /internal/jobs/sweep-subscriptions
upsert_job sweep-engagement    "23 3 * * *"   /internal/jobs/sweep-engagement

echo ">> done. Test one now with:"
echo "     gcloud scheduler jobs run sweep-attempts --location=$REGION"
