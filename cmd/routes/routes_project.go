package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

// ProjectPublic registers public project read routes.
func ProjectPublic(e *echo.Echo, h handlers.ProjectPublicHandler) {
	g := e.Group("/api/v1/public/projects")

	g.GET("", h.ListPublished)
	g.GET("/:slug", h.GetBySlug)
}

func ProjectPrivate(e *echo.Echo, assistant handlers.ProjectAssistantHandlerContract, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/private/projects", middlewares...)

	g.POST("/:slug/assistant/messages", assistant.CreateMessage)
}

// ProjectAdmin defines admin routes for project enrichment.
func ProjectAdmin(e *echo.Echo, h *handlers.ProjectAdminHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/projects", middlewares...)

	g.PUT("/:id/enrichment", h.UpdateProjectEnrichment)
	g.GET("/:id/localizations", h.GetProjectTranslations)
	g.PUT("/:id/localizations/:locale", h.SaveProjectTranslations)
}
