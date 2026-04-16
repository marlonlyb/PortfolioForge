package services

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type loginServiceUserStub struct {
	verifyUser  model.User
	verifyErr   error
	verifyEmail string
	verifyCode  string
}

func (s *loginServiceUserStub) AdminLogin(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *loginServiceUserStub) LoginWithGoogle(model.GoogleIdentity) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *loginServiceUserStub) RequestEmailLogin(string) (model.EmailVerificationDispatchResult, error) {
	return model.EmailVerificationDispatchResult{}, errors.New("not implemented")
}

func (s *loginServiceUserStub) VerifyEmailLogin(email, code string) (model.User, error) {
	s.verifyEmail = email
	s.verifyCode = code
	return s.verifyUser, s.verifyErr
}

func TestLoginVerifyEmailLoginIssuesAuthenticatedSession(t *testing.T) {
	serviceUser := &loginServiceUserStub{verifyUser: model.User{
		ID:                     uuid.New(),
		Email:                  "ada@example.com",
		AuthProvider:           "local",
		EmailVerified:          true,
		FullName:               "Ada Lovelace",
		Company:                "Analytical Engines",
		ProfileCompleted:       true,
		AssistantEligible:      true,
		CanUseProjectAssistant: true,
		CreatedAt:              time.Now().Unix(),
	}}

	service := NewLogin(serviceUser, nil)

	userData, tokenSigned, err := service.VerifyEmailLogin("ada@example.com", "123456", "secret")
	if err != nil {
		t.Fatalf("VerifyEmailLogin() error = %v", err)
	}
	if serviceUser.verifyEmail != "ada@example.com" || serviceUser.verifyCode != "123456" {
		t.Fatalf("VerifyEmailLogin() called with (%q, %q), want (%q, %q)", serviceUser.verifyEmail, serviceUser.verifyCode, "ada@example.com", "123456")
	}
	if tokenSigned == "" {
		t.Fatalf("expected signed token")
	}
	if userData.Password != "" {
		t.Fatalf("password = %q, want empty", userData.Password)
	}

	parsed, err := jwt.ParseWithClaims(tokenSigned, &model.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		t.Fatalf("ParseWithClaims() error = %v", err)
	}
	claims, ok := parsed.Claims.(*model.JWTCustomClaims)
	if !ok || !parsed.Valid {
		t.Fatalf("expected valid JWT custom claims")
	}
	if claims.Email != "ada@example.com" {
		t.Fatalf("claims email = %q, want ada@example.com", claims.Email)
	}
	if !claims.EmailVerified || !claims.CanUseProjectAssistant {
		t.Fatalf("unexpected auth claims: %#v", claims)
	}
}
