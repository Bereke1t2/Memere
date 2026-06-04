// Package pagination provides cursor-based (keyset) pagination helpers used by
// repositories and the delivery layer (spec §5.4: cursor pagination with
// limit/after).
package pagination

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultLimit is applied when a request omits or zeroes the limit.
	DefaultLimit = 20
	// MaxLimit caps the page size to protect the database.
	MaxLimit = 100
)

// Cursor is an opaque keyset position. Listings order by (created_at, id)
// descending, so a cursor carries the last seen row's sort keys.
type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

// Encode serializes the cursor to a URL-safe opaque string for the `after`
// query parameter.
func (c Cursor) Encode() string {
	b, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(b)
}

// Decode parses an opaque cursor string. An empty string yields (nil, nil),
// meaning "from the beginning".
func Decode(s string) (*Cursor, error) {
	if s == "" {
		return nil, nil
	}
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// NormalizeLimit clamps a requested limit into [1, MaxLimit], defaulting when
// non-positive.
func NormalizeLimit(limit int) int {
	switch {
	case limit <= 0:
		return DefaultLimit
	case limit > MaxLimit:
		return MaxLimit
	default:
		return limit
	}
}
