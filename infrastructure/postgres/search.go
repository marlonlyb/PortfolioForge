package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

// SearchRepository implements search.SearchRepository against PostgreSQL
// with FTS, pg_trgm, and optional pgvector support.
type SearchRepository struct {
	db              *pgxpool.Pool
	semanticEnabled bool
}

// NewSearchRepository creates a new SearchRepository.
func NewSearchRepository(db *pgxpool.Pool, semanticEnabled bool) *SearchRepository {
	return &SearchRepository{db: db, semanticEnabled: semanticEnabled}
}

// Search performs a hybrid retrieval search combining lexical (FTS),
// fuzzy (pg_trgm), and optionally semantic (pgvector) signals.
// When params.Query is empty, it returns all published projects (unranked).
func (r *SearchRepository) Search(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	if params.Query == "" {
		return r.listAllPublished(ctx, params)
	}
	return r.searchWithQuery(ctx, params)
}

// listAllPublished returns all published projects when no query is provided.
func (r *SearchRepository) listAllPublished(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	limit := params.PageSize
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.name, p.description, p.category, COALESCE(p.brand, '') AS client_name
		FROM products p
		WHERE p.active = TRUE
		ORDER BY p.created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.listAllPublished: %w", err)
	}
	defer rows.Close()

	var results []model.SearchResult
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Category, &p.ClientName); err != nil {
			return nil, fmt.Errorf("postgres.SearchRepository.listAllPublished scan: %w", err)
		}
		p.Active = true
		p.Status = "published"

		results = append(results, model.SearchResult{
			Project: p,
			Score:   0,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.listAllPublished rows: %w", err)
	}

	return results, nil
}

// searchWithQuery performs the CTE-based hybrid search with raw scores.
func (r *SearchRepository) searchWithQuery(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	query, args := r.buildSearchQuery(params)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.searchWithQuery: %w", err)
	}
	defer rows.Close()

	var results []model.SearchResult
	for rows.Next() {
		var (
			projectID     uuid.UUID
			name          string
			description   string
			category      string
			clientName    string
			lexicalScore  float64
			fuzzyScore    float64
			semanticScore float64
		)

		if err := rows.Scan(
			&projectID, &name, &description, &category, &clientName,
			&lexicalScore, &fuzzyScore, &semanticScore,
		); err != nil {
			return nil, fmt.Errorf("postgres.SearchRepository.searchWithQuery scan: %w", err)
		}

		results = append(results, model.SearchResult{
			Project: model.Project{
				ID:          projectID,
				Name:        name,
				Description: description,
				Category:    category,
				ClientName:  clientName,
				Active:      true,
				Status:      "published",
			},
			LexicalScore:  lexicalScore,
			FuzzyScore:    fuzzyScore,
			SemanticScore: semanticScore,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.searchWithQuery rows: %w", err)
	}

	return results, nil
}

// buildSearchQuery constructs the CTE-based SQL query with optional technology AND filtering.
// When semantic search is disabled, the semantic CTE is omitted and weights are re-normalized.
// When technologies filter is provided, only projects having ALL specified technology slugs are returned.
func (r *SearchRepository) buildSearchQuery(params model.SearchParams) (string, []interface{}) {
	limit := params.PageSize
	if limit <= 0 {
		limit = 20
	}

	// Build technology AND filter clause if technologies are specified
	techFilter := ""
	var techArgs []interface{}
	if len(params.Technologies) > 0 {
		// We'll use a subquery that ensures the project has ALL specified technology slugs.
		// The subquery is appended via param placeholders that we'll assign after the base args.
		techFilter = fmt.Sprintf(
			` AND p.id IN (SELECT pt.project_id FROM project_technologies pt JOIN technologies t ON t.id = pt.technology_id WHERE t.slug = ANY($N) GROUP BY pt.project_id HAVING COUNT(DISTINCT t.slug) = %d)`,
			len(params.Technologies),
		)
		techArgs = append(techArgs, params.Technologies)
	}

	if r.semanticEnabled {
		// With semantic: query text ($1), embedding vector ($2), limit ($3)
		// Technologies: if present, $4 = tech slugs array
		args := []interface{}{params.Query, nil, limit}
		if len(techArgs) > 0 {
			args = append(args, techArgs...)
			// Replace $N with the correct placeholder
			techFilter = strings.Replace(techFilter, "$N", fmt.Sprintf("$%d", len(args)-len(techArgs)+1), 1)
		}

		return `
WITH lexical AS (
    SELECT psd.project_id, ts_rank_cd(psd.search_document, query) AS lexical_score
    FROM project_search_documents psd, plainto_tsquery('simple', unaccent($1)) query
    WHERE psd.search_document @@ query
),
fuzzy AS (
    SELECT psd.project_id, similarity(psd.search_trgm, $1) AS fuzzy_score
    FROM project_search_documents psd
    WHERE psd.search_trgm % $1
),
semantic AS (
    SELECT psd.project_id, 1 - (psd.search_embedding <=> $2::vector) AS semantic_score
    FROM project_search_documents psd
    WHERE psd.search_embedding IS NOT NULL
)
SELECT p.id, p.name, p.description, p.category, COALESCE(p.brand, '') AS client_name,
    COALESCE(l.lexical_score, 0) AS lexical_score,
    COALESCE(f.fuzzy_score, 0) AS fuzzy_score,
    COALESCE(s.semantic_score, 0) AS semantic_score
FROM products p
LEFT JOIN lexical l ON l.project_id = p.id
LEFT JOIN fuzzy f ON f.project_id = p.id
LEFT JOIN semantic s ON s.project_id = p.id
WHERE p.active = TRUE
  AND (l.lexical_score IS NOT NULL OR f.fuzzy_score IS NOT NULL OR s.semantic_score IS NOT NULL)` + techFilter + `
ORDER BY (0.45 * COALESCE(l.lexical_score, 0) + 0.25 * COALESCE(f.fuzzy_score, 0) + 0.30 * COALESCE(s.semantic_score, 0)) DESC
LIMIT $3`, args
	}

	// Without semantic: query text ($1), limit ($2)
	// Technologies: if present, $3 = tech slugs array
	args := []interface{}{params.Query, limit}
	if len(techArgs) > 0 {
		args = append(args, techArgs...)
		techFilter = strings.Replace(techFilter, "$N", fmt.Sprintf("$%d", len(args)-len(techArgs)+1), 1)
	}

	return `
WITH lexical AS (
    SELECT psd.project_id, ts_rank_cd(psd.search_document, query) AS lexical_score
    FROM project_search_documents psd, plainto_tsquery('simple', unaccent($1)) query
    WHERE psd.search_document @@ query
),
fuzzy AS (
    SELECT psd.project_id, similarity(psd.search_trgm, $1) AS fuzzy_score
    FROM project_search_documents psd
    WHERE psd.search_trgm % $1
)
SELECT p.id, p.name, p.description, p.category, COALESCE(p.brand, '') AS client_name,
    COALESCE(l.lexical_score, 0) AS lexical_score,
    COALESCE(f.fuzzy_score, 0) AS fuzzy_score,
    0.0 AS semantic_score
FROM products p
LEFT JOIN lexical l ON l.project_id = p.id
LEFT JOIN fuzzy f ON f.project_id = p.id
WHERE p.active = TRUE
  AND (l.lexical_score IS NOT NULL OR f.fuzzy_score IS NOT NULL)` + techFilter + `
ORDER BY (0.60 * COALESCE(l.lexical_score, 0) + 0.40 * COALESCE(f.fuzzy_score, 0)) DESC
LIMIT $2`, args
}

// RefreshSearchDocument recomposes the search document for a single project
// by calling the compose_project_search_doc() database function.
// It computes a SHA-256 hash of the composed content and compares it with the
// stored hash. If unchanged, it skips re-embedding and returns early.
func (r *SearchRepository) RefreshSearchDocument(ctx context.Context, projectID uuid.UUID) error {
	_, _, err := RefreshProjectSearchDocument(ctx, r.db, projectID)
	if err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshSearchDocument update: %w", err)
	}
	return nil
}

// RefreshAllDocuments recomposes search documents for all projects.
func (r *SearchRepository) RefreshAllDocuments(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		UPDATE project_search_documents
		SET search_document = compose_project_search_doc(project_id),
		    search_trgm = compose_project_search_trgm(project_id)
		WHERE project_id IN (SELECT id FROM products WHERE active = TRUE)`)
	if err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshAllDocuments: %w", err)
	}
	return nil
}
