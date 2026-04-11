-- PortfolioForge: Search Extensions
-- Required PostgreSQL extensions for hybrid retrieval search (FTS + trigram + vector)
-- Date: 2026-04-12

-- Full-text search accent stripping
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Trigram similarity for fuzzy matching
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- pgvector for semantic (embedding) similarity search
-- NOTE: Requires pgvector to be installed on the PostgreSQL server.
-- If not available, this will fail and semantic search will be unavailable.
-- Install: sudo apt-get install -y postgresql-16-pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- Custom text search configuration that applies unaccent before simple tokenization.
-- This ensures searches like "SCADA" match "scada" and accent-insensitive matching.
-- NOTE: "IF NOT EXISTS" for CREATE TEXT SEARCH CONFIGURATION requires PostgreSQL 17+.
-- For PG 16 compatibility, use a DO block with exception handling.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'simple_unaccent') THEN
        CREATE TEXT SEARCH CONFIGURATION simple_unaccent (COPY = simple);
    END IF;
END $$;

ALTER TEXT SEARCH CONFIGURATION simple_unaccent
    ALTER MAPPING FOR asciiword, word WITH unaccent, simple;
