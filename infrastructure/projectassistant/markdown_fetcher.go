package projectassistant

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/internal/markdownpolicy"
)

var (
	ErrMarkdownHostRejected = errors.New("markdown host rejected")
	ErrMarkdownFetchFailed  = errors.New("markdown fetch failed")
	ErrMarkdownTooLarge     = errors.New("markdown too large")
)

const (
	defaultMarkdownTimeout   = 30 * time.Second
	defaultMarkdownMaxBytes  = 512 * 1024
	defaultMarkdownTTL       = 10 * time.Minute
	defaultMarkdownRedirects = 2
	defaultCurlMaxTime       = 45
	defaultCurlConnectTime   = 12
)

type MarkdownFetcher struct {
	client             *http.Client
	cache              *MarkdownCache
	maxBytes           int64
	sourceURLValidator func(string) error
}

func NewMarkdownFetcher(cache *MarkdownCache) *MarkdownFetcher {
	fetcher := &MarkdownFetcher{cache: cache, maxBytes: defaultMarkdownMaxBytes}
	fetcher.sourceURLValidator = markdownpolicy.ValidateSourceURL
	fetcher.client = &http.Client{
		Timeout: defaultMarkdownTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > defaultMarkdownRedirects {
				return fmt.Errorf("%w: too many redirects", ErrMarkdownFetchFailed)
			}
			if err := fetcher.sourceURLValidator(req.URL.String()); err != nil {
				return mapMarkdownValidationError(err, "redirect host")
			}
			return nil
		},
	}

	return fetcher
}

func NewDefaultMarkdownCache() *MarkdownCache {
	return NewMarkdownCache(defaultMarkdownTTL)
}

func (f *MarkdownFetcher) Fetch(ctx context.Context, projectID string, sourceURL string) ([]services.MarkdownChunkAlias, error) {
	if f.cache != nil {
		if cached, ok := f.cache.Get(projectID, sourceURL); ok {
			return cached, nil
		}
	}

	normalizedURL := strings.TrimSpace(sourceURL)
	if err := f.sourceURLValidator(normalizedURL); err != nil {
		return nil, mapMarkdownValidationError(err, "source url")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, normalizedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: build request", ErrMarkdownFetchFailed)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		body, curlErr := f.fetchViaCurl(ctx, normalizedURL)
		if curlErr != nil {
			if f.cache != nil {
				if cached, ok := f.cache.GetStale(projectID, sourceURL); ok {
					return cached, nil
				}
			}
			return nil, fmt.Errorf("%w: %v", ErrMarkdownFetchFailed, err)
		}
		return f.parseAndCache(projectID, sourceURL, body)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, curlErr := f.fetchViaCurl(ctx, normalizedURL)
		if curlErr != nil {
			if f.cache != nil {
				if cached, ok := f.cache.GetStale(projectID, sourceURL); ok {
					return cached, nil
				}
			}
			return nil, fmt.Errorf("%w: status %d", ErrMarkdownFetchFailed, resp.StatusCode)
		}
		return f.parseAndCache(projectID, sourceURL, body)
	}

	limited := io.LimitReader(resp.Body, f.maxBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		fallbackBody, curlErr := f.fetchViaCurl(ctx, normalizedURL)
		if curlErr != nil {
			if f.cache != nil {
				if cached, ok := f.cache.GetStale(projectID, sourceURL); ok {
					return cached, nil
				}
			}
			return nil, fmt.Errorf("%w: read body", ErrMarkdownFetchFailed)
		}
		return f.parseAndCache(projectID, sourceURL, fallbackBody)
	}

	return f.parseAndCache(projectID, sourceURL, body)
}

func mapMarkdownValidationError(err error, detail string) error {
	if errors.Is(err, markdownpolicy.ErrSourceHostRejected) {
		return fmt.Errorf("%w: %s", ErrMarkdownHostRejected, detail)
	}

	return fmt.Errorf("%w: %s", ErrMarkdownFetchFailed, detail)
}

func (f *MarkdownFetcher) parseAndCache(projectID string, sourceURL string, body []byte) ([]services.MarkdownChunkAlias, error) {
	if int64(len(body)) > f.maxBytes {
		return nil, ErrMarkdownTooLarge
	}

	chunks := splitMarkdownIntoChunks(string(body))
	if f.cache != nil {
		f.cache.Set(projectID, sourceURL, chunks)
	}

	return chunks, nil
}

func (f *MarkdownFetcher) fetchViaCurl(ctx context.Context, sourceURL string) ([]byte, error) {
	cmd := exec.CommandContext(
		ctx,
		"curl",
		"-fsSL",
		"--connect-timeout", fmt.Sprintf("%d", defaultCurlConnectTime),
		"--max-time", fmt.Sprintf("%d", defaultCurlMaxTime),
		sourceURL,
	)

	body, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: curl fallback failed", ErrMarkdownFetchFailed)
	}

	return body, nil
}

func splitMarkdownIntoChunks(markdown string) []services.MarkdownChunkAlias {
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	chunks := make([]services.MarkdownChunkAlias, 0)
	currentHeading := "Overview"
	currentLines := make([]string, 0)
	flush := func() {
		body := strings.TrimSpace(strings.Join(currentLines, "\n"))
		if body == "" {
			return
		}
		chunks = append(chunks, services.MarkdownChunkAlias{Heading: currentHeading, Body: body})
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			flush()
			currentHeading = strings.TrimSpace(strings.TrimLeft(line, "#"))
			if currentHeading == "" {
				currentHeading = "Section"
			}
			currentLines = currentLines[:0]
			continue
		}
		if line == "" && len(currentLines) > 0 && currentLines[len(currentLines)-1] == "" {
			continue
		}
		currentLines = append(currentLines, line)
	}
	flush()
	if len(chunks) == 0 && strings.TrimSpace(markdown) != "" {
		chunks = append(chunks, services.MarkdownChunkAlias{Heading: "Overview", Body: strings.TrimSpace(markdown)})
	}
	return chunks
}
