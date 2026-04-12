import { useEffect, useState, useCallback } from 'react';
import { Link, useSearchParams } from 'react-router-dom';

import { searchProjects } from './api';
import { SearchBar } from './SearchBar';
import { SearchResultCard } from './SearchResultCard';
import { SearchFilters } from './SearchFilters';
import type { SearchResult, SearchFilters as SearchFiltersType } from '../../shared/types/search';
import { AppError } from '../../shared/api/errors';

const PAGE_SIZE = 10;

export function SearchResultsPage() {
  const [searchParams, setSearchParams] = useSearchParams();

  const query = searchParams.get('q') ?? '';
  const categoryParam = searchParams.get('category') ?? null;
  const techParam = searchParams.get('technologies') ?? '';

  const [results, setResults] = useState<SearchResult[]>([]);
  const [total, setTotal] = useState(0);
  const [cursor, setCursor] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const filters: SearchFiltersType = {
    category: categoryParam,
    client: null,
    technologies: techParam ? techParam.split(',').filter(Boolean) : [],
  };

  const doSearch = useCallback(
    async (q: string, f: SearchFiltersType, nextCursor?: string) => {
      const trimmed = q.trim();
      if (trimmed.length < 2) {
        setResults([]);
        setTotal(0);
        setCursor(null);
        return;
      }

      const isLoadMore = Boolean(nextCursor);
      if (isLoadMore) {
        setLoadingMore(true);
      } else {
        setLoading(true);
      }
      setError(null);

      try {
        const response = await searchProjects({
          q: trimmed,
          category: f.category ?? undefined,
          technologies: f.technologies.length > 0 ? f.technologies.join(',') : undefined,
          pageSize: PAGE_SIZE,
          cursor: nextCursor ?? undefined,
        });

        if (isLoadMore) {
          setResults((prev) => [...prev, ...response.data]);
        } else {
          setResults(response.data);
        }
        setTotal(response.meta.total);
        setCursor(response.meta.cursor);
      } catch (err: unknown) {
        setError(err instanceof AppError ? err.message : 'Error al buscar proyectos.');
      } finally {
        setLoading(false);
        setLoadingMore(false);
      }
    },
    [],
  );

  // Search when query or filters change
  useEffect(() => {
    doSearch(query, filters);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query, categoryParam, techParam]);

  function handleSearch(newQuery: string) {
    setSearchParams((prev) => {
      prev.set('q', newQuery);
      prev.delete('category');
      prev.delete('technologies');
      return prev;
    });
  }

  function handleFilterChange(newFilters: SearchFiltersType) {
    setSearchParams((prev) => {
      if (newFilters.category) {
        prev.set('category', newFilters.category);
      } else {
        prev.delete('category');
      }
      if (newFilters.technologies.length > 0) {
        prev.set('technologies', newFilters.technologies.join(','));
      } else {
        prev.delete('technologies');
      }
      return prev;
    });
  }

  function handleLoadMore() {
    if (cursor) {
      doSearch(query, filters, cursor);
    }
  }

  return (
    <section className="search-results">
      <div className="search-results__header">
        <SearchBar initialQuery={query} onSearch={handleSearch} loading={loading} />
        {query.trim().length >= 2 && !loading && (
          <p className="search-results__count">
            {total} proyecto{total !== 1 ? 's' : ''} encontrado{total !== 1 ? 's' : ''}
          </p>
        )}
      </div>

      {error && (
        <div className="card card--muted" style={{ borderColor: '#e8c4c4', marginBottom: '1rem' }}>
          <p style={{ color: '#b44040' }}>{error}</p>
        </div>
      )}

      {loading && (
        <div className="search-results__empty">
          <span className="search-bar__spinner" style={{ position: 'static', transform: 'none', margin: '0 auto 1rem' }} />
          <p>Buscando…</p>
        </div>
      )}

      {!loading && query.trim().length >= 2 && results.length === 0 && !error && (
        <div className="search-results__empty">
          <p>No se encontraron proyectos para &lsquo;{query}&rsquo;</p>
          <Link className="btn btn--ghost" to="/" style={{ marginTop: '1rem' }}>
            Ver catálogo completo
          </Link>
        </div>
      )}

      {!loading && results.length > 0 && (
        <div className="search-results__layout">
          <div className="search-results__list">
            {results.map((result, index) => (
              <SearchResultCard key={result.id} result={result} index={index} />
            ))}
            {cursor && (
              <div className="search-results__load-more">
                <button
                  className="btn btn--ghost"
                  type="button"
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                >
                  {loadingMore ? 'Cargando…' : 'Cargar más'}
                </button>
              </div>
            )}
          </div>
          <SearchFilters filters={filters} results={results} onFilterChange={handleFilterChange} />
        </div>
      )}

      {!loading && query.trim().length < 2 && (
        <div className="search-results__empty">
          <p>Escribe al menos 2 caracteres para buscar</p>
        </div>
      )}
    </section>
  );
}
