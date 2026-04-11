package model

import "github.com/google/uuid"

// ReadinessLevel constants define the three search readiness tiers.
const (
	ReadinessIncomplete = "incomplete"
	ReadinessBasic      = "basic"
	ReadinessComplete   = "complete"

	// MinDescriptionLen is the minimum description length for meaningful FTS.
	MinDescriptionLen = 20
)

// ProjectReadiness describes how search-ready a project is.
type ProjectReadiness struct {
	ProjectID          uuid.UUID `json:"project_id"`
	Level              string    `json:"level"` // incomplete, basic, complete
	MissingFields      []string  `json:"missing_fields,omitempty"`
	HasName            bool      `json:"has_name"`
	HasDescription     bool      `json:"has_description"`
	HasCategory        bool      `json:"has_category"`
	HasTechnologies    bool      `json:"has_technologies"`
	HasSolutionSummary bool      `json:"has_solution_summary"`
}

// ComputeReadiness evaluates a project's search readiness based on its
// data completeness. A project is "search-ready" when it has:
//   - name (required)
//   - description (required, min 20 chars for meaningful FTS)
//   - category (required for structured filter)
//   - at least 1 technology linked (boosts fuzzy + semantic)
//   - solution_summary in profile (highest-weight FTS field)
//
// Readiness levels:
//   - Incomplete: missing name or description
//   - Basic: has name + description but no profile/technologies
//   - Complete: has name + description + profile + technologies
func ComputeReadiness(p Project, techs []Technology) ProjectReadiness {
	r := ProjectReadiness{
		ProjectID: p.ID,
	}

	// Evaluate individual fields
	r.HasName = p.Name != ""
	r.HasDescription = len(p.Description) >= MinDescriptionLen
	r.HasCategory = p.Category != ""
	r.HasTechnologies = len(techs) > 0

	r.HasSolutionSummary = false
	if p.Profile != nil && p.Profile.SolutionSummary != "" {
		r.HasSolutionSummary = true
	}

	// Collect missing fields
	var missing []string
	if !r.HasName {
		missing = append(missing, "name")
	}
	if !r.HasDescription {
		missing = append(missing, "description")
	}
	if !r.HasCategory {
		missing = append(missing, "category")
	}
	if !r.HasTechnologies {
		missing = append(missing, "technologies")
	}
	if !r.HasSolutionSummary {
		missing = append(missing, "solution_summary")
	}
	r.MissingFields = missing

	// Determine readiness level
	if !r.HasName || !r.HasDescription {
		r.Level = ReadinessIncomplete
	} else if !r.HasSolutionSummary && !r.HasTechnologies {
		r.Level = ReadinessBasic
	} else {
		r.Level = ReadinessComplete
	}

	return r
}
