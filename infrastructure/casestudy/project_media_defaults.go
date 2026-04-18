package casestudy

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

const (
	defaultProjectMediaCount = 7
	projectMediaBaseURL      = "https://mlbautomation.com/dev/portfolioforge"
	projectMediaFallbackURL  = projectMediaBaseURL + "/imagen_fallback/Logo_500_500.png"
)

func ensureMinimumProjectMedia(items []model.ProjectMedia, slug string, minimum int) []model.ProjectMedia {
	if minimum <= 0 {
		minimum = defaultProjectMediaCount
	}

	trimmedSlug := strings.TrimSpace(slug)
	result := make([]model.ProjectMedia, 0, maxInt(len(items), minimum))
	for _, item := range items {
		result = append(result, item)
	}

	for len(result) < minimum {
		index := len(result)
		result = append(result, buildDefaultProjectMedia(trimmedSlug, index))
	}

	featuredFound := false
	for index := range result {
		result[index].SortOrder = index
		if result[index].Featured {
			if featuredFound {
				result[index].Featured = false
				continue
			}
			featuredFound = true
		}
	}
	if len(result) > 0 && !featuredFound {
		result[0].Featured = true
	}

	return result
}

func assignProjectIDToMedia(items []model.ProjectMedia, projectID uuid.UUID) []model.ProjectMedia {
	assigned := make([]model.ProjectMedia, 0, len(items))
	for _, item := range items {
		item.ProjectID = projectID
		assigned = append(assigned, item)
	}
	return assigned
}

func buildDefaultProjectMedia(slug string, index int) model.ProjectMedia {
	imageNumber := index + 1
	base := strings.TrimSpace(slug)
	if base != "" {
		base = fmt.Sprintf("%s/%s/imagen%02d", projectMediaBaseURL, base, imageNumber)
	}

	item := model.ProjectMedia{
		MediaType:   "image",
		FallbackURL: projectMediaFallbackURL,
		SortOrder:   index,
		Featured:    index == 0,
	}

	if base != "" {
		item.LowURL = base + "_low.webp"
		item.MediumURL = base + "_medium.webp"
		item.HighURL = base + "_high.webp"
	}

	return item
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
