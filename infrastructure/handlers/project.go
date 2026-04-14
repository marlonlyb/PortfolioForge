package handlers

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
)

// ProjectPublic handles public read-only project API requests.
type ProjectPublic struct {
	service             services.Project
	localizationService *localization.Service
}

// NewProjectPublic creates a new ProjectPublic handler.
func NewProjectPublic(service services.Project, localizationService *localization.Service) *ProjectPublic {
	return &ProjectPublic{service: service, localizationService: localizationService}
}

// GetBySlug handles GET /api/v1/public/projects/:slug
func (h *ProjectPublic) GetBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return response.ContractError(400, "validation_error", "El slug del proyecto es requerido")
	}

	project, err := h.service.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Proyecto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el proyecto")
	}

	project, err = h.localizationService.LocalizeProject(c.Request().Context(), project, c.QueryParam("lang"))
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible localizar el proyecto")
	}

	return c.JSON(response.ContractOK(project))
}

// ListPublished handles GET /api/v1/public/projects
func (h *ProjectPublic) ListPublished(c echo.Context) error {
	projects, err := h.service.ListPublished(c.Request().Context())
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener los proyectos")
	}

	projects, err = h.localizationService.LocalizeProjects(c.Request().Context(), projects, c.QueryParam("lang"))
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible localizar los proyectos")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": projects}))
}
