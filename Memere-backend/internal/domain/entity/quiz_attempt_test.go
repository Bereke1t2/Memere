package entity

import "testing"

func TestAttemptStatusCanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from AttemptStatus
		to   AttemptStatus
		want bool
	}{
		{"in_progress to submitted", AttemptInProgress, AttemptSubmitted, true},
		{"in_progress to expired", AttemptInProgress, AttemptExpired, true},
		{"in_progress to graded", AttemptInProgress, AttemptGraded, true},
		{"in_progress to in_progress", AttemptInProgress, AttemptInProgress, false},
		{"submitted to graded", AttemptSubmitted, AttemptGraded, true},
		{"submitted to expired", AttemptSubmitted, AttemptExpired, false},
		{"submitted to in_progress", AttemptSubmitted, AttemptInProgress, false},
		{"expired to graded", AttemptExpired, AttemptGraded, true},
		{"expired to submitted", AttemptExpired, AttemptSubmitted, false},
		{"expired to in_progress", AttemptExpired, AttemptInProgress, false},
		{"graded is terminal to in_progress", AttemptGraded, AttemptInProgress, false},
		{"graded is terminal to submitted", AttemptGraded, AttemptSubmitted, false},
		{"graded is terminal to graded", AttemptGraded, AttemptGraded, false},
		{"unknown source rejects", AttemptStatus("bogus"), AttemptGraded, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("%s.CanTransitionTo(%s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
