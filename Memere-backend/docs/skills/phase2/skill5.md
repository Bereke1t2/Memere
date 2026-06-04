# Phase 2 · Skill 5 — HTTP Delivery & Wiring (Quiz, Exam, Analytics)

> **Prerequisite:** Phase 2 Skills 1–4 done (data layer, quiz engine, exam engine,
> sweeper + analytics — all tested in isolation). Read
> [`docs/skill.md`](../../skill.md) §2 and reuse the Phase 1 delivery patterns
> (DTOs, response helpers, middleware stack, router, wiring).
>
> **Spec references:** `memere_Design_Specification.md` §5.4 (API conventions),
> README "Quiz & Exam Endpoints", §9 (all rules), §7.2 (RBAC).

---

## Goal

Expose the quiz, exam, and analytics functionality over the existing Gin API and
wire the background sweeper into the app lifecycle. After this skill, Phase 2 is
**fully runnable end-to-end**: a teacher builds a quiz/exam, a student takes it
under a server-enforced timer, submits, and gets graded results + analytics — with
correct answers never crossing the wire.

This skill **extends** the Phase 1 delivery layer; it does not rebuild it. Reuse
`apperror` response mapping, the middleware stack, `actorFromContext`, the
paginated envelope, and the constructor-wiring style in `cmd/api/main.go`.

---

## API surface for Phase 2 (from README + spec)

Base path `/api/v1`. The **authoring** routes are teacher/admin; the **taking**
routes are student (enrolled). Optional-auth where listing published content.

**Quiz authoring** (`quiz_handler.go`)
| Method | Path | Auth |
|---|---|---|
| POST | `/courses/:id/quizzes` | bearer + owner/admin |
| POST | `/quizzes/:id/questions` | bearer + owner/admin |
| POST | `/questions/:id/answers` | bearer + owner/admin |
| PUT | `/quizzes/:id` | bearer + owner/admin |
| POST | `/quizzes/:id/publish` | bearer + owner/admin |

**Quiz taking**
| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/quizzes/:id` | bearer (enrolled) | metadata + question count, NO answers |
| POST | `/quizzes/:id/attempts` | bearer | start attempt → questions (no keys) + remaining time |
| PATCH | `/quiz-attempts/:id` | bearer (owner) | auto-save answers (every ~30s) |
| POST | `/quiz-attempts/:id/submit` | bearer (owner) | grade + feedback |
| GET | `/quiz-attempts/:id/result` | bearer (owner) | result (post-submit only) |

**Exam** (`exam_handler.go`) — mirrors quiz
| Method | Path | Auth |
|---|---|---|
| POST | `/courses/:id/exams` | bearer + owner/admin |
| POST | `/exams/:id/questions` | bearer + owner/admin |
| GET | `/mock-exams` | optional bearer (published only) |
| POST | `/mock-exams/:id/start` | bearer |
| PATCH | `/exam-attempts/:id` | bearer (owner) |
| POST | `/exam-attempts/:id/submit` | bearer (owner) |
| GET | `/exam-attempts/:id/results` | bearer (owner) |

**Analytics** (`analytics_handler.go`)
| Method | Path | Auth |
|---|---|---|
| GET | `/exam-attempts/:id/analytics` | bearer (owner) |
| GET | `/me/trend?subject=` | bearer |
| GET | `/exams/:id/stats` | bearer + owner/admin |

(Video/payment/progress/notification endpoints remain later phases.)

---

## Tasks

### 5.1 — DTOs (answer-key-free by construction)

Add request/response DTOs (with `json` tags — delivery layer only):

- **Critical:** the client-facing question DTO (`QuestionClientResponse`) has
  **no** `is_correct` field and **no** correct-answer field. Build it only from
  the Skill 1 `GetQuestionsForClient` data. The compiler must make leaking
  impossible — there is simply no field to populate.
- Quiz: `CreateQuizRequest`, `AddQuestionRequest`, `AddAnswerRequest`,
  `QuizClientResponse`, `StartAttemptResponse` (questions + `expires_at` +
  `remaining_seconds`), `SaveProgressRequest`, `SubmitAttemptRequest`,
  `AttemptResultResponse` (score, percentage, passed, per-question feedback +
  `explanation` + correct answer **only here, post-grade**).
- Exam: analogous DTOs.
- Analytics: `AttemptAnalyticsResponse` (subject breakdown, weak areas,
  percentile), `TrendResponse`, `ExamStatsResponse`.

### 5.2 — Handlers

- `quiz_handler.go`, `exam_handler.go`, `analytics_handler.go` — bind+validate,
  pull `actor` from context, call the Phase 2 usecases, map to DTOs, set status
  codes (`201` start/create, `200` submit/result/save).
- Parse `:id` as UUID with clean `400 INVALID_ID`.
- Auto-save (`PATCH .../attempts/:id`) returns `204 No Content` or `200` with
  remaining time — pick one and be consistent.
- All ownership/state errors surface through the standard `apperror` envelope
  (e.g. `INVALID_ATTEMPT_STATE` → `409`, IDOR → `403`).

### 5.3 — Routes & RBAC

Register the routes above in the existing router, attaching:
- `RequireAuth` + `RequireRole(teacher, admin)` on authoring routes.
- `RequireAuth` on taking routes (ownership enforced in the usecase, not
  middleware).
- `OptionalAuth` on `GET /mock-exams` (published visibility).
Reuse the Phase 1 middleware verbatim — do not fork it.

### 5.4 — Wire the sweeper + new deps into `main.go`

Extend the Phase 1 wiring:
- Construct quiz/exam/attempt repos, redis attempt-state + leaderboard stores,
  the grading package, and the quiz/exam/analytics usecases.
- Construct handlers; register routes.
- **Start the attempt sweeper** (Skill 4) in its own goroutine after the server
  starts; ensure it stops on the same `SIGINT/SIGTERM` shutdown path (respect
  `SWEEPER_ENABLED`).
- Keep `main` thin; factor an `app.Build(cfg) (*App, error)` if it grows.

### 5.5 — End-to-end smoke test (`scripts/smoke_phase2.sh`)

Against a running stack (extends Phase 1's smoke):
1. login as teacher; create a course (Phase 1) → create a quiz with 3 questions
   (1 each type) → publish.
2. login as student; `GET /quizzes/:id` → returns count, **no answer keys**
   (assert the JSON has no `is_correct`).
3. start attempt → returns questions in randomized order, `remaining_seconds`
   present, **no keys**.
4. auto-save partial answers; submit → graded result with score + explanations.
5. fetch result → matches; fetch analytics → subject breakdown + percentile.
6. create a timed exam (short duration); start; **wait past expiry without
   submitting**; confirm the sweeper auto-grades it to `expired/graded`
   (poll the result endpoint).
7. attempt to exceed `max_attempts` → rejected.
8. IDOR: student B tries to read student A's attempt → `403`.

### 5.6 — Docs

- Update `README.md` endpoint tables / "Getting Started" if commands changed.
- Optionally extend the `api/` OpenAPI stub with the Phase 2 endpoints.

---

## Definition of Done

- [ ] `go build ./...` clean; `golangci-lint run` clean; `go test ./...` passes.
- [ ] `make up && make migrate-up && make run` serves the Phase 2 endpoints.
- [ ] `scripts/smoke_phase2.sh` passes every assertion, including the
      answer-key-absence checks and the sweeper auto-grade step.
- [ ] No client-facing response on any quiz/exam path contains `is_correct` or a
      correct-answer field before submission (grep + smoke assertion).
- [ ] Server-enforced timer proven: an abandoned attempt is graded by the sweeper
      with no client call.
- [ ] `INVALID_ATTEMPT_STATE` (409) and IDOR (403) returned correctly via the
      standard error envelope.
- [ ] Sweeper starts with the app and stops cleanly on shutdown.

## Verification commands

```bash
make up && make migrate-up
make run &
bash scripts/smoke_phase2.sh
go test ./... && golangci-lint run
# Final leak sweep across the whole delivery layer:
grep -rn 'is_correct\|IsCorrect\|correct_answer' internal/delivery/http/dto 2>/dev/null \
  && echo "FAIL: answer key in a DTO" || echo "OK: no answer keys in client DTOs"
```

---

## 🎉 Phase 2 complete — what now

When this Definition of Done passes, the quiz & exam engines are live: server-side
grading, server-enforced timers (sync + background sweeper), randomized
answer-key-free delivery, attempt limits, and §9.3 analytics — all on top of the
Phase 1 foundation and honoring every Non-Negotiable.

**To proceed to Phase 3:**
1. Report "Phase 2 complete" to the user with a one-paragraph summary + the
   smoke-test result.
2. Ask Claude to author the **Phase 3** skill files into `docs/skills/phase3/` —
   the **Video pipeline**: pre-signed S3 upload URLs, the HLS transcode
   trigger/metadata flow, and streaming/download URL issuance (spec §8), plus the
   `videos` table Go layer left stubbed in Phase 1.
3. Do **not** scaffold Phase 3 code before its skills are written and reviewed.
