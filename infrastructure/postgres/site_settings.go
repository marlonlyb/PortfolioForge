package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

type SiteSettingsRepository struct {
	db *pgxpool.Pool
}

const ensureSiteSettingsTableSQL = `
CREATE TABLE IF NOT EXISTS site_settings (
	singleton_key BOOLEAN NOT NULL DEFAULT TRUE,
	public_hero_logo_url TEXT,
	public_hero_logo_alt TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT site_settings_singleton_pk PRIMARY KEY (singleton_key),
	CONSTRAINT site_settings_singleton_ck CHECK (singleton_key = TRUE)
)`

func NewSiteSettingsRepository(db *pgxpool.Pool) *SiteSettingsRepository {
	return &SiteSettingsRepository{db: db}
}

func (r *SiteSettingsRepository) ensureSchema(ctx context.Context) error {
	if _, err := r.db.Exec(ctx, ensureSiteSettingsTableSQL); err != nil {
		return fmt.Errorf("postgres.SiteSettingsRepository.ensureSchema: %w", err)
	}

	return nil
}

func (r *SiteSettingsRepository) Get(ctx context.Context) (model.SiteSettings, error) {
	var settings model.SiteSettings

	if err := r.ensureSchema(ctx); err != nil {
		return model.SiteSettings{}, err
	}

	err := r.db.QueryRow(ctx, `
		SELECT
			COALESCE(public_hero_logo_url, ''),
			COALESCE(public_hero_logo_alt, ''),
			COALESCE(EXTRACT(EPOCH FROM updated_at)::bigint, 0)
		FROM site_settings
		WHERE singleton_key = TRUE`).Scan(
		&settings.PublicHeroLogoURL,
		&settings.PublicHeroLogoAlt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.SiteSettings{}, nil
		}
		return model.SiteSettings{}, fmt.Errorf("postgres.SiteSettingsRepository.Get: %w", err)
	}

	return settings, nil
}

func (r *SiteSettingsRepository) Save(ctx context.Context, settings model.SiteSettings) (model.SiteSettings, error) {
	trimmedURL := strings.TrimSpace(settings.PublicHeroLogoURL)
	trimmedAlt := strings.TrimSpace(settings.PublicHeroLogoAlt)

	if err := r.ensureSchema(ctx); err != nil {
		return model.SiteSettings{}, err
	}

	err := r.db.QueryRow(ctx, `
		INSERT INTO site_settings (
			singleton_key,
			public_hero_logo_url,
			public_hero_logo_alt,
			updated_at
		) VALUES (TRUE, $1, $2, NOW())
		ON CONFLICT (singleton_key) DO UPDATE SET
			public_hero_logo_url = EXCLUDED.public_hero_logo_url,
			public_hero_logo_alt = EXCLUDED.public_hero_logo_alt,
			updated_at = EXCLUDED.updated_at
		RETURNING
			COALESCE(public_hero_logo_url, ''),
			COALESCE(public_hero_logo_alt, ''),
			COALESCE(EXTRACT(EPOCH FROM updated_at)::bigint, 0)`,
		trimmedURL,
		trimmedAlt,
	).Scan(
		&settings.PublicHeroLogoURL,
		&settings.PublicHeroLogoAlt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return model.SiteSettings{}, fmt.Errorf("postgres.SiteSettingsRepository.Save: %w", err)
	}

	return settings, nil
}
