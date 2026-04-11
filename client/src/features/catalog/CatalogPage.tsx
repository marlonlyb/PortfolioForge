import { useEffect, useState, useRef } from 'react';
import { Link } from 'react-router-dom';

import { fetchProducts } from './api';
import { searchProjects } from '../search/api';
import type { ProductSummary } from '../../shared/types/product';
import type { SearchResult } from '../../shared/types/search';
import { AppError } from '../../shared/api/errors';

interface FilterState {
  search: string;
  category: string;
}

const EMPTY_FILTERS: FilterState = {
  search: '',
  category: '',
};

function uniqueSorted(values: (string | undefined)[]): string[] {
  return [...new Set(values.filter((value): value is string => Boolean(value?.trim())))].sort();
}

function summarize(text?: string): string {
  if (!text) return 'Interactive case study coming soon.';
  return text.length > 140 ? `${text.slice(0, 137)}...` : text;
}

export function CatalogPage() {
  const [projects, setProjects] = useState<ProductSummary[]>([]);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<FilterState>(EMPTY_FILTERS);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  // Initial load of all products
  useEffect(() => {
    let cancelled = false;

    fetchProducts()
      .then((response) => {
        if (!cancelled) {
          setProjects(response.items);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Could not load projects.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  // Server-side search when search text changes
  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);

    const trimmed = filters.search.trim();
    if (trimmed.length < 2) {
      setSearchResults([]);
      setSearching(false);
      return;
    }

    setSearching(true);
    debounceRef.current = setTimeout(() => {
      searchProjects({ q: trimmed, category: filters.category || undefined })
        .then((response) => {
          setSearchResults(response.data);
          setSearching(false);
        })
        .catch(() => {
          setSearchResults([]);
          setSearching(false);
        });
    }, 300);

    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [filters.search, filters.category]);

  const categories = uniqueSorted(projects.map((project) => project.category));

  // When searching with results, show search results; otherwise, show filtered catalog
  const isSearchActive = filters.search.trim().length >= 2;
  const filtered = isSearchActive
    ? [] // Don't use client-side filter when server search is active
    : projects.filter((project) => {
        if (filters.category && project.category !== filters.category) {
          return false;
        }
        return true;
      });

  if (loading) {
    return (
      <section className="catalog">
        <p className="catalog__loading">Loading projects…</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="catalog">
        <div className="card card--muted">
          <p className="eyebrow">Error</p>
          <p>{error}</p>
        </div>
      </section>
    );
  }

  return (
    <section className="catalog">
      <h2>Projects</h2>

      <div className="catalog__filters">
        <input
          className="catalog__filter-input"
          type="search"
          placeholder="Search by project name…"
          value={filters.search}
          onChange={(e) => setFilters((prev) => ({ ...prev, search: e.target.value }))}
          aria-label="Search projects"
        />

        <select
          className="catalog__filter-select"
          value={filters.category}
          onChange={(e) => setFilters((prev) => ({ ...prev, category: e.target.value }))}
          aria-label="Filter by category"
        >
          <option value="">Category</option>
          {categories.map((category) => (
            <option key={category} value={category}>
              {category}
            </option>
          ))}
        </select>

        {(filters.search || filters.category) && (
          <button
            className="catalog__filter-clear"
            type="button"
            onClick={() => setFilters(EMPTY_FILTERS)}
          >
            Clear filters
          </button>
        )}
      </div>

      {isSearchActive ? (
        <>
          <p className="catalog__count">
            {searching ? 'Buscando…' : `${searchResults.length} resultado${searchResults.length !== 1 ? 's' : ''}`}
          </p>
          {searchResults.length === 0 && !searching && (
            <div className="card card--muted">
              <p>No projects match your search.</p>
            </div>
          )}
          <div className="catalog__grid">
            {searchResults.map((result) => (
              <Link key={result.id} className="catalog__card" to={`/projects/${result.slug}`}>
                {result.hero_image ? (
                  <img
                    className="catalog__card-img"
                    src={result.hero_image}
                    alt={result.title}
                    loading="lazy"
                  />
                ) : (
                  <div className="catalog__card-img catalog__card-img--placeholder">
                    No image
                  </div>
                )}
                <div className="catalog__card-body">
                  {result.category ? <p className="eyebrow">{result.category}</p> : null}
                  <h3>{result.title}</h3>
                  {result.client_name && (
                    <p className="detail__brand">{result.client_name}</p>
                  )}
                  <p>{summarize(result.summary ?? undefined)}</p>
                </div>
              </Link>
            ))}
          </div>
        </>
      ) : (
        <>
          <p className="catalog__count">
            {filtered.length} project{filtered.length !== 1 ? 's' : ''}
          </p>
          {filtered.length === 0 ? (
            <div className="card card--muted">
              <p>No projects match the current filters.</p>
            </div>
          ) : (
            <div className="catalog__grid">
              {filtered.map((project) => {
                const image = project.images[0] || project.variants?.[0]?.image_url;

                return (
                  <Link key={project.id} className="catalog__card" to={`/projects/${project.id}`}>
                    {image ? (
                      <img className="catalog__card-img" src={image} alt={project.name} loading="lazy" />
                    ) : (
                      <div className="catalog__card-img catalog__card-img--placeholder">
                        No image
                      </div>
                    )}

                    <div className="catalog__card-body">
                      {project.category ? <p className="eyebrow">{project.category}</p> : null}
                      <h3>{project.name}</h3>
                      {project.brand ? <p className="detail__brand">{project.brand}</p> : null}
                      <p>{summarize(project.description)}</p>
                    </div>
                  </Link>
                );
              })}
            </div>
          )}
        </>
      )}
    </section>
  );
}
