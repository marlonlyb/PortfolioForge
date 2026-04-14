package services

import (
	"context"
	"errors"
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

func TestProjectAssistantAnswerSelectsRelevantSections(t *testing.T) {
	provider := &stubAssistantProvider{resp: "Detailed grounded answer."}
	service := NewProjectAssistant(
		&stubAssistantRepo{context: model.ProjectAssistantContext{ID: uuid.New(), Name: "PortfolioForge", Active: true, SourceMarkdownURL: "https://mlbautomation.com/docs.md"}},
		&stubAssistantRetriever{chunks: []MarkdownChunkAlias{{Heading: "Architecture", Body: "Uses Go, PostgreSQL and Echo."}, {Heading: "Results", Body: "Improved conversion by 18%."}}},
		provider,
	)

	response, err := service.Answer(context.Background(), "portfolioforge", model.ProjectAssistantRequest{Question: "How is the architecture implemented?", Lang: "en"})
	if err != nil {
		t.Fatalf("Answer() error = %v", err)
	}
	if response.Answer != "Detailed grounded answer." {
		t.Fatalf("answer = %q", response.Answer)
	}
	if len(provider.input.Sections) == 0 || provider.input.Sections[0].Heading != "Architecture" {
		t.Fatalf("selected sections = %#v, want Architecture first", provider.input.Sections)
	}
	if provider.input.Language != "en" {
		t.Fatalf("language = %q, want en", provider.input.Language)
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
