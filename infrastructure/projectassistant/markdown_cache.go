package projectassistant

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/marlonlyb/portfolioforge/domain/services"
)

type cachedMarkdownDocument struct {
	chunks    []services.MarkdownChunkAlias
	expiresAt time.Time
}

type persistedMarkdownDocument struct {
	Chunks    []services.MarkdownChunkAlias `json:"chunks"`
	ExpiresAt time.Time                     `json:"expires_at"`
}

type MarkdownCache struct {
	mu      sync.RWMutex
	ttl     time.Duration
	entries map[string]cachedMarkdownDocument
	baseDir string
}

func NewMarkdownCache(ttl time.Duration) *MarkdownCache {
	baseDir := filepath.Join(os.TempDir(), "portfolioforge-projectassistant-cache")
	_ = os.MkdirAll(baseDir, 0o755)
	return &MarkdownCache{ttl: ttl, entries: map[string]cachedMarkdownDocument{}, baseDir: baseDir}
}

func (c *MarkdownCache) Get(projectID string, sourceURL string) ([]services.MarkdownChunkAlias, bool) {
	if c == nil {
		return nil, false
	}

	cacheKey := buildMarkdownCacheKey(projectID, sourceURL)

	c.mu.RLock()
	entry, ok := c.entries[cacheKey]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		if ok {
			c.mu.Lock()
			delete(c.entries, cacheKey)
			c.mu.Unlock()
		}
		return c.getPersisted(cacheKey, false)
	}

	return cloneChunks(entry.chunks), true
}

func (c *MarkdownCache) GetStale(projectID string, sourceURL string) ([]services.MarkdownChunkAlias, bool) {
	if c == nil {
		return nil, false
	}

	cacheKey := buildMarkdownCacheKey(projectID, sourceURL)

	c.mu.RLock()
	entry, ok := c.entries[cacheKey]
	c.mu.RUnlock()
	if ok && len(entry.chunks) > 0 {
		return cloneChunks(entry.chunks), true
	}

	return c.getPersisted(cacheKey, true)
}

func (c *MarkdownCache) Set(projectID string, sourceURL string, chunks []services.MarkdownChunkAlias) {
	if c == nil {
		return
	}

	cacheKey := buildMarkdownCacheKey(projectID, sourceURL)

	c.mu.Lock()
	document := cachedMarkdownDocument{chunks: cloneChunks(chunks), expiresAt: time.Now().Add(c.ttl)}
	c.entries[cacheKey] = document
	c.mu.Unlock()

	c.persist(cacheKey, document)
}

func buildMarkdownCacheKey(projectID string, sourceURL string) string {
	return strings.TrimSpace(projectID) + "|" + strings.TrimSpace(sourceURL)
}

func (c *MarkdownCache) cacheFilePath(cacheKey string) string {
	if c == nil || c.baseDir == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(cacheKey))
	return filepath.Join(c.baseDir, hex.EncodeToString(sum[:])+".json")
}

func (c *MarkdownCache) persist(cacheKey string, document cachedMarkdownDocument) {
	path := c.cacheFilePath(cacheKey)
	if path == "" {
		return
	}
	payload, err := json.Marshal(persistedMarkdownDocument{
		Chunks:    cloneChunks(document.chunks),
		ExpiresAt: document.expiresAt,
	})
	if err != nil {
		return
	}
	_ = os.WriteFile(path, payload, 0o644)
}

func (c *MarkdownCache) getPersisted(cacheKey string, allowExpired bool) ([]services.MarkdownChunkAlias, bool) {
	path := c.cacheFilePath(cacheKey)
	if path == "" {
		return nil, false
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var payload persistedMarkdownDocument
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, false
	}
	if !allowExpired && time.Now().After(payload.ExpiresAt) {
		return nil, false
	}
	if len(payload.Chunks) == 0 {
		return nil, false
	}

	document := cachedMarkdownDocument{chunks: cloneChunks(payload.Chunks), expiresAt: payload.ExpiresAt}
	c.mu.Lock()
	c.entries[cacheKey] = document
	c.mu.Unlock()

	return cloneChunks(payload.Chunks), true
}

func cloneChunks(chunks []services.MarkdownChunkAlias) []services.MarkdownChunkAlias {
	if len(chunks) == 0 {
		return nil
	}
	cloned := make([]services.MarkdownChunkAlias, 0, len(chunks))
	for _, chunk := range chunks {
		cloned = append(cloned, services.MarkdownChunkAlias{Heading: chunk.Heading, Body: chunk.Body})
	}
	return cloned
}
