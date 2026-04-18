package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/infrastructure/postgres"
)

func runLocalizationBackfill(args []string) error {
	fs := flag.NewFlagSet("localization-backfill", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	projectID := fs.String("project-id", "", "Project UUID to regenerate first. Recommended for smoke verification.")
	localesRaw := fs.String("locale", "", "Optional comma-separated locale subset (ca,en,de). Defaults to all derived locales.")

	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: go run ./cmd localization-backfill [--project-id <uuid>] [--locale ca,en]")
		fmt.Fprintln(fs.Output(), "")
		fmt.Fprintln(fs.Output(), "Checklist before running in batch:")
		fmt.Fprintln(fs.Output(), "  1. Run first with --project-id for one known project.")
		fmt.Fprintln(fs.Output(), "  2. Verify public detail returns localized client_name with Spanish fallback.")
		fmt.Fprintln(fs.Output(), "  3. Verify public search localizes display fields only; do not expect multilingual matching/ranking changes.")
		fmt.Fprintln(fs.Output(), "  4. Verify admin translations expose client_name and preserve any manual overrides.")
		fmt.Fprintln(fs.Output(), "")
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
	projectRepo := postgres.NewProjectRepository(dbPool)
	projectCatalogRepo := postgres.NewProjectCatalogRepository(dbPool)
	localizationService := localization.NewService(
		postgres.NewProjectLocalizationRepository(dbPool),
		localization.NewOpenAITranslator(os.Getenv("OPENAI_API_KEY")),
	)
	backfillService := localization.NewBackfillService(projectRepo, localizationService)

	locales, err := parseBackfillLocales(*localesRaw)
	if err != nil {
		return err
	}

	projectIDs, err := selectBackfillProjects(strings.TrimSpace(*projectID), projectCatalogRepo)
	if err != nil {
		return err
	}

	log.Printf("localization-backfill: regenerating %d project(s) for locales=%v", len(projectIDs), effectiveBackfillLocales(locales))
	for _, id := range projectIDs {
		if err := backfillService.BackfillProject(ctx, id, locales); err != nil {
			return err
		}
		log.Printf("localization-backfill: regenerated project=%s", id)
	}

	log.Println("localization-backfill: done")
	return nil
}

func validateBackfillEnvironment() error {
	required := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME", "DB_SSL_MODE"}
	for _, key := range required {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			return fmt.Errorf("%s es obligatoria para localization-backfill", key)
		}
	}
	return nil
}

func parseBackfillLocales(raw string) ([]string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	return localization.NormalizeLocalesForBackfill(strings.Split(trimmed, ","))
}

func effectiveBackfillLocales(locales []string) []string {
	resolved, err := localization.NormalizeLocalesForBackfill(locales)
	if err != nil {
		return locales
	}
	return resolved
}

func selectBackfillProjects(projectID string, repo postgres.ProjectCatalogRepository) ([]uuid.UUID, error) {
	if projectID != "" {
		id, err := uuid.Parse(projectID)
		if err != nil {
			return nil, fmt.Errorf("project-id inválido: %w", err)
		}
		return []uuid.UUID{id}, nil
	}

	projects, err := repo.GetAdminAll()
	if err != nil {
		return nil, fmt.Errorf("load admin projects: %w", err)
	}

	ids := make([]uuid.UUID, 0, len(projects))
	for _, project := range projects {
		if project.ID == uuid.Nil {
			continue
		}
		ids = append(ids, project.ID)
	}
	return ids, nil
}
