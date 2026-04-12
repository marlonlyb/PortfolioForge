import { Link } from 'react-router-dom';

import type { SearchResult } from '../../shared/types/search';

const MAX_TECHS_VISIBLE = 5;
const MAX_SUMMARY_LENGTH = 120;

interface SearchResultCardProps {
  result: SearchResult;
  index: number;
}

function truncateSummary(text: string | null): string {
  if (!text) return '';
  if (text.length <= MAX_SUMMARY_LENGTH) return text;
  return `${text.slice(0, MAX_SUMMARY_LENGTH)}…`;
}

export function SearchResultCard({ result, index }: SearchResultCardProps) {
  const visibleTechs = result.technologies.slice(0, MAX_TECHS_VISIBLE);
  const overflowCount = result.technologies.length - MAX_TECHS_VISIBLE;
  const delay = Math.min(index * 50, 500);

  return (
    <Link
      className="search-results__card"
      to={`/projects/${result.id}`}
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="search-results__card-body">
        <p className="eyebrow">{result.category}</p>
        <h3>{result.title}</h3>
        {result.client_name && (
          <p className="search-results__card-client">Cliente: {result.client_name}</p>
        )}
        {result.summary && (
          <p className="search-results__card-summary">{truncateSummary(result.summary)}</p>
        )}
        {result.technologies.length > 0 && (
          <div className="search-results__card-techs">
            {visibleTechs.map((tech) => (
              <span key={tech.id} className="search-results__card-tech">
                {tech.color && (
                  <span
                    className="search-results__card-tech-dot"
                    style={{ backgroundColor: tech.color }}
                  />
                )}
                {tech.name}
              </span>
            ))}
            {overflowCount > 0 && (
              <span className="search-results__card-tech-more">+{overflowCount} más</span>
            )}
          </div>
        )}
        {result.explanation && (
          <p className="search-results__card-explanation">{result.explanation}</p>
        )}
      </div>
    </Link>
  );
}
