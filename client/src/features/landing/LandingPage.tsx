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
    <>
      <article className="card landing-hero landing-hero--search landing-page">
        <div className="landing-hero__content">
          <div className="landing-hero__search-focus">
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
          </div>
        </div>
      </article>

      <div id="projects" className="landing-catalog landing-page">
        <CatalogPage
          searchQuery={searchQuery}
          onSearchQueryChange={setSearchQuery}
          onSuggestionsChange={setSuggestions}
          renderSearchControls={false}
        />
      </div>
    </>
  );
}
