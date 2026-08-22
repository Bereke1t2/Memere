// Package video holds the Phase 3 video-pipeline usecases: requesting a
// pre-signed upload URL, confirming an upload, reporting processing status, and
// retrying a failed transcode. Like every usecase package it is pure
// orchestration over domain ports (VideoRepository, LessonRepository,
// CourseRepository, ObjectStore, JobQueue) — no HTTP, no AWS SDK, no SQL. The
// caller's identity arrives explicitly as an *Actor so course ownership is
// enforced here and client-sent IDs are never trusted (IDOR rule, spec §7.2).
package video

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// Actor is the authenticated caller passed into every usecase method. A nil
// *Actor is an anonymous request. Uploads/mutations still reject anonymous
// callers (teacher/admin only); reads (status/stream/download) now tolerate a
// nil actor so unregistered guests can watch published/free content, with the
// access check (assertCanWatch) doing the real gating.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// GuestUserID is the system identifier bound to a download token issued to an
// anonymous / unregistered caller. Mirrors quiz.GuestUserID and exam.GuestUserID.
var GuestUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func (a *Actor) isAdmin() bool {
	return a != nil && a.Role == entity.RoleAdmin
}

// isTeacherOrAdmin reports whether the actor may author/own course content.
func (a *Actor) isTeacherOrAdmin() bool {
	return a != nil && (a.Role == entity.RoleTeacher || a.Role == entity.RoleAdmin)
}
