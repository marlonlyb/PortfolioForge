ALTER TABLE users
    ADD COLUMN IF NOT EXISTS auth_provider VARCHAR(24),
    ADD COLUMN IF NOT EXISTS provider_subject VARCHAR(255),
    ADD COLUMN IF NOT EXISTS email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS full_name VARCHAR(160),
    ADD COLUMN IF NOT EXISTS company VARCHAR(160),
    ADD COLUMN IF NOT EXISTS last_login_at INTEGER;

UPDATE users
SET auth_provider = 'local'
WHERE auth_provider IS NULL;

ALTER TABLE users
    ALTER COLUMN auth_provider SET NOT NULL,
    ALTER COLUMN auth_provider SET DEFAULT 'local',
    ALTER COLUMN password DROP NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_provider_subject
    ON users (auth_provider, provider_subject)
    WHERE provider_subject IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email_normalized
    ON users (LOWER(email));

CREATE INDEX IF NOT EXISTS ix_users_auth_provider
    ON users (auth_provider);
