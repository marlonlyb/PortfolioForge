package localization

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	projectPorts "github.com/marlonlyb/portfolioforge/domain/ports/project"
)

type BackfillService struct {
	projectReader projectPorts.ProjectReader
	service       *Service
}

func NewBackfillService(projectReader projectPorts.ProjectReader, service *Service) *BackfillService {
	return &BackfillService{projectReader: projectReader, service: service}
}

func (s *BackfillService) BackfillProject(ctx context.Context, projectID uuid.UUID, locales []string) error {
	if s == nil || s.projectReader == nil || s.service == nil {
		return nil
	}

	project, err := s.projectReader.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("load project %s: %w", projectID, err)
	}

	if err := s.service.RegenerateFromSpanish(ctx, projectID, BuildProjectFieldMap(project), locales); err != nil {
		return fmt.Errorf("regenerate localizations for %s: %w", projectID, err)
	}

	return nil
}
