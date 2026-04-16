package main

import (
	"log"
	"os"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	"github.com/marlonlyb/portfolioforge/domain/ports/mailer"
	searchPorts "github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/domain/services"
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

	httpServer := NewServer(uService, uHandlers, verificationHandlers, projectCatalogHandlers, productPublicCompatHandlers, lHandlers, projHandlers, assistantHandlers, searchHandlers, searchAdminHandlers, techHandlers, projAdminHandlers, siteSettingsHandlers)
	httpServer.Initialize()

}
