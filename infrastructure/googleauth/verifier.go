package googleauth

import (
	"context"
	"errors"
	"os"
	"strings"

	"google.golang.org/api/idtoken"

	"github.com/marlonlyb/portfolioforge/model"
)

type Verifier struct {
	clientID string
}

func NewVerifier(clientID string) *Verifier {
	return &Verifier{clientID: strings.TrimSpace(clientID)}
}

func NewVerifierFromEnv() *Verifier {
	return NewVerifier(os.Getenv("GOOGLE_CLIENT_ID"))
}

func (v *Verifier) Verify(ctx context.Context, rawIDToken string) (model.GoogleIdentity, error) {
	if strings.TrimSpace(v.clientID) == "" {
		return model.GoogleIdentity{}, model.ErrGoogleAuthUnavailable
	}

	payload, err := idtoken.Validate(ctx, strings.TrimSpace(rawIDToken), v.clientID)
	if err != nil {
		return model.GoogleIdentity{}, err
	}

	email, _ := payload.Claims["email"].(string)
	fullName, _ := payload.Claims["name"].(string)

	verified, ok := extractEmailVerified(payload.Claims["email_verified"])
	if !ok || !verified {
		return model.GoogleIdentity{}, model.ErrGoogleUnverifiedEmail
	}

	if strings.TrimSpace(email) == "" || strings.TrimSpace(payload.Subject) == "" {
		return model.GoogleIdentity{}, errors.New("google token is missing required claims")
	}

	return model.GoogleIdentity{
		Subject:       payload.Subject,
		Email:         email,
		EmailVerified: verified,
		FullName:      fullName,
	}, nil
}

func extractEmailVerified(value any) (bool, bool) {
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		trimmed := strings.TrimSpace(strings.ToLower(typed))
		if trimmed == "true" {
			return true, true
		}
		if trimmed == "false" {
			return false, true
		}
	}

	return false, false
}
