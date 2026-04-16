package mailer

import "github.com/marlonlyb/portfolioforge/model"

type VerificationMailer interface {
	SendEmailVerificationOTP(message model.EmailVerificationMessage) error
}
