package model

import "github.com/google/uuid"

type EmailVerificationChallenge struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	CodeHash          string    `json:"-"`
	AttemptCount      int       `json:"attempt_count"`
	MaxAttempts       int       `json:"max_attempts"`
	ResendAvailableAt int64     `json:"resend_available_at"`
	ExpiresAt         int64     `json:"expires_at"`
	ConsumedAt        int64     `json:"consumed_at,omitempty"`
	CreatedAt         int64     `json:"created_at"`
	UpdatedAt         int64     `json:"updated_at"`
}

type EmailVerificationRequest struct {
	Email string `json:"email"`
}

type EmailVerificationVerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type EmailVerificationDispatchResult struct {
	VerificationRequired bool   `json:"verification_required"`
	Message              string `json:"message"`
	CooldownSeconds      int    `json:"cooldown_seconds"`
}

type EmailVerificationMessage struct {
	ToEmail         string
	OTPCode         string
	ExpiresInMinute int
	SupportLabel    string
}
