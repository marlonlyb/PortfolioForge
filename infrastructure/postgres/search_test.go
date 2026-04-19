package postgres

import (
	"strings"
	"testing"

	"github.com/marlonlyb/portfolioforge/model"
)

func TestBuildSearchQueryIncludesPayloadAndFilters(t *testing.T) {
	repo := &SearchRepository{semanticEnabled: true}
	query, args := repo.buildSearchQuery(model.SearchParams{
		Query:          "portfolio",
		QueryEmbedding: []float32{0.12, 0.34},
		Category:       "platform",
		Client:         "Acme",
		Technologies:   []string{"react", "go"},
	})

	checks := []string{
		"solution_summary",
		"COALESCE(p.industry_type, '') AS industry_type",
		"COALESCE(p.final_product, '') AS final_product",
		"COALESCE(p.images, '[]'::jsonb) AS images",
		"LOWER(COALESCE(p.category, '')) = LOWER($2)",
		"LOWER(COALESCE(NULLIF(p.client_name, ''), p.brand, '')) = LOWER($3)",
		"t.slug = ANY($4)",
		"search_embedding <=> $5::vector",
	}
	for _, check := range checks {
		if !strings.Contains(query, check) {
			t.Fatalf("query missing %q\n%s", check, query)
		}
	}
	if len(args) != 5 {
		t.Fatalf("args length = %d, want 5", len(args))
	}
}

func TestBuildSearchQueryOmitsSemanticClauseWithoutEmbedding(t *testing.T) {
	repo := &SearchRepository{semanticEnabled: true}
	query, args := repo.buildSearchQuery(model.SearchParams{Query: "portfolio"})

	if strings.Contains(query, "search_embedding <=>") {
		t.Fatalf("query should omit semantic clause when no embedding is available\n%s", query)
	}
	if len(args) != 1 {
		t.Fatalf("args length = %d, want 1", len(args))
	}
}
