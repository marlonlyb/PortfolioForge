package search

import (
	"context"

	"github.com/marlonlyb/portfolioforge/model"
)

// ExplanationProvider generates human-readable explanations for why a project
// matched a search query, given the evidence trace and the original query terms.
type ExplanationProvider interface {
	Explain(ctx context.Context, project model.Project, evidence model.EvidenceTrace, query string) (string, error)
}
