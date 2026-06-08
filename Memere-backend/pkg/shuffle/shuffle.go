// Package shuffle provides a deterministic, per-attempt ordering of quiz/exam
// questions and their answer options (spec §9.1: "Shuffle on attempt creation").
//
// The order is seeded from the attempt ID, so it is:
//   - stable for the whole attempt (re-deriving from the same ID yields the same
//     order), which lets the server reconstruct the order if the Redis snapshot
//     is lost, falling back to the durable copy logic without storing anything;
//   - different per attempt (different attempt IDs → different seeds → different
//     orders), so two attempts at the same quiz see different orders.
//
// It uses math/rand seeded deterministically — NOT for any security purpose
// (the answer key never leaves the server regardless of order), purely for
// reproducible presentation.
package shuffle

import (
	"encoding/binary"
	"math/rand"

	"github.com/google/uuid"
)

// seedFromAttempt derives a stable int64 seed from an attempt UUID.
func seedFromAttempt(attemptID uuid.UUID) int64 {
	b := attemptID[:]
	return int64(binary.LittleEndian.Uint64(b[:8]))
}

// newRand builds a deterministic RNG for an attempt.
func newRand(attemptID uuid.UUID) *rand.Rand {
	return rand.New(rand.NewSource(seedFromAttempt(attemptID)))
}

// Order returns a permutation of [0, n) deterministic for the given attempt and
// salt. The salt distinguishes independent shuffles within one attempt (e.g.
// the question order vs. each question's answer order) so they don't all share
// the same permutation.
func Order(attemptID uuid.UUID, salt uint64, n int) []int {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	if n < 2 {
		return perm
	}
	r := newRand(saltedAttempt(attemptID, salt))
	r.Shuffle(n, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
	return perm
}

// Identity returns the unshuffled order [0, n) — used when randomization is off.
func Identity(n int) []int {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	return perm
}

// saltedAttempt folds a salt into the attempt UUID so distinct salts produce
// independent permutations while remaining deterministic.
func saltedAttempt(attemptID uuid.UUID, salt uint64) uuid.UUID {
	var out uuid.UUID
	copy(out[:], attemptID[:])
	var s [8]byte
	binary.LittleEndian.PutUint64(s[:], salt)
	for i := 0; i < 8; i++ {
		out[i] ^= s[i]
	}
	return out
}
