package handlers

import (
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/login"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	service   login.Service
	responser response.API
}

func NewLogin(us login.Service) *Login {
	return &Login{service: us}
}

func (h *Login) Login(c echo.Context) error {

	m := loginRequest{}

	err := c.Bind(&m)
	if err != nil {
		return response.ContractError(400, "validation_error", "Los datos enviados no son válidos")
	}

	userModel, tokenSigned, err := h.service.Login(m.Email, m.Password, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		if strings.Contains(err.Error(), "crypto/bcrypt: hashedPassword is not the hash of the given password") ||
			strings.Contains(err.Error(), "no rows in result set") {
			return response.ContractError(401, "invalid_credentials", "Email o contraseña inválidos")
		}
		return response.ContractError(500, "unexpected_error", "No fue posible iniciar sesión")
	}

	data := map[string]interface{}{
		"user": model.StoreUser{
			ID:        userModel.ID,
			Email:     userModel.Email,
			IsAdmin:   userModel.IsAdmin,
			CreatedAt: time.Unix(userModel.CreatedAt, 0).UTC(),
		},
		"token":      tokenSigned,
		"expires_in": int((12 * time.Hour).Seconds()),
	}
	return c.JSON(response.ContractOK(data))
}
