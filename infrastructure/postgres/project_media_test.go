package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubProjectMediaQueryer struct {
	query string
	args  []interface{}
	rows  pgx.Rows
	err   error
}

func (s *stubProjectMediaQueryer) Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	s.query = sql
	s.args = args
	if s.err != nil {
		return nil, s.err
	}
	return s.rows, nil
}

type stubProjectMediaExecer struct {
	queries []string
	args    [][]interface{}
	errAt   map[int]error
}

func (s *stubProjectMediaExecer) Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	s.queries = append(s.queries, sql)
	s.args = append(s.args, args)
	idx := len(s.queries) - 1
	if err, exists := s.errAt[idx]; exists {
		return pgconn.CommandTag{}, err
	}
	return pgconn.CommandTag{}, nil
}

type stubProjectMediaRows struct {
	index int
	rows  [][]any
	err   error
}

func (s *stubProjectMediaRows) Close() {}

func (s *stubProjectMediaRows) Err() error { return s.err }

func (s *stubProjectMediaRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }

func (s *stubProjectMediaRows) FieldDescriptions() []pgconn.FieldDescription { return nil }

func (s *stubProjectMediaRows) Next() bool {
	if s.index >= len(s.rows) {
		return false
	}
	s.index++
	return true
}

func (s *stubProjectMediaRows) Scan(dest ...any) error {
	if s.index == 0 || s.index > len(s.rows) {
		return errors.New("scan called without active row")
	}

	row := s.rows[s.index-1]
	if len(dest) != len(row) {
		return errors.New("scan destination count mismatch")
	}

	for i, value := range row {
		switch target := dest[i].(type) {
		case *uuid.UUID:
			v, ok := value.(uuid.UUID)
			if !ok {
				return errors.New("expected uuid value")
			}
			*target = v
		case *string:
			v, ok := value.(string)
			if !ok {
				return errors.New("expected string value")
			}
			*target = v
		case *int:
			v, ok := value.(int)
			if !ok {
				return errors.New("expected int value")
			}
			*target = v
		case *bool:
			v, ok := value.(bool)
			if !ok {
				return errors.New("expected bool value")
			}
			*target = v
		default:
			return errors.New("unsupported scan destination")
		}
	}

	return nil
}

func (s *stubProjectMediaRows) Values() ([]any, error) {
	if s.index == 0 || s.index > len(s.rows) {
		return nil, errors.New("values called without active row")
	}
	return s.rows[s.index-1], nil
}

func (s *stubProjectMediaRows) RawValues() [][]byte { return nil }

func (s *stubProjectMediaRows) Conn() *pgx.Conn { return nil }

func TestFetchProjectMediaMapsLegacySQLColumnsToCanonicalFields(t *testing.T) {
	projectID := uuid.New()
	mediaID := uuid.New()
	queryer := &stubProjectMediaQueryer{
		rows: &stubProjectMediaRows{rows: [][]any{{
			mediaID,
			projectID,
			"image",
			"https://cdn.example.com/project-fallback.webp",
			"https://cdn.example.com/project-low.webp",
			"https://cdn.example.com/project-medium.webp",
			"https://cdn.example.com/project-high.webp",
			"Hero shot",
			"PortfolioForge hero",
			0,
			true,
		}}},
	}

	media, err := fetchProjectMedia(context.Background(), queryer, projectID)
	if err != nil {
		t.Fatalf("fetchProjectMedia() error = %v", err)
	}
	if len(media) != 1 {
		t.Fatalf("len(media) = %d, want 1", len(media))
	}

	item := media[0]
	if item.LowURL != "https://cdn.example.com/project-low.webp" {
		t.Fatalf("LowURL = %q", item.LowURL)
	}
	if item.MediumURL != "https://cdn.example.com/project-medium.webp" {
		t.Fatalf("MediumURL = %q", item.MediumURL)
	}
	if item.HighURL != "https://cdn.example.com/project-high.webp" {
		t.Fatalf("HighURL = %q", item.HighURL)
	}
	if item.FallbackURL != "https://cdn.example.com/project-fallback.webp" {
		t.Fatalf("FallbackURL = %q", item.FallbackURL)
	}
	if len(queryer.args) != 1 || queryer.args[0] != projectID {
		t.Fatalf("query args = %#v, want project id", queryer.args)
	}
	for _, want := range []string{"COALESCE(url, '')", "COALESCE(thumbnail_url, '')", "COALESCE(medium_url, '')", "COALESCE(full_url, '')"} {
		if !strings.Contains(queryer.query, want) {
			t.Fatalf("query %q does not include %q", queryer.query, want)
		}
	}
}

func TestReplaceProjectMediaMapsCanonicalFieldsToLegacySQLColumns(t *testing.T) {
	projectID := uuid.New()
	mediaID := uuid.New()
	execer := &stubProjectMediaExecer{}

	err := replaceProjectMedia(context.Background(), execer, projectID, []model.ProjectMedia{{
		ID:          mediaID,
		ProjectID:   projectID,
		MediaType:   "image",
		LowURL:      "https://cdn.example.com/project-low.webp",
		MediumURL:   "https://cdn.example.com/project-medium.webp",
		HighURL:     "https://cdn.example.com/project-high.webp",
		FallbackURL: "https://cdn.example.com/project-fallback.webp",
		Caption:     "Hero shot",
		AltText:     "PortfolioForge hero",
		SortOrder:   3,
		Featured:    true,
	}})
	if err != nil {
		t.Fatalf("replaceProjectMedia() error = %v", err)
	}
	if len(execer.queries) != 2 {
		t.Fatalf("exec count = %d, want 2", len(execer.queries))
	}
	if !strings.Contains(execer.queries[1], "url, thumbnail_url, medium_url, full_url") {
		t.Fatalf("insert query does not preserve SQL column names: %q", execer.queries[1])
	}

	insertArgs := execer.args[1]
	if len(insertArgs) != 11 {
		t.Fatalf("insert args len = %d, want 11", len(insertArgs))
	}
	if insertArgs[0] != mediaID {
		t.Fatalf("insert id = %#v", insertArgs[0])
	}
	if insertArgs[1] != projectID {
		t.Fatalf("insert project id = %#v", insertArgs[1])
	}
	if insertArgs[3] != "https://cdn.example.com/project-fallback.webp" {
		t.Fatalf("url arg = %#v", insertArgs[3])
	}
	if insertArgs[4] != "https://cdn.example.com/project-low.webp" {
		t.Fatalf("thumbnail_url arg = %#v", insertArgs[4])
	}
	if insertArgs[5] != "https://cdn.example.com/project-medium.webp" {
		t.Fatalf("medium_url arg = %#v", insertArgs[5])
	}
	if insertArgs[6] != "https://cdn.example.com/project-high.webp" {
		t.Fatalf("full_url arg = %#v", insertArgs[6])
	}
}
