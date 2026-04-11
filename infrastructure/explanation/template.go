package explanation

import (
	"context"
	"fmt"
	"strings"

	"github.com/marlonlyb/portfolioforge/model"
)

// TemplateExplanationProvider generates static, template-based explanations
// for why a project matched a search query.
type TemplateExplanationProvider struct{}

// NewTemplateExplanationProvider creates a new TemplateExplanationProvider.
func NewTemplateExplanationProvider() *TemplateExplanationProvider {
	return &TemplateExplanationProvider{}
}

// Explain generates a template-based explanation string from the evidence trace.
// The explanation includes the original query terms and matched evidence fields,
// referencing the concepts searched by the user (REQ-4.1).
func (t *TemplateExplanationProvider) Explain(_ context.Context, _ model.Project, evidence model.EvidenceTrace, query string) (string, error) {
	if len(evidence.Fields) == 0 {
		if query != "" {
			return fmt.Sprintf("Este proyecto coincide con la búsqueda «%s».", query), nil
		}
		return "Este proyecto coincide con la búsqueda.", nil
	}

	var parts []string
	for _, f := range evidence.Fields {
		displayName := fieldDisplayName(f.Field)
		if f.MatchedText != "" {
			parts = append(parts, fmt.Sprintf("%s («%s»)", displayName, f.MatchedText))
		} else {
			parts = append(parts, displayName)
		}
	}

	if query != "" {
		return fmt.Sprintf("Para la búsqueda «%s», este proyecto coincide en: %s.", query, strings.Join(parts, ", ")), nil
	}
	return fmt.Sprintf("Este proyecto coincide con la búsqueda en: %s.", strings.Join(parts, ", ")), nil
}

// fieldDisplayName maps internal field names to human-readable Spanish names.
func fieldDisplayName(field string) string {
	names := map[string]string{
		"name":             "nombre",
		"description":      "descripción",
		"client_name":      "cliente",
		"solution_summary": "resumen de solución",
		"architecture":     "arquitectura",
		"business_goal":    "objetivo de negocio",
		"ai_usage":         "uso de IA",
		"technology":       "tecnología",
		"category":         "categoría",
	}
	if name, ok := names[field]; ok {
		return name
	}
	return field
}
