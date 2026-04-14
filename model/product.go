package model

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

// Product is the legacy persistence model backed by the `products` table.
// Prefer Project/AdminProject contracts outside the compat layer.
type Product struct {
	ID                uuid.UUID       `json:"id"`
	ProductName       string          `json:"product_name"`
	Price             float64         `json:"price"`
	Images            json.RawMessage `json:"images"`
	Description       string          `json:"description"`
	Features          json.RawMessage `json:"features"`
	SourceMarkdownURL string          `json:"source_markdown_url,omitempty"`
	Name              string          `json:"name,omitempty"`
	Slug              string          `json:"slug,omitempty"`
	Category          string          `json:"category,omitempty"`
	Brand             string          `json:"brand,omitempty"`
	Active            bool            `json:"active"`
	CreatedAt         int64           `json:"created_at"`
	UpdatedAt         int64           `json:"updated_at"`
}

// SetStoreFields sets the extended store fields on the product.
func (p *Product) SetStoreFields(name, category, brand string, active bool) {
	p.Name = name
	p.Category = category
	p.Brand = brand
	p.Active = active

	lowerName := strings.ToLower(strings.TrimSpace(name))
	if lowerName == "" {
		lowerName = strings.ToLower(strings.TrimSpace(p.ProductName))
	}
	p.Slug = slugRegex.ReplaceAllString(lowerName, "-")
	p.Slug = strings.Trim(p.Slug, "-")
}

func (p Product) HasID() bool {
	return p.ID != uuid.Nil
}

type Products []Product

func (p Products) IsEmpty() bool {
	return len(p) == 0
}
