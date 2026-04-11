package postgres

import (
	"context"
	"database/sql"

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
	"created_at",
	"updated_at",
}

var (
	uPsqlInsert = BuildSQLInsert(uTable, uFields)
	uPsqlGetAll = BuildSQLSelect(uTable, uFields)
)

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
		m.CreatedAt,
		Int64ToNull(m.UpdatedAt),
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetByID(ID uuid.UUID) (model.User, error) {
	query := uPsqlGetAll + " WHERE id = $1"
	row := u.db.QueryRow(
		context.Background(),
		query,
		ID,
	)

	return u.scanRow(row, false)
}

func (u *User) GetByEmail(email string) (model.User, error) {
	query := uPsqlGetAll + " WHERE email = $1"
	row := u.db.QueryRow(
		context.Background(),
		query,
		email,
	)

	return u.scanRow(row, true)
}

func (u *User) GetAll() (model.Users, error) {
	rows, err := u.db.Query(
		context.Background(),
		uPsqlGetAll,
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

func (u *User) scanRow(s pgx.Row, withPassword bool) (model.User, error) {
	m := model.User{}

	updateAtNull := sql.NullInt64{}

	err := s.Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.IsAdmin,
		&m.Details,
		&m.CreatedAt,
		&updateAtNull,
	)
	if err != nil {
		return m, err
	}

	m.UpdatedAt = updateAtNull.Int64

	if !withPassword {
		m.Password = ""
	}

	return m, nil
}
