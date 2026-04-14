package postgres

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/marlonlyb/portfolioforge/model"
)

type projectMediaQueryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type projectMediaExecer interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func fetchProjectMedia(ctx context.Context, queryer projectMediaQueryer, projectID uuid.UUID) ([]model.ProjectMedia, error) {
	rows, err := queryer.Query(ctx, `
		SELECT id, project_id, media_type, COALESCE(url, ''),
			COALESCE(thumbnail_url, ''), COALESCE(medium_url, ''), COALESCE(full_url, ''),
			COALESCE(caption, ''), COALESCE(alt_text, ''), sort_order, featured
		FROM project_media
		WHERE project_id = $1
		ORDER BY featured DESC, sort_order ASC, created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	media := make([]model.ProjectMedia, 0)
	for rows.Next() {
		var item model.ProjectMedia
		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.MediaType,
			&item.URL,
			&item.ThumbnailURL,
			&item.MediumURL,
			&item.FullURL,
			&item.Caption,
			&item.AltText,
			&item.SortOrder,
			&item.Featured,
		); err != nil {
			return nil, err
		}
		media = append(media, normalizeProjectMedia(item))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return media, nil
}

func replaceProjectMedia(ctx context.Context, execer projectMediaExecer, projectID uuid.UUID, media []model.ProjectMedia) error {
	if _, err := execer.Exec(ctx, "DELETE FROM project_media WHERE project_id = $1", projectID); err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO project_media (
			id, project_id, media_type, url, thumbnail_url, medium_url, full_url,
			caption, alt_text, sort_order, featured, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`

	featuredAssigned := false
	for index, raw := range media {
		item := normalizeProjectMedia(raw)
		if item.MediaType == "" {
			item.MediaType = "image"
		}
		if item.SortOrder == 0 && raw.SortOrder == 0 {
			item.SortOrder = index
		}
		if item.ID == uuid.Nil {
			newID, err := uuid.NewUUID()
			if err != nil {
				return err
			}
			item.ID = newID
		}
		if item.ProjectID == uuid.Nil {
			item.ProjectID = projectID
		}

		if item.Featured && !featuredAssigned {
			featuredAssigned = true
		} else {
			item.Featured = false
		}

		if _, err := execer.Exec(ctx, insertQuery,
			item.ID,
			projectID,
			item.MediaType,
			item.URL,
			item.ThumbnailURL,
			item.MediumURL,
			item.FullURL,
			item.Caption,
			item.AltText,
			item.SortOrder,
			item.Featured,
		); err != nil {
			return err
		}
	}

	return nil
}

func normalizeProjectMedia(item model.ProjectMedia) model.ProjectMedia {
	item.MediaType = strings.TrimSpace(item.MediaType)
	item.URL = strings.TrimSpace(item.URL)
	item.ThumbnailURL = strings.TrimSpace(item.ThumbnailURL)
	item.MediumURL = strings.TrimSpace(item.MediumURL)
	item.FullURL = strings.TrimSpace(item.FullURL)
	item.Caption = strings.TrimSpace(item.Caption)
	item.AltText = strings.TrimSpace(item.AltText)

	if item.ThumbnailURL == "" {
		item.ThumbnailURL = firstMediaValue(item.MediumURL, item.FullURL, item.URL)
	}
	if item.MediumURL == "" {
		item.MediumURL = firstMediaValue(item.FullURL, item.ThumbnailURL, item.URL)
	}
	if item.FullURL == "" {
		item.FullURL = firstMediaValue(item.URL, item.MediumURL, item.ThumbnailURL)
	}
	if item.URL == "" {
		item.URL = item.FullURL
	}
	if item.ThumbnailURL == "" {
		item.ThumbnailURL = firstMediaValue(item.MediumURL, item.FullURL, item.URL)
	}

	return item
}

func firstMediaValue(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}
