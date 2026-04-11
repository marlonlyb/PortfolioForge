import { Link } from 'react-router-dom';

import { SearchBar } from '../search/SearchBar';

export function LandingPage() {
  return (
    <section className="card-stack landing-page">
      <article className="card">
        <p className="eyebrow">PortfolioForge</p>
        <h2>AI-assisted portfolio platform</h2>
        <p>
          Public landing + interactive project catalog + private admin console to curate case
          studies, technical decisions, media, and future contact leads.
        </p>
        <SearchBar />
        <div className="landing__sub-actions">
          <Link className="btn btn--ghost" to="/projects">
            Explore projects
          </Link>
          <Link className="btn btn--ghost" to="/login">
            Admin access
          </Link>
        </div>
      </article>

      <article className="card">
        <h3>Planned portfolio profile</h3>
        <p>
          Each project will evolve from a simple catalog entry into a full case study with
          architecture, AI usage, results, media, and management context.
        </p>
      </article>
    </section>
  );
}
