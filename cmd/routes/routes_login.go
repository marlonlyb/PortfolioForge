package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func LoginPublic(e *echo.Echo, h handlers.LoginHandler) {
	g := e.Group("/api/v1/public/login")

	g.POST("/google", h.LoginWithGoogle)
	g.POST("/email/request", h.RequestEmailLogin)
	g.POST("/email/verify", h.VerifyEmailLogin)
}

func LoginAdmin(e *echo.Echo, h handlers.LoginHandler) {
	e.POST("/api/v1/admin/login", h.AdminLogin)
}
