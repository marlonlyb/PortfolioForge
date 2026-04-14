package product

import (
	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type Repository interface {
	Create(m *model.Product) error
	Update(m *model.Product) error
	Delete(ID uuid.UUID) error
	UpdateActive(ID uuid.UUID, active bool) error
	UpdateVariant(v model.StoreProductVariant) error
	DeleteVariantByID(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (model.Product, error)
	GetAll() (model.Products, error)
	GetStoreByID(ID uuid.UUID) (model.StoreProduct, error)
	GetStoreByIDAdmin(ID uuid.UUID) (model.StoreProduct, error)
	GetStoreAll() ([]model.StoreProduct, error)
	GetStoreAllAdmin() ([]model.StoreProduct, error)

	CreateVariants(productID uuid.UUID, variants []model.StoreProductVariant) error
	ReplaceMedia(productID uuid.UUID, media []model.ProjectMedia) error
	DeleteVariantsByProductID(productID uuid.UUID) error
}

type Service interface {
	Create(m *model.Product) error
	Update(m *model.Product) error
	Delete(ID uuid.UUID) error
	UpdateStatus(ID uuid.UUID, active bool) (model.StoreProduct, error)
	CreateVariants(productID uuid.UUID, variants []model.StoreProductVariant) error
	ReplaceVariants(productID uuid.UUID, variants []model.StoreProductVariant) error
	ReplaceMedia(productID uuid.UUID, media []model.ProjectMedia) error

	GetByID(ID uuid.UUID) (model.Product, error)
	GetAll() (model.Products, error)
	GetStoreByID(ID uuid.UUID) (model.StoreProduct, error)
	GetStoreByIDAdmin(ID uuid.UUID) (model.StoreProduct, error)
	GetStoreAll() ([]model.StoreProduct, error)
	GetStoreAllAdmin() ([]model.StoreProduct, error)
}
