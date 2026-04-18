package casestudy

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	projectPorts "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/model"
)

type ProjectImporter struct {
	service projectPorts.AdminCatalogService
}

func NewProjectImporter(service projectPorts.AdminCatalogService) *ProjectImporter {
	return &ProjectImporter{service: service}
}

func (i *ProjectImporter) ImportFromCanonical(ctx context.Context, source PublishTarget, canonicalURL string) (uuid.UUID, error) {
	_ = ctx
	if i == nil || i.service == nil {
		return uuid.Nil, fmt.Errorf("project importer not configured")
	}

	base, err := loadCanonicalBaseContent(source.LocalFile, source.Slug)
	if err != nil {
		return uuid.Nil, err
	}

	existing, err := i.findBySlug(source.Slug)
	if err != nil {
		return uuid.Nil, err
	}

	if existing != nil {
		payload := model.AdminProjectWrite{
			ID:                existing.ID,
			Name:              base.Name,
			Description:       base.Description,
			Category:          base.Category,
			ClientName:        base.ClientName,
			SourceMarkdownURL: canonicalURL,
			Media:             adminMediaFromProject(existing.Media),
			Variants:          adminVariantsFromProject(existing.Variants),
			Active:            boolPtr(existing.Active),
		}
		if len(existing.Images) > 0 {
			payload.Images = mustMarshalRaw(existing.Images)
		}
		if err := i.service.Update(&payload); err != nil {
			return uuid.Nil, fmt.Errorf("update existing project %s: %w", existing.ID, err)
		}
		return existing.ID, nil
	}

	payload := model.AdminProjectWrite{
		Name:              base.Name,
		Description:       base.Description,
		Category:          base.Category,
		ClientName:        base.ClientName,
		SourceMarkdownURL: canonicalURL,
		Images:            mustMarshalRaw([]string{}),
		Features:          mustMarshalRaw([]string{}),
		Active:            boolPtr(false),
	}
	if err := i.service.Create(&payload); err != nil {
		return uuid.Nil, fmt.Errorf("create project from canonical %s: %w", source.Slug, err)
	}
	return payload.ID, nil
}

func (i *ProjectImporter) findBySlug(slug string) (*model.AdminProject, error) {
	projects, err := i.service.GetAdminAll()
	if err != nil {
		return nil, fmt.Errorf("load admin projects: %w", err)
	}

	for _, project := range projects {
		if strings.EqualFold(strings.TrimSpace(project.Slug), strings.TrimSpace(slug)) {
			copy := project
			return &copy, nil
		}
	}

	return nil, nil
}

type canonicalBaseContent struct {
	Name        string
	Description string
	Category    string
	ClientName  string
}

func loadCanonicalBaseContent(markdownPath, fallbackSlug string) (canonicalBaseContent, error) {
	handle, err := os.Open(markdownPath)
	if err != nil {
		return canonicalBaseContent{}, fmt.Errorf("open canonical markdown %s: %w", markdownPath, err)
	}
	defer handle.Close()

	base := canonicalBaseContent{Category: "case-study"}
	frontmatter := map[string]string{}
	collectFrontmatter := false
	frontmatterDone := false
	paragraphs := make([]string, 0, 2)

	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if !frontmatterDone && line == "---" {
			collectFrontmatter = !collectFrontmatter
			if !collectFrontmatter {
				frontmatterDone = true
			}
			continue
		}

		if collectFrontmatter {
			key, value, ok := strings.Cut(line, ":")
			if ok {
				frontmatter[strings.ToLower(strings.TrimSpace(key))] = strings.Trim(strings.TrimSpace(value), `"'`)
			}
			continue
		}

		if strings.HasPrefix(line, "# ") && base.Name == "" {
			base.Name = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}

		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "!") && !strings.HasPrefix(line, "[") {
			paragraphs = append(paragraphs, line)
			if len(paragraphs) >= 2 {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return canonicalBaseContent{}, fmt.Errorf("read canonical markdown %s: %w", markdownPath, err)
	}

	if title := firstNonEmpty(frontmatter["title"], frontmatter["name"]); title != "" {
		base.Name = title
	}
	if description := firstNonEmpty(frontmatter["description"], frontmatter["summary"], strings.Join(paragraphs, " ")); description != "" {
		base.Description = description
	}
	if category := firstNonEmpty(frontmatter["category"], frontmatter["type"]); category != "" {
		base.Category = category
	}
	base.ClientName = firstNonEmpty(frontmatter["client_name"], frontmatter["client"], frontmatter["brand"])

	if strings.TrimSpace(base.Name) == "" {
		base.Name = humanizeSlug(filepath.Base(fallbackSlug))
	}
	if strings.TrimSpace(base.Description) == "" {
		base.Description = fmt.Sprintf("Imported from canonical case-study source %s.", fallbackSlug)
	}

	return base, nil
}

func adminMediaFromProject(items []model.ProjectMedia) []model.AdminProjectMediaInput {
	media := make([]model.AdminProjectMediaInput, 0, len(items))
	for _, item := range items {
		media = append(media, model.AdminProjectMediaInput{
			ID:          item.ID.String(),
			MediaType:   item.MediaType,
			FallbackURL: item.FallbackURL,
			LowURL:      item.LowURL,
			MediumURL:   item.MediumURL,
			HighURL:     item.HighURL,
			Caption:     item.Caption,
			AltText:     item.AltText,
			SortOrder:   item.SortOrder,
			Featured:    item.Featured,
		})
	}
	return media
}

func adminVariantsFromProject(items []model.AdminProjectVariant) []model.AdminProjectVariantInput {
	variants := make([]model.AdminProjectVariantInput, 0, len(items))
	for _, item := range items {
		variants = append(variants, model.AdminProjectVariantInput{
			ID:       item.ID.String(),
			SKU:      item.SKU,
			Color:    item.Color,
			Size:     item.Size,
			Price:    item.Price,
			Stock:    item.Stock,
			ImageURL: item.ImageURL,
		})
	}
	return variants
}

func humanizeSlug(slug string) string {
	parts := strings.FieldsFunc(strings.TrimSpace(slug), func(r rune) bool {
		return r == '-' || r == '_'
	})
	for index := range parts {
		if parts[index] == "" {
			continue
		}
		parts[index] = strings.ToUpper(parts[index][:1]) + strings.ToLower(parts[index][1:])
	}
	if len(parts) == 0 {
		return "Imported Case Study"
	}
	return strings.Join(parts, " ")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func boolPtr(value bool) *bool {
	return &value
}

func mustMarshalRaw(value interface{}) []byte {
	encoded, _ := json.Marshal(value)
	return encoded
}
