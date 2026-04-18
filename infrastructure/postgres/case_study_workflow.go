package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

type CaseStudyWorkflowRepository struct {
	db *pgxpool.Pool
}

func NewCaseStudyWorkflowRepository(db *pgxpool.Pool) *CaseStudyWorkflowRepository {
	return &CaseStudyWorkflowRepository{db: db}
}

func (r *CaseStudyWorkflowRepository) SaveRun(ctx context.Context, run model.CaseStudyWorkflowRun) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin case study workflow tx: %w", err)
	}
	defer tx.Rollback(ctx)

	sourceJSON, err := json.Marshal(run.Source)
	if err != nil {
		return fmt.Errorf("marshal workflow source: %w", err)
	}
	optionsJSON, err := json.Marshal(run.Options)
	if err != nil {
		return fmt.Errorf("marshal workflow options: %w", err)
	}
	scopeJSON, err := json.Marshal(run.GenerationScope)
	if err != nil {
		return fmt.Errorf("marshal workflow scope: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO case_study_workflow_runs (
			id, status, source_payload, options_payload, canonical_url, project_id, last_error, generation_scope_payload, created_at, updated_at
		) VALUES ($1, $2, $3::jsonb, $4::jsonb, $5, $6, $7, $8::jsonb, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			source_payload = EXCLUDED.source_payload,
			options_payload = EXCLUDED.options_payload,
			canonical_url = EXCLUDED.canonical_url,
			project_id = EXCLUDED.project_id,
			last_error = EXCLUDED.last_error,
			generation_scope_payload = EXCLUDED.generation_scope_payload,
			updated_at = EXCLUDED.updated_at`,
		run.ID,
		run.Status,
		sourceJSON,
		optionsJSON,
		NullIfEmpty(run.CanonicalURL),
		run.ProjectID,
		NullIfEmpty(run.LastError),
		scopeJSON,
		run.CreatedAt,
		run.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert case study workflow run: %w", err)
	}

	for _, step := range run.Steps {
		_, err := tx.Exec(ctx, `
			INSERT INTO case_study_workflow_steps (
				run_id, step_name, status, requires_confirmation, confirmation_granted_at, started_at, finished_at, attempt_count, error_message, output_payload
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb)
			ON CONFLICT (run_id, step_name) DO UPDATE SET
				status = EXCLUDED.status,
				requires_confirmation = EXCLUDED.requires_confirmation,
				confirmation_granted_at = EXCLUDED.confirmation_granted_at,
				started_at = EXCLUDED.started_at,
				finished_at = EXCLUDED.finished_at,
				attempt_count = EXCLUDED.attempt_count,
				error_message = EXCLUDED.error_message,
				output_payload = EXCLUDED.output_payload`,
			step.RunID,
			step.Step,
			step.Status,
			step.RequiresConfirmation,
			step.ConfirmationGrantedAt,
			step.StartedAt,
			step.FinishedAt,
			step.AttemptCount,
			NullIfEmpty(step.ErrorMessage),
			normalizeJSONB(step.Output),
		)
		if err != nil {
			return fmt.Errorf("upsert case study workflow step %s: %w", step.Step, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit case study workflow tx: %w", err)
	}
	return nil
}

func (r *CaseStudyWorkflowRepository) GetRun(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	var (
		run            model.CaseStudyWorkflowRun
		projectID      *uuid.UUID
		sourcePayload  []byte
		optionsPayload []byte
		scopePayload   []byte
		canonicalURL   *string
		lastError      *string
	)

	err := r.db.QueryRow(ctx, `
		SELECT id, status, source_payload, options_payload, canonical_url, project_id, last_error, generation_scope_payload, created_at, updated_at
		FROM case_study_workflow_runs
		WHERE id = $1`, runID,
	).Scan(&run.ID, &run.Status, &sourcePayload, &optionsPayload, &canonicalURL, &projectID, &lastError, &scopePayload, &run.CreatedAt, &run.UpdatedAt)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("load case study workflow run: %w", err)
	}

	_ = json.Unmarshal(sourcePayload, &run.Source)
	_ = json.Unmarshal(optionsPayload, &run.Options)
	_ = json.Unmarshal(scopePayload, &run.GenerationScope)
	if canonicalURL != nil {
		run.CanonicalURL = *canonicalURL
	}
	if projectID != nil {
		run.ProjectID = projectID
	}
	if lastError != nil {
		run.LastError = *lastError
	}

	rows, err := r.db.Query(ctx, `
		SELECT run_id, step_name, status, requires_confirmation, confirmation_granted_at, started_at, finished_at, attempt_count, COALESCE(error_message, ''), COALESCE(output_payload, '{}'::jsonb)
		FROM case_study_workflow_steps
		WHERE run_id = $1`, runID)
	if err != nil {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("load case study workflow steps: %w", err)
	}
	defer rows.Close()

	stepsByName := map[string]model.CaseStudyWorkflowStep{}
	for rows.Next() {
		var step model.CaseStudyWorkflowStep
		if err := rows.Scan(&step.RunID, &step.Step, &step.Status, &step.RequiresConfirmation, &step.ConfirmationGrantedAt, &step.StartedAt, &step.FinishedAt, &step.AttemptCount, &step.ErrorMessage, &step.Output); err != nil {
			return model.CaseStudyWorkflowRun{}, fmt.Errorf("scan case study workflow step: %w", err)
		}
		stepsByName[step.Step] = step
	}
	if err := rows.Err(); err != nil {
		return model.CaseStudyWorkflowRun{}, fmt.Errorf("iterate case study workflow steps: %w", err)
	}

	steps := make([]model.CaseStudyWorkflowStep, 0, len(model.CaseStudyWorkflowOrderedSteps))
	for _, stepName := range model.CaseStudyWorkflowOrderedSteps {
		if step, ok := stepsByName[stepName]; ok {
			steps = append(steps, step)
		}
	}
	run.Steps = steps
	return run, nil
}

func (r *CaseStudyWorkflowRepository) ListLogs(ctx context.Context, runID uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, run_id, step_name, level, message, created_at
		FROM case_study_workflow_logs
		WHERE run_id = $1
		ORDER BY id ASC`, runID)
	if err != nil {
		return nil, fmt.Errorf("load case study workflow logs: %w", err)
	}
	defer rows.Close()

	logs := make([]model.CaseStudyWorkflowLogEntry, 0)
	for rows.Next() {
		var entry model.CaseStudyWorkflowLogEntry
		if err := rows.Scan(&entry.ID, &entry.RunID, &entry.Step, &entry.Level, &entry.Message, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan case study workflow log: %w", err)
		}
		logs = append(logs, entry)
	}
	return logs, rows.Err()
}

func (r *CaseStudyWorkflowRepository) AppendLog(ctx context.Context, entry model.CaseStudyWorkflowLogEntry) error {
	createdAt := entry.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO case_study_workflow_logs (run_id, step_name, level, message, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		entry.RunID,
		entry.Step,
		entry.Level,
		entry.Message,
		createdAt,
	)
	if err != nil {
		return fmt.Errorf("insert case study workflow log: %w", err)
	}
	return nil
}

func normalizeJSONB(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "{}"
	}
	return string(raw)
}

func sortSteps(steps []model.CaseStudyWorkflowStep) {
	order := map[string]int{}
	for index, step := range model.CaseStudyWorkflowOrderedSteps {
		order[step] = index
	}
	sort.Slice(steps, func(i, j int) bool {
		return order[steps[i].Step] < order[steps[j].Step]
	})
}

var _ = sortSteps
var _ pgx.Row
