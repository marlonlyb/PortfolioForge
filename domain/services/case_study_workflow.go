package services

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	searchPorts "github.com/marlonlyb/portfolioforge/domain/ports/search"
	workflowPorts "github.com/marlonlyb/portfolioforge/domain/ports/workflow"
	"github.com/marlonlyb/portfolioforge/model"
)

type CaseStudyPublisher interface {
	ResolvePublishTarget(inputPath, explicitSlug string) (CaseStudyPublishTarget, error)
	CollectFiles(localDir string) ([]string, error)
	Publish(ctx context.Context, target CaseStudyPublishTarget, files []string) error
	Verify(ctx context.Context, publicURL string) error
}

type CaseStudyPublishTarget struct {
	Slug      string
	LocalDir  string
	LocalFile string
	PublicURL string
}

type CaseStudyProjectImporter interface {
	ImportFromCanonical(ctx context.Context, source CaseStudyPublishTarget, canonicalURL string) (uuid.UUID, error)
}

type CaseStudyLocalizationBackfiller interface {
	BackfillProject(ctx context.Context, projectID uuid.UUID, locales []string) error
}

type CaseStudyWorkflowConfig struct {
	AllowedSourceRoots []string
}

type CaseStudyWorkflowService struct {
	repository   workflowPorts.Repository
	publisher    CaseStudyPublisher
	importer     CaseStudyProjectImporter
	localization CaseStudyLocalizationBackfiller
	searchRepo   searchPorts.SearchRepository
	allowedRoots []string
	now          func() time.Time
}

func NewCaseStudyWorkflowService(
	repository workflowPorts.Repository,
	publisher CaseStudyPublisher,
	importer CaseStudyProjectImporter,
	localization CaseStudyLocalizationBackfiller,
	searchRepo searchPorts.SearchRepository,
	config CaseStudyWorkflowConfig,
) *CaseStudyWorkflowService {
	allowedRoots := make([]string, 0, len(config.AllowedSourceRoots))
	for _, root := range config.AllowedSourceRoots {
		trimmed := strings.TrimSpace(root)
		if trimmed == "" {
			continue
		}
		allowedRoots = append(allowedRoots, trimmed)
	}

	return &CaseStudyWorkflowService{
		repository:   repository,
		publisher:    publisher,
		importer:     importer,
		localization: localization,
		searchRepo:   searchRepo,
		allowedRoots: allowedRoots,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (s *CaseStudyWorkflowService) StartRun(ctx context.Context, req model.StartCaseStudyWorkflowRunRequest) (model.CaseStudyWorkflowRun, error) {
	req.Normalize()
	if strings.TrimSpace(req.SourcePath) == "" {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("source_path es obligatorio")
	}
	if len(s.allowedRoots) == 0 {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("no hay raíces allowlist configuradas para el workflow")
	}

	source, publishTarget, err := s.resolveSource(req.SourcePath, req.Slug)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}

	run := model.NewCaseStudyWorkflowRun(source, req.ResolveOptions())
	resolveStep := run.StepByName(model.CaseStudyWorkflowStepResolveSource)
	now := s.now()
	resolveStep.Status = model.CaseStudyWorkflowStatusSucceeded
	resolveStep.StartedAt = &now
	resolveStep.FinishedAt = &now
	resolveStep.AttemptCount = 1
	resolveStep.Output = mustMarshalJSON(map[string]any{
		"canonical_directory":     source.CanonicalDirectory,
		"canonical_markdown_path": source.CanonicalMarkdownPath,
		"slug":                    source.Slug,
	})

	s.publishNextReadySteps(&run)
	run.RecomputeStatus()
	if err := s.repository.SaveRun(ctx, run); err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	if err := s.repository.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{
		RunID:     run.ID,
		Step:      model.CaseStudyWorkflowStepResolveSource,
		Level:     model.CaseStudyWorkflowLogInfo,
		Message:   fmt.Sprintf("Canonical source resolved: %s", publishTarget.LocalDir),
		CreatedAt: now,
	}); err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}

	return s.repository.GetRun(ctx, run.ID)
}

func (s *CaseStudyWorkflowService) GetRun(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	return s.repository.GetRun(ctx, runID)
}

func (s *CaseStudyWorkflowService) ListLogs(ctx context.Context, runID uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error) {
	return s.repository.ListLogs(ctx, runID)
}

func (s *CaseStudyWorkflowService) ConfirmStep(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error) {
	run, err := s.repository.GetRun(ctx, runID)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}

	step := run.StepByName(stepName)
	if step == nil {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("step %s no existe", stepName)
	}
	if !step.RequiresConfirmation {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("step %s no requiere confirmación", stepName)
	}
	now := s.now()
	step.ConfirmationGrantedAt = &now
	if step.Status == model.CaseStudyWorkflowStatusAwaitingConfirmation || step.Status == model.CaseStudyWorkflowStatusBlocked {
		step.Status = model.CaseStudyWorkflowStatusPending
	}
	s.publishNextReadySteps(&run)
	run.RecomputeStatus()
	if err := s.repository.SaveRun(ctx, run); err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	if err := s.repository.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{RunID: run.ID, Step: stepName, Level: model.CaseStudyWorkflowLogWarn, Message: "Operator confirmed the step.", CreatedAt: now}); err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	return s.repository.GetRun(ctx, runID)
}

func (s *CaseStudyWorkflowService) StartStep(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error) {
	run, err := s.repository.GetRun(ctx, runID)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	if err := s.executeStep(ctx, &run, stepName, false); err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	return s.repository.GetRun(ctx, runID)
}

func (s *CaseStudyWorkflowService) RetryStep(ctx context.Context, runID uuid.UUID, stepName string) (model.CaseStudyWorkflowRun, error) {
	run, err := s.repository.GetRun(ctx, runID)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	step := run.StepByName(stepName)
	if step == nil {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("step %s no existe", stepName)
	}
	if step.Status != model.CaseStudyWorkflowStatusFailed {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("solo se puede reintentar un step fallido")
	}
	if err := s.executeStep(ctx, &run, stepName, true); err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}
	return s.repository.GetRun(ctx, runID)
}

func (s *CaseStudyWorkflowService) Resume(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	run, err := s.repository.GetRun(ctx, runID)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, err
	}

	for _, stepName := range model.CaseStudyWorkflowOrderedSteps {
		step := run.StepByName(stepName)
		if step == nil {
			continue
		}
		if step.Status == model.CaseStudyWorkflowStatusFailed {
			if err := s.executeStep(ctx, &run, stepName, true); err != nil {
				return model.CaseStudyWorkflowRun{}, err
			}
			return s.repository.GetRun(ctx, runID)
		}
		if s.canStartStep(run, stepName) {
			if err := s.executeStep(ctx, &run, stepName, false); err != nil {
				return model.CaseStudyWorkflowRun{}, err
			}
		}
		current := run.StepByName(stepName)
		if current != nil && (current.Status == model.CaseStudyWorkflowStatusAwaitingConfirmation || current.Status == model.CaseStudyWorkflowStatusBlocked || current.Status == model.CaseStudyWorkflowStatusFailed) {
			break
		}
	}

	return s.repository.GetRun(ctx, runID)
}

func (s *CaseStudyWorkflowService) executeStep(ctx context.Context, run *model.CaseStudyWorkflowRun, stepName string, retry bool) error {
	step := run.StepByName(stepName)
	if step == nil {
		return fmt.Errorf("step %s no existe", stepName)
	}
	if !retry && !s.canStartStep(*run, stepName) {
		return fmt.Errorf("step %s todavía no está listo para ejecutarse", stepName)
	}
	if retry && step.Status != model.CaseStudyWorkflowStatusFailed {
		return fmt.Errorf("step %s no es reintentable", stepName)
	}

	now := s.now()
	step.Status = model.CaseStudyWorkflowStatusRunning
	step.StartedAt = &now
	step.FinishedAt = nil
	step.AttemptCount++
	step.ErrorMessage = ""
	run.RecomputeStatus()
	if err := s.repository.SaveRun(ctx, *run); err != nil {
		return err
	}
	if err := s.repository.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{RunID: run.ID, Step: stepName, Level: model.CaseStudyWorkflowLogInfo, Message: "Step started.", CreatedAt: now}); err != nil {
		return err
	}

	var execErr error
	switch stepName {
	case model.CaseStudyWorkflowStepPublishCanonical:
		execErr = s.executePublishStep(ctx, run)
	case model.CaseStudyWorkflowStepImportProject:
		execErr = s.executeImportStep(ctx, run)
	case model.CaseStudyWorkflowStepLocalization:
		execErr = s.executeLocalizationStep(ctx, run)
	case model.CaseStudyWorkflowStepReembed:
		execErr = s.executeReembedStep(ctx, run)
	default:
		execErr = fmt.Errorf("step %s no se puede ejecutar manualmente", stepName)
	}

	finishedAt := s.now()
	step = run.StepByName(stepName)
	step.FinishedAt = &finishedAt
	if execErr != nil {
		step.Status = model.CaseStudyWorkflowStatusFailed
		step.ErrorMessage = execErr.Error()
		run.RecomputeStatus()
		if err := s.repository.SaveRun(ctx, *run); err != nil {
			return err
		}
		if err := s.repository.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{RunID: run.ID, Step: stepName, Level: model.CaseStudyWorkflowLogError, Message: execErr.Error(), CreatedAt: finishedAt}); err != nil {
			return err
		}
		return execErr
	}

	step.Status = model.CaseStudyWorkflowStatusSucceeded
	step.ErrorMessage = ""
	s.publishNextReadySteps(run)
	run.RecomputeStatus()
	if err := s.repository.SaveRun(ctx, *run); err != nil {
		return err
	}
	if err := s.repository.AppendLog(ctx, model.CaseStudyWorkflowLogEntry{RunID: run.ID, Step: stepName, Level: model.CaseStudyWorkflowLogInfo, Message: "Step completed successfully.", CreatedAt: finishedAt}); err != nil {
		return err
	}
	return nil
}

func (s *CaseStudyWorkflowService) executePublishStep(ctx context.Context, run *model.CaseStudyWorkflowRun) error {
	target, err := s.publisher.ResolvePublishTarget(run.Source.NormalizedPath, run.Source.Slug)
	if err != nil {
		return err
	}
	files, err := s.publisher.CollectFiles(target.LocalDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no se encontraron archivos para publicar en %s", target.LocalDir)
	}
	if err := s.publisher.Publish(ctx, target, files); err != nil {
		return err
	}
	if err := s.publisher.Verify(ctx, target.PublicURL); err != nil {
		return err
	}
	run.CanonicalURL = target.PublicURL
	run.Source.CanonicalDirectory = target.LocalDir
	run.Source.CanonicalMarkdownPath = target.LocalFile
	run.Source.Slug = target.Slug
	run.StepByName(model.CaseStudyWorkflowStepPublishCanonical).Output = mustMarshalJSON(map[string]any{
		"canonical_url": target.PublicURL,
		"files":         len(files),
	})
	return nil
}

func (s *CaseStudyWorkflowService) executeImportStep(ctx context.Context, run *model.CaseStudyWorkflowRun) error {
	if strings.TrimSpace(run.CanonicalURL) == "" {
		return fmt.Errorf("publish_canonical debe completarse antes de importar")
	}
	target, err := s.publisher.ResolvePublishTarget(run.Source.NormalizedPath, run.Source.Slug)
	if err != nil {
		return err
	}
	projectID, err := s.importer.ImportFromCanonical(ctx, target, run.CanonicalURL)
	if err != nil {
		return err
	}
	run.ProjectID = &projectID
	run.StepByName(model.CaseStudyWorkflowStepImportProject).Output = mustMarshalJSON(map[string]any{
		"project_id": projectID.String(),
	})
	return nil
}

func (s *CaseStudyWorkflowService) executeLocalizationStep(ctx context.Context, run *model.CaseStudyWorkflowRun) error {
	if run.ProjectID == nil {
		return fmt.Errorf("import_or_update_project debe completarse antes de localization_backfill")
	}
	if err := s.localization.BackfillProject(ctx, *run.ProjectID, run.Options.Locales); err != nil {
		return err
	}
	run.StepByName(model.CaseStudyWorkflowStepLocalization).Output = mustMarshalJSON(map[string]any{
		"project_id": run.ProjectID.String(),
		"locales":    run.Options.Locales,
	})
	return nil
}

func (s *CaseStudyWorkflowService) executeReembedStep(ctx context.Context, run *model.CaseStudyWorkflowRun) error {
	if run.ProjectID == nil {
		return fmt.Errorf("import_or_update_project debe completarse antes de reembed")
	}
	if err := s.searchRepo.RefreshSearchDocument(ctx, *run.ProjectID); err != nil {
		return err
	}
	run.StepByName(model.CaseStudyWorkflowStepReembed).Output = mustMarshalJSON(map[string]any{
		"project_id": run.ProjectID.String(),
	})
	return nil
}

func (s *CaseStudyWorkflowService) resolveSource(inputPath, explicitSlug string) (model.CaseStudyWorkflowSource, CaseStudyPublishTarget, error) {
	normalizedPath, allowedRoot, err := s.normalizeWithinAllowlist(inputPath)
	if err != nil {
		return model.CaseStudyWorkflowSource{}, CaseStudyPublishTarget{}, err
	}
	target, err := s.publisher.ResolvePublishTarget(normalizedPath, explicitSlug)
	if err != nil {
		return model.CaseStudyWorkflowSource{}, CaseStudyPublishTarget{}, err
	}
	return model.CaseStudyWorkflowSource{
		AllowedRoot:           allowedRoot,
		RequestedPath:         inputPath,
		NormalizedPath:        normalizedPath,
		CanonicalRootPath:     filepath.Join(filepath.Dir(target.LocalDir), ".."),
		CanonicalDirectory:    target.LocalDir,
		CanonicalMarkdownPath: target.LocalFile,
		Slug:                  target.Slug,
	}, target, nil
}

func (s *CaseStudyWorkflowService) normalizeWithinAllowlist(inputPath string) (string, string, error) {
	cleaned := filepath.Clean(strings.TrimSpace(inputPath))
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", "", fmt.Errorf("no se pudo normalizar source_path: %w", err)
	}
	for _, root := range s.allowedRoots {
		rootAbs, err := filepath.Abs(filepath.Clean(root))
		if err != nil {
			continue
		}
		rel, err := filepath.Rel(rootAbs, abs)
		if err != nil {
			continue
		}
		if rel == "." || (!strings.HasPrefix(rel, "..") && rel != "") {
			return abs, rootAbs, nil
		}
	}
	return "", "", fmt.Errorf("la ruta %s está fuera de las raíces allowlist configuradas", abs)
}

func (s *CaseStudyWorkflowService) publishNextReadySteps(run *model.CaseStudyWorkflowRun) {
	resolve := run.StepByName(model.CaseStudyWorkflowStepResolveSource)
	if resolve != nil && resolve.Status == model.CaseStudyWorkflowStatusSucceeded {
		publish := run.StepByName(model.CaseStudyWorkflowStepPublishCanonical)
		if publish != nil && publish.Status == model.CaseStudyWorkflowStatusBlocked {
			if publish.RequiresConfirmation && publish.ConfirmationGrantedAt == nil {
				publish.Status = model.CaseStudyWorkflowStatusAwaitingConfirmation
			} else {
				publish.Status = model.CaseStudyWorkflowStatusPending
			}
		}
	}

	publish := run.StepByName(model.CaseStudyWorkflowStepPublishCanonical)
	if publish != nil && publish.Status == model.CaseStudyWorkflowStatusSucceeded {
		importStep := run.StepByName(model.CaseStudyWorkflowStepImportProject)
		if importStep != nil && importStep.Status == model.CaseStudyWorkflowStatusBlocked {
			if importStep.RequiresConfirmation && importStep.ConfirmationGrantedAt == nil {
				importStep.Status = model.CaseStudyWorkflowStatusAwaitingConfirmation
			} else {
				importStep.Status = model.CaseStudyWorkflowStatusPending
			}
		}
	}

	importStep := run.StepByName(model.CaseStudyWorkflowStepImportProject)
	if importStep != nil && importStep.Status == model.CaseStudyWorkflowStatusSucceeded {
		localization := run.StepByName(model.CaseStudyWorkflowStepLocalization)
		if localization != nil && localization.Status == model.CaseStudyWorkflowStatusBlocked {
			localization.Status = model.CaseStudyWorkflowStatusPending
		}
		if localization == nil || localization.Status == model.CaseStudyWorkflowStatusSucceeded || localization.Status == model.CaseStudyWorkflowStatusSkipped {
			reembed := run.StepByName(model.CaseStudyWorkflowStepReembed)
			if reembed != nil && reembed.Status == model.CaseStudyWorkflowStatusBlocked {
				reembed.Status = model.CaseStudyWorkflowStatusPending
			}
		}
	}

	localization := run.StepByName(model.CaseStudyWorkflowStepLocalization)
	if localization != nil && (localization.Status == model.CaseStudyWorkflowStatusSucceeded || localization.Status == model.CaseStudyWorkflowStatusSkipped) {
		reembed := run.StepByName(model.CaseStudyWorkflowStepReembed)
		if reembed != nil && reembed.Status == model.CaseStudyWorkflowStatusBlocked {
			reembed.Status = model.CaseStudyWorkflowStatusPending
		}
	}
}

func (s *CaseStudyWorkflowService) canStartStep(run model.CaseStudyWorkflowRun, stepName string) bool {
	step := run.StepByName(stepName)
	if step == nil {
		return false
	}
	if step.Status != model.CaseStudyWorkflowStatusPending {
		return false
	}
	if step.RequiresConfirmation && step.ConfirmationGrantedAt == nil {
		return false
	}
	for _, candidate := range model.CaseStudyWorkflowOrderedSteps {
		if candidate == stepName {
			return true
		}
		previous := run.StepByName(candidate)
		if previous == nil {
			continue
		}
		if previous.Status == model.CaseStudyWorkflowStatusSkipped {
			continue
		}
		if previous.Status != model.CaseStudyWorkflowStatusSucceeded {
			return false
		}
	}
	return true
}

func mustMarshalJSON(value interface{}) json.RawMessage {
	encoded, _ := json.Marshal(value)
	return encoded
}
