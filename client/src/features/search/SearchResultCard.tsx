import { Link } from 'react-router-dom';

import type { SearchResult } from '../../shared/types/search';
import {
  buildSearchMatchContext,
  formatEvidenceField,
  formatMatchType,
  hasMatchedText,
  truncateEvidenceText,
} from './matchContext';

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
  const searchMatchContext = buildSearchMatchContext(result);
  const evidence = searchMatchContext?.evidence ?? [];
  const hasExplanation = Boolean(searchMatchContext?.explanation?.trim());
  const hasEvidence = evidence.length > 0;
  const showMatchDetails = hasExplanation || hasEvidence;

  return (
    <Link
      className="search-results__card"
      to={`/projects/${result.slug}`}
      state={searchMatchContext ? { searchMatchContext } : undefined}
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="search-results__card-media">
        {result.hero_image ? (
          <img className="search-results__card-image" src={result.hero_image} alt={result.title} loading="lazy" />
        ) : (
          <div className="search-results__card-image search-results__card-image--placeholder">
            Project visual
          </div>
        )}
      </div>

      <div className="search-results__card-body">
        <div className="search-results__card-meta">
          <p className="eyebrow">{result.category}</p>
          {result.client_name && <p className="search-results__card-client">{result.client_name}</p>}
        </div>
        <h3>{result.title}</h3>
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
        {showMatchDetails && (
          <section className="search-results__card-match-details" aria-label="Detalles de coincidencia">
            <p className="search-results__card-match-title">Por qué coincide</p>

            {hasExplanation && searchMatchContext?.explanation && (
              <p className="search-results__card-explanation">{searchMatchContext.explanation}</p>
            )}

            {hasEvidence && (
              <div className="search-results__card-evidence">
                <p className="search-results__card-evidence-title">Evidencia utilizada</p>
                <ul className="search-results__card-evidence-list">
                  {evidence.map((evidenceItem, evidenceIndex) => (
                    <li
                      key={`${result.id}-${evidenceItem.field}-${evidenceItem.match_type}-${evidenceIndex}`}
                      className="search-results__card-evidence-item"
                    >
                      <div className="search-results__card-evidence-meta">
                        <span className="search-results__card-evidence-field">
                          {formatEvidenceField(evidenceItem.field)}
                        </span>
                        <span className="search-results__card-evidence-type">
                          {formatMatchType(evidenceItem.match_type)}
                        </span>
                      </div>

                      {hasMatchedText(evidenceItem) && (
                        <p className="search-results__card-evidence-text">
                          “{truncateEvidenceText(evidenceItem.matched_text)}”
                        </p>
                      )}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </section>
        )}

        <span className="search-results__card-link">Open case study</span>
      </div>
    </Link>
  );
}
