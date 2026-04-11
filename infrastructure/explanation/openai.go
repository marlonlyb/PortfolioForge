package explanation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/model"
	"github.com/sashabaranov/go-openai"
)

type OpenAIExplanationProvider struct {
	client   chatCompletionClient
	fallback search.ExplanationProvider
}

type chatCompletionClient interface {
	CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func NewOpenAIExplanationProvider(apiKey string, fallback search.ExplanationProvider) search.ExplanationProvider {
	return newOpenAIExplanationProvider(openai.NewClient(apiKey), fallback)
}

func newOpenAIExplanationProvider(client chatCompletionClient, fallback search.ExplanationProvider) search.ExplanationProvider {
	return &OpenAIExplanationProvider{
		client:   client,
		fallback: fallback,
	}
}

func (p *OpenAIExplanationProvider) Explain(ctx context.Context, project model.Project, evidence model.EvidenceTrace, query string) (string, error) {
	// If query is empty or no evidence, defer to fallback
	if query == "" || len(evidence.Fields) == 0 {
		return p.fallback.Explain(ctx, project, evidence, query)
	}

	evidenceBytes, err := json.Marshal(evidence.Fields)
	if err != nil {
		log.Printf("Error marshaling evidence for OpenAI explanation: %v", err)
		return p.fallback.Explain(ctx, project, evidence, query)
	}

	systemInstruction := "Responde en español con exactamente una sola oración y devuelve únicamente esa oración, sin prefijos, listas, comillas ni explicación adicional. Usa solo la evidencia proporcionada; no inventes, no extrapoles y no afirmes nada que no esté explícitamente respaldado por los campos entregados."
	userPrompt := fmt.Sprintf("Query: %s\nProject Name: %s\nEvidence Trace: %s\n", query, project.Name, string(evidenceBytes))

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemInstruction,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		MaxTokens:   60,
		Temperature: 0.1,
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("Error generating explanation from OpenAI: %v", err)
		return p.fallback.Explain(ctx, project, evidence, query)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		log.Println("OpenAI returned empty explanation")
		return p.fallback.Explain(ctx, project, evidence, query)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
