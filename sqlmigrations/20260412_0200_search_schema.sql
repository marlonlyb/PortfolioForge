-- PortfolioForge: Search Schema
-- Adds search-related columns to products, creates project_search_documents table,
-- indexes (GIN + HNSW), and the compose_project_search_doc() function.
-- Date: 2026-04-12
-- Depends on: 20260412_0100_search_extensions.sql, 20260411_0900_portfolioforge_extension.sql

-- ===========================================================================
-- 1. Add new columns to products table
-- ===========================================================================

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS client_name VARCHAR(160),
    ADD COLUMN IF NOT EXISTS status VARCHAR(24) NOT NULL DEFAULT 'published',
    ADD COLUMN IF NOT EXISTS featured BOOLEAN NOT NULL DEFAULT FALSE;

-- ===========================================================================
-- 2. Create project_search_documents table
-- ===========================================================================

CREATE TABLE IF NOT EXISTS project_search_documents (
    project_id UUID NOT NULL PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    search_document TSVECTOR,
    search_embedding VECTOR(1536),
    search_trgm TEXT,
    search_content_hash VARCHAR(64),
    search_composed_at TIMESTAMPTZ
);

-- NOTE: search_trgm is a plain TEXT column, NOT a GENERATED column.
-- PostgreSQL does not allow subqueries in GENERATED column expressions.
-- It is maintained by the compose_project_search_trgm() function and
-- should be updated whenever project data changes (via trigger or app logic).

-- ===========================================================================
-- 3. Create compose_project_search_doc() function
-- ===========================================================================
-- Builds a weighted TSVECTOR from project data across products, project_profiles,
-- and technologies. Weight priorities:
--   A = name, brand, solution_summary  (highest relevance)
--   B = architecture, description, technology names
--   C = business_goal, ai_usage        (lowest relevance)

CREATE OR REPLACE FUNCTION compose_project_search_doc(p_id UUID)
RETURNS TSVECTOR AS $$
SELECT
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.solution_summary, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.name, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.brand, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.architecture, ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.description, ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(
        string_agg(t.name, ' '), ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.business_goal, ''))), 'C') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.ai_usage, ''))), 'C')
FROM products p
LEFT JOIN project_profiles pp ON pp.project_id = p.id
LEFT JOIN project_technologies pt ON pt.project_id = p.id
LEFT JOIN technologies t ON t.id = pt.technology_id
WHERE p.id = p_id
GROUP BY p.id, p.name, p.brand, p.description, pp.solution_summary,
         pp.architecture, pp.business_goal, pp.ai_usage;
$$ LANGUAGE SQL STABLE;

-- ===========================================================================
-- 3b. Create compose_project_search_trgm() function
-- ===========================================================================
-- Builds the concatenated text used for pg_trgm similarity matching.
-- Includes project name, brand, and technology names.

CREATE OR REPLACE FUNCTION compose_project_search_trgm(p_id UUID)
RETURNS TEXT AS $$
SELECT
    COALESCE(p.name, '') || ' ' ||
    COALESCE(p.brand, '') || ' ' ||
    COALESCE(string_agg(t.name, ' '), '')
FROM products p
LEFT JOIN project_technologies pt ON pt.project_id = p.id
LEFT JOIN technologies t ON t.id = pt.technology_id
WHERE p.id = p_id
GROUP BY p.id, p.name, p.brand;
$$ LANGUAGE SQL STABLE;

-- ===========================================================================
-- 4. Create indexes
-- ===========================================================================

-- GIN index on tsvector for full-text search
CREATE INDEX IF NOT EXISTS ix_project_search_document ON project_search_documents
    USING GIN (search_document);

-- GIN index on trigram for fuzzy matching
CREATE INDEX IF NOT EXISTS ix_project_search_trgm ON project_search_documents
    USING GIN (search_trgm gin_trgm_ops);

-- HNSW index on vector for semantic similarity (cosine distance)
-- NOTE: Requires the pgvector extension. This will fail if pgvector is not installed.
-- Install pgvector first: sudo apt-get install -y postgresql-16-pgvector
-- Then run: CREATE EXTENSION IF NOT EXISTS vector;
-- Then run the CREATE INDEX below manually.
-- CREATE INDEX IF NOT EXISTS ix_project_search_embedding ON project_search_documents
--     USING hnsw (search_embedding vector_cosine_ops)
--     WITH (m = 16, ef_construction = 64);

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'vector') THEN
        CREATE INDEX IF NOT EXISTS ix_project_search_embedding ON project_search_documents
            USING hnsw (search_embedding vector_cosine_ops)
            WITH (m = 16, ef_construction = 64);
    END IF;
END $$;

-- ===========================================================================
-- 5. pg_trgm similarity threshold
-- ===========================================================================
-- NOTE: pg_trgm.similarity_threshold is a session-level setting.
-- It should be set at application level when opening connections.
-- Recommended value: 0.3
-- To set globally (requires superuser):
--   ALTER SYSTEM SET pg_trgm.similarity_threshold = 0.3;
--   SELECT pg_reload_conf();
