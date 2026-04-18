package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marlonlyb/portfolioforge/infrastructure/casestudy"
)

type CaseStudyWorkflowEnvConfig struct {
	Configured         bool
	AllowedSourceRoots []string
	Publish            casestudy.PublishConfig
	Reason             string
	Diagnostic         string
}

const caseStudyWorkflowUnavailableReason = "Case-study workflow is not configured in this environment."

func LoadCaseStudyWorkflowEnvConfig() CaseStudyWorkflowEnvConfig {
	config, _ := loadCaseStudyWorkflowEnvConfigFrom(os.Getenv, false)
	return config
}

func LoadRequiredCaseStudyWorkflowEnvConfig() (CaseStudyWorkflowEnvConfig, error) {
	return loadCaseStudyWorkflowEnvConfigFrom(os.Getenv, true)
}

func loadCaseStudyWorkflowEnvConfigFrom(getenv func(string) string, strict bool) (CaseStudyWorkflowEnvConfig, error) {
	publishConfig, err := casestudy.LoadPublishConfigFromEnv(getenv)
	if err != nil {
		return unavailableCaseStudyWorkflowEnvConfig(err, strict)
	}

	roots, err := loadCaseStudyWorkflowAllowedRoots(getenv("PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS"))
	if err != nil {
		return unavailableCaseStudyWorkflowEnvConfig(err, strict)
	}

	return CaseStudyWorkflowEnvConfig{
		Configured:         true,
		AllowedSourceRoots: roots,
		Publish:            publishConfig,
	}, nil
}

func unavailableCaseStudyWorkflowEnvConfig(err error, strict bool) (CaseStudyWorkflowEnvConfig, error) {
	if strict {
		return CaseStudyWorkflowEnvConfig{}, err
	}

	return CaseStudyWorkflowEnvConfig{
		Configured: false,
		Reason:     caseStudyWorkflowUnavailableReason,
		Diagnostic: err.Error(),
	}, nil
}

func loadCaseStudyWorkflowAllowedRoots(raw string) ([]string, error) {
	rootsRaw := strings.TrimSpace(raw)
	if rootsRaw == "" {
		return nil, fmt.Errorf("PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS es obligatoria para el workflow de case studies")
	}

	parts := strings.Split(rootsRaw, ",")
	roots := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		abs, err := filepath.Abs(filepath.Clean(trimmed))
		if err != nil {
			return nil, fmt.Errorf("invalid allowed source root %q: %w", trimmed, err)
		}
		roots = append(roots, abs)
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS no contiene rutas válidas")
	}

	return roots, nil
}
