package handlers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/product"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type Product struct {
	service   product.Service
	responser response.API
}

func NewProduct(ps product.Service) *Product {
	return &Product{service: ps}
}

func (h *Product) Create(c echo.Context) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Brand       string `json:"brand"`
		Active      *bool  `json:"active"`
		// Legacy fields still accepted for backward compat
		ProductName string          `json:"product_name"`
		Price       float64         `json:"price"`
		Images      json.RawMessage `json:"images"`
		Features    json.RawMessage `json:"features"`
		Variants    []struct {
			SKU      string  `json:"sku"`
			Color    string  `json:"color"`
			Size     string  `json:"size"`
			Price    float64 `json:"price"`
			Stock    int     `json:"stock"`
			ImageURL string  `json:"image_url"`
		} `json:"variants"`
	}

	if err := c.Bind(&req); err != nil {
		return h.responser.BindFailed(c, "handlers-Product-Create-c.Bind()", err)
	}

	name := req.Name
	if name == "" {
		name = req.ProductName
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	m := &model.Product{
		ProductName: name,
		Description: req.Description,
		Images:      req.Images,
		Features:    req.Features,
	}

	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	// Set extended fields via the service
	m.SetStoreFields(name, req.Category, req.Brand, active)

	err := h.service.Create(m)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Create-h.service.Create()", err)
	}

	// Create variants if provided
	if len(req.Variants) > 0 {
		variants := make([]model.StoreProductVariant, 0, len(req.Variants))
		for _, v := range req.Variants {
			variants = append(variants, model.StoreProductVariant{
				ProductID: m.ID,
				SKU:       v.SKU,
				Color:     v.Color,
				Size:      v.Size,
				Price:     v.Price,
				Stock:     v.Stock,
				ImageURL:  v.ImageURL,
			})
		}
		err = h.service.CreateVariants(m.ID, variants)
		if err != nil {
			return h.responser.Error(c, "handlers-Product-Create-h.service.CreateVariants()", err)
		}
	}

	// Return the full StoreProduct
	productData, err := h.service.GetStoreByIDAdmin(m.ID)
	if err != nil {
		return c.JSON(h.responser.Created(m))
	}

	return c.JSON(response.ContractCreated(productData))
}

func (h *Product) Update(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del producto no es válido")
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Brand       string `json:"brand"`
		Active      *bool  `json:"active"`
		// Legacy fields still accepted for backward compat
		ProductName string          `json:"product_name"`
		Price       float64         `json:"price"`
		Images      json.RawMessage `json:"images"`
		Features    json.RawMessage `json:"features"`
		Variants    []struct {
			ID       string  `json:"id"`
			SKU      string  `json:"sku"`
			Color    string  `json:"color"`
			Size     string  `json:"size"`
			Price    float64 `json:"price"`
			Stock    int     `json:"stock"`
			ImageURL string  `json:"image_url"`
		} `json:"variants"`
	}

	if err = c.Bind(&req); err != nil {
		return h.responser.BindFailed(c, "handlers-Product-Update-c.Bind()", err)
	}

	name := req.Name
	if name == "" {
		name = req.ProductName
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	m := &model.Product{
		ID:          ID,
		ProductName: name,
		Description: req.Description,
		Images:      req.Images,
		Features:    req.Features,
	}

	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	m.SetStoreFields(name, req.Category, req.Brand, active)

	err = h.service.Update(m)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Update-h.service.Update()", err)
	}

	// Replace variants if provided
	if req.Variants != nil {
		variants := make([]model.StoreProductVariant, 0, len(req.Variants))
		for _, v := range req.Variants {
			variantID, _ := uuid.Parse(v.ID)
			variants = append(variants, model.StoreProductVariant{
				ID:        variantID,
				ProductID: ID,
				SKU:       v.SKU,
				Color:     v.Color,
				Size:      v.Size,
				Price:     v.Price,
				Stock:     v.Stock,
				ImageURL:  v.ImageURL,
			})
		}
		err = h.service.ReplaceVariants(ID, variants)
		if err != nil {
			return h.responser.Error(c, "handlers-Product-Update-h.service.ReplaceVariants()", err)
		}
	}

	// Return the full StoreProduct
	productData, err := h.service.GetStoreByIDAdmin(ID)
	if err != nil {
		return c.JSON(h.responser.Updated(m))
	}

	return c.JSON(response.ContractOK(productData))
}

func (h *Product) Delete(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Delete-uuid.Parse(c.Param('id'))", err)
	}

	err = h.service.Delete(ID)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Delete-h.service.Delete(ID)", err)
	}

	return c.JSON(h.responser.Deleted(nil))
}

func (h *Product) GetByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del producto no es válido")
	}

	productData, err := h.service.GetStoreByIDAdmin(ID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Producto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el producto")
	}

	return c.JSON(response.ContractOK(productData))
}

func (h *Product) GetStoreByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del producto no es válido")
	}

	productData, err := h.service.GetStoreByID(ID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Producto no encontrado")
		}
		if strings.Contains(strings.ToLower(err.Error()), "inactive") {
			return response.ContractError(404, "not_found", "Producto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el producto")
	}

	return c.JSON(response.ContractOK(productData))
}

/* Paginar el GETALL, en el query param del endpoint recibimos:
limit(cuantos registros quieren recibir) y page (en que páquina quieren mostrar)
offset: se genera limit*pag -limit */

func (h *Product) GetAll(c echo.Context) error {
	products, err := h.service.GetAll()
	if err != nil {
		return h.responser.Error(c, "handlers-Product-GetAll-h.service.GetAll()", err)
	}

	return c.JSON(h.responser.OK(products))
}

func (h *Product) GetStoreAll(c echo.Context) error {
	products, err := h.service.GetStoreAll()
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el catálogo")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": products}))
}

// UpdateStatus changes the active status of a product (admin only).
func (h *Product) UpdateStatus(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del producto no es válido")
	}

	var body struct {
		Active *bool `json:"active"`
	}
	if err = c.Bind(&body); err != nil {
		return response.ContractError(400, "validation_error", "Los datos enviados no son válidos")
	}

	if body.Active == nil {
		return response.ContractError(400, "validation_error", "El campo active es requerido")
	}

	productData, err := h.service.UpdateStatus(ID, *body.Active)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Producto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible actualizar el estado del producto")
	}

	return c.JSON(response.ContractOK(productData))
}

// GetAllStore returns all products including inactive ones (admin only).
func (h *Product) GetAllStore(c echo.Context) error {
	products, err := h.service.GetStoreAllAdmin()
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener los productos")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": products}))
}
