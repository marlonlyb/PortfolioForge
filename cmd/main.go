package main

import (
	"log"
	"os"

	"github.com/marlonlyb/portfolioforge/domain/ports/embedding"
	searchPorts "github.com/marlonlyb/portfolioforge/domain/ports/search"
	"github.com/marlonlyb/portfolioforge/domain/services"
	infraEmbedding "github.com/marlonlyb/portfolioforge/infrastructure/embedding"
	"github.com/marlonlyb/portfolioforge/infrastructure/explanation"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
	"github.com/marlonlyb/portfolioforge/infrastructure/localization"
	"github.com/marlonlyb/portfolioforge/infrastructure/postgres"
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

	dbPool, err := NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	uRepository := postgres.NewUser(dbPool)
	uService := services.NewUser(uRepository)
	uHandlers := handlers.NewUser(uService)

	pRepository := postgres.NewProduct(dbPool)
	pService := services.NewProduct(pRepository)

	lService := services.NewLogin(uService)
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
	pHandlers := handlers.NewProduct(pService, projectLocalizationService)
	projHandlers := handlers.NewProjectPublic(projService, projectLocalizationService)

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

	httpServer := NewServer(uHandlers, pHandlers, lHandlers, projHandlers, searchHandlers, searchAdminHandlers, techHandlers, projAdminHandlers, siteSettingsHandlers)
	httpServer.Initialize()

}
