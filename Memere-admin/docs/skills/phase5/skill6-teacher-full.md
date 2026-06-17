# Phase 5 · Skill 6 — Full Teacher Authoring Portal

> **Context:** The original design spec (§1.2, §8) explicitly excluded teachers
> from this panel — they were meant to use the Flutter mobile app. This skill
> extends the panel to give teachers a full web-based authoring experience covering
> everything the backend exposes for the teacher role.

---

## Complete teacher feature matrix

| Area | Endpoints | Status |
|---|---|---|
| Auth (login/logout) | POST /auth/login, POST /auth/logout | ✅ Done (skill 1) |
| My Courses (list/create/edit/publish/delete) | GET+POST /courses, PUT+DELETE /courses/:id, POST /courses/:id/publish | ✅ Done (skill 3) |
| Earnings dashboard | GET /me/earnings, GET /courses/:id/sales | ✅ Done (skill 4) |
| Sections (list/add) | GET+POST /courses/:id/sections | ✅ Done (skill 5) |
| **Lessons** (list/add) | GET /sections/:id/lessons, POST /sections/:id/lessons | ❌ Missing |
| **Quizzes** (create/update/add questions) | POST /courses/:id/quizzes, PUT /quizzes/:id, POST /quizzes/:id/questions | ❌ Missing |
| **Exams** (create/add questions/publish/stats) | POST /courses/:id/exams, POST /exams/:id/questions, POST /exams/:id/publish, GET /exams/:id/stats | ❌ Missing |
| **Video pipeline** (upload-url/status/confirm/retry) | POST /lessons/:id/videos/upload-url, GET /videos/:id/status, POST /videos/:id/confirm, POST /videos/:id/retry | ❌ Missing |

---

## Tasks

### 6.1 — Lessons: list and add

**Schemas** (`lib/api/schemas.ts`):
```ts
LessonSchema: { id, section_id, title, type, order_index, duration_seconds,
  is_free_preview, is_published, video_id?, created_at, updated_at }
LessonListResponseSchema: { data: LessonSchema[] }
AddLessonInputSchema: { title, type: enum["video","text","pdf"],
  is_free_preview, duration_seconds?, is_published? }
```

**Endpoints** (`lib/api/endpoints.ts`):
- `listLessons(sectionId)` → GET /sections/:id/lessons
- `addLesson(sectionId, input)` → POST /sections/:id/lessons

**Route Handlers**:
- `app/api/teacher/sections/[id]/lessons/route.ts` — GET + POST (teacher-only)

**UI** (`components/courses/lessons-list.tsx`):
- Render inside course detail page under each section (accordion or flat list below sections)
- Each section row expandable to show its lessons
- "Add Lesson" button per section → dialog (title, type select, free preview toggle, duration)
- Lesson row: order, title, type badge (video/text/pdf), duration, published badge, video status if video type

### 6.2 — Quizzes: create, update, add questions

**Schemas**:
```ts
QuizSchema: { id, course_id, title, time_limit_seconds, pass_percentage,
  randomize_questions, max_attempts, question_count?, created_at }
QuizListResponseSchema: { data: QuizSchema[] }  // or possibly bare array — check
CreateQuizInputSchema: { title, time_limit_seconds, pass_percentage,
  randomize_questions, max_attempts }
QuizQuestionSchema: { id, text, type, points, explanation?,
  order_index, subject?, topic?, answers: AnswerSchema[] }
AnswerSchema: { text, is_correct }
AddQuizQuestionInputSchema: { text, type: enum["multiple_choice","true_false"],
  points, explanation?, order_index, subject?, topic?, answers: AnswerSchema[] }
```

**Endpoints**:
- `createQuiz(courseId, input)` → POST /courses/:id/quizzes
- `updateQuiz(quizId, input)` → PUT /quizzes/:id
- `addQuizQuestion(quizId, input)` → POST /quizzes/:id/questions
- `listQuizzes(courseId)` — if endpoint exists, else show from course detail

**Route Handlers**:
- `app/api/teacher/courses/[id]/quizzes/route.ts` — POST
- `app/api/teacher/quizzes/[id]/route.ts` — PUT
- `app/api/teacher/quizzes/[id]/questions/route.ts` — POST

**UI** (within course detail `/my-courses/[id]`):
- "Quizzes" tab or section below sections list
- Quiz row: title, questions count, pass%, time limit, edit button
- "Add Quiz" → dialog (title, time limit, pass%, randomize toggle, max attempts)
- Click quiz row → quiz detail (inline or `/my-courses/:id/quizzes/:qid`)
- Add question form: text, type, points, explanation, answer options with is_correct checkbox

### 6.3 — Exams: create, add questions, publish, stats

**Schemas**:
```ts
ExamSchema: { id, course_id, title, subject, grade, duration_minutes,
  pass_marks, instructions?, is_published, created_at }
CreateExamInputSchema: { title, subject, grade, duration_minutes,
  pass_marks, instructions? }
ExamQuestionInputSchema: { question_id?, marks, order_index,
  text?, type?, answers?: AnswerSchema[] }
ExamStatsSchema: { exam_id, total_attempts, avg_score, pass_rate,
  highest_score, lowest_score }
```

**Endpoints**:
- `createExam(courseId, input)` → POST /courses/:id/exams
- `addExamQuestion(examId, input)` → POST /exams/:id/questions
- `publishExam(examId)` → POST /exams/:id/publish
- `getExamStats(examId)` → GET /exams/:id/stats

**Route Handlers**:
- `app/api/teacher/courses/[id]/exams/route.ts` — POST
- `app/api/teacher/exams/[id]/questions/route.ts` — POST
- `app/api/teacher/exams/[id]/publish/route.ts` — POST
- `app/api/teacher/exams/[id]/stats/route.ts` — GET

**UI**: Exams section in course detail (similar layout to Quizzes). Stats shown
in an inline card when exam is published.

### 6.4 — Video pipeline: upload-url, confirm, status, retry

The video pipeline is 4 steps:
1. Request presigned upload URL → get `upload_url` + `video_id`
2. Upload the file directly to S3/MinIO using the presigned URL (from the browser)
3. Confirm the upload → triggers backend transcode job
4. Poll status → `pending | processing | ready | failed`; retry if failed

**Schemas**:
```ts
VideoUploadRequestSchema: { file_name, content_type, size_bytes }
VideoUploadResponseSchema: { video_id, upload_url, expires_at }
VideoStatusSchema: { id, lesson_id?, status, duration_seconds?,
  thumbnail_url?, hls_url?, created_at }
```

**Endpoints** (server-side for URL request/confirm/retry; status can be client polled):
- `requestVideoUpload(lessonId, input)` → POST /lessons/:id/videos/upload-url
- `confirmVideoUpload(videoId)` → POST /videos/:id/confirm
- `retryTranscode(videoId)` → POST /videos/:id/retry
- `getVideoStatus(videoId)` → GET /videos/:id/status

**Route Handlers**:
- `app/api/teacher/lessons/[id]/videos/upload-url/route.ts` — POST
- `app/api/teacher/videos/[id]/confirm/route.ts` — POST
- `app/api/teacher/videos/[id]/retry/route.ts` — POST
- `app/api/teacher/videos/[id]/status/route.ts` — GET

**UI** (`components/courses/video-uploader.tsx`):
- Shown on lesson rows with type="video"
- File input → POST to Route Handler to get presigned URL → PUT the file directly
  to the presigned URL from the browser (this is fine — the presigned URL is
  MinIO/S3, not the Go API; no token leaks)
- Show upload progress bar (XHR with `onprogress`)
- After upload → POST confirm → poll status every 5s until `ready` or `failed`
- `failed` → show "Retry transcode" button

---

## Navigation update

Add to the teacher detail page (`/my-courses/[id]`) tabbed or sectioned layout:
```
Course: [title]
├── Details (metadata)
├── Sections & Lessons
├── Quizzes
└── Exams
```

---

## Definition of Done

- [ ] Teacher can add lessons to a section (all types: video, text, pdf)
- [ ] Teacher can upload a video for a video lesson and see it reach "ready" status
- [ ] Teacher can create a quiz with questions and multiple-choice answers
- [ ] Teacher can create an exam with questions and publish it
- [ ] Teacher can view exam stats once published
- [ ] All Route Handlers verify `user.role === "teacher"` (403 otherwise)
- [ ] `pnpm build` + `pnpm typecheck` pass
