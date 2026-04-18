package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// AdminProjectWrite is the canonical admin write contract for projects.
type AdminProjectWrite struct {
	ID                uuid.UUID                  `json:"id,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Category          string                     `json:"category"`
	ClientName        string                     `json:"client_name,omitempty"`
	SourceMarkdownURL string                     `json:"source_markdown_url,omitempty"`
	Active            *bool                      `json:"active"`
	Images            json.RawMessage            `json:"images"`
	Media             []AdminProjectMediaInput   `json:"media"`
	Features          json.RawMessage            `json:"features"`
	Variants          []AdminProjectVariantInput `json:"variants"`

	// Legacy request aliases kept for backward compatibility.
	ProductName string `json:"product_name,omitempty"`
	Brand       string `json:"brand,omitempty"`

	CreatedAt int64 `json:"-"`
	UpdatedAt int64 `json:"-"`
}

type AdminProjectMediaInput struct {
	ID          string `json:"id"`
	MediaType   string `json:"media_type"`
	FallbackURL string `json:"fallback_url"`
	LowURL      string `json:"low_url"`
	MediumURL   string `json:"medium_url"`
	HighURL     string `json:"high_url"`
	Caption     string `json:"caption"`
	AltText     string `json:"alt_text"`
	SortOrder   int    `json:"sort_order"`
	Featured    bool   `json:"featured"`
}

type LegacyProjectMediaKeyError struct {
	Key         string
	Replacement string
}

func (e *LegacyProjectMediaKeyError) Error() string {
	return fmt.Sprintf("legacy project media key %q is not supported; use %q", e.Key, e.Replacement)
}

func (m *AdminProjectMediaInput) UnmarshalJSON(data []byte) error {
	type adminProjectMediaInputAlias AdminProjectMediaInput
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	legacyKeys := map[string]string{
		"url":           "fallback_url",
		"thumbnail_url": "low_url",
		"full_url":      "high_url",
	}
	for key, replacement := range legacyKeys {
		if _, exists := raw[key]; exists {
			return &LegacyProjectMediaKeyError{Key: key, Replacement: replacement}
		}
	}

	var alias adminProjectMediaInputAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	*m = AdminProjectMediaInput(alias)
	return nil
}

type AdminProjectVariantInput struct {
	ID       string  `json:"id"`
	SKU      string  `json:"sku"`
	Color    string  `json:"color"`
	Size     string  `json:"size"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	ImageURL string  `json:"image_url"`
}

// AdminProject is the canonical admin read contract for projects.
// Legacy aliases remain serialized during the compatibility window.
type AdminProject struct {
	ID                uuid.UUID             `json:"id"`
	Name              string                `json:"name"`
	Slug              string                `json:"slug"`
	Description       string                `json:"description"`
	Category          string                `json:"category"`
	ClientName        string                `json:"client_name,omitempty"`
	Brand             string                `json:"brand,omitempty"`
	SourceMarkdownURL string                `json:"source_markdown_url,omitempty"`
	Images            []string              `json:"images"`
	Media             []ProjectMedia        `json:"media,omitempty"`
	Profile           *ProjectProfile       `json:"profile,omitempty"`
	Technologies      []Technology          `json:"technologies,omitempty"`
	Active            bool                  `json:"active"`
	PriceFrom         float64               `json:"price_from,omitempty"`
	AvailableColors   []string              `json:"available_colors,omitempty"`
	AvailableSizes    []string              `json:"available_sizes,omitempty"`
	Variants          []AdminProjectVariant `json:"variants,omitempty"`
}

type AdminProjectVariant struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	ProductID uuid.UUID `json:"product_id,omitempty"`
	SKU       string    `json:"sku"`
	Color     string    `json:"color"`
	Size      string    `json:"size"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ImageURL  string    `json:"image_url,omitempty"`
}

func (p *AdminProjectWrite) Normalize() {
	if strings.TrimSpace(p.Name) == "" {
		p.Name = p.ProductName
	}
	if strings.TrimSpace(p.ClientName) == "" {
		p.ClientName = p.Brand
	}

	p.Name = strings.TrimSpace(p.Name)
	p.Category = strings.TrimSpace(p.Category)
	p.ClientName = strings.TrimSpace(p.ClientName)
	p.SourceMarkdownURL = strings.TrimSpace(p.SourceMarkdownURL)
}

func (p AdminProjectWrite) ResolveActive(defaultValue bool) bool {
	if p.Active == nil {
		return defaultValue
	}
	return *p.Active
}

func (p AdminProjectWrite) ToLegacyProduct(defaultActive bool) *Product {
	legacy := &Product{
		ID:                p.ID,
		ProductName:       p.Name,
		Images:            p.Images,
		Description:       p.Description,
		Features:          p.Features,
		SourceMarkdownURL: p.SourceMarkdownURL,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
	legacy.SetStoreFields(p.Name, p.Category, p.ClientName, p.ResolveActive(defaultActive))
	return legacy
}

func (v AdminProjectVariantInput) ToLegacy(projectID uuid.UUID) StoreProductVariant {
	variantID, _ := uuid.Parse(v.ID)
	return StoreProductVariant{
		ID:        variantID,
		ProductID: projectID,
		SKU:       v.SKU,
		Color:     v.Color,
		Size:      v.Size,
		Price:     v.Price,
		Stock:     v.Stock,
		ImageURL:  v.ImageURL,
	}
}

func AdminProjectFromStoreProduct(store StoreProduct) AdminProject {
	variants := make([]AdminProjectVariant, 0, len(store.Variants))
	for _, variant := range store.Variants {
		variants = append(variants, AdminProjectVariant{
			ID:        variant.ID,
			ProjectID: variant.ProductID,
			ProductID: variant.ProductID,
			SKU:       variant.SKU,
			Color:     variant.Color,
			Size:      variant.Size,
			Price:     variant.Price,
			Stock:     variant.Stock,
			ImageURL:  variant.ImageURL,
		})
	}

	return AdminProject{
		ID:                store.ID,
		Name:              store.Name,
		Slug:              store.Slug,
		Description:       store.Description,
		Category:          store.Category,
		ClientName:        store.Brand,
		Brand:             store.Brand,
		SourceMarkdownURL: store.SourceMarkdownURL,
		Images:            store.Images,
		Media:             store.Media,
		Active:            store.Active,
		PriceFrom:         store.PriceFrom,
		AvailableColors:   store.AvailableColors,
		AvailableSizes:    store.AvailableSizes,
		Variants:          variants,
	}
}

func AdminProjectsFromStoreProducts(items []StoreProduct) []AdminProject {
	projects := make([]AdminProject, 0, len(items))
	for _, item := range items {
		projects = append(projects, AdminProjectFromStoreProduct(item))
	}
	return projects
}

func (p *AdminProject) ApplyEnrichment(project Project) {
	p.Profile = project.Profile
	p.Technologies = project.Technologies
}

func (p AdminProject) ToProject() Project {
	images, _ := json.Marshal(p.Images)
	return Project{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Description:  p.Description,
		Category:     p.Category,
		ClientName:   p.ClientName,
		Active:       p.Active,
		Images:       images,
		Profile:      p.Profile,
		Technologies: p.Technologies,
		Media:        p.Media,
	}
}
