CREATE OR REPLACE FUNCTION pf_normalize_profile_whitespace(input TEXT)
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT NULLIF(regexp_replace(trim(COALESCE(input, '')), '\s+', ' ', 'g'), '')
$$;

CREATE OR REPLACE FUNCTION pf_normalize_profile_key(input TEXT)
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT regexp_replace(lower(COALESCE(input, '')), '[\s-]+', '_', 'g')
$$;

CREATE OR REPLACE FUNCTION pf_profile_scalar_text(item JSONB)
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT CASE jsonb_typeof(item)
        WHEN 'string' THEN pf_normalize_profile_whitespace(trim(BOTH '"' FROM item::text))
        WHEN 'number' THEN pf_normalize_profile_whitespace(item::text)
        WHEN 'boolean' THEN pf_normalize_profile_whitespace(item::text)
        ELSE NULL
    END
$$;

CREATE OR REPLACE FUNCTION pf_profile_field(item JSONB, field_name TEXT)
RETURNS TEXT
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
    normalized_field TEXT := pf_normalize_profile_key(field_name);
    raw TEXT;
    segment TEXT;
    segment_key TEXT;
    segment_value TEXT;
    entry RECORD;
BEGIN
    IF item IS NULL THEN
        RETURN NULL;
    END IF;

    IF jsonb_typeof(item) = 'object' THEN
        FOR entry IN SELECT key, value FROM jsonb_each_text(item) LOOP
            IF pf_normalize_profile_key(entry.key) = normalized_field THEN
                RETURN pf_normalize_profile_whitespace(entry.value);
            END IF;
        END LOOP;

        RETURN NULL;
    END IF;

    raw := pf_profile_scalar_text(item);
    IF raw IS NULL THEN
        RETURN NULL;
    END IF;

    FOR segment IN SELECT regexp_split_to_table(raw, '\|') LOOP
        IF position(':' IN segment) = 0 THEN
            CONTINUE;
        END IF;

        segment_key := split_part(segment, ':', 1);
        segment_value := substr(segment, position(':' IN segment) + 1);

        IF pf_normalize_profile_key(segment_key) = normalized_field THEN
            RETURN pf_normalize_profile_whitespace(segment_value);
        END IF;
    END LOOP;

    RETURN NULL;
END;
$$;

CREATE OR REPLACE FUNCTION pf_profile_sentence(input TEXT)
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT CASE
        WHEN pf_normalize_profile_whitespace(input) IS NULL THEN NULL
        WHEN pf_normalize_profile_whitespace(input) ~ '[.!?]$' THEN
            upper(left(pf_normalize_profile_whitespace(input), 1)) || substr(pf_normalize_profile_whitespace(input), 2)
        ELSE
            upper(left(pf_normalize_profile_whitespace(input), 1)) || substr(pf_normalize_profile_whitespace(input), 2) || '.'
    END
$$;

CREATE OR REPLACE FUNCTION pf_profile_status_summary(input TEXT)
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT CASE
        WHEN pf_normalize_profile_whitespace(input) IS NULL THEN NULL
        WHEN lower(pf_normalize_profile_whitespace(input)) LIKE '%partial%' THEN 'Quedó parcialmente resuelto.'
        WHEN lower(pf_normalize_profile_whitespace(input)) LIKE '%parcial%' THEN 'Quedó parcialmente resuelto.'
        WHEN lower(pf_normalize_profile_whitespace(input)) LIKE '%pend%' THEN 'Quedó pendiente.'
        WHEN lower(pf_normalize_profile_whitespace(input)) LIKE '%block%' THEN 'Quedó bloqueado.'
        WHEN lower(pf_normalize_profile_whitespace(input)) LIKE '%risk%' THEN 'Siguió siendo un riesgo abierto.'
        ELSE NULL
    END
$$;

CREATE OR REPLACE FUNCTION pf_join_profile_sentences(parts TEXT[])
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT NULLIF(string_agg(pf_profile_sentence(part), ' '), '')
    FROM unnest(parts) AS part
    WHERE pf_normalize_profile_whitespace(part) IS NOT NULL
$$;

CREATE OR REPLACE FUNCTION pf_compact_profile_item(item JSONB, kind TEXT)
RETURNS TEXT
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
    raw TEXT := pf_profile_scalar_text(item);
    primary_text TEXT;
    secondary_text TEXT;
BEGIN
    CASE kind
        WHEN 'integrations' THEN
            primary_text := COALESCE(pf_profile_field(item, 'name'), raw);
            RETURN pf_join_profile_sentences(ARRAY[primary_text]);
        WHEN 'technical_decisions' THEN
            primary_text := COALESCE(pf_profile_field(item, 'decision'), raw);
            RETURN pf_join_profile_sentences(ARRAY[primary_text]);
        WHEN 'challenges' THEN
            primary_text := COALESCE(pf_profile_field(item, 'challenge'), raw);
            secondary_text := pf_profile_status_summary(pf_profile_field(item, 'status'));
            RETURN pf_join_profile_sentences(ARRAY[primary_text, secondary_text]);
        WHEN 'results' THEN
            primary_text := COALESCE(pf_profile_field(item, 'result'), raw);
            secondary_text := pf_profile_field(item, 'impact');
            RETURN pf_join_profile_sentences(ARRAY[primary_text, secondary_text]);
        WHEN 'timeline' THEN
            primary_text := COALESCE(pf_profile_field(item, 'phase'), raw);
            secondary_text := pf_profile_field(item, 'outcome');
            RETURN pf_join_profile_sentences(ARRAY[primary_text, secondary_text]);
        ELSE
            RETURN pf_join_profile_sentences(ARRAY[raw]);
    END CASE;
END;
$$;

CREATE OR REPLACE FUNCTION pf_compact_profile_list(items JSONB, kind TEXT)
RETURNS JSONB
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
    result JSONB := '[]'::jsonb;
    item JSONB;
    compacted TEXT;
BEGIN
    IF items IS NULL OR jsonb_typeof(items) <> 'array' THEN
        RETURN '[]'::jsonb;
    END IF;

    FOR item IN SELECT value FROM jsonb_array_elements(items) LOOP
        compacted := pf_compact_profile_item(item, kind);
        IF pf_normalize_profile_whitespace(compacted) IS NOT NULL THEN
            result := result || jsonb_build_array(compacted);
        END IF;
    END LOOP;

    RETURN result;
END;
$$;

WITH normalized AS (
    SELECT
        project_id,
        pf_compact_profile_list(integrations, 'integrations') AS integrations,
        pf_compact_profile_list(technical_decisions, 'technical_decisions') AS technical_decisions,
        pf_compact_profile_list(challenges, 'challenges') AS challenges,
        pf_compact_profile_list(results, 'results') AS results,
        pf_compact_profile_list(timeline, 'timeline') AS timeline
    FROM project_profiles
), changed AS (
    SELECT n.*
    FROM normalized n
    JOIN project_profiles p ON p.project_id = n.project_id
    WHERE p.integrations IS DISTINCT FROM n.integrations
       OR p.technical_decisions IS DISTINCT FROM n.technical_decisions
       OR p.challenges IS DISTINCT FROM n.challenges
       OR p.results IS DISTINCT FROM n.results
       OR p.timeline IS DISTINCT FROM n.timeline
)
UPDATE project_profiles p
SET integrations = c.integrations,
    technical_decisions = c.technical_decisions,
    challenges = c.challenges,
    results = c.results,
    timeline = c.timeline,
    updated_at = NOW()
FROM changed c
WHERE p.project_id = c.project_id;

DROP FUNCTION IF EXISTS pf_compact_profile_list(JSONB, TEXT);
DROP FUNCTION IF EXISTS pf_compact_profile_item(JSONB, TEXT);
DROP FUNCTION IF EXISTS pf_join_profile_sentences(TEXT[]);
DROP FUNCTION IF EXISTS pf_profile_status_summary(TEXT);
DROP FUNCTION IF EXISTS pf_profile_sentence(TEXT);
DROP FUNCTION IF EXISTS pf_profile_field(JSONB, TEXT);
DROP FUNCTION IF EXISTS pf_profile_scalar_text(JSONB);
DROP FUNCTION IF EXISTS pf_normalize_profile_key(TEXT);
DROP FUNCTION IF EXISTS pf_normalize_profile_whitespace(TEXT);
