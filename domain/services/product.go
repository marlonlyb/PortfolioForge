package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/domain/ports/product"
	"github.com/marlonlyb/portfolioforge/model"
)

type Product struct {
	Repository product.Repository
}

func NewProduct(pr product.Repository) Product {
	return Product{Repository: pr}
}

func (p Product) Create(m *model.Product) error {

	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}
	m.ID = ID

	if m.ProductName == "" {
		return fmt.Errorf("%s", "product name is empty!")
	}

	//Price could be 0 for free

	//Images could be empty?
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

func (p Product) Update(m *model.Product) error {
	if !m.HasID() {
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

func (p Product) Delete(ID uuid.UUID) error {
	err := p.Repository.Delete(ID)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Delete(ID)", err)
	}
	return nil
}

func (p Product) UpdateStatus(ID uuid.UUID, active bool) (model.StoreProduct, error) {
	err := p.Repository.UpdateActive(ID, active)
	if err != nil {
		return model.StoreProduct{}, fmt.Errorf("%s %w", "Repository.UpdateActive(ID, active)", err)
	}

	productData, err := p.Repository.GetStoreByIDAdmin(ID)
	if err != nil {
		return model.StoreProduct{}, fmt.Errorf("%s %w", "Repository.GetStoreByIDAdmin(ID)", err)
	}

	return productData, nil
}

func (p Product) GetByID(ID uuid.UUID) (model.Product, error) {
	product, err := p.Repository.GetByID(ID)
	if err != nil {
		return model.Product{}, fmt.Errorf("%s %w", "Repository.GetByID(ID)", err)
	}
	return product, nil
}

func (p Product) GetStoreByID(ID uuid.UUID) (model.StoreProduct, error) {
	storeProduct, err := p.Repository.GetStoreByID(ID)
	if err != nil {
		return model.StoreProduct{}, fmt.Errorf("%s %w", "Repository.GetStoreByID(ID)", err)
	}

	if !storeProduct.Active {
		return model.StoreProduct{}, errors.New("product inactive")
	}

	return storeProduct, nil
}

func (p Product) GetStoreByIDAdmin(ID uuid.UUID) (model.StoreProduct, error) {
	storeProduct, err := p.Repository.GetStoreByIDAdmin(ID)
	if err != nil {
		return model.StoreProduct{}, fmt.Errorf("%s %w", "Repository.GetStoreByIDAdmin(ID)", err)
	}

	return storeProduct, nil
}

func (p Product) GetAll() (model.Products, error) {
	products, err := p.Repository.GetAll()
	if err != nil {
		return model.Products{}, fmt.Errorf("%s %w", "Repository.GetAll()", err)
	}
	return products, nil
}

func (p Product) GetStoreAll() ([]model.StoreProduct, error) {
	products, err := p.Repository.GetStoreAll()
	if err != nil {
		return nil, fmt.Errorf("%s %w", "Repository.GetStoreAll()", err)
	}

	return products, nil
}

func (p Product) GetStoreAllAdmin() ([]model.StoreProduct, error) {
	products, err := p.Repository.GetStoreAllAdmin()
	if err != nil {
		return nil, fmt.Errorf("%s %w", "Repository.GetStoreAllAdmin()", err)
	}

	return products, nil
}

func (p Product) CreateVariants(productID uuid.UUID, variants []model.StoreProductVariant) error {
	err := p.Repository.CreateVariants(productID, variants)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.CreateVariants()", err)
	}
	return nil
}

func (p Product) ReplaceVariants(productID uuid.UUID, variants []model.StoreProductVariant) error {
	existingProduct, err := p.Repository.GetStoreByIDAdmin(productID)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.GetStoreByIDAdmin(productID)", err)
	}

	existingByID := make(map[uuid.UUID]model.StoreProductVariant, len(existingProduct.Variants))
	for _, existing := range existingProduct.Variants {
		existingByID[existing.ID] = existing
	}

	incomingIDs := make(map[uuid.UUID]struct{}, len(variants))
	newVariants := make([]model.StoreProductVariant, 0)

	for _, variant := range variants {
		variant.ProductID = productID

		if variant.ID != uuid.Nil {
			incomingIDs[variant.ID] = struct{}{}
			if _, exists := existingByID[variant.ID]; exists {
				err = p.Repository.UpdateVariant(variant)
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
		err = p.Repository.CreateVariants(productID, newVariants)
		if err != nil {
			return fmt.Errorf("%s %w", "Repository.CreateVariants()", err)
		}
	}

	return nil
}

func (p Product) ReplaceMedia(productID uuid.UUID, media []model.ProjectMedia) error {
	err := p.Repository.ReplaceMedia(productID, media)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.ReplaceMedia()", err)
	}

	return nil
}
