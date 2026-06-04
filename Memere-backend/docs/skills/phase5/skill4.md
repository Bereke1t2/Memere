# Phase 5 · Skill 4 — Certificates & Engagement Sweeper

> **Prerequisite:** Phase 5 Skills 1–3 done (progress, notifications, admin). Read
> [`docs/skill.md`](../../skill.md) §2.
>
> **Spec references:** `memere_Design_Specification.md` §1.4 (FR-10 certificates),
> §11.2 (`streak_warning`, `certificate_ready`), §8.3 (pre-signed URL pattern —
> reused for certificate download), README "Certificates".

---

## Goal

Two remaining engine pieces before Phase 5's API:

1. **Certificates** — when a student completes a course (100%), generate a
   downloadable completion certificate (PDF), store it, and expose it via a
   short-lived signed URL (reusing the Phase 3 `ObjectStore`/signer).
2. **Engagement sweeper** — a scheduled job that fires `streak_warning`
   notifications to students who haven't studied in 2 days (spec §11.2) and
   maintains streak hygiene.

---

## Tasks

### 4.1 — Migration (`migrations/0016_certificates`)

```sql
-- migrations/0016_certificates.up.sql
CREATE TABLE IF NOT EXISTS courses.certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES auth.users(id),
    course_id  UUID NOT NULL REFERENCES courses.courses(id),
    serial     TEXT NOT NULL UNIQUE,         -- public verification code
    file_key   TEXT,                          -- S3 key of the rendered PDF
    issued_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (student_id, course_id)            -- one cert per course completion
);
CREATE INDEX IF NOT EXISTS certificates_student_idx ON courses.certificates(student_id);
```

### 4.2 — Certificate entity + repo

`entity.Certificate` (no tags), `CertificateRepository`
(`CreateIfNotExists(student, course, serial)`, `GetByID`,
`GetByStudentCourse`, `GetBySerial`, `SetFileKey`, `ListByStudent`).

Serial = a non-guessable code (e.g. `MEM-{year}-{base32(random)}`); generated
server-side, used for public verification (4.5).

### 4.3 — Certificate generation (`internal/usecase/certificate/`)

```go
type Service struct {
    repo    repository.CertificateRepository
    progress repository.ProgressRepository
    course  repository.CourseRepository
    users   repository.UserRepository
    store   service.ObjectStore
    signer  storage.URLSigner
    render  Renderer   // PDF renderer port
    notify  notification.Hooks
    cfg     Config
}

// Issue is called when a course hits 100% (from Phase 5 Skill 1 MarkLessonComplete)
// or on-demand by the student. Idempotent: one certificate per (student, course).
func (s *Service) Issue(ctx context.Context, actor Actor, courseID uuid.UUID) (*entity.Certificate, error) {
    cp, err := s.progress.GetCourseProgress(ctx, actor.UserID, courseID)
    if err != nil { return nil, err }
    if cp.PercentComplete.LessThan(decimal.NewFromInt(100)) {
        return nil, apperror.Conflict("COURSE_NOT_COMPLETE", "")
    }
    serial := s.newSerial()
    cert, created, err := s.repo.CreateIfNotExists(ctx, actor.UserID, courseID, serial)
    if err != nil { return nil, err }
    if created {
        // render async or inline; store PDF; set file_key; notify certificate_ready
        go s.renderAndStore(context.Background(), cert) // or enqueue a job
    }
    return cert, nil
}
```

- `GetDownloadURL(ctx, actor, certID)` — ownership-checked; signed URL (e.g. 1h)
  to the stored PDF via the Phase 3 signer. Regenerate the PDF if `file_key` is
  missing.
- `renderAndStore` — render PDF, `store.Put` under `certificates/{cert_id}.pdf`,
  `SetFileKey`, fire `notify.CertificateReady`.

### 4.4 — PDF renderer port (`internal/infrastructure/pdf/`)

```go
// Renderer abstracts PDF generation so the impl is swappable.
type Renderer interface {
    Certificate(ctx context.Context, data CertData) ([]byte, error)
}
type CertData struct { StudentName, CourseTitle, Serial, IssuedDate string }
```

Impl with a pure-Go library (e.g. `jung-kurt/gofpdf` or `go-pdf/fpdf`) — no
headless-browser dependency. Keep a simple template (logo, name, course, date,
serial). Add the lib to `go.mod`.

### 4.5 — Public verification

`VerifyCertificate(ctx, serial) (*PublicCertView, error)` — **public** (no auth):
given a serial, return `{ student_name, course_title, issued_at, valid: true }` or
not-found. No PII beyond name + course (document this privacy choice). Used by a
future "verify this certificate" public page.

### 4.6 — Engagement sweeper (`internal/worker/engagement_sweeper.go`)

```go
func (w *EngagementSweeper) tick(ctx context.Context) error {
    cutoff := w.clock.Now().AddDate(0,0,-2)
    // students whose last_study_date < cutoff and current_streak > 0:
    atRisk, err := w.progress.ListStreakAtRisk(ctx, cutoff) // paged
    for _, s := range atRisk {
        w.notify.StreakWarning(ctx, s.StudentID) // push (§11.2)
    }
    return err
}
```

- Runs daily (`ENGAGEMENT_SWEEP_INTERVAL`, default 24h, off-peak hour). Paged —
  no unbounded scan. Idempotent (don't double-warn the same day; track via a
  `last_warned_date` or check the notification log). Starts/stops with the app.

### 4.7 — Tests

Cert: issue blocked under 100% (`COURSE_NOT_COMPLETE`); idempotent (one per
student/course); download URL ownership-checked + signed; verification by serial
returns minimal data; unknown serial → not found. Sweeper: at-risk students warned
once per day (fake clock); no warning if studied within 2 days; paged.

---

## Definition of Done

- [ ] `make migrate-up` applies `0016`; down reverses; unique
      `(student_id, course_id)` + unique `serial`.
- [ ] `go build ./...` clean; `golangci-lint run` clean.
- [ ] `go test ./internal/usecase/certificate/... ./internal/worker/...` passes.
- [ ] Certificate issuance is idempotent and gated on 100% completion.
- [ ] Certificate download is ownership-checked and via a short-lived signed URL
      (no public key exposure); PDF renders with name/course/serial/date.
- [ ] Public verification returns minimal, non-sensitive data by serial.
- [ ] Engagement sweeper warns at-risk students once/day, paged, idempotent, and
      stops with the app.

## Verification commands

```bash
make migrate-up && go build ./... && golangci-lint run
go test ./internal/usecase/certificate/... ./internal/worker/... -v
```

## Hand-off to Skill 5

All Phase 5 engines exist. Skill 5 exposes progress, notifications, admin, and
certificate APIs over HTTP, wires the notification + engagement sweepers into
`main.go`, and runs the Phase 5 end-to-end smoke test.
