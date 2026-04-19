package model

import "github.com/google/uuid"

// SearchParams holds the input parameters for a hybrid retrieval search.
type SearchParams struct {
	Query          string
	QueryEmbedding []float32
	Category       string
	Client         string
	Technologies   []string
	Cursor         string
	PageSize       int
}

// SearchResult represents a single project match with relevance metadata.
type SearchResult struct {
	Project       Project
	Score         float64
	LexicalScore  float64
	FuzzyScore    float64
	SemanticScore float64
	Explanation   string
	Evidence      EvidenceTrace
}

// EvidenceTrace captures which fields and terms contributed to a search match.
type EvidenceTrace struct {
	ProjectID uuid.UUID
	Fields    []EvidenceField
}

// MatchType constants define the possible match types for an evidence field.
const (
	MatchTypeFTS        = "fts"
	MatchTypeFuzzy      = "fuzzy"
	MatchTypeSemantic   = "semantic"
	MatchTypeStructured = "structured"
)

// MatchType is the union of valid match type strings.
type MatchType = string

// EvidenceField describes a single field's contribution to the evidence trace.
type EvidenceField struct {
	Field       string    `json:"field"`
	MatchedText string    `json:"matched_text"`
	MatchType   MatchType `json:"match_type"`
	Score       float64   `json:"score"`
}

// SearchResponse is the API-level response envelope for search results.
type SearchResponse struct {
	Data []SearchResultItem `json:"data"`
	Meta SearchMeta         `json:"meta"`
}

// SearchResultItem is a single project in the search API response.
type SearchResultItem struct {
	ID           string              `json:"id"`
	Slug         string              `json:"slug"`
	Title        string              `json:"title"`
	Category     string              `json:"category"`
	ClientName   *string             `json:"client_name"`
	IndustryType *string             `json:"industry_type,omitempty"`
	FinalProduct *string             `json:"final_product,omitempty"`
	Summary      *string             `json:"summary"`
	Technologies []TechnologySummary `json:"technologies"`
	HeroImage    *string             `json:"hero_image"`
	Score        float64             `json:"score"`
	Explanation  *string             `json:"explanation"`
	Evidence     []EvidenceField     `json:"evidence"`
}

// TechnologySummary is a lightweight technology representation for search results.
type TechnologySummary struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Color *string `json:"color"`
}

// SearchMeta contains pagination and filter metadata for the search response.
type SearchMeta struct {
	Total          int            `json:"total"`
	PageSize       int            `json:"page_size"`
	Cursor         *string        `json:"cursor"`
	Query          string         `json:"query"`
	FiltersApplied FiltersApplied `json:"filters_applied"`
}

// FiltersApplied records which filters were active during the search.
type FiltersApplied struct {
	Category     *string  `json:"category"`
	Client       *string  `json:"client"`
	Technologies []string `json:"technologies"`
}
