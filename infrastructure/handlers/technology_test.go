package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/cmd/routes"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
	"github.com/marlonlyb/portfolioforge/model"
)

type stubTechnologyRepo struct {
	tech         model.Technology
	technologies []model.Technology
	requestedID  uuid.UUID
	deletedIDs   []uuid.UUID
	getByIDHits  int
	getAllHits   int
	deleteHits   int
}

func (s *stubTechnologyRepo) Create(_ *model.Technology) error { return nil }
func (s *stubTechnologyRepo) Update(_ *model.Technology) error { return nil }
func (s *stubTechnologyRepo) Delete(id uuid.UUID) error {
	s.deleteHits++
	s.deletedIDs = append(s.deletedIDs, id)
	filtered := make([]model.Technology, 0, len(s.technologies))
	for _, tech := range s.technologies {
		if tech.ID != id {
			filtered = append(filtered, tech)
		}
	}
	s.technologies = filtered
	return nil
}
func (s *stubTechnologyRepo) GetAll() ([]model.Technology, error) {
	s.getAllHits++
	items := make([]model.Technology, len(s.technologies))
	copy(items, s.technologies)
	return items, nil
}
func (s *stubTechnologyRepo) GetByID(id uuid.UUID) (model.Technology, error) {
	s.requestedID = id
	s.getByIDHits++
	return s.tech, nil
}

func TestTechnologyAdmin_GetByIDRoute(t *testing.T) {
	e := echo.New()
	techID := uuid.New()
	repo := &stubTechnologyRepo{tech: model.Technology{ID: techID, Name: "Go", Slug: "go", Category: "backend"}}
	handler := handlers.NewTechnologyHandler(repo)
	routes.TechnologyAdmin(e, handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/technologies/"+techID.String(), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if repo.getByIDHits != 1 {
		t.Fatalf("GetByID hits = %d, want 1", repo.getByIDHits)
	}
	if repo.requestedID != techID {
		t.Fatalf("requested ID = %s, want %s", repo.requestedID, techID)
	}

	var body struct {
		Data model.Technology `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.ID != techID {
		t.Fatalf("response technology ID = %s, want %s", body.Data.ID, techID)
	}
}

func TestTechnologyAdmin_ListThenDeleteRemovesTechnologyFromAdminListing(t *testing.T) {
	e := echo.New()
	goID := uuid.New()
	tsID := uuid.New()
	repo := &stubTechnologyRepo{technologies: []model.Technology{
		{ID: goID, Name: "Go", Slug: "go", Category: "backend"},
		{ID: tsID, Name: "TypeScript", Slug: "typescript", Category: "frontend"},
	}}
	handler := handlers.NewTechnologyHandler(repo)
	routes.TechnologyAdmin(e, handler)

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/technologies", nil)
	listRec := httptest.NewRecorder()
	e.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("initial list status = %d, want %d", listRec.Code, http.StatusOK)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/technologies/"+goID.String(), nil)
	deleteRec := httptest.NewRecorder()
	e.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("delete status = %d, want %d", deleteRec.Code, http.StatusOK)
	}
	if repo.deleteHits != 1 {
		t.Fatalf("Delete hits = %d, want 1", repo.deleteHits)
	}
	if len(repo.deletedIDs) != 1 || repo.deletedIDs[0] != goID {
		t.Fatalf("deleted IDs = %#v, want [%s]", repo.deletedIDs, goID)
	}

	listAfterReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/technologies", nil)
	listAfterRec := httptest.NewRecorder()
	e.ServeHTTP(listAfterRec, listAfterReq)

	if listAfterRec.Code != http.StatusOK {
		t.Fatalf("list after delete status = %d, want %d", listAfterRec.Code, http.StatusOK)
	}

	var body struct {
		Data struct {
			Items []model.Technology `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(listAfterRec.Body).Decode(&body); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if repo.getAllHits != 2 {
		t.Fatalf("GetAll hits = %d, want 2", repo.getAllHits)
	}
	if len(body.Data.Items) != 1 {
		t.Fatalf("listed technologies = %d, want 1", len(body.Data.Items))
	}
	if body.Data.Items[0].ID != tsID {
		t.Fatalf("remaining technology ID = %s, want %s", body.Data.Items[0].ID, tsID)
	}
}
