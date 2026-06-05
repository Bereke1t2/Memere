package course

import (
	"crypto/rand"
	"encoding/hex"
)

// shortSuffix returns a 4-hex-char random string used to disambiguate a slug on
// collision (e.g. "calculus" → "calculus-1a2b").
func shortSuffix() string {
	b := make([]byte, 2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
