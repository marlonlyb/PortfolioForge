package handlers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

// TechnologyRepo defines the contract we need from postgres
type TechnologyRepo interface {
	Create(m *model.Technology) error
	Update(m *model.Technology) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (model.Technology, error)
	GetAll() ([]model.Technology, error)
}

type TechnologyHandler struct {
	repo TechnologyRepo
}

func NewTechnologyHandler(repo TechnologyRepo) *TechnologyHandler {
	return &TechnologyHandler{repo: repo}
}

func (h *TechnologyHandler) Create(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Category string `json:"category"`
		Icon     string `json:"icon"`
		Color    string `json:"color"`
	}

	if err := c.Bind(&req); err != nil {
		return response.ContractError(400, "validation_error", "Datos inválidos")
	}

	m := &model.Technology{
		Name:     req.Name,
		Slug:     req.Slug,
		Category: req.Category,
		Icon:     req.Icon,
		Color:    req.Color,
	}

	if err := h.repo.Create(m); err != nil {
		return response.ContractError(500, "unexpected_error", "No se pudo crear la tecnología: "+err.Error())
	}

	return c.JSON(response.ContractCreated(m))
}

func (h *TechnologyHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "ID de tecnología inválido")
	}

	var req struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Category string `json:"category"`
		Icon     string `json:"icon"`
		Color    string `json:"color"`
	}

	if err := c.Bind(&req); err != nil {
		return response.ContractError(400, "validation_error", "Datos inválidos")
	}

	m := &model.Technology{
		ID:       id,
		Name:     req.Name,
		Slug:     req.Slug,
		Category: req.Category,
		Icon:     req.Icon,
		Color:    req.Color,
	}

	if err := h.repo.Update(m); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Tecnología no encontrada")
		}
		return response.ContractError(500, "unexpected_error", "No se pudo actualizar la tecnología: "+err.Error())
	}

	return c.JSON(response.ContractOK(m))
}

func (h *TechnologyHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "ID de tecnología inválido")
	}

	if err := h.repo.Delete(id); err != nil {
		return response.ContractError(500, "unexpected_error", "No se pudo eliminar la tecnología: "+err.Error())
	}

	return c.JSON(response.ContractOK(map[string]string{"message": "deleted"}))
}

func (h *TechnologyHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "ID de tecnología inválido")
	}

	tech, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(strings.ToLower(err.Error()), "no rows") {
			return response.ContractError(404, "not_found", "Tecnología no encontrada")
		}

		return response.ContractError(500, "unexpected_error", "No se pudo obtener la tecnología")
	}

	return c.JSON(response.ContractOK(tech))
}

func (h *TechnologyHandler) GetAll(c echo.Context) error {
	techs, err := h.repo.GetAll()
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No se pudieron obtener las tecnologías")
	}

	if techs == nil {
		techs = []model.Technology{}
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": techs}))
}
