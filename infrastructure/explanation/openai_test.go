package explanation

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubChatCompletionClient struct {
	request  openai.ChatCompletionRequest
	response openai.ChatCompletionResponse
	err      error
	called   bool
}

func (s *stubChatCompletionClient) CreateChatCompletion(_ context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	s.called = true
	s.request = request
	return s.response, s.err
}

type stubExplanationFallback struct {
	response string
	called   bool
}

func (s *stubExplanationFallback) Explain(_ context.Context, _ model.Project, _ model.EvidenceTrace, _ string) (string, error) {
	s.called = true
	return s.response, nil
}

func TestOpenAIExplanationProvider_Explain(t *testing.T) {
	ctx := context.Background()
	project := model.Project{ID: uuid.New(), Name: "PortfolioForge"}
	evidence := model.EvidenceTrace{
		ProjectID: project.ID,
		Fields: []model.EvidenceField{{
			Field:       "problem_statement",
			MatchedText: "búsqueda semántica inconsistente",
			MatchType:   model.MatchTypeFTS,
			Score:       0.9,
		}},
	}

	tests := []struct {
		name             string
		client           *stubChatCompletionClient
		fallbackResponse string
		query            string
		want             string
		wantFallback     bool
		assertRequest    func(t *testing.T, req openai.ChatCompletionRequest)
	}{
		{
			name: "builds bounded prompt and trims model output",
			client: &stubChatCompletionClient{response: openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{Content: "  Este proyecto resuelve la búsqueda semántica inconsistente con evidencia del problema reportado.  "},
			}}}},
			fallbackResponse: "fallback",
			query:            "búsqueda semántica",
			want:             "Este proyecto resuelve la búsqueda semántica inconsistente con evidencia del problema reportado.",
			assertRequest: func(t *testing.T, req openai.ChatCompletionRequest) {
				t.Helper()
				if req.Model != openai.GPT4oMini {
					t.Fatalf("Model = %q, want %q", req.Model, openai.GPT4oMini)
				}
				if req.MaxTokens != 60 {
					t.Fatalf("MaxTokens = %d, want 60", req.MaxTokens)
				}
				if req.Temperature != 0.1 {
					t.Fatalf("Temperature = %v, want 0.1", req.Temperature)
				}
				if len(req.Messages) != 2 {
					t.Fatalf("Messages len = %d, want 2", len(req.Messages))
				}
				system := req.Messages[0].Content
				for _, fragment := range []string{
					"exactamente una sola oración",
					"devuelve únicamente esa oración",
					"no inventes",
					"no extrapoles",
					"no afirmes nada que no esté explícitamente respaldado",
				} {
					if !strings.Contains(system, fragment) {
						t.Fatalf("system prompt missing fragment %q: %q", fragment, system)
					}
				}
				if !strings.Contains(req.Messages[1].Content, "problem_statement") {
					t.Fatalf("user prompt missing evidence field: %q", req.Messages[1].Content)
				}
			},
		},
		{
			name:             "falls back when openai errors",
			client:           &stubChatCompletionClient{err: errors.New("upstream timeout")},
			fallbackResponse: "fallback explanation",
			query:            "búsqueda semántica",
			want:             "fallback explanation",
			wantFallback:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fallback := &stubExplanationFallback{response: tt.fallbackResponse}
			provider := newOpenAIExplanationProvider(tt.client, fallback)

			got, err := provider.Explain(ctx, project, evidence, tt.query)
			if err != nil {
				t.Fatalf("Explain() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Explain() = %q, want %q", got, tt.want)
			}
			if fallback.called != tt.wantFallback {
				t.Fatalf("fallback called = %v, want %v", fallback.called, tt.wantFallback)
			}
			if tt.assertRequest != nil {
				tt.assertRequest(t, tt.client.request)
			}
		})
	}
}
