package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCanonicalSlugDirFromCaseDirectory(t *testing.T) {
	root := t.TempDir()
	slugDir := filepath.Join(root, "90. dev_portfolioforge", "can-bus-crane-monitoring")
	if err := os.MkdirAll(slugDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(slugDir, "can-bus-crane-monitoring.md"), []byte("# demo"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	resolvedDir, slug, err := resolveCanonicalSlugDir(root, "")
	if err != nil {
		t.Fatalf("resolveCanonicalSlugDir: %v", err)
	}
	if slug != "can-bus-crane-monitoring" {
		t.Fatalf("slug = %q", slug)
	}
	if resolvedDir != slugDir {
		t.Fatalf("resolvedDir = %q, want %q", resolvedDir, slugDir)
	}
}

func TestResolveCanonicalSlugDirWithExplicitSlug(t *testing.T) {
	root := t.TempDir()
	canonicalRoot := filepath.Join(root, "90. dev_portfolioforge")
	slugDir := filepath.Join(canonicalRoot, "printer-05-controls-migration")
	if err := os.MkdirAll(slugDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(slugDir, "printer-05-controls-migration.md"), []byte("# demo"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	resolvedDir, slug, err := resolveCanonicalSlugDir(canonicalRoot, "printer-05-controls-migration")
	if err != nil {
		t.Fatalf("resolveCanonicalSlugDir: %v", err)
	}
	if slug != "printer-05-controls-migration" {
		t.Fatalf("slug = %q", slug)
	}
	if resolvedDir != slugDir {
		t.Fatalf("resolvedDir = %q, want %q", resolvedDir, slugDir)
	}
}

func TestResolveCanonicalPublishTargetBuildsPublicURL(t *testing.T) {
	root := t.TempDir()
	slugDir := filepath.Join(root, "90. dev_portfolioforge", "tableros-y-pupitres-linea-bradbury")
	if err := os.MkdirAll(slugDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(slugDir, "tableros-y-pupitres-linea-bradbury.md"), []byte("# demo"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	target, err := resolveCanonicalPublishTarget(root, "", canonicalPublishConfig{
		PublicBase: "https://mlbautomation.com/dev/portfolioforge",
		RemoteBase: "/",
	})
	if err != nil {
		t.Fatalf("resolveCanonicalPublishTarget: %v", err)
	}

	wantURL := "https://mlbautomation.com/dev/portfolioforge/tableros-y-pupitres-linea-bradbury/tableros-y-pupitres-linea-bradbury.md"
	if target.PublicURL != wantURL {
		t.Fatalf("PublicURL = %q, want %q", target.PublicURL, wantURL)
	}
	if target.RemoteFile != "/tableros-y-pupitres-linea-bradbury/tableros-y-pupitres-linea-bradbury.md" {
		t.Fatalf("RemoteFile = %q", target.RemoteFile)
	}
}
