import { useState } from 'react';

import { CatalogPage } from '../catalog/CatalogPage';
import { SearchBar } from '../search/SearchBar';

export function LandingPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);

  function scrollToProjects() {
    document.getElementById('projects')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  return (
    <section className="card-stack landing-page">
      <article className="card">
        <p className="eyebrow">PortfolioForge</p>
        <h2>AI-assisted portfolio platform</h2>
        <p>
          Public landing + interactive project catalog + private admin console to curate case
          studies, technical decisions, media, and future contact leads.
        </p>
        <SearchBar
          value={searchQuery}
          onQueryChange={setSearchQuery}
          onSearch={(query) => {
            setSearchQuery(query);
            scrollToProjects();
          }}
          showSubmit={false}
          suggestions={suggestions}
          onSuggestionSelect={(suggestion) => {
            setSearchQuery(suggestion);
            scrollToProjects();
          }}
        />
        <div className="landing__sub-actions">
          <button className="btn btn--ghost" type="button" onClick={scrollToProjects}>
            Explore projects
          </button>
        </div>
      </article>

      <div id="projects">
        <CatalogPage
          searchQuery={searchQuery}
          onSearchQueryChange={setSearchQuery}
          onSuggestionsChange={setSuggestions}
          renderSearchControls={false}
        />
      </div>
    </section>
  );
}
