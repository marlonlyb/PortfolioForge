package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/model"
)

// Project provides read-side operations for the public project portfolio.
type Project struct {
	Repository project.ProjectReader
}

// NewProject creates a new Project service with the given read-only repository.
func NewProject(pr project.ProjectReader) Project {
	return Project{Repository: pr}
}

// GetByID returns a single project by its ID.
func (s Project) GetByID(ctx context.Context, id uuid.UUID) (model.Project, error) {
	p, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return model.Project{}, fmt.Errorf("services.Project.GetByID: %w", err)
	}
	return p, nil
}

// GetBySlug returns a single published project by its slug.
func (s Project) GetBySlug(ctx context.Context, slug string) (model.Project, error) {
	p, err := s.Repository.GetBySlug(ctx, slug)
	if err != nil {
		return model.Project{}, fmt.Errorf("services.Project.GetBySlug: %w", err)
	}
	return p, nil
}

// ListPublished returns all published (active) projects.
func (s Project) ListPublished(ctx context.Context) ([]model.Project, error) {
	projects, err := s.Repository.ListPublished(ctx)
	if err != nil {
		return nil, fmt.Errorf("services.Project.ListPublished: %w", err)
	}
	return projects, nil
}
