package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func SiteSettingsPublic(e *echo.Echo, h *handlers.SiteSettingsHandler) {
	g := e.Group("/api/v1/public/site-settings")

	g.GET("", h.GetPublic)
}

func SiteSettingsAdmin(e *echo.Echo, h *handlers.SiteSettingsHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/site-settings", middlewares...)

	g.GET("", h.GetAdmin)
	g.PUT("", h.SaveAdmin)
}
