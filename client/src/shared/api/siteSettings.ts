import { httpGet, httpPut } from './http';
import type { SaveSiteSettingsPayload, SiteSettings } from '../types/siteSettings';

export function fetchAdminSiteSettings(): Promise<SiteSettings> {
  return httpGet<SiteSettings>('/api/v1/admin/site-settings');
}

export function updateAdminSiteSettings(payload: SaveSiteSettingsPayload): Promise<SiteSettings> {
  return httpPut<SiteSettings>('/api/v1/admin/site-settings', payload);
}

export function fetchPublicSiteSettings(): Promise<SiteSettings> {
  return httpGet<SiteSettings>('/api/v1/public/site-settings');
}
