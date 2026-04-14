CREATE TABLE IF NOT EXISTS site_settings (
    singleton_key BOOLEAN NOT NULL DEFAULT TRUE,
    public_hero_logo_url TEXT,
    public_hero_logo_alt TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT site_settings_singleton_pk PRIMARY KEY (singleton_key),
    CONSTRAINT site_settings_singleton_ck CHECK (singleton_key = TRUE)
);
