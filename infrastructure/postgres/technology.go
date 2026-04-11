package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marlonlyb/portfolioforge/model"
)

const techTable = "technologies"

var techFields = []string{
	"id",
	"name",
	"slug",
	"category",
	"icon",
	"color",
}

var (
	techPsqlInsert = BuildSQLInsert(techTable, techFields)
	techPsqlUpdate = BuildSQLUpdatedByID(techTable, techFields)
	techPsqlDelete = BuildSQLDelete(techTable)
	techPsqlGetAll = BuildSQLSelect(techTable, techFields)
)

type TechnologyRepository struct {
	db *pgxpool.Pool
}

func NewTechnologyRepository(db *pgxpool.Pool) *TechnologyRepository {
	return &TechnologyRepository{db: db}
}

func (r *TechnologyRepository) Create(m *model.Technology) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.Slug == "" {
		m.Slug = slugifyTech(m.Name)
	}

	_, err := r.db.Exec(
		context.Background(),
		techPsqlInsert,
		m.ID,
		m.Name,
		m.Slug,
		m.Category,
		NullIfEmpty(m.Icon),
		NullIfEmpty(m.Color),
	)
	return err
}

func (r *TechnologyRepository) Update(m *model.Technology) error {
	if m.Slug == "" {
		m.Slug = slugifyTech(m.Name)
	}

	_, err := r.db.Exec(
		context.Background(),
		techPsqlUpdate,
		m.Name,
		m.Slug,
		m.Category,
		NullIfEmpty(m.Icon),
		NullIfEmpty(m.Color),
		m.ID,
	)
	return err
}

func (r *TechnologyRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(
		context.Background(),
		techPsqlDelete,
		id,
	)
	return err
}

func (r *TechnologyRepository) GetByID(id uuid.UUID) (model.Technology, error) {
	query := techPsqlGetAll + " WHERE id = $1"
	row := r.db.QueryRow(context.Background(), query, id)
	return r.scanRow(row)
}

func (r *TechnologyRepository) GetAll() ([]model.Technology, error) {
	query := techPsqlGetAll + " ORDER BY name ASC"
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []model.Technology
	for rows.Next() {
		m, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

func (r *TechnologyRepository) scanRow(s pgx.Row) (model.Technology, error) {
	var m model.Technology
	iconNull := sql.NullString{}
	colorNull := sql.NullString{}

	err := s.Scan(
		&m.ID,
		&m.Name,
		&m.Slug,
		&m.Category,
		&iconNull,
		&colorNull,
	)
	if err != nil {
		return m, err
	}

	m.Icon = iconNull.String
	m.Color = colorNull.String
	return m, nil
}

func slugifyTech(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	for strings.Contains(value, "--") {
		value = strings.ReplaceAll(value, "--", "-")
	}
	return strings.Trim(value, "-")
}
