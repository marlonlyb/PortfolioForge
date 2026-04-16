package model

import "errors"

var (
	ErrInvalidID                 = errors.New("the ID is not valid")
	ErrProviderConflict          = errors.New("user provider conflict")
	ErrGoogleUnverifiedEmail     = errors.New("google email is not verified")
	ErrGoogleAuthUnavailable     = errors.New("google authentication is unavailable")
	ErrAssistantIneligible       = errors.New("user is not eligible for the assistant")
	ErrEmailVerificationRequired = errors.New("email verification is required")
	ErrOTPInvalid                = errors.New("otp is invalid")
	ErrOTPExpired                = errors.New("otp is expired")
	ErrOTPRateLimited            = errors.New("otp is rate limited")
	ErrAdminUserUpdateScope      = errors.New("admin user update scope is invalid")
	ErrAdminUserProtected        = errors.New("admin user target is protected")
	ErrAdminSelfDelete           = errors.New("admin user self delete is forbidden")
)
