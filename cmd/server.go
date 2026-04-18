package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/cmd/routes"
	userports "github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/middle"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
)

type Server struct {
	uService       userports.Service
	uHandler       handlers.UserHandler
	verificationH  handlers.EmailVerificationHandler
	projectCatH    handlers.ProjectAdminCatalogHandler
	productCompatH handlers.ProductPublicCompatHandler
	lHandler       handlers.LoginHandler
	projHandler    handlers.ProjectPublicHandler
	assistantH     handlers.ProjectAssistantHandlerContract
	sHandler       handlers.SearchHandler
	searchAdminH   handlers.SearchAdminHandler
	techHandler    *handlers.TechnologyHandler
	projAdminH     *handlers.ProjectAdminHandler
	siteConfigH    *handlers.SiteSettingsHandler
	workflowH      *handlers.CaseStudyWorkflowHandler
}

func NewServer(
	uService userports.Service,
	uHandler handlers.UserHandler,
	verificationH handlers.EmailVerificationHandler,
	projectCatH handlers.ProjectAdminCatalogHandler,
	productCompatH handlers.ProductPublicCompatHandler,
	lHandler handlers.LoginHandler,
	projHandler handlers.ProjectPublicHandler,
	assistantH handlers.ProjectAssistantHandlerContract,
	sHandler handlers.SearchHandler,
	searchAdminH handlers.SearchAdminHandler,
	techHandler *handlers.TechnologyHandler,
	projAdminH *handlers.ProjectAdminHandler,
	siteConfigH *handlers.SiteSettingsHandler,
	workflowH *handlers.CaseStudyWorkflowHandler,
) *Server {

	return &Server{
		uService:       uService,
		uHandler:       uHandler,
		verificationH:  verificationH,
		projectCatH:    projectCatH,
		productCompatH: productCompatH,
		lHandler:       lHandler,
		projHandler:    projHandler,
		assistantH:     assistantH,
		sHandler:       sHandler,
		searchAdminH:   searchAdminH,
		techHandler:    techHandler,
		projAdminH:     projAdminH,
		siteConfigH:    siteConfigH,
		workflowH:      workflowH,
	}
}

func (s *Server) Initialize() {

	e := NewHTTP(response.HTTPErrorHandler)

	health(e) //esto es para verificar que el servicio está funcionando

	authMiddleware := middle.New(s.uService)

	routes.UserAdmin(e, s.uHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.UserPublic(e, s.uHandler, s.verificationH)
	routes.UserPrivate(e, s.uHandler, authMiddleware.IsValid)

	routes.ProductAdmin(e, s.projectCatH, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.ProjectAdminCatalog(e, s.projectCatH, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.ProductPublicCompat(e, s.productCompatH)

	routes.TechnologyAdmin(e, s.techHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.TechnologyPublic(e, s.techHandler)

	routes.LoginPublic(e, s.lHandler)
	routes.LoginAdmin(e, s.lHandler)

	routes.ProjectPublic(e, s.projHandler)
	routes.ProjectPrivate(e, s.assistantH, authMiddleware.IsValid)
	routes.ProjectAdmin(e, s.projAdminH, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.SiteSettingsPublic(e, s.siteConfigH)
	routes.SiteSettingsAdmin(e, s.siteConfigH, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.CaseStudyWorkflowAdmin(e, s.workflowH, authMiddleware.IsValid, authMiddleware.IsAdmin)

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
