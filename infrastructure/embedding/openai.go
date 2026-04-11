package embedding

import (
	"context"
	"fmt"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	"github.com/sashabaranov/go-openai"
)

type embeddingClient interface {
	CreateEmbeddings(ctx context.Context, request openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error)
}

type OpenAIEmbeddingProvider struct {
	client embeddingClient
}

func NewOpenAIEmbeddingProvider(apiKey string) embedding.EmbeddingProvider {
	return newOpenAIEmbeddingProvider(openai.NewClient(apiKey))
}

func newOpenAIEmbeddingProvider(client embeddingClient) embedding.EmbeddingProvider {
	return &OpenAIEmbeddingProvider{
		client: client,
	}
}

func (p *OpenAIEmbeddingProvider) Generate(ctx context.Context, text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3, // text-embedding-3-small
	}

	resp, err := p.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding from OpenAI: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("openai returned empty embedding data")
	}

	return resp.Data[0].Embedding, nil
}

func (p *OpenAIEmbeddingProvider) Dimension() int {
	return 1536 // Standard for text-embedding-3-small
}
