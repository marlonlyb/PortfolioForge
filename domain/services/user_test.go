package services

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/marlonlyb/portfolioforge/model"
)

type userRepositoryStub struct {
	createdUser              *model.User
	userByEmail              map[string]model.User
	userByID                 map[uuid.UUID]string
	challengeByUserID        map[uuid.UUID]model.EmailVerificationChallenge
	challengeByEmail         map[string]model.EmailVerificationChallenge
	createdChallenge         *model.EmailVerificationChallenge
	updatedChallengeID       uuid.UUID
	updatedChallengeAttempts int
	consumedChallengeID      uuid.UUID
	markedVerifiedUserID     uuid.UUID
	updatedLastLoginUserID   uuid.UUID
	updatedLastLoginAt       int64
}

func (s *userRepositoryStub) Create(user *model.User) error {
	clone := *user
	s.createdUser = &clone
	if s.userByEmail == nil {
		s.userByEmail = map[string]model.User{}
	}
	s.userByEmail[user.Email] = clone
	s.indexUser(clone)
	return nil
}

func (s *userRepositoryStub) UpsertPasswordlessPublicUser(email, passwordHash string, now int64) (model.User, error) {
	if s.userByEmail == nil {
		s.userByEmail = map[string]model.User{}
	}
	if userData, ok := s.userByEmail[email]; ok {
		if userData.DeletedAt > 0 {
			return model.User{}, pgx.ErrNoRows
		}
		return userData, nil
	}

	userData := model.User{
		ID:           uuid.New(),
		Email:        email,
		Password:     passwordHash,
		AuthProvider: "local",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.userByEmail[email] = userData
	s.indexUser(userData)
	return userData, nil
}

func (s *userRepositoryStub) GetByID(id uuid.UUID) (model.User, error) {
	for _, user := range s.userByEmail {
		if user.ID == id && user.DeletedAt == 0 {
			return user, nil
		}
	}
	return model.User{}, pgx.ErrNoRows
}

func (s *userRepositoryStub) GetByEmail(email string) (model.User, error) {
	user, ok := s.userByEmail[email]
	if !ok || user.DeletedAt > 0 {
		return model.User{}, pgx.ErrNoRows
	}
	return user, nil
}

func (s *userRepositoryStub) GetByProviderSubject(string, string) (model.User, error) {
	return model.User{}, pgx.ErrNoRows
}

func (s *userRepositoryStub) UpsertGoogleUser(model.GoogleIdentity) (model.User, error) {
	return model.User{}, nil
}

func (s *userRepositoryStub) UpdateProfile(id uuid.UUID, fullName, company string) (model.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return model.User{}, err
	}
	user.FullName = fullName
	user.Company = company
	user.UpdatedAt = time.Now().Unix()
	s.userByEmail[user.Email] = user
	s.indexUser(user)
	return user, nil
}

func (s *userRepositoryStub) UpdateLastLogin(id uuid.UUID, lastLoginAt, updatedAt int64) (model.User, error) {
	userData, err := s.GetByID(id)
	if err != nil {
		return model.User{}, err
	}
	userData.LastLoginAt = lastLoginAt
	userData.UpdatedAt = updatedAt
	s.updatedLastLoginUserID = id
	s.updatedLastLoginAt = lastLoginAt
	s.userByEmail[userData.Email] = userData
	s.indexUser(userData)
	return userData, nil
}

func (s *userRepositoryStub) AdminList() (model.Users, error) {
	users := model.Users{}
	for _, user := range s.userByEmail {
		if user.DeletedAt == 0 {
			users = append(users, user)
		}
	}
	return users, nil
}

func (s *userRepositoryStub) AdminGetByID(id uuid.UUID) (model.User, error) {
	return s.GetByID(id)
}

func (s *userRepositoryStub) AdminSetIsAdmin(id uuid.UUID, isAdmin bool, updatedAt int64) (model.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return model.User{}, err
	}
	user.IsAdmin = isAdmin
	user.UpdatedAt = updatedAt
	s.userByEmail[user.Email] = user
	s.indexUser(user)
	return user, nil
}

func (s *userRepositoryStub) AdminSoftDelete(id uuid.UUID, deletedAt int64) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}
	user.DeletedAt = deletedAt
	user.UpdatedAt = deletedAt
	s.userByEmail[user.Email] = user
	s.indexUser(user)
	return nil
}

func (s *userRepositoryStub) CreateEmailVerificationChallenge(challenge *model.EmailVerificationChallenge) error {
	clone := *challenge
	s.createdChallenge = &clone
	if s.challengeByUserID == nil {
		s.challengeByUserID = map[uuid.UUID]model.EmailVerificationChallenge{}
	}
	if s.challengeByEmail == nil {
		s.challengeByEmail = map[string]model.EmailVerificationChallenge{}
	}
	s.challengeByUserID[challenge.UserID] = clone
	for email, user := range s.userByEmail {
		if user.ID == challenge.UserID {
			s.challengeByEmail[email] = clone
		}
	}
	return nil
}

func (s *userRepositoryStub) GetLatestEmailVerificationChallengeByUserID(userID uuid.UUID) (model.EmailVerificationChallenge, error) {
	challenge, ok := s.challengeByUserID[userID]
	if !ok {
		return model.EmailVerificationChallenge{}, pgx.ErrNoRows
	}
	return challenge, nil
}

func (s *userRepositoryStub) GetLatestEmailVerificationChallengeByEmail(email string) (model.EmailVerificationChallenge, error) {
	challenge, ok := s.challengeByEmail[email]
	if !ok {
		return model.EmailVerificationChallenge{}, pgx.ErrNoRows
	}
	return challenge, nil
}

func (s *userRepositoryStub) UpdateEmailVerificationChallengeAttempt(challengeID uuid.UUID, attemptCount int, updatedAt int64) error {
	s.updatedChallengeID = challengeID
	s.updatedChallengeAttempts = attemptCount
	challenge := s.challengeByEmail[findEmailByChallengeID(s.challengeByEmail, challengeID)]
	challenge.AttemptCount = attemptCount
	challenge.UpdatedAt = updatedAt
	for email, current := range s.challengeByEmail {
		if current.ID == challengeID {
			s.challengeByEmail[email] = challenge
		}
	}
	for userID, current := range s.challengeByUserID {
		if current.ID == challengeID {
			challenge.UserID = userID
			s.challengeByUserID[userID] = challenge
		}
	}
	return nil
}

func (s *userRepositoryStub) MarkEmailVerificationChallengeConsumed(challengeID uuid.UUID, consumedAt int64) error {
	s.consumedChallengeID = challengeID
	for email, challenge := range s.challengeByEmail {
		if challenge.ID == challengeID {
			challenge.ConsumedAt = consumedAt
			s.challengeByEmail[email] = challenge
		}
	}
	for userID, challenge := range s.challengeByUserID {
		if challenge.ID == challengeID {
			challenge.ConsumedAt = consumedAt
			s.challengeByUserID[userID] = challenge
		}
	}
	return nil
}

func (s *userRepositoryStub) MarkEmailVerified(userID uuid.UUID, verifiedAt int64) (model.User, error) {
	s.markedVerifiedUserID = userID
	for email, user := range s.userByEmail {
		if user.ID == userID {
			user.EmailVerified = true
			user.UpdatedAt = verifiedAt
			s.userByEmail[email] = user
			return user, nil
		}
	}
	return model.User{}, pgx.ErrNoRows
}

func (s *userRepositoryStub) GetAll() (model.Users, error) { return s.AdminList() }

type verificationMailerStub struct {
	lastMessage model.EmailVerificationMessage
	callCount   int
}

func (s *verificationMailerStub) SendEmailVerificationOTP(message model.EmailVerificationMessage) error {
	s.lastMessage = message
	s.callCount++
	return nil
}

func TestCreateLocalUserStartsPendingVerification(t *testing.T) {
	repo := &userRepositoryStub{userByEmail: map[string]model.User{}}
	mailer := &verificationMailerStub{}
	service := NewUser(repo, mailer)
	details, _ := json.Marshal(map[string]string{})

	user := &model.User{Email: "Ada@example.com", Password: "secret-123", Details: details}
	if err := service.Create(user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.EmailVerified {
		t.Fatalf("email_verified = true, want false")
	}
	if repo.createdUser == nil {
		t.Fatalf("expected user to be persisted")
	}
	if repo.createdUser.LocalAuthState != "ready" {
		t.Fatalf("local_auth_state = %q, want ready", repo.createdUser.LocalAuthState)
	}
	if repo.createdChallenge == nil {
		t.Fatalf("expected email verification challenge to be created")
	}
	if mailer.callCount != 1 {
		t.Fatalf("mailer call count = %d, want 1", mailer.callCount)
	}
	if len(mailer.lastMessage.OTPCode) != 6 {
		t.Fatalf("otp length = %d, want 6", len(mailer.lastMessage.OTPCode))
	}
}

func TestPublicLoginRejectsPasswordSetupRequiredLocalUser(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("legacy-placeholder"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:             uuid.New(),
			Email:          "ada@example.com",
			Password:       string(passwordHash),
			AuthProvider:   "local",
			LocalAuthState: "password_setup_required",
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := NewUser(repo, &verificationMailerStub{})

	_, err = service.PublicLogin("ada@example.com", "secret-123")
	if !errors.Is(err, model.ErrPasswordSetupRequired) {
		t.Fatalf("PublicLogin() error = %v, want %v", err, model.ErrPasswordSetupRequired)
	}
	if repo.updatedLastLoginUserID != uuid.Nil {
		t.Fatalf("unexpected UpdateLastLogin() call for password-setup-required user")
	}
}

func TestPublicLoginAllowsReadyLocalUser(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret-123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	userID := uuid.New()
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:             userID,
			Email:          "ada@example.com",
			Password:       string(passwordHash),
			AuthProvider:   "local",
			LocalAuthState: "ready",
			EmailVerified:  true,
			FullName:       "Ada Lovelace",
			Company:        "Analytical Engines",
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := NewUser(repo, &verificationMailerStub{})

	loggedInUser, err := service.PublicLogin("ada@example.com", "secret-123")
	if err != nil {
		t.Fatalf("PublicLogin() error = %v", err)
	}
	if loggedInUser.Email != "ada@example.com" {
		t.Fatalf("email = %q, want ada@example.com", loggedInUser.Email)
	}
	if loggedInUser.Password != "" {
		t.Fatalf("password = %q, want empty", loggedInUser.Password)
	}
	if repo.updatedLastLoginUserID != userID {
		t.Fatalf("updated last login user = %s, want %s", repo.updatedLastLoginUserID, userID)
	}
}

func TestPublicLoginAllowsAdminUserAndPersistsLastLogin(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret-123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	userID := uuid.New()
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"admin@example.com": {
			ID:             userID,
			Email:          "admin@example.com",
			Password:       string(passwordHash),
			AuthProvider:   "local",
			LocalAuthState: "ready",
			IsAdmin:        true,
			EmailVerified:  true,
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := NewUser(repo, &verificationMailerStub{})

	loggedInUser, err := service.PublicLogin("admin@example.com", "secret-123")
	if err != nil {
		t.Fatalf("PublicLogin() error = %v", err)
	}
	if !loggedInUser.IsAdmin {
		t.Fatalf("is_admin = false, want true")
	}
	if repo.updatedLastLoginUserID != userID {
		t.Fatalf("updated last login user = %s, want %s", repo.updatedLastLoginUserID, userID)
	}
}

func TestAdminLoginRejectsNonAdminLocalUser(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret-123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	userID := uuid.New()
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:             userID,
			Email:          "ada@example.com",
			Password:       string(passwordHash),
			AuthProvider:   "local",
			LocalAuthState: "ready",
			EmailVerified:  true,
			CreatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:      time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := NewUser(repo, &verificationMailerStub{})

	_, err = service.AdminLogin("ada@example.com", "secret-123")
	if !errorsIs(err, model.ErrForbidden) {
		t.Fatalf("AdminLogin() error = %v, want ErrForbidden", err)
	}
	if repo.updatedLastLoginUserID != uuid.Nil {
		t.Fatalf("unexpected UpdateLastLogin() call for non-admin admin login")
	}
}

func TestVerifyEmailVerificationMarksUserVerified(t *testing.T) {
	userID := uuid.New()
	codeHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &userRepositoryStub{
		userByEmail: map[string]model.User{
			"ada@example.com": {ID: userID, Email: "ada@example.com", AuthProvider: "local"},
		},
		challengeByEmail: map[string]model.EmailVerificationChallenge{
			"ada@example.com": {
				ID:          uuid.New(),
				UserID:      userID,
				CodeHash:    string(codeHash),
				MaxAttempts: 5,
				ExpiresAt:   time.Now().Add(10 * time.Minute).Unix(),
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			},
		},
		challengeByUserID: map[uuid.UUID]model.EmailVerificationChallenge{},
	}
	repo.challengeByUserID[userID] = repo.challengeByEmail["ada@example.com"]

	service := NewUser(repo, &verificationMailerStub{})
	verifiedUser, err := service.VerifyEmailVerification("ada@example.com", "123456")
	if err != nil {
		t.Fatalf("VerifyEmailVerification() error = %v", err)
	}

	if !verifiedUser.EmailVerified {
		t.Fatalf("email_verified = false, want true")
	}
	if repo.consumedChallengeID == uuid.Nil {
		t.Fatalf("expected challenge to be consumed")
	}
	if repo.markedVerifiedUserID != userID {
		t.Fatalf("marked verified user = %s, want %s", repo.markedVerifiedUserID, userID)
	}
}

func TestVerifyEmailVerificationRejectsWrongCodeAndTracksAttempts(t *testing.T) {
	userID := uuid.New()
	codeHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	challengeID := uuid.New()
	repo := &userRepositoryStub{
		userByEmail: map[string]model.User{
			"ada@example.com": {ID: userID, Email: "ada@example.com", AuthProvider: "local"},
		},
		challengeByEmail: map[string]model.EmailVerificationChallenge{
			"ada@example.com": {
				ID:          challengeID,
				UserID:      userID,
				CodeHash:    string(codeHash),
				MaxAttempts: 5,
				ExpiresAt:   time.Now().Add(10 * time.Minute).Unix(),
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			},
		},
		challengeByUserID: map[uuid.UUID]model.EmailVerificationChallenge{},
	}
	repo.challengeByUserID[userID] = repo.challengeByEmail["ada@example.com"]

	service := NewUser(repo, &verificationMailerStub{})
	_, err = service.VerifyEmailVerification("ada@example.com", "000000")
	if err == nil || !errorsIs(err, model.ErrOTPInvalid) {
		t.Fatalf("VerifyEmailVerification() error = %v, want ErrOTPInvalid", err)
	}
	if repo.updatedChallengeID != challengeID {
		t.Fatalf("updated challenge id = %s, want %s", repo.updatedChallengeID, challengeID)
	}
	if repo.updatedChallengeAttempts != 1 {
		t.Fatalf("attempt count = %d, want 1", repo.updatedChallengeAttempts)
	}
}

func TestVerifyEmailVerificationRejectsExpiredChallenge(t *testing.T) {
	userID := uuid.New()
	codeHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &userRepositoryStub{
		userByEmail: map[string]model.User{
			"ada@example.com": {ID: userID, Email: "ada@example.com", AuthProvider: "local"},
		},
		challengeByEmail: map[string]model.EmailVerificationChallenge{
			"ada@example.com": {
				ID:          uuid.New(),
				UserID:      userID,
				CodeHash:    string(codeHash),
				MaxAttempts: 5,
				ExpiresAt:   time.Now().Add(-1 * time.Minute).Unix(),
				CreatedAt:   time.Now().Add(-11 * time.Minute).Unix(),
				UpdatedAt:   time.Now().Add(-1 * time.Minute).Unix(),
			},
		},
		challengeByUserID: map[uuid.UUID]model.EmailVerificationChallenge{},
	}
	repo.challengeByUserID[userID] = repo.challengeByEmail["ada@example.com"]

	service := NewUser(repo, &verificationMailerStub{})
	_, err = service.VerifyEmailVerification("ada@example.com", "123456")
	if err == nil || !errorsIs(err, model.ErrOTPExpired) {
		t.Fatalf("VerifyEmailVerification() error = %v, want ErrOTPExpired", err)
	}
	if repo.updatedChallengeID != uuid.Nil {
		t.Fatalf("updated challenge id = %s, want nil", repo.updatedChallengeID)
	}
	if repo.markedVerifiedUserID != uuid.Nil {
		t.Fatalf("marked verified user = %s, want nil", repo.markedVerifiedUserID)
	}
}

func TestVerifyEmailVerificationWrongCodeExhaustsMaxAttempts(t *testing.T) {
	userID := uuid.New()
	codeHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	challengeID := uuid.New()
	repo := &userRepositoryStub{
		userByEmail: map[string]model.User{
			"ada@example.com": {ID: userID, Email: "ada@example.com", AuthProvider: "local"},
		},
		challengeByEmail: map[string]model.EmailVerificationChallenge{
			"ada@example.com": {
				ID:           challengeID,
				UserID:       userID,
				CodeHash:     string(codeHash),
				AttemptCount: 4,
				MaxAttempts:  5,
				ExpiresAt:    time.Now().Add(10 * time.Minute).Unix(),
				CreatedAt:    time.Now().Unix(),
				UpdatedAt:    time.Now().Unix(),
			},
		},
		challengeByUserID: map[uuid.UUID]model.EmailVerificationChallenge{},
	}
	repo.challengeByUserID[userID] = repo.challengeByEmail["ada@example.com"]

	service := NewUser(repo, &verificationMailerStub{})
	_, err = service.VerifyEmailVerification("ada@example.com", "000000")
	if err == nil || !errorsIs(err, model.ErrOTPExpired) {
		t.Fatalf("VerifyEmailVerification() error = %v, want ErrOTPExpired", err)
	}
	if repo.updatedChallengeID != challengeID {
		t.Fatalf("updated challenge id = %s, want %s", repo.updatedChallengeID, challengeID)
	}
	if repo.updatedChallengeAttempts != 5 {
		t.Fatalf("attempt count = %d, want 5", repo.updatedChallengeAttempts)
	}
	if repo.markedVerifiedUserID != uuid.Nil {
		t.Fatalf("marked verified user = %s, want nil", repo.markedVerifiedUserID)
	}
}

func TestEmailVerificationRequestAndResendStayNeutralForIneligibleAccounts(t *testing.T) {
	tests := []struct {
		name   string
		email  string
		user   *model.User
		invoke func(*User, string) (model.EmailVerificationDispatchResult, error)
	}{
		{
			name:  "request ignores nonexistent local account",
			email: "missing@example.com",
			invoke: func(service *User, email string) (model.EmailVerificationDispatchResult, error) {
				return service.RequestEmailVerification(email)
			},
		},
		{
			name:  "request ignores already verified local account",
			email: "verified@example.com",
			user:  &model.User{ID: uuid.New(), Email: "verified@example.com", AuthProvider: "local", EmailVerified: true},
			invoke: func(service *User, email string) (model.EmailVerificationDispatchResult, error) {
				return service.RequestEmailVerification(email)
			},
		},
		{
			name:  "resend ignores non local provider",
			email: "google@example.com",
			user:  &model.User{ID: uuid.New(), Email: "google@example.com", AuthProvider: "google", EmailVerified: true},
			invoke: func(service *User, email string) (model.EmailVerificationDispatchResult, error) {
				return service.ResendEmailVerification(email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepositoryStub{userByEmail: map[string]model.User{}}
			if tt.user != nil {
				repo.userByEmail[tt.email] = *tt.user
			}
			mailer := &verificationMailerStub{}
			service := NewUser(repo, mailer)

			result, err := tt.invoke(service, tt.email)
			if err != nil {
				t.Fatalf("dispatch error = %v, want nil", err)
			}
			assertNeutralDispatchResult(t, result)
			if repo.createdChallenge != nil {
				t.Fatalf("created challenge = %#v, want nil", repo.createdChallenge)
			}
			if mailer.callCount != 0 {
				t.Fatalf("mailer call count = %d, want 0", mailer.callCount)
			}
		})
	}
}

func TestRequestEmailLoginCreatesPasswordlessUserForFirstTimeEmail(t *testing.T) {
	repo := &userRepositoryStub{userByEmail: map[string]model.User{}}
	mailer := &verificationMailerStub{}
	service := NewUser(repo, mailer)

	result, err := service.RequestEmailLogin("Ada@example.com")
	if err != nil {
		t.Fatalf("RequestEmailLogin() error = %v", err)
	}

	assertNeutralDispatchResult(t, result)
	if repo.createdChallenge == nil {
		t.Fatalf("expected login OTP challenge to be created")
	}
	createdUser, ok := repo.userByEmail["ada@example.com"]
	if !ok {
		t.Fatalf("expected passwordless local user to be created")
	}
	if createdUser.Password == "" {
		t.Fatalf("expected placeholder password hash to be stored")
	}
	if mailer.callCount != 1 {
		t.Fatalf("mailer call count = %d, want 1", mailer.callCount)
	}
}

func TestRequestEmailLoginReusesExistingNonAdminLocalUser(t *testing.T) {
	userID := uuid.New()
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:           userID,
			Email:        "ada@example.com",
			Password:     "existing-placeholder-hash",
			AuthProvider: "local",
			CreatedAt:    time.Now().Add(-1 * time.Hour).Unix(),
			UpdatedAt:    time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	mailer := &verificationMailerStub{}
	service := NewUser(repo, mailer)

	result, err := service.RequestEmailLogin(" Ada@example.com ")
	if err != nil {
		t.Fatalf("RequestEmailLogin() error = %v", err)
	}

	assertNeutralDispatchResult(t, result)
	if repo.createdUser != nil {
		t.Fatalf("unexpected sign-up style Create() call for existing local user")
	}
	if repo.createdChallenge == nil {
		t.Fatalf("expected challenge for existing local user")
	}
	if repo.createdChallenge.UserID != userID {
		t.Fatalf("challenge user id = %s, want %s", repo.createdChallenge.UserID, userID)
	}
	if mailer.callCount != 1 {
		t.Fatalf("mailer call count = %d, want 1", mailer.callCount)
	}
	if repo.userByEmail["ada@example.com"].Password != "existing-placeholder-hash" {
		t.Fatalf("existing user password hash was unexpectedly replaced")
	}
	if _, ok := repo.userByEmail[" Ada@example.com "]; ok {
		t.Fatalf("unexpected non-normalized email key created")
	}
}

func TestRequestEmailLoginStaysNeutralForAdminAndProviderConflicts(t *testing.T) {
	tests := []struct {
		name string
		user model.User
	}{
		{
			name: "admin account",
			user: model.User{ID: uuid.New(), Email: "admin@example.com", IsAdmin: true, AuthProvider: "local"},
		},
		{
			name: "google account",
			user: model.User{ID: uuid.New(), Email: "google@example.com", AuthProvider: "google", EmailVerified: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepositoryStub{userByEmail: map[string]model.User{tt.user.Email: tt.user}}
			mailer := &verificationMailerStub{}
			service := NewUser(repo, mailer)

			result, err := service.RequestEmailLogin(tt.user.Email)
			if err != nil {
				t.Fatalf("RequestEmailLogin() error = %v", err)
			}

			assertNeutralDispatchResult(t, result)
			if repo.createdChallenge != nil {
				t.Fatalf("created challenge = %#v, want nil", repo.createdChallenge)
			}
			if mailer.callCount != 0 {
				t.Fatalf("mailer call count = %d, want 0", mailer.callCount)
			}
		})
	}
}

func TestAdminUpdatePromotesOnlyActiveNonAdminUsers(t *testing.T) {
	targetID := uuid.New()
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:           targetID,
			Email:        "ada@example.com",
			AuthProvider: "local",
			FullName:     "Ada Lovelace",
			Company:      "Analytical Engines",
			CreatedAt:    time.Now().Add(-1 * time.Hour).Unix(),
		},
	}}
	service := NewUser(repo, &verificationMailerStub{})
	makeAdmin := true

	updated, err := service.AdminUpdate(targetID, model.AdminUserUpdateRequest{IsAdmin: &makeAdmin})
	if err != nil {
		t.Fatalf("AdminUpdate() error = %v", err)
	}
	if !updated.IsAdmin {
		t.Fatalf("is_admin = false, want true")
	}
	if updated.FullName != "Ada Lovelace" || updated.Company != "Analytical Engines" {
		t.Fatalf("unexpected profile mutation: %#v", updated)
	}
}

func TestAdminUpdateRejectsExistingAdminUsers(t *testing.T) {
	targetID := uuid.New()
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"admin@example.com": {ID: targetID, Email: "admin@example.com", IsAdmin: true, AuthProvider: "local"},
	}}
	service := NewUser(repo, &verificationMailerStub{})
	makeAdmin := false

	_, err := service.AdminUpdate(targetID, model.AdminUserUpdateRequest{IsAdmin: &makeAdmin})
	if !errorsIs(err, model.ErrAdminUserProtected) {
		t.Fatalf("AdminUpdate() error = %v, want ErrAdminUserProtected", err)
	}
}

func TestAdminSoftDeleteRejectsSelfAndOtherAdmins(t *testing.T) {
	adminID := uuid.New()
	otherAdminID := uuid.New()
	service := NewUser(&userRepositoryStub{userByEmail: map[string]model.User{
		"self@example.com":  {ID: adminID, Email: "self@example.com", IsAdmin: true, AuthProvider: "local"},
		"other@example.com": {ID: otherAdminID, Email: "other@example.com", IsAdmin: true, AuthProvider: "local"},
	}}, &verificationMailerStub{})

	if err := service.AdminSoftDelete(adminID, adminID); !errorsIs(err, model.ErrAdminSelfDelete) {
		t.Fatalf("AdminSoftDelete(self) error = %v, want ErrAdminSelfDelete", err)
	}
	if err := service.AdminSoftDelete(adminID, otherAdminID); !errorsIs(err, model.ErrAdminUserProtected) {
		t.Fatalf("AdminSoftDelete(other admin) error = %v, want ErrAdminUserProtected", err)
	}
}

func TestDeletedUsersCannotRestoreAuthenticatedReads(t *testing.T) {
	deletedID := uuid.New()
	service := NewUser(&userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:        deletedID,
			Email:     "ada@example.com",
			DeletedAt: time.Now().Unix(),
		},
	}}, &verificationMailerStub{})

	_, err := service.GetByID(deletedID)
	if err == nil || !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		t.Fatalf("GetByID() error = %v, want wrapped no rows", err)
	}
}

func TestRequestEmailLoginStaysNeutralForDeletedIdentity(t *testing.T) {
	repo := &userRepositoryStub{userByEmail: map[string]model.User{
		"ada@example.com": {
			ID:           uuid.New(),
			Email:        "ada@example.com",
			AuthProvider: "local",
			DeletedAt:    time.Now().Unix(),
		},
	}}
	service := NewUser(repo, &verificationMailerStub{})

	result, err := service.RequestEmailLogin("ada@example.com")
	if err != nil {
		t.Fatalf("RequestEmailLogin() error = %v", err)
	}
	assertNeutralDispatchResult(t, result)
	if repo.createdChallenge != nil {
		t.Fatalf("created challenge = %#v, want nil", repo.createdChallenge)
	}
}

func TestSanitizeUserAllowsVerifiedLocalAssistantEligibility(t *testing.T) {
	userData := sanitizeUser(model.User{
		Email:            "ada@example.com",
		AuthProvider:     "local",
		EmailVerified:    true,
		FullName:         "Ada Lovelace",
		Company:          "Analytical Engines",
		ProfileCompleted: true,
	})

	if !userData.AssistantEligible || !userData.CanUseProjectAssistant {
		t.Fatalf("expected verified local user to be assistant eligible: %#v", userData)
	}
}

func findEmailByChallengeID(challenges map[string]model.EmailVerificationChallenge, challengeID uuid.UUID) string {
	for email, challenge := range challenges {
		if challenge.ID == challengeID {
			return email
		}
	}
	return ""
}

func (s *userRepositoryStub) indexUser(user model.User) {
	if s.userByID == nil {
		s.userByID = map[uuid.UUID]string{}
	}
	s.userByID[user.ID] = user.Email
}

func errorsIs(err error, target error) bool {
	return err != nil && target != nil && err.Error() == target.Error()
}

func assertNeutralDispatchResult(t *testing.T, result model.EmailVerificationDispatchResult) {
	t.Helper()
	if !result.VerificationRequired {
		t.Fatalf("verification_required = false, want true")
	}
	if result.Message != "If the account is eligible, a verification code will be sent shortly." {
		t.Fatalf("message = %q, want neutral message", result.Message)
	}
	if result.CooldownSeconds != 60 {
		t.Fatalf("cooldown_seconds = %d, want 60", result.CooldownSeconds)
	}
}
