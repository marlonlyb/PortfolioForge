package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

type projectLocalizationQueryer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type ProjectLocalizationRepository struct {
	db projectLocalizationQueryer
}

func NewProjectLocalizationRepository(db *pgxpool.Pool) *ProjectLocalizationRepository {
	return &ProjectLocalizationRepository{db: db}
}

func (r *ProjectLocalizationRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.ProjectLocalization, error) {
	rows, err := r.db.Query(ctx, `
		SELECT project_id, locale, field_key, value, mode, COALESCE(source_hash, ''), COALESCE(EXTRACT(EPOCH FROM updated_at)::bigint, 0)
		FROM project_localizations
		WHERE project_id = $1
		ORDER BY locale, field_key`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.ProjectLocalization
	for rows.Next() {
		var item model.ProjectLocalization
		if err := rows.Scan(&item.ProjectID, &item.Locale, &item.FieldKey, &item.Value, &item.Mode, &item.SourceHash, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ProjectLocalizationRepository) ListByProjectIDsAndLocale(ctx context.Context, projectIDs []uuid.UUID, locale string) (map[uuid.UUID][]model.ProjectLocalization, error) {
	result := map[uuid.UUID][]model.ProjectLocalization{}
	if len(projectIDs) == 0 {
		return result, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT project_id, locale, field_key, value, mode, COALESCE(source_hash, ''), COALESCE(EXTRACT(EPOCH FROM updated_at)::bigint, 0)
		FROM project_localizations
		WHERE locale = $1 AND project_id = ANY($2)
		ORDER BY project_id, field_key`, locale, projectIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ProjectLocalization
		if err := rows.Scan(&item.ProjectID, &item.Locale, &item.FieldKey, &item.Value, &item.Mode, &item.SourceHash, &item.UpdatedAt); err != nil {
			return nil, err
		}
		result[item.ProjectID] = append(result[item.ProjectID], item)
	}

	return result, rows.Err()
}

func (r *ProjectLocalizationRepository) UpsertAuto(ctx context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage, sourceHashes map[string]string) error {
	return r.upsert(ctx, projectID, locale, fields, model.LocalizationModeAuto, sourceHashes)
}

func (r *ProjectLocalizationRepository) UpsertManual(ctx context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage) error {
	return r.upsert(ctx, projectID, locale, fields, model.LocalizationModeManual, nil)
}

func (r *ProjectLocalizationRepository) upsert(ctx context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage, mode string, sourceHashes map[string]string) error {
	locale = strings.ToLower(strings.TrimSpace(locale))
	for fieldKey, value := range fields {
		sourceHash := ""
		if sourceHashes != nil {
			sourceHash = sourceHashes[fieldKey]
		}

		_, err := r.db.Exec(ctx, `
			INSERT INTO project_localizations (project_id, locale, field_key, value, mode, source_hash, updated_at)
			VALUES ($1, $2, $3, $4::jsonb, $5, NULLIF($6, ''), NOW())
			ON CONFLICT (project_id, locale, field_key) DO UPDATE SET
				value = EXCLUDED.value,
				mode = EXCLUDED.mode,
				source_hash = EXCLUDED.source_hash,
				updated_at = NOW()`, projectID, locale, fieldKey, string(value), mode, sourceHash)
		if err != nil {
			return fmt.Errorf("upsert localization %s/%s: %w", locale, fieldKey, err)
		}
	}

	return nil
}
