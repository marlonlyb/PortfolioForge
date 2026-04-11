import type { SearchFilters as SearchFiltersType, SearchResult } from '../../shared/types/search';

interface SearchFiltersProps {
  filters: SearchFiltersType;
  results: SearchResult[];
  onFilterChange: (filters: SearchFiltersType) => void;
}

function uniqueSorted(values: (string | null | undefined)[]): string[] {
  return [...new Set(values.filter((v): v is string => Boolean(v?.trim())))].sort();
}

export function SearchFilters({ filters, results, onFilterChange }: SearchFiltersProps) {
  const categories = uniqueSorted(results.map((r) => r.category));
  const clients = uniqueSorted(results.map((r) => r.client_name));
  const techs = uniqueSorted(
    results.flatMap((r) => r.technologies.map((t) => t.name)),
  );

  const hasActiveFilters =
    filters.category !== null || filters.client !== null || filters.technologies.length > 0;

  const techBadgeCount = filters.technologies.length;
  const totalBadgeCount =
    (filters.category ? 1 : 0) + (filters.client ? 1 : 0) + techBadgeCount;

  function handleCategoryChange(value: string) {
    onFilterChange({
      ...filters,
      category: value || null,
    });
  }

  function handleClientChange(value: string) {
    onFilterChange({
      ...filters,
      client: value || null,
    });
  }

  function handleTechnologyToggle(tech: string) {
    const current = filters.technologies;
    const next = current.includes(tech)
      ? current.filter((t) => t !== tech)
      : [...current, tech];
    onFilterChange({
      ...filters,
      technologies: next,
    });
  }

  function handleClear() {
    onFilterChange({
      category: null,
      client: null,
      technologies: [],
    });
  }

  return (
    <aside className="search-filters">
      <h3>
        Filtros
        {totalBadgeCount > 0 && (
          <span className="chip chip--active" style={{ marginLeft: '0.5rem', fontSize: '0.75rem' }}>
            {totalBadgeCount}
          </span>
        )}
      </h3>

      {categories.length > 0 && (
        <div className="search-filters__section">
          <p className="search-filters__label">Categoría</p>
          <div className="search-filters__options">
            {categories.map((cat) => (
              <button
                key={cat}
                className={`chip${filters.category === cat ? ' chip--active' : ''}`}
                type="button"
                onClick={() => handleCategoryChange(filters.category === cat ? '' : cat)}
              >
                {cat}
              </button>
            ))}
          </div>
        </div>
      )}

      {clients.length > 0 && (
        <div className="search-filters__section">
          <p className="search-filters__label">Cliente</p>
          <div className="search-filters__options">
            {clients.map((client) => (
              <button
                key={client}
                className={`chip${filters.client === client ? ' chip--active' : ''}`}
                type="button"
                onClick={() => handleClientChange(filters.client === client ? '' : client)}
              >
                {client}
              </button>
            ))}
          </div>
        </div>
      )}

      {techs.length > 0 && (
        <div className="search-filters__section">
          <p className="search-filters__label">
            Tecnologías
            {techBadgeCount > 0 && (
              <span style={{ marginLeft: '0.35rem', fontWeight: 400 }}>({techBadgeCount})</span>
            )}
          </p>
          <div className="search-filters__options">
            {techs.map((tech) => (
              <button
                key={tech}
                className={`chip${filters.technologies.includes(tech) ? ' chip--active' : ''}`}
                type="button"
                onClick={() => handleTechnologyToggle(tech)}
              >
                {tech}
              </button>
            ))}
          </div>
        </div>
      )}

      {hasActiveFilters && (
        <button
          className="btn btn--ghost btn--small search-filters__clear"
          type="button"
          onClick={handleClear}
        >
          Limpiar filtros
        </button>
      )}
    </aside>
  );
}
