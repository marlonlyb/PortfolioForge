package model

import (
	"strings"
	"testing"
)

func TestNormalizeIndustryTypeInputSupportsLegacyKeysAndEditorialText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "legacy key", input: "metalworking", want: "metalurgia"},
		{name: "legacy key trimmed", input: "  MATERIAL-HANDLING  ", want: "movimiento de materiales"},
		{name: "editorial text", input: "  automatización   industrial  ", want: "automatización industrial"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeIndustryTypeInput(tt.input); got != tt.want {
				t.Fatalf("NormalizeIndustryTypeInput() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateIndustryTypeRejectsEmptyAndOverLimit(t *testing.T) {
	if err := ValidateIndustryType("   "); err == nil || err.Error() != "industry_type is required" {
		t.Fatalf("empty industry_type error = %v", err)
	}

	overLimit := strings.Repeat("a", IndustryTypeMaxLength+1)
	if err := ValidateIndustryType(overLimit); err == nil || err.Error() != "industry_type exceeds max length 160" {
		t.Fatalf("over-limit industry_type error = %v", err)
	}
}

func TestValidateFinalProductNormalizesWhitespace(t *testing.T) {
	value := "  Panel   HMI para   diagnóstico y monitoreo  "
	if err := ValidateFinalProduct(value); err != nil {
		t.Fatalf("ValidateFinalProduct() error = %v", err)
	}
	if got := NormalizeFinalProduct(value); got != "Panel HMI para diagnóstico y monitoreo" {
		t.Fatalf("NormalizeFinalProduct() = %q", got)
	}
}
