package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	Create(c echo.Context) error
	Register(c echo.Context) error
	GetAll(c echo.Context) error
	Me(c echo.Context) error
}

type ProductHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetByID(c echo.Context) error
	GetAll(c echo.Context) error
	GetStoreByID(c echo.Context) error
	GetStoreAll(c echo.Context) error
	GetAllStore(c echo.Context) error
	UpdateStatus(c echo.Context) error
}

type LoginHandler interface {
	Login(c echo.Context) error
}

// SearchHandler defines the interface for the public search endpoint.
type SearchHandler interface {
	Search(c echo.Context) error
}

// ProjectPublicHandler defines the interface for public project read endpoints.
type ProjectPublicHandler interface {
	GetBySlug(c echo.Context) error
	ListPublished(c echo.Context) error
}

// SearchAdminHandler defines the interface for admin search management endpoints.
type SearchAdminHandler interface {
	GetReadiness(c echo.Context) error
	ReembedProject(c echo.Context) error
	ReembedStale(c echo.Context) error
}
