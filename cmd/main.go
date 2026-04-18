package main

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	"github.com/marlonlyb/portfolioforge/domain/ports/mailer"
	searchPorts "github.com/marlonlyb/portfolioforge/domain/ports/search"
	workflowPorts "github.com/marlonlyb/portfolioforge/domain/ports/workflow"
	"github.com/marlonlyb/portfolioforge/domain/services"
	"github.com/marlonlyb/portfolioforge/infrastructure/casestudy"
	infraemail "github.com/marlonlyb/portfolioforge/infrastructure/email"
	infraEmbedding "github.com/marlonlyb/portfolioforge/infrastructure/embedding"
	"github.com/marlonlyb/portfolioforge/infrastructure/explanation"
	"github.com/marlonlyb/portfolioforge/infrastructure/googleauth"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/infrastructure/postgres"
	"github.com/marlonlyb/portfolioforge/infrastructure/projectassistant"
)

func main() {

	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "localization-backfill" {
		if err := runLocalizationBackfill(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "canonical-publish" {
		if err := runCanonicalPublish(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		return
	}

	err = validateEnvironments()
	if err != nil {
		log.Fatal(err)
	}

	smtpConfig, smtpEnabled, err := LoadSMTPConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	var verificationMailer mailer.VerificationMailer = infraemail.NewNoopMailer()
	if smtpEnabled {
		verificationMailer = infraemail.NewSMTPMailer(smtpConfig)
	}

	uRepository := postgres.NewUser(dbPool)
	uService := services.NewUser(uRepository, verificationMailer)
	uHandlers := handlers.NewUser(uService)
	verificationHandlers := handlers.NewEmailVerification(uService)

	projectCatalogRepository := postgres.NewProjectCatalogRepository(dbPool)
	projectCatalogService := services.NewProjectCatalog(projectCatalogRepository)
	productPublicCompatService := services.NewPublicProductCompat(projectCatalogRepository)

	googleVerifier := googleauth.NewVerifierFromEnv()
	lService := services.NewLogin(uService, googleVerifier)
	lHandlers := handlers.NewLogin(lService)

	// Project (public read-side)
	projRepository := postgres.NewProjectRepository(dbPool)
	projService := services.NewProject(projRepository)
	siteSettingsRepository := postgres.NewSiteSettingsRepository(dbPool)

	// Technology
	techRepository := postgres.NewTechnologyRepository(dbPool)
	techHandlers := handlers.NewTechnologyHandler(techRepository)

	// Search
	semanticEnabled := IsSemanticSearchEnabled()
	openAIKey := os.Getenv("OPENAI_API_KEY")
	projectLocalizationRepo := postgres.NewProjectLocalizationRepository(dbPool)
	projectLocalizationService := localization.NewService(projectLocalizationRepo, localization.NewOpenAITranslator(openAIKey))
	projectCatalogHandlers := handlers.NewProjectCatalog(projectCatalogService, projectLocalizationService)
	productPublicCompatHandlers := handlers.NewProductPublicCompat(productPublicCompatService)
	projHandlers := handlers.NewProjectPublic(projService, projectLocalizationService)
	markdownCache := projectassistant.NewDefaultMarkdownCache()
	assistantFetcher := projectassistant.NewMarkdownFetcher(markdownCache)
	assistantProvider := projectassistant.NewOpenAIProvider(openAIKey)
	assistantService := services.NewProjectAssistant(projRepository, assistantFetcher, assistantProvider)
	assistantHandlers := handlers.NewProjectAssistantHandler(assistantService)

	var embeddingProv embedding.EmbeddingProvider = infraEmbedding.NewNoOpEmbeddingProvider()
	var explProv searchPorts.ExplanationProvider = explanation.NewTemplateExplanationProvider()
	semanticDegraded := semanticEnabled

	if openAIKey != "" {
		embeddingProv = infraEmbedding.NewOpenAIEmbeddingProvider(openAIKey)
		explProv = explanation.NewOpenAIExplanationProvider(openAIKey, explanation.NewTemplateExplanationProvider())
		semanticDegraded = false
	} else if semanticEnabled {
		log.Println("WARNING: ENABLE_SEMANTIC_SEARCH is true but OPENAI_API_KEY is not set. Falling back to NoOp/Template providers.")
	}

	searchRepo := postgres.NewSearchRepository(dbPool, semanticEnabled)
	searchService := services.NewSearch(searchRepo, projRepository, embeddingProv, explProv, semanticEnabled)

	searchHandlers := handlers.NewSearch(searchService, semanticDegraded, projectLocalizationService)

	// Search Admin
	searchAdminHandlers := handlers.NewSearchAdmin(projRepository, searchRepo)

	// Project Admin
	projAdminHandlers := handlers.NewProjectAdminHandler(dbPool, embeddingProv, semanticEnabled, projRepository, projectLocalizationService)
	siteSettingsHandlers := handlers.NewSiteSettingsHandler(siteSettingsRepository)
	workflowEnvConfig := LoadCaseStudyWorkflowEnvConfig()
	workflowHandlers, workflowDisabledReason := buildCaseStudyWorkflowHandler(workflowEnvConfig, caseStudyWorkflowDependencies{
		repository:   postgres.NewCaseStudyWorkflowRepository(dbPool),
		publisher:    workflowPublisherAdapter{config: workflowEnvConfig.Publish},
		importer:     workflowImporterAdapter{importer: casestudy.NewProjectImporter(projectCatalogService)},
		localization: localization.NewBackfillService(projRepository, projectLocalizationService),
		searchRepo:   searchRepo,
	})
	if workflowDisabledReason != "" {
		log.Printf("WARNING: case-study workflow disabled: %s", workflowDisabledReason)
	}

	httpServer := NewServer(uService, uHandlers, verificationHandlers, projectCatalogHandlers, productPublicCompatHandlers, lHandlers, projHandlers, assistantHandlers, searchHandlers, searchAdminHandlers, techHandlers, projAdminHandlers, siteSettingsHandlers, workflowHandlers)
	httpServer.Initialize()

}

type workflowPublisherAdapter struct {
	config casestudy.PublishConfig
}

type caseStudyWorkflowDependencies struct {
	repository   workflowPorts.Repository
	publisher    services.CaseStudyPublisher
	importer     services.CaseStudyProjectImporter
	localization services.CaseStudyLocalizationBackfiller
	searchRepo   searchPorts.SearchRepository
}

func buildCaseStudyWorkflowHandler(config CaseStudyWorkflowEnvConfig, deps caseStudyWorkflowDependencies) (*handlers.CaseStudyWorkflowHandler, string) {
	if !config.Configured {
		return handlers.NewCaseStudyWorkflowHandler(services.NewUnavailableCaseStudyWorkflowService(config.Reason)), config.Diagnostic
	}

	workflowService := services.NewCaseStudyWorkflowService(
		deps.repository,
		deps.publisher,
		deps.importer,
		deps.localization,
		deps.searchRepo,
		services.CaseStudyWorkflowConfig{AllowedSourceRoots: config.AllowedSourceRoots},
	)
	return handlers.NewCaseStudyWorkflowHandler(workflowService), ""
}

func (p workflowPublisherAdapter) ResolvePublishTarget(inputPath, explicitSlug string) (services.CaseStudyPublishTarget, error) {
	target, err := casestudy.ResolvePublishTarget(inputPath, explicitSlug, p.config)
	if err != nil {
		return services.CaseStudyPublishTarget{}, err
	}
	return services.CaseStudyPublishTarget{
		Slug:      target.Slug,
		LocalDir:  target.LocalDir,
		LocalFile: target.LocalFile,
		PublicURL: target.PublicURL,
	}, nil
}

func (p workflowPublisherAdapter) CollectFiles(localDir string) ([]string, error) {
	return casestudy.CollectPublishFiles(localDir)
}

func (p workflowPublisherAdapter) Publish(ctx context.Context, target services.CaseStudyPublishTarget, files []string) error {
	resolved, err := casestudy.ResolvePublishTarget(target.LocalDir, target.Slug, p.config)
	if err != nil {
		return err
	}
	return casestudy.PublishFiles(resolved, files, p.config)
}

func (p workflowPublisherAdapter) Verify(ctx context.Context, publicURL string) error {
	return casestudy.VerifyPublishedCanonical(ctx, publicURL)
}

type workflowImporterAdapter struct {
	importer *casestudy.ProjectImporter
}

func (a workflowImporterAdapter) ImportFromCanonical(ctx context.Context, source services.CaseStudyPublishTarget, canonicalURL string) (uuid.UUID, error) {
	target := casestudy.PublishTarget{
		Slug:      source.Slug,
		LocalDir:  source.LocalDir,
		LocalFile: source.LocalFile,
		PublicURL: source.PublicURL,
	}
	return a.importer.ImportFromCanonical(ctx, target, canonicalURL)
}
