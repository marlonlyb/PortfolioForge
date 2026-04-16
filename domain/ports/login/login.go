package login

import "github.com/marlonlyb/portfolioforge/model"

type Service interface {
	AdminLogin(email, password, jwtSecretKey string) (model.User, string, error)
	PublicLogin(email, password, jwtSecretKey string) (model.User, string, error)
	PublicSignup(email, password, jwtSecretKey string) (model.EmailVerificationDispatchResult, error)
	LoginWithGoogle(idToken, jwtSecretKey string) (model.User, string, error)
}
