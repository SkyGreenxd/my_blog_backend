CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    refresh_token_hash TEXT UNIQUE NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);