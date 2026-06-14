package entity

import (
	"time"

	"github.com/google/uuid"
)

// AdminAuditLog is one immutable record of a privileged admin action.
type AdminAuditLog struct {
	ID         uuid.UUID
	ActorID    uuid.UUID
	Action     string
	TargetType string
	TargetID   *uuid.UUID
	Details    map[string]any
	CreatedAt  time.Time
}
