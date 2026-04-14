import { useState } from 'react';

import { CatalogPage } from '../catalog/CatalogPage';
import { SearchBar } from '../search/SearchBar';
import { useLocale } from '../../app/providers/LocaleProvider';

const LANDING_SEARCH_COPY = {
  es: {
    eyebrow: 'Buscar proyectos',
    lead: 'Explora por tecnología, cliente, industria o problema.',
    hint: 'Tecnología · Cliente · Problema',
    topics: ['PLC', 'Robotica', 'Motion', 'Safety', 'AI'],
  },
  ca: {
    eyebrow: 'Cercar projectes',
    lead: 'Explora per tecnologia, client, indústria o problema.',
    hint: 'Tecnologia · Client · Problema',
    topics: ['PLC', 'Robotica', 'Motion', 'Safety', 'AI'],
  },
  en: {
    eyebrow: 'Search projects',
    lead: 'Explore by technology, client, industry, or problem.',
    hint: 'Technology · Client · Problem',
    topics: ['PLC', 'Robotica', 'Motion', 'Safety', 'AI'],
  },
  de: {
    eyebrow: 'Projekte suchen',
    lead: 'Nach Technologie, Kunde, Branche oder Problem erkunden.',
    hint: 'Technologie · Kunde · Problem',
    topics: ['PLC', 'Robotica', 'Motion', 'Safety', 'AI'],
  },
} as const;

export function LandingPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const { locale } = useLocale();
  const searchCopy = LANDING_SEARCH_COPY[locale] ?? LANDING_SEARCH_COPY.es;
  const shouldShowQuickTopics = searchQuery.trim().length === 0;

  function scrollToProjects() {
    document.getElementById('projects')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  function applyQuickSearch(topic: string) {
    setSearchQuery(topic);
    scrollToProjects();
  }

  return (
    <>
      <article className="card landing-hero landing-hero--search landing-page">
        <div className="landing-hero__content">
          <div className="landing-hero__search-focus">
            <div className="landing-search-panel">
              <div className="landing-search-panel__header">
                <p className="eyebrow">{searchCopy.eyebrow}</p>
                <p className="landing-search-panel__lead">{searchCopy.lead}</p>
              </div>

              <SearchBar
                value={searchQuery}
                onQueryChange={setSearchQuery}
                onSearch={(query) => {
                  setSearchQuery(query);
                  scrollToProjects();
                }}
                showSubmit={false}
                suggestions={suggestions}
                contextHint={searchCopy.hint}
                onSuggestionSelect={(suggestion) => {
                  setSearchQuery(suggestion);
                  scrollToProjects();
                }}
              />

              {shouldShowQuickTopics ? (
                <div className="landing-search-panel__topics" aria-label={searchCopy.eyebrow}>
                  {searchCopy.topics.map((topic) => (
                    <button
                      key={topic}
                      className="landing-search-panel__topic"
                      type="button"
                      onClick={() => applyQuickSearch(topic)}
                    >
                      {topic}
                    </button>
                  ))}
                </div>
              ) : null}
            </div>
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
