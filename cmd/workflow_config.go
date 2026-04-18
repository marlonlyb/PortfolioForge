package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marlonlyb/portfolioforge/infrastructure/casestudy"
)

type CaseStudyWorkflowEnvConfig struct {
	AllowedSourceRoots []string
	Publish            casestudy.PublishConfig
}

func LoadCaseStudyWorkflowEnvConfig() (CaseStudyWorkflowEnvConfig, error) {
	publishConfig, err := casestudy.LoadPublishConfigFromEnv(os.Getenv)
	if err != nil {
		return CaseStudyWorkflowEnvConfig{}, err
	}

	rootsRaw := strings.TrimSpace(os.Getenv("PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS"))
	if rootsRaw == "" {
		return CaseStudyWorkflowEnvConfig{}, fmt.Errorf("PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS es obligatoria para el workflow de case studies")
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
			return CaseStudyWorkflowEnvConfig{}, fmt.Errorf("invalid allowed source root %q: %w", trimmed, err)
		}
		roots = append(roots, abs)
	}
	if len(roots) == 0 {
		return CaseStudyWorkflowEnvConfig{}, fmt.Errorf("PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS no contiene rutas válidas")
	}

	return CaseStudyWorkflowEnvConfig{AllowedSourceRoots: roots, Publish: publishConfig}, nil
}
