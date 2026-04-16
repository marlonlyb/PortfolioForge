package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type stubAssistantHandlerRepo struct {
	context model.ProjectAssistantContext
	err     error
}

func (s *stubAssistantHandlerRepo) GetAssistantContextBySlug(context.Context, string) (model.ProjectAssistantContext, error) {
	return s.context, s.err
}

type stubAssistantHandlerRetriever struct {
	chunks []services.MarkdownChunkAlias
	err    error
}

func (s *stubAssistantHandlerRetriever) Fetch(context.Context, string, string) ([]services.MarkdownChunkAlias, error) {
	return s.chunks, s.err
}

type stubAssistantHandlerProvider struct {
	input services.ProjectAssistantAnswerInput
	resp  string
	err   error
}

func (s *stubAssistantHandlerProvider) GenerateAnswer(_ context.Context, input services.ProjectAssistantAnswerInput) (string, error) {
	s.input = input
	return s.resp, s.err
}

func TestProjectAssistantHandlerCreateMessageReturnsAnswerOnly(t *testing.T) {
	provider := &stubAssistantHandlerProvider{resp: "Respuesta aterrizada."}
	handler := NewProjectAssistantHandler(services.NewProjectAssistant(
		&stubAssistantHandlerRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}},
		&stubAssistantHandlerRetriever{chunks: []services.MarkdownChunkAlias{{Heading: "Architecture", Body: "Uses Go."}}},
		provider,
	))

	body := `{"question":"¿Cómo está implementada la arquitectura?","history":[{"role":"assistant","content":"Resumen previo"}],"lang":"ca"}`
	rec := performAssistantRequest(t, handler, body, model.User{ID: uuid.New(), AuthProvider: "google", EmailVerified: true, FullName: "Ada Lovelace", Company: "Analytical Engines", ProfileCompleted: true, AssistantEligible: true, CanUseProjectAssistant: true})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := payload["data"]
	if data["answer"] != "Respuesta aterrizada." {
		t.Fatalf("answer = %#v", data["answer"])
	}
	if len(data) != 1 {
		t.Fatalf("response data leaked extra fields: %#v", data)
	}
	if provider.input.Language != "ca" {
		t.Fatalf("language = %q, want ca", provider.input.Language)
	}
	if provider.input.Question != "¿Cómo está implementada la arquitectura?" {
		t.Fatalf("question = %q", provider.input.Question)
	}
	if len(provider.input.History) != 1 || provider.input.History[0].Content != "Resumen previo" {
		t.Fatalf("history = %#v", provider.input.History)
	}
}

func TestProjectAssistantHandlerCreateMessageRejectsUnauthenticatedAccess(t *testing.T) {
	handler := NewProjectAssistantHandler(services.NewProjectAssistant(
		&stubAssistantHandlerRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}},
		&stubAssistantHandlerRetriever{chunks: []services.MarkdownChunkAlias{{Heading: "Architecture", Body: "Uses Go."}}},
		&stubAssistantHandlerProvider{resp: "Should not be used."},
	))

	rec := performAssistantRequest(t, handler, `{"question":"How does it work?"}`, model.User{})

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var payload model.APIErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Error.Code != "authentication_required" {
		t.Fatalf("error code = %q, want authentication_required", payload.Error.Code)
	}
}

func TestProjectAssistantHandlerCreateMessageRejectsIneligibleUser(t *testing.T) {
	handler := NewProjectAssistantHandler(services.NewProjectAssistant(
		&stubAssistantHandlerRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}},
		&stubAssistantHandlerRetriever{chunks: []services.MarkdownChunkAlias{{Heading: "Architecture", Body: "Uses Go."}}},
		&stubAssistantHandlerProvider{resp: "Should not be used."},
	))

	rec := performAssistantRequest(t, handler, `{"question":"How does it work?"}`, model.User{ID: uuid.New(), AuthProvider: "google", EmailVerified: true})

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}

	var payload model.APIErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Error.Code != "assistant_ineligible" {
		t.Fatalf("error code = %q, want assistant_ineligible", payload.Error.Code)
	}
}

func TestProjectAssistantHandlerCreateMessageMapsErrors(t *testing.T) {
	projectContext := model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}
	tests := []struct {
		name       string
		body       string
		service    services.ProjectAssistant
		wantStatus int
		wantCode   string
	}{
		{
			name:       "400 invalid question",
			body:       `{"question":"a"}`,
			service:    services.NewProjectAssistant(&stubAssistantHandlerRepo{}, &stubAssistantHandlerRetriever{}, &stubAssistantHandlerProvider{}),
			wantStatus: http.StatusBadRequest,
			wantCode:   "validation_error",
		},
		{
			name:       "404 missing project",
			body:       `{"question":"¿Existe?"}`,
			service:    services.NewProjectAssistant(&stubAssistantHandlerRepo{err: errors.New("no rows in result set")}, &stubAssistantHandlerRetriever{}, &stubAssistantHandlerProvider{}),
			wantStatus: http.StatusNotFound,
			wantCode:   "not_found",
		},
		{
			name:       "409 assistant unavailable",
			body:       `{"question":"¿Existe contexto?"}`,
			service:    services.NewProjectAssistant(&stubAssistantHandlerRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true}}, &stubAssistantHandlerRetriever{}, &stubAssistantHandlerProvider{}),
			wantStatus: http.StatusConflict,
			wantCode:   "assistant_unavailable",
		},
		{
			name:       "502 upstream failure",
			body:       `{"question":"¿Qué tradeoffs hubo?"}`,
			service:    services.NewProjectAssistant(&stubAssistantHandlerRepo{context: projectContext}, &stubAssistantHandlerRetriever{err: errors.New("upstream down")}, &stubAssistantHandlerProvider{}),
			wantStatus: http.StatusBadGateway,
			wantCode:   "assistant_upstream_error",
		},
		{
			name:       "400 invalid json",
			body:       `{`,
			service:    services.NewProjectAssistant(&stubAssistantHandlerRepo{}, &stubAssistantHandlerRetriever{}, &stubAssistantHandlerProvider{}),
			wantStatus: http.StatusBadRequest,
			wantCode:   "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performAssistantRequest(t, NewProjectAssistantHandler(tt.service), tt.body, model.User{ID: uuid.New(), AuthProvider: "google", EmailVerified: true, FullName: "Ada Lovelace", Company: "Analytical Engines", ProfileCompleted: true, AssistantEligible: true, CanUseProjectAssistant: true})
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var payload model.APIErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if payload.Error.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", payload.Error.Code, tt.wantCode)
			}
		})
	}
}

func TestPerformAssistantRequestTargetsPrivateAssistantRoute(t *testing.T) {
	provider := &stubAssistantHandlerProvider{resp: "ok"}
	handler := NewProjectAssistantHandler(services.NewProjectAssistant(
		&stubAssistantHandlerRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}},
		&stubAssistantHandlerRetriever{chunks: []services.MarkdownChunkAlias{{Heading: "Architecture", Body: "Uses Go."}}},
		provider,
	))

	rec := performAssistantRequest(t, handler, `{"question":"How does it work?"}`, model.User{ID: uuid.New(), AuthProvider: "google", EmailVerified: true, FullName: "Ada Lovelace", Company: "Analytical Engines", ProfileCompleted: true, AssistantEligible: true, CanUseProjectAssistant: true})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func performAssistantRequest(t *testing.T, handler *ProjectAssistantHandler, body string, currentUser model.User) *httptest.ResponseRecorder {
	t.Helper()
	privatePath := "/api/v1/private/projects/portfolioforge/assistant/messages"
	req := httptest.NewRequest(http.MethodPost, privatePath, bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) { c.Error(err) }
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/private/projects/:slug/assistant/messages")
	c.SetParamNames("slug")
	c.SetParamValues("portfolioforge")
	if req.URL.Path != privatePath {
		t.Fatalf("request path = %q, want %q", req.URL.Path, privatePath)
	}
	c.Set("currentUser", currentUser)

	if err := handler.CreateMessage(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	return rec
}
