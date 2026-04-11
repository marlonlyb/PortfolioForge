package model

import "github.com/google/uuid"

type Technology struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Category string    `json:"category"`
	Icon     string    `json:"icon,omitempty"`
	Color    string    `json:"color,omitempty"`
}
