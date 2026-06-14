package notification

import (
	"context"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// Hooks is the typed facade wired into every bounded domain as its notification
// hook. Each method builds a NotifyEvent and delegates to the central Dispatcher.
// It satisfies the narrow notifier interfaces in every bounded domain package:
//   - payment.Notifier   (PurchaseConfirmed)
//   - progress.Notifier  (CertificateReady)
//   - worker.SubscriptionNotifier (SubscriptionExpired)
//   - worker.NotifyPort  (VideoReady)
//   - grading.Notifier   (ExamGraded)
//
// Methods are best-effort: errors from Dispatcher.Notify are silently swallowed
// because provider failure must never surface to the caller's request.
type Hooks struct {
	n service.Notifier
}

func NewHooks(n service.Notifier) *Hooks { return &Hooks{n: n} }

// PurchaseConfirmed satisfies payment.Notifier.
func (h *Hooks) PurchaseConfirmed(ctx context.Context, p *entity.Payment) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		UserID: p.StudentID.String(),
		Type:   "purchase_confirmed",
		Title:  "Payment received",
		Body:   "Your enrollment is confirmed. Happy studying!",
		Data:   map[string]string{"payment_id": p.ID.String()},
		Channels: []service.Channel{
			service.ChannelEmail,
			service.ChannelInApp,
		},
	})
}

// CertificateReady satisfies progress.Notifier.
func (h *Hooks) CertificateReady(ctx context.Context, studentID, courseID uuid.UUID) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		UserID: studentID.String(),
		Type:   "certificate_ready",
		Title:  "Course complete!",
		Body:   "Your certificate is ready. Download it from your profile.",
		Data: map[string]string{
			"course_id": courseID.String(),
		},
		Channels: []service.Channel{
			service.ChannelPush,
			service.ChannelEmail,
			service.ChannelInApp,
		},
	})
}

// SubscriptionExpired satisfies worker.SubscriptionNotifier.
func (h *Hooks) SubscriptionExpired(ctx context.Context, sub *entity.Subscription) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		UserID: sub.StudentID.String(),
		Type:   "subscription_expired",
		Title:  "Subscription expired",
		Body:   "Your subscription has expired. Renew to keep learning.",
		Data:   map[string]string{"subscription_id": sub.ID.String()},
		Channels: []service.Channel{
			service.ChannelPush,
			service.ChannelEmail,
			service.ChannelInApp,
		},
	})
}

// VideoReady satisfies worker.NotifyPort. teacherID and courseID are strings
// (the worker passes UUIDs as strings).
func (h *Hooks) VideoReady(ctx context.Context, videoID, courseID string) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		// Video-ready events target the teacher who uploaded the video. In Phase 5
		// we don't have a direct teacher-lookup, so we use courseID as the recipient
		// hint and let the worker resolve the real userID via the course owner.
		// For now we set UserID to courseID as a placeholder and mark it teacher-facing.
		// Skill 5 wiring replaces this with the real teacher lookup.
		UserID: courseID,
		Type:   "video_ready",
		Title:  "Video processed",
		Body:   "Your video is ready for students.",
		Data: map[string]string{
			"video_id":  videoID,
			"course_id": courseID,
		},
		Channels: []service.Channel{service.ChannelPush, service.ChannelInApp},
	})
}

// ExamGraded satisfies grading.Notifier.
func (h *Hooks) ExamGraded(ctx context.Context, studentID uuid.UUID, attemptID uuid.UUID, score float64) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		UserID: studentID.String(),
		Type:   "exam_graded",
		Title:  "Exam graded",
		Body:   "Your exam has been graded. Check your results.",
		Data: map[string]string{
			"attempt_id": attemptID.String(),
		},
		Channels: []service.Channel{
			service.ChannelPush,
			service.ChannelInApp,
		},
	})
}

// StreakWarning satisfies worker.StreakNotifier.
// Sent when a student hasn't studied in 2 days and their streak is at risk.
func (h *Hooks) StreakWarning(ctx context.Context, studentID uuid.UUID) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		UserID: studentID.String(),
		Type:   "streak_warning",
		Title:  "Your streak is at risk!",
		Body:   "You haven't studied in 2 days. Keep your streak alive!",
		Data:   map[string]string{"student_id": studentID.String()},
		Channels: []service.Channel{
			service.ChannelPush,
			service.ChannelInApp,
		},
	})
}

// LessonPublished satisfies course.Notifier (teacher publishes a lesson).
func (h *Hooks) LessonPublished(ctx context.Context, lessonID, courseID uuid.UUID) {
	_ = h.n.Notify(ctx, service.NotifyEvent{
		// Enrolled students should receive this. Skill 5 HTTP wiring resolves the
		// enrolled student list; here we send to courseID as a placeholder.
		UserID: courseID.String(),
		Type:   "lesson_published",
		Title:  "New lesson available",
		Body:   "A new lesson has been published in your course.",
		Data: map[string]string{
			"lesson_id": lessonID.String(),
			"course_id": courseID.String(),
		},
		Channels: []service.Channel{service.ChannelPush, service.ChannelInApp},
	})
}
