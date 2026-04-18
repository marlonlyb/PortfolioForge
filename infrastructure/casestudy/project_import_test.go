package casestudy

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	projectPorts "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/model"
)

var _ projectPorts.AdminCatalogService = (*stubAdminCatalogService)(nil)

type stubAdminCatalogService struct {
	projects            []model.AdminProject
	created             *model.AdminProjectWrite
	updated             *model.AdminProjectWrite
	replacedMedia       []model.ProjectMedia
	replacedMediaID     uuid.UUID
	replaceMediaInvoked int
}

func (s *stubAdminCatalogService) Create(m *model.AdminProjectWrite) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	}
	copy := *m
	s.created = &copy
	return nil
}
func (s *stubAdminCatalogService) Update(m *model.AdminProjectWrite) error {
	copy := *m
	s.updated = &copy
	return nil
}
func (s *stubAdminCatalogService) Delete(id uuid.UUID) error { return nil }
func (s *stubAdminCatalogService) UpdateStatus(id uuid.UUID, active bool) (model.AdminProject, error) {
	return model.AdminProject{}, nil
}
func (s *stubAdminCatalogService) CreateVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error {
	return nil
}
func (s *stubAdminCatalogService) ReplaceVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error {
	return nil
}
func (s *stubAdminCatalogService) ReplaceMedia(projectID uuid.UUID, media []model.ProjectMedia) error {
	s.replacedMediaID = projectID
	s.replaceMediaInvoked++
	s.replacedMedia = append([]model.ProjectMedia(nil), media...)
	return nil
}
func (s *stubAdminCatalogService) GetAdminByID(id uuid.UUID) (model.AdminProject, error) {
	return model.AdminProject{}, nil
}
func (s *stubAdminCatalogService) GetAdminAll() ([]model.AdminProject, error) {
	return append([]model.AdminProject(nil), s.projects...), nil
}

func TestImportFromCanonicalCreatesSevenDefaultMediaItems(t *testing.T) {
	service := &stubAdminCatalogService{}
	importer := NewProjectImporter(service)
	canonicalPath := writeCanonicalMarkdown(t, "printer-05-controls-migration", "# Printer 05 Controls Migration\n\nResumen del caso.")

	projectID, err := importer.ImportFromCanonical(context.Background(), PublishTarget{
		Slug:      "printer-05-controls-migration",
		LocalFile: canonicalPath,
	}, "https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/printer-05-controls-migration.md")
	if err != nil {
		t.Fatalf("ImportFromCanonical() error = %v", err)
	}
	if projectID == uuid.Nil {
		t.Fatalf("projectID is nil")
	}
	if service.created == nil {
		t.Fatalf("expected Create to be called")
	}
	if service.replaceMediaInvoked != 1 {
		t.Fatalf("expected ReplaceMedia to be called once, got %d", service.replaceMediaInvoked)
	}
	if len(service.replacedMedia) != defaultProjectMediaCount {
		t.Fatalf("len(replacedMedia) = %d, want %d", len(service.replacedMedia), defaultProjectMediaCount)
	}
	if !service.replacedMedia[0].Featured {
		t.Fatalf("first media item should be featured")
	}
	if got := service.replacedMedia[0].LowURL; got != "https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/imagen01_low.webp" {
		t.Fatalf("first low url = %q", got)
	}
	if got := service.replacedMedia[6].HighURL; got != "https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/imagen07_high.webp" {
		t.Fatalf("seventh high url = %q", got)
	}
	if got := service.replacedMedia[0].FallbackURL; got != projectMediaFallbackURL {
		t.Fatalf("fallback = %q", got)
	}
	if len(service.created.Media) != defaultProjectMediaCount {
		t.Fatalf("payload media len = %d", len(service.created.Media))
	}
}

func TestImportFromCanonicalBackfillsExistingProjectMediaToSeven(t *testing.T) {
	existingID := uuid.MustParse("165cd616-d471-4b4c-9814-ed5a09bfc31e")
	service := &stubAdminCatalogService{
		projects: []model.AdminProject{{
			ID:   existingID,
			Slug: "printer-05-controls-migration",
			Name: "Printer 05 Controls Migration",
			Media: []model.ProjectMedia{{
				ProjectID:   existingID,
				MediaType:   "image",
				LowURL:      "https://custom.example.com/imagen01_low.webp",
				MediumURL:   "https://custom.example.com/imagen01_medium.webp",
				HighURL:     "https://custom.example.com/imagen01_high.webp",
				FallbackURL: projectMediaFallbackURL,
				SortOrder:   0,
				Featured:    true,
			}},
			Active: true,
		}},
	}
	importer := NewProjectImporter(service)
	canonicalPath := writeCanonicalMarkdown(t, "printer-05-controls-migration", "# Printer 05 Controls Migration\n\nResumen del caso.")

	_, err := importer.ImportFromCanonical(context.Background(), PublishTarget{
		Slug:      "printer-05-controls-migration",
		LocalFile: canonicalPath,
	}, "https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/printer-05-controls-migration.md")
	if err != nil {
		t.Fatalf("ImportFromCanonical() error = %v", err)
	}
	if service.updated == nil {
		t.Fatalf("expected Update to be called")
	}
	if len(service.replacedMedia) != defaultProjectMediaCount {
		t.Fatalf("len(replacedMedia) = %d, want %d", len(service.replacedMedia), defaultProjectMediaCount)
	}
	if service.replacedMedia[0].LowURL != "https://custom.example.com/imagen01_low.webp" {
		t.Fatalf("existing media should be preserved, got %q", service.replacedMedia[0].LowURL)
	}
	if service.replacedMedia[6].LowURL != "https://mlbautomation.com/dev/portfolioforge/printer-05-controls-migration/imagen07_low.webp" {
		t.Fatalf("generated low url = %q", service.replacedMedia[6].LowURL)
	}
	if service.replacedMediaID != existingID {
		t.Fatalf("ReplaceMedia project id = %s, want %s", service.replacedMediaID, existingID)
	}
	if len(service.updated.Media) != defaultProjectMediaCount {
		t.Fatalf("updated payload media len = %d", len(service.updated.Media))
	}
	if len(service.updated.Images) == 0 {
		t.Fatalf("expected legacy images to be generated")
	}
}

func writeCanonicalMarkdown(t *testing.T, slug, content string) string {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, slug+".md")
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("write canonical markdown: %v", err)
	}
	return file
}
