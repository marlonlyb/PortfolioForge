package login

import "github.com/marlonlyb/portfolioforge/model"

type Service interface {
	Login(email, password, jwtSecretKey string) (model.User, string, error)
}
