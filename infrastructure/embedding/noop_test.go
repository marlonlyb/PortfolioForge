package embedding

import (
	"context"
	"testing"
)

func TestNoOpEmbeddingProvider_Generate(t *testing.T) {
	p := NewNoOpEmbeddingProvider()
	vec, err := p.Generate(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vec != nil {
		t.Fatalf("expected nil vector, got %v", vec)
	}
}

func TestNoOpEmbeddingProvider_Dimension(t *testing.T) {
	p := NewNoOpEmbeddingProvider()
	if d := p.Dimension(); d != 1536 {
		t.Fatalf("expected dimension 1536, got %d", d)
	}
}
