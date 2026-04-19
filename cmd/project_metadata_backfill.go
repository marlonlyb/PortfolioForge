package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/infrastructure/casestudy"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/infrastructure/postgres"
	"github.com/marlonlyb/portfolioforge/internal/markdownpolicy"
)

type activeProjectMetadataRow struct {
	ID                uuid.UUID
	Slug              string
	SourceMarkdownURL string
	IndustryType      string
	FinalProduct      string
}

func runProjectMetadataBackfill(args []string) error {
	fs := flag.NewFlagSet("project-metadata-backfill", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	projectID := fs.String("project-id", "", "Optional UUID to backfill a single project.")
	force := fs.Bool("force", false, "Overwrite existing industry_type/final_product values.")
	localesRaw := fs.String("locale", "", "Optional comma-separated locale subset for industry_type/final_product localization regeneration.")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: go run ./cmd project-metadata-backfill [--project-id <uuid>] [--force] [--locale ca,en]")
		fmt.Fprintln(fs.Output(), "")
		fmt.Fprintln(fs.Output(), "Backfills active projects from their remote source_markdown_url canonical markdown.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := validateBackfillEnvironment(); err != nil {
		return err
	}

	dbPool, err := NewDBConnection()
	if err != nil {
		return err
	}
	defer dbPool.Close()

	ctx := context.Background()
	rows, err := loadActiveProjectsForMetadataBackfill(ctx, dbPool, strings.TrimSpace(*projectID))
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		log.Println("project-metadata-backfill: no active projects matched")
		return nil
	}

	locales, err := parseBackfillLocales(*localesRaw)
	if err != nil {
		return err
	}
	projectRepo := postgres.NewProjectRepository(dbPool)
	localizationService := localization.NewService(
		postgres.NewProjectLocalizationRepository(dbPool),
		localization.NewOpenAITranslator(os.Getenv("OPENAI_API_KEY")),
	)
	localizationBackfill := localization.NewBackfillService(projectRepo, localizationService)
	searchRepo := postgres.NewSearchRepository(dbPool, IsSemanticSearchEnabled(), nil)

	client := &http.Client{Timeout: 30 * time.Second}
	flagged := make([]string, 0)
	updated := 0
	for _, row := range rows {
		if !*force && strings.TrimSpace(row.IndustryType) != "" && strings.TrimSpace(row.FinalProduct) != "" {
			continue
		}
		metadata, err := fetchCanonicalProjectMetadata(ctx, client, row.SourceMarkdownURL, row.Slug)
		if err != nil {
			flagged = append(flagged, fmt.Sprintf("%s (%s): %v", row.Slug, row.ID, err))
			continue
		}
		if _, err := dbPool.Exec(ctx, `
			UPDATE products
			SET industry_type = $1,
			    final_product = $2,
			    updated_at = EXTRACT(EPOCH FROM NOW())::int
			WHERE id = $3`, metadata.IndustryType, metadata.FinalProduct, row.ID); err != nil {
			flagged = append(flagged, fmt.Sprintf("%s (%s): update failed: %v", row.Slug, row.ID, err))
			continue
		}
		if err := localizationBackfill.BackfillProject(ctx, row.ID, locales); err != nil {
			flagged = append(flagged, fmt.Sprintf("%s (%s): localization failed: %v", row.Slug, row.ID, err))
			continue
		}
		if err := searchRepo.RefreshSearchDocument(ctx, row.ID); err != nil {
			flagged = append(flagged, fmt.Sprintf("%s (%s): search refresh failed: %v", row.Slug, row.ID, err))
			continue
		}
		updated++
		log.Printf("project-metadata-backfill: updated %s (%s)", row.Slug, row.ID)
	}

	log.Printf("project-metadata-backfill: updated=%d flagged=%d", updated, len(flagged))
	for _, item := range flagged {
		log.Printf("project-metadata-backfill: flagged %s", item)
	}
	if len(flagged) > 0 {
		return fmt.Errorf("project-metadata-backfill finished with %d flagged project(s)", len(flagged))
	}
	return nil
}

func loadActiveProjectsForMetadataBackfill(ctx context.Context, db *pgxpool.Pool, projectID string) ([]activeProjectMetadataRow, error) {
	query := `
		SELECT id,
		       COALESCE(NULLIF(slug, ''), regexp_replace(lower(COALESCE(NULLIF(name, ''), product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
		       COALESCE(source_markdown_url, '') AS source_markdown_url,
		       COALESCE(industry_type, '') AS industry_type,
		       COALESCE(final_product, '') AS final_product
		FROM products
		WHERE active = TRUE`
	args := []interface{}{}
	if projectID != "" {
		id, err := uuid.Parse(projectID)
		if err != nil {
			return nil, fmt.Errorf("project-id inválido: %w", err)
		}
		query += ` AND id = $1`
		args = append(args, id)
	}
	query += ` ORDER BY created_at DESC, id DESC`

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]activeProjectMetadataRow, 0)
	for rows.Next() {
		var item activeProjectMetadataRow
		if err := rows.Scan(&item.ID, &item.Slug, &item.SourceMarkdownURL, &item.IndustryType, &item.FinalProduct); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func fetchCanonicalProjectMetadata(ctx context.Context, client *http.Client, sourceURL string, fallbackSlug string) (casestudy.CanonicalProjectMetadata, error) {
	normalizedURL := strings.TrimSpace(sourceURL)
	if err := markdownpolicy.ValidateSourceURL(normalizedURL); err != nil {
		return casestudy.CanonicalProjectMetadata{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, normalizedURL, nil)
	if err != nil {
		return casestudy.CanonicalProjectMetadata{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return casestudy.CanonicalProjectMetadata{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return casestudy.CanonicalProjectMetadata{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024+1))
	if err != nil {
		return casestudy.CanonicalProjectMetadata{}, err
	}
	return casestudy.ParseCanonicalProjectMetadata(string(body), fallbackSlug)
}
