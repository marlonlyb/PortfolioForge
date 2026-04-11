import { Link } from 'react-router-dom';

export function NotFoundPage() {
  return (
    <section className="card-stack">
      <article className="card">
        <p className="eyebrow">404</p>
        <h2>Page not found</h2>
        <p>The store SPA shell is running, but this route does not exist.</p>
        <Link className="nav-link nav-link--active" to="/">
          Go back home
        </Link>
      </article>
    </section>
  );
}
