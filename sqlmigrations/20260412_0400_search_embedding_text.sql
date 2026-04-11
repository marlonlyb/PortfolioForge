-- Adds raw composed text helper for embedding input and ensures
-- problem_statement participates in search document composition.

CREATE OR REPLACE FUNCTION compose_project_embedding_text(p_id UUID)
RETURNS TEXT AS $$
SELECT trim(concat_ws(' ',
    COALESCE(p.name, ''),
    COALESCE(p.brand, ''),
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
GROUP BY p.id, p.name, p.brand, p.description, pp.solution_summary,
         pp.architecture, pp.business_goal, pp.problem_statement, pp.ai_usage;
$$ LANGUAGE SQL STABLE;
