import { useEffect, useState, type KeyboardEvent } from 'react';
import { useNavigate } from 'react-router-dom';

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
        <input
          className={showSubmit ? 'search-bar__input' : 'search-bar__input search-bar__input--compact'}
          type="text"
          placeholder="Busca proyectos por tecnología, cliente o concepto…"
          value={query}
        onChange={(e) => updateQuery(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={loading}
        aria-label="Buscar proyectos"
      />

      <div className="search-bar__actions">
        {query && !loading && (
          <button
            className="search-bar__clear"
            onClick={handleClear}
            type="button"
            aria-label="Limpiar búsqueda"
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
            Buscar
          </button>
        )}

        {loading && <span className="search-bar__spinner" aria-label="Buscando…" />}
      </div>

      {suggestions.length > 0 && query.trim() && (
        <div className="search-bar__suggestions" aria-label="Sugerencias de búsqueda">
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
