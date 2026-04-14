package model

// SiteSettings stores global configurable content for the public portfolio.
type SiteSettings struct {
	PublicHeroLogoURL string `json:"public_hero_logo_url,omitempty"`
	PublicHeroLogoAlt string `json:"public_hero_logo_alt,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}
