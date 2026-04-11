package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/model"
)

const loginTokenTTL = 12 * time.Hour

type Login struct {
	/* services.Login usa un puerto creado en user
	del tipo Service específico para Login que solo
	tiene el Login(email, password string) (model.User, error) */
	ServiceUser user.ServiceLogin
}

func NewLogin(usl user.ServiceLogin) Login {
	return Login{ServiceUser: usl}
}

func (l Login) Login(email, password, jwtSecretKey string) (model.User, string, error) {
	user, err := l.ServiceUser.Login(email, password)
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s %w", "ServiceUser.Login()", err)
	}

	claims := model.JWTCustomClaims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(loginTokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSigned, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s %w", "token.SignedString()", err)
	}

	user.Password = ""

	return user, tokenSigned, nil
}
