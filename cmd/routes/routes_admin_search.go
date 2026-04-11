package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

// SearchAdmin registers admin search management routes behind auth middleware.
// Endpoints:
//
//	GET  /api/v1/admin/projects/:id/readiness — project search readiness
//	POST /api/v1/admin/projects/:id/reembed   — refresh single project search doc
//	POST /api/v1/admin/projects/reembed-stale  — refresh all search docs
func SearchAdmin(e *echo.Echo, h handlers.SearchAdminHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/projects", middlewares...)

	// NOTE: reembed-stale must be registered before /:id routes to avoid
	// Echo matching "reembed-stale" as an :id parameter.
	g.POST("/reembed-stale", h.ReembedStale)
	g.GET("/:id/readiness", h.GetReadiness)
	g.POST("/:id/reembed", h.ReembedProject)
}
