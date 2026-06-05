// Package password provides bcrypt-based password hashing and verification.
// Plaintext passwords are never logged, returned, or stored — only their bcrypt
// hashes (spec §7.3, Non-Negotiable #8).
package password

import "golang.org/x/crypto/bcrypt"

// Cost is the bcrypt work factor. The spec requires ≥ 12; raising it later only
// affects newly hashed passwords (existing hashes carry their own cost).
const Cost = 12

// Hash returns the bcrypt hash of a plaintext password. The cost is embedded in
// the returned string, so Compare needs no external configuration.
func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), Cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare reports whether plain matches the bcrypt hash. It returns nil on a
// match and a non-nil error otherwise; the comparison is constant-time with
// respect to the hash, so it does not leak timing information about the value.
func Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
