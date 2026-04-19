package casestudy

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	projectPorts "github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/model"
)

const canonicalMetadataHeading = "metadata"

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
		projectMedia := assignProjectIDToMedia(ensureMinimumProjectMedia(existing.Media, source.Slug, defaultProjectMediaCount), existing.ID)
		legacyImages := model.BuildProjectImageList(projectMedia, existing.Images)

		payload := model.AdminProjectWrite{
			ID:                existing.ID,
			Name:              base.Name,
			Description:       base.Description,
			Category:          base.Category,
			ClientName:        base.ClientName,
			IndustryType:      base.IndustryType,
			FinalProduct:      base.FinalProduct,
			SourceMarkdownURL: canonicalURL,
			Media:             adminMediaFromProject(projectMedia),
			Variants:          adminVariantsFromProject(existing.Variants),
			Active:            boolPtr(existing.Active),
			Images:            mustMarshalRaw(legacyImages),
		}
		if err := i.service.Update(&payload); err != nil {
			return uuid.Nil, fmt.Errorf("update existing project %s: %w", existing.ID, err)
		}
		if err := i.service.ReplaceMedia(existing.ID, projectMedia); err != nil {
			return uuid.Nil, fmt.Errorf("replace project media %s: %w", existing.ID, err)
		}
		return existing.ID, nil
	}

	projectMedia := assignProjectIDToMedia(ensureMinimumProjectMedia(nil, source.Slug, defaultProjectMediaCount), uuid.Nil)
	legacyImages := model.BuildProjectImageList(projectMedia, nil)

	payload := model.AdminProjectWrite{
		Name:              base.Name,
		Description:       base.Description,
		Category:          base.Category,
		ClientName:        base.ClientName,
		IndustryType:      base.IndustryType,
		FinalProduct:      base.FinalProduct,
		SourceMarkdownURL: canonicalURL,
		Media:             adminMediaFromProject(projectMedia),
		Images:            mustMarshalRaw(legacyImages),
		Features:          mustMarshalRaw([]string{}),
		Active:            boolPtr(false),
	}
	if err := i.service.Create(&payload); err != nil {
		return uuid.Nil, fmt.Errorf("create project from canonical %s: %w", source.Slug, err)
	}
	projectMedia = assignProjectIDToMedia(projectMedia, payload.ID)
	if err := i.service.ReplaceMedia(payload.ID, projectMedia); err != nil {
		return uuid.Nil, fmt.Errorf("replace project media %s: %w", payload.ID, err)
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

type CanonicalProjectMetadata struct {
	Name         string
	Description  string
	Category     string
	ClientName   string
	IndustryType string
	FinalProduct string
}

func loadCanonicalBaseContent(markdownPath, fallbackSlug string) (CanonicalProjectMetadata, error) {
	body, err := os.ReadFile(markdownPath)
	if err != nil {
		return CanonicalProjectMetadata{}, fmt.Errorf("open canonical markdown %s: %w", markdownPath, err)
	}
	return ParseCanonicalProjectMetadataReader(bytesReader(body), fallbackSlug)
}

func ParseCanonicalProjectMetadata(markdown string, fallbackSlug string) (CanonicalProjectMetadata, error) {
	return ParseCanonicalProjectMetadataReader(strings.NewReader(markdown), fallbackSlug)
}

func ParseCanonicalProjectMetadataReader(reader io.Reader, fallbackSlug string) (CanonicalProjectMetadata, error) {
	base := CanonicalProjectMetadata{Category: "case-study"}
	frontmatter := map[string]string{}
	collectFrontmatter := false
	frontmatterDone := false
	collectMetadata := false
	metadata := map[string]string{}
	paragraphs := make([]string, 0, 2)

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)
		if line == "" {
			if collectMetadata {
				collectMetadata = false
			}
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

		if strings.HasPrefix(line, "## ") {
			collectMetadata = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(line, "## ")), canonicalMetadataHeading)
			continue
		}

		if collectMetadata {
			trimmedBullet := strings.TrimSpace(strings.TrimLeft(line, "-*)) "))
			key, value, ok := strings.Cut(trimmedBullet, ":")
			if ok {
				metadata[normalizeCanonicalMetadataKey(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
				continue
			}
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
		return CanonicalProjectMetadata{}, fmt.Errorf("read canonical markdown %s: %w", fallbackSlug, err)
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
	base.ClientName = firstNonEmpty(frontmatter["client_name"], frontmatter["client"], frontmatter["brand"], metadata["client_context"])
	base.IndustryType = model.NormalizeIndustryTypeInput(firstNonEmpty(frontmatter["industry_type"], metadata["industry_type"]))
	base.FinalProduct = model.NormalizeFinalProduct(firstNonEmpty(frontmatter["final_product"], metadata["final_product"]))

	if strings.TrimSpace(base.Name) == "" {
		base.Name = humanizeSlug(filepath.Base(fallbackSlug))
	}
	if strings.TrimSpace(base.Description) == "" {
		base.Description = fmt.Sprintf("Imported from canonical case-study source %s.", fallbackSlug)
	}
	if err := model.ValidateIndustryType(base.IndustryType); err != nil {
		return CanonicalProjectMetadata{}, fmt.Errorf("canonical markdown %s: %w", fallbackSlug, err)
	}
	if err := model.ValidateFinalProduct(base.FinalProduct); err != nil {
		return CanonicalProjectMetadata{}, fmt.Errorf("canonical markdown %s: %w", fallbackSlug, err)
	}

	return base, nil
}

func bytesReader(body []byte) io.Reader {
	return strings.NewReader(string(body))
}

func normalizeCanonicalMetadataKey(key string) string {
	normalized := strings.ToLower(strings.TrimSpace(key))
	normalized = strings.ReplaceAll(normalized, "/", " ")
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = strings.Join(strings.Fields(normalized), "_")
	return normalized
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
