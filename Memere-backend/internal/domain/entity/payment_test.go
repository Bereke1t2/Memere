package entity

import "testing"

// TestPaymentStatusCanTransitionTo pins the payment state machine: pending
// settles to completed or fails; completed only refunds; failed/refunded are
// terminal. The same rule is enforced by the guarded UPDATE in the repo.
func TestPaymentStatusCanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from PaymentStatus
		to   PaymentStatus
		want bool
	}{
		{"pending to completed", PayPending, PayCompleted, true},
		{"pending to failed", PayPending, PayFailed, true},
		{"pending to refunded", PayPending, PayRefunded, false},
		{"pending to pending", PayPending, PayPending, false},
		{"completed to refunded", PayCompleted, PayRefunded, true},
		{"completed to failed", PayCompleted, PayFailed, false},
		{"completed to completed", PayCompleted, PayCompleted, false},
		{"failed is terminal", PayFailed, PayCompleted, false},
		{"refunded is terminal", PayRefunded, PayCompleted, false},
		{"unknown source rejects", PaymentStatus("bogus"), PayCompleted, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("%s.CanTransitionTo(%s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
