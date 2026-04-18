package model

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

const (
	LocaleES = "es"
	LocaleCA = "ca"
	LocaleEN = "en"
	LocaleDE = "de"

	LocalizationModeAuto   = "auto"
	LocalizationModeManual = "manual"
)

var SupportedTranslationLocales = []string{LocaleCA, LocaleEN, LocaleDE}

var SupportedPublicLocales = []string{LocaleES, LocaleCA, LocaleEN, LocaleDE}

var TranslatableProjectFieldKeys = []string{
	"name",
	"description",
	"category",
	"client_name",
	"business_goal",
	"problem_statement",
	"solution_summary",
	"delivery_scope",
	"responsibility_scope",
	"architecture",
	"ai_usage",
	"integrations",
	"technical_decisions",
	"challenges",
	"results",
	"metrics",
	"timeline",
}

type ProjectLocalization struct {
	ProjectID  uuid.UUID       `json:"project_id"`
	Locale     string          `json:"locale"`
	FieldKey   string          `json:"field_key"`
	Value      json.RawMessage `json:"value"`
	Mode       string          `json:"mode"`
	SourceHash string          `json:"source_hash,omitempty"`
	UpdatedAt  int64           `json:"updated_at,omitempty"`
}

func IsSupportedPublicLocale(locale string) bool {
	locale = strings.ToLower(strings.TrimSpace(locale))
	for _, candidate := range SupportedPublicLocales {
		if candidate == locale {
			return true
		}
	}
	return false
}

func IsSupportedTranslationLocale(locale string) bool {
	locale = strings.ToLower(strings.TrimSpace(locale))
	for _, candidate := range SupportedTranslationLocales {
		if candidate == locale {
			return true
		}
	}
	return false
}
