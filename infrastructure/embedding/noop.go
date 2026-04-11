package embedding

import "context"

// NoOpEmbeddingProvider is a no-op implementation of EmbeddingProvider.
// It returns nil vectors and reports dimension 1536 (OpenAI text-embedding-ada-002 compatible).
type NoOpEmbeddingProvider struct{}

// NewNoOpEmbeddingProvider creates a new NoOpEmbeddingProvider.
func NewNoOpEmbeddingProvider() *NoOpEmbeddingProvider {
	return &NoOpEmbeddingProvider{}
}

// Generate returns a nil vector (no embedding generated).
func (n *NoOpEmbeddingProvider) Generate(_ context.Context, _ string) ([]float32, error) {
	return nil, nil
}

// Dimension returns 1536, the standard OpenAI embedding dimension.
func (n *NoOpEmbeddingProvider) Dimension() int {
	return 1536
}
