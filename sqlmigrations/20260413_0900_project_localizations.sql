CREATE TABLE IF NOT EXISTS project_localizations (
    project_id UUID NOT NULL,
    locale VARCHAR(8) NOT NULL,
    field_key VARCHAR(64) NOT NULL,
    value JSONB NOT NULL,
    mode VARCHAR(16) NOT NULL DEFAULT 'auto',
    source_hash VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT project_localizations_pk PRIMARY KEY (project_id, locale, field_key),
    CONSTRAINT project_localizations_project_fk FOREIGN KEY (project_id)
        REFERENCES products (id) ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT project_localizations_locale_ck CHECK (locale IN ('ca', 'en', 'de')),
    CONSTRAINT project_localizations_mode_ck CHECK (mode IN ('auto', 'manual'))
);

CREATE INDEX IF NOT EXISTS ix_project_localizations_project_locale
    ON project_localizations (project_id, locale, field_key);
