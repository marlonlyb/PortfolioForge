package handlers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
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
	currentUser, ok := c.Get("currentUser").(model.User)
	if !ok || currentUser.ID == uuid.Nil {
		return response.ContractError(401, "authentication_required", "You must sign in to continue")
	}
	if !currentUser.CanUseProjectAssistant {
		return response.ContractError(403, "assistant_ineligible", "Complete the required sign-in and profile steps to use the assistant")
	}

	slug := strings.TrimSpace(c.Param("slug"))
	var request model.ProjectAssistantRequest
	if err := c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Invalid assistant payload")
	}

	answer, err := h.service.Answer(c.Request().Context(), slug, request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAssistantInvalidQuestion):
			return response.ContractError(400, "validation_error", "The question must contain between 2 and 2000 characters")
		case errors.Is(err, services.ErrAssistantProjectNotFound):
			return response.ContractError(404, "not_found", "Project not found")
		case errors.Is(err, services.ErrAssistantUnavailable):
			return response.ContractError(409, "assistant_unavailable", "The assistant is not available for this project")
		default:
			return response.ContractError(502, "assistant_upstream_error", "Unable to generate the assistant response")
		}
	}

	return c.JSON(response.ContractOK(answer))
}
