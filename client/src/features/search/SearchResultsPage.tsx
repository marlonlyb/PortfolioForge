import { useEffect, useRef, useState } from 'react';
import { Link, useLocation, useSearchParams } from 'react-router-dom';

import { useLocale } from '../../app/providers/LocaleProvider';
import { searchProjects } from './api';
import { SearchBar } from './SearchBar';
import { SearchResultCard } from './SearchResultCard';
import { SearchFilters } from './SearchFilters';
import type { SearchResult, SearchFilters as SearchFiltersType } from '../../shared/types/search';
import { AppError } from '../../shared/api/errors';
import {
  buildSearchResultsLocationState,
  matchesActiveSearchState,
  type ProjectDetailLocationState,
} from './matchContext';

const PAGE_SIZE = 10;

export function SearchResultsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const location = useLocation();
  const { locale, t } = useLocale();

  const query = searchParams.get('q') ?? '';
  const categoryParam = searchParams.get('category') ?? null;
  const clientParam = searchParams.get('client') ?? null;
  const techParam = searchParams.get('technologies') ?? '';

  const filters: SearchFiltersType = {
    category: categoryParam,
    client: clientParam,
    technologies: techParam ? techParam.split(',').filter(Boolean) : [],
  };
  const locationState = (location.state as ProjectDetailLocationState | null) ?? null;
  const canRestoreFromLocation = matchesActiveSearchState(locationState, query, filters);
  const restoredSnapshot = canRestoreFromLocation ? locationState?.searchResultsSnapshot : undefined;

  const [results, setResults] = useState<SearchResult[]>(() => restoredSnapshot?.results ?? []);
  const [total, setTotal] = useState(() => restoredSnapshot?.total ?? 0);
  const [cursor, setCursor] = useState<string | null>(() => restoredSnapshot?.cursor ?? null);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const skipInitialSearchRef = useRef(canRestoreFromLocation);

  async function doSearch(q: string, f: SearchFiltersType, nextCursor?: string) {
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
        client: f.client ?? undefined,
        technologies: f.technologies.length > 0 ? f.technologies.join(',') : undefined,
        pageSize: PAGE_SIZE,
        cursor: nextCursor ?? undefined,
        lang: locale,
      });

      if (isLoadMore) {
        setResults((prev) => [...prev, ...response.data]);
      } else {
        setResults(response.data);
      }
      setTotal(response.meta.total);
      setCursor(response.meta.cursor);
    } catch (err: unknown) {
      setError(err instanceof AppError ? err.message : t.searchResultsError);
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }

  // Search when query or filters change
  useEffect(() => {
    if (skipInitialSearchRef.current) {
      skipInitialSearchRef.current = false;
      return;
    }

    doSearch(query, filters);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryParam, clientParam, locale, query, techParam]);

  const detailState = buildSearchResultsLocationState(query, filters, {
    results,
    total,
    cursor,
  });

  function handleSearch(newQuery: string) {
    setSearchParams((prev) => {
      prev.set('q', newQuery);
      prev.delete('category');
      prev.delete('client');
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
      if (newFilters.client) {
        prev.set('client', newFilters.client);
      } else {
        prev.delete('client');
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
      <article className="card search-results__hero">
        <div className="search-results__header">
          <div>
            <p className="eyebrow">{t.searchResultsEyebrow}</p>
            <h2>{t.searchResultsTitle}</h2>
            <p className="search-results__intro">
              {t.searchResultsIntro}
            </p>
          </div>
          <SearchBar initialQuery={query} onSearch={handleSearch} loading={loading} />
        </div>
        {query.trim().length >= 2 && !loading && (
          <p className="search-results__count">
            {total} {total === 1 ? t.searchResultsCountSingular : t.searchResultsCountPlural}
          </p>
        )}
      </article>

      {error && (
        <div className="card card--muted" style={{ borderColor: '#e8c4c4', marginBottom: '1rem' }}>
          <p style={{ color: '#b44040' }}>{error}</p>
        </div>
      )}

      {loading && (
        <div className="search-results__empty">
          <span className="search-bar__spinner" style={{ position: 'static', transform: 'none', margin: '0 auto 1rem' }} />
          <p>{t.searchResultsSearching}…</p>
        </div>
      )}

      {!loading && query.trim().length >= 2 && results.length === 0 && !error && (
        <div className="search-results__empty">
          <p>{t.searchResultsNoResults} &lsquo;{query}&rsquo;</p>
          <Link className="btn btn--ghost" to="/" style={{ marginTop: '1rem' }}>
            {t.searchResultsViewCatalog}
          </Link>
        </div>
      )}

      {!loading && results.length > 0 && (
        <div className="search-results__layout">
          <div className="search-results__list">
            {results.map((result, index) => (
              <SearchResultCard key={result.id} result={result} index={index} detailState={detailState} />
            ))}
            {cursor && (
              <div className="search-results__load-more">
                <button
                  className="btn btn--ghost"
                  type="button"
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                >
                  {loadingMore ? t.searchResultsLoadingMore : t.searchResultsLoadMore}
                </button>
              </div>
            )}
          </div>
          <SearchFilters filters={filters} results={results} onFilterChange={handleFilterChange} />
        </div>
      )}

      {!loading && query.trim().length < 2 && (
        <div className="search-results__empty">
          <p>{t.searchResultsMinCharacters}</p>
        </div>
      )}
    </section>
  );
}
