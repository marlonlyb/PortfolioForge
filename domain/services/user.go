package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/marlonlyb/portfolioforge/domain/ports/mailer"
	"github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/model"
)

const (
	emailVerificationCodeLength   = 6
	emailVerificationTTL          = 10 * time.Minute
	emailVerificationCooldown     = 60 * time.Second
	emailVerificationMaxAttempts  = 5
	passwordlessLoginSupportLabel = "PortfolioForge"
)

type User struct {
	Repository user.Repository
	Mailer     mailer.VerificationMailer
}

func NewUser(ur user.Repository, verificationMailer mailer.VerificationMailer) *User {
	return &User{Repository: ur, Mailer: verificationMailer}
}

func (u *User) Create(m *model.User) error {
	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}
	m.ID = ID

	if m.Email == "" {
		return fmt.Errorf("%s", "email is empty!")
	}

	if m.Password == "" {
		return fmt.Errorf("%s", "password is empty!")
	}

	m.Email = normalizeEmail(m.Email)
	m.AuthProvider = "local"
	m.LocalAuthState = "ready"
	m.EmailVerified = false

	password, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s %w", "bcrypt.GenerateFromPassword()", err)
	}

	m.Password = string(password)

	if m.Details == nil {
		m.Details = []byte("{}")
	}

	now := time.Now().Unix()
	m.CreatedAt = now
	m.UpdatedAt = now

	err = u.Repository.Create(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Create()", err)
	}

	if !m.IsAdmin {
		if _, err := u.issueEmailVerificationChallenge(*m, now); err != nil {
			return fmt.Errorf("%s %w", "issueEmailVerificationChallenge()", err)
		}
	}

	m.Password = ""
	return nil
}

func (u *User) GetByID(ID uuid.UUID) (model.User, error) {
	user, err := u.Repository.GetByID(ID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.GetByID()", err)
	}

	return sanitizeUser(user), nil
}

func (u *User) GetByEmail(email string) (model.User, error) {
	user, err := u.Repository.GetByEmail(normalizeEmail(email))
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.GetByEmail()", err)
	}
	return sanitizeUser(user), nil
}

func (u *User) GetAll() (model.Users, error) {
	users, err := u.Repository.GetAll()
	if err != nil {
		return model.Users{}, fmt.Errorf("%s %w", "Repository.GetAll()", err)
	}
	for index := range users {
		users[index] = sanitizeUser(users[index])
	}
	return users, nil
}

func (u *User) AdminList() ([]model.AdminUserSummary, error) {
	users, err := u.Repository.AdminList()
	if err != nil {
		return nil, fmt.Errorf("%s %w", "Repository.AdminList()", err)
	}

	items := make([]model.AdminUserSummary, 0, len(users))
	for _, userData := range users {
		items = append(items, toAdminUserSummary(sanitizeUser(userData)))
	}

	return items, nil
}

func (u *User) AdminGetByID(ID uuid.UUID) (model.AdminUserDetail, error) {
	userData, err := u.Repository.AdminGetByID(ID)
	if err != nil {
		return model.AdminUserDetail{}, fmt.Errorf("%s %w", "Repository.AdminGetByID()", err)
	}

	return toAdminUserDetail(sanitizeUser(userData)), nil
}

func (u *User) AdminUpdate(ID uuid.UUID, request model.AdminUserUpdateRequest) (model.AdminUserDetail, error) {
	if request.IsAdmin == nil {
		return model.AdminUserDetail{}, model.ErrAdminUserUpdateScope
	}

	targetUser, err := u.Repository.AdminGetByID(ID)
	if err != nil {
		return model.AdminUserDetail{}, fmt.Errorf("%s %w", "Repository.AdminGetByID()", err)
	}
	if targetUser.IsAdmin {
		return model.AdminUserDetail{}, model.ErrAdminUserProtected
	}

	updatedUser, err := u.Repository.AdminSetIsAdmin(ID, *request.IsAdmin, time.Now().Unix())
	if err != nil {
		return model.AdminUserDetail{}, fmt.Errorf("%s %w", "Repository.AdminSetIsAdmin()", err)
	}

	return toAdminUserDetail(sanitizeUser(updatedUser)), nil
}

func (u *User) AdminSoftDelete(actorID, targetID uuid.UUID) error {
	targetUser, err := u.Repository.AdminGetByID(targetID)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.AdminGetByID()", err)
	}
	if targetUser.ID == actorID {
		return model.ErrAdminSelfDelete
	}
	if targetUser.IsAdmin {
		return model.ErrAdminUserProtected
	}

	if err := u.Repository.AdminSoftDelete(targetID, time.Now().Unix()); err != nil {
		return fmt.Errorf("%s %w", "Repository.AdminSoftDelete()", err)
	}

	return nil
}

func (u *User) AdminLogin(email, password string) (model.User, error) {
	m, err := u.Repository.GetByEmail(normalizeEmail(email))
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "GetByEmail()", err)
	}
	if strings.TrimSpace(m.AuthProvider) == "" {
		m.AuthProvider = "local"
	}
	if m.AuthProvider != "local" {
		return model.User{}, model.ErrProviderConflict
	}
	if !m.IsAdmin {
		return model.User{}, model.ErrProviderConflict
	}

	//aquí comparo los passwords, pero no que sean iguales, comparo sus
	//comportamientos de cambio ya que estamos usando bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "CompareHashAndPassword()", err)
	}

	m.LastLoginAt = time.Now().Unix()
	return sanitizeUser(m), nil
}

func (u *User) PublicSignup(email, password string) (model.EmailVerificationDispatchResult, error) {
	normalizedEmail := normalizeEmail(email)
	trimmedPassword := strings.TrimSpace(password)
	if normalizedEmail == "" || trimmedPassword == "" {
		return model.EmailVerificationDispatchResult{}, model.ErrInvalidCredentials
	}

	if existingUser, err := u.Repository.GetByEmail(normalizedEmail); err == nil {
		authProvider := strings.TrimSpace(existingUser.AuthProvider)
		if authProvider == "" {
			authProvider = "local"
		}
		if authProvider != "local" {
			return model.EmailVerificationDispatchResult{}, model.ErrProviderConflict
		}
		return model.EmailVerificationDispatchResult{}, model.ErrEmailAlreadyInUse
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return model.EmailVerificationDispatchResult{}, fmt.Errorf("%s %w", "Repository.GetByEmail()", err)
	}

	userData := &model.User{
		Email:    normalizedEmail,
		Password: trimmedPassword,
		IsAdmin:  false,
		Details:  []byte("{}"),
	}
	if err := u.Create(userData); err != nil {
		return model.EmailVerificationDispatchResult{}, err
	}

	return model.EmailVerificationDispatchResult{
		VerificationRequired: true,
		Message:              "Account created. Check your email for the verification code.",
		CooldownSeconds:      int(emailVerificationCooldown / time.Second),
	}, nil
}

func (u *User) PublicLogin(email, password string) (model.User, error) {
	normalizedEmail := normalizeEmail(email)
	trimmedPassword := strings.TrimSpace(password)
	if normalizedEmail == "" || trimmedPassword == "" {
		return model.User{}, model.ErrInvalidCredentials
	}

	userData, err := u.Repository.GetByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrInvalidCredentials
		}
		return model.User{}, fmt.Errorf("%s %w", "Repository.GetByEmail()", err)
	}

	authProvider := strings.TrimSpace(userData.AuthProvider)
	if authProvider == "" {
		authProvider = "local"
	}
	if authProvider != "local" {
		return model.User{}, model.ErrProviderConflict
	}
	if userData.IsAdmin {
		return model.User{}, model.ErrInvalidCredentials
	}
	if strings.TrimSpace(userData.Password) == "" || userData.LocalAuthState == "password_setup_required" {
		return model.User{}, model.ErrPasswordSetupRequired
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(trimmedPassword)); err != nil {
		return model.User{}, model.ErrInvalidCredentials
	}

	now := time.Now().Unix()
	updatedUser, err := u.Repository.UpdateLastLogin(userData.ID, now, now)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.UpdateLastLogin()", err)
	}

	return sanitizeUser(updatedUser), nil
}

func (u *User) LoginWithGoogle(identity model.GoogleIdentity) (model.User, error) {
	if !identity.EmailVerified {
		return model.User{}, model.ErrGoogleUnverifiedEmail
	}

	userData, err := u.Repository.UpsertGoogleUser(model.GoogleIdentity{
		Subject:       strings.TrimSpace(identity.Subject),
		Email:         normalizeEmail(identity.Email),
		EmailVerified: identity.EmailVerified,
		FullName:      strings.TrimSpace(identity.FullName),
	})
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.UpsertGoogleUser()", err)
	}

	return sanitizeUser(userData), nil
}

func (u *User) UpdateProfile(ID uuid.UUID, fullName, company string) (model.User, error) {
	trimmedFullName := strings.TrimSpace(fullName)
	trimmedCompany := strings.TrimSpace(company)
	if trimmedFullName == "" || trimmedCompany == "" {
		return model.User{}, errors.New("profile fields are required")
	}

	userData, err := u.Repository.UpdateProfile(ID, trimmedFullName, trimmedCompany)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.UpdateProfile()", err)
	}

	return sanitizeUser(userData), nil
}

func (u *User) RequestEmailVerification(email string) (model.EmailVerificationDispatchResult, error) {
	return u.dispatchEmailVerification(email, false)
}

func (u *User) RequestEmailLogin(email string) (model.EmailVerificationDispatchResult, error) {
	result := neutralEmailVerificationDispatchResult()
	normalizedEmail := normalizeEmail(email)
	if normalizedEmail == "" {
		return result, nil
	}

	now := time.Now().Unix()
	passwordHash, err := hashPasswordlessPlaceholder(normalizedEmail, now)
	if err != nil {
		return result, err
	}

	userData, err := u.Repository.UpsertPasswordlessPublicUser(normalizedEmail, passwordHash, now)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, nil
		}
		return result, fmt.Errorf("%s %w", "Repository.UpsertPasswordlessPublicUser()", err)
	}

	authProvider := strings.TrimSpace(userData.AuthProvider)
	if authProvider == "" {
		authProvider = "local"
	}
	if userData.IsAdmin || authProvider != "local" {
		return result, nil
	}

	latestChallenge, err := u.Repository.GetLatestEmailVerificationChallengeByUserID(userData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("%s %w", "Repository.GetLatestEmailVerificationChallengeByUserID()", err)
	}
	if err == nil && latestChallenge.ResendAvailableAt > now {
		return result, nil
	}

	if _, err := u.issueEmailVerificationChallenge(userData, now); err != nil {
		return result, err
	}

	return result, nil
}

func (u *User) ResendEmailVerification(email string) (model.EmailVerificationDispatchResult, error) {
	return u.dispatchEmailVerification(email, true)
}

func (u *User) VerifyEmailVerification(email, code string) (model.User, error) {
	return u.verifyEmailChallenge(email, code)
}

func (u *User) VerifyEmailLogin(email, code string) (model.User, error) {
	return u.verifyEmailChallenge(email, code)
}

func (u *User) verifyEmailChallenge(email, code string) (model.User, error) {
	normalizedEmail := normalizeEmail(email)
	trimmedCode := strings.TrimSpace(code)
	if normalizedEmail == "" || len(trimmedCode) != emailVerificationCodeLength {
		return model.User{}, model.ErrOTPInvalid
	}

	challenge, err := u.Repository.GetLatestEmailVerificationChallengeByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrOTPInvalid
		}
		return model.User{}, fmt.Errorf("%s %w", "Repository.GetLatestEmailVerificationChallengeByEmail()", err)
	}

	now := time.Now().Unix()
	if challenge.ConsumedAt > 0 || challenge.ExpiresAt <= now || challenge.AttemptCount >= challenge.MaxAttempts {
		return model.User{}, model.ErrOTPExpired
	}

	if err := bcrypt.CompareHashAndPassword([]byte(challenge.CodeHash), []byte(trimmedCode)); err != nil {
		attemptCount := challenge.AttemptCount + 1
		if updateErr := u.Repository.UpdateEmailVerificationChallengeAttempt(challenge.ID, attemptCount, now); updateErr != nil {
			return model.User{}, fmt.Errorf("%s %w", "Repository.UpdateEmailVerificationChallengeAttempt()", updateErr)
		}
		if attemptCount >= challenge.MaxAttempts {
			return model.User{}, model.ErrOTPExpired
		}
		return model.User{}, model.ErrOTPInvalid
	}

	if err := u.Repository.MarkEmailVerificationChallengeConsumed(challenge.ID, now); err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.MarkEmailVerificationChallengeConsumed()", err)
	}

	userData, err := u.Repository.MarkEmailVerified(challenge.UserID, now)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.MarkEmailVerified()", err)
	}

	return sanitizeUser(userData), nil
}

func (u *User) ToStoreUser(userData model.User) model.StoreUser {
	storeUser := model.StoreUser{
		ID:                     userData.ID,
		Email:                  userData.Email,
		IsAdmin:                userData.IsAdmin,
		AuthProvider:           userData.AuthProvider,
		EmailVerified:          userData.EmailVerified,
		FullName:               userData.FullName,
		Company:                userData.Company,
		ProfileCompleted:       userData.ProfileCompleted,
		AssistantEligible:      userData.AssistantEligible,
		CanUseProjectAssistant: userData.CanUseProjectAssistant,
		CreatedAt:              time.Unix(userData.CreatedAt, 0).UTC(),
	}

	if userData.LastLoginAt > 0 {
		storeUser.LastLoginAt = time.Unix(userData.LastLoginAt, 0).UTC()
	}

	return storeUser
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (u *User) dispatchEmailVerification(email string, allowCooldownBypass bool) (model.EmailVerificationDispatchResult, error) {
	result := neutralEmailVerificationDispatchResult()
	normalizedEmail := normalizeEmail(email)
	if normalizedEmail == "" {
		return result, nil
	}

	userData, err := u.Repository.GetByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, nil
		}
		return result, fmt.Errorf("%s %w", "Repository.GetByEmail()", err)
	}

	authProvider := strings.TrimSpace(userData.AuthProvider)
	if authProvider == "" {
		authProvider = "local"
	}
	if userData.IsAdmin || authProvider != "local" || userData.EmailVerified {
		return result, nil
	}

	now := time.Now().Unix()
	latestChallenge, err := u.Repository.GetLatestEmailVerificationChallengeByUserID(userData.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("%s %w", "Repository.GetLatestEmailVerificationChallengeByUserID()", err)
	}
	if err == nil && latestChallenge.ResendAvailableAt > now && !allowCooldownBypass {
		return result, nil
	}
	if err == nil && latestChallenge.ResendAvailableAt > now && allowCooldownBypass {
		return result, model.ErrOTPRateLimited
	}

	if _, err := u.issueEmailVerificationChallenge(userData, now); err != nil {
		return result, err
	}

	return result, nil
}

func (u *User) issueEmailVerificationChallenge(userData model.User, now int64) (model.EmailVerificationChallenge, error) {
	code, err := generateOTPCode(emailVerificationCodeLength)
	if err != nil {
		return model.EmailVerificationChallenge{}, err
	}

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return model.EmailVerificationChallenge{}, err
	}

	challengeID, err := uuid.NewUUID()
	if err != nil {
		return model.EmailVerificationChallenge{}, err
	}

	challenge := model.EmailVerificationChallenge{
		ID:                challengeID,
		UserID:            userData.ID,
		CodeHash:          string(codeHash),
		AttemptCount:      0,
		MaxAttempts:       emailVerificationMaxAttempts,
		ResendAvailableAt: now + int64(emailVerificationCooldown.Seconds()),
		ExpiresAt:         now + int64(emailVerificationTTL.Seconds()),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := u.Repository.CreateEmailVerificationChallenge(&challenge); err != nil {
		return model.EmailVerificationChallenge{}, fmt.Errorf("%s %w", "Repository.CreateEmailVerificationChallenge()", err)
	}

	if u.Mailer != nil {
		message := model.EmailVerificationMessage{
			ToEmail:         userData.Email,
			OTPCode:         code,
			ExpiresInMinute: int(emailVerificationTTL / time.Minute),
			SupportLabel:    passwordlessLoginSupportLabel,
		}
		if err := u.Mailer.SendEmailVerificationOTP(message); err != nil {
			return model.EmailVerificationChallenge{}, err
		}
	}

	return challenge, nil
}

func hashPasswordlessPlaceholder(email string, now int64) (string, error) {
	return hashPasswordValue(fmt.Sprintf("passwordless:%s:%d", email, now))
}

func hashPasswordValue(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func neutralEmailVerificationDispatchResult() model.EmailVerificationDispatchResult {
	return model.EmailVerificationDispatchResult{
		VerificationRequired: true,
		Message:              "If the account is eligible, a verification code will be sent shortly.",
		CooldownSeconds:      int(emailVerificationCooldown / time.Second),
	}
}

func generateOTPCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("otp length must be positive")
	}

	max := intPow(10, length)
	buffer := make([]byte, 4)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	value := int(buffer[0])<<24 | int(buffer[1])<<16 | int(buffer[2])<<8 | int(buffer[3])
	if value < 0 {
		value = -value
	}

	return fmt.Sprintf("%0*d", length, value%max), nil
}

func intPow(base, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
}

func sanitizeUser(userData model.User) model.User {
	userData.Password = ""
	if strings.TrimSpace(userData.AuthProvider) == "" {
		userData.AuthProvider = "local"
	}
	userData.FullName = strings.TrimSpace(userData.FullName)
	userData.Company = strings.TrimSpace(userData.Company)
	userData.ProfileCompleted = userData.FullName != "" && userData.Company != ""
	userData.AssistantEligible = userData.IsAdmin || ((userData.AuthProvider == "google" || userData.AuthProvider == "local") && userData.EmailVerified && userData.ProfileCompleted)
	userData.CanUseProjectAssistant = userData.AssistantEligible
	return userData
}

func toAdminUserSummary(userData model.User) model.AdminUserSummary {
	return model.AdminUserSummary{
		ID:            userData.ID,
		Email:         userData.Email,
		IsAdmin:       userData.IsAdmin,
		AuthProvider:  userData.AuthProvider,
		EmailVerified: userData.EmailVerified,
		FullName:      userData.FullName,
		Company:       userData.Company,
		CreatedAt:     time.Unix(userData.CreatedAt, 0).UTC(),
		UpdatedAt:     unixPtrToTime(userData.UpdatedAt),
		LastLoginAt:   unixPtrToTime(userData.LastLoginAt),
		DeletedAt:     unixPtrToTime(userData.DeletedAt),
	}
}

func toAdminUserDetail(userData model.User) model.AdminUserDetail {
	return model.AdminUserDetail{
		ID:            userData.ID,
		Email:         userData.Email,
		IsAdmin:       userData.IsAdmin,
		AuthProvider:  userData.AuthProvider,
		EmailVerified: userData.EmailVerified,
		FullName:      userData.FullName,
		Company:       userData.Company,
		CreatedAt:     time.Unix(userData.CreatedAt, 0).UTC(),
		UpdatedAt:     unixPtrToTime(userData.UpdatedAt),
		LastLoginAt:   unixPtrToTime(userData.LastLoginAt),
		DeletedAt:     unixPtrToTime(userData.DeletedAt),
	}
}

func unixPtrToTime(value int64) *time.Time {
	if value <= 0 {
		return nil
	}
	timestamp := time.Unix(value, 0).UTC()
	return &timestamp
}
