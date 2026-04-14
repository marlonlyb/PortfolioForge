package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"

	project "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/model"
)

// ProjectCatalog is the canonical admin handler for portfolio projects.
// Legacy product routes can still reuse it during the compatibility window.
type ProjectCatalog struct {
	service      project.AdminCatalogService
	responser    response.API
	localization *localization.Service
}

func NewProjectCatalog(ps project.AdminCatalogService, localizationService *localization.Service) *ProjectCatalog {
	return &ProjectCatalog{service: ps, localization: localizationService}
}

func (h *ProjectCatalog) Create(c echo.Context) error {
	var req model.AdminProjectWrite

	if err := c.Bind(&req); err != nil {
		return h.responser.BindFailed(c, "handlers-Product-Create-c.Bind()", err)
	}

	req.Normalize()

	if contractErr := validateCreateProductRequest(req.Name, req.Category, req.ClientName); contractErr != nil {
		return contractErr
	}

	projectMedia := buildProjectMediaPayload(uuid.Nil, req.Media, req.Images)
	legacyImages, err := marshalProjectLegacyImages(projectMedia, req.Images)
	if err != nil {
		return response.ContractError(400, "validation_error", "Debes enviar una lista válida de imágenes")
	}
	req.Images = legacyImages

	if len(req.Images) == 0 {
		req.Images = []byte(`[]`)
	}
	if len(req.Features) == 0 {
		req.Features = []byte(`[]`)
	}

	err = h.service.Create(&req)
	if err != nil {
		return mapCreateProductError(err)
	}

	projectMedia = assignProjectMediaProjectID(req.ID, projectMedia)
	if err = h.service.ReplaceMedia(req.ID, projectMedia); err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible guardar la galería del proyecto")
	}

	if len(req.Variants) > 0 {
		err = h.service.CreateVariants(req.ID, req.Variants)
		if err != nil {
			return mapCreateProductError(err)
		}
	}

	if syncErr := h.syncSpanishProjectFields(c.Request().Context(), uuid.Nil, req.ID); syncErr != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible generar las traducciones automáticas")
	}

	projectData, err := h.service.GetAdminByID(req.ID)
	if err != nil {
		return c.JSON(h.responser.Created(req))
	}

	return c.JSON(response.ContractCreated(projectData))
}

func validateCreateProductRequest(name, category, clientName string) *model.ContractError {
	if name == "" {
		return response.ContractError(400, "validation_error", "El nombre del proyecto es requerido", model.APIErrorDetail{Field: "name", Issue: "required"})
	}

	if len([]rune(name)) > 128 {
		return response.ContractError(400, "validation_error", "El nombre del proyecto no puede exceder 128 caracteres", model.APIErrorDetail{Field: "name", Issue: "max_length:128"})
	}

	if len([]rune(category)) > 80 {
		return response.ContractError(400, "validation_error", "La categoría no puede exceder 80 caracteres", model.APIErrorDetail{Field: "category", Issue: "max_length:80"})
	}

	if len([]rune(clientName)) > 80 {
		return response.ContractError(400, "validation_error", "El cliente/contexto no puede exceder 80 caracteres", model.APIErrorDetail{Field: "client_name", Issue: "max_length:80"})
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
		return "client_name", "El cliente/contexto no puede exceder 80 caracteres"
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

func (h *ProjectCatalog) Update(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del proyecto no es válido")
	}

	previousProjectID := ID
	previousProject, _ := h.service.GetAdminByID(previousProjectID)
	previousLocalized := adminProjectToLocalizedProject(previousProject)

	var req model.AdminProjectWrite

	if err = c.Bind(&req); err != nil {
		return h.responser.BindFailed(c, "handlers-Product-Update-c.Bind()", err)
	}

	req.ID = ID
	req.Normalize()

	projectMedia := buildProjectMediaPayload(ID, req.Media, req.Images)
	legacyImages, marshalErr := marshalProjectLegacyImages(projectMedia, req.Images)
	if marshalErr != nil {
		return response.ContractError(400, "validation_error", "Debes enviar una lista válida de imágenes")
	}
	req.Images = legacyImages

	if len(req.Images) == 0 {
		req.Images = []byte(`[]`)
	}
	if len(req.Features) == 0 {
		req.Features = []byte(`[]`)
	}

	err = h.service.Update(&req)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Update-h.service.Update()", err)
	}

	if err = h.service.ReplaceMedia(ID, projectMedia); err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible guardar la galería del proyecto")
	}

	if req.Variants != nil {
		err = h.service.ReplaceVariants(ID, req.Variants)
		if err != nil {
			return h.responser.Error(c, "handlers-Product-Update-h.service.ReplaceVariants()", err)
		}
	}

	currentProject, syncErr := h.service.GetAdminByID(ID)
	if syncErr == nil {
		if err := h.localization.SyncFromSpanish(c.Request().Context(), ID, localization.BuildProjectFieldMap(previousLocalized), localization.BuildProjectFieldMap(adminProjectToLocalizedProject(currentProject))); err != nil {
			return response.ContractError(500, "unexpected_error", "No fue posible actualizar las traducciones automáticas")
		}
	}

	projectData, err := h.service.GetAdminByID(ID)
	if err != nil {
		return c.JSON(h.responser.Updated(req))
	}

	return c.JSON(response.ContractOK(projectData))
}

func (h *ProjectCatalog) syncSpanishProjectFields(ctx context.Context, previousProjectID uuid.UUID, currentProjectID uuid.UUID) error {
	if h.localization == nil {
		return nil
	}

	previous := model.Project{}
	if previousProjectID != uuid.Nil {
		if previousAdmin, err := h.service.GetAdminByID(previousProjectID); err == nil {
			previous = adminProjectToLocalizedProject(previousAdmin)
		}
	}

	currentProject, err := h.service.GetAdminByID(currentProjectID)
	if err != nil {
		return err
	}

	return h.localization.SyncFromSpanish(ctx, currentProjectID, localization.BuildProjectFieldMap(previous), localization.BuildProjectFieldMap(adminProjectToLocalizedProject(currentProject)))
}

func adminProjectToLocalizedProject(project model.AdminProject) model.Project {
	return project.ToProject()
}

func buildProjectMediaPayload(projectID uuid.UUID, rawMedia []model.AdminProjectMediaInput, fallbackImages json.RawMessage) []model.ProjectMedia {
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

func (h *ProjectCatalog) Delete(c echo.Context) error {
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

func (h *ProjectCatalog) GetByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del proyecto no es válido")
	}

	projectData, err := h.service.GetAdminByID(ID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || strings.Contains(err.Error(), "no rows") {
			return response.ContractError(404, "not_found", "Proyecto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible obtener el proyecto")
	}

	return c.JSON(response.ContractOK(projectData))
}

// UpdateStatus changes the active status of a product (admin only).
func (h *ProjectCatalog) UpdateStatus(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador del proyecto no es válido")
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
			return response.ContractError(404, "not_found", "Proyecto no encontrado")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible actualizar el estado del proyecto")
	}

	return c.JSON(response.ContractOK(productData))
}

// GetAllStore returns all products including inactive ones (admin only).
func (h *ProjectCatalog) GetAllStore(c echo.Context) error {
	projects, err := h.service.GetAdminAll()
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener los proyectos")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": projects}))
}
