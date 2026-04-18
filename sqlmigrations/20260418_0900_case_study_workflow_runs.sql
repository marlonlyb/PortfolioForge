CREATE TABLE IF NOT EXISTS case_study_workflow_runs (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    source_payload JSONB NOT NULL,
    options_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    canonical_url TEXT NULL,
    project_id UUID NULL REFERENCES products(id) ON DELETE SET NULL,
    last_error TEXT NULL,
    generation_scope_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS case_study_workflow_steps (
    run_id UUID NOT NULL REFERENCES case_study_workflow_runs(id) ON DELETE CASCADE,
    step_name TEXT NOT NULL,
    status TEXT NOT NULL,
    requires_confirmation BOOLEAN NOT NULL DEFAULT FALSE,
    confirmation_granted_at TIMESTAMPTZ NULL,
    started_at TIMESTAMPTZ NULL,
    finished_at TIMESTAMPTZ NULL,
    attempt_count INTEGER NOT NULL DEFAULT 0,
    error_message TEXT NULL,
    output_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY (run_id, step_name)
);

CREATE TABLE IF NOT EXISTS case_study_workflow_logs (
    id BIGSERIAL PRIMARY KEY,
    run_id UUID NOT NULL REFERENCES case_study_workflow_runs(id) ON DELETE CASCADE,
    step_name TEXT NOT NULL,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_case_study_workflow_runs_updated_at
    ON case_study_workflow_runs (updated_at DESC);

CREATE INDEX IF NOT EXISTS ix_case_study_workflow_logs_run_id
    ON case_study_workflow_logs (run_id, id ASC);
