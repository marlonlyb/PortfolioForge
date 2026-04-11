import { httpGet } from '../../shared/api/http';
import type { SearchResponse } from '../../shared/types/search';

interface SearchParams {
  q?: string;
  category?: string;
  client?: string;
  technologies?: string;
  pageSize?: number;
  cursor?: string;
}

export function searchProjects(params: SearchParams = {}): Promise<SearchResponse> {
  const query = new URLSearchParams();
  if (params.q) query.set('q', params.q);
  if (params.category) query.set('category', params.category);
  if (params.client) query.set('client', params.client);
  if (params.technologies) query.set('technologies', params.technologies);
  if (params.pageSize) query.set('pageSize', String(params.pageSize));
  if (params.cursor) query.set('cursor', params.cursor);

  const qs = query.toString();
  return httpGet<SearchResponse>(`/api/v1/public/search${qs ? `?${qs}` : ''}`);
}
