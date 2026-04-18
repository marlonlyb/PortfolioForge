package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubCaseStudyWorkflowService struct {
	startRunFn    func(context.Context, model.StartCaseStudyWorkflowRunRequest) (model.CaseStudyWorkflowRun, error)
	getRunFn      func(context.Context, uuid.UUID) (model.CaseStudyWorkflowRun, error)
	listLogsFn    func(context.Context, uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error)
	confirmStepFn func(context.Context, uuid.UUID, string) (model.CaseStudyWorkflowRun, error)
	startStepFn   func(context.Context, uuid.UUID, string) (model.CaseStudyWorkflowRun, error)
	retryStepFn   func(context.Context, uuid.UUID, string) (model.CaseStudyWorkflowRun, error)
	resumeFn      func(context.Context, uuid.UUID) (model.CaseStudyWorkflowRun, error)
}

func (s stubCaseStudyWorkflowService) StartRun(ctx context.Context, req model.StartCaseStudyWorkflowRunRequest) (model.CaseStudyWorkflowRun, error) {
	if s.startRunFn != nil {
		return s.startRunFn(ctx, req)
	}
	return model.CaseStudyWorkflowRun{}, nil
}
func (s stubCaseStudyWorkflowService) GetRun(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	if s.getRunFn != nil {
		return s.getRunFn(ctx, runID)
	}
	return model.CaseStudyWorkflowRun{}, nil
}
func (s stubCaseStudyWorkflowService) ListLogs(_ context.Context, runID uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error) {
	if s.listLogsFn != nil {
		return s.listLogsFn(context.Background(), runID)
	}
	return []model.CaseStudyWorkflowLogEntry{{ID: 1, RunID: runID, Step: model.CaseStudyWorkflowStepPublishCanonical, Level: model.CaseStudyWorkflowLogInfo, Message: "published"}}, nil
}
func (s stubCaseStudyWorkflowService) ConfirmStep(ctx context.Context, runID uuid.UUID, step string) (model.CaseStudyWorkflowRun, error) {
	if s.confirmStepFn != nil {
		return s.confirmStepFn(ctx, runID, step)
	}
	return model.CaseStudyWorkflowRun{}, nil
}
func (s stubCaseStudyWorkflowService) StartStep(ctx context.Context, runID uuid.UUID, step string) (model.CaseStudyWorkflowRun, error) {
	if s.startStepFn != nil {
		return s.startStepFn(ctx, runID, step)
	}
	return model.CaseStudyWorkflowRun{}, nil
}
func (s stubCaseStudyWorkflowService) RetryStep(ctx context.Context, runID uuid.UUID, step string) (model.CaseStudyWorkflowRun, error) {
	if s.retryStepFn != nil {
		return s.retryStepFn(ctx, runID, step)
	}
	return model.CaseStudyWorkflowRun{}, nil
}
func (s stubCaseStudyWorkflowService) Resume(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	if s.resumeFn != nil {
		return s.resumeFn(ctx, runID)
	}
	return model.CaseStudyWorkflowRun{}, nil
}

func TestCaseStudyWorkflowHandler_GetLogs(t *testing.T) {
	runID := uuid.New()
	handler := NewCaseStudyWorkflowHandler(stubCaseStudyWorkflowService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/settings/case-study-runs/"+runID.String()+"/logs", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/settings/case-study-runs/:id/logs")
	c.SetParamNames("id")
	c.SetParamValues(runID.String())

	if err := handler.GetLogs(c); err != nil {
		t.Fatalf("GetLogs() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data struct {
			Items []model.CaseStudyWorkflowLogEntry `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload.Data.Items) != 1 || payload.Data.Items[0].RunID != runID {
		t.Fatalf("unexpected logs payload: %#v", payload.Data.Items)
	}
}

func TestCaseStudyWorkflowHandler_StartRun(t *testing.T) {
	var received model.StartCaseStudyWorkflowRunRequest
	handler := NewCaseStudyWorkflowHandler(stubCaseStudyWorkflowService{
		startRunFn: func(_ context.Context, req model.StartCaseStudyWorkflowRunRequest) (model.CaseStudyWorkflowRun, error) {
			received = req
			return model.CaseStudyWorkflowRun{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")}, nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/settings/case-study-runs", strings.NewReader(`{"source_path":"/safe/root/demo","locales":["ca","en"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	if err := handler.StartRun(c); err != nil {
		t.Fatalf("StartRun() error = %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if received.SourcePath != "/safe/root/demo" || len(received.Locales) != 2 {
		t.Fatalf("received request = %#v", received)
	}
}

func TestCaseStudyWorkflowHandler_ConfirmStep(t *testing.T) {
	runID := uuid.New()
	called := false
	handler := NewCaseStudyWorkflowHandler(stubCaseStudyWorkflowService{
		confirmStepFn: func(_ context.Context, gotRunID uuid.UUID, step string) (model.CaseStudyWorkflowRun, error) {
			called = true
			if gotRunID != runID {
				t.Fatalf("runID = %s, want %s", gotRunID, runID)
			}
			if step != model.CaseStudyWorkflowStepPublishCanonical {
				t.Fatalf("step = %s", step)
			}
			return model.CaseStudyWorkflowRun{ID: gotRunID}, nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/settings/case-study-runs/"+runID.String()+"/steps/"+model.CaseStudyWorkflowStepPublishCanonical+"/confirm", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/settings/case-study-runs/:id/steps/:step/confirm")
	c.SetParamNames("id", "step")
	c.SetParamValues(runID.String(), model.CaseStudyWorkflowStepPublishCanonical)

	if err := handler.ConfirmStep(c); err != nil {
		t.Fatalf("ConfirmStep() error = %v", err)
	}
	if !called {
		t.Fatal("expected confirm step service call")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCaseStudyWorkflowHandler_ResumeMapsValidationErrors(t *testing.T) {
	runID := uuid.New()
	handler := NewCaseStudyWorkflowHandler(stubCaseStudyWorkflowService{
		resumeFn: func(context.Context, uuid.UUID) (model.CaseStudyWorkflowRun, error) {
			return model.CaseStudyWorkflowRun{}, errors.New("step publish_canonical todavía no está listo para ejecutarse")
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/settings/case-study-runs/"+runID.String()+"/resume", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/settings/case-study-runs/:id/resume")
	c.SetParamNames("id")
	c.SetParamValues(runID.String())

	err := handler.Resume(c)
	if err == nil {
		t.Fatal("expected validation error")
	}
	contractErr := err.(*model.ContractError)
	if contractErr.StatusHTTP != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", contractErr.StatusHTTP, http.StatusBadRequest)
	}
	if !strings.Contains(contractErr.Response.Error.Message, "todavía no está listo") {
		t.Fatalf("message = %v", contractErr.Response.Error.Message)
	}
}
