package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

// SearchPublic registers public search routes with rate limiting (30 requests/minute)
// and CORS header (Access-Control-Allow-Origin: *) for public access.
func SearchPublic(e *echo.Echo, h handlers.SearchHandler) {
	g := e.Group("/api/v1/public/search")

	// CORS: allow all origins on the public search endpoint
	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions},
	}))

	// Rate limit: 30 requests per minute per IP
	// Rate is in requests/second, so 30/60 = 0.5 req/s
	store := middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      0.5,
			Burst:     30,
			ExpiresIn: 1 * time.Minute,
		},
	)

	g.Use(middleware.RateLimiter(store))
	g.GET("", h.Search)
}
