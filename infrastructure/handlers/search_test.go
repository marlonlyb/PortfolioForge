package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/embedding"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/model"
)

type stubSearchPublicRepo struct {
	results []model.SearchResult
	params  []model.SearchParams
}

func (s *stubSearchPublicRepo) Search(_ context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	s.params = append(s.params, params)
	return s.results, nil
}

func (s *stubSearchPublicRepo) RefreshSearchDocument(context.Context, uuid.UUID) error {
	return nil
}

func (s *stubSearchPublicRepo) RefreshAllDocuments(context.Context) error {
	return nil
}

type stubSearchProjectReader struct{}

func (s *stubSearchProjectReader) GetByID(context.Context, uuid.UUID) (model.Project, error) {
	return model.Project{}, nil
}
func (s *stubSearchProjectReader) GetBySlug(context.Context, string) (model.Project, error) {
	return model.Project{}, nil
}
func (s *stubSearchProjectReader) ListPublished(context.Context) ([]model.Project, error) {
	return nil, nil
}
func (s *stubSearchProjectReader) GetTechnologiesByProjectID(context.Context, uuid.UUID) ([]model.Technology, error) {
	return nil, nil
}
func (s *stubSearchProjectReader) GetAssistantContextBySlug(context.Context, string) (model.ProjectAssistantContext, error) {
	return model.ProjectAssistantContext{}, nil
}

type stubExplanationProvider struct{}

func (s *stubExplanationProvider) Explain(context.Context, model.Project, model.EvidenceTrace, string) (string, error) {
	return "", nil
}

func TestSearchHandlerLocalizesClientNameInResults(t *testing.T) {
	projectID := uuid.New()
	searchService := services.NewSearch(
		&stubSearchPublicRepo{results: []model.SearchResult{{
			Project: model.Project{
				ID:          projectID,
				Slug:        "portfolioforge",
				Name:        "Proyecto base",
				Description: "Descripción base",
				Category:    "platform",
				ClientName:  "Cliente base",
				Profile:     &model.ProjectProfile{SolutionSummary: "Resumen base"},
			},
			Score: 1,
		}}},
		&stubSearchProjectReader{},
		embedding.NewNoOpEmbeddingProvider(),
		&stubExplanationProvider{},
		false,
	)
	handler := NewSearch(
		searchService,
		false,
		localization.NewService(&stubProjectLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
			projectID: {{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "client_name", Value: json.RawMessage(`"Acme Industries"`)}},
		}}, nil),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/search?lang=en", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	if err := handler.Search(c); err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	var body struct {
		Data struct {
			Data []struct {
				ClientName *string `json:"client_name"`
			} `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Data) != 1 || body.Data.Data[0].ClientName == nil || *body.Data.Data[0].ClientName != "Acme Industries" {
		t.Fatalf("client_name = %#v", body.Data.Data)
	}
}

func TestSearchHandlerPreservesCursorAndFiltersAcrossServiceBoundary(t *testing.T) {
	projectID := uuid.New()
	searchRepo := &stubSearchPublicRepo{results: []model.SearchResult{{
		Project: model.Project{
			ID:          projectID,
			Slug:        "portfolioforge",
			Name:        "PortfolioForge",
			Description: "Descripción base",
			Category:    "platform",
			ClientName:  "Acme",
			Profile:     &model.ProjectProfile{ProjectID: projectID, SolutionSummary: "Resumen base"},
		},
		LexicalScore: 1,
	}}}
	searchService := services.NewSearch(
		searchRepo,
		&stubSearchProjectReader{},
		embedding.NewNoOpEmbeddingProvider(),
		&stubExplanationProvider{},
		false,
	)
	handler := NewSearch(searchService, false, localization.NewService(nil, nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/search?q=portfolio&category=platform&client=Acme&technologies=react,node&pageSize=5&cursor=2", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	if err := handler.Search(c); err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if len(searchRepo.params) != 1 {
		t.Fatalf("repo calls = %d, want 1", len(searchRepo.params))
	}
	got := searchRepo.params[0]
	if got.Category != "platform" || got.Client != "Acme" {
		t.Fatalf("filters = %#v", got)
	}
	if got.Cursor != "2" {
		t.Fatalf("cursor = %q, want 2", got.Cursor)
	}
	if len(got.Technologies) != 2 || got.Technologies[0] != "react" || got.Technologies[1] != "node" {
		t.Fatalf("technologies = %#v", got.Technologies)
	}

	var body struct {
		Data struct {
			Meta struct {
				FiltersApplied struct {
					Client       *string  `json:"client"`
					Category     *string  `json:"category"`
					Technologies []string `json:"technologies"`
				} `json:"filters_applied"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Meta.FiltersApplied.Client == nil || *body.Data.Meta.FiltersApplied.Client != "Acme" {
		t.Fatalf("client filter = %#v", body.Data.Meta.FiltersApplied.Client)
	}
	if body.Data.Meta.FiltersApplied.Category == nil || *body.Data.Meta.FiltersApplied.Category != "platform" {
		t.Fatalf("category filter = %#v", body.Data.Meta.FiltersApplied.Category)
	}
	if len(body.Data.Meta.FiltersApplied.Technologies) != 2 {
		t.Fatalf("technology filters = %#v", body.Data.Meta.FiltersApplied.Technologies)
	}
}
