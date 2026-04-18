package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type workflowMemoryRepository struct {
	runs map[uuid.UUID]model.CaseStudyWorkflowRun
	logs map[uuid.UUID][]model.CaseStudyWorkflowLogEntry
}

func newWorkflowMemoryRepository() *workflowMemoryRepository {
	return &workflowMemoryRepository{runs: map[uuid.UUID]model.CaseStudyWorkflowRun{}, logs: map[uuid.UUID][]model.CaseStudyWorkflowLogEntry{}}
}

func (r *workflowMemoryRepository) SaveRun(_ context.Context, run model.CaseStudyWorkflowRun) error {
	r.runs[run.ID] = run
	return nil
}

func (r *workflowMemoryRepository) GetRun(_ context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	run, ok := r.runs[runID]
	if !ok {
		return model.CaseStudyWorkflowRun{}, errors.New("run no existe")
	}
	return run, nil
}

func (r *workflowMemoryRepository) ListLogs(_ context.Context, runID uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error) {
	return append([]model.CaseStudyWorkflowLogEntry(nil), r.logs[runID]...), nil
}

func (r *workflowMemoryRepository) AppendLog(_ context.Context, entry model.CaseStudyWorkflowLogEntry) error {
	r.logs[entry.RunID] = append(r.logs[entry.RunID], entry)
	return nil
}

type stubWorkflowPublisher struct {
	target           CaseStudyPublishTarget
	collectFilesResp []string
}

func (s stubWorkflowPublisher) ResolvePublishTarget(string, string) (CaseStudyPublishTarget, error) {
	return s.target, nil
}

func (s stubWorkflowPublisher) CollectFiles(string) ([]string, error) {
	return s.collectFilesResp, nil
}

func (s stubWorkflowPublisher) Publish(context.Context, CaseStudyPublishTarget, []string) error {
	return nil
}

func (s stubWorkflowPublisher) Verify(context.Context, string) error { return nil }

type stubWorkflowImporter struct {
	projectID uuid.UUID
	fails     int
	hits      int
}

func (s *stubWorkflowImporter) ImportFromCanonical(context.Context, CaseStudyPublishTarget, string) (uuid.UUID, error) {
	s.hits++
	if s.hits <= s.fails {
		return uuid.Nil, errors.New("import failed")
	}
	return s.projectID, nil
}

type stubWorkflowLocalization struct{ hits int }

func (s *stubWorkflowLocalization) BackfillProject(context.Context, uuid.UUID, []string) error {
	s.hits++
	return nil
}

type stubWorkflowSearchRepo struct{ hits int }

func (s *stubWorkflowSearchRepo) RefreshSearchDocument(context.Context, uuid.UUID) error {
	s.hits++
	return nil
}
func (s *stubWorkflowSearchRepo) Search(context.Context, model.SearchParams) ([]model.SearchResult, error) {
	return nil, nil
}
func (s *stubWorkflowSearchRepo) RefreshAllDocuments(context.Context) error { return nil }

func TestCaseStudyWorkflowService_StartRunRejectsUnsafePath(t *testing.T) {
	repo := newWorkflowMemoryRepository()
	service := NewCaseStudyWorkflowService(
		repo,
		stubWorkflowPublisher{},
		&stubWorkflowImporter{projectID: uuid.New()},
		&stubWorkflowLocalization{},
		&stubWorkflowSearchRepo{},
		CaseStudyWorkflowConfig{AllowedSourceRoots: []string{"/safe/root"}},
	)

	_, err := service.StartRun(context.Background(), model.StartCaseStudyWorkflowRunRequest{SourcePath: "/tmp/escape"})
	if err == nil {
		t.Fatal("expected allowlist validation error")
	}
}

func TestCaseStudyWorkflowService_RequiresConfirmationBeforePublish(t *testing.T) {
	repo := newWorkflowMemoryRepository()
	service := NewCaseStudyWorkflowService(
		repo,
		stubWorkflowPublisher{target: CaseStudyPublishTarget{Slug: "demo", LocalDir: "/safe/root/90. dev_portfolioforge/demo", LocalFile: "/safe/root/90. dev_portfolioforge/demo/demo.md", PublicURL: "https://example.com/demo/demo.md"}, collectFilesResp: []string{"a.md"}},
		&stubWorkflowImporter{projectID: uuid.New()},
		&stubWorkflowLocalization{},
		&stubWorkflowSearchRepo{},
		CaseStudyWorkflowConfig{AllowedSourceRoots: []string{"/safe/root"}},
	)
	service.now = func() time.Time { return time.Date(2026, 4, 18, 9, 0, 0, 0, time.UTC) }

	run, err := service.StartRun(context.Background(), model.StartCaseStudyWorkflowRunRequest{SourcePath: "/safe/root/90. dev_portfolioforge/demo"})
	if err != nil {
		t.Fatalf("StartRun() error = %v", err)
	}

	publishStep := run.StepByName(model.CaseStudyWorkflowStepPublishCanonical)
	if publishStep == nil || publishStep.Status != model.CaseStudyWorkflowStatusAwaitingConfirmation {
		t.Fatalf("publish status = %#v, want awaiting_confirmation", publishStep)
	}

	if _, err := service.StartStep(context.Background(), run.ID, model.CaseStudyWorkflowStepPublishCanonical); err == nil {
		t.Fatal("expected start publish without confirmation to fail")
	}

	run, err = service.ConfirmStep(context.Background(), run.ID, model.CaseStudyWorkflowStepPublishCanonical)
	if err != nil {
		t.Fatalf("ConfirmStep() error = %v", err)
	}
	run, err = service.StartStep(context.Background(), run.ID, model.CaseStudyWorkflowStepPublishCanonical)
	if err != nil {
		t.Fatalf("StartStep(publish) error = %v", err)
	}

	if run.CanonicalURL == "" {
		t.Fatal("expected canonical_url after publish")
	}
	importStep := run.StepByName(model.CaseStudyWorkflowStepImportProject)
	if importStep == nil || importStep.Status != model.CaseStudyWorkflowStatusAwaitingConfirmation {
		t.Fatalf("import status = %#v, want awaiting_confirmation", importStep)
	}
}

func TestCaseStudyWorkflowService_RetryImportKeepsPublishArtifacts(t *testing.T) {
	repo := newWorkflowMemoryRepository()
	importer := &stubWorkflowImporter{projectID: uuid.New(), fails: 1}
	service := NewCaseStudyWorkflowService(
		repo,
		stubWorkflowPublisher{target: CaseStudyPublishTarget{Slug: "demo", LocalDir: "/safe/root/90. dev_portfolioforge/demo", LocalFile: "/safe/root/90. dev_portfolioforge/demo/demo.md", PublicURL: "https://example.com/demo/demo.md"}, collectFilesResp: []string{"a.md"}},
		importer,
		&stubWorkflowLocalization{},
		&stubWorkflowSearchRepo{},
		CaseStudyWorkflowConfig{AllowedSourceRoots: []string{"/safe/root"}},
	)

	run, err := service.StartRun(context.Background(), model.StartCaseStudyWorkflowRunRequest{SourcePath: "/safe/root/90. dev_portfolioforge/demo"})
	if err != nil {
		t.Fatalf("StartRun() error = %v", err)
	}
	if _, err := service.ConfirmStep(context.Background(), run.ID, model.CaseStudyWorkflowStepPublishCanonical); err != nil {
		t.Fatalf("Confirm publish: %v", err)
	}
	run, err = service.StartStep(context.Background(), run.ID, model.CaseStudyWorkflowStepPublishCanonical)
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if _, err := service.ConfirmStep(context.Background(), run.ID, model.CaseStudyWorkflowStepImportProject); err != nil {
		t.Fatalf("Confirm import: %v", err)
	}
	if _, err := service.StartStep(context.Background(), run.ID, model.CaseStudyWorkflowStepImportProject); err == nil {
		t.Fatal("expected import failure on first attempt")
	}

	run, err = service.GetRun(context.Background(), run.ID)
	if err != nil {
		t.Fatalf("GetRun() error = %v", err)
	}
	if run.CanonicalURL == "" {
		t.Fatal("expected canonical URL to remain available after import failure")
	}
	if run.StepByName(model.CaseStudyWorkflowStepPublishCanonical).Status != model.CaseStudyWorkflowStatusSucceeded {
		t.Fatal("expected publish step to remain succeeded after downstream failure")
	}

	run, err = service.RetryStep(context.Background(), run.ID, model.CaseStudyWorkflowStepImportProject)
	if err != nil {
		t.Fatalf("RetryStep(import) error = %v", err)
	}
	if run.ProjectID == nil || *run.ProjectID != importer.projectID {
		t.Fatalf("project_id = %v, want %s", run.ProjectID, importer.projectID)
	}
}
