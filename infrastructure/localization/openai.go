package localization

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const localizationTranslationMaxTokens = 4000

type chatCompletionClient interface {
	CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

type OpenAITranslator struct {
	client chatCompletionClient
}

func NewOpenAITranslator(apiKey string) *OpenAITranslator {
	if strings.TrimSpace(apiKey) == "" {
		return nil
	}
	return &OpenAITranslator{client: openai.NewClient(apiKey)}
}

func (t *OpenAITranslator) TranslateFields(ctx context.Context, sourceLocale string, targetLocale string, fields map[string]json.RawMessage) (map[string]json.RawMessage, error) {
	if t == nil || t.client == nil {
		return nil, fmt.Errorf("openai translator is not configured")
	}

	payloadBytes, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("marshal translation payload: %w", err)
	}

	request := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You translate structured portfolio content. Return only valid JSON with the exact same top-level keys and the same nested structure. Translate human-facing Spanish text into the requested language, keep proper nouns and technical terms when that is more natural, never mix languages unnecessarily, and never add commentary.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Source locale: %s\nTarget locale: %s\nTranslate this JSON value-by-value and preserve keys exactly:\n%s", sourceLocale, targetLocale, string(payloadBytes)),
			},
		},
		Temperature: 0.2,
		MaxTokens:   localizationTranslationMaxTokens,
	}

	response, err := t.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no translation choices")
	}

	rawContent := strings.TrimSpace(response.Choices[0].Message.Content)
	rawContent = strings.TrimPrefix(rawContent, "```json")
	rawContent = strings.TrimPrefix(rawContent, "```")
	rawContent = strings.TrimSuffix(rawContent, "```")
	rawContent = strings.TrimSpace(rawContent)

	var translated map[string]json.RawMessage
	if err := json.Unmarshal([]byte(rawContent), &translated); err != nil {
		return nil, fmt.Errorf("parse translated JSON: %w", err)
	}

	return translated, nil
}
