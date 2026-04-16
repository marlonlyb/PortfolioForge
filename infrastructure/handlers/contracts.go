package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	Create(c echo.Context) error
	GetAll(c echo.Context) error
	AdminGetByID(c echo.Context) error
	AdminUpdate(c echo.Context) error
	AdminDelete(c echo.Context) error
	Me(c echo.Context) error
	UpdateProfile(c echo.Context) error
}

type EmailVerificationHandler interface {
	Request(c echo.Context) error
	Resend(c echo.Context) error
	Verify(c echo.Context) error
}

// ProjectAdminCatalogHandler defines canonical admin project catalog operations.
type ProjectAdminCatalogHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetByID(c echo.Context) error
	GetAllStore(c echo.Context) error
	UpdateStatus(c echo.Context) error
}

// ProductPublicCompatHandler isolates legacy `/public/products` reads.
type ProductPublicCompatHandler interface {
	GetStoreByID(c echo.Context) error
	GetStoreAll(c echo.Context) error
}

type LoginHandler interface {
	AdminLogin(c echo.Context) error
	PublicLogin(c echo.Context) error
	PublicSignup(c echo.Context) error
	LoginWithGoogle(c echo.Context) error
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

type ProjectAssistantHandlerContract interface {
	CreateMessage(c echo.Context) error
}

// SearchAdminHandler defines the interface for admin search management endpoints.
type SearchAdminHandler interface {
	GetReadiness(c echo.Context) error
	ReembedProject(c echo.Context) error
	ReembedStale(c echo.Context) error
}

type SiteSettingsHandlerContract interface {
	GetPublic(c echo.Context) error
	GetAdmin(c echo.Context) error
	SaveAdmin(c echo.Context) error
}
