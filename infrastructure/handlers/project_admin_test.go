package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubProjectAdminRow struct {
	scan func(dest ...interface{}) error
}

func (s stubProjectAdminRow) Scan(dest ...interface{}) error {
	return s.scan(dest...)
}

type stubProjectAdminTx struct {
	queryRow   func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	execs      []string
	committed  bool
	rolledBack bool
}

func (s *stubProjectAdminTx) Exec(_ context.Context, sql string, _ ...interface{}) (pgconn.CommandTag, error) {
	s.execs = append(s.execs, sql)
	return pgconn.NewCommandTag("OK"), nil
}

func (s *stubProjectAdminTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return s.queryRow(ctx, sql, args...)
}

func (s *stubProjectAdminTx) Commit(_ context.Context) error {
	s.committed = true
	return nil
}

func (s *stubProjectAdminTx) Rollback(_ context.Context) error {
	s.rolledBack = true
	return nil
}

type stubEmbeddingProvider struct {
	text string
	err  error
}

func (s *stubEmbeddingProvider) Generate(_ context.Context, text string) ([]float32, error) {
	s.text = text
	if s.err != nil {
		return nil, s.err
	}
	return []float32{0.1, 0.2}, nil
}

func (s *stubEmbeddingProvider) Dimension() int { return 1536 }

func TestProjectAdminHandler_UpdateProjectEnrichmentFailsWhenEmbeddingGenerationFails(t *testing.T) {
	projectID := uuid.New()
	rawSearchText := "PortfolioForge semantic search búsqueda semántica inconsistente Go PostgreSQL"
	tx := &stubProjectAdminTx{
		queryRow: func(_ context.Context, sql string, _ ...interface{}) pgx.Row {
			switch {
			case strings.Contains(sql, "compose_project_embedding_text"):
				return stubProjectAdminRow{scan: func(dest ...interface{}) error {
					*(dest[0].(*string)) = rawSearchText
					return nil
				}}
			case strings.Contains(sql, "SELECT search_content_hash"):
				return stubProjectAdminRow{scan: func(dest ...interface{}) error {
					return pgx.ErrNoRows
				}}
			default:
				return stubProjectAdminRow{scan: func(dest ...interface{}) error {
					return errors.New("unexpected query")
				}}
			}
		},
	}
	embeddingProv := &stubEmbeddingProvider{err: errors.New("openai unavailable")}
	handler := &ProjectAdminHandler{
		beginTx: func(context.Context) (projectAdminTx, error) {
			return tx, nil
		},
		embeddingProv:   embeddingProv,
		semanticEnabled: true,
	}

	body := `{"profile":{"business_goal":"mejorar conversión","problem_statement":"búsqueda semántica inconsistente","solution_summary":"normaliza evidencia","architecture":"go api","ai_usage":"embeddings"},"technology_ids":["` + uuid.New().String() + `"]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/projects/"+projectID.String()+"/enrichment", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/projects/:id/enrichment")
	c.SetParamNames("id")
	c.SetParamValues(projectID.String())

	err := handler.UpdateProjectEnrichment(c)
	if err == nil {
		t.Fatal("expected error when embedding generation fails")
	}

	contractErr, ok := err.(*model.ContractError)
	if !ok {
		t.Fatalf("error type = %T, want *model.ContractError", err)
	}
	if contractErr.StatusHTTP != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", contractErr.StatusHTTP, http.StatusInternalServerError)
	}
	if tx.committed {
		t.Fatal("transaction committed despite embedding failure")
	}
	if !tx.rolledBack {
		t.Fatal("transaction rollback was not triggered")
	}
	if embeddingProv.text != rawSearchText {
		t.Fatalf("embedding input = %q, want raw text %q", embeddingProv.text, rawSearchText)
	}
	for _, sql := range tx.execs {
		if strings.Contains(sql, "search_embedding") {
			t.Fatalf("search embedding update executed despite generation failure: %q", sql)
		}
	}
}

func TestProjectAdminHandler_UpdateProjectEnrichmentSucceedsAndPersistsReembedding(t *testing.T) {
	projectID := uuid.New()
	technologyID := uuid.New()
	problemStatement := "búsqueda semántica inconsistente"
	rawSearchText := "PortfolioForge mejora la búsqueda semántica inconsistente con Go, PostgreSQL y embeddings"
	tx := &stubProjectAdminTx{
		queryRow: func(_ context.Context, sql string, _ ...interface{}) pgx.Row {
			switch {
			case strings.Contains(sql, "compose_project_embedding_text"):
				return stubProjectAdminRow{scan: func(dest ...interface{}) error {
					*(dest[0].(*string)) = rawSearchText
					return nil
				}}
			case strings.Contains(sql, "SELECT search_content_hash"):
				return stubProjectAdminRow{scan: func(dest ...interface{}) error {
					return pgx.ErrNoRows
				}}
			default:
				return stubProjectAdminRow{scan: func(dest ...interface{}) error {
					return errors.New("unexpected query")
				}}
			}
		},
	}
	embeddingProv := &stubEmbeddingProvider{}
	handler := &ProjectAdminHandler{
		beginTx: func(context.Context) (projectAdminTx, error) {
			return tx, nil
		},
		embeddingProv:   embeddingProv,
		semanticEnabled: true,
	}

	body := `{"profile":{"business_goal":"mejorar conversión","problem_statement":"` + problemStatement + `","solution_summary":"normaliza evidencia","architecture":"go api","ai_usage":"embeddings"},"technology_ids":["` + technologyID.String() + `"]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/projects/"+projectID.String()+"/enrichment", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/admin/projects/:id/enrichment")
	c.SetParamNames("id")
	c.SetParamValues(projectID.String())

	err := handler.UpdateProjectEnrichment(c)
	if err != nil {
		t.Fatalf("UpdateProjectEnrichment() error = %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !tx.committed {
		t.Fatal("transaction was not committed")
	}
	if embeddingProv.text != rawSearchText {
		t.Fatalf("embedding input = %q, want raw text %q", embeddingProv.text, rawSearchText)
	}
	if !strings.Contains(embeddingProv.text, problemStatement) {
		t.Fatalf("embedding input %q does not include problem_statement %q", embeddingProv.text, problemStatement)
	}

	assertExecContains(t, tx.execs, []string{
		"INSERT INTO project_profiles",
		"DELETE FROM project_technologies",
		"INSERT INTO project_technologies",
		"INSERT INTO project_search_documents",
		"UPDATE project_search_documents SET search_embedding",
	})
}

func assertExecContains(t *testing.T, execs []string, fragments []string) {
	t.Helper()
	for _, fragment := range fragments {
		found := false
		for _, sql := range execs {
			if strings.Contains(sql, fragment) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected executed SQL containing %q, got %#v", fragment, execs)
		}
	}
}
