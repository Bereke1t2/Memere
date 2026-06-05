#!/usr/bin/env bash
#
# End-to-end smoke test for the Phase 1 API (skill 5 §5.7).
#
# Drives a running stack through the full happy path and the key security
# guarantees: unpublished courses are hidden from anonymous callers, refresh
# rotates tokens, and the login rate limiter returns 429.
#
# Usage:
#   make up && make migrate-up && make run &   # start the stack
#   bash scripts/smoke.sh                       # then run this
#
# Requires: curl, jq. Override the target with BASE_URL (default :8080).

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
API="${BASE_URL}/api/v1"
# A unique email per run so re-runs don't collide on the unique index.
EMAIL="teacher+$(date +%s)@example.com"
PASSWORD="sup3rsecret!"

pass() { printf '  \033[32mPASS\033[0m %s\n' "$1"; }
fail() { printf '  \033[31mFAIL\033[0m %s\n' "$1"; exit 1; }

# req METHOD PATH [JSON_BODY] [AUTH_TOKEN] -> sets $HTTP_CODE and $BODY
req() {
  local method="$1" path="$2" body="${3:-}" token="${4:-}"
  local args=(-sS -o /tmp/smoke_body -w '%{http_code}' -X "$method" "${API}${path}")
  [[ -n "$body" ]] && args+=(-H 'Content-Type: application/json' -d "$body")
  [[ -n "$token" ]] && args+=(-H "Authorization: Bearer ${token}")
  HTTP_CODE="$(curl "${args[@]}")"
  BODY="$(cat /tmp/smoke_body)"
}

echo "==> Target: ${API}"

echo "[1] register teacher"
req POST /auth/register "$(jq -nc --arg e "$EMAIL" --arg p "$PASSWORD" \
  '{email:$e, password:$p, first_name:"Test", last_name:"Teacher", role:"teacher"}')"
[[ "$HTTP_CODE" == "201" ]] || fail "register expected 201, got $HTTP_CODE ($BODY)"
pass "registered ($EMAIL)"

echo "[2] login"
req POST /auth/login "$(jq -nc --arg e "$EMAIL" --arg p "$PASSWORD" '{email:$e, password:$p}')"
[[ "$HTTP_CODE" == "200" ]] || fail "login expected 200, got $HTTP_CODE ($BODY)"
ACCESS="$(jq -r '.access_token' <<<"$BODY")"
REFRESH="$(jq -r '.refresh_token' <<<"$BODY")"
[[ -n "$ACCESS" && "$ACCESS" != "null" ]] || fail "no access token"
pass "logged in"

echo "[3] GET /auth/me"
req GET /auth/me "" "$ACCESS"
[[ "$HTTP_CODE" == "200" ]] || fail "me expected 200, got $HTTP_CODE ($BODY)"
grep -q password_hash <<<"$BODY" && fail "user response leaked password_hash"
pass "me ok, no hash leak"

echo "[4] create course (unpublished)"
req POST /courses "$(jq -nc '{title:"Smoke Course", description:"d", subject:"Math", grade:12, price:0, level:"beginner"}')" "$ACCESS"
[[ "$HTTP_CODE" == "201" ]] || fail "create expected 201, got $HTTP_CODE ($BODY)"
COURSE_ID="$(jq -r '.id' <<<"$BODY")"
[[ "$(jq -r '.is_published' <<<"$BODY")" == "false" ]] || fail "new course should be unpublished"
pass "course created ($COURSE_ID)"

echo "[5] anonymous list excludes unpublished"
req GET /courses
[[ "$HTTP_CODE" == "200" ]] || fail "list expected 200, got $HTTP_CODE"
if jq -e --arg id "$COURSE_ID" '.data[] | select(.id == $id)' <<<"$BODY" >/dev/null; then
  fail "unpublished course visible to anonymous"
fi
pass "unpublished hidden from anonymous"

echo "[6] publish, then anonymous list includes it"
req POST "/courses/${COURSE_ID}/publish" "" "$ACCESS"
[[ "$HTTP_CODE" == "200" ]] || fail "publish expected 200, got $HTTP_CODE ($BODY)"
req GET /courses
jq -e --arg id "$COURSE_ID" '.data[] | select(.id == $id)' <<<"$BODY" >/dev/null \
  || fail "published course not listed for anonymous"
pass "published course now public"

echo "[7] add section + lesson, counters update"
req POST "/courses/${COURSE_ID}/sections" "$(jq -nc '{title:"S1", is_published:true}')" "$ACCESS"
[[ "$HTTP_CODE" == "201" ]] || fail "add section expected 201, got $HTTP_CODE ($BODY)"
SECTION_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/sections/${SECTION_ID}/lessons" "$(jq -nc '{title:"L1", type:"video", duration_seconds:120, is_published:true}')" "$ACCESS"
[[ "$HTTP_CODE" == "201" ]] || fail "add lesson expected 201, got $HTTP_CODE ($BODY)"
req GET "/courses/${COURSE_ID}"
[[ "$(jq -r '.total_lessons' <<<"$BODY")" == "1" ]] || fail "total_lessons != 1"
[[ "$(jq -r '.total_duration_seconds' <<<"$BODY")" == "120" ]] || fail "total_duration_seconds != 120"
pass "counters updated (1 lesson, 120s)"

echo "[8] refresh rotates tokens; old refresh is revoked"
req POST /auth/refresh "$(jq -nc --arg t "$REFRESH" '{refresh_token:$t}')"
[[ "$HTTP_CODE" == "200" ]] || fail "refresh expected 200, got $HTTP_CODE ($BODY)"
NEW_REFRESH="$(jq -r '.refresh_token' <<<"$BODY")"
[[ "$NEW_REFRESH" != "$REFRESH" ]] || fail "refresh token was not rotated"
req POST /auth/refresh "$(jq -nc --arg t "$REFRESH" '{refresh_token:$t}')"
[[ "$HTTP_CODE" == "401" ]] || fail "old refresh token should be revoked (got $HTTP_CODE)"
pass "tokens rotated, old refresh revoked"

echo "[9] login rate limit (6th bad attempt -> 429)"
RL_EMAIL="rl+$(date +%s)@example.com"
req POST /auth/register "$(jq -nc --arg e "$RL_EMAIL" --arg p "$PASSWORD" \
  '{email:$e, password:$p, first_name:"R", last_name:"L", role:"student"}')"
[[ "$HTTP_CODE" == "201" ]] || fail "rl register expected 201, got $HTTP_CODE"
got429=0
for i in $(seq 1 6); do
  req POST /auth/login "$(jq -nc --arg e "$RL_EMAIL" '{email:$e, password:"wrongpass"}')"
  if [[ "$HTTP_CODE" == "429" ]]; then got429=1; break; fi
done
[[ "$got429" == "1" ]] || fail "expected a 429 within 6 login attempts, never got one"
pass "login rate limit enforced (429)"

echo ""
echo "==> ALL SMOKE CHECKS PASSED"
