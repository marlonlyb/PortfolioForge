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
