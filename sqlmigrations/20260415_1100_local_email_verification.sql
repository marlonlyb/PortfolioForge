CREATE TABLE IF NOT EXISTS email_verification_challenges (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL,
    attempt_count INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5,
    resend_available_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    consumed_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM now())::int,
    updated_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM now())::int,
    CONSTRAINT chk_email_verification_attempt_count_non_negative CHECK (attempt_count >= 0),
    CONSTRAINT chk_email_verification_max_attempts_positive CHECK (max_attempts > 0)
);

CREATE INDEX IF NOT EXISTS ix_email_verification_challenges_user_created
    ON email_verification_challenges (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS ix_email_verification_challenges_user_active
    ON email_verification_challenges (user_id, consumed_at, expires_at);

CREATE INDEX IF NOT EXISTS ix_email_verification_challenges_expires_at
    ON email_verification_challenges (expires_at);
