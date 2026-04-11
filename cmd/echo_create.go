package main

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTP(errorHandler echo.HTTPErrorHandler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) //para el panic

	//permitir a fronend hacer pruebas
	corsConfig := middleware.CORSConfig{
		//separamos por coma ","
		AllowOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowMethods: strings.Split(os.Getenv("ALLOWED_METHODS"), ","),
	}

	e.Use(middleware.CORSWithConfig(corsConfig))

	// HTTPErrorHandler is a centralized HTTP error handler.
	// type HTTPErrorHandler func(err error, c Context)
	e.HTTPErrorHandler = errorHandler

	return e
}
