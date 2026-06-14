#!/usr/bin/env bash
#
# End-to-end smoke test for the Phase 4 payment flow (skill5 §5.6).
#
# Drives a running stack through the full purchase path and the Phase 4
# Non-Negotiables:
#   - idempotent initiate: a retried Idempotency-Key never double-charges,
#   - signature-verified, raw-body provider webhook settles the payment,
#   - webhook dedup + guarded tx: a re-delivered event never double-grants,
#   - access gating: a paid course's quiz is 403 NOT_ENROLLED until purchase,
#   - subscription purchase grants an all-access subscription,
#   - revenue: teacher earnings == gross * the configured share (0.70),
#   - no provider secret or raw payload is returned to the client.
#
# It uses the test-only "mock" provider, which settles OFFLINE and signs its own
# webhook with the same HMAC-SHA256 scheme as the real providers — so the
# signature/dedup/fulfillment paths are exercised without a live provider. The
# server MUST be started with the mock provider enabled:
#
#   PAYMENT_MOCK_ENABLED=true PAYMENT_MOCK_WEBHOOK_SECRET=mock-secret \
#   SUBSCRIPTION_SWEEP_INTERVAL=2s make run &
#
# Prereqs (run first, in another shell):
#   make up && make migrate-up
#   PAYMENT_MOCK_ENABLED=true make run &
#   bash scripts/smoke_phase4.sh
#
# Requires: curl, jq, openssl. Override the target with BASE_URL (default :8080)
# and the mock secret with MOCK_SECRET (default mock-secret). Steps needing an
# admin (refund) or a coupon row use docker-compose psql and are skipped with a
# NOTE if the DB is not reachable that way.

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
API="${BASE_URL}/api/v1"
MOCK_SECRET="${MOCK_SECRET:-mock-secret}"
TS="$(date +%s)"
PASSWORD="sup3rsecret!"

pass() { printf '  \033[32mPASS\033[0m %s\n' "$1"; }
fail() { printf '  \033[31mFAIL\033[0m %s\n' "$1"; exit 1; }
note() { printf '  \033[33mNOTE\033[0m %s\n' "$1"; }

# req METHOD PATH [JSON_BODY] [AUTH_TOKEN] [IDEMPOTENCY_KEY] -> sets $HTTP_CODE, $BODY
req() {
  local method="$1" path="$2" body="${3:-}" token="${4:-}" idem="${5:-}"
  local args=(-sS -o /tmp/smoke4_body -w '%{http_code}' -X "$method" "${API}${path}")
  [[ -n "$body" ]] && args+=(-H 'Content-Type: application/json' -d "$body")
  [[ -n "$token" ]] && args+=(-H "Authorization: Bearer ${token}")
  [[ -n "$idem" ]] && args+=(-H "Idempotency-Key: ${idem}")
  HTTP_CODE="$(curl "${args[@]}")"
  BODY="$(cat /tmp/smoke4_body)"
}

# send_webhook JSON_BODY -> signs the EXACT bytes and POSTs to the mock webhook.
# Sets $HTTP_CODE, $BODY. The signature is over the raw body (no trailing
# newline) so it matches what the server reads via io.LimitReader.
send_webhook() {
  local payload="$1" sig
  sig="$(printf '%s' "$payload" | openssl dgst -sha256 -hmac "$MOCK_SECRET" | awk '{print $NF}')"
  HTTP_CODE="$(curl -sS -o /tmp/smoke4_body -w '%{http_code}' -X POST \
    -H 'Content-Type: application/json' -H "X-Mock-Signature: ${sig}" \
    -d "$payload" "${API}/webhooks/payments/mock")"
  BODY="$(cat /tmp/smoke4_body)"
}

# db_exec SQL -> runs SQL via docker-compose psql. Returns non-zero (silently) if
# the DB is not reachable that way, so callers can degrade to a NOTE.
db_exec() {
  docker compose exec -T postgres psql -U "${DB_USER:-postgres}" -d "${DB_NAME:-memere}" \
    -tAc "$1" 2>/dev/null
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

# assert_no_secret BODY — fail if a client response leaks a webhook/raw payload
# or a provider secret field.
assert_no_secret() {
  grep -qiE 'webhook_secret|secret_key|raw_payload|"metadata"|x-mock-signature' <<<"$1" \
    && fail "response leaked a secret/raw payload: $1" || true
}

echo "==> Target: ${API} (mock provider, secret '${MOCK_SECRET}')"

# Preflight: the mock provider must be enabled, else initiate(provider=mock) 400s.
req GET /subscription-plans
[[ "$HTTP_CODE" == "200" ]] || fail "GET /subscription-plans expected 200, got $HTTP_CODE ($BODY) — is the server up?"

echo "[1] provision teacher + student"
TEACHER="$(register_login teacher payteacher)"
STUDENT="$(register_login student paystudent)"
pass "teacher + student provisioned"

echo "[2] teacher creates a PAID course + published quiz (1 question)"
req POST /courses "$(jq -nc '{title:"Phase4 Smoke", description:"d", subject:"Math", grade:12, price:100, is_free:false, level:"beginner"}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create course expected 201, got $HTTP_CODE ($BODY)"
COURSE_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/courses/${COURSE_ID}/quizzes" "$(jq -nc '{title:"Q1", pass_percentage:50, time_limit_seconds:1800, max_attempts:3}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create quiz expected 201, got $HTTP_CODE ($BODY)"
QUIZ_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/quizzes/${QUIZ_ID}/questions" '{"text":"2+2?","type":"multiple_choice","points":2,"subject":"Math","topic":"Arithmetic","order_index":0,"answers":[{"text":"3","is_correct":false,"order_index":0},{"text":"4","is_correct":true,"order_index":1}]}' "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "add question expected 201, got $HTTP_CODE ($BODY)"
req POST "/courses/${COURSE_ID}/publish" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "publish course expected 200, got $HTTP_CODE ($BODY)"
pass "paid course + quiz created + published ($COURSE_ID)"

echo "[3] student is BLOCKED from the paid quiz before purchase -> 403"
req POST "/quizzes/${QUIZ_ID}/attempts" "" "$STUDENT"
[[ "$HTTP_CODE" == "403" ]] || fail "pre-purchase quiz attempt expected 403, got $HTTP_CODE ($BODY)"
pass "student blocked from paid quiz (403 — Skill 2 access hook live)"

echo "[4] initiate payment (provider=mock) WITHOUT Idempotency-Key -> 400"
req POST /payments/initiate "$(jq -nc --arg c "$COURSE_ID" '{course_id:$c, provider:"mock"}')" "$STUDENT"
[[ "$HTTP_CODE" == "400" ]] || fail "initiate w/o idem-key expected 400, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.code' <<<"$BODY")" == "IDEMPOTENCY_KEY_REQUIRED" ]] || fail "expected IDEMPOTENCY_KEY_REQUIRED ($BODY)"
pass "initiate without Idempotency-Key rejected (400)"

echo "[5] initiate payment WITH Idempotency-Key -> 201 (payment_id + redirect_url)"
IDEM="course-${TS}"
req POST /payments/initiate "$(jq -nc --arg c "$COURSE_ID" '{course_id:$c, provider:"mock"}')" "$STUDENT" "$IDEM"
[[ "$HTTP_CODE" == "201" ]] || fail "initiate expected 201, got $HTTP_CODE ($BODY)"
PAYMENT_ID="$(jq -r '.payment_id' <<<"$BODY")"
[[ -n "$PAYMENT_ID" && "$PAYMENT_ID" != "null" ]] || fail "no payment_id ($BODY)"
[[ "$(jq -r '.redirect_url' <<<"$BODY")" != "null" ]] || fail "no redirect_url ($BODY)"
[[ "$(jq -r '.amount' <<<"$BODY")" == "100" ]] || fail "amount expected 100, got $(jq -r '.amount' <<<"$BODY")"
assert_no_secret "$BODY"
pass "checkout initiated ($PAYMENT_ID), amount=100"

echo "[6] repeat the SAME request with the SAME key -> same payment_id (no double-charge)"
req POST /payments/initiate "$(jq -nc --arg c "$COURSE_ID" '{course_id:$c, provider:"mock"}')" "$STUDENT" "$IDEM"
[[ "$HTTP_CODE" == "201" ]] || fail "idempotent replay expected 201, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.payment_id' <<<"$BODY")" == "$PAYMENT_ID" ]] || fail "idempotent replay returned a different payment_id ($BODY)"
pass "idempotent initiate returns the same payment_id (no double-charge)"

echo "[7] payment appears in the student's history (pending)"
req GET /payments "" "$STUDENT"
[[ "$HTTP_CODE" == "200" ]] || fail "list payments expected 200, got $HTTP_CODE ($BODY)"
jq -e --arg id "$PAYMENT_ID" '.data[] | select(.payment_id == $id and .status == "pending")' <<<"$BODY" >/dev/null \
  || fail "pending payment not in history ($BODY)"
pass "payment listed as pending in /payments"

echo "[8] deliver the (signed) provider webhook: success -> 200, payment completed"
WH="$(jq -nc --arg ev "evt-${TS}-1" --arg ref "$PAYMENT_ID" \
  '{event_id:$ev, type:"charge.success", tx_ref:$ref, status:"success", transaction_id:("txn-"+$ref)}')"
send_webhook "$WH"
[[ "$HTTP_CODE" == "200" ]] || fail "webhook expected 200, got $HTTP_CODE ($BODY)"
req GET "/payments/${PAYMENT_ID}/status" "" "$STUDENT"
[[ "$(jq -r '.status' <<<"$BODY")" == "completed" ]] || fail "payment status expected completed, got $(jq -r '.status' <<<"$BODY")"
pass "webhook settled the payment (completed)"

echo "[9] webhook signature is enforced: a bad signature -> 401"
BAD_CODE="$(curl -sS -o /dev/null -w '%{http_code}' -X POST \
  -H 'Content-Type: application/json' -H 'X-Mock-Signature: deadbeef' \
  -d "$WH" "${API}/webhooks/payments/mock")"
[[ "$BAD_CODE" == "401" ]] || fail "bad-signature webhook expected 401, got $BAD_CODE"
pass "tampered/unsigned webhook rejected (401)"

echo "[10] re-deliver the SAME webhook -> 200 and STILL exactly one enrollment (dedup)"
send_webhook "$WH"
[[ "$HTTP_CODE" == "200" ]] || fail "webhook re-delivery expected 200, got $HTTP_CODE ($BODY)"
req GET /me/enrollments "" "$STUDENT"
[[ "$HTTP_CODE" == "200" ]] || fail "list enrollments expected 200, got $HTTP_CODE ($BODY)"
COUNT="$(jq --arg c "$COURSE_ID" '[.data[] | select(.course_id == $c)] | length' <<<"$BODY")"
[[ "$COUNT" == "1" ]] || fail "expected exactly 1 enrollment for the course, got $COUNT ($BODY)"
pass "webhook dedup + guarded tx: exactly one enrollment after re-delivery"

echo "[11] student now PASSES the access gate -> quiz attempt 200"
req POST "/quizzes/${QUIZ_ID}/attempts" "" "$STUDENT"
[[ "$HTTP_CODE" == "200" || "$HTTP_CODE" == "201" ]] || fail "post-purchase quiz attempt expected 200/201, got $HTTP_CODE ($BODY)"
pass "student can take the paid quiz after purchase"

echo "[12] a SECOND initiate for the now-owned course -> 409 ALREADY_ENROLLED"
req POST /payments/initiate "$(jq -nc --arg c "$COURSE_ID" '{course_id:$c, provider:"mock"}')" "$STUDENT" "course-${TS}-again"
[[ "$HTTP_CODE" == "409" ]] || fail "re-purchase of owned course expected 409, got $HTTP_CODE ($BODY)"
pass "owned-course re-purchase rejected (409 ALREADY_ENROLLED)"

echo "[13] subscription-plans is public and lists monthly + annual"
req GET /subscription-plans
jq -e '.data[] | select(.plan == "monthly")' <<<"$BODY" >/dev/null || fail "monthly plan missing ($BODY)"
jq -e '.data[] | select(.plan == "annual")'  <<<"$BODY" >/dev/null || fail "annual plan missing ($BODY)"
pass "plan catalogue public (monthly + annual)"

echo "[14] subscribe (monthly) -> initiate + webhook -> all-access subscription"
req POST /subscriptions "$(jq -nc '{plan:"monthly", provider:"mock"}')" "$STUDENT" "sub-${TS}"
[[ "$HTTP_CODE" == "201" ]] || fail "subscribe expected 201, got $HTTP_CODE ($BODY)"
SUB_PAYMENT="$(jq -r '.payment_id' <<<"$BODY")"
[[ -n "$SUB_PAYMENT" && "$SUB_PAYMENT" != "null" ]] || fail "no subscription payment_id ($BODY)"
WHS="$(jq -nc --arg ev "evt-${TS}-sub" --arg ref "$SUB_PAYMENT" \
  '{event_id:$ev, type:"charge.success", tx_ref:$ref, status:"success", transaction_id:("txn-"+$ref)}')"
send_webhook "$WHS"
[[ "$HTTP_CODE" == "200" ]] || fail "subscription webhook expected 200, got $HTTP_CODE ($BODY)"
req GET /me/subscription "" "$STUDENT"
[[ "$HTTP_CODE" == "200" ]] || fail "get subscription expected 200, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.status' <<<"$BODY")" == "active" ]] || fail "subscription not active ($BODY)"
[[ "$(jq -r '.plan' <<<"$BODY")" == "monthly" ]] || fail "subscription plan != monthly ($BODY)"
SUB_ID="$(jq -r '.id' <<<"$BODY")"
assert_no_secret "$BODY"
pass "monthly subscription active ($SUB_ID)"

echo "[15] cancel the subscription -> 204 (no-renew; access kept until period end)"
req POST "/subscriptions/${SUB_ID}/cancel" "" "$STUDENT"
[[ "$HTTP_CODE" == "204" ]] || fail "cancel expected 204, got $HTTP_CODE ($BODY)"
pass "subscription canceled (will not renew)"

echo "[16] teacher earnings == gross * 0.70 (configured share)"
req GET /me/earnings "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "earnings expected 200, got $HTTP_CODE ($BODY)"
GROSS="$(jq -r '.gross' <<<"$BODY")"
EARN="$(jq -r '.earnings' <<<"$BODY")"
SHARE="$(jq -r '.teacher_share' <<<"$BODY")"
EXPECT="$(awk -v g="$GROSS" 'BEGIN{printf "%.2f", g*0.70}')"
[[ "$(awk -v g="$GROSS" 'BEGIN{print (g>0)}')" == "1" ]] || fail "teacher gross not > 0 ($BODY)"
[[ "$EARN" == "$EXPECT" ]] || fail "earnings $EARN != gross*0.70 ($EXPECT) [gross=$GROSS share=$SHARE]"
pass "earnings $EARN == gross $GROSS * 0.70"

echo "[17] a student CANNOT view teacher earnings -> 403"
req GET /me/earnings "" "$STUDENT"
[[ "$HTTP_CODE" == "403" ]] || fail "student earnings expected 403, got $HTTP_CODE ($BODY)"
pass "earnings gated to teacher/admin (student 403)"

echo "[18] a non-admin CANNOT refund -> 403"
req POST "/payments/${PAYMENT_ID}/refund" "" "$STUDENT"
[[ "$HTTP_CODE" == "403" ]] || fail "non-admin refund expected 403, got $HTTP_CODE ($BODY)"
pass "refund gated to admin (student 403)"

echo "[19] admin refund -> completed payment becomes refunded (needs DB to mint an admin)"
if ADMIN_EMAIL="payadmin+${TS}@example.com" && \
   ADMIN="$(register_login student payadmin 2>/dev/null)" && \
   db_exec "UPDATE auth.users SET role='admin' WHERE email='${ADMIN_EMAIL}';" >/dev/null; then
  # Re-login so the JWT carries the elevated role.
  req POST /auth/login "$(jq -nc --arg e "$ADMIN_EMAIL" --arg p "$PASSWORD" '{email:$e, password:$p}')"
  ADMIN="$(jq -r '.access_token' <<<"$BODY")"
  req POST "/payments/${PAYMENT_ID}/refund" "" "$ADMIN"
  [[ "$HTTP_CODE" == "200" ]] || fail "admin refund expected 200, got $HTTP_CODE ($BODY)"
  [[ "$(jq -r '.status' <<<"$BODY")" == "refunded" ]] || fail "payment not refunded ($BODY)"
  pass "admin refunded the completed payment"
else
  note "skipped admin-refund (could not promote an admin via docker-compose psql)"
fi

echo "[20] coupon: a discounted second purchase; coupon use increments only on success"
COUPON_CODE="SMOKE${TS}"
if db_exec "INSERT INTO payments.coupons (code, discount_type, discount_value, max_uses, applicable_to) VALUES ('${COUPON_CODE}', 'percentage', 25, 100, 'all');" >/dev/null; then
  req POST /courses "$(jq -nc '{title:"Phase4 Smoke 2", description:"d", subject:"Math", grade:12, price:200, is_free:false, level:"beginner"}')" "$TEACHER"
  [[ "$HTTP_CODE" == "201" ]] || fail "create course2 expected 201, got $HTTP_CODE ($BODY)"
  COURSE2="$(jq -r '.id' <<<"$BODY")"
  req POST "/courses/${COURSE2}/publish" "" "$TEACHER"
  [[ "$HTTP_CODE" == "200" ]] || fail "publish course2 expected 200, got $HTTP_CODE ($BODY)"
  req POST /payments/initiate "$(jq -nc --arg c "$COURSE2" --arg code "$COUPON_CODE" '{course_id:$c, provider:"mock", coupon_code:$code}')" "$STUDENT" "coupon-${TS}"
  [[ "$HTTP_CODE" == "201" ]] || fail "coupon initiate expected 201, got $HTTP_CODE ($BODY)"
  CPAY="$(jq -r '.payment_id' <<<"$BODY")"
  AMT="$(jq -r '.amount' <<<"$BODY")"
  [[ "$(awk -v a="$AMT" 'BEGIN{print (a==150)}')" == "1" ]] || fail "discounted amount expected 150 (25% off 200), got $AMT"
  [[ "$(db_exec "SELECT used_count FROM payments.coupons WHERE code='${COUPON_CODE}';")" == "0" ]] \
    || fail "coupon used_count incremented before payment success"
  WHC="$(jq -nc --arg ev "evt-${TS}-coupon" --arg ref "$CPAY" '{event_id:$ev, type:"charge.success", tx_ref:$ref, status:"success", transaction_id:("txn-"+$ref)}')"
  send_webhook "$WHC"
  [[ "$HTTP_CODE" == "200" ]] || fail "coupon webhook expected 200, got $HTTP_CODE ($BODY)"
  [[ "$(db_exec "SELECT used_count FROM payments.coupons WHERE code='${COUPON_CODE}';")" == "1" ]] \
    || fail "coupon used_count did not increment to 1 after success"
  pass "coupon discounted to 150; used_count 0 -> 1 only after success"
else
  note "skipped coupon test (could not insert a coupon via docker-compose psql)"
fi

echo ""
echo "==> ALL PHASE 4 SMOKE CHECKS PASSED"
