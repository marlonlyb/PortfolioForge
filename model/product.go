package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID       `json:"id"`
	ProductName string          `json:"product_name"`
	Price       float64         `json:"price"`
	Images      json.RawMessage `json:"images"`
	Description string          `json:"description"`
	Features    json.RawMessage `json:"features"`
	Name        string          `json:"name,omitempty"`
	Slug        string          `json:"slug,omitempty"`
	Category    string          `json:"category,omitempty"`
	Brand       string          `json:"brand,omitempty"`
	Active      bool            `json:"active"`
	CreatedAt   int64           `json:"created_at"`
	UpdatedAt   int64           `json:"updated_at"`
}

// SetStoreFields sets the extended store fields on the product.
func (p *Product) SetStoreFields(name, category, brand string, active bool) {
	p.Name = name
	p.Category = category
	p.Brand = brand
	p.Active = active
}

func (p Product) HasID() bool {
	return p.ID != uuid.Nil
}

type Products []Product

func (p Products) IsEmpty() bool {
	return len(p) == 0
}
