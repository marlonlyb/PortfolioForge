package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	searchPorts "github.com/marlonlyb/portfolioforge/domain/ports/search"
	workflowPorts "github.com/marlonlyb/portfolioforge/domain/ports/workflow"
	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/model"
)

type workflowBootstrapRepository struct{}

func (workflowBootstrapRepository) SaveRun(context.Context, model.CaseStudyWorkflowRun) error {
	return nil
}
func (workflowBootstrapRepository) GetRun(context.Context, uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, nil
}
func (workflowBootstrapRepository) ListLogs(context.Context, uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error) {
	return nil, nil
}
func (workflowBootstrapRepository) AppendLog(context.Context, model.CaseStudyWorkflowLogEntry) error {
	return nil
}

type workflowBootstrapPublisher struct{}

func (workflowBootstrapPublisher) ResolvePublishTarget(string, string) (services.CaseStudyPublishTarget, error) {
	return services.CaseStudyPublishTarget{}, nil
}
func (workflowBootstrapPublisher) CollectFiles(string) ([]string, error) { return nil, nil }
func (workflowBootstrapPublisher) Publish(context.Context, services.CaseStudyPublishTarget, []string) error {
	return nil
}
func (workflowBootstrapPublisher) Verify(context.Context, string) error { return nil }

type workflowBootstrapImporter struct{}

func (workflowBootstrapImporter) ImportFromCanonical(context.Context, services.CaseStudyPublishTarget, string) (uuid.UUID, error) {
	return uuid.New(), nil
}

type workflowBootstrapLocalization struct{}

func (workflowBootstrapLocalization) BackfillProject(context.Context, uuid.UUID, []string) error {
	return nil
}

type workflowBootstrapSearchRepo struct{}

func (workflowBootstrapSearchRepo) RefreshSearchDocument(context.Context, uuid.UUID) error {
	return nil
}
func (workflowBootstrapSearchRepo) Search(context.Context, model.SearchParams) ([]model.SearchResult, error) {
	return nil, nil
}
func (workflowBootstrapSearchRepo) RefreshAllDocuments(context.Context) error { return nil }

var _ workflowPorts.Repository = workflowBootstrapRepository{}
var _ services.CaseStudyPublisher = workflowBootstrapPublisher{}
var _ services.CaseStudyProjectImporter = workflowBootstrapImporter{}
var _ services.CaseStudyLocalizationBackfiller = workflowBootstrapLocalization{}
var _ searchPorts.SearchRepository = workflowBootstrapSearchRepo{}

func TestBuildCaseStudyWorkflowHandler_UsesUnavailableServiceWhenWorkflowIsDisabled(t *testing.T) {
	handler, reason := buildCaseStudyWorkflowHandler(CaseStudyWorkflowEnvConfig{
		Configured: false,
		Reason:     caseStudyWorkflowUnavailableReason,
		Diagnostic: "missing workflow env",
	}, caseStudyWorkflowDependencies{})
	if reason == "" {
		t.Fatal("expected disabled diagnostic reason")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/settings/case-study-workflow", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	if err := handler.GetAvailability(c); err != nil {
		t.Fatalf("GetAvailability() error = %v", err)
	}

	var payload struct {
		Data model.CaseStudyWorkflowAvailability `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Data.Configured {
		t.Fatal("expected workflow to be disabled")
	}

	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/settings/case-study-runs", nil)
	postRec := httptest.NewRecorder()
	postCtx := echo.New().NewContext(postReq, postRec)
	err := handler.StartRun(postCtx)
	if err == nil {
		t.Fatal("expected workflow_unavailable contract error")
	}
	contractErr := err.(*model.ContractError)
	if contractErr.StatusHTTP != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", contractErr.StatusHTTP, http.StatusServiceUnavailable)
	}
}

func TestBuildCaseStudyWorkflowHandler_UsesConfiguredServiceWhenWorkflowEnvIsPresent(t *testing.T) {
	handler, reason := buildCaseStudyWorkflowHandler(CaseStudyWorkflowEnvConfig{
		Configured:         true,
		AllowedSourceRoots: []string{"/safe/root"},
		Reason:             caseStudyWorkflowUnavailableReason,
	}, caseStudyWorkflowDependencies{
		repository:   workflowBootstrapRepository{},
		publisher:    workflowBootstrapPublisher{},
		importer:     workflowBootstrapImporter{},
		localization: workflowBootstrapLocalization{},
		searchRepo:   workflowBootstrapSearchRepo{},
	})
	if reason != "" {
		t.Fatalf("unexpected disabled reason: %s", reason)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/settings/case-study-workflow", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	if err := handler.GetAvailability(c); err != nil {
		t.Fatalf("GetAvailability() error = %v", err)
	}

	var payload struct {
		Data model.CaseStudyWorkflowAvailability `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !payload.Data.Configured {
		t.Fatal("expected configured workflow availability")
	}
}
