package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func registerAdminCatalogRoutes(g *echo.Group, h handlers.ProjectAdminCatalogHandler) {
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.GetAllStore)
	g.GET("/:id", h.GetByID)
	g.PATCH("/:id/status", h.UpdateStatus)
}

// ProductAdmin keeps the legacy admin-products contract alive during transition.
func ProductAdmin(e *echo.Echo, h handlers.ProjectAdminCatalogHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/products", middlewares...)
	registerAdminCatalogRoutes(g, h)
}

// ProjectAdminCatalog exposes the canonical admin-projects contract while storage
// still flows through the legacy `products` persistence layer.
func ProjectAdminCatalog(e *echo.Echo, h handlers.ProjectAdminCatalogHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/projects", middlewares...)
	registerAdminCatalogRoutes(g, h)
}

// ProductPublicCompat keeps `/api/v1/public/products` alive as an explicitly
// isolated compatibility surface. Canonical portfolio reads live in
// `/api/v1/public/projects`.
func ProductPublicCompat(e *echo.Echo, h handlers.ProductPublicCompatHandler) {
	g := e.Group("/api/v1/public/products")

	g.GET("", h.GetStoreAll)
	g.GET("/:id", h.GetStoreByID)
}
