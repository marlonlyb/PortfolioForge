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
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/model"
)

type stubProjectLocalizationRepo struct {
	rowsByProject map[uuid.UUID][]model.ProjectLocalization
}

func (s *stubProjectLocalizationRepo) ListByProjectID(context.Context, uuid.UUID) ([]model.ProjectLocalization, error) {
	return nil, nil
}

func (s *stubProjectLocalizationRepo) ListByProjectIDsAndLocale(_ context.Context, projectIDs []uuid.UUID, locale string) (map[uuid.UUID][]model.ProjectLocalization, error) {
	rows := map[uuid.UUID][]model.ProjectLocalization{}
	for _, projectID := range projectIDs {
		for _, row := range s.rowsByProject[projectID] {
			if row.Locale == locale {
				rows[projectID] = append(rows[projectID], row)
			}
		}
	}
	return rows, nil
}

func (s *stubProjectLocalizationRepo) UpsertAuto(context.Context, uuid.UUID, string, map[string]json.RawMessage, map[string]string) error {
	return nil
}

func (s *stubProjectLocalizationRepo) UpsertManual(context.Context, uuid.UUID, string, map[string]json.RawMessage) error {
	return nil
}

type stubPublicProjectRepo struct {
	project model.Project
}

func (s *stubPublicProjectRepo) GetByID(context.Context, uuid.UUID) (model.Project, error) {
	return s.project, nil
}
func (s *stubPublicProjectRepo) GetBySlug(context.Context, string) (model.Project, error) {
	return s.project, nil
}
func (s *stubPublicProjectRepo) ListPublished(context.Context) ([]model.Project, error) {
	return []model.Project{s.project}, nil
}
func (s *stubPublicProjectRepo) GetTechnologiesByProjectID(context.Context, uuid.UUID) ([]model.Technology, error) {
	return nil, nil
}
func (s *stubPublicProjectRepo) GetAssistantContextBySlug(context.Context, string) (model.ProjectAssistantContext, error) {
	return model.ProjectAssistantContext{}, nil
}

func TestProjectPublicGetBySlugReturnsAssistantAvailabilityOnly(t *testing.T) {
	handler := NewProjectPublic(services.NewProject(&stubPublicProjectRepo{project: model.Project{
		ID:                 uuid.New(),
		Name:               "PortfolioForge",
		Slug:               "portfolioforge",
		Description:        "Detailed project",
		Category:           "platform",
		Status:             "published",
		Active:             true,
		AssistantAvailable: true,
		Media: []model.ProjectMedia{{
			ID:          uuid.New(),
			ProjectID:   uuid.New(),
			MediaType:   "image",
			LowURL:      "https://cdn.example.com/project-low.webp",
			MediumURL:   "https://cdn.example.com/project-medium.webp",
			HighURL:     "https://cdn.example.com/project-high.webp",
			FallbackURL: "https://cdn.example.com/project-fallback.webp",
			SortOrder:   0,
			Featured:    true,
		}},
	}}), localization.NewService(nil, nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/projects/portfolioforge", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/public/projects/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("portfolioforge")

	if err := handler.GetBySlug(c); err != nil {
		t.Fatalf("GetBySlug() error = %v", err)
	}

	var body map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := body["data"]
	if data["assistant_available"] != true {
		t.Fatalf("assistant_available = %#v, want true", data["assistant_available"])
	}
	if _, exists := data["source_markdown_url"]; exists {
		t.Fatal("source_markdown_url leaked in public response")
	}

	media, ok := data["media"].([]any)
	if !ok || len(media) != 1 {
		t.Fatalf("media = %#v", data["media"])
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
		t.Fatalf("thumbnail_url leaked in public media response: %#v", item)
	}
	if _, exists := item["full_url"]; exists {
		t.Fatalf("full_url leaked in public media response: %#v", item)
	}
	if _, exists := item["url"]; exists {
		t.Fatalf("url leaked in public media response: %#v", item)
	}
}

func TestProjectPublicGetBySlugReturnsAssistantUnavailableWhenMarkdownMissing(t *testing.T) {
	handler := NewProjectPublic(services.NewProject(&stubPublicProjectRepo{project: model.Project{
		ID:                 uuid.New(),
		Name:               "PortfolioForge",
		Slug:               "portfolioforge",
		Description:        "Detailed project",
		Category:           "platform",
		Status:             "published",
		Active:             true,
		AssistantAvailable: false,
	}}), localization.NewService(nil, nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/projects/portfolioforge", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/public/projects/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("portfolioforge")

	if err := handler.GetBySlug(c); err != nil {
		t.Fatalf("GetBySlug() error = %v", err)
	}

	var body map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := body["data"]
	if data["assistant_available"] != false {
		t.Fatalf("assistant_available = %#v, want false", data["assistant_available"])
	}
	if _, exists := data["source_markdown_url"]; exists {
		t.Fatal("source_markdown_url leaked in public response")
	}
}

func TestProjectPublicGetBySlugLocalizesPublicFields(t *testing.T) {
	projectID := uuid.New()
	handler := NewProjectPublic(
		services.NewProject(&stubPublicProjectRepo{project: model.Project{
			ID:                 projectID,
			Name:               "PortfolioForge",
			Slug:               "portfolioforge",
			Description:        "Proyecto original",
			Category:           "platform",
			Status:             "published",
			Active:             true,
			AssistantAvailable: true,
		}}),
		localization.NewService(&stubProjectLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
			projectID: {
				{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "name", Value: json.RawMessage(`"PortfolioForge EN"`)},
				{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "description", Value: json.RawMessage(`"Localized description"`)},
			},
		}}, nil),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/projects/portfolioforge?lang=en", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/public/projects/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("portfolioforge")

	if err := handler.GetBySlug(c); err != nil {
		t.Fatalf("GetBySlug() error = %v", err)
	}

	var body map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := body["data"]
	if data["name"] != "PortfolioForge EN" {
		t.Fatalf("name = %#v, want localized value", data["name"])
	}
	if data["description"] != "Localized description" {
		t.Fatalf("description = %#v, want localized value", data["description"])
	}
	if _, exists := data["source_markdown_url"]; exists {
		t.Fatal("source_markdown_url leaked in localized public response")
	}
}
