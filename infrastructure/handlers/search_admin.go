package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

// SearchAdmin handles admin search management endpoints: readiness checks
// and search document re-embedding.
type SearchAdmin struct {
	projectReader project.ProjectReader
	searchRepo    search.SearchRepository
}

// NewSearchAdmin creates a new SearchAdmin handler.
func NewSearchAdmin(projectReader project.ProjectReader, searchRepo search.SearchRepository) *SearchAdmin {
	return &SearchAdmin{
		projectReader: projectReader,
		searchRepo:    searchRepo,
	}
}

// GetReadiness handles GET /api/v1/admin/projects/:id/readiness
// Returns the search readiness assessment for a project.
func (h *SearchAdmin) GetReadiness(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del proyecto no es válido")
	}

	p, err := h.projectReader.GetByID(c.Request().Context(), ID)
	if err != nil {
		return response.ContractError(404, "not_found", "Proyecto no encontrado")
	}

	techs, err := h.projectReader.GetTechnologiesByProjectID(c.Request().Context(), ID)
	if err != nil {
		// Technologies are optional for readiness — continue with empty slice
		techs = nil
	}

	readiness := model.ComputeReadiness(p, techs)
	return c.JSON(response.ContractOK(readiness))
}

// ReembedProject handles POST /api/v1/admin/projects/:id/reembed
// Refreshes the search document (tsvector) for a single project.
func (h *SearchAdmin) ReembedProject(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del proyecto no es válido")
	}

	// Verify project exists
	_, err = h.projectReader.GetByID(c.Request().Context(), ID)
	if err != nil {
		return response.ContractError(404, "not_found", "Proyecto no encontrado")
	}

	err = h.searchRepo.RefreshSearchDocument(c.Request().Context(), ID)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible actualizar el documento de búsqueda")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"message":    "Documento de búsqueda actualizado correctamente",
		"project_id": ID.String(),
	}))
}

// ReembedStale handles POST /api/v1/admin/projects/reembed-stale
// Refreshes search documents for all active projects.
func (h *SearchAdmin) ReembedStale(c echo.Context) error {
	err := h.searchRepo.RefreshAllDocuments(c.Request().Context())
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible actualizar los documentos de búsqueda")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"message": "Todos los documentos de búsqueda han sido actualizados",
	}))
}
