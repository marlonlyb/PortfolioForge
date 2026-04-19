package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	"github.com/marlonlyb/portfolioforge/domain/ports/project"
	"github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/model"
)

// Search provides hybrid retrieval search over projects.
type Search struct {
	searchRepo      search.SearchRepository
	projectRepo     project.ProjectReader
	embeddingProv   embedding.EmbeddingProvider
	explanationProv search.ExplanationProvider
	semanticEnabled bool
	constants       model.ScoreConstants
}

// NewSearch creates a new Search service.
func NewSearch(
	sr search.SearchRepository,
	pr project.ProjectReader,
	ep embedding.EmbeddingProvider,
	exp search.ExplanationProvider,
	semanticEnabled bool,
) *Search {
	var constants model.ScoreConstants
	if semanticEnabled {
		constants = model.DefaultScoreConstants()
	} else {
		constants = model.NoSemanticScoreConstants()
	}

	return &Search{
		searchRepo:      sr,
		projectRepo:     pr,
		embeddingProv:   ep,
		explanationProv: exp,
		semanticEnabled: semanticEnabled,
		constants:       constants,
	}
}

// Search performs a hybrid search: normalize → query → fuse scores → threshold → evidence → explain.
func (s *Search) Search(ctx context.Context, params model.SearchParams) (model.SearchResponse, error) {
	// 1. Normalize query
	normalizedQuery := NormalizeQuery(params.Query)
	params.Query = normalizedQuery
	if s.semanticEnabled && normalizedQuery != "" && s.embeddingProv != nil {
		embeddingVec, err := s.embeddingProv.Generate(ctx, normalizedQuery)
		if err != nil {
			return model.SearchResponse{}, fmt.Errorf("services.Search.Search generate embedding: %w", err)
		}
		params.QueryEmbedding = embeddingVec
	}

	// 2. Call SearchRepository (CTE query or all-published)
	results, err := s.searchRepo.Search(ctx, params)
	if err != nil {
		return model.SearchResponse{}, fmt.Errorf("services.Search.Search: %w", err)
	}

	// 3. Apply structured boosts and fuse scores
	results = s.fuseAndBoost(results, normalizedQuery)

	// 4. Apply threshold (skip for empty query — all published)
	if normalizedQuery != "" {
		results = s.applyThreshold(results)
	}

	total := len(results)
	pagedResults, nextCursor := paginateResults(results, params)

	// 5. For each result, extract evidence and generate explanation
	for i := range pagedResults {
		// Fetch technologies for evidence extraction
		techs, _ := s.projectRepo.GetTechnologiesByProjectID(ctx, pagedResults[i].Project.ID)
		pagedResults[i].Project.Technologies = techs

		// Extract evidence
		pagedResults[i].Evidence = s.extractEvidence(pagedResults[i], normalizedQuery, techs)

		// Generate explanation
		explanation, err := s.explanationProv.Explain(ctx, pagedResults[i].Project, pagedResults[i].Evidence, normalizedQuery)
		if err == nil && explanation != "" {
			pagedResults[i].Explanation = explanation
		}
	}

	// 6. Build response
	return s.buildResponse(pagedResults, params, total, nextCursor), nil
}

func paginateResults(results []model.SearchResult, params model.SearchParams) ([]model.SearchResult, *string) {
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	start := 0
	if params.Cursor != "" {
		parsed, err := strconv.Atoi(params.Cursor)
		if err == nil && parsed > 0 {
			start = parsed
		}
	}
	if start >= len(results) {
		return []model.SearchResult{}, nil
	}

	end := start + pageSize
	if end > len(results) {
		end = len(results)
	}

	var nextCursor *string
	if end < len(results) {
		cursor := strconv.Itoa(end)
		nextCursor = &cursor
	}

	return results[start:end], nextCursor
}

// fuseAndBoost normalizes raw scores and applies structured boosts.
func (s *Search) fuseAndBoost(results []model.SearchResult, query string) []model.SearchResult {
	if len(results) == 0 {
		return results
	}

	// Find max lexical score for normalization (floor at 0.1)
	maxLexical := 0.1
	for _, r := range results {
		if r.LexicalScore > maxLexical {
			maxLexical = r.LexicalScore
		}
	}

	for i := range results {
		r := &results[i]

		// Normalize per-signal
		lexicalNorm := r.LexicalScore / maxLexical
		fuzzyNorm := r.FuzzyScore       // already 0-1
		semanticNorm := r.SemanticScore // already 0-1

		// Clamp
		if lexicalNorm > 1.0 {
			lexicalNorm = 1.0
		}
		if fuzzyNorm > 1.0 {
			fuzzyNorm = 1.0
		}
		if semanticNorm > 1.0 {
			semanticNorm = 1.0
		}

		// Weighted fusion
		fusedScore := s.constants.LexicalWeight*lexicalNorm +
			s.constants.FuzzyWeight*fuzzyNorm +
			s.constants.SemanticWeight*semanticNorm

		// Structured boosts
		boost := 0.0
		if query != "" {
			normalizedQuery := strings.ToLower(query)
			if strings.ToLower(r.Project.Category) == normalizedQuery {
				boost += s.constants.CategoryBoost
			}
			if strings.ToLower(r.Project.IndustryType) == normalizedQuery {
				boost += s.constants.CategoryBoost
			}
			if strings.ToLower(r.Project.FinalProduct) == normalizedQuery {
				boost += s.constants.TechBoost
			}
			if strings.ToLower(r.Project.ClientName) == normalizedQuery {
				boost += s.constants.TechBoost
			}
		}

		// Cap boost
		if boost > s.constants.MaxBoost {
			boost = s.constants.MaxBoost
		}

		r.Score = fusedScore + boost
	}

	return results
}

// applyThreshold removes results below the minimum score threshold.
func (s *Search) applyThreshold(results []model.SearchResult) []model.SearchResult {
	filtered := make([]model.SearchResult, 0, len(results))
	for _, r := range results {
		if r.Score >= s.constants.Threshold {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// extractEvidence checks which fields of a project contain the normalized query terms.
func (s *Search) extractEvidence(result model.SearchResult, query string, techs []model.Technology) model.EvidenceTrace {
	if query == "" {
		return model.EvidenceTrace{ProjectID: result.Project.ID}
	}

	terms := strings.Fields(query)
	evidence := model.EvidenceTrace{
		ProjectID: result.Project.ID,
		Fields:    []model.EvidenceField{},
	}

	// Check text fields
	fieldTexts := map[string]string{
		"name":          result.Project.Name,
		"description":   result.Project.Description,
		"client_name":   result.Project.ClientName,
		"category":      result.Project.Category,
		"industry_type": result.Project.IndustryType,
		"final_product": result.Project.FinalProduct,
	}

	if result.Project.Profile != nil {
		fieldTexts["solution_summary"] = result.Project.Profile.SolutionSummary
		fieldTexts["architecture"] = result.Project.Profile.Architecture
		fieldTexts["business_goal"] = result.Project.Profile.BusinessGoal
		fieldTexts["ai_usage"] = result.Project.Profile.AIUsage
	}

	for field, text := range fieldTexts {
		matchedTerms := findMatchingTerms(text, terms)
		if len(matchedTerms) > 0 {
			matchType := model.MatchTypeFTS
			score := result.LexicalScore
			if result.SemanticScore > 0 && result.LexicalScore == 0 && result.FuzzyScore == 0 {
				matchType = model.MatchTypeSemantic
				score = result.SemanticScore
			} else if result.LexicalScore == 0 && result.FuzzyScore > 0 {
				matchType = model.MatchTypeFuzzy
				score = result.FuzzyScore
			}
			evidence.Fields = append(evidence.Fields, model.EvidenceField{
				Field:       field,
				MatchedText: extractSnippet(text, matchedTerms[0], 80),
				MatchType:   matchType,
				Score:       score,
			})
		}
	}

	// Check technology names/slugs
	var matchedTechTerms []string
	for _, t := range techs {
		lowerName := strings.ToLower(t.Name)
		lowerSlug := strings.ToLower(t.Slug)
		for _, term := range terms {
			if strings.Contains(lowerName, term) || strings.Contains(lowerSlug, term) {
				matchedTechTerms = append(matchedTechTerms, term)
			}
		}
	}
	if len(matchedTechTerms) > 0 {
		seen := map[string]bool{}
		unique := []string{}
		for _, t := range matchedTechTerms {
			if !seen[t] {
				seen[t] = true
				unique = append(unique, t)
			}
		}
		evidence.Fields = append(evidence.Fields, model.EvidenceField{
			Field:       "technology",
			MatchedText: strings.Join(unique, ", "),
			MatchType:   model.MatchTypeStructured,
			Score:       result.LexicalScore,
		})
	}

	return evidence
}

// findMatchingTerms returns which terms appear in the text.
func findMatchingTerms(text string, terms []string) []string {
	lowerText := strings.ToLower(text)
	var matched []string
	for _, term := range terms {
		if strings.Contains(lowerText, term) {
			matched = append(matched, term)
		}
	}
	return matched
}

// extractSnippet returns a short snippet of text around the first occurrence of term.
func extractSnippet(text string, term string, maxLen int) string {
	lowerText := strings.ToLower(text)
	idx := strings.Index(lowerText, term)
	if idx == -1 {
		if len(text) > maxLen {
			return text[:maxLen] + "..."
		}
		return text
	}

	start := idx - maxLen/2
	if start < 0 {
		start = 0
	}
	end := start + maxLen
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}
	return snippet
}

// getHeroImage extracts the first image from a project's images JSON array.
func getHeroImage(images json.RawMessage) *string {
	if len(images) == 0 {
		return nil
	}
	var imgList []string
	if err := json.Unmarshal(images, &imgList); err != nil {
		return nil
	}
	if len(imgList) == 0 {
		return nil
	}
	return &imgList[0]
}

// buildResponse converts search results into the API response format.
func (s *Search) buildResponse(results []model.SearchResult, params model.SearchParams, total int, cursor *string) model.SearchResponse {
	items := make([]model.SearchResultItem, 0, len(results))

	for _, r := range results {
		item := model.SearchResultItem{
			ID:           r.Project.ID.String(),
			Slug:         r.Project.Slug,
			Title:        r.Project.Name,
			Category:     r.Project.Category,
			Score:        math.Round(r.Score*1000) / 1000, // 3 decimal places
			Evidence:     r.Evidence.Fields,
			Technologies: []model.TechnologySummary{},
		}

		if r.Project.ClientName != "" {
			item.ClientName = &r.Project.ClientName
		}
		if r.Project.IndustryType != "" {
			item.IndustryType = &r.Project.IndustryType
		}
		if r.Project.FinalProduct != "" {
			item.FinalProduct = &r.Project.FinalProduct
		}

		if r.Project.Profile != nil && r.Project.Profile.SolutionSummary != "" {
			item.Summary = &r.Project.Profile.SolutionSummary
		}

		for _, technology := range r.Project.Technologies {
			tech := model.TechnologySummary{
				ID:   technology.ID.String(),
				Name: technology.Name,
				Slug: technology.Slug,
			}
			if strings.TrimSpace(technology.Color) != "" {
				tech.Color = &technology.Color
			}
			item.Technologies = append(item.Technologies, tech)
		}

		// Hero image
		if heroImg := getHeroImage(r.Project.Images); heroImg != nil {
			item.HeroImage = heroImg
		}

		// Explanation
		if r.Explanation != "" {
			item.Explanation = &r.Explanation
		}

		items = append(items, item)
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	var category *string
	if params.Category != "" {
		category = &params.Category
	}
	var client *string
	if params.Client != "" {
		client = &params.Client
	}

	return model.SearchResponse{
		Data: items,
		Meta: model.SearchMeta{
			Total:    total,
			PageSize: pageSize,
			Cursor:   cursor,
			Query:    params.Query,
			FiltersApplied: model.FiltersApplied{
				Category:     category,
				Client:       client,
				Technologies: params.Technologies,
			},
		},
	}
}
