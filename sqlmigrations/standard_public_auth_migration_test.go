package sqlmigrations

import (
	"os"
	"strings"
	"testing"
)

const passwordlessPublicAuthStartedAtUnix int64 = 1776247200

type migrationUserFixture struct {
	AuthProvider  string
	IsAdmin       bool
	EmailVerified bool
	FullName      string
	Company       string
	LastLoginAt   int64
	CreatedAt     int64
	DetailsJSON   string
}

func TestStandardPublicAuthMigrationFlagsOnlyLegacyPasswordlessLocals(t *testing.T) {
	legacyPasswordlessLocal := migrationUserFixture{
		AuthProvider: "local",
		CreatedAt:    passwordlessPublicAuthStartedAtUnix + 300,
		DetailsJSON:  `{}`,
	}

	if !migrationShouldRequirePasswordSetup(legacyPasswordlessLocal) {
		t.Fatalf("legacy passwordless local user should require password setup")
	}
}

func TestStandardPublicAuthMigrationKeepsExistingPasswordUsersReady(t *testing.T) {
	existingPasswordUser := migrationUserFixture{
		AuthProvider: "local",
		CreatedAt:    passwordlessPublicAuthStartedAtUnix - 300,
		DetailsJSON:  `{}`,
	}

	if migrationShouldRequirePasswordSetup(existingPasswordUser) {
		t.Fatalf("existing local password user should stay ready")
	}
}

func TestStandardPublicAuthMigrationSQLIncludesSafeLegacyPredicate(t *testing.T) {
	sqlBytes, err := os.ReadFile("20260416_1400_standard_public_auth.sql")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	sql := string(sqlBytes)
	requiredFragments := []string{
		"created_at >= 1776247200",
		"email_verified = FALSE",
		"COALESCE(NULLIF(TRIM(full_name), ''), '') = ''",
		"COALESCE(NULLIF(TRIM(company), ''), '') = ''",
		"COALESCE(last_login_at, 0) = 0",
		"COALESCE(details, '{}'::jsonb) = '{}'::jsonb",
	}

	for _, fragment := range requiredFragments {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("migration SQL missing safety fragment %q", fragment)
		}
	}

	if strings.Contains(sql, "AND is_admin = FALSE THEN 'password_setup_required'") {
		t.Fatalf("migration SQL still contains the broad local-user catch-all")
	}
}

func migrationShouldRequirePasswordSetup(user migrationUserFixture) bool {
	authProvider := strings.TrimSpace(user.AuthProvider)
	if authProvider == "" {
		authProvider = "local"
	}

	return authProvider == "local" &&
		!user.IsAdmin &&
		user.CreatedAt >= passwordlessPublicAuthStartedAtUnix &&
		!user.EmailVerified &&
		strings.TrimSpace(user.FullName) == "" &&
		strings.TrimSpace(user.Company) == "" &&
		user.LastLoginAt == 0 &&
		strings.TrimSpace(user.DetailsJSON) == `{}`
}
