package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Project struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	Slug               string          `json:"slug"`
	Description        string          `json:"description"`
	Category           string          `json:"category"`
	ClientName         string          `json:"client_name,omitempty"`
	IndustryType       string          `json:"industry_type,omitempty"`
	FinalProduct       string          `json:"final_product,omitempty"`
	Status             string          `json:"status"` // draft|published|archived
	Featured           bool            `json:"featured"`
	Active             bool            `json:"active"`
	AssistantAvailable bool            `json:"assistant_available"`
	Images             json.RawMessage `json:"images"`
	CreatedAt          int64           `json:"created_at"`
	UpdatedAt          int64           `json:"updated_at"`
	// Joined data (not in projects table)
	Profile      *ProjectProfile `json:"profile,omitempty"`
	Technologies []Technology    `json:"technologies,omitempty"`
	Media        []ProjectMedia  `json:"media,omitempty"`
}

type Projects []Project

func (p Projects) IsEmpty() bool {
	return len(p) == 0
}
