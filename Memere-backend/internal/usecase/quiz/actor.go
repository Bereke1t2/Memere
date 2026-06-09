// Package quiz holds the quiz-engine usecases (spec §9.1): teacher authoring of
// quizzes/questions/answers, and the student attempt lifecycle — start, auto-
// save, submit, and server-side grading with feedback. Like the other usecase
// packages it is pure orchestration over the domain repository interfaces; it
// imports no infrastructure or delivery code (dependency rule).
//
// Two phase-defining non-negotiables are enforced here and must never be
// relaxed: the correct-answer key never leaves the server (client views are
// built from answer-key-free queries), and grading is server-side only (the
// client-reported score is ignored).
package quiz

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// Actor is the authenticated caller's identity passed into every usecase
// method. A nil *Actor is an anonymous request, which may read published
// content only and may never start or submit an attempt.
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
