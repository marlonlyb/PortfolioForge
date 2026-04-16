package login

import "github.com/marlonlyb/portfolioforge/model"

type Service interface {
	AdminLogin(email, password, jwtSecretKey string) (model.User, string, error)
	RequestEmailLogin(email string) (model.EmailVerificationDispatchResult, error)
	VerifyEmailLogin(email, code, jwtSecretKey string) (model.User, string, error)
	LoginWithGoogle(idToken, jwtSecretKey string) (model.User, string, error)
}
