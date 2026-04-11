interface PagePlaceholderProps {
  title: string;
  description: string;
  aside?: string;
}

export function PagePlaceholder({ title, description, aside }: PagePlaceholderProps) {
  return (
    <section className="card-stack">
      <article className="card">
        <p className="eyebrow">Phase 1 scaffold</p>
        <h2>{title}</h2>
        <p>{description}</p>
      </article>

      <article className="card card--muted">
        <h3>Next phase handoff</h3>
        <p>
          Router, layouts and entrypoint are ready so feature work can focus on API integration and
          persisted state.
        </p>
        {aside ? <p className="aside-note">{aside}</p> : null}
      </article>
    </section>
  );
}
