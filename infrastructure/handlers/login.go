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

type publicSignupRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
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

func (h *Login) PublicLogin(c echo.Context) error {
	request, err := bindPasswordLoginRequest(c)
	if err != nil {
		return err
	}

	userModel, tokenSigned, serviceErr := h.service.PublicLogin(request.Email, request.Password, os.Getenv("JWT_SECRET_KEY"))
	if serviceErr != nil {
		switch {
		case errors.Is(serviceErr, model.ErrInvalidCredentials):
			return response.ContractError(401, "invalid_credentials", "Invalid email or password")
		case errors.Is(serviceErr, model.ErrProviderConflict):
			return response.ContractError(409, "account_provider_conflict", "That email already uses Google sign-in. Continue with Google.")
		case errors.Is(serviceErr, model.ErrPasswordSetupRequired):
			return response.ContractError(409, "password_setup_required", "This account still needs a password setup or reset before it can log in.")
		default:
			return response.ContractError(500, "unexpected_error", "Unable to sign in")
		}
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"user":       toStoreUser(userModel),
		"token":      tokenSigned,
		"expires_in": int((12 * time.Hour).Seconds()),
	}))
}

func (h *Login) PublicSignup(c echo.Context) error {
	request, err := bindPublicSignupRequest(c)
	if err != nil {
		return err
	}

	result, serviceErr := h.service.PublicSignup(request.Email, request.Password, os.Getenv("JWT_SECRET_KEY"))
	if serviceErr != nil {
		switch {
		case errors.Is(serviceErr, model.ErrProviderConflict):
			return response.ContractError(409, "account_provider_conflict", "That email already uses Google sign-in. Continue with Google.")
		case errors.Is(serviceErr, model.ErrEmailAlreadyInUse):
			return response.ContractError(409, "email_already_in_use", "That email already belongs to an existing local account. Log in instead.")
		default:
			return response.ContractError(500, "unexpected_error", "Unable to create the account")
		}
	}

	return c.JSON(response.ContractOK(result))
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

func bindPasswordLoginRequest(c echo.Context) (loginRequest, error) {
	var request loginRequest
	if err := c.Bind(&request); err != nil {
		return request, response.ContractError(400, "validation_error", "Invalid sign-in payload")
	}
	if strings.TrimSpace(request.Email) == "" {
		return request, response.ContractError(400, "validation_error", "Email is required", model.APIErrorDetail{Field: "email", Issue: "required"})
	}
	if len(strings.TrimSpace(request.Password)) < 8 {
		return request, response.ContractError(400, "validation_error", "Password must contain at least 8 characters", model.APIErrorDetail{Field: "password", Issue: "invalid"})
	}
	return request, nil
}

func bindPublicSignupRequest(c echo.Context) (publicSignupRequest, error) {
	var request publicSignupRequest
	if err := c.Bind(&request); err != nil {
		return request, response.ContractError(400, "validation_error", "Invalid sign-up payload")
	}
	if strings.TrimSpace(request.Email) == "" {
		return request, response.ContractError(400, "validation_error", "Email is required", model.APIErrorDetail{Field: "email", Issue: "required"})
	}
	if len(strings.TrimSpace(request.Password)) < 8 {
		return request, response.ContractError(400, "validation_error", "Password must contain at least 8 characters", model.APIErrorDetail{Field: "password", Issue: "invalid"})
	}
	if strings.TrimSpace(request.ConfirmPassword) == "" {
		return request, response.ContractError(400, "validation_error", "Confirm password is required", model.APIErrorDetail{Field: "confirm_password", Issue: "required"})
	}
	if request.Password != request.ConfirmPassword {
		return request, response.ContractError(400, "validation_error", "Passwords must match", model.APIErrorDetail{Field: "confirm_password", Issue: "mismatch"})
	}
	return request, nil
}
