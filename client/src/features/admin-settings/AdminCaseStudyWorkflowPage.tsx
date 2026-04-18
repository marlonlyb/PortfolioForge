import { useEffect, useState, type FormEvent } from 'react';
import { Link, useSearchParams } from 'react-router-dom';

import { AppError } from '../../shared/api/errors';
import {
  confirmCaseStudyWorkflowStep,
  fetchCaseStudyWorkflowAvailability,
  fetchCaseStudyWorkflowLogs,
  fetchCaseStudyWorkflowRun,
  isWorkflowUnavailableError,
  resumeCaseStudyWorkflowRun,
  retryCaseStudyWorkflowStep,
  startCaseStudyWorkflowRun,
  startCaseStudyWorkflowStep,
  type CaseStudyWorkflowAvailability,
  type CaseStudyWorkflowLogEntry,
  type CaseStudyWorkflowRun,
  type CaseStudyWorkflowStatus,
  type CaseStudyWorkflowStep,
  type CaseStudyWorkflowStepName,
} from './api';

const STORAGE_KEY = 'admin.case-study-workflow.last-run-id';

const STEP_LABELS: Record<CaseStudyWorkflowStepName, string> = {
  resolve_source: 'Resolve canonical source',
  publish_canonical: 'Publish canonical files',
  import_or_update_project: 'Create or update PortfolioForge project',
  localization_backfill: 'Localization backfill',
  reembed: 'Refresh search embedding document',
};

const STATUS_LABELS: Record<CaseStudyWorkflowStatus, string> = {
  pending: 'Ready',
  blocked: 'Blocked',
  awaiting_confirmation: 'Needs confirmation',
  running: 'Running',
  succeeded: 'Done',
  failed: 'Failed',
  skipped: 'Skipped',
};

export function AdminCaseStudyWorkflowPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [sourcePath, setSourcePath] = useState('');
  const [slug, setSlug] = useState('');
  const [runLocalizationBackfill, setRunLocalizationBackfill] = useState(true);
  const [runReembed, setRunReembed] = useState(true);
  const [localesRaw, setLocalesRaw] = useState('ca,en,de');
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [run, setRun] = useState<CaseStudyWorkflowRun | null>(null);
  const [logs, setLogs] = useState<CaseStudyWorkflowLogEntry[]>([]);
  const [availability, setAvailability] = useState<CaseStudyWorkflowAvailability | null>(null);
  const [availabilityLoading, setAvailabilityLoading] = useState(true);

  const runId = searchParams.get('run') ?? sessionStorage.getItem(STORAGE_KEY) ?? '';

  useEffect(() => {
    let cancelled = false;

    fetchCaseStudyWorkflowAvailability()
      .then((status) => {
        if (cancelled) return;
        setAvailability(status);
        if (!status.configured) {
          sessionStorage.removeItem(STORAGE_KEY);
          setRun(null);
          setLogs([]);
          if (searchParams.get('run')) {
            setSearchParams({}, { replace: true });
          }
        }
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        setAvailability({ configured: true });
        setError(err instanceof AppError ? err.message : 'Failed to load workflow availability.');
      })
      .finally(() => {
        if (!cancelled) setAvailabilityLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [searchParams, setSearchParams]);

  useEffect(() => {
    if (availabilityLoading || availability?.configured === false) return;
    if (!runId) return;

    let cancelled = false;
    setRefreshing(true);
    Promise.all([fetchCaseStudyWorkflowRun(runId), fetchCaseStudyWorkflowLogs(runId)])
      .then(([workflowRun, logPayload]) => {
        if (cancelled) return;
        setRun(workflowRun);
        setLogs(logPayload.items);
        sessionStorage.setItem(STORAGE_KEY, workflowRun.id);
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        setError(err instanceof AppError ? err.message : 'Failed to load workflow run.');
      })
      .finally(() => {
        if (!cancelled) setRefreshing(false);
      });

    return () => {
      cancelled = true;
    };
  }, [availability?.configured, availabilityLoading, runId]);

  const nextActionableStep = run?.steps.find((step) => step.status === 'pending' || step.status === 'awaiting_confirmation' || step.status === 'failed') ?? null;

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (availability?.configured === false) return;
    setLoading(true);
    setError(null);

    try {
      const created = await startCaseStudyWorkflowRun({
        source_path: sourcePath,
        slug: slug || undefined,
        run_localization_backfill: runLocalizationBackfill,
        run_reembed: runReembed,
        locales: parseLocales(localesRaw),
      });
      sessionStorage.setItem(STORAGE_KEY, created.id);
      setSearchParams({ run: created.id });
      setRun(created);
      const logPayload = await fetchCaseStudyWorkflowLogs(created.id);
      setLogs(logPayload.items);
    } catch (err: unknown) {
      if (isWorkflowUnavailableError(err)) {
        sessionStorage.removeItem(STORAGE_KEY);
        setAvailability({ configured: false, reason: err.message });
        setRun(null);
        setLogs([]);
        setSearchParams({}, { replace: true });
        setError(null);
      } else {
        setError(err instanceof AppError ? err.message : 'Failed to start workflow run.');
      }
    } finally {
      setLoading(false);
    }
  }

  async function refreshRun(updatedRun?: CaseStudyWorkflowRun) {
    if (availability?.configured === false) return;
    const targetRunId = updatedRun?.id ?? run?.id;
    if (!targetRunId) return;

    setRefreshing(true);
    setError(null);
    try {
      const [currentRun, currentLogs] = await Promise.all([
        updatedRun ? Promise.resolve(updatedRun) : fetchCaseStudyWorkflowRun(targetRunId),
        fetchCaseStudyWorkflowLogs(targetRunId),
      ]);
      setRun(currentRun);
      setLogs(currentLogs.items);
    } catch (err: unknown) {
      if (isWorkflowUnavailableError(err)) {
        sessionStorage.removeItem(STORAGE_KEY);
        setAvailability({ configured: false, reason: err.message });
        setRun(null);
        setLogs([]);
        setSearchParams({}, { replace: true });
        setError(null);
      } else {
        setError(err instanceof AppError ? err.message : 'Failed to refresh workflow state.');
      }
    } finally {
      setRefreshing(false);
    }
  }

  async function handleConfirm(step: CaseStudyWorkflowStepName) {
    if (!run) return;
    const updated = await confirmCaseStudyWorkflowStep(run.id, step);
    await refreshRun(updated);
  }

  async function handleStart(step: CaseStudyWorkflowStepName) {
    if (!run) return;
    const updated = await startCaseStudyWorkflowStep(run.id, step);
    await refreshRun(updated);
  }

  async function handleRetry(step: CaseStudyWorkflowStepName) {
    if (!run) return;
    const updated = await retryCaseStudyWorkflowStep(run.id, step);
    await refreshRun(updated);
  }

  async function handleResume() {
    if (!run) return;
    const updated = await resumeCaseStudyWorkflowRun(run.id);
    await refreshRun(updated);
  }

  if (availabilityLoading) {
    return (
      <section className="card-stack">
        <p className="admin__loading">Loading workflow availability…</p>
      </section>
    );
  }

  if (availability?.configured === false) {
    return (
      <section className="card-stack">
        <article className="card">
          <p className="eyebrow">Settings</p>
          <h2>Case-study workflow</h2>
          <p className="admin__helper-copy">
            Workflow unavailable. {availability.reason ?? 'Case-study workflow is not configured in this environment.'}
          </p>
          <p className="admin__helper-copy">
            Configure the workflow environment before running publish/import actions.
          </p>
          <div className="admin__form-actions">
            <Link className="btn btn--secondary" to="/admin/settings">
              Back to settings
            </Link>
          </div>
        </article>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <p className="eyebrow">Settings</p>
        <h2>Case-study workflow</h2>
        <p className="admin__helper-copy">
          Start from an already canonical source under <code>90. dev_portfolioforge/&lt;slug&gt;/</code>,
          then guide publish, import, localization, and re-embed step by step.
        </p>
        <p className="admin__helper-copy">
          Raw folder → canonical generation is intentionally out of scope for this MVP.
        </p>

        {error ? <div className="admin__error" role="alert">{error}</div> : null}

        <form className="admin__form" onSubmit={handleSubmit}>
          <label className="admin__label">
            Canonical source path
            <input
              className="admin__input"
              placeholder="/allowed/root/90. dev_portfolioforge/my-case-study"
              value={sourcePath}
              onChange={(event) => setSourcePath(event.target.value)}
            />
          </label>

          <label className="admin__label">
            Slug override (optional)
            <input
              className="admin__input"
              placeholder="my-case-study"
              value={slug}
              onChange={(event) => setSlug(event.target.value)}
            />
          </label>

          <label className="admin__label">
            Localization locales (optional comma-separated subset)
            <input
              className="admin__input"
              placeholder="ca,en,de"
              value={localesRaw}
              onChange={(event) => setLocalesRaw(event.target.value)}
            />
          </label>

          <div className="admin__form-section">
            <label className="admin__checkbox">
              <input
                checked={runLocalizationBackfill}
                onChange={(event) => setRunLocalizationBackfill(event.target.checked)}
                type="checkbox"
              />
              Run localization backfill after import
            </label>

            <label className="admin__checkbox">
              <input
                checked={runReembed}
                onChange={(event) => setRunReembed(event.target.checked)}
                type="checkbox"
              />
              Refresh search document after import/localization
            </label>
          </div>

          <div className="admin__form-actions">
            <button className="btn btn--primary" disabled={loading} type="submit">
              {loading ? 'Starting…' : 'Start workflow run'}
            </button>
            <Link className="btn btn--secondary" to="/admin/settings">
              Back to settings
            </Link>
          </div>
        </form>
      </article>

      {run ? (
        <article className="card">
          <div className="admin__form-actions" style={{ justifyContent: 'space-between' }}>
            <div>
              <p className="eyebrow">Run</p>
              <h3>{run.source.slug}</h3>
              <p className="admin__helper-copy">
                Status: <strong>{STATUS_LABELS[run.status]}</strong>
              </p>
            </div>
            <button className="btn btn--secondary" disabled={refreshing} onClick={() => void refreshRun()} type="button">
              {refreshing ? 'Refreshing…' : 'Refresh status'}
            </button>
          </div>

          <dl className="admin__details-grid">
            <div>
              <dt>Source</dt>
              <dd>{run.source.normalized_path}</dd>
            </div>
            <div>
              <dt>Canonical markdown</dt>
              <dd>{run.source.canonical_markdown_path}</dd>
            </div>
            <div>
              <dt>Published URL</dt>
              <dd>{run.canonical_url ?? 'Not published yet'}</dd>
            </div>
            <div>
              <dt>Project ID</dt>
              <dd>{run.project_id ?? 'Project not created yet'}</dd>
            </div>
          </dl>

          <p className="admin__helper-copy">{run.generation_scope.canonical_generation_message}</p>

          {nextActionableStep ? (
            <div className="admin__form-actions">
              <button className="btn btn--secondary" onClick={() => void handleResume()} type="button">
                Continue from latest checkpoint
              </button>
            </div>
          ) : null}

          <div className="card-stack">
            {run.steps.map((step) => (
              <StepCard
                key={step.step}
                step={step}
                onConfirm={() => void handleConfirm(step.step)}
                onRetry={() => void handleRetry(step.step)}
                onStart={() => void handleStart(step.step)}
              />
            ))}
          </div>
        </article>
      ) : null}

      {run ? (
        <article className="card">
          <p className="eyebrow">Operator log</p>
          <h3>Run timeline</h3>
          {logs.length === 0 ? <p className="admin__helper-copy">No logs yet.</p> : null}
          <ul className="admin__list">
            {logs.map((entry) => (
              <li key={entry.id}>
                <strong>{STEP_LABELS[entry.step]}</strong> · {entry.level.toUpperCase()} · {entry.message}
              </li>
            ))}
          </ul>
        </article>
      ) : null}
    </section>
  );
}

function StepCard({
  step,
  onConfirm,
  onStart,
  onRetry,
}: {
  step: CaseStudyWorkflowStep;
  onConfirm: () => void;
  onStart: () => void;
  onRetry: () => void;
}) {
  const outputEntries = Object.entries(step.output ?? {});

  return (
    <article className="card" data-step={step.step}>
      <div className="admin__form-actions" style={{ justifyContent: 'space-between' }}>
        <div>
          <h4>{STEP_LABELS[step.step]}</h4>
          <p className="admin__helper-copy">{STATUS_LABELS[step.status]}</p>
        </div>
        <div className="admin__form-actions">
          {step.status === 'awaiting_confirmation' ? (
            <button className="btn btn--secondary" onClick={onConfirm} type="button">
              Confirm
            </button>
          ) : null}
          {step.status === 'pending' ? (
            <button className="btn btn--primary" onClick={onStart} type="button">
              Run step
            </button>
          ) : null}
          {step.status === 'failed' ? (
            <button className="btn btn--secondary" onClick={onRetry} type="button">
              Retry step
            </button>
          ) : null}
        </div>
      </div>

      {step.error_message ? <div className="admin__error">{step.error_message}</div> : null}

      {outputEntries.length > 0 ? (
        <ul className="admin__list">
          {outputEntries.map(([key, value]) => (
            <li key={key}>
              <strong>{key}:</strong> {String(value)}
            </li>
          ))}
        </ul>
      ) : null}
    </article>
  );
}

function parseLocales(value: string): string[] {
  return value
    .split(',')
    .map((locale) => locale.trim().toLowerCase())
    .filter(Boolean);
}
