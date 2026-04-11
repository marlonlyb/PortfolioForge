package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubProjectReader struct {
	project          model.Project
	requestedID      uuid.UUID
	getByIDHits      int
	technologiesResp []model.Technology
}

func (s *stubProjectReader) GetByID(_ context.Context, id uuid.UUID) (model.Project, error) {
	s.requestedID = id
	s.getByIDHits++
	return s.project, nil
}

func (s *stubProjectReader) GetBySlug(context.Context, string) (model.Project, error) {
	return model.Project{}, nil
}

func (s *stubProjectReader) ListPublished(context.Context) ([]model.Project, error) {
	return nil, nil
}

func (s *stubProjectReader) GetTechnologiesByProjectID(context.Context, uuid.UUID) ([]model.Technology, error) {
	return s.technologiesResp, nil
}

type stubSearchRepository struct {
	refreshedProjectID uuid.UUID
	refreshHits        int
}

func (s *stubSearchRepository) Search(context.Context, model.SearchParams) ([]model.SearchResult, error) {
	return nil, nil
}

func (s *stubSearchRepository) RefreshSearchDocument(_ context.Context, projectID uuid.UUID) error {
	s.refreshedProjectID = projectID
	s.refreshHits++
	return nil
}

func (s *stubSearchRepository) RefreshAllDocuments(context.Context) error {
	return nil
}

func TestSearchAdmin_ReembedProject(t *testing.T) {
	projectID := uuid.New()
	handler := NewSearchAdmin(
		&stubProjectReader{project: model.Project{ID: projectID, Name: "PortfolioForge"}},
		&stubSearchRepository{},
	)
	projectReader := handler.projectReader.(*stubProjectReader)
	searchRepo := handler.searchRepo.(*stubSearchRepository)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/projects/"+projectID.String()+"/reembed", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/projects/:id/reembed")
	c.SetParamNames("id")
	c.SetParamValues(projectID.String())

	err := handler.ReembedProject(c)
	if err != nil {
		t.Fatalf("ReembedProject() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if projectReader.getByIDHits != 1 {
		t.Fatalf("GetByID hits = %d, want 1", projectReader.getByIDHits)
	}
	if projectReader.requestedID != projectID {
		t.Fatalf("requested project ID = %s, want %s", projectReader.requestedID, projectID)
	}
	if searchRepo.refreshHits != 1 {
		t.Fatalf("RefreshSearchDocument hits = %d, want 1", searchRepo.refreshHits)
	}
	if searchRepo.refreshedProjectID != projectID {
		t.Fatalf("reembedded project ID = %s, want %s", searchRepo.refreshedProjectID, projectID)
	}

	var body struct {
		Data struct {
			Message   string `json:"message"`
			ProjectID string `json:"project_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.ProjectID != projectID.String() {
		t.Fatalf("response project_id = %q, want %q", body.Data.ProjectID, projectID.String())
	}
	if body.Data.Message == "" {
		t.Fatal("response message is empty")
	}
}
