package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// esto es la información que va encriptada en el token
type JWTCustomClaims struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	IsAdmin bool      `json:"is_admin"`
	jwt.StandardClaims
}
