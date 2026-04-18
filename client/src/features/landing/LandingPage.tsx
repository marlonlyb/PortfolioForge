import { useState } from 'react';

import { CatalogPage } from '../catalog/CatalogPage';
import { SearchBar } from '../search/SearchBar';
import { useLocale } from '../../app/providers/LocaleProvider';

export function LandingPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const { t } = useLocale();
  const shouldShowQuickTopics = searchQuery.trim().length === 0;

  function scrollToProjects() {
    document.getElementById('projects')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  function submitLandingQuery(query: string) {
    setSearchQuery(query);
    scrollToProjects();
  }

  return (
    <>
      <article className="card landing-hero landing-hero--search landing-page">
        <div className="landing-hero__content">
          <div className="landing-hero__search-focus">
            <div className="landing-search-panel">
              <div className="landing-search-panel__header">
                <p className="eyebrow">{t.landingSearchEyebrow}</p>
                {t.landingSearchTitle ? <h1 className="landing-search-panel__title">{t.landingSearchTitle}</h1> : null}
                <p className="landing-search-panel__lead">{t.landingSearchLead}</p>
              </div>

              <SearchBar
                value={searchQuery}
                onQueryChange={setSearchQuery}
                onSearch={submitLandingQuery}
                showSubmit={false}
                suggestions={suggestions}
                contextHint={t.landingSearchContextHint}
                placeholderOverride={t.landingSearchPlaceholder}
                onSuggestionSelect={submitLandingQuery}
              />

              {shouldShowQuickTopics ? (
                <div className="landing-search-panel__topics" aria-label={t.landingSearchEyebrow}>
                  {t.landingQuickPrompts.map((prompt) => (
                    <button
                      key={prompt.label}
                      className="landing-search-panel__topic"
                      type="button"
                      onClick={() => submitLandingQuery(prompt.query)}
                    >
                      <span>{prompt.label}</span>
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
