#!/usr/bin/env bash
#
# End-to-end smoke test for the Phase 3 video pipeline (skill 5 §5.5).
#
# Drives a running stack through the full video path and the Phase 3
# Non-Negotiables:
#   - direct-to-storage pre-signed upload (server never proxies the bytes),
#   - FFmpeg -> HLS transcoding via the async worker (status reaches 'ready'),
#   - NO raw storage key ever appears in a client response,
#   - streaming/download URLs are SIGNED and time-limited (carry a signature),
#   - access control: a paid course's video is 403 for a student, allowed once
#     the course is free / the lesson is a preview,
#   - the download token is SINGLE-USE (second redeem fails).
#
# Prereqs (run these first, in another shell):
#   make up && make minio-bucket && make migrate-up
#   make run &            # ffmpeg + ffprobe must be on PATH for the worker
#   bash scripts/smoke_phase3.sh
#
# Requires: curl, jq, ffmpeg/ffprobe (for the committed sample). Override the
# target with BASE_URL (default :8080); READY_WAIT (default 60) bounds how long
# we poll for transcoding to finish.

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
API="${BASE_URL}/api/v1"
READY_WAIT="${READY_WAIT:-60}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SAMPLE="${SCRIPT_DIR}/../testdata/sample.mp4"
TS="$(date +%s)"
PASSWORD="sup3rsecret!"

pass() { printf '  \033[32mPASS\033[0m %s\n' "$1"; }
fail() { printf '  \033[31mFAIL\033[0m %s\n' "$1"; exit 1; }
note() { printf '  \033[33mNOTE\033[0m %s\n' "$1"; }

[[ -r "$SAMPLE" ]] || fail "sample clip not found at $SAMPLE"

# req METHOD PATH [JSON_BODY] [AUTH_TOKEN] -> sets $HTTP_CODE and $BODY
req() {
  local method="$1" path="$2" body="${3:-}" token="${4:-}"
  local args=(-sS -o /tmp/smoke3_body -w '%{http_code}' -X "$method" "${API}${path}")
  [[ -n "$body" ]] && args+=(-H 'Content-Type: application/json' -d "$body")
  [[ -n "$token" ]] && args+=(-H "Authorization: Bearer ${token}")
  HTTP_CODE="$(curl "${args[@]}")"
  BODY="$(cat /tmp/smoke3_body)"
}

# register_login ROLE EMAILPREFIX -> echoes the access token
register_login() {
  local role="$1" email="$2+${TS}@example.com"
  req POST /auth/register "$(jq -nc --arg e "$email" --arg p "$PASSWORD" --arg r "$role" \
    '{email:$e, password:$p, first_name:"T", last_name:"U", role:$r}')"
  [[ "$HTTP_CODE" == "201" ]] || fail "register $role expected 201, got $HTTP_CODE ($BODY)"
  req POST /auth/login "$(jq -nc --arg e "$email" --arg p "$PASSWORD" '{email:$e, password:$p}')"
  [[ "$HTTP_CODE" == "200" ]] || fail "login $role expected 200, got $HTTP_CODE ($BODY)"
  jq -r '.access_token' <<<"$BODY"
}

# assert_no_raw_key BODY — fail if a response exposes a bare storage-key field.
# Pre-signed URLs necessarily contain the object path in their *_url fields, which
# is expected; we strip those and assert no separate raw-key field/value remains.
assert_no_raw_key() {
  local cleaned
  cleaned="$(jq -c 'with_entries(select(.key | endswith("_url") | not))' <<<"$1" 2>/dev/null || echo "$1")"
  grep -qiE 'originals/|hls/|thumbnails/|file_key|_key"[[:space:]]*:' <<<"$cleaned" \
    && fail "response exposed a raw storage key: $cleaned" || true
}

echo "==> Target: ${API}"

echo "[1] provision teacher + student"
TEACHER="$(register_login teacher vteacher)"
STUDENT="$(register_login student vstudent)"
pass "teacher + student provisioned"

echo "[2] teacher creates a PAID course + section + video lesson, publishes"
req POST /courses "$(jq -nc '{title:"Phase3 Smoke", description:"d", subject:"Physics", grade:12, price:100, is_free:false, level:"beginner"}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create course expected 201, got $HTTP_CODE ($BODY)"
COURSE_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/courses/${COURSE_ID}/sections" "$(jq -nc '{title:"S1", order_index:0, is_published:true}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create section expected 201, got $HTTP_CODE ($BODY)"
SECTION_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/sections/${SECTION_ID}/lessons" "$(jq -nc '{title:"L1", type:"video", is_free_preview:false, order_index:0, is_published:true}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create lesson expected 201, got $HTTP_CODE ($BODY)"
LESSON_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/courses/${COURSE_ID}/publish" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "publish course expected 200, got $HTTP_CODE ($BODY)"
pass "course/section/lesson created + published ($COURSE_ID)"

echo "[3] request a pre-signed upload URL (teacher)"
req POST "/lessons/${LESSON_ID}/videos/upload-url" \
  "$(jq -nc '{file_name:"sample.mp4", content_type:"video/mp4", size_bytes:47226}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "upload-url expected 201, got $HTTP_CODE ($BODY)"
VIDEO_ID="$(jq -r '.video_id' <<<"$BODY")"
UPLOAD_URL="$(jq -r '.upload_url' <<<"$BODY")"
[[ -n "$UPLOAD_URL" && "$UPLOAD_URL" != "null" ]] || fail "no upload_url in response ($BODY)"
assert_no_raw_key "$BODY"
pass "got upload URL + video_id ($VIDEO_ID); no raw key in response"

echo "[4] student CANNOT request an upload URL (role gate) -> 403"
req POST "/lessons/${LESSON_ID}/videos/upload-url" \
  "$(jq -nc '{file_name:"x.mp4", content_type:"video/mp4", size_bytes:10}')" "$STUDENT"
[[ "$HTTP_CODE" == "403" ]] || fail "student upload-url expected 403, got $HTTP_CODE ($BODY)"
pass "student blocked from requesting upload URL (403)"

echo "[5] PUT the bytes directly to storage (pre-signed, server not in the path)"
PUT_CODE="$(curl -sS -o /dev/null -w '%{http_code}' -X PUT \
  --upload-file "$SAMPLE" -H 'Content-Type: video/mp4' "$UPLOAD_URL")"
[[ "$PUT_CODE" == "200" ]] || fail "direct PUT to storage expected 200, got $PUT_CODE"
pass "uploaded sample.mp4 directly to object storage (200)"

echo "[6] confirm upload -> processing + enqueue"
req POST "/videos/${VIDEO_ID}/confirm" "" "$TEACHER"
[[ "$HTTP_CODE" == "202" ]] || fail "confirm expected 202, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.status' <<<"$BODY")" == "processing" ]] || fail "status after confirm != processing ($BODY)"
pass "upload confirmed; status=processing"

echo "[7] poll status until 'ready' (worker runs FFmpeg->HLS), max ${READY_WAIT}s"
deadline=$(( $(date +%s) + READY_WAIT ))
STATUS=""
while :; do
  req GET "/videos/${VIDEO_ID}/status" "" "$TEACHER"
  [[ "$HTTP_CODE" == "200" ]] || fail "status expected 200, got $HTTP_CODE ($BODY)"
  STATUS="$(jq -r '.status' <<<"$BODY")"
  [[ "$STATUS" == "ready" ]] && break
  [[ "$STATUS" == "failed" ]] && fail "transcode failed: $(jq -r '.error // "?"' <<<"$BODY")"
  [[ "$(date +%s)" -ge "$deadline" ]] && fail "video not ready after ${READY_WAIT}s (status=$STATUS)"
  sleep 2
done
assert_no_raw_key "$BODY"
[[ "$(jq -r '.duration_seconds' <<<"$BODY")" -ge 1 ]] || fail "ready video has duration_seconds < 1 ($BODY)"
pass "video ready; duration_seconds=$(jq -r '.duration_seconds' <<<"$BODY"); no raw key in status"

echo "[8] owner gets a SIGNED stream URL"
req GET "/videos/${VIDEO_ID}/stream" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "owner stream expected 200, got $HTTP_CODE ($BODY)"
MASTER_URL="$(jq -r '.master_url' <<<"$BODY")"
assert_no_raw_key "$BODY"
grep -qiE 'X-Amz-Signature|Signature=|X-Amz-Expires|Expires=' <<<"$MASTER_URL" \
  || fail "stream master_url is not a signed/expiring URL: $MASTER_URL"
[[ "$(jq -r '.expires_in' <<<"$BODY")" -gt 0 ]] || fail "stream expires_in not > 0"
pass "owner stream URL is signed + expiring (expires_in=$(jq -r '.expires_in' <<<"$BODY")s)"

echo "[9] the signed master manifest actually fetches (200) from storage"
GET_CODE="$(curl -sS -o /dev/null -w '%{http_code}' "$MASTER_URL")"
[[ "$GET_CODE" == "200" ]] || fail "signed master_url fetch expected 200, got $GET_CODE"
pass "signed master.m3u8 fetched from storage (200)"

echo "[10] student is BLOCKED from the paid course's video -> 403"
req GET "/videos/${VIDEO_ID}/stream" "" "$STUDENT"
[[ "$HTTP_CODE" == "403" ]] || fail "student stream of paid content expected 403, got $HTTP_CODE ($BODY)"
pass "student blocked from paid content (403)"

echo "[11] make the course free -> student now gets a stream URL"
req PUT "/courses/${COURSE_ID}" "$(jq -nc '{is_free:true}')" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "update course is_free expected 200, got $HTTP_CODE ($BODY)"
req GET "/videos/${VIDEO_ID}/stream" "" "$STUDENT"
[[ "$HTTP_CODE" == "200" ]] || fail "student stream of free course expected 200, got $HTTP_CODE ($BODY)"
pass "student allowed once course is free"

echo "[12] download-url -> single-use token; redeem once (302) then fail"
req GET "/videos/${VIDEO_ID}/download-url" "" "$STUDENT"
[[ "$HTTP_CODE" == "200" ]] || fail "download-url expected 200, got $HTTP_CODE ($BODY)"
TOKEN="$(jq -r '.token' <<<"$BODY")"
[[ -n "$TOKEN" && "$TOKEN" != "null" ]] || fail "no download token ($BODY)"
# First redeem: 302 redirect to the signed manifest (curl does not follow it).
RED_CODE="$(curl -sS -o /dev/null -w '%{http_code}' -H "Authorization: Bearer ${STUDENT}" \
  "${API}/videos/download/${TOKEN}")"
[[ "$RED_CODE" == "302" ]] || fail "first token redeem expected 302, got $RED_CODE"
# Second redeem: token consumed -> 404.
req GET "/videos/download/${TOKEN}" "" "$STUDENT"
[[ "$HTTP_CODE" == "404" ]] || fail "second token redeem expected 404, got $HTTP_CODE ($BODY)"
pass "download token is single-use (302 then 404)"

echo ""
echo "==> ALL PHASE 3 SMOKE CHECKS PASSED"
