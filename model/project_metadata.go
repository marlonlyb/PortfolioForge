package model

import (
	"fmt"
	"strings"
)

const (
	IndustryTypeMaxLength = 160
	FinalProductMaxLength = 160
)

var legacyIndustryTypeLabels = map[string]string{
	"food":                "alimentación",
	"beverages":           "bebidas",
	"construction":        "construcción",
	"plastics":            "plásticos",
	"cardboard":           "cartón",
	"metalworking":        "metalurgia",
	"material-handling":   "movimiento de materiales",
	"industrial-services": "servicios industriales",
	"other":               "otras industrias",
}

func NormalizeShortEditorialText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func NormalizeIndustryTypeInput(value string) string {
	normalized := NormalizeShortEditorialText(value)
	if normalized == "" {
		return ""
	}
	if legacy, ok := legacyIndustryTypeLabels[strings.ToLower(normalized)]; ok {
		return legacy
	}
	return normalized
}

func NormalizeIndustryTypeKey(value string) string {
	return NormalizeIndustryTypeInput(value)
}

func ValidateIndustryType(value string) error {
	normalized := NormalizeIndustryTypeInput(value)
	if normalized == "" {
		return fmt.Errorf("industry_type is required")
	}
	if len([]rune(normalized)) > IndustryTypeMaxLength {
		return fmt.Errorf("industry_type exceeds max length %d", IndustryTypeMaxLength)
	}
	return nil
}

func NormalizeFinalProduct(value string) string {
	return NormalizeShortEditorialText(value)
}

func ValidateFinalProduct(value string) error {
	normalized := NormalizeFinalProduct(value)
	if normalized == "" {
		return fmt.Errorf("final_product is required")
	}
	if len([]rune(normalized)) > FinalProductMaxLength {
		return fmt.Errorf("final_product exceeds max length %d", FinalProductMaxLength)
	}
	return nil
}
