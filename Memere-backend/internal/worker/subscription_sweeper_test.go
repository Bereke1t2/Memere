package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
)

// fakeExpirer returns a scripted batch of expired subscriptions and records calls.
type fakeExpirer struct {
	mu      sync.Mutex
	calls   int
	expired []*entity.Subscription
	err     error
}

func (f *fakeExpirer) SweepExpired(context.Context, int) ([]*entity.Subscription, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return f.expired, f.err
}

func (f *fakeExpirer) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// recordingNotifier records which subscriptions it was told expired.
type recordingNotifier struct {
	mu       sync.Mutex
	notified []uuid.UUID
}

func (n *recordingNotifier) SubscriptionExpired(_ context.Context, sub *entity.Subscription) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notified = append(n.notified, sub.ID)
}

func (n *recordingNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.notified)
}

func TestSubscriptionSweeper_SweepsOnStartAndNotifies(t *testing.T) {
	exp := &fakeExpirer{expired: []*entity.Subscription{
		{ID: uuid.New(), Status: entity.SubExpired},
		{ID: uuid.New(), Status: entity.SubExpired},
	}}
	notify := &recordingNotifier{}
	sw := NewSubscriptionSweeper(time.Hour, exp, notify) // long interval: only the initial sweep fires

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		sw.Run(ctx)
		close(done)
	}()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if exp.callCount() >= 1 && notify.count() >= 2 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if exp.callCount() != 1 {
		t.Fatalf("want one initial sweep, got %d", exp.callCount())
	}
	if notify.count() != 2 {
		t.Fatalf("want 2 expiry notifications, got %d", notify.count())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sweeper did not stop on context cancel")
	}
}

func TestSubscriptionSweeper_ContinuesPastError(t *testing.T) {
	exp := &fakeExpirer{err: errors.New("boom")}
	sw := NewSubscriptionSweeper(time.Hour, exp, nil) // nil notify -> no-op default

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sw.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if exp.callCount() >= 1 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	// A sweep error must not crash the worker; it keeps running for the next tick.
	if exp.callCount() < 1 {
		t.Errorf("sweep should have been attempted; calls=%d", exp.callCount())
	}
}
