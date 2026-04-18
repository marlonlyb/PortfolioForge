package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

func CaseStudyWorkflowAdmin(e *echo.Echo, h *handlers.CaseStudyWorkflowHandler, middlewares ...echo.MiddlewareFunc) {
	e.GET("/api/v1/admin/settings/case-study-workflow", h.GetAvailability, middlewares...)

	g := e.Group("/api/v1/admin/settings/case-study-runs", middlewares...)

	g.POST("", h.StartRun)
	g.GET("/:id", h.GetRun)
	g.GET("/:id/logs", h.GetLogs)
	g.POST("/:id/resume", h.Resume)
	g.POST("/:id/steps/:step/confirm", h.ConfirmStep)
	g.POST("/:id/steps/:step/start", h.StartStep)
	g.POST("/:id/steps/:step/retry", h.RetryStep)
}
