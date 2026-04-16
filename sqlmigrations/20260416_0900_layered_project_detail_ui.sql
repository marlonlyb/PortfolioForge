ALTER TABLE project_profiles
    ADD COLUMN IF NOT EXISTS delivery_scope TEXT,
    ADD COLUMN IF NOT EXISTS responsibility_scope TEXT;
