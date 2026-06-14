package notification

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/apperror"
)

// ---- fake notification repo -------------------------------------------------

type fakeNotificationRepo struct {
	mu    sync.Mutex
	rows  []*entity.Notification
	byID  map[uuid.UUID]*entity.Notification
}

func newFakeNotificationRepo() *fakeNotificationRepo {
	return &fakeNotificationRepo{byID: map[uuid.UUID]*entity.Notification{}}
}

func (f *fakeNotificationRepo) Create(_ context.Context, n *entity.Notification) (*entity.Notification, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	n.CreatedAt = time.Now()
	cp := *n
	f.rows = append(f.rows, &cp)
	f.byID[n.ID] = &cp
	return &cp, nil
}

func (f *fakeNotificationRepo) ListByUser(_ context.Context, userID uuid.UUID, limit int) ([]*entity.Notification, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*entity.Notification
	for _, n := range f.rows {
		if n.UserID == userID {
			cp := *n
			out = append(out, &cp)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakeNotificationRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.Notification, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if n, ok := f.byID[id]; ok {
		cp := *n
		return &cp, nil
	}
	return nil, apperror.NotFound("notification not found", nil)
}

func (f *fakeNotificationRepo) MarkRead(_ context.Context, id, userID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.byID[id]
	if !ok || n.UserID != userID {
		return apperror.NotFound("notification not found", nil)
	}
	if n.ReadAt == nil {
		now := time.Now()
		n.ReadAt = &now
	}
	return nil
}

func (f *fakeNotificationRepo) MarkAllRead(_ context.Context, userID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	for _, n := range f.rows {
		if n.UserID == userID && n.ReadAt == nil {
			t := now
			n.ReadAt = &t
		}
	}
	return nil
}

func (f *fakeNotificationRepo) UnreadCount(_ context.Context, userID uuid.UUID) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	count := 0
	for _, n := range f.rows {
		if n.UserID == userID && n.ReadAt == nil {
			count++
		}
	}
	return count, nil
}

// ---- fake device token repo -------------------------------------------------

type fakeDeviceTokenRepo struct {
	mu     sync.Mutex
	byUser map[uuid.UUID][]*entity.DeviceToken
}

func newFakeDeviceTokenRepo() *fakeDeviceTokenRepo {
	return &fakeDeviceTokenRepo{byUser: map[uuid.UUID][]*entity.DeviceToken{}}
}

func (f *fakeDeviceTokenRepo) Upsert(_ context.Context, t *entity.DeviceToken) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	tokens := f.byUser[t.UserID]
	for _, existing := range tokens {
		if existing.FCMToken == t.FCMToken {
			existing.Platform = t.Platform
			return nil
		}
	}
	cp := *t
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}
	f.byUser[t.UserID] = append(tokens, &cp)
	return nil
}

func (f *fakeDeviceTokenRepo) Delete(_ context.Context, fcmToken string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for uid, tokens := range f.byUser {
		var keep []*entity.DeviceToken
		for _, t := range tokens {
			if t.FCMToken != fcmToken {
				keep = append(keep, t)
			}
		}
		f.byUser[uid] = keep
	}
	return nil
}

func (f *fakeDeviceTokenRepo) DeleteInvalid(_ context.Context, tokens []string) error {
	set := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		set[t] = true
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for uid, tks := range f.byUser {
		var keep []*entity.DeviceToken
		for _, t := range tks {
			if !set[t.FCMToken] {
				keep = append(keep, t)
			}
		}
		f.byUser[uid] = keep
	}
	return nil
}

func (f *fakeDeviceTokenRepo) ListByUser(_ context.Context, userID uuid.UUID) ([]*entity.DeviceToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*entity.DeviceToken, len(f.byUser[userID]))
	copy(out, f.byUser[userID])
	return out, nil
}

// ---- fake preference repo ---------------------------------------------------

type fakePreferenceRepo struct {
	mu    sync.Mutex
	prefs map[uuid.UUID]*entity.NotificationPreference
}

func newFakePreferenceRepo() *fakePreferenceRepo {
	return &fakePreferenceRepo{prefs: map[uuid.UUID]*entity.NotificationPreference{}}
}

func (f *fakePreferenceRepo) Get(_ context.Context, userID uuid.UUID) (*entity.NotificationPreference, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.prefs[userID]; ok {
		cp := *p
		return &cp, nil
	}
	// Default: all enabled.
	return &entity.NotificationPreference{
		UserID:       userID,
		PushEnabled:  true,
		EmailEnabled: true,
	}, nil
}

func (f *fakePreferenceRepo) Upsert(_ context.Context, p *entity.NotificationPreference) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *p
	f.prefs[p.UserID] = &cp
	return nil
}

// ---- fake job queue ---------------------------------------------------------

type fakeJobQueue struct {
	mu   sync.Mutex
	jobs []service.NotificationJob
}

func (q *fakeJobQueue) EnqueueTranscode(context.Context, service.TranscodeJob) error { return nil }

func (q *fakeJobQueue) EnqueueNotification(_ context.Context, job service.NotificationJob) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
	return nil
}

func (q *fakeJobQueue) count() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.jobs)
}

func (q *fakeJobQueue) last() (service.NotificationJob, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.jobs) == 0 {
		return service.NotificationJob{}, false
	}
	return q.jobs[len(q.jobs)-1], true
}

// ---- fake push / email senders (worker tests) --------------------------------

type fakePushSender struct {
	mu    sync.Mutex
	calls []pushCall
	err   error
}

type pushCall struct {
	tokens []string
	ev     service.NotifyEvent
}

func (s *fakePushSender) Send(_ context.Context, tokens []string, ev service.NotifyEvent) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, pushCall{tokens: tokens, ev: ev})
	return nil, s.err
}

func (s *fakePushSender) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

type fakeEmailSender struct {
	mu    sync.Mutex
	calls []emailCall
	err   error
}

type emailCall struct {
	to      string
	subject string
}

func (s *fakeEmailSender) Send(_ context.Context, to, subject, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, emailCall{to: to, subject: subject})
	return s.err
}

func (s *fakeEmailSender) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

// ---- compile-time interface checks ------------------------------------------

var _ repository.NotificationRepository = (*fakeNotificationRepo)(nil)
var _ repository.DeviceTokenRepository = (*fakeDeviceTokenRepo)(nil)
var _ repository.PreferenceRepository = (*fakePreferenceRepo)(nil)
var _ service.JobQueue = (*fakeJobQueue)(nil)
var _ service.PushSender = (*fakePushSender)(nil)
var _ service.EmailSender = (*fakeEmailSender)(nil)
