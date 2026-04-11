package explanation

import (
	"context"

	"github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/model"
)

// AIExplanationProvider generates AI-powered explanations for search matches.
// In the MVP, it delegates to the fallback template provider.
// Future integration: bounded prompt, 2s timeout, validation (≥1 evidence term), template fallback.
type AIExplanationProvider struct {
	fallback search.ExplanationProvider
}

// NewAIExplanationProvider creates a new AIExplanationProvider with a fallback provider.
func NewAIExplanationProvider(fallback search.ExplanationProvider) *AIExplanationProvider {
	return &AIExplanationProvider{fallback: fallback}
}

// Explain generates an explanation for why a project matched.
// MVP: delegates to the fallback template provider.
func (a *AIExplanationProvider) Explain(ctx context.Context, project model.Project, evidence model.EvidenceTrace, query string) (string, error) {
	// MVP: delegate to template fallback
	// Future: call AI API with bounded prompt, 2s timeout, validation (≥1 evidence term)
	return a.fallback.Explain(ctx, project, evidence, query)
}
