package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	projectport "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type stubProjectCatalogService struct {
	current model.AdminProject
	updated *model.AdminProjectWrite
}

var _ projectport.AdminCatalogService = (*stubProjectCatalogService)(nil)

func (s *stubProjectCatalogService) Create(*model.AdminProjectWrite) error { return nil }
func (s *stubProjectCatalogService) Delete(uuid.UUID) error                { return nil }
func (s *stubProjectCatalogService) UpdateStatus(uuid.UUID, bool) (model.AdminProject, error) {
	return s.current, nil
}
func (s *stubProjectCatalogService) CreateVariants(uuid.UUID, []model.AdminProjectVariantInput) error {
	return nil
}
func (s *stubProjectCatalogService) ReplaceVariants(uuid.UUID, []model.AdminProjectVariantInput) error {
	return nil
}
func (s *stubProjectCatalogService) ReplaceMedia(uuid.UUID, []model.ProjectMedia) error { return nil }
func (s *stubProjectCatalogService) GetAdminAll() ([]model.AdminProject, error) {
	return []model.AdminProject{s.current}, nil
}

func (s *stubProjectCatalogService) Update(req *model.AdminProjectWrite) error {
	copyReq := *req
	s.updated = &copyReq
	s.current = model.AdminProject{
		ID:                req.ID,
		Name:              req.Name,
		Slug:              "portfolioforge",
		Description:       req.Description,
		Category:          req.Category,
		SourceMarkdownURL: req.SourceMarkdownURL,
		Active:            req.ResolveActive(true),
	}
	return nil
}

func (s *stubProjectCatalogService) GetAdminByID(id uuid.UUID) (model.AdminProject, error) {
	if s.current.ID == uuid.Nil {
		s.current = model.AdminProject{ID: id, Name: "PortfolioForge", Slug: "portfolioforge", Description: "Original", Category: "platform", Active: true}
	}
	return s.current, nil
}

func TestProjectCatalogUpdatePersistsSourceMarkdownURL(t *testing.T) {
	projectID := uuid.New()
	service := &stubProjectCatalogService{current: model.AdminProject{ID: projectID, Name: "PortfolioForge", Slug: "portfolioforge", Description: "Original", Category: "platform", Active: true}}
	handler := &ProjectCatalog{service: service, responser: *response.New()}

	rec := performProjectCatalogUpdate(t, handler, projectID, map[string]any{
		"name":                "PortfolioForge",
		"description":         "Updated project",
		"category":            "platform",
		"source_markdown_url": "https://mlbautomation.com/docs.md",
		"images":              []string{},
		"active":              true,
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if service.updated == nil {
		t.Fatal("update request was not passed to service")
	}
	if service.updated.SourceMarkdownURL != "https://mlbautomation.com/docs.md" {
		t.Fatalf("source_markdown_url = %q", service.updated.SourceMarkdownURL)
	}

	var payload map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["data"]["source_markdown_url"] != "https://mlbautomation.com/docs.md" {
		t.Fatalf("response source_markdown_url = %#v", payload["data"]["source_markdown_url"])
	}
}

func TestProjectCatalogUpdateClearsSourceMarkdownURL(t *testing.T) {
	projectID := uuid.New()
	service := &stubProjectCatalogService{current: model.AdminProject{ID: projectID, Name: "PortfolioForge", Slug: "portfolioforge", Description: "Original", Category: "platform", SourceMarkdownURL: "https://mlbautomation.com/docs.md", Active: true}}
	handler := &ProjectCatalog{service: service, responser: *response.New()}

	rec := performProjectCatalogUpdate(t, handler, projectID, map[string]any{
		"name":                "PortfolioForge",
		"description":         "Updated project",
		"category":            "platform",
		"source_markdown_url": "",
		"images":              []string{},
		"active":              true,
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if service.updated == nil {
		t.Fatal("update request was not passed to service")
	}
	if service.updated.SourceMarkdownURL != "" {
		t.Fatalf("source_markdown_url = %q, want empty string", service.updated.SourceMarkdownURL)
	}

	var payload map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, exists := payload["data"]["source_markdown_url"]; exists {
		t.Fatalf("response still includes source_markdown_url after clear: %#v", payload["data"])
	}
}

func TestProjectCatalogUpdateRejectsLegacyMediaKeys(t *testing.T) {
	tests := []struct {
		name        string
		legacyKey   string
		legacyValue string
		wantMessage string
	}{
		{
			name:        "thumbnail url",
			legacyKey:   "thumbnail_url",
			legacyValue: "https://cdn.example.com/project-low.webp",
			wantMessage: "La variante de media \"thumbnail_url\" ya no es válida. Usa \"low_url\".",
		},
		{
			name:        "full url",
			legacyKey:   "full_url",
			legacyValue: "https://cdn.example.com/project-high.webp",
			wantMessage: "La variante de media \"full_url\" ya no es válida. Usa \"high_url\".",
		},
		{
			name:        "url",
			legacyKey:   "url",
			legacyValue: "https://cdn.example.com/project-fallback.webp",
			wantMessage: "La variante de media \"url\" ya no es válida. Usa \"fallback_url\".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID := uuid.New()
			service := &stubProjectCatalogService{current: model.AdminProject{ID: projectID, Name: "PortfolioForge", Slug: "portfolioforge", Description: "Original", Category: "platform", Active: true}}
			handler := &ProjectCatalog{service: service, responser: *response.New()}

			rec := performProjectCatalogUpdate(t, handler, projectID, map[string]any{
				"name":        "PortfolioForge",
				"description": "Updated project",
				"category":    "platform",
				"images":      []string{},
				"active":      true,
				"media": []map[string]any{{
					"id":         uuid.NewString(),
					"media_type": "image",
					tt.legacyKey: tt.legacyValue,
				}},
			})

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}

			var payload map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			errorPayload, ok := payload["error"].(map[string]any)
			if !ok {
				t.Fatalf("error payload = %#v", payload)
			}
			if errorPayload["code"] != "validation_error" {
				t.Fatalf("error code = %#v", errorPayload["code"])
			}
			if errorPayload["message"] != tt.wantMessage {
				t.Fatalf("error message = %#v", errorPayload["message"])
			}
		})
	}
}

func TestProjectCatalogGetByIDReturnsCanonicalMediaKeys(t *testing.T) {
	projectID := uuid.New()
	service := &stubProjectCatalogService{current: model.AdminProject{
		ID:          projectID,
		Name:        "PortfolioForge",
		Slug:        "portfolioforge",
		Description: "Original",
		Category:    "platform",
		Active:      true,
		Media: []model.ProjectMedia{{
			ID:          uuid.New(),
			ProjectID:   projectID,
			MediaType:   "image",
			LowURL:      "https://cdn.example.com/project-low.webp",
			MediumURL:   "https://cdn.example.com/project-medium.webp",
			HighURL:     "https://cdn.example.com/project-high.webp",
			FallbackURL: "https://cdn.example.com/project-fallback.webp",
			SortOrder:   0,
			Featured:    true,
		}},
	}}
	handler := &ProjectCatalog{service: service, responser: *response.New()}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/projects/"+projectID.String(), nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/projects/:id")
	c.SetParamNames("id")
	c.SetParamValues(projectID.String())

	if err := handler.GetByID(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	media, ok := payload["data"]["media"].([]any)
	if !ok || len(media) != 1 {
		t.Fatalf("media = %#v", payload["data"]["media"])
	}
	item, ok := media[0].(map[string]any)
	if !ok {
		t.Fatalf("media item = %#v", media[0])
	}
	if item["low_url"] != "https://cdn.example.com/project-low.webp" {
		t.Fatalf("low_url = %#v", item["low_url"])
	}
	if item["medium_url"] != "https://cdn.example.com/project-medium.webp" {
		t.Fatalf("medium_url = %#v", item["medium_url"])
	}
	if item["high_url"] != "https://cdn.example.com/project-high.webp" {
		t.Fatalf("high_url = %#v", item["high_url"])
	}
	if item["fallback_url"] != "https://cdn.example.com/project-fallback.webp" {
		t.Fatalf("fallback_url = %#v", item["fallback_url"])
	}
	if _, exists := item["thumbnail_url"]; exists {
		t.Fatalf("thumbnail_url leaked in admin media response: %#v", item)
	}
	if _, exists := item["full_url"]; exists {
		t.Fatalf("full_url leaked in admin media response: %#v", item)
	}
	if _, exists := item["url"]; exists {
		t.Fatalf("url leaked in admin media response: %#v", item)
	}
}

func performProjectCatalogUpdate(t *testing.T, handler *ProjectCatalog, projectID uuid.UUID, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/projects/"+projectID.String(), bytes.NewBuffer(encoded))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/projects/:id")
	c.SetParamNames("id")
	c.SetParamValues(projectID.String())

	if err := handler.Update(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	return rec
}
