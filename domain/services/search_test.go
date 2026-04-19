package services

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	"github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/model"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockSearchRepo struct {
	results []model.SearchResult
	err     error
	params  []model.SearchParams
}

func (m *mockSearchRepo) Search(_ context.Context, params model.SearchParams) ([]model.SearchResult, error) {
	m.params = append(m.params, params)
	return m.results, m.err
}
func (m *mockSearchRepo) RefreshSearchDocument(_ context.Context, _ uuid.UUID) error {
	return nil
}
func (m *mockSearchRepo) RefreshAllDocuments(_ context.Context) error { return nil }

type mockProjectReader struct {
	techsByProject map[uuid.UUID][]model.Technology
}

func (m *mockProjectReader) GetByID(_ context.Context, _ uuid.UUID) (model.Project, error) {
	return model.Project{}, nil
}
func (m *mockProjectReader) GetBySlug(_ context.Context, _ string) (model.Project, error) {
	return model.Project{}, nil
}
func (m *mockProjectReader) ListPublished(_ context.Context) ([]model.Project, error) {
	return nil, nil
}
func (m *mockProjectReader) GetTechnologiesByProjectID(_ context.Context, projectID uuid.UUID) ([]model.Technology, error) {
	return m.techsByProject[projectID], nil
}

func (m *mockProjectReader) GetAssistantContextBySlug(_ context.Context, _ string, _ string) (model.ProjectAssistantContext, error) {
	return model.ProjectAssistantContext{}, nil
}

type mockEmbeddingProv struct {
	vectors []float32
	err     error
	inputs  []string
}

func (m *mockEmbeddingProv) Generate(_ context.Context, input string) ([]float32, error) {
	m.inputs = append(m.inputs, input)
	return m.vectors, m.err
}
func (m *mockEmbeddingProv) Dimension() int { return 1536 }

// Verify interfaces at compile time.
var (
	_ search.SearchRepository     = (*mockSearchRepo)(nil)
	_ project.ProjectReader       = (*mockProjectReader)(nil)
	_ embedding.EmbeddingProvider = (*mockEmbeddingProv)(nil)
)

// ---------------------------------------------------------------------------
// Helper: build a minimal Search service with mocked deps.
// ---------------------------------------------------------------------------

func newTestSearch(semanticEnabled bool, repoResults []model.SearchResult) *Search {
	return NewSearch(
		&mockSearchRepo{results: repoResults},
		&mockProjectReader{techsByProject: map[uuid.UUID][]model.Technology{}},
		&mockEmbeddingProv{},
		NewNoOpExplanationProvider(),
		semanticEnabled,
	)
}

// NewNoOpExplanationProvider is a local helper that returns empty explanations.
// We use a simple inline stub instead of importing infrastructure.
type noOpExplanationProv struct{}

func (n *noOpExplanationProv) Explain(_ context.Context, _ model.Project, _ model.EvidenceTrace, _ string) (string, error) {
	return "", nil
}

func NewNoOpExplanationProvider() search.ExplanationProvider {
	return &noOpExplanationProv{}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestFuseAndBoost_ScoreFusionWithSemantic(t *testing.T) {
	// With semantic enabled: LexicalWeight=0.45, FuzzyWeight=0.25, SemanticWeight=0.30
	// Given a single result where LexicalScore=0.8, FuzzyScore=0.6, SemanticScore=0.85,
	// maxLexical=0.8 so lexicalNorm=0.8/0.8=1.0, fuzzyNorm=0.6, semanticNorm=0.85
	// fusedScore = 0.45*1.0 + 0.25*0.6 + 0.30*0.85 = 0.45 + 0.15 + 0.255 = 0.855
	svc := newTestSearch(true, nil)

	projectID := uuid.New()
	results := []model.SearchResult{
		{
			Project:       model.Project{ID: projectID, Name: "Test", Category: "web"},
			LexicalScore:  0.8,
			FuzzyScore:    0.6,
			SemanticScore: 0.85,
		},
	}

	got := svc.fuseAndBoost(results, "react")

	wantScore := 0.45*1.0 + 0.25*0.6 + 0.30*0.85 // 0.855
	if diff := got[0].Score - wantScore; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("fused score = %v, want %v", got[0].Score, wantScore)
	}
}

func TestFuseAndBoost_NoSemantic(t *testing.T) {
	// Without semantic: LexicalWeight=0.60, FuzzyWeight=0.40, SemanticWeight=0
	// maxLexical=0.5 so lexicalNorm=0.5/0.5=1.0, fuzzyNorm=0.4, semanticNorm=0
	// fusedScore = 0.60*1.0 + 0.40*0.4 + 0*0 = 0.60 + 0.16 = 0.76
	svc := newTestSearch(false, nil)

	projectID := uuid.New()
	results := []model.SearchResult{
		{
			Project:       model.Project{ID: projectID, Name: "Test", Category: "web"},
			LexicalScore:  0.5,
			FuzzyScore:    0.4,
			SemanticScore: 0.9, // should be ignored
		},
	}

	got := svc.fuseAndBoost(results, "react")

	wantScore := 0.60*1.0 + 0.40*0.4 // 0.76
	if diff := got[0].Score - wantScore; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("fused score = %v, want %v", got[0].Score, wantScore)
	}
}

func TestFuseAndBoost_CategoryBoost(t *testing.T) {
	// Category match adds CategoryBoost=0.15
	svc := newTestSearch(true, nil)

	projectID := uuid.New()
	results := []model.SearchResult{
		{
			Project:       model.Project{ID: projectID, Name: "Test", Category: "microservices"},
			LexicalScore:  0.8,
			FuzzyScore:    0.6,
			SemanticScore: 0.85,
		},
	}

	got := svc.fuseAndBoost(results, "microservices")

	baseScore := 0.45*1.0 + 0.25*0.6 + 0.30*0.85 // 0.855
	wantScore := baseScore + 0.15                // category boost
	if diff := got[0].Score - wantScore; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("fused score with category boost = %v, want %v", got[0].Score, wantScore)
	}
}

func TestFuseAndBoost_EmptyResults(t *testing.T) {
	svc := newTestSearch(true, nil)
	got := svc.fuseAndBoost([]model.SearchResult{}, "react")
	if len(got) != 0 {
		t.Errorf("expected empty results, got %d", len(got))
	}
}

func TestApplyThreshold_ExcludesLowScores(t *testing.T) {
	svc := newTestSearch(true, nil)

	// Threshold is 0.10
	results := []model.SearchResult{
		{Project: model.Project{ID: uuid.New()}, Score: 0.50},
		{Project: model.Project{ID: uuid.New()}, Score: 0.09},  // below threshold
		{Project: model.Project{ID: uuid.New()}, Score: 0.10},  // exactly threshold
		{Project: model.Project{ID: uuid.New()}, Score: 0.001}, // well below
	}

	got := svc.applyThreshold(results)

	if len(got) != 2 {
		t.Fatalf("expected 2 results above threshold, got %d", len(got))
	}
	if got[0].Score != 0.50 {
		t.Errorf("first result score = %v, want 0.50", got[0].Score)
	}
	if got[1].Score != 0.10 {
		t.Errorf("second result score = %v, want 0.10", got[1].Score)
	}
}

func TestFindMatchingTerms(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		terms []string
		want  []string
	}{
		{
			name:  "single match",
			text:  "microservices architecture",
			terms: []string{"microservices"},
			want:  []string{"microservices"},
		},
		{
			name:  "multiple matches",
			text:  "react and node microservices",
			terms: []string{"react", "node", "python"},
			want:  []string{"react", "node"},
		},
		{
			name:  "no matches",
			text:  "hello world",
			terms: []string{"react"},
			want:  nil,
		},
		{
			name:  "case insensitive",
			text:  "REACT App",
			terms: []string{"react"},
			want:  []string{"react"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMatchingTerms(tt.text, tt.terms)
			if len(got) != len(tt.want) {
				t.Fatalf("findMatchingTerms() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("term[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExtractSnippet(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		term   string
		maxLen int
		want   string
	}{
		{
			name:   "term found at start",
			text:   "react application with microservices",
			term:   "react",
			maxLen: 20,
			want:   "react application wi...",
		},
		{
			name:   "term not found",
			text:   "hello world",
			term:   "react",
			maxLen: 10,
			want:   "hello worl...",
		},
		{
			name:   "short text no truncation",
			text:   "react",
			term:   "react",
			maxLen: 80,
			want:   "react",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSnippet(tt.text, tt.term, tt.maxLen)
			if got != tt.want {
				t.Errorf("extractSnippet() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractEvidence(t *testing.T) {
	svc := newTestSearch(true, nil)
	projectID := uuid.New()

	tests := []struct {
		name   string
		result model.SearchResult
		query  string
		techs  []model.Technology
		want   []string // expected field names (order may vary)
	}{
		{
			name: "match in name and category",
			result: model.SearchResult{
				Project: model.Project{
					ID:          projectID,
					Name:        "React Dashboard",
					Category:    "web",
					Description: "A dashboard app",
				},
			},
			query: "react web",
			techs: nil,
			want:  []string{"name", "category"},
		},
		{
			name: "match in technology",
			result: model.SearchResult{
				Project: model.Project{
					ID:          projectID,
					Name:        "Dashboard",
					Category:    "web",
					Description: "A dashboard app",
				},
			},
			query: "react",
			techs: []model.Technology{
				{ID: uuid.New(), Name: "React", Slug: "react"},
			},
			want: []string{"technology"},
		},
		{
			name: "empty query returns empty evidence",
			result: model.SearchResult{
				Project: model.Project{ID: projectID},
			},
			query: "",
			techs: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.extractEvidence(tt.result, tt.query, tt.techs)

			if tt.query == "" {
				if len(got.Fields) != 0 {
					t.Errorf("expected no fields for empty query, got %d", len(got.Fields))
				}
				return
			}

			gotFields := map[string]bool{}
			for _, f := range got.Fields {
				gotFields[f.Field] = true
			}
			for _, w := range tt.want {
				if !gotFields[w] {
					t.Errorf("expected field %q in evidence, not found; got fields: %v", w, gotFields)
				}
			}
		})
	}
}

func TestSearchGeneratesSemanticEmbeddingAndPassesItToRepository(t *testing.T) {
	projectID := uuid.New()
	searchRepo := &mockSearchRepo{results: []model.SearchResult{{
		Project:      model.Project{ID: projectID, Name: "PortfolioForge", Category: "platform"},
		LexicalScore: 1,
	}}}
	projectReader := &mockProjectReader{techsByProject: map[uuid.UUID][]model.Technology{}}
	embeddingProvider := &mockEmbeddingProv{vectors: []float32{0.11, 0.22, 0.33}}
	svc := NewSearch(searchRepo, projectReader, embeddingProvider, NewNoOpExplanationProvider(), true)

	_, err := svc.Search(context.Background(), model.SearchParams{Query: "  PortfolioForge  ", PageSize: 10})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(embeddingProvider.inputs) != 1 || embeddingProvider.inputs[0] != "portfolioforge" {
		t.Fatalf("embedding inputs = %#v", embeddingProvider.inputs)
	}
	if len(searchRepo.params) != 1 {
		t.Fatalf("repo calls = %d, want 1", len(searchRepo.params))
	}
	if len(searchRepo.params[0].QueryEmbedding) != 3 {
		t.Fatalf("query embedding length = %d, want 3", len(searchRepo.params[0].QueryEmbedding))
	}
}

func TestSearchPaginatesWithoutDuplicatesAndBuildsCompleteResultItems(t *testing.T) {
	projectID1 := uuid.New()
	projectID2 := uuid.New()
	projectID3 := uuid.New()
	searchRepo := &mockSearchRepo{results: []model.SearchResult{
		{Project: model.Project{ID: projectID1, Slug: "alpha", Name: "Alpha", Category: "platform", ClientName: "Acme", Profile: &model.ProjectProfile{ProjectID: projectID1, SolutionSummary: "Summary alpha"}, Images: []byte(`["https://img/alpha.png"]`)}, LexicalScore: 0.9},
		{Project: model.Project{ID: projectID2, Slug: "beta", Name: "Beta", Category: "platform", ClientName: "Acme", Profile: &model.ProjectProfile{ProjectID: projectID2, SolutionSummary: "Summary beta"}, Images: []byte(`["https://img/beta.png"]`)}, LexicalScore: 0.8},
		{Project: model.Project{ID: projectID3, Slug: "gamma", Name: "Gamma", Category: "platform", ClientName: "Acme", Profile: &model.ProjectProfile{ProjectID: projectID3, SolutionSummary: "Summary gamma"}, Images: []byte(`["https://img/gamma.png"]`)}, LexicalScore: 0.7},
	}}
	projectReader := &mockProjectReader{techsByProject: map[uuid.UUID][]model.Technology{
		projectID1: {{ID: uuid.New(), Name: "React", Slug: "react", Color: "#61dafb"}},
		projectID2: {{ID: uuid.New(), Name: "Go", Slug: "go", Color: "#00add8"}},
		projectID3: {{ID: uuid.New(), Name: "PostgreSQL", Slug: "postgresql", Color: "#336791"}},
	}}
	svc := NewSearch(searchRepo, projectReader, &mockEmbeddingProv{}, NewNoOpExplanationProvider(), false)

	firstPage, err := svc.Search(context.Background(), model.SearchParams{Query: "platform", PageSize: 2})
	if err != nil {
		t.Fatalf("Search() first page error = %v", err)
	}
	if len(firstPage.Data) != 2 {
		t.Fatalf("first page items = %d, want 2", len(firstPage.Data))
	}
	if firstPage.Meta.Total != 3 {
		t.Fatalf("first page total = %d, want 3", firstPage.Meta.Total)
	}
	if firstPage.Meta.Cursor == nil || *firstPage.Meta.Cursor != "2" {
		t.Fatalf("first page cursor = %#v, want 2", firstPage.Meta.Cursor)
	}
	if firstPage.Data[0].Slug != "alpha" || firstPage.Data[0].HeroImage == nil || *firstPage.Data[0].HeroImage != "https://img/alpha.png" {
		t.Fatalf("first page item missing slug/hero image: %#v", firstPage.Data[0])
	}
	if len(firstPage.Data[0].Technologies) != 1 || firstPage.Data[0].Technologies[0].Slug != "react" {
		t.Fatalf("first page technologies = %#v", firstPage.Data[0].Technologies)
	}

	secondPage, err := svc.Search(context.Background(), model.SearchParams{Query: "platform", PageSize: 2, Cursor: "2"})
	if err != nil {
		t.Fatalf("Search() second page error = %v", err)
	}
	if len(secondPage.Data) != 1 || secondPage.Data[0].Slug != "gamma" {
		t.Fatalf("second page data = %#v", secondPage.Data)
	}
	if secondPage.Meta.Cursor != nil {
		t.Fatalf("second page cursor = %#v, want nil", secondPage.Meta.Cursor)
	}
}
