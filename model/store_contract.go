package model

import (
	"time"

	"github.com/google/uuid"
)

type APIErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

type APIErrorPayload struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	Details   []APIErrorDetail `json:"details,omitempty"`
	RequestID string           `json:"request_id,omitempty"`
}

type APIErrorResponse struct {
	Error APIErrorPayload `json:"error"`
}

type StoreUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

// StoreProductVariant is the legacy storage-shaped variant contract.
// Prefer AdminProjectVariant/AdminProjectVariantInput in admin flows.
type StoreProductVariant struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	SKU       string    `json:"sku"`
	Color     string    `json:"color"`
	Size      string    `json:"size"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ImageURL  string    `json:"image_url,omitempty"`
}

// StoreProduct is the legacy storage-shaped read contract used by compat routes.
// Prefer Project/AdminProject contracts outside the compat layer.
type StoreProduct struct {
	ID                uuid.UUID             `json:"id"`
	Name              string                `json:"name"`
	Slug              string                `json:"slug"`
	Description       string                `json:"description"`
	Category          string                `json:"category"`
	Brand             string                `json:"brand,omitempty"`
	SourceMarkdownURL string                `json:"source_markdown_url,omitempty"`
	Images            []string              `json:"images"`
	Media             []ProjectMedia        `json:"media,omitempty"`
	Active            bool                  `json:"active"`
	PriceFrom         float64               `json:"price_from,omitempty"`
	AvailableColors   []string              `json:"available_colors,omitempty"`
	AvailableSizes    []string              `json:"available_sizes,omitempty"`
	Variants          []StoreProductVariant `json:"variants,omitempty"`
}
