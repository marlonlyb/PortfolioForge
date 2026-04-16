package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/model"
)

const loginTokenTTL = 12 * time.Hour

type Login struct {
	ServiceUser    user.ServiceLogin
	GoogleVerifier GoogleTokenVerifier
}

type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (model.GoogleIdentity, error)
}

func NewLogin(usl user.ServiceLogin, verifier GoogleTokenVerifier) Login {
	return Login{ServiceUser: usl, GoogleVerifier: verifier}
}

func (l Login) AdminLogin(email, password, jwtSecretKey string) (model.User, string, error) {
	user, err := l.ServiceUser.AdminLogin(email, password)
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s %w", "ServiceUser.AdminLogin()", err)
	}
	return l.issueToken(user, jwtSecretKey)
}

func (l Login) RequestEmailLogin(email string) (model.EmailVerificationDispatchResult, error) {
	result, err := l.ServiceUser.RequestEmailLogin(email)
	if err != nil {
		return model.EmailVerificationDispatchResult{}, fmt.Errorf("%s %w", "ServiceUser.RequestEmailLogin()", err)
	}
	return result, nil
}

func (l Login) VerifyEmailLogin(email, code, jwtSecretKey string) (model.User, string, error) {
	user, err := l.ServiceUser.VerifyEmailLogin(email, code)
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s %w", "ServiceUser.VerifyEmailLogin()", err)
	}
	return l.issueToken(user, jwtSecretKey)
}

func (l Login) LoginWithGoogle(idToken, jwtSecretKey string) (model.User, string, error) {
	if l.GoogleVerifier == nil {
		return model.User{}, "", model.ErrGoogleAuthUnavailable
	}

	identity, err := l.GoogleVerifier.Verify(context.Background(), strings.TrimSpace(idToken))
	if err != nil {
		return model.User{}, "", err
	}

	user, err := l.ServiceUser.LoginWithGoogle(identity)
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s %w", "ServiceUser.LoginWithGoogle()", err)
	}

	return l.issueToken(user, jwtSecretKey)
}

func (l Login) issueToken(user model.User, jwtSecretKey string) (model.User, string, error) {
	claims := model.JWTCustomClaims{
		UserID:                 user.ID,
		Email:                  user.Email,
		IsAdmin:                user.IsAdmin,
		AuthProvider:           user.AuthProvider,
		EmailVerified:          user.EmailVerified,
		ProfileCompleted:       user.ProfileCompleted,
		AssistantEligible:      user.AssistantEligible,
		CanUseProjectAssistant: user.CanUseProjectAssistant,
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
