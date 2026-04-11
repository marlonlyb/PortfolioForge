package model

import "github.com/google/uuid"

type ProjectMedia struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	MediaType string    `json:"media_type"` // image|video|diagram|document
	URL       string    `json:"url"`
	Caption   string    `json:"caption,omitempty"`
	SortOrder int       `json:"sort_order"`
}
