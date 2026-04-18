package workflow

import (
	"context"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type Repository interface {
	SaveRun(ctx context.Context, run model.CaseStudyWorkflowRun) error
	GetRun(ctx context.Context, runID uuid.UUID) (model.CaseStudyWorkflowRun, error)
	ListLogs(ctx context.Context, runID uuid.UUID) ([]model.CaseStudyWorkflowLogEntry, error)
	AppendLog(ctx context.Context, entry model.CaseStudyWorkflowLogEntry) error
}
