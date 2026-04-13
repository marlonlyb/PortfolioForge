import { useEffect, useMemo, useRef, useState } from 'react';
import { Link } from 'react-router-dom';

import { fetchProjects } from './api';
import { searchProjects } from '../search/api';
import type { Project } from '../../shared/types/project';
import type { SearchResult } from '../../shared/types/search';
import { AppError } from '../../shared/api/errors';

interface CatalogPageProps {
  searchQuery?: string;
  onSearchQueryChange?: (query: string) => void;
  onSuggestionsChange?: (suggestions: string[]) => void;
  renderSearchControls?: boolean;
}

interface FilterState {
  searchDraft: string;
  category: string;
}

interface SuggestionEntry {
  value: string;
  count: number;
}

const EMPTY_FILTERS: FilterState = {
  searchDraft: '',
  category: '',
};

const SEARCH_DEBOUNCE_MS = 250;
const MIN_REMOTE_QUERY_LENGTH = 2;
const MAX_SUGGESTIONS = 6;

const STOP_WORDS = new Set([
  'and',
  'case',
  'con',
  'for',
  'from',
  'para',
  'por',
  'project',
  'projects',
  'study',
  'that',
  'the',
  'this',
]);

function uniqueSorted(values: (string | undefined)[]): string[] {
  return [...new Set(values.filter((value): value is string => Boolean(value?.trim())))].sort();
}

function summarize(text?: string): string {
  if (!text) return 'No summary available.';
  return text.length > 140 ? `${text.slice(0, 137)}...` : text;
}

function normalizeText(value?: string): string {
  return (value ?? '')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim();
}

function splitKeywords(value?: string): string[] {
  return normalizeText(value)
    .split(/[^a-z0-9+#.]+/)
    .filter((token) => token.length >= 2 && !STOP_WORDS.has(token));
}

function getProjectSearchCorpus(project: Project): string[] {
  return [
    project.name,
    project.slug,
    project.category,
    project.client_name,
    project.description,
    project.profile?.business_goal,
    project.profile?.problem_statement,
    project.profile?.solution_summary,
    project.profile?.architecture,
    project.profile?.ai_usage,
    ...(project.technologies?.flatMap((technology) => [technology.name, technology.slug, technology.category]) ?? []),
  ].filter((value): value is string => Boolean(value?.trim()));
}

function matchesProject(project: Project, query: string): boolean {
  const normalizedQuery = normalizeText(query);
  if (!normalizedQuery) return true;

  return getProjectSearchCorpus(project).some((value) => normalizeText(value).includes(normalizedQuery));
}

function buildSuggestions(projects: Project[], category: string, query: string): string[] {
  const normalizedQuery = normalizeText(query);
  if (!normalizedQuery) return [];

  const matches = new Map<string, SuggestionEntry>();

  projects.forEach((project) => {
    if (category && project.category !== category) return;

    const keywords = [
      project.name,
      project.category,
      project.client_name,
      ...(project.technologies?.map((technology) => technology.name) ?? []),
      ...splitKeywords(project.description),
      ...splitKeywords(project.profile?.business_goal),
      ...splitKeywords(project.profile?.problem_statement),
      ...splitKeywords(project.profile?.solution_summary),
      ...splitKeywords(project.profile?.architecture),
      ...splitKeywords(project.profile?.ai_usage),
    ].filter((value): value is string => Boolean(value?.trim()));

    const seenInProject = new Set<string>();

    keywords.forEach((keyword) => {
      const normalizedKeyword = normalizeText(keyword);
      if (!normalizedKeyword || seenInProject.has(normalizedKeyword) || !normalizedKeyword.includes(normalizedQuery)) {
        return;
      }

      seenInProject.add(normalizedKeyword);
      const existing = matches.get(normalizedKeyword);
      matches.set(normalizedKeyword, {
        value: existing?.value ?? keyword.trim(),
        count: (existing?.count ?? 0) + 1,
      });
    });
  });

  return [...matches.values()]
    .sort((left, right) => {
      const leftStartsWith = normalizeText(left.value).startsWith(normalizedQuery) ? 1 : 0;
      const rightStartsWith = normalizeText(right.value).startsWith(normalizedQuery) ? 1 : 0;

      if (leftStartsWith !== rightStartsWith) return rightStartsWith - leftStartsWith;
      if (left.count !== right.count) return right.count - left.count;
      if (left.value.length !== right.value.length) return left.value.length - right.value.length;
      return left.value.localeCompare(right.value);
    })
    .slice(0, MAX_SUGGESTIONS)
    .map((entry) => entry.value);
}

function renderProductCard(project: Project) {
  const image = project.images[0];

  return (
    <Link key={project.id} className="catalog__card" to={`/projects/${project.slug}`}>
      {image ? (
        <img className="catalog__card-img" src={image} alt={project.name} loading="lazy" />
      ) : (
        <div className="catalog__card-img catalog__card-img--placeholder">No image</div>
      )}

      <div className="catalog__card-body">
        {project.category ? <p className="eyebrow">{project.category}</p> : null}
        <h3>{project.name}</h3>
        {project.client_name ? <p className="detail__brand">{project.client_name}</p> : null}
        <p>{summarize(project.description)}</p>
      </div>
    </Link>
  );
}

function renderSearchCard(result: SearchResult) {
  return (
    <Link key={result.id} className="catalog__card" to={`/projects/${result.slug}`}>
      {result.hero_image ? (
        <img className="catalog__card-img" src={result.hero_image} alt={result.title} loading="lazy" />
      ) : (
        <div className="catalog__card-img catalog__card-img--placeholder">No image</div>
      )}
      <div className="catalog__card-body">
        {result.category ? <p className="eyebrow">{result.category}</p> : null}
        <h3>{result.title}</h3>
        {result.client_name ? <p className="detail__brand">{result.client_name}</p> : null}
        <p>{summarize(result.summary ?? undefined)}</p>
      </div>
    </Link>
  );
}

export function CatalogPage({
  searchQuery,
  onSearchQueryChange,
  onSuggestionsChange,
  renderSearchControls = true,
}: CatalogPageProps) {
  const [projects, setProjects] = useState<Project[]>([]);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<FilterState>(EMPTY_FILTERS);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const searchRequestRef = useRef(0);

  useEffect(() => {
    let cancelled = false;

    fetchProjects()
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

  const isSearchControlled = searchQuery !== undefined;
  const searchDraft = isSearchControlled ? searchQuery : filters.searchDraft;
  const trimmedSearchDraft = searchDraft.trim();

  function updateSearchDraft(nextQuery: string) {
    if (isSearchControlled) {
      onSearchQueryChange?.(nextQuery);
      return;
    }

    setFilters((prev) => ({ ...prev, searchDraft: nextQuery }));
  }

  function clearFilters() {
    if (isSearchControlled) {
      onSearchQueryChange?.('');
      setFilters((prev) => ({ ...prev, category: '' }));
      return;
    }

    setFilters(EMPTY_FILTERS);
  }

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);

    const requestId = searchRequestRef.current + 1;
    searchRequestRef.current = requestId;

    if (trimmedSearchDraft.length < MIN_REMOTE_QUERY_LENGTH) {
      setSearchResults([]);
      setSearching(false);
      return;
    }

    setSearchResults([]);
    setSearching(true);

    debounceRef.current = setTimeout(() => {
      searchProjects({ q: trimmedSearchDraft, category: filters.category || undefined })
        .then((response) => {
          if (searchRequestRef.current !== requestId) return;
          setSearchResults(response.data);
          setSearching(false);
        })
        .catch(() => {
          if (searchRequestRef.current !== requestId) return;
          setSearchResults([]);
          setSearching(false);
        });
    }, SEARCH_DEBOUNCE_MS);

    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [filters.category, trimmedSearchDraft]);

  const categories = useMemo(
    () => uniqueSorted(projects.map((project) => project.category)),
    [projects],
  );

  const filteredProjects = useMemo(
    () =>
      projects.filter((project) => {
        if (filters.category && project.category !== filters.category) {
          return false;
        }
        return true;
      }),
    [projects, filters.category],
  );

  const localMatches = useMemo(
    () =>
      trimmedSearchDraft
        ? filteredProjects.filter((project) => matchesProject(project, trimmedSearchDraft))
        : filteredProjects,
    [filteredProjects, trimmedSearchDraft],
  );

  const suggestions = useMemo(
    () => buildSuggestions(projects, filters.category, trimmedSearchDraft),
    [projects, filters.category, trimmedSearchDraft],
  );

  useEffect(() => {
    onSuggestionsChange?.(suggestions);
  }, [onSuggestionsChange, suggestions]);

  const shouldUseRemoteResults = trimmedSearchDraft.length >= MIN_REMOTE_QUERY_LENGTH && searchResults.length > 0;
  const displayedCount = shouldUseRemoteResults ? searchResults.length : localMatches.length;

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
        {renderSearchControls && (
          <input
            className="catalog__filter-input"
            type="text"
            placeholder="Search by project name…"
            value={searchDraft}
            onChange={(event) => updateSearchDraft(event.target.value)}
            aria-label="Search projects"
          />
        )}

        <select
          className="catalog__filter-select"
          value={filters.category}
          onChange={(event) => setFilters((prev) => ({ ...prev, category: event.target.value }))}
          aria-label="Filter by category"
        >
          <option value="">Category</option>
          {categories.map((category) => (
            <option key={category} value={category}>
              {category}
            </option>
          ))}
        </select>

        {(searchDraft || filters.category) && (
          <button className="catalog__filter-clear" type="button" onClick={clearFilters}>
            Clear filters
          </button>
        )}
      </div>

      {renderSearchControls && trimmedSearchDraft && suggestions.length > 0 && (
        <div className="catalog__suggestions" aria-label="Suggested keywords">
          {suggestions.map((suggestion) => (
            <button
              key={suggestion}
              className="catalog__suggestion"
              type="button"
              onClick={() => updateSearchDraft(suggestion)}
            >
              {suggestion}
            </button>
          ))}
        </div>
      )}

      {trimmedSearchDraft ? (
        <>
          <p className="catalog__count">
            {displayedCount} resultado{displayedCount !== 1 ? 's' : ''}
            {searching ? ' · refinando búsqueda…' : ''}
          </p>

          {displayedCount === 0 && !searching && (
            <div className="card card--muted">
              <p>No projects match your search.</p>
            </div>
          )}

          {displayedCount > 0 && (
            <div className="catalog__grid">
              {shouldUseRemoteResults
                ? searchResults.map((result) => renderSearchCard(result))
                : localMatches.map((project) => renderProductCard(project))}
            </div>
          )}
        </>
      ) : (
        <>
          <p className="catalog__count">
            {filteredProjects.length} project{filteredProjects.length !== 1 ? 's' : ''}
          </p>

          {filteredProjects.length === 0 ? (
            <div className="card card--muted">
              <p>No projects match the current filters.</p>
            </div>
          ) : (
            <div className="catalog__grid">{filteredProjects.map((project) => renderProductCard(project))}</div>
          )}
        </>
      )}
    </section>
  );
}
