package sqlcgen

// toJSONString formats a []byte JSON parameter as a string for queries running with pgx simple protocol.
// In simple query protocol (such as with Neon / PgBouncer), []byte is serialized as hex bytea by pgx,
// which causes PostgreSQL $N::jsonb cast to fail with SQLSTATE 22P02 (invalid input syntax for type json).
// Passing as a Go string ensures PostgreSQL parses it directly as JSON text.
func toJSONString(b []byte) string {
	if len(b) == 0 {
		return "{}"
	}
	return string(b)
}
