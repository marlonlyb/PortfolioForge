package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/product"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/model"
)

type Product struct {
	service      product.Service
	responser    response.API
	localization *localization.Service
}

func NewProduct(ps product.Service, localizationService *localization.Service) *Product {
	return &Product{service: ps, localization: localizationService}
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
		Media       []struct {
			ID           string `json:"id"`
			MediaType    string `json:"media_type"`
			URL          string `json:"url"`
			ThumbnailURL string `json:"thumbnail_url"`
			MediumURL    string `json:"medium_url"`
			FullURL      string `json:"full_url"`
			Caption      string `json:"caption"`
			AltText      string `json:"alt_text"`
			SortOrder    int    `json:"sort_order"`
			Featured     bool   `json:"featured"`
		} `json:"media"`
		Features json.RawMessage `json:"features"`
		Variants []struct {
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
	name = strings.TrimSpace(name)
	req.Category = strings.TrimSpace(req.Category)
	req.Brand = strings.TrimSpace(req.Brand)

	if contractErr := validateCreateProductRequest(name, req.Category, req.Brand); contractErr != nil {
		return contractErr
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	m := &model.Product{
		ProductName: name,
		Description: req.Description,
		Features:    req.Features,
	}

	projectMedia := buildProjectMediaPayload(m.ID, req.Media, req.Images)
	legacyImages, err := marshalProjectLegacyImages(projectMedia, req.Images)
	if err != nil {
		return response.ContractError(400, "validation_error", "Debes enviar una lista válida de imágenes")
	}
	m.Images = legacyImages

	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	// Set extended fields via the service
	m.SetStoreFields(name, req.Category, req.Brand, active)

	err = h.service.Create(m)
	if err != nil {
		return mapCreateProductError(err)
	}

	projectMedia = assignProjectMediaProjectID(m.ID, projectMedia)
	if err = h.service.ReplaceMedia(m.ID, projectMedia); err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible guardar la galería del proyecto")
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
			return mapCreateProductError(err)
		}
	}

	if syncErr := h.syncSpanishProjectFields(c.Request().Context(), uuid.Nil, m.ID); syncErr != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible generar las traducciones automáticas")
	}

	// Return the full StoreProduct
	productData, err := h.service.GetStoreByIDAdmin(m.ID)
	if err != nil {
		return c.JSON(h.responser.Created(m))
	}

	return c.JSON(response.ContractCreated(productData))
}

func validateCreateProductRequest(name, category, brand string) *model.ContractError {
	if name == "" {
		return response.ContractError(400, "validation_error", "El nombre del proyecto es requerido", model.APIErrorDetail{Field: "name", Issue: "required"})
	}

	if len([]rune(name)) > 128 {
		return response.ContractError(400, "validation_error", "El nombre del proyecto no puede exceder 128 caracteres", model.APIErrorDetail{Field: "name", Issue: "max_length:128"})
	}

	if len([]rune(category)) > 80 {
		return response.ContractError(400, "validation_error", "La categoría no puede exceder 80 caracteres", model.APIErrorDetail{Field: "category", Issue: "max_length:80"})
	}

	if len([]rune(brand)) > 80 {
		return response.ContractError(400, "validation_error", "La marca no puede exceder 80 caracteres", model.APIErrorDetail{Field: "brand", Issue: "max_length:80"})
	}

	return nil
}

func mapCreateProductError(err error) error {
	if err == nil {
		return nil
	}

	lowerErr := strings.ToLower(err.Error())
	if strings.Contains(lowerErr, "product name is empty") {
		return response.ContractError(400, "validation_error", "El nombre del proyecto es requerido", model.APIErrorDetail{Field: "name", Issue: "required"})
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "22001":
			field, message := productLengthViolationMessage(pgErr.ColumnName)
			if field != "" {
				return response.ContractError(400, "validation_error", message, model.APIErrorDetail{Field: field, Issue: "max_length"})
			}
			return response.ContractError(400, "validation_error", "Uno de los textos enviados excede la longitud permitida")
		case "23502":
			field, message := productRequiredFieldMessage(pgErr.ColumnName)
			if field != "" {
				return response.ContractError(400, "validation_error", message, model.APIErrorDetail{Field: field, Issue: "required"})
			}
			return response.ContractError(400, "validation_error", "Faltan datos requeridos para crear el proyecto")
		case "23505":
			field, message := productUniqueViolationMessage(pgErr.ConstraintName)
			if field != "" {
				return response.ContractError(409, "validation_error", message, model.APIErrorDetail{Field: field, Issue: "unique"})
			}
			return response.ContractError(409, "validation_error", "Ya existe un registro con uno de los valores enviados")
		case "23503":
			return response.ContractError(400, "validation_error", "Uno de los datos relacionados no existe o ya no está disponible")
		case "23514":
			return response.ContractError(400, "validation_error", "Uno de los valores enviados no cumple las restricciones permitidas")
		case "22P02":
			return response.ContractError(400, "validation_error", "Uno de los valores enviados tiene un formato inválido")
		}
	}

	return response.ContractError(500, "unexpected_error", "No fue posible crear el proyecto")
}

func productLengthViolationMessage(column string) (field string, message string) {
	switch column {
	case "product_name", "name":
		return "name", "El nombre del proyecto no puede exceder 128 caracteres"
	case "category":
		return "category", "La categoría no puede exceder 80 caracteres"
	case "brand":
		return "brand", "La marca no puede exceder 80 caracteres"
	case "slug":
		return "name", "El nombre del proyecto genera un slug demasiado largo"
	case "sku":
		return "variants", "El SKU de una variante no puede exceder 120 caracteres"
	case "color":
		return "variants", "El color de una variante no puede exceder 60 caracteres"
	case "size":
		return "variants", "La talla de una variante no puede exceder 30 caracteres"
	default:
		return "", ""
	}
}

func productRequiredFieldMessage(column string) (field string, message string) {
	switch column {
	case "product_name", "name":
		return "name", "El nombre del proyecto es requerido"
	case "description":
		return "description", "La descripción del proyecto es requerida"
	case "images":
		return "images", "Debes enviar una lista válida de imágenes"
	case "features":
		return "features", "Debes enviar una lista válida de features"
	case "category":
		return "category", "La categoría del proyecto es requerida"
	default:
		return "", ""
	}
}

func productUniqueViolationMessage(constraint string) (field string, message string) {
	switch constraint {
	case "ix_products_slug", "products_slug_uk":
		return "name", "Ya existe un proyecto con un slug equivalente. Usa un nombre distinto"
	case "product_variants_sku_uk":
		return "variants", "El SKU de una de las variantes ya existe"
	case "product_variants_product_color_size_uk":
		return "variants", "Ya existe una variante con la misma combinación de color y talla"
	default:
		return "", ""
	}
}

func (h *Product) Update(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del producto no es válido")
	}

	previousProjectID := ID
	previousProject, _ := h.service.GetStoreByIDAdmin(previousProjectID)
	previousLocalized := storeProductToLocalizedProject(previousProject)

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
		Media       []struct {
			ID           string `json:"id"`
			MediaType    string `json:"media_type"`
			URL          string `json:"url"`
			ThumbnailURL string `json:"thumbnail_url"`
			MediumURL    string `json:"medium_url"`
			FullURL      string `json:"full_url"`
			Caption      string `json:"caption"`
			AltText      string `json:"alt_text"`
			SortOrder    int    `json:"sort_order"`
			Featured     bool   `json:"featured"`
		} `json:"media"`
		Features json.RawMessage `json:"features"`
		Variants []struct {
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
		Features:    req.Features,
	}

	projectMedia := buildProjectMediaPayload(ID, req.Media, req.Images)
	legacyImages, marshalErr := marshalProjectLegacyImages(projectMedia, req.Images)
	if marshalErr != nil {
		return response.ContractError(400, "validation_error", "Debes enviar una lista válida de imágenes")
	}
	m.Images = legacyImages

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

	if err = h.service.ReplaceMedia(ID, projectMedia); err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible guardar la galería del proyecto")
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

	currentProject, syncErr := h.service.GetStoreByIDAdmin(ID)
	if syncErr == nil {
		if err := h.localization.SyncFromSpanish(c.Request().Context(), ID, localization.BuildProjectFieldMap(previousLocalized), localization.BuildProjectFieldMap(storeProductToLocalizedProject(currentProject))); err != nil {
			return response.ContractError(500, "unexpected_error", "No fue posible actualizar las traducciones automáticas")
		}
	}

	// Return the full StoreProduct
	productData, err := h.service.GetStoreByIDAdmin(ID)
	if err != nil {
		return c.JSON(h.responser.Updated(m))
	}

	return c.JSON(response.ContractOK(productData))
}

func (h *Product) syncSpanishProjectFields(ctx context.Context, previousProjectID uuid.UUID, currentProjectID uuid.UUID) error {
	if h.localization == nil {
		return nil
	}

	previous := model.Project{}
	if previousProjectID != uuid.Nil {
		if previousStore, err := h.service.GetStoreByIDAdmin(previousProjectID); err == nil {
			previous = storeProductToLocalizedProject(previousStore)
		}
	}

	currentStore, err := h.service.GetStoreByIDAdmin(currentProjectID)
	if err != nil {
		return err
	}

	return h.localization.SyncFromSpanish(ctx, currentProjectID, localization.BuildProjectFieldMap(previous), localization.BuildProjectFieldMap(storeProductToLocalizedProject(currentStore)))
}

func storeProductToLocalizedProject(storeProduct model.StoreProduct) model.Project {
	return model.Project{
		ID:          storeProduct.ID,
		Name:        storeProduct.Name,
		Description: storeProduct.Description,
		Category:    storeProduct.Category,
	}
}

func buildProjectMediaPayload(projectID uuid.UUID, rawMedia []struct {
	ID           string `json:"id"`
	MediaType    string `json:"media_type"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	MediumURL    string `json:"medium_url"`
	FullURL      string `json:"full_url"`
	Caption      string `json:"caption"`
	AltText      string `json:"alt_text"`
	SortOrder    int    `json:"sort_order"`
	Featured     bool   `json:"featured"`
}, fallbackImages json.RawMessage) []model.ProjectMedia {
	if len(rawMedia) == 0 {
		var legacyImages []string
		if len(fallbackImages) > 0 {
			_ = json.Unmarshal(fallbackImages, &legacyImages)
		}
		return model.BuildLegacyProjectMedia(projectID, legacyImages)
	}

	media := make([]model.ProjectMedia, 0, len(rawMedia))
	for _, item := range rawMedia {
		if strings.TrimSpace(item.ThumbnailURL) == "" && strings.TrimSpace(item.MediumURL) == "" && strings.TrimSpace(item.FullURL) == "" && strings.TrimSpace(item.URL) == "" {
			continue
		}

		mediaID, _ := uuid.Parse(item.ID)
		media = append(media, model.ProjectMedia{
			ID:           mediaID,
			ProjectID:    projectID,
			MediaType:    item.MediaType,
			URL:          item.URL,
			ThumbnailURL: item.ThumbnailURL,
			MediumURL:    item.MediumURL,
			FullURL:      item.FullURL,
			Caption:      item.Caption,
			AltText:      item.AltText,
			SortOrder:    item.SortOrder,
			Featured:     item.Featured,
		})
	}

	if len(media) == 0 {
		var legacyImages []string
		if len(fallbackImages) > 0 {
			_ = json.Unmarshal(fallbackImages, &legacyImages)
		}
		return model.BuildLegacyProjectMedia(projectID, legacyImages)
	}

	return media
}

func marshalProjectLegacyImages(media []model.ProjectMedia, fallback json.RawMessage) (json.RawMessage, error) {
	var fallbackImages []string
	if len(fallback) > 0 {
		if err := json.Unmarshal(fallback, &fallbackImages); err != nil {
			return nil, err
		}
	}

	images := model.BuildProjectImageList(media, fallbackImages)
	if len(images) == 0 {
		images = []string{}
	}

	encoded, err := json.Marshal(images)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func assignProjectMediaProjectID(projectID uuid.UUID, media []model.ProjectMedia) []model.ProjectMedia {
	for index := range media {
		media[index].ProjectID = projectID
	}

	return media
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
