package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	project "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/model"
)

// ProjectCatalog is the canonical admin service for portfolio projects.
// It still operates on the legacy products storage model for compatibility.
type ProjectCatalog struct {
	Repository project.AdminCatalogRepository
}

func NewProjectCatalog(pr project.AdminCatalogRepository) ProjectCatalog {
	return ProjectCatalog{Repository: pr}
}

// PublicProductCompat isolates the legacy `/public/products` read flow.
type PublicProductCompat struct {
	Repository project.PublicProductCompatRepository
}

func NewPublicProductCompat(pr project.PublicProductCompatRepository) PublicProductCompat {
	return PublicProductCompat{Repository: pr}
}

func (p ProjectCatalog) Create(m *model.AdminProjectWrite) error {

	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}
	m.ID = ID
	m.Normalize()

	if m.Name == "" {
		return fmt.Errorf("%s", "product name is empty!")
	}

	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}

	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	m.CreatedAt = time.Now().Unix()

	err = p.Repository.Create(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Create(m)", err)
	}

	return nil
}

func (p ProjectCatalog) Update(m *model.AdminProjectWrite) error {
	m.Normalize()
	if m.ID == uuid.Nil {
		return fmt.Errorf("product: %w", model.ErrInvalidID)
	}

	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	m.UpdatedAt = time.Now().Unix()

	err := p.Repository.Update(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Update(m)", err)
	}
	return nil
}

func (p ProjectCatalog) Delete(ID uuid.UUID) error {
	err := p.Repository.Delete(ID)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Delete(ID)", err)
	}
	return nil
}

func (p ProjectCatalog) UpdateStatus(ID uuid.UUID, active bool) (model.AdminProject, error) {
	err := p.Repository.UpdateActive(ID, active)
	if err != nil {
		return model.AdminProject{}, fmt.Errorf("%s %w", "Repository.UpdateActive(ID, active)", err)
	}

	projectData, err := p.Repository.GetAdminByID(ID)
	if err != nil {
		return model.AdminProject{}, fmt.Errorf("%s %w", "Repository.GetAdminByID(ID)", err)
	}

	return projectData, nil
}

func (p ProjectCatalog) GetAdminByID(ID uuid.UUID) (model.AdminProject, error) {
	adminProject, err := p.Repository.GetAdminByID(ID)
	if err != nil {
		return model.AdminProject{}, fmt.Errorf("%s %w", "Repository.GetAdminByID(ID)", err)
	}

	return adminProject, nil
}

func (p ProjectCatalog) GetAdminAll() ([]model.AdminProject, error) {
	products, err := p.Repository.GetAdminAll()
	if err != nil {
		return nil, fmt.Errorf("%s %w", "Repository.GetAdminAll()", err)
	}

	return products, nil
}

func (p ProjectCatalog) CreateVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error {
	err := p.Repository.CreateVariants(projectID, variants)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.CreateVariants()", err)
	}
	return nil
}

func (p ProjectCatalog) ReplaceVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error {
	existingProject, err := p.Repository.GetAdminByID(projectID)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.GetAdminByID(projectID)", err)
	}

	existingByID := make(map[uuid.UUID]model.AdminProjectVariant, len(existingProject.Variants))
	for _, existing := range existingProject.Variants {
		existingByID[existing.ID] = existing
	}

	incomingIDs := make(map[uuid.UUID]struct{}, len(variants))
	newVariants := make([]model.AdminProjectVariantInput, 0)

	for _, variant := range variants {
		variantID, _ := uuid.Parse(variant.ID)

		if variantID != uuid.Nil {
			incomingIDs[variantID] = struct{}{}
			if _, exists := existingByID[variantID]; exists {
				err = p.Repository.UpdateVariant(variant, projectID)
				if err != nil {
					return fmt.Errorf("%s %w", "Repository.UpdateVariant()", err)
				}
				continue
			}
		}

		newVariants = append(newVariants, variant)
	}

	for existingID := range existingByID {
		if _, keep := incomingIDs[existingID]; keep {
			continue
		}

		err = p.Repository.DeleteVariantByID(existingID)
		if err != nil {
			return fmt.Errorf("%s %w", "Repository.DeleteVariantByID()", err)
		}
	}

	if len(newVariants) > 0 {
		err = p.Repository.CreateVariants(projectID, newVariants)
		if err != nil {
			return fmt.Errorf("%s %w", "Repository.CreateVariants()", err)
		}
	}

	return nil
}

func (p PublicProductCompat) GetStoreByID(ID uuid.UUID) (model.StoreProduct, error) {
	storeProduct, err := p.Repository.GetStoreByID(ID)
	if err != nil {
		return model.StoreProduct{}, fmt.Errorf("%s %w", "Repository.GetStoreByID(ID)", err)
	}

	if !storeProduct.Active {
		return model.StoreProduct{}, errors.New("product inactive")
	}

	return storeProduct, nil
}

func (p PublicProductCompat) GetStoreAll() ([]model.StoreProduct, error) {
	products, err := p.Repository.GetStoreAll()
	if err != nil {
		return nil, fmt.Errorf("%s %w", "Repository.GetStoreAll()", err)
	}

	return products, nil
}

func (p ProjectCatalog) ReplaceMedia(projectID uuid.UUID, media []model.ProjectMedia) error {
	err := p.Repository.ReplaceMedia(projectID, media)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.ReplaceMedia()", err)
	}

	return nil
}
