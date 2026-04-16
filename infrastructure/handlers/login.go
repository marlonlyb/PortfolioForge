package handlers

import (
	"errors"
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

type googleLoginRequest struct {
	IDToken string `json:"id_token"`
}

type Login struct {
	service   login.Service
	responser response.API
}

func NewLogin(us login.Service) *Login {
	return &Login{service: us}
}

func (h *Login) AdminLogin(c echo.Context) error {

	m := loginRequest{}

	err := c.Bind(&m)
	if err != nil {
		return response.ContractError(400, "validation_error", "Invalid sign-in payload")
	}

	userModel, tokenSigned, err := h.service.AdminLogin(m.Email, m.Password, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		if strings.Contains(err.Error(), "crypto/bcrypt: hashedPassword is not the hash of the given password") ||
			strings.Contains(err.Error(), "no rows in result set") ||
			errors.Is(err, model.ErrProviderConflict) {
			return response.ContractError(401, "invalid_credentials", "Invalid email or password")
		}
		return response.ContractError(500, "unexpected_error", "Unable to sign in")
	}

	data := map[string]interface{}{
		"user":       toStoreUser(userModel),
		"token":      tokenSigned,
		"expires_in": int((12 * time.Hour).Seconds()),
	}
	return c.JSON(response.ContractOK(data))
}

func (h *Login) RequestEmailLogin(c echo.Context) error {
	var request model.EmailVerificationRequest
	if err := c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Invalid email login payload")
	}
	if strings.TrimSpace(request.Email) == "" {
		return response.ContractError(400, "validation_error", "Email is required", model.APIErrorDetail{Field: "email", Issue: "required"})
	}

	result, err := h.service.RequestEmailLogin(request.Email)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "Unable to process the email login request")
	}

	return c.JSON(response.ContractOK(result))
}

func (h *Login) VerifyEmailLogin(c echo.Context) error {
	var request model.EmailVerificationVerifyRequest
	if err := c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Invalid verification payload")
	}

	if strings.TrimSpace(request.Email) == "" {
		return response.ContractError(400, "validation_error", "Email is required", model.APIErrorDetail{Field: "email", Issue: "required"})
	}
	if len(strings.TrimSpace(request.Code)) != 6 {
		return response.ContractError(400, "validation_error", "Verification code must contain 6 digits", model.APIErrorDetail{Field: "code", Issue: "invalid"})
	}

	userModel, tokenSigned, err := h.service.VerifyEmailLogin(request.Email, request.Code, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOTPInvalid):
			return response.ContractError(400, "otp_invalid", "The verification code is invalid")
		case errors.Is(err, model.ErrOTPExpired):
			return response.ContractError(410, "otp_expired", "The verification code expired or must be renewed")
		default:
			return response.ContractError(500, "unexpected_error", "Unable to verify the email code")
		}
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"user":       toStoreUser(userModel),
		"token":      tokenSigned,
		"expires_in": int((12 * time.Hour).Seconds()),
	}))
}

func (h *Login) LoginWithGoogle(c echo.Context) error {
	var request googleLoginRequest
	if err := c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Invalid Google sign-in payload")
	}

	userModel, tokenSigned, err := h.service.LoginWithGoogle(request.IDToken, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrGoogleAuthUnavailable):
			return response.ContractError(503, "google_auth_unavailable", "Google sign-in is not configured")
		case errors.Is(err, model.ErrGoogleUnverifiedEmail):
			return response.ContractError(403, "google_email_unverified", "Google sign-in requires a verified email")
		case errors.Is(err, model.ErrProviderConflict):
			return response.ContractError(409, "account_provider_conflict", "That email already belongs to a local account. Use the existing sign-in method.")
		default:
			return response.ContractError(401, "invalid_google_token", "Unable to verify the Google sign-in")
		}
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"user":       toStoreUser(userModel),
		"token":      tokenSigned,
		"expires_in": int((12 * time.Hour).Seconds()),
	}))
}
