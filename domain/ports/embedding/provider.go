package embedding

import "context"

// EmbeddingProvider generates vector embeddings from text input.
type EmbeddingProvider interface {
	Generate(ctx context.Context, text string) ([]float32, error)
	Dimension() int
}
