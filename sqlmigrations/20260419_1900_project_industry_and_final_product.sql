ALTER TABLE products
    ADD COLUMN IF NOT EXISTS industry_type VARCHAR(160),
    ADD COLUMN IF NOT EXISTS final_product VARCHAR(160);

ALTER TABLE products
    DROP CONSTRAINT IF EXISTS products_industry_type_ck;

ALTER TABLE products
    ALTER COLUMN industry_type TYPE VARCHAR(160),
    ALTER COLUMN final_product TYPE VARCHAR(160);

UPDATE products
SET industry_type = CASE lower(trim(industry_type))
    WHEN 'food' THEN 'alimentación'
    WHEN 'beverages' THEN 'bebidas'
    WHEN 'construction' THEN 'construcción'
    WHEN 'plastics' THEN 'plásticos'
    WHEN 'cardboard' THEN 'cartón'
    WHEN 'metalworking' THEN 'metalurgia'
    WHEN 'material-handling' THEN 'movimiento de materiales'
    WHEN 'industrial-services' THEN 'servicios industriales'
    WHEN 'other' THEN 'otras industrias'
    ELSE regexp_replace(trim(COALESCE(industry_type, '')), '\s+', ' ', 'g')
END
WHERE industry_type IS NOT NULL;

CREATE OR REPLACE FUNCTION collect_project_localized_field_text(p_id UUID, p_field_key TEXT)
RETURNS TEXT AS $$
    SELECT COALESCE(string_agg(DISTINCT trim(both '"' from value::text), ' ' ORDER BY trim(both '"' from value::text)), '')
    FROM project_localizations
    WHERE project_id = p_id
      AND field_key = p_field_key
      AND jsonb_typeof(value) = 'string'
      AND trim(both '"' from value::text) <> '';
$$ LANGUAGE SQL STABLE;

CREATE OR REPLACE FUNCTION compose_project_embedding_text(p_id UUID)
RETURNS TEXT AS $$
SELECT trim(concat_ws(' ',
    COALESCE(p.name, ''),
    COALESCE(p.brand, ''),
    COALESCE(p.industry_type, ''),
    collect_project_localized_field_text(p.id, 'industry_type'),
    COALESCE(p.final_product, ''),
    collect_project_localized_field_text(p.id, 'final_product'),
    COALESCE(p.description, ''),
    COALESCE(pp.solution_summary, ''),
    COALESCE(pp.architecture, ''),
    COALESCE(pp.business_goal, ''),
    COALESCE(pp.problem_statement, ''),
    COALESCE(pp.ai_usage, ''),
    COALESCE(tech.names, '')
))
FROM products p
LEFT JOIN project_profiles pp ON pp.project_id = p.id
LEFT JOIN LATERAL (
    SELECT string_agg(t.name, ' ' ORDER BY t.name) AS names
    FROM project_technologies pt
    JOIN technologies t ON t.id = pt.technology_id
    WHERE pt.project_id = p.id
) tech ON TRUE
WHERE p.id = p_id;
$$ LANGUAGE SQL STABLE;

CREATE OR REPLACE FUNCTION compose_project_search_doc(p_id UUID)
RETURNS TSVECTOR AS $$
SELECT
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.solution_summary, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.name, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.brand, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.final_product, ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(collect_project_localized_field_text(p.id, 'final_product'), ''))), 'A') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.industry_type, ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(collect_project_localized_field_text(p.id, 'industry_type'), ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.architecture, ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(p.description, ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(string_agg(t.name, ' '), ''))), 'B') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.business_goal, ''))), 'C') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.problem_statement, ''))), 'C') ||
    setweight(to_tsvector('simple', unaccent(COALESCE(pp.ai_usage, ''))), 'C')
FROM products p
LEFT JOIN project_profiles pp ON pp.project_id = p.id
LEFT JOIN project_technologies pt ON pt.project_id = p.id
LEFT JOIN technologies t ON t.id = pt.technology_id
WHERE p.id = p_id
GROUP BY p.id, p.name, p.brand, p.industry_type, p.final_product, p.description, pp.solution_summary,
         pp.architecture, pp.business_goal, pp.problem_statement, pp.ai_usage;
$$ LANGUAGE SQL STABLE;

CREATE OR REPLACE FUNCTION compose_project_search_trgm(p_id UUID)
RETURNS TEXT AS $$
SELECT trim(concat_ws(' ',
    COALESCE(p.name, ''),
    COALESCE(p.brand, ''),
    COALESCE(p.industry_type, ''),
    collect_project_localized_field_text(p.id, 'industry_type'),
    COALESCE(p.final_product, ''),
    collect_project_localized_field_text(p.id, 'final_product'),
    COALESCE(string_agg(t.name, ' '), '')
))
FROM products p
LEFT JOIN project_technologies pt ON pt.project_id = p.id
LEFT JOIN technologies t ON t.id = pt.technology_id
WHERE p.id = p_id
GROUP BY p.id, p.name, p.brand, p.industry_type, p.final_product;
$$ LANGUAGE SQL STABLE;
