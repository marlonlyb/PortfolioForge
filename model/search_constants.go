package model

// ScoreConstants holds the weight configuration for search score fusion.
type ScoreConstants struct {
	LexicalWeight  float64
	FuzzyWeight    float64
	SemanticWeight float64
	Threshold      float64
	CategoryBoost  float64
	TechBoost      float64
	MaxBoost       float64
}

// DefaultScoreConstants returns the default scoring configuration with semantic enabled.
func DefaultScoreConstants() ScoreConstants {
	return ScoreConstants{
		LexicalWeight:  0.45,
		FuzzyWeight:    0.25,
		SemanticWeight: 0.30,
		Threshold:      0.10,
		CategoryBoost:  0.15,
		TechBoost:      0.10,
		MaxBoost:       0.20,
	}
}

// NoSemanticScoreConstants returns scoring configuration without semantic search.
func NoSemanticScoreConstants() ScoreConstants {
	return ScoreConstants{
		LexicalWeight:  0.60,
		FuzzyWeight:    0.40,
		SemanticWeight: 0,
		Threshold:      0.10,
		CategoryBoost:  0.15,
		TechBoost:      0.10,
		MaxBoost:       0.20,
	}
}
