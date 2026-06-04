-- auth.users (spec §4.2.1)
CREATE TABLE auth.users (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                     VARCHAR(255) UNIQUE NOT NULL,
    phone                     VARCHAR(20) UNIQUE,
    password_hash             VARCHAR(255) NOT NULL,
    role                      VARCHAR(20) NOT NULL CHECK (role IN ('student', 'teacher', 'admin')),
    first_name                VARCHAR(100) NOT NULL,
    last_name                 VARCHAR(100) NOT NULL,
    avatar_url                TEXT,
    is_active                 BOOLEAN NOT NULL DEFAULT true,
    is_email_verified         BOOLEAN NOT NULL DEFAULT false,
    email_verification_token  VARCHAR(255),
    password_reset_token      VARCHAR(255),
    password_reset_expires_at TIMESTAMPTZ,
    last_login_at             TIMESTAMPTZ,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at                TIMESTAMPTZ
);

CREATE INDEX idx_users_role ON auth.users (role);
CREATE INDEX idx_users_deleted_at ON auth.users (deleted_at);
CREATE INDEX idx_users_created_at ON auth.users (created_at);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON auth.users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- auth.refresh_tokens (spec §7.1) — only the hash of the token is stored.
CREATE TABLE auth.refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) UNIQUE NOT NULL,
    device_info TEXT,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_user_id ON auth.refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON auth.refresh_tokens (expires_at);
