package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ObjectIndex maps stable storage keys (e.g. "lessons/<id>/notes.pdf") to the
// opaque Google Drive file ids that back them. It is the source of truth for
// locating a Drive-backed object; Drive file names are not unique, so objects
// are never resolved by name.
type ObjectIndex interface {
	// Put records (or replaces) the key -> driveFileID mapping.
	Put(ctx context.Context, key, driveFileID, contentType string, size int64) error
	// Lookup returns the Drive file id for key; ok is false when key is unknown.
	Lookup(ctx context.Context, key string) (driveFileID string, ok bool, err error)
	// Delete removes the mapping for key (idempotent).
	Delete(ctx context.Context, key string) error
}

// PgObjectIndex is a Postgres-backed ObjectIndex over storage.objects. It uses
// raw pgx (not sqlc) so the Drive backend is self-contained and needs no code
// generation.
type PgObjectIndex struct {
	pool *pgxpool.Pool
}

// NewPgObjectIndex builds a PgObjectIndex over the shared pgx pool.
func NewPgObjectIndex(pool *pgxpool.Pool) *PgObjectIndex { return &PgObjectIndex{pool: pool} }

var _ ObjectIndex = (*PgObjectIndex)(nil)

func (p *PgObjectIndex) Put(ctx context.Context, key, driveFileID, contentType string, size int64) error {
	const q = `
INSERT INTO storage.objects (object_key, drive_file_id, content_type, size_bytes, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
ON CONFLICT (object_key) DO UPDATE
   SET drive_file_id = EXCLUDED.drive_file_id,
       content_type  = EXCLUDED.content_type,
       size_bytes    = EXCLUDED.size_bytes,
       updated_at    = now()`
	if _, err := p.pool.Exec(ctx, q, key, driveFileID, contentType, size); err != nil {
		return fmt.Errorf("storage: index put %q: %w", key, err)
	}
	return nil
}

func (p *PgObjectIndex) Lookup(ctx context.Context, key string) (string, bool, error) {
	const q = `SELECT drive_file_id FROM storage.objects WHERE object_key = $1`
	var id string
	err := p.pool.QueryRow(ctx, q, key).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("storage: index lookup %q: %w", key, err)
	}
	return id, true, nil
}

func (p *PgObjectIndex) Delete(ctx context.Context, key string) error {
	const q = `DELETE FROM storage.objects WHERE object_key = $1`
	if _, err := p.pool.Exec(ctx, q, key); err != nil {
		return fmt.Errorf("storage: index delete %q: %w", key, err)
	}
	return nil
}
