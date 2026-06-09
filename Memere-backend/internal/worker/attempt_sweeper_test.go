package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// fakeSweeper records how many times it was swept and returns a scripted count.
type fakeSweeper struct {
	mu     sync.Mutex
	calls  int
	graded int
	err    error
}

func (f *fakeSweeper) SweepExpired(_ context.Context, _ time.Time, _ int) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return f.graded, f.err
}

func (f *fakeSweeper) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func TestAttemptSweeper_SweepsOnStartAndStopsOnCancel(t *testing.T) {
	q := &fakeSweeper{graded: 2}
	e := &fakeSweeper{graded: 0}
	sw := NewAttemptSweeper(time.Hour, q, e) // long interval: only the initial sweep fires

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		sw.Run(ctx)
		close(done)
	}()

	// The immediate on-start sweep should hit both engines exactly once.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if q.callCount() >= 1 && e.callCount() >= 1 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if q.callCount() != 1 || e.callCount() != 1 {
		t.Fatalf("want one initial sweep per engine, got q=%d e=%d", q.callCount(), e.callCount())
	}

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("sweeper did not stop on context cancel")
	}
}

func TestAttemptSweeper_ContinuesPastEngineError(t *testing.T) {
	failing := &fakeSweeper{err: errors.New("boom")}
	ok := &fakeSweeper{graded: 1}
	sw := NewAttemptSweeper(time.Hour, failing, ok)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sw.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if ok.callCount() >= 1 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	// A failing first engine must not prevent the second from being swept.
	if ok.callCount() < 1 {
		t.Errorf("second engine should still be swept after first errors; calls=%d", ok.callCount())
	}
}
