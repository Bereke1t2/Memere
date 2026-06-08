// Package postgres holds the concrete, sqlc-backed implementations of the domain
// repository interfaces. The mappers here are the *only* place sqlc's pgtype
// values cross into domain entities — usecases never see a pgtype (dependency
// rule): the conversions are centralized so every repository maps the database's
// nullable representation to the domain's pointer/value representation the same
// way.
package postgres

import (
	"encoding/json"
	"strconv"
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

// toPgUUIDPtr wraps a *uuid.UUID; a nil pointer maps to SQL NULL. Used for
// nullable FKs (quizzes.lesson_id, exams.course_id).
func toPgUUIDPtr(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// fromPgUUIDPtr unwraps a nullable pgtype.UUID to a *uuid.UUID (nil when NULL).
func fromPgUUIDPtr(v pgtype.UUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	id := uuid.UUID(v.Bytes)
	return &id
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

// toPgNumeric converts a float64 to a pgtype.Numeric for DECIMAL columns
// (price, rating_avg). Phase 1 does no money arithmetic, so float round-tripping
// is acceptable; formatting through a fixed 2-decimal string avoids binary-float
// drift for the values these columns hold.
func toPgNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Scan never errors for a well-formed decimal string.
	_ = n.Scan(strconv.FormatFloat(f, 'f', 2, 64))
	return n
}

// fromPgNumeric converts a pgtype.Numeric to float64 for display/storage
// round-tripping. A NULL or unset value yields 0.
func fromPgNumeric(v pgtype.Numeric) float64 {
	if !v.Valid {
		return 0
	}
	f, err := v.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

// fromPgNumericPtr unwraps a nullable numeric to a *float64 (nil when NULL).
// Used for attempt score/percentage, which stay NULL until grading completes.
func fromPgNumericPtr(v pgtype.Numeric) *float64 {
	if !v.Valid {
		return nil
	}
	f := fromPgNumeric(v)
	return &f
}

// toPgNumericPtr converts a *float64 to a pgtype.Numeric; a nil pointer maps to
// SQL NULL.
func toPgNumericPtr(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	return toPgNumeric(*f)
}

// toPgInt4Ptr converts a *int to a nullable int4; nil maps to SQL NULL.
func toPgInt4Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}
	v := int32(*i)
	return &v
}

// fromPgInt4Ptr unwraps a nullable int4 to a *int (nil when NULL).
func fromPgInt4Ptr(v *int32) *int {
	if v == nil {
		return nil
	}
	i := int(*v)
	return &i
}

// toJSONB marshals a metadata map to the []byte sqlc expects for a JSONB column.
// A nil map maps to SQL NULL (nil []byte).
func toJSONB(m map[string]any) []byte {
	if m == nil {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

// fromJSONB unmarshals a JSONB column into a metadata map. NULL/empty yields nil.
func fromJSONB(b []byte) map[string]any {
	if len(b) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}
