ALTER TABLE project_media
    ADD COLUMN IF NOT EXISTS thumbnail_url TEXT,
    ADD COLUMN IF NOT EXISTS medium_url TEXT,
    ADD COLUMN IF NOT EXISTS full_url TEXT,
    ADD COLUMN IF NOT EXISTS alt_text TEXT,
    ADD COLUMN IF NOT EXISTS featured BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE project_media
SET thumbnail_url = COALESCE(NULLIF(thumbnail_url, ''), url),
    medium_url = COALESCE(NULLIF(medium_url, ''), url),
    full_url = COALESCE(NULLIF(full_url, ''), url),
    alt_text = COALESCE(alt_text, ''),
    featured = COALESCE(featured, FALSE)
WHERE thumbnail_url IS NULL
   OR medium_url IS NULL
   OR full_url IS NULL
   OR alt_text IS NULL;

CREATE INDEX IF NOT EXISTS ix_project_media_project_featured_sort
    ON project_media (project_id, featured DESC, sort_order ASC);
