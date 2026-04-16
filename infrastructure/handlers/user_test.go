package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	userport "github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

type stubUserService struct {
	getByIDResp     model.User
	getByIDErr      error
	updateErr       error
	createErr       error
	createdUser     model.User
	requestResp     model.EmailVerificationDispatchResult
	requestErr      error
	resendResp      model.EmailVerificationDispatchResult
	resendErr       error
	verifyResp      model.User
	verifyErr       error
	adminListResp   []model.AdminUserSummary
	adminListErr    error
	adminDetailResp model.AdminUserDetail
	adminDetailErr  error
	adminUpdateResp model.AdminUserDetail
	adminUpdateErr  error
	adminDeleteErr  error
}

func (s *stubUserService) Create(user *model.User) error {
	if s.createErr != nil {
		return s.createErr
	}

	s.createdUser = *user
	return nil
}

func (s *stubUserService) GetByID(id uuid.UUID) (model.User, error) {
	if id == uuid.Nil {
		return model.User{}, model.ErrInvalidID
	}
	return s.getByIDResp, s.getByIDErr
}

func (s *stubUserService) GetByEmail(string) (model.User, error) { return model.User{}, nil }

func (s *stubUserService) AdminLogin(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *stubUserService) PublicSignup(string, string) (model.EmailVerificationDispatchResult, error) {
	return model.EmailVerificationDispatchResult{}, errors.New("not implemented")
}

func (s *stubUserService) PublicLogin(string, string) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *stubUserService) LoginWithGoogle(model.GoogleIdentity) (model.User, error) {
	return model.User{}, errors.New("not implemented")
}

func (s *stubUserService) UpdateProfile(_ uuid.UUID, fullName, company string) (model.User, error) {
	if s.updateErr != nil {
		return model.User{}, s.updateErr
	}

	s.getByIDResp.FullName = strings.TrimSpace(fullName)
	s.getByIDResp.Company = strings.TrimSpace(company)
	s.getByIDResp.ProfileCompleted = s.getByIDResp.FullName != "" && s.getByIDResp.Company != ""
	s.getByIDResp.AssistantEligible = s.getByIDResp.IsAdmin || ((s.getByIDResp.AuthProvider == "google" || s.getByIDResp.AuthProvider == "local") && s.getByIDResp.EmailVerified && s.getByIDResp.ProfileCompleted)
	s.getByIDResp.CanUseProjectAssistant = s.getByIDResp.AssistantEligible
	return s.getByIDResp, nil
}

func (s *stubUserService) RequestEmailVerification(string) (model.EmailVerificationDispatchResult, error) {
	return s.requestResp, s.requestErr
}

func (s *stubUserService) ResendEmailVerification(string) (model.EmailVerificationDispatchResult, error) {
	return s.resendResp, s.resendErr
}

func (s *stubUserService) VerifyEmailVerification(string, string) (model.User, error) {
	return s.verifyResp, s.verifyErr
}

func (s *stubUserService) GetAll() (model.Users, error) { return nil, nil }

func (s *stubUserService) AdminList() ([]model.AdminUserSummary, error) {
	return s.adminListResp, s.adminListErr
}

func (s *stubUserService) AdminGetByID(uuid.UUID) (model.AdminUserDetail, error) {
	return s.adminDetailResp, s.adminDetailErr
}

func (s *stubUserService) AdminUpdate(_ uuid.UUID, _ model.AdminUserUpdateRequest) (model.AdminUserDetail, error) {
	return s.adminUpdateResp, s.adminUpdateErr
}

func (s *stubUserService) AdminSoftDelete(_, _ uuid.UUID) error {
	return s.adminDeleteErr
}

func (s *stubUserService) ToStoreUser(userData model.User) model.StoreUser {
	return toStoreUser(userData)
}

var _ userport.Service = (*stubUserService)(nil)

func TestUserMeReturnsEligibilityPayload(t *testing.T) {
	userID := uuid.New()
	handler := NewUser(&stubUserService{getByIDResp: model.User{
		ID:                     userID,
		Email:                  "ada@example.com",
		AuthProvider:           "google",
		EmailVerified:          true,
		FullName:               "Ada Lovelace",
		Company:                "Analytical Engines",
		ProfileCompleted:       true,
		AssistantEligible:      true,
		CanUseProjectAssistant: true,
		CreatedAt:              time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix(),
		LastLoginAt:            time.Date(2026, 4, 15, 1, 0, 0, 0, time.UTC).Unix(),
	}})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/private/me", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	if err := handler.Me(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data model.StoreUser `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Data.AuthProvider != "google" {
		t.Fatalf("auth_provider = %q, want google", payload.Data.AuthProvider)
	}
	if !payload.Data.EmailVerified || !payload.Data.ProfileCompleted || !payload.Data.AssistantEligible || !payload.Data.CanUseProjectAssistant {
		t.Fatalf("unexpected eligibility payload: %#v", payload.Data)
	}
	if payload.Data.FullName != "Ada Lovelace" || payload.Data.Company != "Analytical Engines" {
		t.Fatalf("unexpected profile payload: %#v", payload.Data)
	}
}

func TestUserMeRejectsUnauthenticatedAccess(t *testing.T) {
	handler := NewUser(&stubUserService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/private/me", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	if err := handler.Me(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var payload model.APIErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Error.Code != "authentication_required" {
		t.Fatalf("error code = %q, want authentication_required", payload.Error.Code)
	}
}

func TestUserMeReturnsAdminAssistantEligibility(t *testing.T) {
	userID := uuid.New()
	handler := NewUser(&stubUserService{getByIDResp: model.User{
		ID:                     userID,
		Email:                  "admin@example.com",
		IsAdmin:                true,
		AuthProvider:           "local",
		ProfileCompleted:       false,
		AssistantEligible:      true,
		CanUseProjectAssistant: true,
		CreatedAt:              time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix(),
		LastLoginAt:            time.Date(2026, 4, 15, 1, 0, 0, 0, time.UTC).Unix(),
	}})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/private/me", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	if err := handler.Me(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data model.StoreUser `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Data.AuthProvider != "local" {
		t.Fatalf("auth_provider = %q, want local", payload.Data.AuthProvider)
	}
	if !payload.Data.IsAdmin {
		t.Fatalf("is_admin = false, want true")
	}
	if !payload.Data.AssistantEligible || !payload.Data.CanUseProjectAssistant {
		t.Fatalf("unexpected admin eligibility payload: %#v", payload.Data)
	}
}

func TestUserMeRestoresLocalPublicUsersWithoutAssistantEligibility(t *testing.T) {
	userID := uuid.New()
	handler := NewUser(&stubUserService{getByIDResp: model.User{
		ID:                     userID,
		Email:                  "ada@example.com",
		IsAdmin:                false,
		AuthProvider:           "local",
		EmailVerified:          false,
		FullName:               "Ada Lovelace",
		Company:                "Analytical Engines",
		ProfileCompleted:       true,
		AssistantEligible:      false,
		CanUseProjectAssistant: false,
		CreatedAt:              time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix(),
		LastLoginAt:            time.Date(2026, 4, 15, 1, 0, 0, 0, time.UTC).Unix(),
	}})

	payload := requestCurrentUser(t, handler, userID)

	if payload.AuthProvider != "local" {
		t.Fatalf("auth_provider = %q, want local", payload.AuthProvider)
	}
	if payload.IsAdmin {
		t.Fatalf("is_admin = true, want false")
	}
	if payload.AssistantEligible {
		t.Fatalf("assistant_eligible = true, want false")
	}
	if payload.CanUseProjectAssistant {
		t.Fatalf("can_use_project_assistant = true, want false")
	}
}

func TestUserUpdateProfileRefreshesAssistantEligibility(t *testing.T) {
	userID := uuid.New()
	service := &stubUserService{getByIDResp: model.User{
		ID:                     userID,
		Email:                  "ada@example.com",
		AuthProvider:           "google",
		EmailVerified:          true,
		ProfileCompleted:       false,
		AssistantEligible:      false,
		CanUseProjectAssistant: false,
		CreatedAt:              time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix(),
	}}
	handler := NewUser(service)

	initialPayload := requestCurrentUser(t, handler, userID)
	if initialPayload.AssistantEligible {
		t.Fatalf("initial assistant_eligible = true, want false")
	}

	updateBody := `{"full_name":"Ada Lovelace","company":"Analytical Engines"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/private/me/profile", bytes.NewBufferString(updateBody))
	updateReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	updateRec := httptest.NewRecorder()
	updateEcho := echo.New()
	updateCtx := updateEcho.NewContext(updateReq, updateRec)
	updateCtx.Set("userID", userID)

	if err := handler.UpdateProfile(updateCtx); err != nil {
		response.HTTPErrorHandler(err, updateCtx)
	}

	if updateRec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updateRec.Code, http.StatusOK)
	}

	var updatePayload struct {
		Data struct {
			User model.StoreUser `json:"user"`
		} `json:"data"`
	}
	if err := json.NewDecoder(updateRec.Body).Decode(&updatePayload); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if !updatePayload.Data.User.AssistantEligible {
		t.Fatalf("updated assistant_eligible = false, want true")
	}

	refreshedPayload := requestCurrentUser(t, handler, userID)
	if !refreshedPayload.ProfileCompleted || !refreshedPayload.AssistantEligible || !refreshedPayload.CanUseProjectAssistant {
		t.Fatalf("unexpected refreshed eligibility payload: %#v", refreshedPayload)
	}
	if refreshedPayload.FullName != "Ada Lovelace" || refreshedPayload.Company != "Analytical Engines" {
		t.Fatalf("unexpected refreshed profile payload: %#v", refreshedPayload)
	}
}

func TestEmailVerificationVerifyReturnsVerifiedUser(t *testing.T) {
	service := &stubUserService{verifyResp: model.User{
		ID:                uuid.New(),
		Email:             "ada@example.com",
		AuthProvider:      "local",
		EmailVerified:     true,
		FullName:          "Ada Lovelace",
		Company:           "Analytical Engines",
		ProfileCompleted:  true,
		AssistantEligible: true,
		CreatedAt:         time.Now().Unix(),
	}}
	handler := NewEmailVerification(service)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/email-verification/verify", bytes.NewBufferString(`{"email":"ada@example.com","code":"123456"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	if err := handler.Verify(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestEmailVerificationResendReturnsNeutralResponseWhenRateLimited(t *testing.T) {
	service := &stubUserService{
		resendResp: model.EmailVerificationDispatchResult{
			VerificationRequired: true,
			Message:              "If the account is eligible, a verification code will be sent shortly.",
			CooldownSeconds:      60,
		},
		resendErr: model.ErrOTPRateLimited,
	}
	handler := NewEmailVerification(service)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/email-verification/resend", bytes.NewBufferString(`{"email":"ada@example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	if err := handler.Resend(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAdminUserListReturnsItems(t *testing.T) {
	handler := NewUser(&stubUserService{adminListResp: []model.AdminUserSummary{{
		ID:            uuid.New(),
		Email:         "ada@example.com",
		AuthProvider:  "local",
		EmailVerified: true,
		CreatedAt:     time.Now().UTC(),
	}}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	if err := handler.GetAll(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAdminUserGetByIDReturnsDetail(t *testing.T) {
	targetID := uuid.New()
	handler := NewUser(&stubUserService{adminDetailResp: model.AdminUserDetail{
		ID:            targetID,
		Email:         "ada@example.com",
		AuthProvider:  "local",
		EmailVerified: true,
		CreatedAt:     time.Now().UTC(),
	}})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/"+targetID.String(), nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(targetID.String())

	if err := handler.AdminGetByID(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAdminUserUpdateRejectsUnknownFields(t *testing.T) {
	targetID := uuid.New()
	handler := NewUser(&stubUserService{})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/"+targetID.String(), bytes.NewBufferString(`{"is_admin":true,"email":"nope@example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(targetID.String())

	if err := handler.AdminUpdate(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAdminUserUpdateRejectsProtectedAdmins(t *testing.T) {
	targetID := uuid.New()
	handler := NewUser(&stubUserService{adminUpdateErr: model.ErrAdminUserProtected})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/"+targetID.String(), bytes.NewBufferString(`{"is_admin":false}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(targetID.String())

	if err := handler.AdminUpdate(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestAdminUserDeleteRejectsSelfDelete(t *testing.T) {
	actorID := uuid.New()
	targetID := uuid.New()
	handler := NewUser(&stubUserService{adminDeleteErr: model.ErrAdminSelfDelete})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/"+targetID.String(), nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("userID", actorID)
	c.SetParamNames("id")
	c.SetParamValues(targetID.String())

	if err := handler.AdminDelete(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestAdminUserDeleteReturnsNoContent(t *testing.T) {
	actorID := uuid.New()
	targetID := uuid.New()
	handler := NewUser(&stubUserService{})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/"+targetID.String(), nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("userID", actorID)
	c.SetParamNames("id")
	c.SetParamValues(targetID.String())

	if err := handler.AdminDelete(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func requestCurrentUser(t *testing.T, handler *User, userID uuid.UUID) model.StoreUser {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/private/me", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	if err := handler.Me(c); err != nil {
		response.HTTPErrorHandler(err, c)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data model.StoreUser `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	return payload.Data
}
