// Package postgres holds the concrete, sqlc-backed implementations of the domain
// repository interfaces. The mappers here are the *only* place sqlc's pgtype
// values cross into domain entities — usecases never see a pgtype (dependency
// rule): the conversions are centralized so every repository maps the database's
// nullable representation to the domain's pointer/value representation the same
// way.
package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// toPgUUID wraps a domain uuid.UUID as a valid pgtype.UUID.
func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// fromPgUUID unwraps a pgtype.UUID. An invalid (NULL) value yields the zero
// uuid.UUID, which the schema never produces for a NOT NULL primary key.
func fromPgUUID(v pgtype.UUID) uuid.UUID {
	if !v.Valid {
		return uuid.Nil
	}
	return v.Bytes
}

// toPgTimestamptz wraps a non-nil time as a valid timestamptz; a nil pointer
// maps to SQL NULL.
func toPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// pgTimestamptzValue wraps a plain (non-nullable) time as a valid timestamptz.
func pgTimestamptzValue(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// fromPgTimestamptz unwraps a nullable timestamptz to a *time.Time (nil when
// NULL).
func fromPgTimestamptz(v pgtype.Timestamptz) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

// fromPgTimestamptzValue unwraps a timestamptz that the schema guarantees NOT
// NULL (e.g. created_at), returning the zero time if it is unexpectedly NULL.
func fromPgTimestamptzValue(v pgtype.Timestamptz) time.Time {
	if !v.Valid {
		return time.Time{}
	}
	return v.Time
}
