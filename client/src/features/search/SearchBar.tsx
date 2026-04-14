import { useEffect, useState, type KeyboardEvent } from 'react';
import { useNavigate } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';

interface SearchBarProps {
  initialQuery?: string;
  value?: string;
  onQueryChange?: (query: string) => void;
  onSearch?: (query: string) => void;
  loading?: boolean;
  showSubmit?: boolean;
  suggestions?: string[];
  onSuggestionSelect?: (query: string) => void;
}

export function SearchBar({
  initialQuery = '',
  value,
  onQueryChange,
  onSearch,
  loading,
  showSubmit = true,
  suggestions = [],
  onSuggestionSelect,
}: SearchBarProps) {
  const [internalQuery, setInternalQuery] = useState(initialQuery);
  const navigate = useNavigate();
  const { t } = useLocale();
  const isControlled = value !== undefined;
  const query = isControlled ? value : internalQuery;

  useEffect(() => {
    if (!isControlled) {
      setInternalQuery(initialQuery);
    }
  }, [initialQuery]);

  function updateQuery(nextQuery: string) {
    if (!isControlled) {
      setInternalQuery(nextQuery);
    }

    onQueryChange?.(nextQuery);
  }

  function executeSearch(value: string) {
    const trimmed = value.trim();
    if (trimmed.length < 2) return;

    if (onSearch) {
      onSearch(trimmed);
    } else {
      navigate(`/search?q=${encodeURIComponent(trimmed)}`);
    }
  }

  function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Escape') {
      updateQuery('');
      (e.target as HTMLInputElement).blur();
    }
  }

  function handleClear() {
    updateQuery('');
  }

  return (
    <form
      className="search-bar"
      onSubmit={(event) => {
        event.preventDefault();
        executeSearch(query);
      }}
    >
      <div className="search-bar__field">
        <span className="search-bar__icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" focusable="false">
            <circle cx="11" cy="11" r="6.5" fill="none" stroke="currentColor" strokeWidth="1.8" />
            <path d="m16 16 4.5 4.5" fill="none" stroke="currentColor" strokeLinecap="round" strokeWidth="1.8" />
          </svg>
        </span>

        <input
          className={showSubmit ? 'search-bar__input' : 'search-bar__input search-bar__input--compact'}
          type="text"
          placeholder={t.searchPlaceholder}
          value={query}
          onChange={(e) => updateQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={loading}
          aria-label={t.searchButton}
        />

        <div className="search-bar__actions">
          {query && !loading && (
            <button
              className="search-bar__clear"
              onClick={handleClear}
              type="button"
              aria-label={t.searchClear}
            >
              ×
            </button>
          )}

          {showSubmit && (
            <button
              className="btn btn--ghost search-bar__submit"
              type="submit"
              disabled={loading || query.trim().length < 2}
            >
              {t.searchButton}
            </button>
          )}

          {loading && <span className="search-bar__spinner" aria-label="Buscando…" />}
        </div>
      </div>

      {suggestions.length > 0 && query.trim() && (
        <div className="search-bar__suggestions" aria-label={t.searchSuggestionsLabel}>
          {suggestions.map((suggestion) => (
            <button
              key={suggestion}
              className="search-bar__suggestion"
              type="button"
              onClick={() => {
                updateQuery(suggestion);
                onSuggestionSelect?.(suggestion);
              }}
            >
              {suggestion}
            </button>
          ))}
        </div>
      )}
    </form>
  );
}
