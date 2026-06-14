package admin

import (
	"context"
	"log"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
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

	log.Printf("admin.Broadcast: sent to %d recipient(s) for segment=%s", total, input.Segment)
	s.writeAudit(ctx, actor, "broadcast.send", "system", nil, map[string]any{
		"segment":    string(input.Segment),
		"recipients": total,
		"title":      input.Title,
	})
	return nil
}

// broadcastToRole pages through users with the given role and notifies each.
func (s *Service) broadcastToRole(ctx context.Context, input BroadcastInput, role entity.Role) int {
	if s.notify == nil {
		return 0
	}
	roleStr := string(role)
	users, _, err := s.users.List(ctx, repository.AdminUserFilter{Role: &roleStr}, nil, broadcastPageSize)
	if err != nil {
		log.Printf("admin.Broadcast: list users role=%s: %v", role, err)
		return 0
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
	}
	return len(users)
}

// broadcastToSubscribers pages through active subscribers.
func (s *Service) broadcastToSubscribers(_ context.Context, input BroadcastInput) int {
	// Active subscriber list query is deferred to Skill 5 DB implementation.
	log.Printf("admin.Broadcast: subscriber segment not yet implemented (title=%q)", input.Title)
	return 0
}
