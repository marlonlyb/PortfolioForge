package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/cmd/routes"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/middle"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
)

type Server struct {
	uHandler     handlers.UserHandler
	pHandler     handlers.ProductHandler
	lHandler     handlers.LoginHandler
	projHandler  handlers.ProjectPublicHandler
	sHandler     handlers.SearchHandler
	searchAdminH handlers.SearchAdminHandler
	techHandler  *handlers.TechnologyHandler
	projAdminH   *handlers.ProjectAdminHandler
}

func NewServer(
	uHandler handlers.UserHandler,
	pHandler handlers.ProductHandler,
	lHandler handlers.LoginHandler,
	projHandler handlers.ProjectPublicHandler,
	sHandler handlers.SearchHandler,
	searchAdminH handlers.SearchAdminHandler,
	techHandler *handlers.TechnologyHandler,
	projAdminH *handlers.ProjectAdminHandler,
) *Server {

	return &Server{
		uHandler:     uHandler,
		pHandler:     pHandler,
		lHandler:     lHandler,
		projHandler:  projHandler,
		sHandler:     sHandler,
		searchAdminH: searchAdminH,
		techHandler:  techHandler,
		projAdminH:   projAdminH,
	}
}

func (s *Server) Initialize() {

	e := NewHTTP(response.HTTPErrorHandler)

	health(e) //esto es para verificar que el servicio está funcionando

	authMiddleware := middle.New()

	routes.UserAdmin(e, s.uHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.UserPublic(e, s.uHandler)
	routes.UserPrivate(e, s.uHandler, authMiddleware.IsValid)

	routes.ProductAdmin(e, s.pHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.ProductPublic(e, s.pHandler)

	routes.TechnologyAdmin(e, s.techHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.TechnologyPublic(e, s.techHandler)

	routes.LoginPublic(e, s.lHandler)

	routes.ProjectPublic(e, s.projHandler)
	routes.ProjectAdmin(e, s.projAdminH, authMiddleware.IsValid, authMiddleware.IsAdmin)

	routes.SearchPublic(e, s.sHandler)
	routes.SearchAdmin(e, s.searchAdminH, authMiddleware.IsValid, authMiddleware.IsAdmin)

	err := e.Start(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}

func health(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"time":         time.Now().String(),
				"message":      "PortfolioForge is running",
				"service_name": "PortfolioForge API",
			},
		)
	})
}
