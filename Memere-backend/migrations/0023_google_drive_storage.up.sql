-- Google Drive storage backend (STORAGE_PROVIDER=gdrive).
-- One admin-owned Drive account backs all object storage. The backend maps its
-- stable object keys (e.g. "lessons/<id>/notes.pdf") to opaque Drive file ids
-- here, and stores the single OAuth refresh token (encrypted) that lets it mint
-- Drive access tokens. Students and teachers never authenticate with Google.

CREATE SCHEMA IF NOT EXISTS storage;

-- 1. Object index: maps a stable storage key to the opaque Google Drive file id.
--    This is the source of truth for "where does key X live in Drive". Drive
--    file names are not unique, so objects are never resolved by name.
CREATE TABLE IF NOT EXISTS storage.objects (
    object_key    TEXT PRIMARY KEY,          -- e.g. "lessons/<id>/notes.pdf"
    drive_file_id TEXT NOT NULL,             -- opaque Google Drive file id
    content_type  TEXT NOT NULL DEFAULT '',
    size_bytes    BIGINT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. Single-row credential store for the admin Drive account's OAuth refresh
--    token. The token is stored ENCRYPTED (AES-256-GCM; key derived from
--    JWT_SECRET) so a DB dump alone cannot be replayed against Google. The
--    single-row invariant is enforced by a fixed primary key (id = 1).
CREATE TABLE IF NOT EXISTS storage.google_credentials (
    id             SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    refresh_token  BYTEA NOT NULL,           -- AES-256-GCM ciphertext (nonce||ciphertext)
    google_email   TEXT,                     -- connected account email, for the admin UI only
    connected_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
