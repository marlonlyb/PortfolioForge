-- PortfolioForge: Search Backfill
-- Populates project_search_documents for existing active projects
-- and migrates legacy product columns to new search-aware columns.
-- Date: 2026-04-12
-- Depends on: 20260412_0200_search_schema.sql

-- ===========================================================================
-- 1. Backfill search_document for existing projects
-- ===========================================================================
-- Uses ON CONFLICT to support re-running: inserts new rows, updates existing ones.

INSERT INTO project_search_documents (project_id, search_document, search_composed_at)
SELECT p.id, compose_project_search_doc(p.id), NOW()
FROM products p
WHERE p.active = TRUE
ON CONFLICT (project_id) DO UPDATE
SET search_document = compose_project_search_doc(EXCLUDED.project_id),
    search_composed_at = NOW();

-- ===========================================================================
-- 1b. Backfill search_trgm for existing projects
-- ===========================================================================

UPDATE project_search_documents psd
SET search_trgm = compose_project_search_trgm(psd.project_id)
WHERE psd.search_trgm IS NULL;

-- ===========================================================================
-- 2. Migrate brand → client_name where client_name is NULL
-- ===========================================================================
-- In the portfolio context, "brand" represents the client/company name.
-- This one-time migration copies brand values to the new client_name column.

UPDATE products SET client_name = brand
WHERE client_name IS NULL
  AND brand IS NOT NULL
  AND brand != '';

-- ===========================================================================
-- 3. Set status based on active flag
-- ===========================================================================
-- Sync the new status column with the existing active flag.
-- active=TRUE → 'published', active=FALSE → 'draft'

UPDATE products SET status = 'published'
WHERE active = TRUE AND status = 'published';

UPDATE products SET status = 'draft'
WHERE active = FALSE AND status = 'published';
