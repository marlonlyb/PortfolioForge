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
	GetAssistantContextBySlug(ctx context.Context, slug string) (model.ProjectAssistantContext, error)
}

// ProjectWriter defines write operations for projects.
type ProjectWriter interface {
	Create(ctx context.Context, m *model.Project) error
	Update(ctx context.Context, m *model.Project) error
}

// AdminCatalogRepository defines the canonical admin catalog contract for projects.
// It still persists against the legacy `products` storage during the transition.
type AdminCatalogRepository interface {
	Create(m *model.AdminProjectWrite) error
	Update(m *model.AdminProjectWrite) error
	Delete(id uuid.UUID) error
	UpdateActive(id uuid.UUID, active bool) error
	UpdateVariant(v model.AdminProjectVariantInput, projectID uuid.UUID) error
	DeleteVariantByID(id uuid.UUID) error

	GetAdminByID(id uuid.UUID) (model.AdminProject, error)
	GetAdminAll() ([]model.AdminProject, error)

	CreateVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error
	ReplaceMedia(projectID uuid.UUID, media []model.ProjectMedia) error
	DeleteVariantsByProductID(projectID uuid.UUID) error
}

// PublicProductCompatRepository isolates legacy `/public/products` reads.
// Storage still comes from the legacy products tables, but the compat flow stays
// out of the canonical admin/public project contracts.
type PublicProductCompatRepository interface {
	GetStoreByID(id uuid.UUID) (model.StoreProduct, error)
	GetStoreAll() ([]model.StoreProduct, error)
}

// AdminCatalogService defines the canonical service contract for admin project catalog operations.
type AdminCatalogService interface {
	Create(m *model.AdminProjectWrite) error
	Update(m *model.AdminProjectWrite) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, active bool) (model.AdminProject, error)
	CreateVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error
	ReplaceVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error
	ReplaceMedia(projectID uuid.UUID, media []model.ProjectMedia) error

	GetAdminByID(id uuid.UUID) (model.AdminProject, error)
	GetAdminAll() ([]model.AdminProject, error)
}

// PublicProductCompatService isolates legacy `/public/products` reads.
type PublicProductCompatService interface {
	GetStoreByID(id uuid.UUID) (model.StoreProduct, error)
	GetStoreAll() ([]model.StoreProduct, error)
}

// ProjectRepository combines read and write operations for projects.
type ProjectRepository interface {
	ProjectReader
	ProjectWriter
}
