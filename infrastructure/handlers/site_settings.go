package handlers

import (
	"context"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type SiteSettingsRepository interface {
	Get(ctx context.Context) (model.SiteSettings, error)
	Save(ctx context.Context, settings model.SiteSettings) (model.SiteSettings, error)
}

type SiteSettingsHandler struct {
	repo SiteSettingsRepository
}

type SaveSiteSettingsRequest struct {
	PublicHeroLogoURL string `json:"public_hero_logo_url"`
	PublicHeroLogoAlt string `json:"public_hero_logo_alt"`
}

func NewSiteSettingsHandler(repo SiteSettingsRepository) *SiteSettingsHandler {
	return &SiteSettingsHandler{repo: repo}
}

func (h *SiteSettingsHandler) GetPublic(c echo.Context) error {
	settings, err := h.repo.Get(c.Request().Context())
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener la configuración pública")
	}

	return c.JSON(response.ContractOK(settings))
}

func (h *SiteSettingsHandler) GetAdmin(c echo.Context) error {
	settings, err := h.repo.Get(c.Request().Context())
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener la configuración del sitio")
	}

	return c.JSON(response.ContractOK(settings))
}

func (h *SiteSettingsHandler) SaveAdmin(c echo.Context) error {
	var req SaveSiteSettingsRequest
	if err := c.Bind(&req); err != nil {
		return response.ContractError(400, "validation_error", "Datos de configuración inválidos")
	}

	urlValue := strings.TrimSpace(req.PublicHeroLogoURL)
	if urlValue != "" {
		parsed, err := url.ParseRequestURI(urlValue)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return response.ContractError(400, "validation_error", "La URL del logo público no es válida")
		}
	}

	settings, err := h.repo.Save(c.Request().Context(), model.SiteSettings{
		PublicHeroLogoURL: urlValue,
		PublicHeroLogoAlt: strings.TrimSpace(req.PublicHeroLogoAlt),
	})
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible guardar la configuración del sitio")
	}

	return c.JSON(response.ContractOK(settings))
}
