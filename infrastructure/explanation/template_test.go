package explanation

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

func TestTemplateExplanationProvider_Explain(t *testing.T) {
	prov := NewTemplateExplanationProvider()
	ctx := context.Background()

	tests := []struct {
		name     string
		project  model.Project
		evidence model.EvidenceTrace
		query    string
		want     string
	}{
		{
			name:    "single field match with matched text and query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "architecture", MatchedText: "microservices", MatchType: model.MatchTypeFTS, Score: 0.5},
				},
			},
			query: "microservices",
			want:  "Para la búsqueda «microservices», este proyecto coincide en: arquitectura («microservices»).",
		},
		{
			name:    "multiple field match with query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "name", MatchedText: "react", MatchType: model.MatchTypeFTS, Score: 0.6},
					{Field: "technology", MatchedText: "react", MatchType: model.MatchTypeStructured, Score: 0.4},
				},
			},
			query: "react",
			want:  "Para la búsqueda «react», este proyecto coincide en: nombre («react»), tecnología («react»).",
		},
		{
			name:    "no fields — generic message without query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields:    []model.EvidenceField{},
			},
			query: "",
			want:  "Este proyecto coincide con la búsqueda.",
		},
		{
			name:    "no fields — generic message with query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields:    []model.EvidenceField{},
			},
			query: "SCADA",
			want:  "Este proyecto coincide con la búsqueda «SCADA».",
		},
		{
			name:    "all known fields with query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "name", MatchedText: "app", MatchType: model.MatchTypeFTS, Score: 0.5},
					{Field: "description", MatchedText: "app", MatchType: model.MatchTypeFTS, Score: 0.4},
					{Field: "client_name", MatchedText: "acme", MatchType: model.MatchTypeFTS, Score: 0.3},
					{Field: "category", MatchedText: "web", MatchType: model.MatchTypeStructured, Score: 0.2},
				},
			},
			query: "app acme",
			want:  "Para la búsqueda «app acme», este proyecto coincide en: nombre («app»), descripción («app»), cliente («acme»), categoría («web»).",
		},
		{
			name:    "all known fields without query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "name", MatchedText: "app", MatchType: model.MatchTypeFTS, Score: 0.5},
					{Field: "description", MatchedText: "app", MatchType: model.MatchTypeFTS, Score: 0.4},
					{Field: "client_name", MatchedText: "acme", MatchType: model.MatchTypeFTS, Score: 0.3},
					{Field: "category", MatchedText: "web", MatchType: model.MatchTypeStructured, Score: 0.2},
				},
			},
			query: "",
			want:  "Este proyecto coincide con la búsqueda en: nombre («app»), descripción («app»), cliente («acme»), categoría («web»).",
		},
		{
			name:    "unknown field falls back to raw name",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "custom_field", MatchedText: "x", MatchType: model.MatchTypeFTS, Score: 0.5},
				},
			},
			query: "x",
			want:  "Para la búsqueda «x», este proyecto coincide en: custom_field («x»).",
		},
		{
			name:    "field without matched text uses display name only",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "architecture", MatchedText: "", MatchType: model.MatchTypeFTS, Score: 0.5},
				},
			},
			query: "test",
			want:  "Para la búsqueda «test», este proyecto coincide en: arquitectura.",
		},
		{
			name:    "profile fields with query",
			project: model.Project{ID: uuid.New()},
			evidence: model.EvidenceTrace{
				ProjectID: uuid.New(),
				Fields: []model.EvidenceField{
					{Field: "solution_summary", MatchedText: "api", MatchType: model.MatchTypeFTS, Score: 0.5},
					{Field: "business_goal", MatchedText: "revenue", MatchType: model.MatchTypeFTS, Score: 0.4},
					{Field: "ai_usage", MatchedText: "gpt", MatchType: model.MatchTypeFTS, Score: 0.3},
				},
			},
			query: "api revenue",
			want:  "Para la búsqueda «api revenue», este proyecto coincide en: resumen de solución («api»), objetivo de negocio («revenue»), uso de IA («gpt»).",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prov.Explain(ctx, tt.project, tt.evidence, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Explain() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFieldDisplayName(t *testing.T) {
	tests := []struct {
		field string
		want  string
	}{
		{"name", "nombre"},
		{"description", "descripción"},
		{"client_name", "cliente"},
		{"solution_summary", "resumen de solución"},
		{"architecture", "arquitectura"},
		{"business_goal", "objetivo de negocio"},
		{"ai_usage", "uso de IA"},
		{"technology", "tecnología"},
		{"category", "categoría"},
		{"unknown_field", "unknown_field"},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			got := fieldDisplayName(tt.field)
			if got != tt.want {
				t.Errorf("fieldDisplayName(%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}
}
