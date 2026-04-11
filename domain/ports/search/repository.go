package search

import (
	"context"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

// SearchRepository defines operations for hybrid retrieval search over projects.
type SearchRepository interface {
	Search(ctx context.Context, params model.SearchParams) ([]model.SearchResult, error)
	RefreshSearchDocument(ctx context.Context, projectID uuid.UUID) error
	RefreshAllDocuments(ctx context.Context) error
}
