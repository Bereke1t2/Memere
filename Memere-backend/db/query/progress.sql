-- Progress tracking queries (Phase 5 Skill 1, spec §4.2.8).
-- All reads filter deleted_at IS NULL. Writes use guarded upserts.

-- UpsertLessonProgress creates or updates a lesson progress row.
-- video_progress_seconds uses GREATEST so it only advances (monotonic).
-- completed_at is set once and never cleared.
-- name: UpsertLessonProgress :one
INSERT INTO progress.progress (
    id, student_id, lesson_id, course_id, is_completed,
    completed_at, video_progress_seconds, last_accessed_at
)
VALUES (
    gen_random_uuid(), $1, $2, $3, $4,
    CASE WHEN $4 THEN now() ELSE NULL END,
    $5, now()
)
ON CONFLICT (student_id, lesson_id) WHERE deleted_at IS NULL
DO UPDATE SET
    is_completed          = EXCLUDED.is_completed OR progress.progress.is_completed,
    completed_at          = COALESCE(progress.progress.completed_at, EXCLUDED.completed_at),
    video_progress_seconds = GREATEST(progress.progress.video_progress_seconds, EXCLUDED.video_progress_seconds),
    last_accessed_at      = now(),
    updated_at            = now()
RETURNING *;

-- GetLessonProgress fetches a single lesson progress row for a student.
-- name: GetLessonProgress :one
SELECT * FROM progress.progress
WHERE student_id = $1 AND lesson_id = $2 AND deleted_at IS NULL;

-- ListByCourse returns all lesson progress rows for a student in a course.
-- name: ListByCourse :many
SELECT * FROM progress.progress
WHERE student_id = $1 AND course_id = $2 AND deleted_at IS NULL
ORDER BY updated_at DESC;

-- CountCourseCompletion returns the total published lessons and the number the
-- student has completed, used by RecomputeCourseProgress.
-- name: CountCourseCompletion :one
SELECT
    (SELECT COUNT(*)::int
     FROM courses.lessons l
     WHERE l.course_id = $2 AND l.deleted_at IS NULL AND l.is_published) AS total,
    (SELECT COUNT(*)::int
     FROM progress.progress p
     WHERE p.student_id = $1 AND p.course_id = $2
       AND p.is_completed AND p.deleted_at IS NULL) AS completed;

-- UpsertCourseProgress writes the denormalized course completion rollup.
-- name: UpsertCourseProgress :exec
INSERT INTO progress.course_progress (
    student_id, course_id, completed_lessons, total_lessons,
    percent_complete, completed_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, now())
ON CONFLICT (student_id, course_id)
DO UPDATE SET
    completed_lessons = EXCLUDED.completed_lessons,
    total_lessons     = EXCLUDED.total_lessons,
    percent_complete  = EXCLUDED.percent_complete,
    completed_at      = COALESCE(progress.course_progress.completed_at, EXCLUDED.completed_at),
    updated_at        = now();

-- GetCourseProgress fetches the rollup for one student+course pair.
-- name: GetCourseProgress :one
SELECT * FROM progress.course_progress
WHERE student_id = $1 AND course_id = $2;

-- GetStreak returns the study streak row for a student, or a zero row if none.
-- name: GetStreak :one
SELECT * FROM progress.study_streaks WHERE student_id = $1;

-- UpsertStreak creates or updates the streak row for a student.
-- name: UpsertStreak :exec
INSERT INTO progress.study_streaks (
    student_id, current_streak, longest_streak, last_study_date, updated_at
)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT (student_id)
DO UPDATE SET
    current_streak  = EXCLUDED.current_streak,
    longest_streak  = EXCLUDED.longest_streak,
    last_study_date = EXCLUDED.last_study_date,
    updated_at      = now();
