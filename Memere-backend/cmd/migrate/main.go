// Command migrate applies (or reverts) the embedded SQL migrations against the
// database configured via the environment. Usage:
//
//	go run ./cmd/migrate -direction up
//	go run ./cmd/migrate -direction down
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
	dbURL := "pgx5://" + strings.TrimPrefix(cfg.DB.DSN(), "postgres://")

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		log.Fatalf("init migrator: %v", err)
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("migrator close: src=%v db=%v", srcErr, dbErr)
		}
	}()

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
