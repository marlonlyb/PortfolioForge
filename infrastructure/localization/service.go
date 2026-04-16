package localization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/model"
)

type Translator interface {
	TranslateFields(ctx context.Context, sourceLocale string, targetLocale string, fields map[string]json.RawMessage) (map[string]json.RawMessage, error)
}

type Repository interface {
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.ProjectLocalization, error)
	ListByProjectIDsAndLocale(ctx context.Context, projectIDs []uuid.UUID, locale string) (map[uuid.UUID][]model.ProjectLocalization, error)
	UpsertAuto(ctx context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage, sourceHashes map[string]string) error
	UpsertManual(ctx context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage) error
}

type Service struct {
	repo       Repository
	translator Translator
}

type AdminFieldState struct {
	Value json.RawMessage `json:"value"`
	Mode  string          `json:"mode"`
}

type AdminLocaleState struct {
	Locale string                     `json:"locale"`
	Fields map[string]AdminFieldState `json:"fields"`
}

type AdminTranslationsResponse struct {
	ProjectID string                      `json:"project_id"`
	Locales   map[string]AdminLocaleState `json:"locales"`
	Base      map[string]json.RawMessage  `json:"base"`
}

func NewService(repo Repository, translator Translator) *Service {
	return &Service{repo: repo, translator: translator}
}

func NormalizeLocale(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if !model.IsSupportedPublicLocale(locale) {
		return model.LocaleES
	}
	return locale
}

func BuildProjectFieldMap(project model.Project) map[string]json.RawMessage {
	fields := map[string]json.RawMessage{
		"name":        mustMarshal(project.Name),
		"description": mustMarshal(project.Description),
		"category":    mustMarshal(project.Category),
	}

	profile := project.Profile
	if profile == nil {
		fields["business_goal"] = mustMarshal("")
		fields["problem_statement"] = mustMarshal("")
		fields["solution_summary"] = mustMarshal("")
		fields["delivery_scope"] = mustMarshal("")
		fields["responsibility_scope"] = mustMarshal("")
		fields["architecture"] = mustMarshal("")
		fields["ai_usage"] = mustMarshal("")
		fields["integrations"] = mustMarshal([]string{})
		fields["technical_decisions"] = mustMarshal([]string{})
		fields["challenges"] = mustMarshal([]string{})
		fields["results"] = mustMarshal([]string{})
		fields["metrics"] = mustMarshal(map[string]string{})
		fields["timeline"] = mustMarshal([]string{})
		return fields
	}

	fields["business_goal"] = normalizeRaw(profile.BusinessGoal)
	fields["problem_statement"] = normalizeRaw(profile.ProblemStatement)
	fields["solution_summary"] = normalizeRaw(profile.SolutionSummary)
	fields["delivery_scope"] = normalizeRaw(profile.DeliveryScope)
	fields["responsibility_scope"] = normalizeRaw(profile.ResponsibilityScope)
	fields["architecture"] = normalizeRaw(profile.Architecture)
	fields["ai_usage"] = normalizeRaw(profile.AIUsage)
	fields["integrations"] = normalizeJSONOrDefault(profile.Integrations, []string{})
	fields["technical_decisions"] = normalizeJSONOrDefault(profile.TechnicalDecisions, []string{})
	fields["challenges"] = normalizeJSONOrDefault(profile.Challenges, []string{})
	fields["results"] = normalizeJSONOrDefault(profile.Results, []string{})
	fields["metrics"] = normalizeJSONOrDefault(profile.Metrics, map[string]string{})
	fields["timeline"] = normalizeJSONOrDefault(profile.Timeline, []string{})

	return fields
}

func (s *Service) SyncFromSpanish(ctx context.Context, projectID uuid.UUID, previousFields map[string]json.RawMessage, currentFields map[string]json.RawMessage) error {
	if s == nil || s.repo == nil {
		return nil
	}

	changedFields := map[string]json.RawMessage{}
	sourceHashes := map[string]string{}
	for _, fieldKey := range model.TranslatableProjectFieldKeys {
		current := normalizeKnownField(currentFields[fieldKey], fieldKey)
		previous := normalizeKnownField(previousFields[fieldKey], fieldKey)
		if jsonEqual(previous, current) {
			continue
		}
		changedFields[fieldKey] = current
		sourceHashes[fieldKey] = hashRaw(current)
	}

	if len(changedFields) == 0 {
		return nil
	}

	existingRows, err := s.repo.ListByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("load existing localizations: %w", err)
	}
	existing := indexRows(existingRows)

	for _, locale := range model.SupportedTranslationLocales {
		fieldsForLocale := map[string]json.RawMessage{}
		for fieldKey, value := range changedFields {
			row, ok := existing[locale][fieldKey]
			if ok && row.Mode == model.LocalizationModeManual {
				continue
			}
			fieldsForLocale[fieldKey] = value
		}

		if len(fieldsForLocale) == 0 {
			continue
		}

		if s.translator == nil {
			return fmt.Errorf("translation provider not configured")
		}

		translatedFields, err := s.translator.TranslateFields(ctx, model.LocaleES, locale, fieldsForLocale)
		if err != nil {
			return fmt.Errorf("translate fields for %s: %w", locale, err)
		}

		filteredHashes := map[string]string{}
		for fieldKey := range translatedFields {
			filteredHashes[fieldKey] = sourceHashes[fieldKey]
		}

		if err := s.repo.UpsertAuto(ctx, projectID, locale, translatedFields, filteredHashes); err != nil {
			return fmt.Errorf("persist auto localizations for %s: %w", locale, err)
		}
	}

	return nil
}

func (s *Service) SaveManualTranslations(ctx context.Context, projectID uuid.UUID, locale string, fields map[string]json.RawMessage) error {
	if s == nil || s.repo == nil {
		return nil
	}

	locale = NormalizeLocale(locale)
	if !model.IsSupportedTranslationLocale(locale) {
		return fmt.Errorf("unsupported locale: %s", locale)
	}

	normalized := map[string]json.RawMessage{}
	for fieldKey, value := range fields {
		if !containsField(fieldKey) {
			continue
		}
		normalized[fieldKey] = normalizeKnownField(value, fieldKey)
	}

	if len(normalized) == 0 {
		return nil
	}

	return s.repo.UpsertManual(ctx, projectID, locale, normalized)
}

func (s *Service) BuildAdminTranslationsResponse(ctx context.Context, project model.Project) (AdminTranslationsResponse, error) {
	response := AdminTranslationsResponse{
		ProjectID: project.ID.String(),
		Locales:   map[string]AdminLocaleState{},
		Base:      BuildProjectFieldMap(project),
	}

	rows, err := s.repo.ListByProjectID(ctx, project.ID)
	if err != nil {
		return response, err
	}
	indexed := indexRows(rows)

	for _, locale := range model.SupportedTranslationLocales {
		localeFields := map[string]AdminFieldState{}
		for _, fieldKey := range model.TranslatableProjectFieldKeys {
			if row, ok := indexed[locale][fieldKey]; ok {
				localeFields[fieldKey] = AdminFieldState{Value: row.Value, Mode: row.Mode}
				continue
			}
			localeFields[fieldKey] = AdminFieldState{Value: response.Base[fieldKey], Mode: model.LocalizationModeAuto}
		}
		response.Locales[locale] = AdminLocaleState{Locale: locale, Fields: localeFields}
	}

	return response, nil
}

func (s *Service) LocalizeProject(ctx context.Context, project model.Project, locale string) (model.Project, error) {
	localized, err := s.LocalizeProjects(ctx, []model.Project{project}, locale)
	if err != nil {
		return project, err
	}
	if len(localized) == 0 {
		return project, nil
	}
	return localized[0], nil
}

func (s *Service) LocalizeProjects(ctx context.Context, projects []model.Project, locale string) ([]model.Project, error) {
	locale = NormalizeLocale(locale)
	if locale == model.LocaleES || len(projects) == 0 || s == nil || s.repo == nil {
		return projects, nil
	}

	projectIDs := make([]uuid.UUID, 0, len(projects))
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}

	rowsByProject, err := s.repo.ListByProjectIDsAndLocale(ctx, projectIDs, locale)
	if err != nil {
		return nil, err
	}

	localized := make([]model.Project, len(projects))
	copy(localized, projects)
	for index := range localized {
		applyLocalizedFields(&localized[index], rowsByProject[localized[index].ID])
	}

	return localized, nil
}

func (s *Service) LocalizeSearchResults(ctx context.Context, items []model.SearchResultItem, locale string) ([]model.SearchResultItem, error) {
	locale = NormalizeLocale(locale)
	if locale == model.LocaleES || len(items) == 0 || s == nil || s.repo == nil {
		return items, nil
	}

	projectIDs := make([]uuid.UUID, 0, len(items))
	ordered := map[uuid.UUID]int{}
	for index, item := range items {
		projectID, err := uuid.Parse(item.ID)
		if err != nil {
			continue
		}
		projectIDs = append(projectIDs, projectID)
		ordered[projectID] = index
	}

	rowsByProject, err := s.repo.ListByProjectIDsAndLocale(ctx, projectIDs, locale)
	if err != nil {
		return nil, err
	}

	localized := make([]model.SearchResultItem, len(items))
	copy(localized, items)
	for projectID, rows := range rowsByProject {
		index, ok := ordered[projectID]
		if !ok {
			continue
		}
		applyLocalizedSearchResult(&localized[index], rows)
	}

	return localized, nil
}

func containsField(fieldKey string) bool {
	for _, candidate := range model.TranslatableProjectFieldKeys {
		if candidate == fieldKey {
			return true
		}
	}
	return false
}

func indexRows(rows []model.ProjectLocalization) map[string]map[string]model.ProjectLocalization {
	indexed := map[string]map[string]model.ProjectLocalization{}
	for _, row := range rows {
		if _, ok := indexed[row.Locale]; !ok {
			indexed[row.Locale] = map[string]model.ProjectLocalization{}
		}
		indexed[row.Locale][row.FieldKey] = row
	}
	return indexed
}

func applyLocalizedFields(project *model.Project, rows []model.ProjectLocalization) {
	for _, row := range rows {
		switch row.FieldKey {
		case "name":
			project.Name = decodeString(row.Value, project.Name)
		case "description":
			project.Description = decodeString(row.Value, project.Description)
		case "category":
			project.Category = decodeString(row.Value, project.Category)
		case "business_goal":
			ensureProfile(project)
			project.Profile.BusinessGoal = decodeString(row.Value, project.Profile.BusinessGoal)
		case "problem_statement":
			ensureProfile(project)
			project.Profile.ProblemStatement = decodeString(row.Value, project.Profile.ProblemStatement)
		case "solution_summary":
			ensureProfile(project)
			project.Profile.SolutionSummary = decodeString(row.Value, project.Profile.SolutionSummary)
		case "delivery_scope":
			ensureProfile(project)
			project.Profile.DeliveryScope = decodeString(row.Value, project.Profile.DeliveryScope)
		case "responsibility_scope":
			ensureProfile(project)
			project.Profile.ResponsibilityScope = decodeString(row.Value, project.Profile.ResponsibilityScope)
		case "architecture":
			ensureProfile(project)
			project.Profile.Architecture = decodeString(row.Value, project.Profile.Architecture)
		case "ai_usage":
			ensureProfile(project)
			project.Profile.AIUsage = decodeString(row.Value, project.Profile.AIUsage)
		case "integrations":
			ensureProfile(project)
			project.Profile.Integrations = cloneRaw(row.Value)
		case "technical_decisions":
			ensureProfile(project)
			project.Profile.TechnicalDecisions = cloneRaw(row.Value)
		case "challenges":
			ensureProfile(project)
			project.Profile.Challenges = cloneRaw(row.Value)
		case "results":
			ensureProfile(project)
			project.Profile.Results = cloneRaw(row.Value)
		case "metrics":
			ensureProfile(project)
			project.Profile.Metrics = cloneRaw(row.Value)
		case "timeline":
			ensureProfile(project)
			project.Profile.Timeline = cloneRaw(row.Value)
		}
	}
}

func applyLocalizedSearchResult(item *model.SearchResultItem, rows []model.ProjectLocalization) {
	for _, row := range rows {
		switch row.FieldKey {
		case "name":
			item.Title = decodeString(row.Value, item.Title)
		case "description":
			summary := decodeString(row.Value, derefString(item.Summary))
			item.Summary = &summary
		case "category":
			item.Category = decodeString(row.Value, item.Category)
		}
	}
}

func decodeString(raw json.RawMessage, fallback string) string {
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return fallback
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func ensureProfile(project *model.Project) {
	if project.Profile == nil {
		project.Profile = &model.ProjectProfile{ProjectID: project.ID}
	}
}

func normalizeRaw(value string) json.RawMessage {
	return mustMarshal(strings.TrimSpace(value))
}

func normalizeJSONOrDefault(raw json.RawMessage, fallback interface{}) json.RawMessage {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return mustMarshal(fallback)
	}

	var decoded interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return mustMarshal(fallback)
	}
	return mustMarshal(decoded)
}

func normalizeKnownField(raw json.RawMessage, fieldKey string) json.RawMessage {
	switch fieldKey {
	case "integrations", "technical_decisions", "challenges", "results", "timeline":
		return normalizeJSONOrDefault(raw, []string{})
	case "metrics":
		return normalizeJSONOrDefault(raw, map[string]string{})
	default:
		if len(raw) == 0 {
			return mustMarshal("")
		}
		var text string
		if err := json.Unmarshal(raw, &text); err == nil {
			return mustMarshal(strings.TrimSpace(text))
		}
		return normalizeJSONOrDefault(raw, "")
	}
}

func cloneRaw(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	cloned := make([]byte, len(raw))
	copy(cloned, raw)
	return cloned
}

func jsonEqual(left, right json.RawMessage) bool {
	return string(left) == string(right)
}

func hashRaw(raw json.RawMessage) string {
	hash := sha256.Sum256(raw)
	return hex.EncodeToString(hash[:])
}

func mustMarshal(value interface{}) json.RawMessage {
	encoded, _ := json.Marshal(value)
	return encoded
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func SortedFieldKeys() []string {
	keys := append([]string(nil), model.TranslatableProjectFieldKeys...)
	sort.Strings(keys)
	return keys
}
