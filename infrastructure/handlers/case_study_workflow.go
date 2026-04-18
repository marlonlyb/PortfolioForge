package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type CaseStudyWorkflowService interface {
	GetAvailability(ctx context.Context) (model.CaseStudyWorkflowAvailability, error)
	StartRun(ctx context.Context, req model.StartCaseStudyWorkflowRunRequest) (model.CaseStudyWorkflowRun, error)
	GetRun(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error)
	ListLogs(ctx context.Context, runID uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error)
	ConfirmStep(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error)
	StartStep(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error)
	RetryStep(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error)
	Resume(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error)
}

type CaseStudyWorkflowHandler struct {
	service CaseStudyWorkflowService
}

func NewCaseStudyWorkflowHandler(service CaseStudyWorkflowService) *CaseStudyWorkflowHandler {
	return &CaseStudyWorkflowHandler{service: service}
}

func (h *CaseStudyWorkflowHandler) GetAvailability(c echo.Context) error {
	availability, err := h.service.GetAvailability(c.Request().Context())
	if err != nil {
		return mapCaseStudyWorkflowError(err)
	}
	return c.JSON(response.ContractOK(availability))
}

func (h *CaseStudyWorkflowHandler) StartRun(c echo.Context) error {
	var req model.StartCaseStudyWorkflowRunRequest
	if err := c.Bind(&req); err != nil {
		return response.ContractError(400, "validation_error", "Datos del workflow inválidos")
	}
	run, err := h.service.StartRun(c.Request().Context(), req)
	if err != nil {
		return mapCaseStudyWorkflowError(err)
	}
	return c.JSON(response.ContractCreated(run))
}

func (h *CaseStudyWorkflowHandler) GetRun(c echo.Context) error {
	runID, err := parseWorkflowRunID(c)
	if err != nil {
		return err
	}
	run, err := h.service.GetRun(c.Request().Context(), runID)
	if err != nil {
		return mapCaseStudyWorkflowError(err)
	}
	return c.JSON(response.ContractOK(run))
}

func (h *CaseStudyWorkflowHandler) GetLogs(c echo.Context) error {
	runID, err := parseWorkflowRunID(c)
	if err != nil {
		return err
	}
	logs, err := h.service.ListLogs(c.Request().Context(), runID)
	if err != nil {
		return mapCaseStudyWorkflowError(err)
	}
	return c.JSON(response.ContractOK(map[string]any{"items": logs}))
}

func (h *CaseStudyWorkflowHandler) ConfirmStep(c echo.Context) error {
	return h.handleStepMutation(c, func(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error) {
		return h.service.ConfirmStep(ctx, runID, stepName)
	})
}

func (h *CaseStudyWorkflowHandler) StartStep(c echo.Context) error {
	return h.handleStepMutation(c, func(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error) {
		return h.service.StartStep(ctx, runID, stepName)
	})
}

func (h *CaseStudyWorkflowHandler) RetryStep(c echo.Context) error {
	return h.handleStepMutation(c, func(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error) {
		return h.service.RetryStep(ctx, runID, stepName)
	})
}

func (h *CaseStudyWorkflowHandler) Resume(c echo.Context) error {
	runID, err := parseWorkflowRunID(c)
	if err != nil {
		return err
	}
	run, err := h.service.Resume(c.Request().Context(), runID)
	if err != nil {
		return mapCaseStudyWorkflowError(err)
	}
	return c.JSON(response.ContractOK(run))
}

func (h *CaseStudyWorkflowHandler) handleStepMutation(c echo.Context, fn func(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error)) error {
	runID, err := parseWorkflowRunID(c)
	if err != nil {
		return err
	}
	stepName := strings.TrimSpace(c.Param("step"))
	if !model.IsCaseStudyWorkflowStep(stepName) {
		return response.ContractError(400, "validation_error", fmt.Sprintf("Step inválido: %s", stepName))
	}
	run, err := fn(c.Request().Context(), runID, stepName)
	if err != nil {
		return mapCaseStudyWorkflowError(err)
	}
	return c.JSON(response.ContractOK(run))
}

func parseWorkflowRunID(c echo.Context) (uuid.UUID, error) {
	runID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, response.ContractError(400, "validation_error", "ID de workflow inválido")
	}
	return runID, nil
}

func mapCaseStudyWorkflowError(err error) error {
	if err == nil {
		return nil
	}
	var unavailableErr *model.CaseStudyWorkflowUnavailableError
	if errors.As(err, &unavailableErr) {
		return response.ContractError(503, "workflow_unavailable", unavailableErr.Error())
	}
	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "obligatorio") || strings.Contains(lower, "inválido") || strings.Contains(lower, "allowlist") || strings.Contains(lower, "fuera de las raíces") || strings.Contains(lower, "todavía no está listo") || strings.Contains(lower, "confirmación") || strings.Contains(lower, "no es reintentable") {
		return response.ContractError(400, "validation_error", err.Error())
	}
	if strings.Contains(lower, "no existe") || strings.Contains(lower, "no encontrado") {
		return response.ContractError(404, "not_found", err.Error())
	}
	return response.ContractError(500, "unexpected_error", err.Error())
}
