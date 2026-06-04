-- Per-domain schemas (spec §4.1).
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS courses;
CREATE SCHEMA IF NOT EXISTS payments;
CREATE SCHEMA IF NOT EXISTS progress;

-- gen_random_uuid() is built into PostgreSQL 13+, so no pgcrypto extension is
-- required for UUID primary keys.

-- Shared trigger function: auto-touch updated_at on every UPDATE.
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
