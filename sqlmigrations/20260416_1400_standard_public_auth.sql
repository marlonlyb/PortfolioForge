ALTER TABLE users
    ADD COLUMN IF NOT EXISTS local_auth_state VARCHAR(32) NOT NULL DEFAULT 'ready';

-- Only legacy passwordless public accounts should require a password setup.
-- Those accounts were auto-created during the OTP-only public auth window that
-- started with the authenticated public-auth rollout on 2026-04-15 10:00 UTC.
UPDATE users
SET local_auth_state = CASE
    WHEN COALESCE(NULLIF(TRIM(auth_provider), ''), 'local') = 'local'
        AND is_admin = FALSE
        AND created_at >= 1776247200
        AND email_verified = FALSE
        AND COALESCE(NULLIF(TRIM(full_name), ''), '') = ''
        AND COALESCE(NULLIF(TRIM(company), ''), '') = ''
        AND COALESCE(last_login_at, 0) = 0
        AND COALESCE(details, '{}'::jsonb) = '{}'::jsonb
    THEN 'password_setup_required'
    ELSE 'ready'
END;
