package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubAssistantRepo struct {
	context model.ProjectAssistantContext
	err     error
}

func (s *stubAssistantRepo) GetAssistantContextBySlug(context.Context, string) (model.ProjectAssistantContext, error) {
	return s.context, s.err
}

type stubAssistantRetriever struct {
	chunks []MarkdownChunkAlias
	err    error
}

func (s *stubAssistantRetriever) Fetch(context.Context, string, string) ([]MarkdownChunkAlias, error) {
	return s.chunks, s.err
}

type stubAssistantProvider struct {
	called bool
	input  ProjectAssistantAnswerInput
	resp   string
	err    error
}

func (s *stubAssistantProvider) GenerateAnswer(_ context.Context, input ProjectAssistantAnswerInput) (string, error) {
	s.called = true
	s.input = input
	return s.resp, s.err
}

func TestNormalizeAssistantHistory(t *testing.T) {
	tests := []struct {
		name    string
		history []model.ProjectAssistantMessage
		want    []model.ProjectAssistantMessage
	}{
		{
			name:    "returns nil for empty history",
			history: nil,
			want:    nil,
		},
		{
			name: "keeps newest valid eight messages",
			history: []model.ProjectAssistantMessage{
				{Role: "user", Content: "Message 1"},
				{Role: "assistant", Content: "Message 2"},
				{Role: "user", Content: "Message 3"},
				{Role: "assistant", Content: "Message 4"},
				{Role: "user", Content: "Message 5"},
				{Role: "assistant", Content: "Message 6"},
				{Role: "user", Content: "Message 7"},
				{Role: "assistant", Content: "Message 8"},
				{Role: "user", Content: "Message 9"},
				{Role: "assistant", Content: "Message 10"},
			},
			want: []model.ProjectAssistantMessage{
				{Role: "user", Content: "Message 3"},
				{Role: "assistant", Content: "Message 4"},
				{Role: "user", Content: "Message 5"},
				{Role: "assistant", Content: "Message 6"},
				{Role: "user", Content: "Message 7"},
				{Role: "assistant", Content: "Message 8"},
				{Role: "user", Content: "Message 9"},
				{Role: "assistant", Content: "Message 10"},
			},
		},
		{
			name: "ignores malformed entries before counting the limit",
			history: []model.ProjectAssistantMessage{
				{Role: "system", Content: "Ignore"},
				{Role: "user", Content: "  Message 1  "},
				{Role: "assistant", Content: "Message 2"},
				{Role: "assistant", Content: "   "},
				{Role: "user", Content: "Message 3"},
				{Role: "assistant", Content: "Message 4"},
				{Role: "user", Content: "Message 5"},
				{Role: "assistant", Content: "Message 6"},
				{Role: "user", Content: "Message 7"},
				{Role: "assistant", Content: "Message 8"},
				{Role: "user", Content: "Message 9"},
				{Role: "assistant", Content: "Message 10"},
			},
			want: []model.ProjectAssistantMessage{
				{Role: "user", Content: "Message 1"},
				{Role: "assistant", Content: "Message 2"},
				{Role: "user", Content: "Message 3"},
				{Role: "assistant", Content: "Message 4"},
				{Role: "user", Content: "Message 5"},
				{Role: "assistant", Content: "Message 6"},
				{Role: "user", Content: "Message 7"},
				{Role: "assistant", Content: "Message 8"},
				{Role: "user", Content: "Message 9"},
				{Role: "assistant", Content: "Message 10"},
			}[2:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeAssistantHistory(tt.history)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("normalizeAssistantHistory() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNormalizeAssistantLanguage(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "defaults to spanish", raw: "", want: "es"},
		{name: "supports catalan alias", raw: "Catalan", want: "ca"},
		{name: "supports english alias", raw: " ENGLISH ", want: "en"},
		{name: "supports german code", raw: "de", want: "de"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeAssistantLanguage(tt.raw); got != tt.want {
				t.Fatalf("normalizeAssistantLanguage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProjectAssistantAnswerUsesNormalizedHistoryForProviderAndRetrieval(t *testing.T) {
	provider := &stubAssistantProvider{resp: "Detailed grounded answer."}
	service := NewProjectAssistant(
		&stubAssistantRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}},
		&stubAssistantRetriever{chunks: []MarkdownChunkAlias{
			{Heading: "Dropped Topic", Body: "This section mentions droppedtopic only."},
			{Heading: "Recent Topic", Body: "This section explains keeptopic and rollout details."},
		}},
		provider,
	)

	history := make([]model.ProjectAssistantMessage, 0, 10)
	for index := 1; index <= 10; index++ {
		role := "user"
		if index%2 == 0 {
			role = "assistant"
		}
		content := fmt.Sprintf("Message %d about keeptopic", index)
		if index <= 2 {
			content = fmt.Sprintf("Message %d about droppedtopic", index)
		}
		history = append(history, model.ProjectAssistantMessage{Role: role, Content: content})
	}

	response, err := service.Answer(context.Background(), "portfolioforge", model.ProjectAssistantRequest{
		Question: "What changed recently?",
		History:  history,
		Lang:     "en",
	})
	if err != nil {
		t.Fatalf("Answer() error = %v", err)
	}
	if response.Answer != "Detailed grounded answer." {
		t.Fatalf("answer = %q", response.Answer)
	}

	if len(provider.input.History) != assistantHistoryLimit {
		t.Fatalf("provider history length = %d, want %d", len(provider.input.History), assistantHistoryLimit)
	}
	if provider.input.History[0].Content != "Message 3 about keeptopic" {
		t.Fatalf("provider first history = %#v, want Message 3 about keeptopic", provider.input.History[0])
	}
	if provider.input.History[len(provider.input.History)-1].Content != "Message 10 about keeptopic" {
		t.Fatalf("provider last history = %#v, want Message 10 about keeptopic", provider.input.History[len(provider.input.History)-1])
	}
	if len(provider.input.Sections) == 0 || provider.input.Sections[0].Heading != "Recent Topic" {
		t.Fatalf("selected sections = %#v, want Recent Topic first", provider.input.Sections)
	}
	if provider.input.Language != "en" {
		t.Fatalf("language = %q, want en", provider.input.Language)
	}
}

func TestSelectRelevantChunksFallsBackWhenHistoryNormalizesToEmpty(t *testing.T) {
	chunks := []markdownChunk{
		{Heading: "Architecture", Body: "Uses Go and Echo."},
		{Heading: "Results", Body: "Improved conversion by 18%."},
		{Heading: "Integrations", Body: "Connects with analytics."},
		{Heading: "Tradeoffs", Body: "Documents implementation tradeoffs."},
		{Heading: "Ignored", Body: "Fifth chunk should be excluded from fallback."},
	}

	selected := selectRelevantChunks("  ", normalizeAssistantHistory([]model.ProjectAssistantMessage{{Role: "system", Content: "ignore"}, {Role: "user", Content: "   "}}), chunks)
	if len(selected) != assistantRelevantSectionLimit {
		t.Fatalf("fallback chunks length = %d, want %d", len(selected), assistantRelevantSectionLimit)
	}
	if selected[0].Heading != "Architecture" || selected[3].Heading != "Tradeoffs" {
		t.Fatalf("fallback chunks = %#v, want first four original chunks", selected)
	}
}

func TestSelectRelevantChunksKeepsHighestSignalBoundedToFourSections(t *testing.T) {
	chunks := []markdownChunk{
		{Heading: "Authentication", Body: "Token auth and session flow details."},
		{Heading: "Deployment", Body: "Deployment steps and rollout plan."},
		{Heading: "Observability", Body: "Monitoring, logs, alerts, and observability dashboards."},
		{Heading: "Architecture", Body: "Architecture decisions, service boundaries, and Go adapters."},
		{Heading: "Tradeoffs", Body: "Tradeoffs for architecture and deployment."},
		{Heading: "Roadmap", Body: "Future ideas only."},
	}

	selected := selectRelevantChunks("Explain the architecture deployment tradeoffs and monitoring", nil, chunks)
	if len(selected) != assistantRelevantSectionLimit {
		t.Fatalf("selected chunks length = %d, want %d", len(selected), assistantRelevantSectionLimit)
	}
	headings := make(map[string]bool, len(selected))
	for _, chunk := range selected {
		headings[chunk.Heading] = true
		if chunk.Heading == "Roadmap" {
			t.Fatalf("unexpected low-signal chunk selected: %#v", selected)
		}
	}
	for _, heading := range []string{"Architecture", "Deployment", "Tradeoffs"} {
		if !headings[heading] {
			t.Fatalf("selected headings = %#v, missing %q", selected, heading)
		}
	}
}

func TestProjectAssistantAnswerRejectsUnavailableProject(t *testing.T) {
	provider := &stubAssistantProvider{}
	service := NewProjectAssistant(
		&stubAssistantRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true}},
		&stubAssistantRetriever{},
		provider,
	)

	_, err := service.Answer(context.Background(), "portfolioforge", model.ProjectAssistantRequest{Question: "What happened?"})
	if !errors.Is(err, ErrAssistantUnavailable) {
		t.Fatalf("error = %v, want ErrAssistantUnavailable", err)
	}
	if provider.called {
		t.Fatal("provider should not be invoked when assistant is unavailable")
	}
}
