package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func TechnologyAdmin(e *echo.Echo, h *handlers.TechnologyHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/technologies", middlewares...)

	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.GetAll)
}

func TechnologyPublic(e *echo.Echo, h *handlers.TechnologyHandler) {
	g := e.Group("/api/v1/public/technologies")

	g.GET("", h.GetAll)
}
