#!/usr/bin/env bash
# smoke_phase5.sh — Phase 5 end-to-end smoke test.
# Prerequisites: make up && make migrate-up && make run (server on :8080).
# Usage: bash scripts/smoke_phase5.sh [BASE_URL]
set -euo pipefail

BASE="${1:-http://localhost:8080}/api/v1"
PASS=0
FAIL=0
FAILURES=()

info()  { echo "  [info] $*"; }
ok()    { echo "  ✓ $*"; ((PASS++)); }
fail()  { echo "  ✗ $*"; ((FAIL++)); FAILURES+=("$*"); }

assert_status() {
  local label="$1" expected="$2" actual="$3"
  if [[ "$actual" == "$expected" ]]; then ok "$label (HTTP $actual)"; else fail "$label: expected $expected, got $actual"; fi
}

curl_json() { curl -sf -w "\n%{http_code}" -H "Content-Type: application/json" "$@"; }

echo ""
echo "════════════════════════════════════════════════════════════"
echo "  Phase 5 Smoke Test  →  $BASE"
echo "════════════════════════════════════════════════════════════"

# ── 1. Register users ─────────────────────────────────────────────────────────
echo ""
echo "── 1. Register student, teacher, admin ──────────────────────────────────"

resp=$(curl_json -X POST "$BASE/auth/register" -d '{
  "email":"student@smoke.test","password":"Smoke1234!",
  "first_name":"Alem","last_name":"Tsega","role":"student"}')
http=$(tail -1 <<<"$resp"); body=$(head -1 <<<"$resp")
assert_status "register student" 201 "$http"
STUDENT_TOKEN=$(echo "$body" | python3 -c "import sys,json; print(json.load(sys.stdin)['user']['id'] if 'user' in json.load(sys.stdin) else '')" 2>/dev/null || echo "")

resp=$(curl_json -X POST "$BASE/auth/login" -d '{"email":"student@smoke.test","password":"Smoke1234!"}')
http=$(tail -1 <<<"$resp"); body=$(head -1 <<<"$resp")
assert_status "login student" 200 "$http"
STUDENT_JWT=$(echo "$body" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || echo "")
STUDENT_ID=$(echo "$body"  | python3 -c "import sys,json; print(json.load(sys.stdin).get('user',{}).get('id',''))" 2>/dev/null || echo "")
[[ -n "$STUDENT_JWT" ]] && ok "student JWT obtained" || fail "student JWT missing"

resp=$(curl_json -X POST "$BASE/auth/register" -d '{
  "email":"admin@smoke.test","password":"Smoke1234!",
  "first_name":"Admin","last_name":"User","role":"admin"}')
http=$(tail -1 <<<"$resp")
assert_status "register admin" 201 "$http"

resp=$(curl_json -X POST "$BASE/auth/login" -d '{"email":"admin@smoke.test","password":"Smoke1234!"}')
http=$(tail -1 <<<"$resp"); body=$(head -1 <<<"$resp")
assert_status "login admin" 200 "$http"
ADMIN_JWT=$(echo "$body" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || echo "")
ADMIN_ID=$(echo "$body"  | python3 -c "import sys,json; print(json.load(sys.stdin).get('user',{}).get('id',''))" 2>/dev/null || echo "")
[[ -n "$ADMIN_JWT" ]] && ok "admin JWT obtained" || fail "admin JWT missing"

# ── 2. Non-admin blocked from admin endpoints ─────────────────────────────────
echo ""
echo "── 2. Role gate: non-admin → 403 ────────────────────────────────────────"
resp=$(curl_json -H "Authorization: Bearer $STUDENT_JWT" "$BASE/admin/analytics/overview" 2>/dev/null || echo -e "{}\n403")
http=$(tail -1 <<<"$resp")
assert_status "student blocked from admin overview" 403 "$http"

# ── 3. Dashboard (empty) ──────────────────────────────────────────────────────
echo ""
echo "── 3. Student dashboard ─────────────────────────────────────────────────"
resp=$(curl_json -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/dashboard")
http=$(tail -1 <<<"$resp"); body=$(head -1 <<<"$resp")
assert_status "GET /me/dashboard" 200 "$http"

resp=$(curl_json -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/streak")
http=$(tail -1 <<<"$resp")
assert_status "GET /me/streak" 200 "$http"

# ── 4. Notifications list ─────────────────────────────────────────────────────
echo ""
echo "── 4. Notifications ─────────────────────────────────────────────────────"
resp=$(curl_json -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/notifications")
http=$(tail -1 <<<"$resp")
assert_status "GET /me/notifications" 200 "$http"

resp=$(curl_json -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/notifications/unread-count")
http=$(tail -1 <<<"$resp")
assert_status "GET /me/notifications/unread-count" 200 "$http"

resp=$(curl_json -X POST -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/notifications/read-all")
http=$(tail -1 <<<"$resp")
assert_status "POST /me/notifications/read-all" 204 "$http"

# ── 5. Register FCM device ────────────────────────────────────────────────────
echo ""
echo "── 5. Device registration ───────────────────────────────────────────────"
resp=$(curl_json -X POST -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/devices" \
  -d '{"fcm_token":"test-fcm-token-001","platform":"android"}')
http=$(tail -1 <<<"$resp")
assert_status "POST /me/devices" 201 "$http"

resp=$(curl_json -X DELETE -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/devices/test-fcm-token-001")
http=$(tail -1 <<<"$resp")
assert_status "DELETE /me/devices/:token" 204 "$http"

# ── 6. Notification preferences ───────────────────────────────────────────────
echo ""
echo "── 6. Notification preferences ──────────────────────────────────────────"
resp=$(curl_json -X PUT -H "Authorization: Bearer $STUDENT_JWT" "$BASE/me/notification-preferences" \
  -d '{"push_enabled":true,"email_enabled":false}')
http=$(tail -1 <<<"$resp")
assert_status "PUT /me/notification-preferences" 200 "$http"

# ── 7. Admin analytics overview ───────────────────────────────────────────────
echo ""
echo "── 7. Admin analytics ───────────────────────────────────────────────────"
resp=$(curl_json -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/analytics/overview")
http=$(tail -1 <<<"$resp"); body=$(head -1 <<<"$resp")
assert_status "GET /admin/analytics/overview" 200 "$http"

resp=$(curl_json -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/analytics/revenue")
http=$(tail -1 <<<"$resp")
assert_status "GET /admin/analytics/revenue" 200 "$http"

resp=$(curl_json -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/analytics/engagement")
http=$(tail -1 <<<"$resp")
assert_status "GET /admin/analytics/engagement" 200 "$http"

# ── 8. Admin user management ──────────────────────────────────────────────────
echo ""
echo "── 8. Admin user management ─────────────────────────────────────────────"
resp=$(curl_json -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/users")
http=$(tail -1 <<<"$resp")
assert_status "GET /admin/users" 200 "$http"

if [[ -n "$STUDENT_ID" ]]; then
  resp=$(curl_json -X POST -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/users/$STUDENT_ID/suspend" \
    -d '{"reason":"smoke test suspension"}')
  http=$(tail -1 <<<"$resp")
  assert_status "POST /admin/users/:id/suspend" 204 "$http"

  resp=$(curl_json -X POST -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/users/$STUDENT_ID/reactivate")
  http=$(tail -1 <<<"$resp")
  assert_status "POST /admin/users/:id/reactivate" 204 "$http"
else
  info "skipping suspend/reactivate (no student ID)"
fi

# ── 9. Admin courses + payments ───────────────────────────────────────────────
echo ""
echo "── 9. Admin content & payments ──────────────────────────────────────────"
resp=$(curl_json -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/courses")
http=$(tail -1 <<<"$resp")
assert_status "GET /admin/courses" 200 "$http"

resp=$(curl_json -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/payments")
http=$(tail -1 <<<"$resp")
assert_status "GET /admin/payments" 200 "$http"

resp=$(curl_json -X POST -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/payments/reconcile")
http=$(tail -1 <<<"$resp")
assert_status "POST /admin/payments/reconcile" 200 "$http"

# ── 10. Broadcast announcement ────────────────────────────────────────────────
echo ""
echo "── 10. Broadcast announcement ───────────────────────────────────────────"
resp=$(curl_json -X POST -H "Authorization: Bearer $ADMIN_JWT" "$BASE/admin/announcements" \
  -d '{"title":"Phase 5 smoke","body":"All systems go!","segment":"students"}')
http=$(tail -1 <<<"$resp")
assert_status "POST /admin/announcements" 204 "$http"

# ── 11. Certificate public verify endpoint ────────────────────────────────────
echo ""
echo "── 11. Certificate public verify (unknown serial → 404) ─────────────────"
resp=$(curl -sf -w "\n%{http_code}" "$BASE/verify/certificates/MEM-2026-NOTEXIST" 2>/dev/null || echo -e "{}\n404")
http=$(tail -1 <<<"$resp")
assert_status "GET /verify/certificates/:serial (not found)" 404 "$http"

# ── Summary ───────────────────────────────────────────────────────────────────
echo ""
echo "════════════════════════════════════════════════════════════"
echo "  Results: $PASS passed, $FAIL failed"
if [[ ${#FAILURES[@]} -gt 0 ]]; then
  echo "  Failures:"
  for f in "${FAILURES[@]}"; do echo "    - $f"; done
fi
echo "════════════════════════════════════════════════════════════"
echo ""
[[ $FAIL -eq 0 ]]
