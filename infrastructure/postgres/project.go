package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

// ProjectRepository implements project.ProjectReader against the PostgreSQL `products` table.
// During the transition period, products.brand maps to Project.ClientName.
type ProjectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new ProjectRepository with the given connection pool.
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetByID returns a single project by its ID, including profile and technologies.
func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Project, error) {
	var p model.Project
	var profile model.ProjectProfile
	var profileUpdatedAt sql.NullInt64

	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.name,
			COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
			p.description, p.category, COALESCE(p.brand, '') AS client_name,
			p.active, (COALESCE(NULLIF(trim(p.source_markdown_url), ''), '') <> '') AS assistant_available, p.images, p.created_at, COALESCE(p.updated_at, 0) AS updated_at,
			COALESCE(pp.business_goal, ''), COALESCE(pp.problem_statement, ''),
			COALESCE(pp.solution_summary, ''), COALESCE(pp.delivery_scope, ''), COALESCE(pp.responsibility_scope, ''), COALESCE(pp.architecture, ''),
			COALESCE(pp.integrations, 'null'), COALESCE(pp.ai_usage, ''),
			COALESCE(pp.technical_decisions, 'null'), COALESCE(pp.challenges, 'null'),
			COALESCE(pp.results, 'null'), COALESCE(pp.metrics, 'null'),
			COALESCE(pp.timeline, 'null'), COALESCE(EXTRACT(EPOCH FROM pp.updated_at)::bigint, 0)
		FROM products p
		LEFT JOIN project_profiles pp ON pp.project_id = p.id
		WHERE p.id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Category, &p.ClientName,
		&p.Active, &p.AssistantAvailable, &p.Images, &p.CreatedAt, &p.UpdatedAt,
		&profile.BusinessGoal, &profile.ProblemStatement,
		&profile.SolutionSummary, &profile.DeliveryScope, &profile.ResponsibilityScope, &profile.Architecture,
		&profile.Integrations, &profile.AIUsage,
		&profile.TechnicalDecisions, &profile.Challenges,
		&profile.Results, &profile.Metrics,
		&profile.Timeline, &profileUpdatedAt,
	)
	if err != nil {
		return model.Project{}, fmt.Errorf("postgres.ProjectRepository.GetByID: %w", err)
	}

	p.Status = "published"
	if !p.Active {
		p.Status = "draft"
	}

	profile.ProjectID = p.ID
	if profileUpdatedAt.Valid {
		profile.UpdatedAt = profileUpdatedAt.Int64
	}
	p.Profile = &profile

	techs, err := r.fetchTechnologies(ctx, p.ID)
	if err == nil {
		p.Technologies = techs
	}

	media, err := fetchProjectMedia(ctx, r.db, p.ID)
	if err == nil {
		p.Media = media
		rebuildProjectImages(&p)
	}

	return p, nil
}

// GetBySlug returns a single published project by its slug, including profile and technologies.
func (r *ProjectRepository) GetBySlug(ctx context.Context, slug string) (model.Project, error) {
	var p model.Project
	var profile model.ProjectProfile
	var profileUpdatedAt sql.NullInt64

	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.name,
			COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
			p.description, p.category, COALESCE(p.brand, '') AS client_name,
			p.active, (COALESCE(NULLIF(trim(p.source_markdown_url), ''), '') <> '') AS assistant_available, p.images, p.created_at, COALESCE(p.updated_at, 0) AS updated_at,
			COALESCE(pp.business_goal, ''), COALESCE(pp.problem_statement, ''),
			COALESCE(pp.solution_summary, ''), COALESCE(pp.delivery_scope, ''), COALESCE(pp.responsibility_scope, ''), COALESCE(pp.architecture, ''),
			COALESCE(pp.integrations, 'null'), COALESCE(pp.ai_usage, ''),
			COALESCE(pp.technical_decisions, 'null'), COALESCE(pp.challenges, 'null'),
			COALESCE(pp.results, 'null'), COALESCE(pp.metrics, 'null'),
			COALESCE(pp.timeline, 'null'), COALESCE(EXTRACT(EPOCH FROM pp.updated_at)::bigint, 0)
		FROM products p
		LEFT JOIN project_profiles pp ON pp.project_id = p.id
		WHERE p.slug = $1 AND p.active = TRUE`, slug).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Category, &p.ClientName,
		&p.Active, &p.AssistantAvailable, &p.Images, &p.CreatedAt, &p.UpdatedAt,
		&profile.BusinessGoal, &profile.ProblemStatement,
		&profile.SolutionSummary, &profile.DeliveryScope, &profile.ResponsibilityScope, &profile.Architecture,
		&profile.Integrations, &profile.AIUsage,
		&profile.TechnicalDecisions, &profile.Challenges,
		&profile.Results, &profile.Metrics,
		&profile.Timeline, &profileUpdatedAt,
	)
	if err != nil {
		return model.Project{}, fmt.Errorf("postgres.ProjectRepository.GetBySlug: %w", err)
	}

	p.Status = "published"

	profile.ProjectID = p.ID
	if profileUpdatedAt.Valid {
		profile.UpdatedAt = profileUpdatedAt.Int64
	}
	p.Profile = &profile

	techs, err := r.fetchTechnologies(ctx, p.ID)
	if err == nil {
		p.Technologies = techs
	}

	media, err := fetchProjectMedia(ctx, r.db, p.ID)
	if err == nil {
		p.Media = media
		rebuildProjectImages(&p)
	}

	return p, nil
}

// ListPublished returns all published (active) projects, including profile and technologies.
func (r *ProjectRepository) ListPublished(ctx context.Context) ([]model.Project, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.name,
			COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
			p.description, p.category, COALESCE(p.brand, '') AS client_name,
			p.active, (COALESCE(NULLIF(trim(p.source_markdown_url), ''), '') <> '') AS assistant_available, p.images, p.created_at, COALESCE(p.updated_at, 0) AS updated_at
		FROM products p
		WHERE p.active = TRUE
		ORDER BY p.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("postgres.ProjectRepository.ListPublished: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.Category, &p.ClientName,
			&p.Active, &p.AssistantAvailable, &p.Images, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres.ProjectRepository.ListPublished scan: %w", err)
		}
		p.Status = "published"
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.ProjectRepository.ListPublished rows: %w", err)
	}

	// For each project, fetch profile and technologies
	for i := range projects {
		r.fetchProfile(ctx, &projects[i])
		techs, err := r.fetchTechnologies(ctx, projects[i].ID)
		if err == nil {
			projects[i].Technologies = techs
		}
		media, err := fetchProjectMedia(ctx, r.db, projects[i].ID)
		if err == nil {
			projects[i].Media = media
			rebuildProjectImages(&projects[i])
		}
	}

	return projects, nil
}

// GetTechnologiesByProjectID returns technologies for a specific project.
func (r *ProjectRepository) GetTechnologiesByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Technology, error) {
	return r.fetchTechnologies(ctx, projectID)
}

func (r *ProjectRepository) GetAssistantContextBySlug(ctx context.Context, slug string) (model.ProjectAssistantContext, error) {
	var project model.ProjectAssistantContext
	err := r.db.QueryRow(ctx, `
		SELECT p.id,
			COALESCE(NULLIF(p.name, ''), p.product_name) AS name,
			COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
			p.active,
			COALESCE(p.source_markdown_url, '') AS source_markdown_url
		FROM products p
		WHERE COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) = $1`, slug,
	).Scan(&project.ID, &project.Name, &project.Slug, &project.Active, &project.SourceMarkdownURL)
	if err != nil {
		return model.ProjectAssistantContext{}, fmt.Errorf("postgres.ProjectRepository.GetAssistantContextBySlug: %w", err)
	}

	return project, nil
}

// fetchProfile populates the ProjectProfile for a given project.
func (r *ProjectRepository) fetchProfile(ctx context.Context, p *model.Project) {
	var profile model.ProjectProfile
	var profileUpdatedAt sql.NullInt64

	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(business_goal, ''), COALESCE(problem_statement, ''),
			COALESCE(solution_summary, ''), COALESCE(delivery_scope, ''), COALESCE(responsibility_scope, ''), COALESCE(architecture, ''),
			COALESCE(integrations, 'null'), COALESCE(ai_usage, ''),
			COALESCE(technical_decisions, 'null'), COALESCE(challenges, 'null'),
			COALESCE(results, 'null'), COALESCE(metrics, 'null'),
			COALESCE(timeline, 'null'), COALESCE(EXTRACT(EPOCH FROM updated_at)::bigint, 0)
		FROM project_profiles
		WHERE project_id = $1`, p.ID).Scan(
		&profile.BusinessGoal, &profile.ProblemStatement,
		&profile.SolutionSummary, &profile.DeliveryScope, &profile.ResponsibilityScope, &profile.Architecture,
		&profile.Integrations, &profile.AIUsage,
		&profile.TechnicalDecisions, &profile.Challenges,
		&profile.Results, &profile.Metrics,
		&profile.Timeline, &profileUpdatedAt,
	)
	if err != nil {
		// No profile found — leave nil
		return
	}

	profile.ProjectID = p.ID
	if profileUpdatedAt.Valid {
		profile.UpdatedAt = profileUpdatedAt.Int64
	}
	p.Profile = &profile
}

// fetchTechnologies returns the technologies associated with a project.
func (r *ProjectRepository) fetchTechnologies(ctx context.Context, projectID uuid.UUID) ([]model.Technology, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.name, t.slug, t.category, COALESCE(t.icon, ''), COALESCE(t.color, '')
		FROM project_technologies pt
		JOIN technologies t ON t.id = pt.technology_id
		WHERE pt.project_id = $1
		ORDER BY t.name`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var techs []model.Technology
	for rows.Next() {
		var t model.Technology
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Category, &t.Icon, &t.Color); err != nil {
			return nil, err
		}
		techs = append(techs, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return techs, nil
}

func rebuildProjectImages(project *model.Project) {
	legacyImages := make([]string, 0)
	if len(project.Images) > 0 {
		_ = json.Unmarshal(project.Images, &legacyImages)
	}

	rebuilt, err := json.Marshal(model.BuildProjectImageList(project.Media, legacyImages))
	if err != nil {
		return
	}

	project.Images = rebuilt
}
