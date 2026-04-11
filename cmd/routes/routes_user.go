package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func UserAdmin(e *echo.Echo, h handlers.UserHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/users", middlewares...)

	g.GET("", h.GetAll)
}

func UserPublic(e *echo.Echo, h handlers.UserHandler) {
	g := e.Group("/api/v1/public/users")

	g.POST("", h.Create)
	e.POST("/api/v1/public/register", h.Register)
}

func UserPrivate(e *echo.Echo, h handlers.UserHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/private", middlewares...)

	g.GET("/me", h.Me)
}
