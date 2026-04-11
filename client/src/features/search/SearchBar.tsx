import { useState, useEffect, useRef, type KeyboardEvent } from 'react';
import { useNavigate } from 'react-router-dom';

interface SearchBarProps {
  initialQuery?: string;
  onSearch?: (query: string) => void;
  loading?: boolean;
}

export function SearchBar({ initialQuery = '', onSearch, loading }: SearchBarProps) {
  const [query, setQuery] = useState(initialQuery);
  const navigate = useNavigate();
  const debounceRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  function executeSearch(value: string) {
    const trimmed = value.trim();
    if (trimmed.length < 2) return;

    if (onSearch) {
      onSearch(trimmed);
    } else {
      navigate(`/search?q=${encodeURIComponent(trimmed)}`);
    }
  }

  function handleChange(value: string) {
    setQuery(value);
    if (debounceRef.current) clearTimeout(debounceRef.current);

    if (value.trim().length >= 2) {
      debounceRef.current = setTimeout(() => executeSearch(value), 300);
    }
  }

  function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter') {
      if (debounceRef.current) clearTimeout(debounceRef.current);
      executeSearch(query);
    } else if (e.key === 'Escape') {
      setQuery('');
      if (debounceRef.current) clearTimeout(debounceRef.current);
      (e.target as HTMLInputElement).blur();
    }
  }

  function handleClear() {
    setQuery('');
    if (debounceRef.current) clearTimeout(debounceRef.current);
  }

  return (
    <div className="search-bar">
      <input
        className="search-bar__input"
        type="search"
        placeholder="Busca proyectos por tecnología, cliente o concepto…"
        value={query}
        onChange={(e) => handleChange(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={loading}
        aria-label="Buscar proyectos"
      />
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
      {loading && <span className="search-bar__spinner" aria-label="Buscando…" />}
    </div>
  );
}
