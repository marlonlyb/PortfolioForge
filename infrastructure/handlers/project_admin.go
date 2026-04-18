package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	projectPorts "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/infrastructure/postgres"
	"github.com/marlonlyb/portfolioforge/model"
)

type projectAdminTx interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type ProjectAdminHandler struct {
	beginTx         func(context.Context) (projectAdminTx, error)
	embeddingProv   embedding.EmbeddingProvider
	semanticEnabled bool
	projectRepo     projectPorts.ProjectReader
	localization    *localization.Service
}

func NewProjectAdminHandler(
	db *pgxpool.Pool,
	embeddingProv embedding.EmbeddingProvider,
	semanticEnabled bool,
	projectRepo projectPorts.ProjectReader,
	localizationService *localization.Service,
) *ProjectAdminHandler {
	return &ProjectAdminHandler{
		beginTx: func(ctx context.Context) (projectAdminTx, error) {
			return db.Begin(ctx)
		},
		embeddingProv:   embeddingProv,
		semanticEnabled: semanticEnabled,
		projectRepo:     projectRepo,
		localization:    localizationService,
	}
}

type EnrichmentProfileReq struct {
	BusinessGoal        string          `json:"business_goal"`
	ProblemStatement    string          `json:"problem_statement"`
	SolutionSummary     string          `json:"solution_summary"`
	DeliveryScope       string          `json:"delivery_scope"`
	ResponsibilityScope string          `json:"responsibility_scope"`
	Architecture        string          `json:"architecture"`
	Integrations        json.RawMessage `json:"integrations"`
	AIUsage             string          `json:"ai_usage"`
	TechnicalDecisions  json.RawMessage `json:"technical_decisions"`
	Challenges          json.RawMessage `json:"challenges"`
	Results             json.RawMessage `json:"results"`
	Metrics             json.RawMessage `json:"metrics"`
	Timeline            json.RawMessage `json:"timeline"`
}

type EnrichmentReq struct {
	Profile       EnrichmentProfileReq `json:"profile"`
	TechnologyIDs []string             `json:"technology_ids"`
}

func (h *ProjectAdminHandler) UpdateProjectEnrichment(c echo.Context) error {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		return response.ContractError(400, "validation_error", "ID de proyecto inválido")
	}

	var req EnrichmentReq
	if err := c.Bind(&req); err != nil {
		return response.ContractError(400, "validation_error", "Datos de entrada inválidos")
	}

	if err := normalizeEnrichmentProfile(&req.Profile); err != nil {
		return response.ContractError(400, "validation_error", err.Error())
	}

	ctx := c.Request().Context()
	previousProject, err := h.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "Error leyendo el estado actual del proyecto: "+err.Error())
	}

	tx, err := h.beginTx(ctx)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "Error actualizando información del proyecto: "+err.Error())
	}
	defer tx.Rollback(ctx)

	err = h.executeEnrichmentTx(ctx, tx, projectID, req)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "Error actualizando información del proyecto: "+err.Error())
	}

	rawSearchText, _, err := postgres.RefreshProjectSearchDocument(ctx, tx, projectID)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "Error al actualizar documento de búsqueda: "+err.Error())
	}

	if h.semanticEnabled {
		if h.embeddingProv == nil {
			return response.ContractError(500, "unexpected_error", "No hay proveedor de embeddings configurado para reindexar el proyecto")
		}

		embeddingVec, embeddingErr := h.embeddingProv.Generate(ctx, rawSearchText)
		if embeddingErr != nil {
			return response.ContractError(500, "unexpected_error", "Error generando embedding del proyecto: "+embeddingErr.Error())
		}

		if err = postgres.UpdateProjectSearchEmbedding(ctx, tx, projectID, embeddingVec); err != nil {
			return response.ContractError(500, "unexpected_error", "Error guardando embedding del proyecto: "+err.Error())
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return response.ContractError(500, "unexpected_error", "Error confirmando actualización del proyecto: "+err.Error())
	}

	nextProject := previousProject
	nextProject.Profile = &model.ProjectProfile{
		ProjectID:           projectID,
		BusinessGoal:        req.Profile.BusinessGoal,
		ProblemStatement:    req.Profile.ProblemStatement,
		SolutionSummary:     req.Profile.SolutionSummary,
		DeliveryScope:       req.Profile.DeliveryScope,
		ResponsibilityScope: req.Profile.ResponsibilityScope,
		Architecture:        req.Profile.Architecture,
		Integrations:        req.Profile.Integrations,
		AIUsage:             req.Profile.AIUsage,
		TechnicalDecisions:  req.Profile.TechnicalDecisions,
		Challenges:          req.Profile.Challenges,
		Results:             req.Profile.Results,
		Metrics:             req.Profile.Metrics,
		Timeline:            req.Profile.Timeline,
	}

	if err := h.localization.SyncFromSpanish(ctx, projectID, localization.BuildProjectFieldMap(previousProject), localization.BuildProjectFieldMap(nextProject)); err != nil {
		return response.ContractError(500, "unexpected_error", "Error actualizando traducciones automáticas: "+err.Error())
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"message": "Proyecto enriquecido exitosamente",
	}))
}

func (h *ProjectAdminHandler) GetProjectTranslations(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "ID de proyecto inválido")
	}

	project, err := h.projectRepo.GetByID(c.Request().Context(), projectID)
	if err != nil {
		return response.ContractError(404, "not_found", "Proyecto no encontrado")
	}

	payload, err := h.localization.BuildAdminTranslationsResponse(c.Request().Context(), project)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener las traducciones del proyecto")
	}

	return c.JSON(response.ContractOK(payload))
}

func (h *ProjectAdminHandler) SaveProjectTranslations(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "ID de proyecto inválido")
	}

	locale := localization.NormalizeLocale(c.Param("locale"))
	if !model.IsSupportedTranslationLocale(locale) {
		return response.ContractError(400, "validation_error", "Idioma de traducción inválido")
	}

	var req struct {
		Fields map[string]json.RawMessage `json:"fields"`
	}
	if err := c.Bind(&req); err != nil {
		return response.ContractError(400, "validation_error", "Datos de traducción inválidos")
	}

	if err := h.localization.SaveManualTranslations(c.Request().Context(), projectID, locale, req.Fields); err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible guardar las traducciones manuales")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"message": "Traducciones guardadas"}))
}

func (h *ProjectAdminHandler) executeEnrichmentTx(ctx context.Context, tx projectAdminTx, projectID uuid.UUID, req EnrichmentReq) error {
	var err error

	// Upsert project_profiles
	// NOTE: Because there might not be a row yet, we use INSERT ON CONFLICT
	upsertQuery := `
		INSERT INTO project_profiles (
			project_id, business_goal, problem_statement, solution_summary,
			delivery_scope, responsibility_scope, architecture, integrations, ai_usage, technical_decisions,
			challenges, results, metrics, timeline, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW())
		ON CONFLICT (project_id) DO UPDATE SET
			business_goal = EXCLUDED.business_goal,
			problem_statement = EXCLUDED.problem_statement,
			solution_summary = EXCLUDED.solution_summary,
			delivery_scope = EXCLUDED.delivery_scope,
			responsibility_scope = EXCLUDED.responsibility_scope,
			architecture = EXCLUDED.architecture,
			integrations = EXCLUDED.integrations,
			ai_usage = EXCLUDED.ai_usage,
			technical_decisions = EXCLUDED.technical_decisions,
			challenges = EXCLUDED.challenges,
			results = EXCLUDED.results,
			metrics = EXCLUDED.metrics,
			timeline = EXCLUDED.timeline,
			updated_at = EXCLUDED.updated_at
	`
	_, err = tx.Exec(ctx, upsertQuery,
		projectID,
		req.Profile.BusinessGoal,
		req.Profile.ProblemStatement,
		req.Profile.SolutionSummary,
		req.Profile.DeliveryScope,
		req.Profile.ResponsibilityScope,
		req.Profile.Architecture,
		req.Profile.Integrations,
		req.Profile.AIUsage,
		req.Profile.TechnicalDecisions,
		req.Profile.Challenges,
		req.Profile.Results,
		req.Profile.Metrics,
		req.Profile.Timeline,
	)
	if err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}

	// Delete old technologies
	_, err = tx.Exec(ctx, "DELETE FROM project_technologies WHERE project_id = $1", projectID)
	if err != nil {
		return fmt.Errorf("delete technologies: %w", err)
	}

	// Insert new technologies
	if len(req.TechnologyIDs) > 0 {
		var valArgs []string
		var args []interface{}
		args = append(args, projectID) // $1 is projectID

		for i, techIDStr := range req.TechnologyIDs {
			techID, parseErr := uuid.Parse(techIDStr)
			if parseErr == nil {
				paramIndex := i + 2
				valArgs = append(valArgs, fmt.Sprintf("($1, $%d)", paramIndex))
				args = append(args, techID)
			}
		}

		if len(valArgs) > 0 {
			insertQuery := fmt.Sprintf("INSERT INTO project_technologies (project_id, technology_id) VALUES %s", strings.Join(valArgs, ","))
			_, err = tx.Exec(ctx, insertQuery, args...)
			if err != nil {
				return fmt.Errorf("insert technologies: %w", err)
			}
		}
	}

	return nil
}

func normalizeEnrichmentProfile(profile *EnrichmentProfileReq) error {
	var err error

	profile.Integrations, err = normalizeJSONArray(profile.Integrations, "integrations")
	if err != nil {
		return err
	}

	profile.TechnicalDecisions, err = normalizeJSONArray(profile.TechnicalDecisions, "technical_decisions")
	if err != nil {
		return err
	}

	profile.Challenges, err = normalizeJSONArray(profile.Challenges, "challenges")
	if err != nil {
		return err
	}

	profile.Results, err = normalizeJSONArray(profile.Results, "results")
	if err != nil {
		return err
	}

	profile.Timeline, err = normalizeJSONArray(profile.Timeline, "timeline")
	if err != nil {
		return err
	}

	profile.Metrics, err = normalizeJSONObject(profile.Metrics, "metrics")
	if err != nil {
		return err
	}

	return nil
}

func normalizeJSONArray(raw json.RawMessage, field string) (json.RawMessage, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return json.RawMessage("[]"), nil
	}

	var value []interface{}
	if err := json.Unmarshal([]byte(trimmed), &value); err != nil {
		return nil, fmt.Errorf("El campo %s debe ser un array JSON válido", field)
	}

	if err := validateStructuredJSONArray(value, field); err != nil {
		return nil, err
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("No se pudo normalizar el campo %s", field)
	}

	return normalized, nil
}

func normalizeJSONObject(raw json.RawMessage, field string) (json.RawMessage, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return json.RawMessage("{}"), nil
	}

	var value map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &value); err != nil {
		return nil, fmt.Errorf("El campo %s debe ser un objeto JSON válido", field)
	}

	if err := validateStructuredJSONObject(value, field); err != nil {
		return nil, err
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("No se pudo normalizar el campo %s", field)
	}

	return normalized, nil
}

func validateStructuredJSONArray(value []interface{}, field string) error {
	for _, item := range value {
		switch typed := item.(type) {
		case string, float64, bool:
			continue
		case map[string]interface{}:
			if err := validateStructuredJSONObject(typed, field); err != nil {
				return err
			}
		default:
			return fmt.Errorf("El campo %s solo admite arrays de valores primitivos u objetos planos", field)
		}
	}

	return nil
}

func validateStructuredJSONObject(value map[string]interface{}, field string) error {
	for _, item := range value {
		switch item.(type) {
		case string, float64, bool:
			continue
		default:
			return fmt.Errorf("El campo %s solo admite objetos planos con valores primitivos", field)
		}
	}

	return nil
}
