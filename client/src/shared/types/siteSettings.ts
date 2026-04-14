export interface SiteSettings {
  public_hero_logo_url?: string;
  public_hero_logo_alt?: string;
  updated_at?: number;
}

export interface SaveSiteSettingsPayload {
  public_hero_logo_url: string;
  public_hero_logo_alt: string;
}
