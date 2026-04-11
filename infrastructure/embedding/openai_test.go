package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/sashabaranov/go-openai"
)

type stubEmbeddingClient struct {
	request  openai.EmbeddingRequest
	response openai.EmbeddingResponse
	err      error
	called   bool
}

func (s *stubEmbeddingClient) CreateEmbeddings(_ context.Context, request openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error) {
	s.called = true
	req, ok := request.(openai.EmbeddingRequest)
	if ok {
		s.request = req
	}
	return s.response, s.err
}

func TestOpenAIEmbeddingProvider_Generate(t *testing.T) {
	tests := []struct {
		name          string
		client        *stubEmbeddingClient
		input         string
		want          []float32
		wantErr       string
		assertRequest func(t *testing.T, req openai.EmbeddingRequest)
	}{
		{
			name: "returns provider vector on success",
			client: &stubEmbeddingClient{response: openai.EmbeddingResponse{Data: []openai.Embedding{{
				Embedding: []float32{0.12, 0.34, 0.56},
			}}}},
			input: "problem statement and architecture evidence",
			want:  []float32{0.12, 0.34, 0.56},
			assertRequest: func(t *testing.T, req openai.EmbeddingRequest) {
				t.Helper()
				if req.Model != openai.SmallEmbedding3 {
					t.Fatalf("Model = %q, want %q", req.Model, openai.SmallEmbedding3)
				}
				inputs, ok := req.Input.([]string)
				if !ok {
					t.Fatalf("Input type = %T, want []string", req.Input)
				}
				if len(inputs) != 1 || inputs[0] != "problem statement and architecture evidence" {
					t.Fatalf("Input = %#v, want single input text", inputs)
				}
			},
		},
		{
			name:    "bubbles up openai failures",
			client:  &stubEmbeddingClient{err: errors.New("upstream timeout")},
			input:   "semantic evidence",
			wantErr: "failed to generate embedding from OpenAI: upstream timeout",
		},
		{
			name:    "fails on empty embedding payload",
			client:  &stubEmbeddingClient{response: openai.EmbeddingResponse{}},
			input:   "semantic evidence",
			wantErr: "openai returned empty embedding data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := newOpenAIEmbeddingProvider(tt.client)

			got, err := provider.Generate(context.Background(), tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("Generate() error = nil, want non-nil")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("Generate() error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !tt.client.called {
				t.Fatal("CreateEmbeddings was not called")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len(Generate()) = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("Generate()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
			if tt.assertRequest != nil {
				tt.assertRequest(t, tt.client.request)
			}
		})
	}
}
