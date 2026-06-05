package entity

// Permission is a discrete capability checked by the authorization layer. The
// set mirrors the RBAC matrix in spec §7.2.
type Permission string

const (
	// Student capabilities.
	PermReadEnrolledCourses Permission = "course:read_enrolled"
	PermSubmitAttempt       Permission = "attempt:submit"
	PermViewOwnProgress     Permission = "progress:view_own"
	PermPurchase            Permission = "purchase"

	// Teacher capabilities (a teacher also has every student capability).
	PermCreateOwnCourse        Permission = "course:create_own"
	PermEditOwnCourse          Permission = "course:edit_own"
	PermViewOwnCourseAnalytics Permission = "analytics:view_own_course"

	// Admin capabilities (an admin has every capability — see Can).
	PermManageAllUsers   Permission = "users:manage_all"
	PermManageAllCourses Permission = "courses:manage_all"
	PermManagePayments   Permission = "payments:manage"
)

// studentPerms are granted to every authenticated role.
var studentPerms = map[Permission]bool{
	PermReadEnrolledCourses: true,
	PermSubmitAttempt:       true,
	PermViewOwnProgress:     true,
	PermPurchase:            true,
}

// teacherPerms are granted to teachers in addition to studentPerms.
var teacherPerms = map[Permission]bool{
	PermCreateOwnCourse:        true,
	PermEditOwnCourse:          true,
	PermViewOwnCourseAnalytics: true,
}

// Can reports whether the role is permitted to perform p. Admin can do
// anything; teacher has teacher+student permissions; student has student
// permissions. This is the coarse role gate — resource-level ownership checks
// (e.g. "this teacher owns this course") are enforced separately by the
// usecase layer in later skills.
func (r Role) Can(p Permission) bool {
	switch r {
	case RoleAdmin:
		return true
	case RoleTeacher:
		return teacherPerms[p] || studentPerms[p]
	case RoleStudent:
		return studentPerms[p]
	default:
		return false
	}
}

// HasRole reports whether r is one of the allowed roles. The middleware uses
// this for endpoints gated purely by role rather than by a specific permission.
func (r Role) HasRole(allowed ...Role) bool {
	for _, a := range allowed {
		if r == a {
			return true
		}
	}
	return false
}
