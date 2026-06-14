package admin

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// ListAllCourses returns every course (including unpublished, any teacher).
// Admin only.
func (s *Service) ListAllCourses(ctx context.Context, actor Actor, filter repository.CourseFilter, cursor *pagination.Cursor, limit int) ([]*entity.Course, *pagination.Cursor, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.courses.List(ctx, filter, cursor, limit)
}

// UnpublishCourse sets is_published=false and notifies the teacher. Admin only.
// Audited.
func (s *Service) UnpublishCourse(ctx context.Context, actor Actor, courseID uuid.UUID, reason string) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	c, err := s.courses.FindByID(ctx, courseID)
	if err != nil {
		return err
	}
	if !c.IsPublished {
		return nil // idempotent
	}
	c.IsPublished = false
	if err := s.courses.Update(ctx, c); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "course.unpublish", "course", &courseID, map[string]any{"reason": reason})

	// Notify the course teacher best-effort.
	if s.notify != nil {
		_ = s.notify.Notify(ctx, service.NotifyEvent{
			UserID: c.TeacherID.String(),
			Type:   "course_unpublished",
			Title:  "Course unpublished",
			Body:   "Your course has been unpublished by an admin: " + reason,
			Data:   map[string]string{"course_id": courseID.String()},
			Channels: []service.Channel{
				service.ChannelPush,
				service.ChannelInApp,
			},
		})
	}
	return nil
}

// DeleteCourse admin-soft-deletes any course. Admin only. Audited.
func (s *Service) DeleteCourse(ctx context.Context, actor Actor, courseID uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if err := s.courses.SoftDelete(ctx, courseID); err != nil {
		return err
	}
	s.writeAudit(ctx, actor, "course.delete", "course", &courseID, nil)
	return nil
}
