// Package migrations embeds the SQL migration files so the migrate runner
// (cmd/migrate) can apply them without shipping the .sql files separately.
package migrations

import "embed"

// FS holds every golang-migrate .up.sql / .down.sql file in this directory.
//
//go:embed *.sql
var FS embed.FS
