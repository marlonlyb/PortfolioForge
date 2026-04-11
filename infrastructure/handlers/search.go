package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

// Search handles public search API requests.
type Search struct {
	service          *services.Search
	semanticDegraded bool // true when semantic search is requested but NoOp provider is in use
}

// NewSearch creates a new Search handler.
// semanticDegraded should be true when ENABLE_SEMANTIC_SEARCH is true but the
// embedding provider is a no-op (pgvector unavailable).
func NewSearch(service *services.Search, semanticDegraded bool) *Search {
	return &Search{service: service, semanticDegraded: semanticDegraded}
}

// Search handles GET /api/v1/public/search?q=...&category=...&client=...&technologies=...&pageSize=...&cursor=...
func (h *Search) Search(c echo.Context) error {
	q := c.QueryParam("q")

	// Validation: query length 2-200 (empty is allowed for all published)
	if q != "" {
		if len(q) < 2 {
			return response.ContractError(400, "validation_error", "La búsqueda debe tener al menos 2 caracteres")
		}
		if len(q) > 200 {
			return response.ContractError(400, "validation_error", "La búsqueda no puede exceder 200 caracteres")
		}
	}

	category := c.QueryParam("category")
	client := c.QueryParam("client")

	var technologies []string
	techParam := c.QueryParam("technologies")
	if techParam != "" {
		for _, t := range strings.Split(techParam, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				technologies = append(technologies, t)
			}
		}
	}

	pageSize := 20
	if ps := c.QueryParam("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	cursor := c.QueryParam("cursor")

	params := model.SearchParams{
		Query:        q,
		Category:     category,
		Client:       client,
		Technologies: technologies,
		Cursor:       cursor,
		PageSize:     pageSize,
	}

	searchResponse, err := h.service.Search(c.Request().Context(), params)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible realizar la búsqueda")
	}

	// 422: Unprocessable query — after normalization the query is empty (all stop words / special chars stripped)
	if q != "" && searchResponse.Meta.Query == "" {
		return response.ContractError(422, "unprocessable_query", "La consulta no contiene términos válidos para la búsqueda")
	}

	// 503: Semantic search requested but unavailable (NoOp provider)
	// Only fires when the user actually submitted a search query — listing all
	// published projects does not require semantic search and must return 200.
	if h.semanticDegraded && q != "" {
		c.Response().Header().Set("X-Search-Degraded", "semantic-unavailable")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"data": searchResponse.Data,
			"meta": searchResponse.Meta,
		})
	}

	return c.JSON(response.ContractOK(searchResponse))
}
