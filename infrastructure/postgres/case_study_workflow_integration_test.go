package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

func TestCaseStudyWorkflowRepository_RoundTripWithRealPostgres(t *testing.T) {
	ctx := context.Background()
	pool := newWorkflowTestDatabase(t)
	repo := NewCaseStudyWorkflowRepository(pool)
	now := time.Date(2026, 4, 18, 9, 0, 0, 0, time.UTC)
	projectID := uuid.New()
	runID := uuid.New()

	if _, err := pool.Exec(ctx, `INSERT INTO products (id) VALUES ($1)`, projectID); err != nil {
		t.Fatalf("seed products: %v", err)
	}

	run := model.CaseStudyWorkflowRun{
		ID:           runID,
		Status:       model.CaseStudyWorkflowStatusFailed,
		CanonicalURL: "https://example.com/demo/demo.md",
		ProjectID:    &projectID,
		LastError:    "import failed",
		Source: model.CaseStudyWorkflowSource{
			AllowedRoot:           "/safe/root",
			RequestedPath:         "/safe/root/90. dev_portfolioforge/demo",
			NormalizedPath:        "/safe/root/90. dev_portfolioforge/demo",
			CanonicalRootPath:     "/safe/root/90. dev_portfolioforge",
			CanonicalDirectory:    "/safe/root/90. dev_portfolioforge/demo",
			CanonicalMarkdownPath: "/safe/root/90. dev_portfolioforge/demo/demo.md",
			Slug:                  "demo",
		},
		Options: model.CaseStudyWorkflowOptions{
			RunLocalizationBackfill: true,
			RunReembed:              true,
			Locales:                 []string{"ca", "en"},
		},
		GenerationScope: model.CaseStudyWorkflowScopeUI{
			CanonicalGenerationAvailable: false,
			CanonicalGenerationMessage:   "MVP starts from an existing canonical source.",
		},
		CreatedAt: now,
		UpdatedAt: now.Add(2 * time.Minute),
		Steps: []model.CaseStudyWorkflowStep{
			{
				RunID:        runID,
				Step:         model.CaseStudyWorkflowStepImportProject,
				Status:       model.CaseStudyWorkflowStatusFailed,
				AttemptCount: 2,
				StartedAt:    pointerToTime(now.Add(90 * time.Second)),
				FinishedAt:   pointerToTime(now.Add(2 * time.Minute)),
				ErrorMessage: "import failed",
			},
			{
				RunID:                 runID,
				Step:                  model.CaseStudyWorkflowStepPublishCanonical,
				Status:                model.CaseStudyWorkflowStatusSucceeded,
				RequiresConfirmation:  true,
				ConfirmationGrantedAt: pointerToTime(now.Add(30 * time.Second)),
				StartedAt:             pointerToTime(now.Add(30 * time.Second)),
				FinishedAt:            pointerToTime(now.Add(time.Minute)),
				AttemptCount:          1,
				Output:                mustJSON(t, map[string]any{"canonical_url": "https://example.com/demo/demo.md", "files": 3}),
			},
			{
				RunID:        runID,
				Step:         model.CaseStudyWorkflowStepResolveSource,
				Status:       model.CaseStudyWorkflowStatusSucceeded,
				AttemptCount: 1,
				StartedAt:    pointerToTime(now),
				FinishedAt:   pointerToTime(now),
				Output:       mustJSON(t, map[string]any{"slug": "demo"}),
			},
		},
	}

	if err := repo.SaveRun(ctx, run); err != nil {
		t.Fatalf("SaveRun() error = %v", err)
	}
	if err := repo.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{RunID: runID, Step: model.CaseStudyWorkflowStepPublishCanonical, Level: model.CaseStudyWorkflowLogInfo, Message: "Published canonical files.", CreatedAt: now.Add(time.Minute)}); err != nil {
		t.Fatalf("AppendLog(publish) error = %v", err)
	}
	if err := repo.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{RunID: runID, Step: model.CaseStudyWorkflowStepImportProject, Level: model.CaseStudyWorkflowLogError, Message: "Import failed.", CreatedAt: now.Add(2 * time.Minute)}); err != nil {
		t.Fatalf("AppendLog(import) error = %v", err)
	}

	loaded, err := repo.GetRun(ctx, runID)
	if err != nil {
		t.Fatalf("GetRun() error = %v", err)
	}
	if loaded.CanonicalURL != run.CanonicalURL {
		t.Fatalf("CanonicalURL = %q, want %q", loaded.CanonicalURL, run.CanonicalURL)
	}
	if loaded.ProjectID == nil || *loaded.ProjectID != projectID {
		t.Fatalf("ProjectID = %v, want %s", loaded.ProjectID, projectID)
	}
	if len(loaded.Steps) != 3 {
		t.Fatalf("len(loaded.Steps) = %d, want 3", len(loaded.Steps))
	}
	if loaded.Steps[0].Step != model.CaseStudyWorkflowStepResolveSource || loaded.Steps[1].Step != model.CaseStudyWorkflowStepPublishCanonical || loaded.Steps[2].Step != model.CaseStudyWorkflowStepImportProject {
		t.Fatalf("loaded steps order = %#v", loaded.Steps)
	}
	var publishOutput map[string]any
	if err := json.Unmarshal(loaded.StepByName(model.CaseStudyWorkflowStepPublishCanonical).Output, &publishOutput); err != nil {
		t.Fatalf("unmarshal publish output: %v", err)
	}
	if publishOutput["canonical_url"] != "https://example.com/demo/demo.md" || int(publishOutput["files"].(float64)) != 3 {
		t.Fatalf("publish output = %#v", publishOutput)
	}

	logs, err := repo.ListLogs(ctx, runID)
	if err != nil {
		t.Fatalf("ListLogs() error = %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("len(logs) = %d, want 2", len(logs))
	}
	if logs[0].Message != "Published canonical files." || logs[1].Message != "Import failed." {
		t.Fatalf("logs = %#v", logs)
	}
}

func newWorkflowTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()
	port := freePort(t)
	baseDir := t.TempDir()
	db := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Version(embeddedpostgres.V16).
			Port(uint32(port)).
			Database("workflow_test").
			Username("postgres").
			Password("postgres").
			RuntimePath(filepath.Join(baseDir, "runtime")).
			DataPath(filepath.Join(baseDir, "data")).
			BinariesPath(filepath.Join(baseDir, "binaries")),
	)
	if err := db.Start(); err != nil {
		t.Fatalf("start embedded postgres: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Stop(); err != nil {
			t.Fatalf("stop embedded postgres: %v", err)
		}
	})

	dsn := fmt.Sprintf("postgres://postgres:postgres@127.0.0.1:%d/workflow_test?sslmode=disable", port)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	t.Cleanup(pool.Close)

	applyWorkflowSchema(t, pool)
	return pool
}

func applyWorkflowSchema(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	migrationPath := filepath.Join("..", "..", "sqlmigrations", "20260418_0900_case_study_workflow_runs.sql")
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("read workflow migration: %v", err)
	}
	statements := []string{
		`CREATE TABLE IF NOT EXISTS products (id UUID PRIMARY KEY)`,
		string(migration),
	}
	for _, statement := range statements {
		if _, err := pool.Exec(context.Background(), statement); err != nil {
			t.Fatalf("apply schema statement %q: %v", statement, err)
		}
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve free port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func pointerToTime(value time.Time) *time.Time {
	return &value
}

func mustJSON(t *testing.T, value map[string]any) []byte {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	return raw
}
