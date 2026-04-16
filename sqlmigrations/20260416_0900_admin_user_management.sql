ALTER TABLE users
    ADD COLUMN IF NOT EXISTS deleted_at INTEGER;

CREATE INDEX IF NOT EXISTS ix_users_active_lookup
    ON users (deleted_at, created_at DESC);
