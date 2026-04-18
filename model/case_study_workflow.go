package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	CaseStudyWorkflowStepResolveSource      = "resolve_source"
	CaseStudyWorkflowStepPublishCanonical   = "publish_canonical"
	CaseStudyWorkflowStepImportProject      = "import_or_update_project"
	CaseStudyWorkflowStepLocalization       = "localization_backfill"
	CaseStudyWorkflowStepReembed            = "reembed"
	CaseStudyWorkflowStepGenerationDeferred = "generate_canonical"
)

const (
	CaseStudyWorkflowStatusPending              = "pending"
	CaseStudyWorkflowStatusBlocked              = "blocked"
	CaseStudyWorkflowStatusAwaitingConfirmation = "awaiting_confirmation"
	CaseStudyWorkflowStatusRunning              = "running"
	CaseStudyWorkflowStatusSucceeded            = "succeeded"
	CaseStudyWorkflowStatusFailed               = "failed"
	CaseStudyWorkflowStatusSkipped              = "skipped"
)

const (
	CaseStudyWorkflowLogInfo  = "info"
	CaseStudyWorkflowLogWarn  = "warn"
	CaseStudyWorkflowLogError = "error"
)

var CaseStudyWorkflowOrderedSteps = []string{
	CaseStudyWorkflowStepResolveSource,
	CaseStudyWorkflowStepPublishCanonical,
	CaseStudyWorkflowStepImportProject,
	CaseStudyWorkflowStepLocalization,
	CaseStudyWorkflowStepReembed,
}

type CaseStudyWorkflowOptions struct {
	RunLocalizationBackfill bool     `json:"run_localization_backfill"`
	RunReembed              bool     `json:"run_reembed"`
	Locales                 []string `json:"locales,omitempty"`
}

type CaseStudyWorkflowSource struct {
	AllowedRoot           string `json:"allowed_root"`
	RequestedPath         string `json:"requested_path"`
	NormalizedPath        string `json:"normalized_path"`
	CanonicalRootPath     string `json:"canonical_root_path"`
	CanonicalDirectory    string `json:"canonical_directory"`
	CanonicalMarkdownPath string `json:"canonical_markdown_path"`
	Slug                  string `json:"slug"`
}

type CaseStudyWorkflowStep struct {
	RunID                 uuid.UUID       `json:"run_id"`
	Step                  string          `json:"step"`
	Status                string          `json:"status"`
	RequiresConfirmation  bool            `json:"requires_confirmation"`
	ConfirmationGrantedAt *time.Time      `json:"confirmation_granted_at,omitempty"`
	StartedAt             *time.Time      `json:"started_at,omitempty"`
	FinishedAt            *time.Time      `json:"finished_at,omitempty"`
	AttemptCount          int             `json:"attempt_count"`
	ErrorMessage          string          `json:"error_message,omitempty"`
	Output                json.RawMessage `json:"output,omitempty"`
}

type CaseStudyWorkflowRun struct {
	ID              uuid.UUID                `json:"id"`
	Status          string                   `json:"status"`
	Source          CaseStudyWorkflowSource  `json:"source"`
	Options         CaseStudyWorkflowOptions `json:"options"`
	CanonicalURL    string                   `json:"canonical_url,omitempty"`
	ProjectID       *uuid.UUID               `json:"project_id,omitempty"`
	Steps           []CaseStudyWorkflowStep  `json:"steps"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	LastError       string                   `json:"last_error,omitempty"`
	GenerationScope CaseStudyWorkflowScopeUI `json:"generation_scope"`
}

type CaseStudyWorkflowScopeUI struct {
	CanonicalGenerationAvailable bool   `json:"canonical_generation_available"`
	CanonicalGenerationMessage   string `json:"canonical_generation_message"`
}

type CaseStudyWorkflowLogEntry struct {
	ID        int64     `json:"id"`
	RunID     uuid.UUID `json:"run_id"`
	Step      string    `json:"step"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type StartCaseStudyWorkflowRunRequest struct {
	SourcePath              string   `json:"source_path"`
	Slug                    string   `json:"slug,omitempty"`
	RunLocalizationBackfill *bool    `json:"run_localization_backfill,omitempty"`
	RunReembed              *bool    `json:"run_reembed,omitempty"`
	Locales                 []string `json:"locales,omitempty"`
}

func (r *StartCaseStudyWorkflowRunRequest) Normalize() {
	r.SourcePath = strings.TrimSpace(r.SourcePath)
	r.Slug = strings.TrimSpace(r.Slug)

	locales := make([]string, 0, len(r.Locales))
	seen := map[string]struct{}{}
	for _, locale := range r.Locales {
		normalized := strings.ToLower(strings.TrimSpace(locale))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		locales = append(locales, normalized)
	}
	r.Locales = locales
}

func (r StartCaseStudyWorkflowRunRequest) ResolveOptions() CaseStudyWorkflowOptions {
	runLocalization := true
	if r.RunLocalizationBackfill != nil {
		runLocalization = *r.RunLocalizationBackfill
	}

	runReembed := true
	if r.RunReembed != nil {
		runReembed = *r.RunReembed
	}

	return CaseStudyWorkflowOptions{
		RunLocalizationBackfill: runLocalization,
		RunReembed:              runReembed,
		Locales:                 append([]string(nil), r.Locales...),
	}
}

func NewCaseStudyWorkflowRun(source CaseStudyWorkflowSource, options CaseStudyWorkflowOptions) CaseStudyWorkflowRun {
	now := time.Now().UTC()
	run := CaseStudyWorkflowRun{
		ID:        uuid.New(),
		Status:    CaseStudyWorkflowStatusPending,
		Source:    source,
		Options:   options,
		CreatedAt: now,
		UpdatedAt: now,
		GenerationScope: CaseStudyWorkflowScopeUI{
			CanonicalGenerationAvailable: false,
			CanonicalGenerationMessage:   "MVP starts from an existing canonical source under 90. dev_portfolioforge/<slug>/. Raw folder generation is not available yet.",
		},
	}

	run.Steps = []CaseStudyWorkflowStep{
		{RunID: run.ID, Step: CaseStudyWorkflowStepResolveSource, Status: CaseStudyWorkflowStatusPending},
		{RunID: run.ID, Step: CaseStudyWorkflowStepPublishCanonical, Status: CaseStudyWorkflowStatusBlocked, RequiresConfirmation: true},
		{RunID: run.ID, Step: CaseStudyWorkflowStepImportProject, Status: CaseStudyWorkflowStatusBlocked, RequiresConfirmation: true},
		{RunID: run.ID, Step: CaseStudyWorkflowStepLocalization, Status: CaseStudyWorkflowStatusBlocked},
		{RunID: run.ID, Step: CaseStudyWorkflowStepReembed, Status: CaseStudyWorkflowStatusBlocked},
	}

	if !options.RunLocalizationBackfill {
		step := run.StepByName(CaseStudyWorkflowStepLocalization)
		step.Status = CaseStudyWorkflowStatusSkipped
	}
	if !options.RunReembed {
		step := run.StepByName(CaseStudyWorkflowStepReembed)
		step.Status = CaseStudyWorkflowStatusSkipped
	}

	return run
}

func (r *CaseStudyWorkflowRun) StepByName(step string) *CaseStudyWorkflowStep {
	for index := range r.Steps {
		if r.Steps[index].Step == step {
			return &r.Steps[index]
		}
	}
	return nil
}

func (r *CaseStudyWorkflowRun) RecomputeStatus() {
	r.UpdatedAt = time.Now().UTC()
	r.LastError = ""
	r.Status = CaseStudyWorkflowStatusPending

	allTerminal := true
	for _, step := range r.Steps {
		switch step.Status {
		case CaseStudyWorkflowStatusFailed:
			r.Status = CaseStudyWorkflowStatusFailed
			r.LastError = step.ErrorMessage
			return
		case CaseStudyWorkflowStatusRunning:
			r.Status = CaseStudyWorkflowStatusRunning
			return
		case CaseStudyWorkflowStatusAwaitingConfirmation:
			r.Status = CaseStudyWorkflowStatusAwaitingConfirmation
			allTerminal = false
		case CaseStudyWorkflowStatusPending, CaseStudyWorkflowStatusBlocked:
			allTerminal = false
		case CaseStudyWorkflowStatusSucceeded, CaseStudyWorkflowStatusSkipped:
		default:
			allTerminal = false
		}
	}

	if allTerminal {
		r.Status = CaseStudyWorkflowStatusSucceeded
	}
}

func IsCaseStudyWorkflowStep(step string) bool {
	for _, candidate := range CaseStudyWorkflowOrderedSteps {
		if candidate == step {
			return true
		}
	}
	return false
}
