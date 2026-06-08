// Package exam holds the mock-exam engine usecases (spec §9.2): teacher authoring
// of exams and their question set, and the student attempt lifecycle driven by
// the §9.2 state machine with a server-enforced timer. Like the other usecase
// packages it is pure orchestration over the domain repository interfaces and
// imports no infrastructure or delivery code (dependency rule).
//
// The phase non-negotiables hold here too: the correct-answer key never leaves
// the server (client views are built from answer-key-free queries) and grading is
// server-side only (the shared internal/usecase/grading core).
package exam

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// Actor is the authenticated caller's identity passed into every usecase method.
// A nil *Actor is an anonymous request, which may read published exams only and
// may never start or submit an attempt.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

func (a *Actor) isTeacherOrAdmin() bool {
	return a != nil && (a.Role == entity.RoleTeacher || a.Role == entity.RoleAdmin)
}

// ownsCourse reports whether the actor may author content under the course: the
// owning teacher, or any admin.
func (a *Actor) ownsCourse(c *entity.Course) bool {
	if a == nil {
		return false
	}
	return a.Role == entity.RoleAdmin || (a.Role == entity.RoleTeacher && c.TeacherID == a.UserID)
}
