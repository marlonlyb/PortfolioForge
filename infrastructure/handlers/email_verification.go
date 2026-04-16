package handlers

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type EmailVerification struct {
	service user.Service
}

func NewEmailVerification(service user.Service) *EmailVerification {
	return &EmailVerification{service: service}
}

func (h *EmailVerification) Request(c echo.Context) error {
	request, err := bindEmailVerificationRequest(c)
	if err != nil {
		return err
	}

	result, serviceErr := h.service.RequestEmailVerification(request.Email)
	if serviceErr != nil && !errors.Is(serviceErr, model.ErrOTPRateLimited) {
		return response.ContractError(500, "unexpected_error", "Unable to process the verification request")
	}

	return c.JSON(response.ContractOK(result))
}

func (h *EmailVerification) Resend(c echo.Context) error {
	request, err := bindEmailVerificationRequest(c)
	if err != nil {
		return err
	}

	result, serviceErr := h.service.ResendEmailVerification(request.Email)
	if serviceErr != nil {
		if errors.Is(serviceErr, model.ErrOTPRateLimited) {
			return c.JSON(response.ContractOK(result))
		}
		return response.ContractError(500, "unexpected_error", "Unable to process the verification request")
	}

	return c.JSON(response.ContractOK(result))
}

func (h *EmailVerification) Verify(c echo.Context) error {
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

	userData, err := h.service.VerifyEmailVerification(request.Email, request.Code)
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
		"user": h.service.ToStoreUser(userData),
	}))
}

func bindEmailVerificationRequest(c echo.Context) (model.EmailVerificationRequest, error) {
	var request model.EmailVerificationRequest
	if err := c.Bind(&request); err != nil {
		return request, response.ContractError(400, "validation_error", "Invalid verification payload")
	}
	if strings.TrimSpace(request.Email) == "" {
		return request, response.ContractError(400, "validation_error", "Email is required", model.APIErrorDetail{Field: "email", Issue: "required"})
	}
	return request, nil
}
