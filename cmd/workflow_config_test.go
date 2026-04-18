package main

import (
	"path/filepath"
	"testing"
)

func TestLoadCaseStudyWorkflowEnvConfig_DisablesWorkflowWithoutRequiredEnv(t *testing.T) {
	config, err := loadCaseStudyWorkflowEnvConfigFrom(func(string) string { return "" }, false)
	if err != nil {
		t.Fatalf("loadCaseStudyWorkflowEnvConfigFrom() error = %v", err)
	}
	if config.Configured {
		t.Fatal("expected configured=false when env is missing")
	}
	if config.Reason == "" || config.Diagnostic == "" {
		t.Fatalf("expected unavailable metadata, got %#v", config)
	}
}

func TestLoadRequiredCaseStudyWorkflowEnvConfig_FailsWithoutRequiredEnv(t *testing.T) {
	_, err := loadCaseStudyWorkflowEnvConfigFrom(func(string) string { return "" }, true)
	if err == nil {
		t.Fatal("expected strict workflow loader to fail")
	}
}

func TestLoadCaseStudyWorkflowEnvConfig_ReturnsConfiguredWhenEnvIsValid(t *testing.T) {
	root := t.TempDir()
	getenv := func(key string) string {
		switch key {
		case "PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS":
			return root
		case "PF_FTP_HOST":
			return "ftp.example.com"
		case "PF_FTP_PORT":
			return "21"
		case "PF_FTP_USER":
			return "demo"
		case "PF_FTP_PASSWORD":
			return "secret"
		case "PF_PUBLIC_BASE":
			return "https://example.com/dev/portfolioforge"
		case "PF_FTP_REMOTE_BASE":
			return "/remote"
		default:
			return ""
		}
	}

	config, err := loadCaseStudyWorkflowEnvConfigFrom(getenv, false)
	if err != nil {
		t.Fatalf("loadCaseStudyWorkflowEnvConfigFrom() error = %v", err)
	}
	if !config.Configured {
		t.Fatalf("expected configured workflow, got %#v", config)
	}
	if len(config.AllowedSourceRoots) != 1 || config.AllowedSourceRoots[0] != filepath.Clean(root) {
		t.Fatalf("allowed roots = %#v", config.AllowedSourceRoots)
	}
}
