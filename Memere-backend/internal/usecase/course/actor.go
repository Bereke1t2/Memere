// Package course holds the course-domain usecases: course CRUD, section and
// lesson management, and assembling the nested course view. Like the auth
// package it is pure orchestration over the domain repository interfaces — no
// infrastructure, no delivery code (dependency rule). The caller's identity
// arrives explicitly as an Actor (from the JWT middleware in Skill 5), never
// inferred from client-sent IDs, so ownership is enforced here (spec §7.2).
package course

import (
	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// Actor is the authenticated caller's identity passed into every usecase
// method. A nil *Actor represents an anonymous (unauthenticated) request, which
// may read published content only.
type Actor struct {
	UserID uuid.UUID
	Role   entity.Role
}

// isAdmin reports whether the actor is an admin (nil actor is never admin).
func (a *Actor) isAdmin() bool {
	return a != nil && a.Role == entity.RoleAdmin
}

// isTeacherOrAdmin reports whether the actor may author courses.
func (a *Actor) isTeacherOrAdmin() bool {
	return a != nil && (a.Role == entity.RoleTeacher || a.Role == entity.RoleAdmin)
}

// canSeeUnpublished reports whether the actor may view unpublished content for
// the given course: the owning teacher, or any admin.
func (a *Actor) canSeeUnpublished(c *entity.Course) bool {
	if a == nil {
		return false
	}
	return a.Role == entity.RoleAdmin || (a.Role == entity.RoleTeacher && c.TeacherID == a.UserID)
}
