package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestProjectMediaVariantFallbackOrder(t *testing.T) {
	media := ProjectMedia{
		LowURL:      "https://cdn.example.com/project-low.webp",
		MediumURL:   "https://cdn.example.com/project-medium.webp",
		HighURL:     "https://cdn.example.com/project-high.webp",
		FallbackURL: "https://cdn.example.com/project-original.jpg",
	}

	if got := media.ThumbnailSrc(); got != media.LowURL {
		t.Fatalf("ThumbnailSrc() = %q, want %q", got, media.LowURL)
	}
	if got := media.MediumSrc(); got != media.MediumURL {
		t.Fatalf("MediumSrc() = %q, want %q", got, media.MediumURL)
	}
	if got := media.FullSrc(); got != media.HighURL {
		t.Fatalf("FullSrc() = %q, want %q", got, media.HighURL)
	}
}

func TestBuildLegacyProjectMediaUsesCanonicalVariantKeys(t *testing.T) {
	projectID := uuid.New()
	media := BuildLegacyProjectMedia(projectID, []string{"https://cdn.example.com/project.webp"})
	if len(media) != 1 {
		t.Fatalf("len(media) = %d, want 1", len(media))
	}
	item := media[0]
	if item.LowURL == "" || item.MediumURL == "" || item.HighURL == "" || item.FallbackURL == "" {
		t.Fatalf("legacy media variants were not fully populated: %#v", item)
	}
}

func TestBuildProjectImageListPrefersMediumThenHighThenLowThenFallback(t *testing.T) {
	images := BuildProjectImageList([]ProjectMedia{{
		LowURL:      "https://cdn.example.com/project-low.webp",
		HighURL:     "https://cdn.example.com/project-high.webp",
		FallbackURL: "https://cdn.example.com/project-original.jpg",
	}}, nil)

	want := []string{"https://cdn.example.com/project-high.webp"}
	if !reflect.DeepEqual(images, want) {
		t.Fatalf("BuildProjectImageList() = %#v, want %#v", images, want)
	}
}
