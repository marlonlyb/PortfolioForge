package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

// Para autorizaciones:
// func productAdminRoutes(e *echo.Echo, h productPorts.ProductHandlers, middlewares ...echo.MiddlewareFunc) {
func ProductAdmin(e *echo.Echo, h handlers.ProductHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/products", middlewares...)

	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.GetAllStore)
	g.GET("/:id", h.GetByID)
	g.PATCH("/:id/status", h.UpdateStatus)
}

func ProductPublic(e *echo.Echo, h handlers.ProductHandler) {
	g := e.Group("/api/v1/public/products")

	g.GET("", h.GetStoreAll)
	g.GET("/:id", h.GetStoreByID)
}
