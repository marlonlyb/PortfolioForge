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
	loginUser     model.User
	loginErr      error
	loginEmail    string
	loginPassword string
}

func (s *loginServiceUserStub) AdminLogin(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *loginServiceUserStub) PublicSignup(string, string) (model.EmailVerificationDispatchResult, error) {
	return model.EmailVerificationDispatchResult{}, errors.New("not implemented")
}

func (s *loginServiceUserStub) PublicLogin(email, password string) (model.User, error) {
	s.loginEmail = email
	s.loginPassword = password
	return s.loginUser, s.loginErr
}

func (s *loginServiceUserStub) LoginWithGoogle(model.GoogleIdentity) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func TestLoginPublicLoginIssuesAuthenticatedSession(t *testing.T) {
	serviceUser := &loginServiceUserStub{loginUser: model.User{
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

	userData, tokenSigned, err := service.PublicLogin("ada@example.com", "secret-123", "secret")
	if err != nil {
		t.Fatalf("PublicLogin() error = %v", err)
	}
	if serviceUser.loginEmail != "ada@example.com" || serviceUser.loginPassword != "secret-123" {
		t.Fatalf("PublicLogin() called with (%q, %q), want (%q, %q)", serviceUser.loginEmail, serviceUser.loginPassword, "ada@example.com", "secret-123")
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
