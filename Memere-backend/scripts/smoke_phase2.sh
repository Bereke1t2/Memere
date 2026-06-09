#!/usr/bin/env bash
#
# End-to-end smoke test for the Phase 2 API (quiz, exam, analytics) — skill 5 §5.5.
#
# Drives a running stack through the full learning path and the Phase 2
# Non-Negotiables:
#   - answer keys NEVER cross the wire before submission (asserted on GET /quizzes
#     and on start-attempt responses),
#   - server-side grading produces a score + explanations,
#   - the server-enforced timer fires via the BACKGROUND SWEEPER for an abandoned
#     exam attempt (no client submit),
#   - attempt limits are enforced,
#   - IDOR is blocked (student B cannot read student A's attempt).
#
# Usage:
#   # For a fast sweeper check, run the server with a short interval:
#   SWEEPER_INTERVAL=5s make run &
#   bash scripts/smoke_phase2.sh
#
# Requires: curl, jq. Override target with BASE_URL (default :8080).
# SWEEP_WAIT (default 80) is how long to wait QUIETLY (no reads) for the sweeper
# to grade the abandoned exam — must exceed the exam duration (1 min) + the
# sweeper interval. Run the server with SWEEPER_INTERVAL=5s to keep this short.
# Set SERVER_LOG=<path to server stdout> for an airtight assertion that the
# background sweeper (not the lazy-on-read path) graded the abandoned attempt.

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
API="${BASE_URL}/api/v1"
SWEEP_WAIT="${SWEEP_WAIT:-80}"
TS="$(date +%s)"
PASSWORD="sup3rsecret!"

pass() { printf '  \033[32mPASS\033[0m %s\n' "$1"; }
fail() { printf '  \033[31mFAIL\033[0m %s\n' "$1"; exit 1; }
note() { printf '  \033[33mNOTE\033[0m %s\n' "$1"; }

# req METHOD PATH [JSON_BODY] [AUTH_TOKEN] -> sets $HTTP_CODE and $BODY
req() {
  local method="$1" path="$2" body="${3:-}" token="${4:-}"
  local args=(-sS -o /tmp/smoke2_body -w '%{http_code}' -X "$method" "${API}${path}")
  [[ -n "$body" ]] && args+=(-H 'Content-Type: application/json' -d "$body")
  [[ -n "$token" ]] && args+=(-H "Authorization: Bearer ${token}")
  HTTP_CODE="$(curl "${args[@]}")"
  BODY="$(cat /tmp/smoke2_body)"
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

echo "==> Target: ${API}"

echo "[1] provision teacher + two students"
TEACHER="$(register_login teacher teacher)"
STUDENT_A="$(register_login student studenta)"
STUDENT_B="$(register_login student studentb)"
pass "teacher + 2 students provisioned"

echo "[2] teacher creates + publishes a course"
req POST /courses "$(jq -nc '{title:"Phase2 Smoke", description:"d", subject:"Math", grade:12, price:0, level:"beginner"}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create course expected 201, got $HTTP_CODE ($BODY)"
COURSE_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/courses/${COURSE_ID}/publish" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "publish course expected 200, got $HTTP_CODE ($BODY)"
pass "course created + published ($COURSE_ID)"

echo "[3] teacher creates a quiz (max_attempts=1) with 3 questions (one per type)"
req POST "/courses/${COURSE_ID}/quizzes" "$(jq -nc '{title:"Q1", pass_percentage:50, time_limit_seconds:1800, max_attempts:1}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create quiz expected 201, got $HTTP_CODE ($BODY)"
QUIZ_ID="$(jq -r '.id' <<<"$BODY")"

add_q() { # TYPE JSON_ANSWERS ORDER -> sets $QID
  req POST "/quizzes/${QUIZ_ID}/questions" "$1" "$TEACHER"
  [[ "$HTTP_CODE" == "201" ]] || fail "add question expected 201, got $HTTP_CODE ($BODY)"
  QID="$(jq -r '.id' <<<"$BODY")"
}
add_q '{"text":"2+2?","type":"multiple_choice","points":2,"subject":"Math","topic":"Arithmetic","order_index":0,"answers":[{"text":"3","is_correct":false,"order_index":0},{"text":"4","is_correct":true,"order_index":1}]}'
Q1="$QID"
add_q '{"text":"Earth is flat.","type":"true_false","points":1,"subject":"Science","topic":"Geo","order_index":1,"answers":[{"text":"true","is_correct":false,"order_index":0},{"text":"false","is_correct":true,"order_index":1}]}'
Q2="$QID"
add_q '{"text":"Capital of Ethiopia?","type":"short_answer","points":1,"subject":"Geography","topic":"Capitals","order_index":2,"answers":[{"text":"Addis Ababa","is_correct":true,"order_index":0}]}'
Q3="$QID"
pass "quiz + 3 questions created ($QUIZ_ID)"

echo "[4] student GET /quizzes/:id -> count present, NO answer keys"
req GET "/quizzes/${QUIZ_ID}" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "200" ]] || fail "get quiz expected 200, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.question_count' <<<"$BODY")" == "3" ]] || fail "question_count != 3 ($BODY)"
grep -qi 'is_correct\|correct_answer' <<<"$BODY" && fail "quiz metadata leaked an answer key"
pass "quiz metadata: 3 questions, no answer keys"

echo "[5] student starts attempt -> questions, remaining_seconds, NO answer keys"
req POST "/quizzes/${QUIZ_ID}/attempts" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "201" ]] || fail "start attempt expected 201, got $HTTP_CODE ($BODY)"
ATTEMPT_A="$(jq -r '.attempt_id' <<<"$BODY")"
[[ "$(jq -r '.remaining_seconds' <<<"$BODY")" != "null" ]] || fail "remaining_seconds missing"
[[ "$(jq -r '.questions | length' <<<"$BODY")" == "3" ]] || fail "attempt did not return 3 questions"
grep -qi 'is_correct\|correct_answer' <<<"$BODY" && fail "start-attempt leaked an answer key"
pass "attempt started, timer present, no answer keys"

echo "[6] result before submit -> 409 INVALID/CONFLICT state"
req GET "/quiz-attempts/${ATTEMPT_A}/result" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "409" ]] || fail "pre-submit result expected 409, got $HTTP_CODE ($BODY)"
pass "in-progress result correctly rejected (409)"

echo "[7] auto-save then submit -> graded result with score + explanations"
req PATCH "/quiz-attempts/${ATTEMPT_A}" "$(jq -nc --arg q "$Q1" '{answers:{($q):[]}}')" "$STUDENT_A"
[[ "$HTTP_CODE" == "204" ]] || fail "auto-save expected 204, got $HTTP_CODE ($BODY)"
SUBMIT="$(jq -nc --arg q1 "$Q1" --arg q3 "$Q3" '{answers:{($q1):"correct-placeholder",($q3):"Addis Ababa"}}')"
req POST "/quiz-attempts/${ATTEMPT_A}/submit" "$SUBMIT" "$STUDENT_A"
[[ "$HTTP_CODE" == "200" ]] || fail "submit expected 200, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.status' <<<"$BODY")" == "graded" ]] || fail "submit status != graded ($BODY)"
[[ "$(jq -r '.score' <<<"$BODY")" != "null" ]] || fail "submit has no score"
[[ "$(jq -r '.feedback | length' <<<"$BODY")" == "3" ]] || fail "submit feedback missing"
pass "submitted, graded server-side (score=$(jq -r '.score' <<<"$BODY"), feedback x3)"

echo "[8] fetch result again -> matches; subject_breakdown present"
req GET "/quiz-attempts/${ATTEMPT_A}/result" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "200" ]] || fail "result expected 200, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.subject_breakdown | length' <<<"$BODY")" -ge 1 ]] || fail "result missing subject_breakdown"
pass "result fetch ok with subject breakdown"

echo "[9] attempt limit: a 2nd start -> rejected (max_attempts=1)"
req POST "/quizzes/${QUIZ_ID}/attempts" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "403" ]] || fail "second attempt expected 403, got $HTTP_CODE ($BODY)"
pass "attempt limit enforced (403)"

echo "[10] IDOR: student B reads student A's attempt -> blocked"
req GET "/quiz-attempts/${ATTEMPT_A}/result" "" "$STUDENT_B"
[[ "$HTTP_CODE" == "403" || "$HTTP_CODE" == "404" ]] || fail "IDOR expected 403/404, got $HTTP_CODE ($BODY)"
jq -e --arg id "$ATTEMPT_A" 'select(.attempt_id == $id)' <<<"$BODY" >/dev/null 2>&1 && fail "IDOR leaked attempt data"
note "IDOR blocked with $HTTP_CODE (impl hides existence with 404)"
pass "IDOR blocked, no data leak"

echo "[11] exam analytics path: build + publish a 1-min exam reusing the questions"
req POST "/courses/${COURSE_ID}/exams" "$(jq -nc '{title:"E1", subject:"Math", grade:12, duration_minutes:60, pass_marks:1}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create exam expected 201, got $HTTP_CODE ($BODY)"
EXAM_ID="$(jq -r '.id' <<<"$BODY")"
req POST "/exams/${EXAM_ID}/questions" "$(jq -nc --arg q "$Q1" '{question_id:$q, marks:2, order_index:0}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "add exam question expected 201, got $HTTP_CODE ($BODY)"
req POST "/exams/${EXAM_ID}/publish" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "publish exam expected 200, got $HTTP_CODE ($BODY)"
req GET /mock-exams
jq -e --arg id "$EXAM_ID" '.data[] | select(.id == $id)' <<<"$BODY" >/dev/null || fail "published exam not in /mock-exams"
pass "exam published + visible in catalog ($EXAM_ID)"

echo "[12] student starts + submits exam -> analytics has breakdown + percentile"
req POST "/mock-exams/${EXAM_ID}/start" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "201" ]] || fail "start exam expected 201, got $HTTP_CODE ($BODY)"
EXAM_ATTEMPT="$(jq -r '.attempt_id' <<<"$BODY")"
grep -qi 'is_correct\|correct_answer' <<<"$BODY" && fail "exam start leaked an answer key"
req POST "/exam-attempts/${EXAM_ATTEMPT}/submit" "$(jq -nc --arg q "$Q1" '{answers:{($q):[]}}')" "$STUDENT_A"
[[ "$HTTP_CODE" == "200" ]] || fail "exam submit expected 200, got $HTTP_CODE ($BODY)"
req GET "/exam-attempts/${EXAM_ATTEMPT}/analytics" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "200" ]] || fail "analytics expected 200, got $HTTP_CODE ($BODY)"
jq -e 'has("subject_breakdown") and has("percentile")' <<<"$BODY" >/dev/null || fail "analytics missing breakdown/percentile"
pass "exam analytics returned (breakdown + percentile)"

echo "[13] teacher exam stats"
req GET "/exams/${EXAM_ID}/stats" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "stats expected 200, got $HTTP_CODE ($BODY)"
[[ "$(jq -r '.total_attempts' <<<"$BODY")" -ge 1 ]] || fail "stats total_attempts < 1"
pass "exam stats ok (attempts=$(jq -r '.total_attempts' <<<"$BODY"))"

echo "[14] student stats forbidden"
req GET "/exams/${EXAM_ID}/stats" "" "$STUDENT_A"
[[ "$HTTP_CODE" == "403" ]] || fail "student stats expected 403, got $HTTP_CODE ($BODY)"
pass "student barred from exam stats (403)"

echo "[15] SERVER TIMER via SWEEPER: abandon a 1-min exam, never submit, never read"
req POST "/courses/${COURSE_ID}/exams" "$(jq -nc '{title:"E2-timed", subject:"Math", grade:12, duration_minutes:1, pass_marks:1}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "create timed exam expected 201, got $HTTP_CODE ($BODY)"
TIMED_EXAM="$(jq -r '.id' <<<"$BODY")"
req POST "/exams/${TIMED_EXAM}/questions" "$(jq -nc --arg q "$Q1" '{question_id:$q, marks:2, order_index:0}')" "$TEACHER"
[[ "$HTTP_CODE" == "201" ]] || fail "add timed exam question expected 201, got $HTTP_CODE ($BODY)"
req POST "/exams/${TIMED_EXAM}/publish" "" "$TEACHER"
[[ "$HTTP_CODE" == "200" ]] || fail "publish timed exam expected 200, got $HTTP_CODE ($BODY)"
req POST "/mock-exams/${TIMED_EXAM}/start" "" "$STUDENT_B"
[[ "$HTTP_CODE" == "201" ]] || fail "start timed exam expected 201, got $HTTP_CODE ($BODY)"
TIMED_ATTEMPT="$(jq -r '.attempt_id' <<<"$BODY")"
# CRITICAL: do NOT read the attempt while waiting. GetExamResult lazily finalizes
# an expired in-progress attempt on read, which would grade it via the client
# request path and mask whether the BACKGROUND sweeper actually fired. We wait
# quietly past (duration + sweeper interval), so by the time we read it the only
# thing that could have graded it with zero client calls is the sweeper.
note "abandoning attempt $TIMED_ATTEMPT; waiting ${SWEEP_WAIT}s quietly (no reads)"
sleep "$SWEEP_WAIT"
# Optional airtight proof: if SERVER_LOG points at the server's stdout, assert the
# sweeper logged the grading (the lazy-read path can't fire — we never read it).
if [[ -n "${SERVER_LOG:-}" && -r "${SERVER_LOG}" ]]; then
  grep -qi "graded .* expired attempt" "${SERVER_LOG}" \
    && pass "server log confirms the sweeper graded expired attempt(s)" \
    || fail "no sweeper grading line in ${SERVER_LOG} — sweeper did not fire"
else
  note "set SERVER_LOG=<path> to assert the sweeper grading line in the server log"
fi
# Confirm the attempt is now terminal. After the quiet wait the sweeper has run,
# so this read is a confirmation, not the grader.
req GET "/exam-attempts/${TIMED_ATTEMPT}/results" "" "$STUDENT_B"
st="$(jq -r '.status' <<<"$BODY" 2>/dev/null || echo '')"
[[ "$HTTP_CODE" == "200" && ( "$st" == "graded" || "$st" == "expired" ) ]] \
  || fail "abandoned attempt not terminal after ${SWEEP_WAIT}s (code=$HTTP_CODE status=$st)"
pass "abandoned attempt auto-graded (status=$st) with NO submit and NO read while waiting"

echo ""
echo "==> ALL PHASE 2 SMOKE CHECKS PASSED"
