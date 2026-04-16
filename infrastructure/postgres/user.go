package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

const uTable = "users"

var uFields = []string{
	"id",
	"email",
	"password",
	"is_admin",
	"details",
	"auth_provider",
	"provider_subject",
	"email_verified",
	"full_name",
	"company",
	"local_auth_state",
	"last_login_at",
	"created_at",
	"updated_at",
	"deleted_at",
}

var (
	uPsqlInsert = BuildSQLInsert(uTable, uFields)
	uPsqlGetAll = BuildSQLSelect(uTable, uFields)
)

var emailVerificationFields = []string{
	"id",
	"user_id",
	"code_hash",
	"attempt_count",
	"max_attempts",
	"resend_available_at",
	"expires_at",
	"consumed_at",
	"created_at",
	"updated_at",
}

var emailVerificationInsert = BuildSQLInsert("email_verification_challenges", emailVerificationFields)

type User struct {
	db *pgxpool.Pool
}

/* como este es un adaptador, una implementación específica de storage,
aqui si podemos usar las librerías porque específicamente va a conectar a postgres */

func NewUser(db *pgxpool.Pool) *User {
	return &User{db}
}

func (u *User) Create(m *model.User) error {
	_, err := u.db.Exec(
		context.Background(),
		uPsqlInsert,
		m.ID,
		m.Email,
		m.Password,
		m.IsAdmin,
		m.Details,
		NullIfEmpty(m.AuthProvider),
		NullIfEmpty(m.ProviderSubject),
		m.EmailVerified,
		NullIfEmpty(m.FullName),
		NullIfEmpty(m.Company),
		NullIfEmpty(m.LocalAuthState),
		Int64ToNull(m.LastLoginAt),
		m.CreatedAt,
		Int64ToNull(m.UpdatedAt),
		Int64ToNull(m.DeletedAt),
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpsertPasswordlessPublicUser(email, passwordHash string, now int64) (model.User, error) {
	userID, err := uuid.NewUUID()
	if err != nil {
		return model.User{}, err
	}

	details, err := json.Marshal(map[string]any{})
	if err != nil {
		return model.User{}, err
	}

	row := u.db.QueryRow(context.Background(), `
		INSERT INTO users (
			id, email, password, is_admin, details, auth_provider, provider_subject, email_verified, full_name, company, local_auth_state, last_login_at, created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, FALSE, $4, 'local', NULL, FALSE, NULL, NULL, 'password_setup_required', NULL, $5, $6, NULL
		)
		ON CONFLICT (email) DO UPDATE
		SET updated_at = users.updated_at
		WHERE users.deleted_at IS NULL
		RETURNING `+strings.Join(uFields, ", "), userID, email, passwordHash, details, now, now)

	return u.scanRow(row, true)
}

func (u *User) GetByID(ID uuid.UUID) (model.User, error) {
	query := uPsqlGetAll + " WHERE id = $1 AND deleted_at IS NULL"
	row := u.db.QueryRow(
		context.Background(),
		query,
		ID,
	)

	return u.scanRow(row, false)
}

func (u *User) GetByEmail(email string) (model.User, error) {
	query := uPsqlGetAll + " WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL"
	row := u.db.QueryRow(
		context.Background(),
		query,
		email,
	)

	return u.scanRow(row, true)
}

func (u *User) GetByProviderSubject(provider, subject string) (model.User, error) {
	query := uPsqlGetAll + " WHERE auth_provider = $1 AND provider_subject = $2 AND deleted_at IS NULL"
	row := u.db.QueryRow(context.Background(), query, provider, subject)
	return u.scanRow(row, true)
}

func (u *User) UpsertGoogleUser(identity model.GoogleIdentity) (model.User, error) {
	ctx := context.Background()
	trimmedSubject := strings.TrimSpace(identity.Subject)
	trimmedEmail := strings.ToLower(strings.TrimSpace(identity.Email))
	trimmedFullName := strings.TrimSpace(identity.FullName)
	now := time.Now().Unix()

	if trimmedSubject == "" || trimmedEmail == "" {
		return model.User{}, errors.New("google identity is incomplete")
	}

	if existing, err := u.GetByProviderSubject("google", trimmedSubject); err == nil {
		row := u.db.QueryRow(ctx, `
			UPDATE users
			SET email = $2,
			    email_verified = $3,
			    full_name = COALESCE(NULLIF($4, ''), full_name),
			    last_login_at = $5,
			    updated_at = $6
			WHERE id = $1
			RETURNING `+strings.Join(uFields, ", "), existing.ID, trimmedEmail, identity.EmailVerified, trimmedFullName, now, now)
		return u.scanRow(row, true)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, err
	}

	if existing, err := u.GetByEmail(trimmedEmail); err == nil {
		authProvider := strings.TrimSpace(existing.AuthProvider)
		if authProvider == "" {
			authProvider = "local"
		}
		if authProvider != "google" {
			return model.User{}, model.ErrProviderConflict
		}

		row := u.db.QueryRow(ctx, `
			UPDATE users
			SET provider_subject = $2,
			    email_verified = $3,
			    full_name = COALESCE(NULLIF($4, ''), full_name),
			    last_login_at = $5,
			    updated_at = $6
			WHERE id = $1
			RETURNING `+strings.Join(uFields, ", "), existing.ID, trimmedSubject, identity.EmailVerified, trimmedFullName, now, now)
		return u.scanRow(row, true)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, err
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return model.User{}, err
	}

	details, err := json.Marshal(map[string]any{})
	if err != nil {
		return model.User{}, err
	}

	row := u.db.QueryRow(ctx, `
		INSERT INTO users (
			id, email, password, is_admin, details, auth_provider, provider_subject, email_verified, full_name, company, local_auth_state, last_login_at, created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, NULL, FALSE, $3, 'google', $4, $5, NULLIF($6, ''), NULL, 'ready', $7, $8, $9, NULL
		)
		RETURNING `+strings.Join(uFields, ", "), userID, trimmedEmail, details, trimmedSubject, identity.EmailVerified, trimmedFullName, now, now, now)

	return u.scanRow(row, true)
}

func (u *User) UpdateProfile(ID uuid.UUID, fullName, company string) (model.User, error) {
	now := time.Now().Unix()
	row := u.db.QueryRow(context.Background(), `
		UPDATE users
		SET full_name = $2,
		    company = $3,
		    updated_at = $4
		WHERE id = $1
		  AND deleted_at IS NULL
		RETURNING `+strings.Join(uFields, ", "), ID, strings.TrimSpace(fullName), strings.TrimSpace(company), now)
	return u.scanRow(row, true)
}

func (u *User) UpdateLastLogin(ID uuid.UUID, lastLoginAt, updatedAt int64) (model.User, error) {
	row := u.db.QueryRow(context.Background(), `
		UPDATE users
		SET last_login_at = $2,
		    updated_at = $3
		WHERE id = $1
		  AND deleted_at IS NULL
		RETURNING `+strings.Join(uFields, ", "), ID, lastLoginAt, updatedAt)
	return u.scanRow(row, true)
}

func (u *User) CreateEmailVerificationChallenge(challenge *model.EmailVerificationChallenge) error {
	_, err := u.db.Exec(
		context.Background(),
		emailVerificationInsert,
		challenge.ID,
		challenge.UserID,
		challenge.CodeHash,
		challenge.AttemptCount,
		challenge.MaxAttempts,
		challenge.ResendAvailableAt,
		challenge.ExpiresAt,
		Int64ToNull(challenge.ConsumedAt),
		challenge.CreatedAt,
		challenge.UpdatedAt,
	)
	return err
}

func (u *User) GetLatestEmailVerificationChallengeByUserID(userID uuid.UUID) (model.EmailVerificationChallenge, error) {
	row := u.db.QueryRow(context.Background(), `
		SELECT `+strings.Join(emailVerificationFields, ", ")+`
		FROM email_verification_challenges
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`, userID)
	return scanEmailVerificationChallenge(row)
}

func (u *User) GetLatestEmailVerificationChallengeByEmail(email string) (model.EmailVerificationChallenge, error) {
	row := u.db.QueryRow(context.Background(), `
		SELECT c.`+strings.Join(emailVerificationFields, ", c.")+`
		FROM email_verification_challenges c
		INNER JOIN users u ON u.id = c.user_id
		WHERE LOWER(u.email) = LOWER($1)
		  AND COALESCE(u.auth_provider, 'local') = 'local'
		  AND u.is_admin = FALSE
		  AND u.deleted_at IS NULL
		ORDER BY c.created_at DESC
		LIMIT 1`, email)
	return scanEmailVerificationChallenge(row)
}

func (u *User) UpdateEmailVerificationChallengeAttempt(challengeID uuid.UUID, attemptCount int, updatedAt int64) error {
	_, err := u.db.Exec(context.Background(), `
		UPDATE email_verification_challenges
		SET attempt_count = $2,
		    updated_at = $3
		WHERE id = $1`, challengeID, attemptCount, updatedAt)
	return err
}

func (u *User) MarkEmailVerificationChallengeConsumed(challengeID uuid.UUID, consumedAt int64) error {
	_, err := u.db.Exec(context.Background(), `
		UPDATE email_verification_challenges
		SET consumed_at = $2,
		    updated_at = $2
		WHERE id = $1`, challengeID, consumedAt)
	return err
}

func (u *User) MarkEmailVerified(userID uuid.UUID, verifiedAt int64) (model.User, error) {
	row := u.db.QueryRow(context.Background(), `
		UPDATE users
		SET email_verified = TRUE,
		    updated_at = $2
		WHERE id = $1
		  AND deleted_at IS NULL
		RETURNING `+strings.Join(uFields, ", "), userID, verifiedAt)
	return u.scanRow(row, true)
}

func (u *User) GetAll() (model.Users, error) {
	rows, err := u.db.Query(
		context.Background(),
		uPsqlGetAll+" WHERE deleted_at IS NULL ORDER BY created_at DESC, email ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := model.Users{}
	for rows.Next() {
		m, err := u.scanRow(rows, false)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (u *User) AdminList() (model.Users, error) {
	return u.GetAll()
}

func (u *User) AdminGetByID(ID uuid.UUID) (model.User, error) {
	return u.GetByID(ID)
}

func (u *User) AdminSetIsAdmin(ID uuid.UUID, isAdmin bool, updatedAt int64) (model.User, error) {
	row := u.db.QueryRow(context.Background(), `
		UPDATE users
		SET is_admin = $2,
		    updated_at = $3
		WHERE id = $1
		  AND deleted_at IS NULL
		RETURNING `+strings.Join(uFields, ", "), ID, isAdmin, updatedAt)
	return u.scanRow(row, true)
}

func (u *User) AdminSoftDelete(ID uuid.UUID, deletedAt int64) error {
	commandTag, err := u.db.Exec(context.Background(), `
		UPDATE users
		SET deleted_at = $2,
		    updated_at = $2
		WHERE id = $1
		  AND deleted_at IS NULL`, ID, deletedAt)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (u *User) scanRow(s pgx.Row, withPassword bool) (model.User, error) {
	m := model.User{}

	password := sql.NullString{}
	providerSubject := sql.NullString{}
	fullName := sql.NullString{}
	company := sql.NullString{}
	localAuthState := sql.NullString{}
	lastLoginAt := sql.NullInt64{}
	updateAtNull := sql.NullInt64{}
	deletedAtNull := sql.NullInt64{}

	err := s.Scan(
		&m.ID,
		&m.Email,
		&password,
		&m.IsAdmin,
		&m.Details,
		&m.AuthProvider,
		&providerSubject,
		&m.EmailVerified,
		&fullName,
		&company,
		&localAuthState,
		&lastLoginAt,
		&m.CreatedAt,
		&updateAtNull,
		&deletedAtNull,
	)
	if err != nil {
		return m, err
	}

	m.Password = password.String
	m.ProviderSubject = providerSubject.String
	m.FullName = fullName.String
	m.Company = company.String
	m.LocalAuthState = localAuthState.String
	m.LastLoginAt = lastLoginAt.Int64
	m.UpdatedAt = updateAtNull.Int64
	m.DeletedAt = deletedAtNull.Int64

	if !withPassword {
		m.Password = ""
	}

	return m, nil
}

func scanEmailVerificationChallenge(s pgx.Row) (model.EmailVerificationChallenge, error) {
	challenge := model.EmailVerificationChallenge{}
	consumedAt := sql.NullInt64{}

	err := s.Scan(
		&challenge.ID,
		&challenge.UserID,
		&challenge.CodeHash,
		&challenge.AttemptCount,
		&challenge.MaxAttempts,
		&challenge.ResendAvailableAt,
		&challenge.ExpiresAt,
		&consumedAt,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
	)
	if err != nil {
		return challenge, err
	}

	challenge.ConsumedAt = consumedAt.Int64
	return challenge, nil
}
