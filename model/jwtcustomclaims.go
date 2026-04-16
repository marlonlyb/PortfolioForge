package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// esto es la información que va encriptada en el token
type JWTCustomClaims struct {
	UserID                 uuid.UUID `json:"user_id"`
	Email                  string    `json:"email"`
	IsAdmin                bool      `json:"is_admin"`
	AuthProvider           string    `json:"auth_provider"`
	EmailVerified          bool      `json:"email_verified"`
	ProfileCompleted       bool      `json:"profile_completed"`
	AssistantEligible      bool      `json:"assistant_eligible"`
	CanUseProjectAssistant bool      `json:"can_use_project_assistant"`
	jwt.StandardClaims
}
