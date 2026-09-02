// Command migrate applies (or reverts) the embedded SQL migrations against the
// database configured via the environment. Usage:
//
//	go run ./cmd/migrate -direction up
//	go run ./cmd/migrate -direction down
//	go run ./cmd/migrate -force 18    // recovery: clear a dirty state, set version
//
// It uses the golang-migrate library with an embedded source, so no external
// `migrate` binary is required.
package main

import (
	"errors"
	"flag"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/Bereke1t2/Memere/memere-backend/migrations"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	force := flag.Int("force", -1, "recovery: set the recorded version and clear the dirty flag without running any migration; -1 disables")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("load embedded migrations: %v", err)
	}

	// The golang-migrate pgx/v5 driver registers the "pgx5" URL scheme.
	rawDSN := cfg.DB.DSN()
	rawDSN = strings.TrimPrefix(rawDSN, "postgresql://")
	rawDSN = strings.TrimPrefix(rawDSN, "postgres://")
	dbURL := "pgx5://" + rawDSN

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		log.Fatalf("init migrator: %v", err)
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("migrator close: src=%v db=%v", srcErr, dbErr)
		}
	}()

	// Recovery path: a failed migration leaves golang-migrate in a "dirty" state
	// that blocks every subsequent up/down. -force clears the dirty flag and sets
	// the recorded version WITHOUT running any migration, so the next `up` resumes
	// from there (e.g. `-force 18` then `-direction up` re-runs migration 0019).
	if *force >= 0 {
		if err := m.Force(*force); err != nil {
			log.Fatalf("migrate force %d: %v", *force, err)
		}
		log.Printf("migrate force: version set to %d, dirty flag cleared", *force)
		return
	}

	switch *direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("unknown direction %q (want up or down)", *direction)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate %s: %v", *direction, err)
	}

	log.Printf("migrate %s: complete", *direction)
}
