package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type UnavailableCaseStudyWorkflowService struct {
	reason string
}

func NewUnavailableCaseStudyWorkflowService(reason string) *UnavailableCaseStudyWorkflowService {
	return &UnavailableCaseStudyWorkflowService{reason: reason}
}

func (s *UnavailableCaseStudyWorkflowService) GetAvailability(context.Context) (model.CaseStudyWorkflowAvailability, error) {
	return model.CaseStudyWorkflowAvailability{
		Configured: false,
		Reason:     s.reason,
	}, nil
}

func (s *UnavailableCaseStudyWorkflowService) StartRun(context.Context, model.StartCaseStudyWorkflowRunRequest) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) GetRun(context.Context, uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) ListLogs(context.Context, uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error) {
	return nil, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) ConfirmStep(context.Context, uuid.UUID, string) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) StartStep(context.Context, uuid.UUID, string) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) RetryStep(context.Context, uuid.UUID, string) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) Resume(context.Context, uuid.UUID) (model.CaseStudyWorkflowRun, error) {
	return model.CaseStudyWorkflowRun{}, s.unavailable()
}

func (s *UnavailableCaseStudyWorkflowService) unavailable() error {
	return &model.CaseStudyWorkflowUnavailableError{Reason: s.reason}
}
