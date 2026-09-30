package admin

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// BroadcastSegment identifies which users receive an announcement.
type BroadcastSegment string

const (
	SegmentAll         BroadcastSegment = "all"
	SegmentStudents    BroadcastSegment = "students"
	SegmentTeachers    BroadcastSegment = "teachers"
	SegmentSubscribers BroadcastSegment = "subscribers"
)

// BroadcastInput carries the announcement payload and target segment.
type BroadcastInput struct {
	Title   string
	Body    string
	Segment BroadcastSegment
	Data    map[string]string
}

const broadcastPageSize = 100

// Broadcast sends an announcement to every user in the segment, paging through
// recipients in batches so the in-memory list is always bounded. Admin only.
// Audited.
func (s *Service) Broadcast(ctx context.Context, actor Actor, input BroadcastInput) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}

	total := 0
	switch input.Segment {
	case SegmentAll:
		total += s.broadcastToRole(ctx, input, entity.RoleStudent)
		total += s.broadcastToRole(ctx, input, entity.RoleTeacher)
	case SegmentStudents:
		total += s.broadcastToRole(ctx, input, entity.RoleStudent)
	case SegmentTeachers:
		total += s.broadcastToRole(ctx, input, entity.RoleTeacher)
	case SegmentSubscribers:
		total += s.broadcastToSubscribers(ctx, input)
	}

	if s.notifications != nil {
		_, _ = s.notifications.SaveAnnouncement(ctx, &entity.Announcement{
			SenderID:       actor.UserID,
			Title:          input.Title,
			Body:           input.Body,
			Segment:        string(input.Segment),
			Data:           input.Data,
			RecipientCount: total,
		})
	}

	if s.notify != nil {
		_ = s.notify.Notify(ctx, service.NotifyEvent{
			UserID:   actor.UserID.String(),
			Type:     "announcement",
			Title:    input.Title,
			Body:     input.Body,
			Data:     input.Data,
			Channels: []service.Channel{service.ChannelInApp},
		})
	}

	log.Printf("admin.Broadcast: sent to %d recipient(s) for segment=%s", total, input.Segment)
	s.writeAudit(ctx, actor, "broadcast.send", "system", nil, map[string]any{
		"segment":    string(input.Segment),
		"recipients": total,
		"title":      input.Title,
		"body":       input.Body,
	})
	return nil
}

// ListAnnouncements lists broadcast announcements with aggregated recipient counts. Admin only.
func (s *Service) ListAnnouncements(ctx context.Context, actor Actor, limit int) ([]*entity.AnnouncementSummary, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if s.notifications == nil {
		return []*entity.AnnouncementSummary{}, nil
	}
	return s.notifications.ListAnnouncements(ctx, limit)
}

// DeleteAnnouncement deletes all notification rows corresponding to the given announcement. Admin only.
// Audited.
func (s *Service) DeleteAnnouncement(ctx context.Context, actor Actor, id uuid.UUID) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if s.notifications == nil {
		return apperror.NotFound("notification repository unavailable", nil)
	}

	if err := s.notifications.DeleteAnnouncementGroup(ctx, id); err != nil {
		return err
	}

	log.Printf("admin.DeleteAnnouncement: removed announcement id=%s", id)
	s.writeAudit(ctx, actor, "announcement.delete", "announcement", &id, map[string]any{
		"announcement_id": id.String(),
	})
	return nil
}

// broadcastToRole pages through users with the given role and notifies each.
func (s *Service) broadcastToRole(ctx context.Context, input BroadcastInput, role entity.Role) int {
	if s.notify == nil {
		return 0
	}
	roleStr := string(role)
	total := 0
	var cursor *pagination.Cursor
	for {
		users, next, err := s.users.List(ctx, repository.AdminUserFilter{Role: &roleStr}, cursor, broadcastPageSize)
		if err != nil {
			log.Printf("admin.Broadcast: list users role=%s: %v", role, err)
			break
		}
		for _, u := range users {
			_ = s.notify.Notify(ctx, service.NotifyEvent{
				UserID:   u.ID.String(),
				Type:     "announcement",
				Title:    input.Title,
				Body:     input.Body,
				Data:     input.Data,
				Channels: []service.Channel{service.ChannelPush, service.ChannelInApp},
			})
			total++
		}
		if next == nil || len(users) == 0 {
			break
		}
		cursor = next
	}
	return total
}

// broadcastToSubscribers pages through active subscribers.
func (s *Service) broadcastToSubscribers(_ context.Context, input BroadcastInput) int {
	// Active subscriber list query is deferred to Skill 5 DB implementation.
	log.Printf("admin.Broadcast: subscriber segment not yet implemented (title=%q)", input.Title)
	return 0
}

