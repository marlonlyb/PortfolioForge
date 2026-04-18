package projectassistant

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/model"
	"github.com/sashabaranov/go-openai"
)

type stubOpenAIChatClient struct {
	request openai.ChatCompletionRequest
	resp    openai.ChatCompletionResponse
	err     error
}

func (s *stubOpenAIChatClient) CreateChatCompletion(_ context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	s.request = request
	return s.resp, s.err
}

func TestBuildAssistantSystemPromptLocksSelectedLocaleAndBrevity(t *testing.T) {
	tests := []struct {
		name      string
		language  string
		fragments []string
	}{
		{
			name:     "spanish default",
			language: "es",
			fragments: []string{
				"Responde solo en español",
				"independientemente del idioma de la pregunta",
				"Mantén la respuesta breve, directa y fundamentada",
				"No reproduzcas markdown crudo",
				"Si la documentación no alcanza, dilo brevemente en español",
			},
		},
		{
			name:     "english",
			language: "en",
			fragments: []string{
				"Answer only in English",
				"regardless of the language used in the question",
				"Keep the answer concise, direct, and grounded",
				"Do not reproduce raw markdown formatting",
				"If the documentation is insufficient, say so briefly in English",
			},
		},
		{
			name:     "catalan",
			language: "ca",
			fragments: []string{
				"Respon només en català",
				"independentment de l’idioma de la pregunta",
				"Mantén la resposta breu, directa i fonamentada",
				"No reprodueixis format markdown cru",
				"Si la documentació és insuficient, digues-ho breument en català",
			},
		},
		{
			name:     "german",
			language: "de",
			fragments: []string{
				"Antworte nur auf Deutsch",
				"unabhängig von der Sprache der Frage",
				"Halte die Antwort kurz, direkt und belastbar",
				"Gib kein rohes Markdown",
				"Wenn die Dokumentation nicht ausreicht, sage das kurz auf Deutsch",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildAssistantSystemPrompt(tt.language)
			for _, fragment := range tt.fragments {
				if !strings.Contains(prompt, fragment) {
					t.Fatalf("prompt %q missing fragment %q", prompt, fragment)
				}
			}
		})
	}
}

func TestBuildAssistantUserPayloadReinforcesSelectedLocaleAndIncludesSections(t *testing.T) {
	tests := []struct {
		name            string
		language        string
		question        string
		wantInstruction string
	}{
		{name: "spanish", language: "es", question: "How does deployment work?", wantInstruction: "Idioma de respuesta obligatorio: español (es)."},
		{name: "english", language: "en", question: "¿Cómo funciona el despliegue?", wantInstruction: "Required response language: English (en)."},
		{name: "catalan", language: "ca", question: "How does deployment work?", wantInstruction: "Idioma de resposta obligatori: català (ca)."},
		{name: "german", language: "de", question: "¿Cómo funciona el despliegue?", wantInstruction: "Erforderliche Antwortsprache: Deutsch (de)."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := buildAssistantUserPayload(services.ProjectAssistantAnswerInput{
				ProjectName: "PortfolioForge",
				Language:    tt.language,
				Question:    tt.question,
				Sections: []services.MarkdownChunkAlias{
					{Heading: "Deployment", Body: "Uses CI jobs and release gates."},
					{Heading: "Observability", Body: "Alerts and dashboards are configured."},
				},
			})

			for _, fragment := range []string{
				tt.wantInstruction,
				"Project: PortfolioForge",
				"Question: " + tt.question,
				"Relevant markdown sections:",
				"## Deployment\nUses CI jobs and release gates.",
				"## Observability\nAlerts and dashboards are configured.",
			} {
				if !strings.Contains(payload, fragment) {
					t.Fatalf("payload %q missing fragment %q", payload, fragment)
				}
			}
		})
	}
}

func TestGenerateAnswerBuildsBoundedChatRequestForEachLocale(t *testing.T) {
	tests := []struct {
		name                   string
		language               string
		question               string
		wantSystemInstruction  string
		wantPayloadInstruction string
	}{
		{name: "spanish", language: "es", question: "How does deployment work?", wantSystemInstruction: "Responde solo en español", wantPayloadInstruction: "Idioma de respuesta obligatorio: español (es)."},
		{name: "english", language: "en", question: "¿Cómo funciona el despliegue?", wantSystemInstruction: "Answer only in English", wantPayloadInstruction: "Required response language: English (en)."},
		{name: "catalan", language: "ca", question: "How does deployment work?", wantSystemInstruction: "Respon només en català", wantPayloadInstruction: "Idioma de resposta obligatori: català (ca)."},
		{name: "german", language: "de", question: "¿Cómo funciona el despliegue?", wantSystemInstruction: "Antworte nur auf Deutsch", wantPayloadInstruction: "Erforderliche Antwortsprache: Deutsch (de)."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &stubOpenAIChatClient{
				resp: openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "  Grounded answer.  "}}},
				},
			}
			provider := &OpenAIProvider{client: client}

			input := services.ProjectAssistantAnswerInput{
				ProjectName: "PortfolioForge",
				Language:    tt.language,
				Question:    tt.question,
				History: []model.ProjectAssistantMessage{
					{Role: "user", Content: "Previous question"},
					{Role: "assistant", Content: "Previous answer"},
				},
				Sections: []services.MarkdownChunkAlias{
					{Heading: "Deployment", Body: "Uses CI jobs and release gates."},
					{Heading: "Observability", Body: "Alerts and dashboards are configured."},
				},
			}

			answer, err := provider.GenerateAnswer(context.Background(), input)
			if err != nil {
				t.Fatalf("GenerateAnswer() error = %v", err)
			}
			if answer != "Grounded answer." {
				t.Fatalf("answer = %q, want %q", answer, "Grounded answer.")
			}

			if client.request.Model != openai.GPT4oMini {
				t.Fatalf("Model = %q, want %q", client.request.Model, openai.GPT4oMini)
			}
			if client.request.Temperature != 0.2 {
				t.Fatalf("Temperature = %v, want 0.2", client.request.Temperature)
			}
			if client.request.MaxTokens != assistantAnswerMaxTokens {
				t.Fatalf("MaxTokens = %d, want %d", client.request.MaxTokens, assistantAnswerMaxTokens)
			}
			if len(client.request.Messages) != 4 {
				t.Fatalf("messages length = %d, want 4", len(client.request.Messages))
			}

			if client.request.Messages[0].Role != openai.ChatMessageRoleSystem || !strings.Contains(client.request.Messages[0].Content, tt.wantSystemInstruction) {
				t.Fatalf("system message = %#v", client.request.Messages[0])
			}
			if client.request.Messages[1].Role != openai.ChatMessageRoleUser || client.request.Messages[1].Content != "Previous question" {
				t.Fatalf("first history message = %#v", client.request.Messages[1])
			}
			if client.request.Messages[2].Role != openai.ChatMessageRoleAssistant || client.request.Messages[2].Content != "Previous answer" {
				t.Fatalf("second history message = %#v", client.request.Messages[2])
			}

			finalMessage := client.request.Messages[3]
			if finalMessage.Role != openai.ChatMessageRoleUser {
				t.Fatalf("final role = %q, want %q", finalMessage.Role, openai.ChatMessageRoleUser)
			}
			if finalMessage.Content != buildAssistantUserPayload(input) {
				t.Fatalf("final message = %q, want %q", finalMessage.Content, buildAssistantUserPayload(input))
			}
			for _, fragment := range []string{
				tt.wantPayloadInstruction,
				"Project: PortfolioForge",
				"Question: " + tt.question,
				"## Deployment\nUses CI jobs and release gates.",
				"## Observability\nAlerts and dashboards are configured.",
			} {
				if !strings.Contains(finalMessage.Content, fragment) {
					t.Fatalf("final message %q missing fragment %q", finalMessage.Content, fragment)
				}
			}
		})
	}
}

func TestGenerateAnswerRequiresConfiguredClient(t *testing.T) {
	provider := &OpenAIProvider{}

	_, err := provider.GenerateAnswer(context.Background(), services.ProjectAssistantAnswerInput{})
	if !errors.Is(err, services.ErrAssistantUnavailable) {
		t.Fatalf("error = %v, want %v", err, services.ErrAssistantUnavailable)
	}
}
