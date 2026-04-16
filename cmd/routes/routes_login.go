package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func LoginPublic(e *echo.Echo, h handlers.LoginHandler) {
	e.POST("/api/v1/public/signup", h.PublicSignup)
	e.POST("/api/v1/public/login", h.PublicLogin)
	e.POST("/api/v1/public/login/google", h.LoginWithGoogle)
}

func LoginAdmin(e *echo.Echo, h handlers.LoginHandler) {
	e.POST("/api/v1/admin/login", h.AdminLogin)
}
