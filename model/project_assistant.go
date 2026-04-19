package model

import "github.com/google/uuid"

type ProjectAssistantContext struct {
	ID                uuid.UUID
	Name              string
	Slug              string
	Active            bool
	SourceMarkdownURL string
	IndustryType      string
	FinalProduct      string
}

type ProjectAssistantMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ProjectAssistantRequest struct {
	Question string                    `json:"question"`
	History  []ProjectAssistantMessage `json:"history"`
	Lang     string                    `json:"lang"`
}

type ProjectAssistantResponse struct {
	Answer string `json:"answer"`
}
