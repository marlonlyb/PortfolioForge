package handlers

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type ProjectAssistantHandler struct {
	service services.ProjectAssistant
}

func NewProjectAssistantHandler(service services.ProjectAssistant) *ProjectAssistantHandler {
	return &ProjectAssistantHandler{service: service}
}

func (h *ProjectAssistantHandler) CreateMessage(c echo.Context) error {
	slug := strings.TrimSpace(c.Param("slug"))
	var request model.ProjectAssistantRequest
	if err := c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Datos inválidos para el assistant")
	}

	answer, err := h.service.Answer(c.Request().Context(), slug, request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAssistantInvalidQuestion):
			return response.ContractError(400, "validation_error", "La pregunta debe tener entre 2 y 2000 caracteres")
		case errors.Is(err, services.ErrAssistantProjectNotFound):
			return response.ContractError(404, "not_found", "Proyecto no encontrado")
		case errors.Is(err, services.ErrAssistantUnavailable):
			return response.ContractError(409, "assistant_unavailable", "El assistant no está disponible para este proyecto")
		default:
			return response.ContractError(502, "assistant_upstream_error", "No fue posible generar la respuesta del assistant")
		}
	}

	return c.JSON(response.ContractOK(answer))
}
