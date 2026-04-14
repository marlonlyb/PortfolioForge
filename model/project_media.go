package model

import (
	"strings"

	"github.com/google/uuid"
)

type ProjectMedia struct {
	ID           uuid.UUID `json:"id"`
	ProjectID    uuid.UUID `json:"project_id"`
	MediaType    string    `json:"media_type"` // image|video|diagram|document
	URL          string    `json:"url,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	MediumURL    string    `json:"medium_url,omitempty"`
	FullURL      string    `json:"full_url,omitempty"`
	Caption      string    `json:"caption,omitempty"`
	AltText      string    `json:"alt_text,omitempty"`
	SortOrder    int       `json:"sort_order"`
	Featured     bool      `json:"featured"`
}

func (m ProjectMedia) ThumbnailSrc() string {
	return firstNonEmpty(m.ThumbnailURL, m.MediumURL, m.FullURL, m.URL)
}

func (m ProjectMedia) MediumSrc() string {
	return firstNonEmpty(m.MediumURL, m.FullURL, m.ThumbnailURL, m.URL)
}

func (m ProjectMedia) FullSrc() string {
	return firstNonEmpty(m.FullURL, m.MediumURL, m.ThumbnailURL, m.URL)
}

func BuildProjectImageList(media []ProjectMedia, fallback []string) []string {
	if len(media) == 0 {
		return fallback
	}

	seen := make(map[string]struct{}, len(media))
	images := make([]string, 0, len(media))

	for _, item := range media {
		src := item.MediumSrc()
		if src == "" {
			src = item.FullSrc()
		}
		if src == "" {
			continue
		}
		if _, exists := seen[src]; exists {
			continue
		}
		seen[src] = struct{}{}
		images = append(images, src)
	}

	if len(images) == 0 {
		return fallback
	}

	return images
}

func BuildLegacyProjectMedia(projectID uuid.UUID, images []string) []ProjectMedia {
	media := make([]ProjectMedia, 0, len(images))
	for index, image := range images {
		trimmed := strings.TrimSpace(image)
		if trimmed == "" {
			continue
		}

		media = append(media, ProjectMedia{
			ProjectID:    projectID,
			MediaType:    "image",
			URL:          trimmed,
			ThumbnailURL: trimmed,
			MediumURL:    trimmed,
			FullURL:      trimmed,
			SortOrder:    index,
			Featured:     index == 0,
		})
	}

	return media
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
