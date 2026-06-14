package entity

import (
	"time"

	"github.com/google/uuid"
)

// Certificate is issued once per student+course upon 100% completion.
// Serial is the non-guessable public verification code.
// FileKey points to the rendered PDF in object storage (nil until rendered).
type Certificate struct {
	ID        uuid.UUID
	StudentID uuid.UUID
	CourseID  uuid.UUID
	Serial    string
	FileKey   *string
	IssuedAt  time.Time
}
