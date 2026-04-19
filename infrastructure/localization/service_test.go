package localization

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type stubLocalizationTranslator struct {
	translate func(ctx context.Context, sourceLocale string, targetLocale string, fields map[string]json.RawMessage) (map[string]json.RawMessage, error)
}

func (s *stubLocalizationTranslator) TranslateFields(ctx context.Context, sourceLocale string, targetLocale string, fields map[string]json.RawMessage) (map[string]json.RawMessage, error) {
	if s.translate != nil {
		return s.translate(ctx, sourceLocale, targetLocale, fields)
	}
	return fields, nil
}

type stubLocalizationRepo struct {
	rowsByProject map[uuid.UUID][]model.ProjectLocalization
	autoUpserts   []stubLocalizationUpsert
}

type stubLocalizationUpsert struct {
	projectID    uuid.UUID
	locale       string
	fields       map[string]json.RawMessage
	sourceHashes map[string]string
}

func (s *stubLocalizationRepo) ListByProjectID(_ context.Context, projectID uuid.UUID) ([]model.ProjectLocalization, error) {
	return append([]model.ProjectLocalization(nil), s.rowsByProject[projectID]...), nil
}

func (s *stubLocalizationRepo) ListByProjectIDsAndLocale(_ context.Context, projectIDs []uuid.UUID, locale string) (map[uuid.UUID][]model.ProjectLocalization, error) {
	rows := map[uuid.UUID][]model.ProjectLocalization{}
	for _, projectID := range projectIDs {
		for _, row := range s.rowsByProject[projectID] {
			if row.Locale == locale {
				rows[projectID] = append(rows[projectID], row)
			}
		}
	}
	return rows, nil
}

func (s *stubLocalizationRepo) UpsertAuto(_ context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage, sourceHashes map[string]string) error {
	clonedFields := map[string]json.RawMessage{}
	for key, value := range fields {
		clonedFields[key] = append(json.RawMessage(nil), value...)
	}
	clonedHashes := map[string]string{}
	for key, value := range sourceHashes {
		clonedHashes[key] = value
	}
	s.autoUpserts = append(s.autoUpserts, stubLocalizationUpsert{projectID: projectID, locale: locale, fields: clonedFields, sourceHashes: clonedHashes})
	return nil
}

func (s *stubLocalizationRepo) UpsertManual(context.Context, uuid.UUID, string, map[string]json.RawMessage) error {
	return nil
}

func TestBuildProjectFieldMapIncludesClientNameWithoutProfile(t *testing.T) {
	fields := BuildProjectFieldMap(model.Project{
		Name:         "Proyecto base",
		Description:  "Descripción base",
		Category:     "automation",
		ClientName:   "Cliente base",
		IndustryType: "automatización industrial",
		FinalProduct: "Panel HMI",
	})

	if got := string(fields["client_name"]); got != `"Cliente base"` {
		t.Fatalf("client_name = %s", got)
	}
	if got := string(fields["business_goal"]); got != `""` {
		t.Fatalf("business_goal default = %s", got)
	}
	if got := string(fields["final_product"]); got != `"Panel HMI"` {
		t.Fatalf("final_product = %s", got)
	}
	if got := string(fields["industry_type"]); got != `"automatización industrial"` {
		t.Fatalf("industry_type = %s", got)
	}
}

func TestLocalizeProjectAppliesClientNameFallback(t *testing.T) {
	projectID := uuid.New()
	service := NewService(&stubLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
		projectID: {
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "client_name", Value: json.RawMessage(`"Acme Industries"`)},
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "industry_type", Value: json.RawMessage(`"industrial automation"`)},
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "final_product", Value: json.RawMessage(`"Operator dashboard"`)},
		},
	}}, nil)

	localized, err := service.LocalizeProject(context.Background(), model.Project{ID: projectID, ClientName: "Cliente base", IndustryType: "automatización industrial", FinalProduct: "Panel HMI"}, model.LocaleEN)
	if err != nil {
		t.Fatalf("LocalizeProject() error = %v", err)
	}
	if localized.ClientName != "Acme Industries" {
		t.Fatalf("ClientName = %q", localized.ClientName)
	}
	if localized.IndustryType != "industrial automation" {
		t.Fatalf("IndustryType = %q", localized.IndustryType)
	}
	if localized.FinalProduct != "Operator dashboard" {
		t.Fatalf("FinalProduct = %q", localized.FinalProduct)
	}

	fallback, err := service.LocalizeProject(context.Background(), model.Project{ID: projectID, ClientName: "Cliente base"}, model.LocaleDE)
	if err != nil {
		t.Fatalf("LocalizeProject() fallback error = %v", err)
	}
	if fallback.ClientName != "Cliente base" {
		t.Fatalf("fallback ClientName = %q", fallback.ClientName)
	}
}

func TestLocalizeProjectFallsBackToSpanishIndustryTypeWhenLocalizationMissing(t *testing.T) {
	projectID := uuid.New()
	service := NewService(&stubLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
		projectID: {
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "industry_type", Value: json.RawMessage(`"industrial automation"`)},
		},
	}}, nil)

	localized, err := service.LocalizeProject(context.Background(), model.Project{
		ID:           projectID,
		IndustryType: "metalurgia",
		FinalProduct: "Panel HMI",
	}, model.LocaleDE)
	if err != nil {
		t.Fatalf("LocalizeProject() error = %v", err)
	}
	if localized.IndustryType != "metalurgia" {
		t.Fatalf("IndustryType fallback = %q", localized.IndustryType)
	}
	if localized.FinalProduct != "Panel HMI" {
		t.Fatalf("FinalProduct fallback = %q", localized.FinalProduct)
	}
}

func TestLocalizeSearchResultsAppliesClientName(t *testing.T) {
	projectID := uuid.New()
	service := NewService(&stubLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
		projectID: {
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "client_name", Value: json.RawMessage(`"Acme Industries"`)},
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "industry_type", Value: json.RawMessage(`"industrial automation"`)},
		},
	}}, nil)

	baseClient := "Cliente base"
	baseIndustry := "automatización industrial"
	localized, err := service.LocalizeSearchResults(context.Background(), []model.SearchResultItem{{
		ID:           projectID.String(),
		Title:        "Proyecto",
		ClientName:   &baseClient,
		IndustryType: &baseIndustry,
		Evidence: []model.EvidenceField{{
			Field:       "industry_type",
			MatchedText: "automatización industrial",
			MatchType:   model.MatchTypeStructured,
			Score:       1,
		}},
	}}, model.LocaleEN)
	if err != nil {
		t.Fatalf("LocalizeSearchResults() error = %v", err)
	}
	if localized[0].ClientName == nil || *localized[0].ClientName != "Acme Industries" {
		t.Fatalf("ClientName = %#v", localized[0].ClientName)
	}
	if localized[0].IndustryType == nil || *localized[0].IndustryType != "industrial automation" {
		t.Fatalf("IndustryType = %#v", localized[0].IndustryType)
	}
	if localized[0].Evidence[0].MatchedText != "industrial automation" {
		t.Fatalf("Evidence = %#v", localized[0].Evidence)
	}
}

func TestBuildAdminTranslationsResponseExposesClientNameContract(t *testing.T) {
	projectID := uuid.New()
	service := NewService(&stubLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
		projectID: {
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "client_name", Mode: model.LocalizationModeManual, Value: json.RawMessage(`"Analytical Engines"`)},
		},
	}}, nil)

	response, err := service.BuildAdminTranslationsResponse(context.Background(), model.Project{
		ID:         projectID,
		Name:       "Proyecto base",
		ClientName: "Cliente base",
	})
	if err != nil {
		t.Fatalf("BuildAdminTranslationsResponse() error = %v", err)
	}

	if got := string(response.Base["client_name"]); got != `"Cliente base"` {
		t.Fatalf("base client_name = %s", got)
	}

	for _, locale := range model.SupportedTranslationLocales {
		localeState, ok := response.Locales[locale]
		if !ok {
			t.Fatalf("locale %s missing from admin response", locale)
		}

		fieldState, ok := localeState.Fields["client_name"]
		if !ok {
			t.Fatalf("client_name missing from locale %s fields", locale)
		}

		if locale == model.LocaleEN {
			if got := string(fieldState.Value); got != `"Analytical Engines"` {
				t.Fatalf("EN client_name = %s", got)
			}
			if fieldState.Mode != model.LocalizationModeManual {
				t.Fatalf("EN client_name mode = %s", fieldState.Mode)
			}
			continue
		}

		if got := string(fieldState.Value); got != `"Cliente base"` {
			t.Fatalf("%s fallback client_name = %s", locale, got)
		}
		if fieldState.Mode != model.LocalizationModeAuto {
			t.Fatalf("%s fallback client_name mode = %s", locale, fieldState.Mode)
		}
	}
}

func TestRegenerateFromSpanishPreservesManualRowsAndBackfillsClientName(t *testing.T) {
	projectID := uuid.New()
	repo := &stubLocalizationRepo{rowsByProject: map[uuid.UUID][]model.ProjectLocalization{
		projectID: {
			{ProjectID: projectID, Locale: model.LocaleEN, FieldKey: "client_name", Mode: model.LocalizationModeManual, Value: json.RawMessage(`"Manual EN"`)},
			{ProjectID: projectID, Locale: model.LocaleCA, FieldKey: "name", Mode: model.LocalizationModeAuto, Value: json.RawMessage(`"Nom antic"`)},
		},
	}}
	translator := &stubLocalizationTranslator{translate: func(_ context.Context, sourceLocale string, targetLocale string, fields map[string]json.RawMessage) (map[string]json.RawMessage, error) {
		translated := map[string]json.RawMessage{}
		for key, value := range fields {
			translated[key] = value
		}
		if targetLocale == model.LocaleEN {
			translated["name"] = json.RawMessage(`"Project EN"`)
			if _, ok := fields["client_name"]; ok {
				translated["client_name"] = json.RawMessage(`"Auto EN"`)
			}
		}
		if targetLocale == model.LocaleCA {
			translated["name"] = json.RawMessage(`"Projecte CA"`)
			if _, ok := fields["client_name"]; ok {
				translated["client_name"] = json.RawMessage(`"Client CA"`)
			}
		}
		return translated, nil
	}}

	service := NewService(repo, translator)
	err := service.RegenerateFromSpanish(context.Background(), projectID, BuildProjectFieldMap(model.Project{
		Name:        "Proyecto",
		Description: "Descripción",
		Category:    "automation",
		ClientName:  "Cliente base",
	}), []string{model.LocaleCA, model.LocaleEN})
	if err != nil {
		t.Fatalf("RegenerateFromSpanish() error = %v", err)
	}
	if len(repo.autoUpserts) != 2 {
		t.Fatalf("auto upserts = %d, want 2", len(repo.autoUpserts))
	}
	for _, upsert := range repo.autoUpserts {
		if upsert.locale == model.LocaleEN {
			if _, ok := upsert.fields["client_name"]; ok {
				t.Fatalf("manual client_name should not be overwritten in EN: %#v", upsert.fields)
			}
			if got := string(upsert.fields["name"]); got != `"Project EN"` {
				t.Fatalf("EN name = %s", got)
			}
		}
		if upsert.locale == model.LocaleCA {
			if got := string(upsert.fields["client_name"]); got != `"Client CA"` {
				t.Fatalf("CA client_name = %s", got)
			}
		}
	}
}
