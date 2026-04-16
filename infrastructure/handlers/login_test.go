package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type mockLoginService struct {
	user   model.User
	token  string
	err    error
	result model.EmailVerificationDispatchResult
}

func (m *mockLoginService) AdminLogin(email, password, jwtSecretKey string) (model.User, string, error) {
	return m.user, m.token, m.err
}

func (m *mockLoginService) RequestEmailLogin(email string) (model.EmailVerificationDispatchResult, error) {
	return m.result, m.err
}

func (m *mockLoginService) VerifyEmailLogin(email, code, jwtSecretKey string) (model.User, string, error) {
	return m.user, m.token, m.err
}

func (m *mockLoginService) LoginWithGoogle(idToken, jwtSecretKey string) (model.User, string, error) {
	return m.user, m.token, m.err
}

func TestLogin_AdminLogin(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	tests := []struct {
		name         string
		reqBody      string
		serviceErr   error
		userResp     model.User
		tokenResp    string
		expectedCode int
	}{
		{
			name:         "success",
			reqBody:      `{"email":"test@example.com","password":"password123"}`,
			userResp:     model.User{Email: "test@example.com"},
			tokenResp:    "valid.token.here",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid json",
			reqBody:      `{invalid}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			reqBody:      `{"email":"test@example.com","password":"wrongpassword"}`,
			serviceErr:   errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "other service error",
			reqBody:      `{"email":"test@example.com","password":"password123"}`,
			serviceErr:   errors.New("internal error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockSvc := &mockLoginService{
				user:  tt.userResp,
				token: tt.tokenResp,
				err:   tt.serviceErr,
			}

			h := &Login{
				service:   mockSvc,
				responser: response.API{},
			}

			err := h.AdminLogin(c)

			if err != nil {
				response.HTTPErrorHandler(err, c)
			}

			if rec.Code != tt.expectedCode {
				t.Errorf("Login() status code = %v, want %v", rec.Code, tt.expectedCode)
			}
		})
	}
}

func TestLogin_LoginWithGoogleReturnsGoogleSession(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	googleUser := model.User{
		ID:                     uuid.New(),
		Email:                  "ada@example.com",
		AuthProvider:           "google",
		EmailVerified:          true,
		FullName:               "Ada Lovelace",
		Company:                "Analytical Engines",
		ProfileCompleted:       true,
		AssistantEligible:      true,
		CanUseProjectAssistant: true,
		CreatedAt:              time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix(),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login/google", bytes.NewBufferString(`{"id_token":"google-id-token"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := &Login{
		service: &mockLoginService{
			user:  googleUser,
			token: "google.jwt.token",
		},
		responser: response.API{},
	}

	if err := h.LoginWithGoogle(c); err != nil {
		t.Fatalf("LoginWithGoogle() returned unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data struct {
			User      model.StoreUser `json:"user"`
			Token     string          `json:"token"`
			ExpiresIn int             `json:"expires_in"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Data.Token != "google.jwt.token" {
		t.Fatalf("token = %q, want google.jwt.token", payload.Data.Token)
	}
	if payload.Data.ExpiresIn <= 0 {
		t.Fatalf("expires_in = %d, want positive value", payload.Data.ExpiresIn)
	}
	if payload.Data.User.AuthProvider != "google" {
		t.Fatalf("auth_provider = %q, want google", payload.Data.User.AuthProvider)
	}
	if payload.Data.User.IsAdmin {
		t.Fatalf("is_admin = true, want false")
	}
	if !payload.Data.User.EmailVerified {
		t.Fatalf("email_verified = false, want true")
	}
}

func TestLogin_RequestEmailLoginReturnsNeutralDispatch(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login/email/request", bytes.NewBufferString(`{"email":"ada@example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := &Login{
		service: &mockLoginService{
			result: model.EmailVerificationDispatchResult{
				VerificationRequired: true,
				Message:              "If the account is eligible, a verification code will be sent shortly.",
				CooldownSeconds:      60,
			},
		},
		responser: response.API{},
	}

	if err := h.RequestEmailLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data struct {
			VerificationRequired bool   `json:"verification_required"`
			Message              string `json:"message"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !payload.Data.VerificationRequired {
		t.Fatalf("verification_required = false, want true")
	}
	if payload.Data.Message == "" {
		t.Fatalf("expected neutral message")
	}
}

func TestLogin_VerifyEmailLoginReturnsAuthenticatedSession(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	verifiedUser := model.User{
		ID:                     uuid.New(),
		Email:                  "ada@example.com",
		AuthProvider:           "local",
		EmailVerified:          true,
		FullName:               "Ada Lovelace",
		Company:                "Analytical Engines",
		ProfileCompleted:       true,
		AssistantEligible:      true,
		CanUseProjectAssistant: true,
		CreatedAt:              time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix(),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login/email/verify", bytes.NewBufferString(`{"email":"ada@example.com","code":"123456"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	h := &Login{
		service:   &mockLoginService{user: verifiedUser, token: "otp.jwt.token"},
		responser: response.API{},
	}

	if err := h.VerifyEmailLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data struct {
			User      model.StoreUser `json:"user"`
			Token     string          `json:"token"`
			ExpiresIn int             `json:"expires_in"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Data.Token != "otp.jwt.token" {
		t.Fatalf("token = %q, want otp.jwt.token", payload.Data.Token)
	}
	if payload.Data.ExpiresIn <= 0 {
		t.Fatalf("expires_in = %d, want positive value", payload.Data.ExpiresIn)
	}
	if payload.Data.User.Email != "ada@example.com" || !payload.Data.User.EmailVerified {
		t.Fatalf("unexpected user payload: %#v", payload.Data.User)
	}
}

func TestLogin_VerifyEmailLoginRejectsInvalidAndExpiredOTP(t *testing.T) {
	tests := []struct {
		name         string
		serviceErr   error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "invalid otp",
			serviceErr:   model.ErrOTPInvalid,
			expectedCode: http.StatusBadRequest,
			expectedBody: "otp_invalid",
		},
		{
			name:         "expired otp",
			serviceErr:   model.ErrOTPExpired,
			expectedCode: http.StatusGone,
			expectedBody: "otp_expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login/email/verify", bytes.NewBufferString(`{"email":"ada@example.com","code":"123456"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)

			h := &Login{
				service:   &mockLoginService{err: tt.serviceErr},
				responser: response.API{},
			}

			if err := h.VerifyEmailLogin(c); err != nil {
				response.HTTPErrorHandler(err, c)
			}

			if rec.Code != tt.expectedCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.expectedCode)
			}

			var payload model.APIErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if payload.Error.Code != tt.expectedBody {
				t.Fatalf("error code = %q, want %q", payload.Error.Code, tt.expectedBody)
			}
		})
	}
}
