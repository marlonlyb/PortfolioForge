package projectassistant

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/marlonlyb/portfolioforge/internal/markdownpolicy"
)

func TestMarkdownFetcherRejectsHostOutsideAllowlist(t *testing.T) {
	fetcher := NewMarkdownFetcher(NewMarkdownCache(time.Minute))
	_, err := fetcher.Fetch(context.Background(), "project-1", "https://example.com/doc.md")
	if !errors.Is(err, ErrMarkdownHostRejected) {
		t.Fatalf("error = %v, want ErrMarkdownHostRejected", err)
	}
}

func TestMarkdownFetcherUsesSharedMarkdownPolicy(t *testing.T) {
	if err := markdownpolicy.ValidateSourceURL("https://example.com/doc.md"); !errors.Is(err, markdownpolicy.ErrSourceHostRejected) {
		t.Fatalf("shared policy error = %v, want ErrSourceHostRejected", err)
	}

	fetcher := NewMarkdownFetcher(NewMarkdownCache(time.Minute))
	_, err := fetcher.Fetch(context.Background(), "project-1", "https://example.com/doc.md")
	if !errors.Is(err, ErrMarkdownHostRejected) {
		t.Fatalf("fetch error = %v, want ErrMarkdownHostRejected", err)
	}
}

func TestMarkdownFetcherCachesChunks(t *testing.T) {
	serverHits := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		serverHits++
		_, _ = w.Write([]byte("# Architecture\nUses Go and PostgreSQL."))
	}))
	defer server.Close()

	fetcher := NewMarkdownFetcher(NewMarkdownCache(time.Minute))
	fetcher.client = server.Client()
	fetcher.sourceURLValidator = func(raw string) error {
		trimmed := strings.TrimSpace(raw)
		if trimmed == server.URL {
			return nil
		}
		return markdownpolicy.ValidateSourceURL(trimmed)
	}

	chunksA, err := fetcher.Fetch(context.Background(), "project-1", server.URL)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	chunksB, err := fetcher.Fetch(context.Background(), "project-1", server.URL)
	if err != nil {
		t.Fatalf("Fetch() second error = %v", err)
	}
	if serverHits != 1 {
		t.Fatalf("server hits = %d, want 1", serverHits)
	}
	if len(chunksA) != 1 || len(chunksB) != 1 || chunksA[0].Heading != "Architecture" {
		t.Fatalf("chunks = %#v / %#v", chunksA, chunksB)
	}
}

func TestMarkdownFetcherCacheKeyIncludesSourceURL(t *testing.T) {
	responses := map[string]string{
		"/first.md":  "# Architecture\nUses Go and PostgreSQL.",
		"/second.md": "# Results\nImproved conversion by 18%.",
	}
	serverHits := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHits++
		_, _ = w.Write([]byte(responses[r.URL.Path]))
	}))
	defer server.Close()

	firstURL := server.URL + "/first.md"
	secondURL := server.URL + "/second.md"

	fetcher := NewMarkdownFetcher(NewMarkdownCache(time.Minute))
	fetcher.client = server.Client()
	fetcher.sourceURLValidator = func(raw string) error {
		trimmed := strings.TrimSpace(raw)
		if trimmed == firstURL || trimmed == secondURL {
			return nil
		}
		return markdownpolicy.ValidateSourceURL(trimmed)
	}

	firstChunks, err := fetcher.Fetch(context.Background(), "project-1", firstURL)
	if err != nil {
		t.Fatalf("Fetch() first error = %v", err)
	}
	secondChunks, err := fetcher.Fetch(context.Background(), "project-1", secondURL)
	if err != nil {
		t.Fatalf("Fetch() second error = %v", err)
	}

	if serverHits != 2 {
		t.Fatalf("server hits = %d, want 2 after source url change", serverHits)
	}
	if len(firstChunks) != 1 || firstChunks[0].Heading != "Architecture" {
		t.Fatalf("first chunks = %#v", firstChunks)
	}
	if len(secondChunks) != 1 || secondChunks[0].Heading != "Results" {
		t.Fatalf("second chunks = %#v", secondChunks)
	}
	if firstChunks[0].Heading == secondChunks[0].Heading {
		t.Fatalf("cache returned stale chunks: %#v / %#v", firstChunks, secondChunks)
	}
}

func TestMarkdownFetcherReturnsFetchFailedOnTimeout(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte("# Architecture\nToo slow."))
	}))
	defer server.Close()

	fetcher := NewMarkdownFetcher(NewMarkdownCache(time.Minute))
	fetcher.client = server.Client()
	fetcher.client.Timeout = 10 * time.Millisecond
	fetcher.sourceURLValidator = func(raw string) error {
		if strings.TrimSpace(raw) == server.URL {
			return nil
		}
		return markdownpolicy.ValidateSourceURL(raw)
	}

	_, err := fetcher.Fetch(context.Background(), "project-1", server.URL)
	if !errors.Is(err, ErrMarkdownFetchFailed) {
		t.Fatalf("error = %v, want ErrMarkdownFetchFailed", err)
	}
}

func TestMarkdownFetcherRejectsResponsesOverSizeCap(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "# Architecture\n1234567890")
	}))
	defer server.Close()

	fetcher := NewMarkdownFetcher(NewMarkdownCache(time.Minute))
	fetcher.client = server.Client()
	fetcher.sourceURLValidator = func(raw string) error {
		if strings.TrimSpace(raw) == server.URL {
			return nil
		}
		return markdownpolicy.ValidateSourceURL(raw)
	}
	fetcher.maxBytes = 8

	_, err := fetcher.Fetch(context.Background(), "project-1", server.URL)
	if !errors.Is(err, ErrMarkdownTooLarge) {
		t.Fatalf("error = %v, want ErrMarkdownTooLarge", err)
	}
}
