package main

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/infrastructure/casestudy"
)

func processActiveProjectMetadataRow(
	ctx context.Context,
	row activeProjectMetadataRow,
	force bool,
	locales []string,
	fetch func(context.Context, string, string) (casestudy.CanonicalProjectMetadata, error),
	update func(context.Context, uuid.UUID, casestudy.CanonicalProjectMetadata) error,
	localize func(context.Context, uuid.UUID, []string) error,
	refresh func(context.Context, uuid.UUID) error,
) (bool, string) {
	if !force && row.IndustryType != "" && row.FinalProduct != "" {
		return false, ""
	}

	metadata, err := fetch(ctx, row.SourceMarkdownURL, row.Slug)
	if err != nil {
		return false, formatProjectMetadataBackfillFlag(row, err)
	}
	if err := update(ctx, row.ID, metadata); err != nil {
		return false, formatProjectMetadataBackfillFlag(row, err)
	}
	if err := localize(ctx, row.ID, locales); err != nil {
		return false, formatProjectMetadataBackfillFlag(row, err)
	}
	if err := refresh(ctx, row.ID); err != nil {
		return false, formatProjectMetadataBackfillFlag(row, err)
	}

	return true, ""
}

func formatProjectMetadataBackfillFlag(row activeProjectMetadataRow, err error) string {
	return row.Slug + " (" + row.ID.String() + "): " + err.Error()
}

func TestProcessActiveProjectMetadataRowFlagsFailures(t *testing.T) {
	projectID := uuid.New()
	row := activeProjectMetadataRow{ID: projectID, Slug: "portfolioforge", SourceMarkdownURL: "https://example.com/portfolioforge.md"}

	updated, flagged := processActiveProjectMetadataRow(
		context.Background(),
		row,
		false,
		nil,
		func(context.Context, string, string) (casestudy.CanonicalProjectMetadata, error) {
			return casestudy.CanonicalProjectMetadata{}, errors.New("unexpected status 404")
		},
		func(context.Context, uuid.UUID, casestudy.CanonicalProjectMetadata) error { return nil },
		func(context.Context, uuid.UUID, []string) error { return nil },
		func(context.Context, uuid.UUID) error { return nil },
	)

	if updated {
		t.Fatal("expected row not to be updated when fetch fails")
	}
	if flagged != "portfolioforge ("+projectID.String()+"): unexpected status 404" {
		t.Fatalf("flagged = %q", flagged)
	}
}

func TestProcessActiveProjectMetadataRowSkipsCompleteRowsWithoutForce(t *testing.T) {
	updated, flagged := processActiveProjectMetadataRow(
		context.Background(),
		activeProjectMetadataRow{
			ID:           uuid.New(),
			Slug:         "portfolioforge",
			IndustryType: "metalworking",
			FinalProduct: "Operator HMI panel",
		},
		false,
		nil,
		func(context.Context, string, string) (casestudy.CanonicalProjectMetadata, error) {
			t.Fatal("fetch should not be called for complete rows without force")
			return casestudy.CanonicalProjectMetadata{}, nil
		},
		func(context.Context, uuid.UUID, casestudy.CanonicalProjectMetadata) error {
			t.Fatal("update should not be called for complete rows without force")
			return nil
		},
		func(context.Context, uuid.UUID, []string) error {
			t.Fatal("localize should not be called for complete rows without force")
			return nil
		},
		func(context.Context, uuid.UUID) error {
			t.Fatal("refresh should not be called for complete rows without force")
			return nil
		},
	)

	if updated {
		t.Fatal("expected row to be skipped")
	}
	if flagged != "" {
		t.Fatalf("flagged = %q, want empty", flagged)
	}
}

func TestProcessActiveProjectMetadataRowBackfillsLegacyIndustryValuesWhenForced(t *testing.T) {
	projectID := uuid.New()
	row := activeProjectMetadataRow{
		ID:           projectID,
		Slug:         "portfolioforge",
		IndustryType: "metalworking",
		FinalProduct: "Panel legado",
	}

	var updatedMetadata casestudy.CanonicalProjectMetadata
	localizeCalled := false
	refreshCalled := false
	updated, flagged := processActiveProjectMetadataRow(
		context.Background(),
		row,
		true,
		[]string{"en"},
		func(context.Context, string, string) (casestudy.CanonicalProjectMetadata, error) {
			return casestudy.ParseCanonicalProjectMetadata(`---
title: PortfolioForge
industry_type: metalworking
final_product: Panel HMI para diagnóstico y monitoreo
---

# PortfolioForge

## Metadata
- Industry Type: metalworking
- Final Product: Panel HMI para diagnóstico y monitoreo
`, row.Slug)
		},
		func(_ context.Context, id uuid.UUID, metadata casestudy.CanonicalProjectMetadata) error {
			if id != projectID {
				t.Fatalf("update project id = %s, want %s", id, projectID)
			}
			updatedMetadata = metadata
			return nil
		},
		func(_ context.Context, id uuid.UUID, locales []string) error {
			if id != projectID {
				t.Fatalf("localize project id = %s, want %s", id, projectID)
			}
			if len(locales) != 1 || locales[0] != "en" {
				t.Fatalf("locales = %#v", locales)
			}
			localizeCalled = true
			return nil
		},
		func(_ context.Context, id uuid.UUID) error {
			if id != projectID {
				t.Fatalf("refresh project id = %s, want %s", id, projectID)
			}
			refreshCalled = true
			return nil
		},
	)

	if !updated {
		t.Fatal("expected legacy row to be backfilled when force is enabled")
	}
	if flagged != "" {
		t.Fatalf("flagged = %q, want empty", flagged)
	}
	if updatedMetadata.IndustryType != "metalurgia" {
		t.Fatalf("updated industry_type = %q, want metalurgia", updatedMetadata.IndustryType)
	}
	if updatedMetadata.IndustryType == "metalworking" {
		t.Fatal("legacy industry_type leaked into backfill update payload")
	}
	if !localizeCalled {
		t.Fatal("expected localization backfill to run")
	}
	if !refreshCalled {
		t.Fatal("expected search refresh to run")
	}
}
