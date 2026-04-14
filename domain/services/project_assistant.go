package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/marlonlyb/portfolioforge/model"
)

var (
	ErrAssistantInvalidQuestion = errors.New("assistant question is invalid")
	ErrAssistantProjectNotFound = errors.New("assistant project not found")
	ErrAssistantUnavailable     = errors.New("assistant unavailable")
	ErrAssistantUpstream        = errors.New("assistant upstream failure")
)

type projectAssistantRepository interface {
	GetAssistantContextBySlug(ctx context.Context, slug string) (model.ProjectAssistantContext, error)
}

type markdownChunk struct {
	Heading string
	Body    string
}

type MarkdownChunkAlias = markdownChunk

type projectAssistantRetriever interface {
	Fetch(ctx context.Context, projectID string, sourceURL string) ([]markdownChunk, error)
}

type projectAssistantProvider interface {
	GenerateAnswer(ctx context.Context, input ProjectAssistantAnswerInput) (string, error)
}

type ProjectAssistantAnswerInput struct {
	ProjectName string
	Language    string
	Question    string
	History     []model.ProjectAssistantMessage
	Sections    []markdownChunk
}

type ProjectAssistant struct {
	repo      projectAssistantRepository
	retriever projectAssistantRetriever
	provider  projectAssistantProvider
}

func NewProjectAssistant(repo projectAssistantRepository, retriever projectAssistantRetriever, provider projectAssistantProvider) ProjectAssistant {
	return ProjectAssistant{repo: repo, retriever: retriever, provider: provider}
}

func (s ProjectAssistant) Answer(ctx context.Context, slug string, request model.ProjectAssistantRequest) (model.ProjectAssistantResponse, error) {
	question := strings.TrimSpace(request.Question)
	if slug = strings.TrimSpace(slug); slug == "" {
		return model.ProjectAssistantResponse{}, fmt.Errorf("%w: slug is required", ErrAssistantInvalidQuestion)
	}
	if len([]rune(question)) < 2 || len([]rune(question)) > 2000 {
		return model.ProjectAssistantResponse{}, fmt.Errorf("%w: question length", ErrAssistantInvalidQuestion)
	}

	project, err := s.repo.GetAssistantContextBySlug(ctx, slug)
	if err != nil {
		return model.ProjectAssistantResponse{}, fmt.Errorf("%w: %v", ErrAssistantProjectNotFound, err)
	}
	if !project.Active {
		return model.ProjectAssistantResponse{}, ErrAssistantProjectNotFound
	}
	if strings.TrimSpace(project.SourceMarkdownURL) == "" {
		return model.ProjectAssistantResponse{}, ErrAssistantUnavailable
	}

	chunks, err := s.retriever.Fetch(ctx, project.ID.String(), project.SourceMarkdownURL)
	if err != nil {
		return model.ProjectAssistantResponse{}, fmt.Errorf("%w: %v", ErrAssistantUpstream, err)
	}
	if len(chunks) == 0 {
		return model.ProjectAssistantResponse{}, fmt.Errorf("%w: markdown without sections", ErrAssistantUpstream)
	}

	answer, err := s.provider.GenerateAnswer(ctx, ProjectAssistantAnswerInput{
		ProjectName: project.Name,
		Language:    normalizeAssistantLanguage(request.Lang),
		Question:    question,
		History:     normalizeAssistantHistory(request.History),
		Sections:    selectRelevantChunks(question, request.History, chunks),
	})
	if err != nil {
		if errors.Is(err, ErrAssistantUnavailable) {
			return model.ProjectAssistantResponse{}, err
		}
		return model.ProjectAssistantResponse{}, fmt.Errorf("%w: %v", ErrAssistantUpstream, err)
	}

	return model.ProjectAssistantResponse{Answer: strings.TrimSpace(answer)}, nil
}

func normalizeAssistantLanguage(raw string) string {
	lang := strings.ToLower(strings.TrimSpace(raw))
	switch lang {
	case "ca", "catalan":
		return "ca"
	case "en", "english":
		return "en"
	case "de", "german":
		return "de"
	default:
		return "es"
	}
}

func normalizeAssistantHistory(history []model.ProjectAssistantMessage) []model.ProjectAssistantMessage {
	if len(history) == 0 {
		return nil
	}
	start := 0
	if len(history) > 8 {
		start = len(history) - 8
	}
	trimmed := make([]model.ProjectAssistantMessage, 0, len(history)-start)
	for _, item := range history[start:] {
		role := strings.TrimSpace(strings.ToLower(item.Role))
		if role != "assistant" && role != "user" {
			continue
		}
		content := strings.TrimSpace(item.Content)
		if content == "" {
			continue
		}
		trimmed = append(trimmed, model.ProjectAssistantMessage{Role: role, Content: content})
	}
	return trimmed
}

func selectRelevantChunks(question string, history []model.ProjectAssistantMessage, chunks []markdownChunk) []markdownChunk {
	queryTerms := tokenizeAssistantText(question)
	for _, item := range history {
		queryTerms = append(queryTerms, tokenizeAssistantText(item.Content)...)
	}
	if len(queryTerms) == 0 {
		if len(chunks) > 4 {
			return chunks[:4]
		}
		return chunks
	}

	type scoredChunk struct {
		chunk markdownChunk
		score int
	}

	scored := make([]scoredChunk, 0, len(chunks))
	for _, chunk := range chunks {
		score := chunkScore(chunk, queryTerms)
		if score == 0 {
			continue
		}
		scored = append(scored, scoredChunk{chunk: chunk, score: score})
	}

	if len(scored) == 0 {
		if len(chunks) > 4 {
			return chunks[:4]
		}
		return chunks
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return len(scored[i].chunk.Body) > len(scored[j].chunk.Body)
		}
		return scored[i].score > scored[j].score
	})

	limit := min(6, len(scored))
	selected := make([]markdownChunk, 0, limit)
	for _, item := range scored[:limit] {
		selected = append(selected, item.chunk)
	}
	return selected
}

func chunkScore(chunk markdownChunk, queryTerms []string) int {
	body := strings.ToLower(chunk.Body)
	heading := strings.ToLower(chunk.Heading)
	score := 0
	for _, term := range queryTerms {
		if term == "" {
			continue
		}
		if strings.Contains(heading, term) {
			score += 3
		}
		if strings.Contains(body, term) {
			score++
		}
	}
	return score
}

func tokenizeAssistantText(value string) []string {
	parts := strings.FieldsFunc(strings.ToLower(value), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	terms := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) < 3 {
			continue
		}
		terms = append(terms, part)
	}
	return terms
}
