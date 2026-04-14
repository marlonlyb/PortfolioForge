package handlers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	project "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

// ProductPublicCompat isolates the legacy `/public/products` route family.
// Canonical portfolio reads must use `/public/projects`.
type ProductPublicCompat struct {
	service project.PublicProductCompatService
}

func NewProductPublicCompat(service project.PublicProductCompatService) *ProductPublicCompat {
	return &ProductPublicCompat{service: service}
}

func (h *ProductPublicCompat) GetStoreByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del proyecto no es válido")
	}

	productData, err := h.service.GetStoreByID(ID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Proyecto no encontrado")
		}
		if strings.Contains(strings.ToLower(err.Error()), "inactive") {
			return response.ContractError(404, "not_found", "Proyecto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el proyecto")
	}

	return c.JSON(response.ContractOK(productData))
}

func (h *ProductPublicCompat) GetStoreAll(c echo.Context) error {
	products, err := h.service.GetStoreAll()
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el catálogo")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": products}))
}
