package project

import (
	"context"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

// ProjectReader defines read-only operations for projects.
type ProjectReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Project, error)
	GetBySlug(ctx context.Context, slug string) (model.Project, error)
	ListPublished(ctx context.Context) ([]model.Project, error)
	GetTechnologiesByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Technology, error)
}

// ProjectWriter defines write operations for projects.
type ProjectWriter interface {
	Create(ctx context.Context, m *model.Project) error
	Update(ctx context.Context, m *model.Project) error
}

// ProjectRepository combines read and write operations for projects.
type ProjectRepository interface {
	ProjectReader
	ProjectWriter
}
