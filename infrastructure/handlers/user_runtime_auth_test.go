package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	userport "github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/middle"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type runtimeAuthUserServiceStub struct {
	users map[uuid.UUID]model.User
}

func (s *runtimeAuthUserServiceStub) Create(*model.User) error {
	return errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) GetByID(id uuid.UUID) (model.User, error) {
	userData, ok := s.users[id]
	if !ok || userData.DeletedAt > 0 {
		return model.User{}, pgx.ErrNoRows
	}

	return userData, nil
}

func (s *runtimeAuthUserServiceStub) GetByEmail(string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) AdminLogin(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) LoginWithGoogle(model.GoogleIdentity) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) RequestEmailLogin(string) (model.EmailVerificationDispatchResult, error) {
	return model.EmailVerificationDispatchResult{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) VerifyEmailLogin(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) UpdateProfile(uuid.UUID, string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) RequestEmailVerification(string) (model.EmailVerificationDispatchResult, error) {
	return model.EmailVerificationDispatchResult{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) ResendEmailVerification(string) (model.EmailVerificationDispatchResult, error) {
	return model.EmailVerificationDispatchResult{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) VerifyEmailVerification(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) GetAll() (model.Users, error) {
	return nil, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) AdminList() ([]model.AdminUserSummary, error) {
	items := make([]model.AdminUserSummary, 0, len(s.users))
	for _, userData := range s.users {
		if userData.DeletedAt > 0 {
			continue
		}

		items = append(items, model.AdminUserSummary{
			ID:            userData.ID,
			Email:         userData.Email,
			IsAdmin:       userData.IsAdmin,
			AuthProvider:  userData.AuthProvider,
			EmailVerified: userData.EmailVerified,
			CreatedAt:     time.Unix(userData.CreatedAt, 0).UTC(),
			UpdatedAt:     nil,
			LastLoginAt:   nil,
			DeletedAt:     nil,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].Email < items[j].Email
		}

		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	return items, nil
}

func (s *runtimeAuthUserServiceStub) AdminGetByID(id uuid.UUID) (model.AdminUserDetail, error) {
	userData, err := s.GetByID(id)
	if err != nil {
		return model.AdminUserDetail{}, err
	}

	return model.AdminUserDetail{
		ID:            userData.ID,
		Email:         userData.Email,
		IsAdmin:       userData.IsAdmin,
		AuthProvider:  userData.AuthProvider,
		EmailVerified: userData.EmailVerified,
		CreatedAt:     time.Unix(userData.CreatedAt, 0).UTC(),
		UpdatedAt:     nil,
		LastLoginAt:   nil,
		DeletedAt:     nil,
	}, nil
}

func (s *runtimeAuthUserServiceStub) AdminUpdate(uuid.UUID, model.AdminUserUpdateRequest) (model.AdminUserDetail, error) {
	return model.AdminUserDetail{}, errors.New("not implemented")
}

func (s *runtimeAuthUserServiceStub) AdminSoftDelete(_, id uuid.UUID) error {
	userData, ok := s.users[id]
	if !ok || userData.DeletedAt > 0 {
		return pgx.ErrNoRows
	}

	userData.DeletedAt = time.Now().Unix()
	s.users[id] = userData
	return nil
}

func (s *runtimeAuthUserServiceStub) ToStoreUser(userData model.User) model.StoreUser {
	return toStoreUser(userData)
}

var _ userport.Service = (*runtimeAuthUserServiceStub)(nil)

func TestAdminRoutesRejectAuthenticatedNonAdminsAtRuntime(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	nonAdmin := model.User{
		ID:           uuid.New(),
		Email:        "ada@example.com",
		IsAdmin:      false,
		AuthProvider: "local",
		CreatedAt:    time.Now().Add(-1 * time.Hour).Unix(),
	}
	targetID := uuid.New()
	service := &runtimeAuthUserServiceStub{users: map[uuid.UUID]model.User{nonAdmin.ID: nonAdmin}}
	server := buildRuntimeAuthTestServer(service)
	token := signRuntimeAuthTestToken(t, nonAdmin)

	tests := []struct {
		name string
		path string
	}{
		{name: "directory", path: "/api/v1/admin/users"},
		{name: "detail", path: "/api/v1/admin/users/" + targetID.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assertAPIError(t, rec, http.StatusForbidden, "forbidden")
		})
	}
}

func TestSoftDeletedUsersCannotReuseOldTokensOnProtectedRoutes(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	deletedAt := time.Now().Unix()
	softDeletedUser := model.User{
		ID:           uuid.New(),
		Email:        "deleted@example.com",
		IsAdmin:      false,
		AuthProvider: "local",
		CreatedAt:    time.Now().Add(-2 * time.Hour).Unix(),
		DeletedAt:    deletedAt,
	}
	softDeletedAdmin := model.User{
		ID:           uuid.New(),
		Email:        "deleted-admin@example.com",
		IsAdmin:      true,
		AuthProvider: "local",
		CreatedAt:    time.Now().Add(-2 * time.Hour).Unix(),
		DeletedAt:    deletedAt,
	}

	tests := []struct {
		name string
		user model.User
		path string
	}{
		{name: "private me restoration", user: softDeletedUser, path: "/api/v1/private/me"},
		{name: "admin route reuse", user: softDeletedAdmin, path: "/api/v1/admin/users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &runtimeAuthUserServiceStub{users: map[uuid.UUID]model.User{tt.user.ID: tt.user}}
			server := buildRuntimeAuthTestServer(service)
			token := signRuntimeAuthTestToken(t, tt.user)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			assertAPIError(t, rec, http.StatusUnauthorized, "authentication_required")
		})
	}
}

func TestAdminDeleteRemovesUserFromSubsequentDirectoryReads(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	adminUser := model.User{
		ID:           uuid.New(),
		Email:        "admin@example.com",
		IsAdmin:      true,
		AuthProvider: "local",
		CreatedAt:    time.Now().Add(-2 * time.Hour).Unix(),
	}
	standardUser := model.User{
		ID:            uuid.New(),
		Email:         "ada@example.com",
		IsAdmin:       false,
		AuthProvider:  "local",
		EmailVerified: true,
		CreatedAt:     time.Now().Add(-1 * time.Hour).Unix(),
	}

	service := &runtimeAuthUserServiceStub{users: map[uuid.UUID]model.User{
		adminUser.ID:    adminUser,
		standardUser.ID: standardUser,
	}}
	server := buildRuntimeAuthTestServer(service)
	token := signRuntimeAuthTestToken(t, adminUser)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/"+standardUser.ID.String(), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRec := httptest.NewRecorder()
	server.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteRec.Code, http.StatusNoContent)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listRec := httptest.NewRecorder()
	server.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listRec.Code, http.StatusOK)
	}

	var listPayload struct {
		Data struct {
			Items []model.AdminUserSummary `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(listRec.Body).Decode(&listPayload); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	for _, userData := range listPayload.Data.Items {
		if userData.ID == standardUser.ID {
			t.Fatalf("deleted user %s still present in admin list", standardUser.ID)
		}
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/"+standardUser.ID.String(), nil)
	detailReq.Header.Set("Authorization", "Bearer "+token)
	detailRec := httptest.NewRecorder()
	server.ServeHTTP(detailRec, detailReq)

	assertAPIError(t, detailRec, http.StatusNotFound, "not_found")
}

func buildRuntimeAuthTestServer(service *runtimeAuthUserServiceStub) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = response.HTTPErrorHandler
	authMiddleware := middle.New(service)
	userHandler := NewUser(service)

	privateGroup := e.Group("/api/v1/private", authMiddleware.IsValid)
	privateGroup.GET("/me", userHandler.Me)

	adminGroup := e.Group("/api/v1/admin/users", authMiddleware.IsValid, authMiddleware.IsAdmin)
	adminGroup.GET("", userHandler.GetAll)
	adminGroup.GET("/:id", userHandler.AdminGetByID)
	adminGroup.DELETE("/:id", userHandler.AdminDelete)

	return e
}

func signRuntimeAuthTestToken(t *testing.T, userData model.User) string {
	t.Helper()

	claims := model.JWTCustomClaims{
		UserID:                 userData.ID,
		Email:                  userData.Email,
		IsAdmin:                userData.IsAdmin,
		AuthProvider:           userData.AuthProvider,
		EmailVerified:          userData.EmailVerified,
		ProfileCompleted:       userData.ProfileCompleted,
		AssistantEligible:      userData.AssistantEligible,
		CanUseProjectAssistant: userData.CanUseProjectAssistant,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	return tokenSigned
}

func assertAPIError(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()

	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
	}

	var payload model.APIErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Error.Code != wantCode {
		t.Fatalf("error code = %q, want %q", payload.Error.Code, wantCode)
	}
}
