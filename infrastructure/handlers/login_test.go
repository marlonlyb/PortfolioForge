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
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	userservice "github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type mockLoginService struct {
	user   model.User
	token  string
	err    error
	result model.EmailVerificationDispatchResult
}

func (m *mockLoginService) AdminLogin(string, string, string) (model.User, string, error) {
	return m.user, m.token, m.err
}

func (m *mockLoginService) PublicLogin(string, string, string) (model.User, string, error) {
	return m.user, m.token, m.err
}

func (m *mockLoginService) PublicSignup(string, string, string) (model.EmailVerificationDispatchResult, error) {
	return m.result, m.err
}

func (m *mockLoginService) LoginWithGoogle(string, string) (model.User, string, error) {
	return m.user, m.token, m.err
}

type runtimeLoginRepositoryStub struct {
	users              map[string]model.User
	updatedLastLoginID uuid.UUID
	updatedLastLoginAt int64
	updatedAtTimestamp int64
}

func (s *runtimeLoginRepositoryStub) Create(*model.User) error { return errors.New("not implemented") }

func (s *runtimeLoginRepositoryStub) UpsertPasswordlessPublicUser(string, string, int64) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) GetByID(uuid.UUID) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) GetByEmail(email string) (model.User, error) {
	userData, ok := s.users[email]
	if !ok || userData.DeletedAt > 0 {
		return model.User{}, pgx.ErrNoRows
	}
	return userData, nil
}

func (s *runtimeLoginRepositoryStub) GetByProviderSubject(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) UpsertGoogleUser(model.GoogleIdentity) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) UpdateLastLogin(id uuid.UUID, lastLoginAt, updatedAt int64) (model.User, error) {
	for email, userData := range s.users {
		if userData.ID != id || userData.DeletedAt > 0 {
			continue
		}
		userData.LastLoginAt = lastLoginAt
		userData.UpdatedAt = updatedAt
		s.users[email] = userData
		s.updatedLastLoginID = id
		s.updatedLastLoginAt = lastLoginAt
		s.updatedAtTimestamp = updatedAt
		return userData, nil
	}
	return model.User{}, pgx.ErrNoRows
}

func (s *runtimeLoginRepositoryStub) UpdateProfile(uuid.UUID, string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) CreateEmailVerificationChallenge(*model.EmailVerificationChallenge) error {
	return errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) GetLatestEmailVerificationChallengeByUserID(uuid.UUID) (model.EmailVerificationChallenge, error) {
	return model.EmailVerificationChallenge{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) GetLatestEmailVerificationChallengeByEmail(string) (model.EmailVerificationChallenge, error) {
	return model.EmailVerificationChallenge{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) UpdateEmailVerificationChallengeAttempt(uuid.UUID, int, int64) error {
	return errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) MarkEmailVerificationChallengeConsumed(uuid.UUID, int64) error {
	return errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) MarkEmailVerified(uuid.UUID, int64) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) GetAll() (model.Users, error) {
	return nil, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) AdminList() (model.Users, error) {
	return nil, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) AdminGetByID(uuid.UUID) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) AdminSetIsAdmin(uuid.UUID, bool, int64) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeLoginRepositoryStub) AdminSoftDelete(uuid.UUID, int64) error {
	return errors.New("not implemented")
}

type noopVerificationMailer struct{}

func (noopVerificationMailer) SendEmailVerificationOTP(model.EmailVerificationMessage) error {
	return nil
}

func hashRuntimeLoginPassword(t *testing.T, raw string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	return string(hash)
}

func TestLogin_AdminLoginUsesUnifiedPasswordFlow(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/login", bytes.NewBufferString(`{"email":"admin@example.com","password":"secret-123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: &mockLoginService{user: model.User{Email: "admin@example.com", IsAdmin: true}, token: "admin.jwt"}, responser: response.API{}}
	if err := h.AdminLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestLogin_AdminLoginRuntimeAllowsNonAdminUserForCompatibilityAlias(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	userID := uuid.New()
	repo := &runtimeLoginRepositoryStub{users: map[string]model.User{
		"ada@example.com": {
			ID:             userID,
			Email:          "ada@example.com",
			Password:       hashRuntimeLoginPassword(t, "secret-123"),
			AuthProvider:   "local",
			LocalAuthState: "ready",
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := userservice.NewLogin(userservice.NewUser(repo, noopVerificationMailer{}), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/login", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: service, responser: response.API{}}
	if err := h.AdminLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if repo.updatedLastLoginID != userID {
		t.Fatalf("updated last login id = %s, want %s", repo.updatedLastLoginID, userID)
	}
}

func TestLogin_PublicLoginMapsProviderConflictAndSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		h := &Login{service: &mockLoginService{user: model.User{Email: "ada@example.com", AuthProvider: "local"}, token: "public.jwt"}, responser: response.API{}}
		if err := h.PublicLogin(c); err != nil {
			response.HTTPErrorHandler(err, c)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("provider conflict", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		h := &Login{service: &mockLoginService{err: model.ErrProviderConflict}, responser: response.API{}}
		if err := h.PublicLogin(c); err != nil {
			response.HTTPErrorHandler(err, c)
		}

		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
		}
	})
}

func TestLogin_PublicLoginRuntimeRejectsPasswordSetupRequiredUser(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	repo := &runtimeLoginRepositoryStub{users: map[string]model.User{
		"ada@example.com": {
			ID:             uuid.New(),
			Email:          "ada@example.com",
			Password:       hashRuntimeLoginPassword(t, "legacy-placeholder"),
			AuthProvider:   "local",
			LocalAuthState: "password_setup_required",
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := userservice.NewLogin(userservice.NewUser(repo, noopVerificationMailer{}), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: service, responser: response.API{}}
	if err := h.PublicLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	assertAPIError(t, rec, http.StatusConflict, "password_setup_required")
	if repo.updatedLastLoginID != uuid.Nil {
		t.Fatalf("unexpected UpdateLastLogin() call for migrated passwordless user")
	}
}

func TestLogin_PublicLoginRuntimeAllowsReadyLocalUser(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	userID := uuid.New()
	repo := &runtimeLoginRepositoryStub{users: map[string]model.User{
		"ada@example.com": {
			ID:                     userID,
			Email:                  "ada@example.com",
			Password:               hashRuntimeLoginPassword(t, "secret-123"),
			AuthProvider:           "local",
			LocalAuthState:         "ready",
			EmailVerified:          true,
			FullName:               "Ada Lovelace",
			Company:                "Analytical Engines",
			ProfileCompleted:       true,
			AssistantEligible:      true,
			CanUseProjectAssistant: true,
			CreatedAt:              time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:              time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := userservice.NewLogin(userservice.NewUser(repo, noopVerificationMailer{}), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: service, responser: response.API{}}
	if err := h.PublicLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if repo.updatedLastLoginID != userID {
		t.Fatalf("updated last login id = %s, want %s", repo.updatedLastLoginID, userID)
	}
}

func TestLogin_PublicLoginRuntimeRejectsInvalidCredentials(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	repo := &runtimeLoginRepositoryStub{users: map[string]model.User{
		"ada@example.com": {
			ID:             uuid.New(),
			Email:          "ada@example.com",
			Password:       hashRuntimeLoginPassword(t, "secret-123"),
			AuthProvider:   "local",
			LocalAuthState: "ready",
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := userservice.NewLogin(userservice.NewUser(repo, noopVerificationMailer{}), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/login", bytes.NewBufferString(`{"email":"ada@example.com","password":"wrong-password"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: service, responser: response.API{}}
	if err := h.PublicLogin(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	assertAPIError(t, rec, http.StatusUnauthorized, "invalid_credentials")
	if repo.updatedLastLoginID != uuid.Nil {
		t.Fatalf("unexpected UpdateLastLogin() call for invalid credentials")
	}
}

func TestLogin_PublicSignupReturnsVerificationContract(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/signup", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123","confirm_password":"secret-123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: &mockLoginService{result: model.EmailVerificationDispatchResult{VerificationRequired: true, Message: "Account created. Check your email for the verification code.", CooldownSeconds: 60}}, responser: response.API{}}
	if err := h.PublicSignup(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data model.EmailVerificationDispatchResult `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !payload.Data.VerificationRequired {
		t.Fatalf("verification_required = false, want true")
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
	c := echo.New().NewContext(req, rec)

	h := &Login{service: &mockLoginService{user: googleUser, token: "google.jwt.token"}, responser: response.API{}}
	if err := h.LoginWithGoogle(c); err != nil {
		t.Fatalf("LoginWithGoogle() returned unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestLogin_PublicSignupRejectsConflict(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/signup", bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123","confirm_password":"secret-123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	h := &Login{service: &mockLoginService{err: model.ErrEmailAlreadyInUse}, responser: response.API{}}
	if err := h.PublicSignup(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}
