package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	"github.com/marlonlyb/portfolioforge/model"
)

// SearchRepository implements search.SearchRepository against PostgreSQL
// with FTS, pg_trgm, and optional pgvector support.
type SearchRepository struct {
	db              *pgxpool.Pool
	embeddingProv   embedding.EmbeddingProvider
	semanticEnabled bool
}

// NewSearchRepository creates a new SearchRepository.
func NewSearchRepository(db *pgxpool.Pool, semanticEnabled bool, embeddingProv embedding.EmbeddingProvider) *SearchRepository {
	return &SearchRepository{db: db, embeddingProv: embeddingProv, semanticEnabled: semanticEnabled}
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
	query, args := r.buildListAllPublishedQuery(params)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.listAllPublished: %w", err)
	}
	defer rows.Close()

	results := make([]model.SearchResult, 0)
	for rows.Next() {
		var (
			project         model.Project
			solutionSummary string
		)
		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Slug,
			&project.Description,
			&project.Category,
			&project.ClientName,
			&project.IndustryType,
			&project.FinalProduct,
			&project.Images,
			&solutionSummary,
		); err != nil {
			return nil, fmt.Errorf("postgres.SearchRepository.listAllPublished scan: %w", err)
		}
		project.Active = true
		project.Status = "published"
		project.Profile = buildSearchProfile(project.ID, solutionSummary)

		results = append(results, model.SearchResult{Project: project})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.listAllPublished rows: %w", err)
	}

	return results, nil
}

func (r *SearchRepository) buildListAllPublishedQuery(params model.SearchParams) (string, []interface{}) {
	args := make([]interface{}, 0, 3)
	conditions := []string{"p.active = TRUE"}
	filterClause, filterArgs := buildProjectFilterClause(params, 1)
	if filterClause != "" {
		conditions = append(conditions, filterClause)
		args = append(args, filterArgs...)
	}

	return `
		SELECT p.id,
			COALESCE(NULLIF(p.name, ''), p.product_name) AS name,
			COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
			COALESCE(p.description, '') AS description,
			COALESCE(p.category, '') AS category,
			COALESCE(NULLIF(p.client_name, ''), p.brand, '') AS client_name,
			COALESCE(p.industry_type, '') AS industry_type,
			COALESCE(p.final_product, '') AS final_product,
			COALESCE(p.images, '[]'::jsonb) AS images,
			COALESCE(pp.solution_summary, '') AS solution_summary
		FROM products p
		LEFT JOIN project_profiles pp ON pp.project_id = p.id
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY p.created_at DESC, p.id DESC`, args
}

// searchWithQuery performs the CTE-based hybrid search with raw scores.
func (r *SearchRepository) searchWithQuery(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	query, args := r.buildSearchQuery(params)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres.SearchRepository.searchWithQuery: %w", err)
	}
	defer rows.Close()

	results := make([]model.SearchResult, 0)
	for rows.Next() {
		var (
			projectID       uuid.UUID
			name            string
			slug            string
			description     string
			category        string
			clientName      string
			industryType    string
			finalProduct    string
			images          json.RawMessage
			solutionSummary string
			lexicalScore    float64
			fuzzyScore      float64
			semanticScore   float64
		)

		if err := rows.Scan(
			&projectID,
			&name,
			&slug,
			&description,
			&category,
			&clientName,
			&industryType,
			&finalProduct,
			&images,
			&solutionSummary,
			&lexicalScore,
			&fuzzyScore,
			&semanticScore,
		); err != nil {
			return nil, fmt.Errorf("postgres.SearchRepository.searchWithQuery scan: %w", err)
		}

		results = append(results, model.SearchResult{
			Project: model.Project{
				ID:           projectID,
				Name:         name,
				Slug:         slug,
				Description:  description,
				Category:     category,
				ClientName:   clientName,
				IndustryType: industryType,
				FinalProduct: finalProduct,
				Images:       images,
				Active:       true,
				Status:       "published",
				Profile:      buildSearchProfile(projectID, solutionSummary),
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

// buildSearchQuery constructs the CTE-based SQL query with optional technology filtering.
func (r *SearchRepository) buildSearchQuery(params model.SearchParams) (string, []interface{}) {
	args := []interface{}{params.Query}
	filterClause, filterArgs := buildProjectFilterClause(params, len(args)+1)
	args = append(args, filterArgs...)

	query := `
WITH lexical AS (
    SELECT psd.project_id, ts_rank_cd(psd.search_document, query) AS lexical_score
    FROM project_search_documents psd, plainto_tsquery('simple', unaccent($1)) query
    WHERE psd.search_document @@ query
),
fuzzy AS (
    SELECT psd.project_id, similarity(psd.search_trgm, $1) AS fuzzy_score
    FROM project_search_documents psd
    WHERE psd.search_trgm % $1
)`

	semanticEnabled := r.semanticEnabled && len(params.QueryEmbedding) > 0
	if semanticEnabled {
		args = append(args, embeddingVector(params.QueryEmbedding))
		query += fmt.Sprintf(`,
semantic AS (
    SELECT psd.project_id, 1 - (psd.search_embedding <=> $%d::vector) AS semantic_score
    FROM project_search_documents psd
    WHERE psd.search_embedding IS NOT NULL
)`, len(args))
	}

	conditions := []string{"p.active = TRUE"}
	if semanticEnabled {
		conditions = append(conditions, "(l.lexical_score IS NOT NULL OR f.fuzzy_score IS NOT NULL OR s.semantic_score IS NOT NULL)")
	} else {
		conditions = append(conditions, "(l.lexical_score IS NOT NULL OR f.fuzzy_score IS NOT NULL)")
	}
	if filterClause != "" {
		conditions = append(conditions, filterClause)
	}

	semanticJoin := ""
	semanticScoreSelect := "0.0 AS semantic_score"
	orderBy := "(0.60 * COALESCE(l.lexical_score, 0) + 0.40 * COALESCE(f.fuzzy_score, 0))"
	if semanticEnabled {
		semanticJoin = "LEFT JOIN semantic s ON s.project_id = p.id"
		semanticScoreSelect = "COALESCE(s.semantic_score, 0) AS semantic_score"
		orderBy = "(0.45 * COALESCE(l.lexical_score, 0) + 0.25 * COALESCE(f.fuzzy_score, 0) + 0.30 * COALESCE(s.semantic_score, 0))"
	}

	query += `
SELECT p.id,
    COALESCE(NULLIF(p.name, ''), p.product_name) AS name,
    COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
    COALESCE(p.description, '') AS description,
    COALESCE(p.category, '') AS category,
    COALESCE(NULLIF(p.client_name, ''), p.brand, '') AS client_name,
    COALESCE(p.industry_type, '') AS industry_type,
    COALESCE(p.final_product, '') AS final_product,
    COALESCE(p.images, '[]'::jsonb) AS images,
    COALESCE(pp.solution_summary, '') AS solution_summary,
    COALESCE(l.lexical_score, 0) AS lexical_score,
    COALESCE(f.fuzzy_score, 0) AS fuzzy_score,
    ` + semanticScoreSelect + `
FROM products p
LEFT JOIN project_profiles pp ON pp.project_id = p.id
LEFT JOIN lexical l ON l.project_id = p.id
LEFT JOIN fuzzy f ON f.project_id = p.id
` + semanticJoin + `
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY ` + orderBy + ` DESC, p.id DESC`

	return query, args
}

// RefreshSearchDocument recomposes the search document for a single project.
func (r *SearchRepository) RefreshSearchDocument(ctx context.Context, projectID uuid.UUID) error {
	contentText, changed, err := RefreshProjectSearchDocument(ctx, r.db, projectID)
	if err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshSearchDocument update: %w", err)
	}
	if !changed || !r.semanticEnabled || r.embeddingProv == nil {
		return nil
	}

	embeddingVec, err := r.embeddingProv.Generate(ctx, contentText)
	if err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshSearchDocument embedding: %w", err)
	}
	if len(embeddingVec) == 0 {
		return nil
	}

	if err := UpdateProjectSearchEmbedding(ctx, r.db, projectID, embeddingVec); err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshSearchDocument persist embedding: %w", err)
	}
	return nil
}

// RefreshAllDocuments recomposes search documents for all projects.
func (r *SearchRepository) RefreshAllDocuments(ctx context.Context) error {
	rows, err := r.db.Query(ctx, `SELECT id FROM products WHERE active = TRUE ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshAllDocuments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var projectID uuid.UUID
		if err := rows.Scan(&projectID); err != nil {
			return fmt.Errorf("postgres.SearchRepository.RefreshAllDocuments scan: %w", err)
		}
		if err := r.RefreshSearchDocument(ctx, projectID); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres.SearchRepository.RefreshAllDocuments rows: %w", err)
	}

	return nil
}

func buildProjectFilterClause(params model.SearchParams, startIndex int) (string, []interface{}) {
	conditions := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	placeholder := startIndex

	if strings.TrimSpace(params.Category) != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(COALESCE(p.category, '')) = LOWER($%d)", placeholder))
		args = append(args, strings.TrimSpace(params.Category))
		placeholder++
	}
	if strings.TrimSpace(params.Client) != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(COALESCE(NULLIF(p.client_name, ''), p.brand, '')) = LOWER($%d)", placeholder))
		args = append(args, strings.TrimSpace(params.Client))
		placeholder++
	}
	if len(params.Technologies) > 0 {
		conditions = append(conditions, fmt.Sprintf(`p.id IN (
			SELECT pt.project_id
			FROM project_technologies pt
			JOIN technologies t ON t.id = pt.technology_id
			WHERE t.slug = ANY($%d)
			GROUP BY pt.project_id
			HAVING COUNT(DISTINCT t.slug) = %d
		)`, placeholder, len(params.Technologies)))
		args = append(args, params.Technologies)
	}

	return strings.Join(conditions, " AND "), args
}

func buildSearchProfile(projectID uuid.UUID, solutionSummary string) *model.ProjectProfile {
	if strings.TrimSpace(solutionSummary) == "" {
		return nil
	}
	return &model.ProjectProfile{ProjectID: projectID, SolutionSummary: solutionSummary}
}

func embeddingVector(values []float32) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, fmt.Sprintf("%f", value))
	}
	return "[" + strings.Join(parts, ",") + "]"
}
