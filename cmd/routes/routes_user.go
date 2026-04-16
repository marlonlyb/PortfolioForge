package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func UserAdmin(e *echo.Echo, h handlers.UserHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/users", middlewares...)

	g.GET("", h.GetAll)
	g.GET("/:id", h.AdminGetByID)
	g.PATCH("/:id", h.AdminUpdate)
	g.DELETE("/:id", h.AdminDelete)
}

func UserPublic(e *echo.Echo, _ handlers.UserHandler, verificationHandler handlers.EmailVerificationHandler) {
	e.POST("/api/v1/public/email-verification/request", verificationHandler.Request)
	e.POST("/api/v1/public/email-verification/resend", verificationHandler.Resend)
	e.POST("/api/v1/public/email-verification/verify", verificationHandler.Verify)
}

func UserPrivate(e *echo.Echo, h handlers.UserHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/private", middlewares...)

	g.GET("/me", h.Me)
	g.PUT("/me/profile", h.UpdateProfile)
}
