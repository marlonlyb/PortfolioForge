package model

import "github.com/google/uuid"

type ContactLead struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Company         string    `json:"company,omitempty"`
	ProjectInterest string    `json:"project_interest,omitempty"`
	Message         string    `json:"message"`
	Status          string    `json:"status"` // new|reviewed|contacted|closed
}
