package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID                     uuid.UUID       `json:"id"`
	Email                  string          `json:"email"`
	Password               string          `json:"password"`
	IsAdmin                bool            `json:"is_admin"`
	Details                json.RawMessage `json:"details"`
	AuthProvider           string          `json:"auth_provider"`
	ProviderSubject        string          `json:"provider_subject,omitempty"`
	EmailVerified          bool            `json:"email_verified"`
	FullName               string          `json:"full_name,omitempty"`
	Company                string          `json:"company,omitempty"`
	ProfileCompleted       bool            `json:"profile_completed"`
	AssistantEligible      bool            `json:"assistant_eligible"`
	CanUseProjectAssistant bool            `json:"can_use_project_assistant"`
	LocalAuthState         string          `json:"-"`
	LastLoginAt            int64           `json:"last_login_at,omitempty"`
	CreatedAt              int64           `json:"created_at"`
	UpdatedAt              int64           `json:"updated_at"`
	DeletedAt              int64           `json:"deleted_at,omitempty"`
}

type Users []User

type GoogleIdentity struct {
	Subject       string
	Email         string
	EmailVerified bool
	FullName      string
}
