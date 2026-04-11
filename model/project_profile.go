package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ProjectProfile struct {
	ProjectID          uuid.UUID       `json:"project_id"`
	BusinessGoal       string          `json:"business_goal,omitempty"`
	ProblemStatement   string          `json:"problem_statement,omitempty"`
	SolutionSummary    string          `json:"solution_summary,omitempty"`
	Architecture       string          `json:"architecture,omitempty"`
	Integrations       json.RawMessage `json:"integrations"`
	AIUsage            string          `json:"ai_usage,omitempty"`
	TechnicalDecisions json.RawMessage `json:"technical_decisions"`
	Challenges         json.RawMessage `json:"challenges"`
	Results            json.RawMessage `json:"results"`
	Metrics            json.RawMessage `json:"metrics"`
	Timeline           json.RawMessage `json:"timeline"`
	UpdatedAt          int64           `json:"updated_at"`
}
