package user

import (
	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type Repository interface {
	Create(m *model.User) error
	UpsertPasswordlessPublicUser(email, passwordHash string, now int64) (model.User, error)
	GetByID(ID uuid.UUID) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetByProviderSubject(provider, subject string) (model.User, error)
	UpsertGoogleUser(identity model.GoogleIdentity) (model.User, error)
	UpdateLastLogin(ID uuid.UUID, lastLoginAt, updatedAt int64) (model.User, error)
	UpdateProfile(ID uuid.UUID, fullName, company string) (model.User, error)
	CreateEmailVerificationChallenge(challenge *model.EmailVerificationChallenge) error
	GetLatestEmailVerificationChallengeByUserID(userID uuid.UUID) (model.EmailVerificationChallenge, error)
	GetLatestEmailVerificationChallengeByEmail(email string) (model.EmailVerificationChallenge, error)
	UpdateEmailVerificationChallengeAttempt(challengeID uuid.UUID, attemptCount int, updatedAt int64) error
	MarkEmailVerificationChallengeConsumed(challengeID uuid.UUID, consumedAt int64) error
	MarkEmailVerified(userID uuid.UUID, verifiedAt int64) (model.User, error)
	GetAll() (model.Users, error)
	AdminList() (model.Users, error)
	AdminGetByID(ID uuid.UUID) (model.User, error)
	AdminSetIsAdmin(ID uuid.UUID, isAdmin bool, updatedAt int64) (model.User, error)
	AdminSoftDelete(ID uuid.UUID, deletedAt int64) error
}

type Service interface {
	Create(m *model.User) error
	GetByID(ID uuid.UUID) (model.User, error)
	GetByEmail(email string) (model.User, error)
	AdminLogin(email, password string) (model.User, error)
	PublicSignup(email, password string) (model.EmailVerificationDispatchResult, error)
	PublicLogin(email, password string) (model.User, error)
	LoginWithGoogle(identity model.GoogleIdentity) (model.User, error)
	UpdateProfile(ID uuid.UUID, fullName, company string) (model.User, error)
	RequestEmailVerification(email string) (model.EmailVerificationDispatchResult, error)
	ResendEmailVerification(email string) (model.EmailVerificationDispatchResult, error)
	VerifyEmailVerification(email, code string) (model.User, error)
	GetAll() (model.Users, error)
	AdminList() ([]model.AdminUserSummary, error)
	AdminGetByID(ID uuid.UUID) (model.AdminUserDetail, error)
	AdminUpdate(ID uuid.UUID, request model.AdminUserUpdateRequest) (model.AdminUserDetail, error)
	AdminSoftDelete(actorID, targetID uuid.UUID) error
	ToStoreUser(userData model.User) model.StoreUser
}

type ServiceLogin interface {
	AdminLogin(email, password string) (model.User, error)
	PublicSignup(email, password string) (model.EmailVerificationDispatchResult, error)
	PublicLogin(email, password string) (model.User, error)
	LoginWithGoogle(identity model.GoogleIdentity) (model.User, error)
}
