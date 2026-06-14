-- Certificate queries (Phase 5 Skill 4).

-- CreateCertificateIfNotExists inserts if no row exists for (student, course).
-- name: CreateCertificateIfNotExists :one
INSERT INTO courses.certificates (student_id, course_id, serial)
VALUES ($1, $2, $3)
ON CONFLICT (student_id, course_id) DO NOTHING
RETURNING id, student_id, course_id, serial, file_key, issued_at;

-- GetCertificateByID returns one certificate by primary key.
-- name: GetCertificateByID :one
SELECT id, student_id, course_id, serial, file_key, issued_at
FROM courses.certificates
WHERE id = $1;

-- GetCertificateByStudentCourse returns the cert for a student+course pair.
-- name: GetCertificateByStudentCourse :one
SELECT id, student_id, course_id, serial, file_key, issued_at
FROM courses.certificates
WHERE student_id = $1 AND course_id = $2;

-- GetCertificateBySerial looks up a cert by its public verification code.
-- name: GetCertificateBySerial :one
SELECT id, student_id, course_id, serial, file_key, issued_at
FROM courses.certificates
WHERE serial = $1;

-- SetCertificateFileKey records the S3 key of the rendered PDF.
-- name: SetCertificateFileKey :exec
UPDATE courses.certificates SET file_key = $2 WHERE id = $1;

-- ListCertificatesByStudent returns all certs for a student, newest first.
-- name: ListCertificatesByStudent :many
SELECT id, student_id, course_id, serial, file_key, issued_at
FROM courses.certificates
WHERE student_id = $1
ORDER BY issued_at DESC;
