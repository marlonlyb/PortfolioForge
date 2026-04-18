import type { PublicLocale } from './config';

export interface LandingQuickPrompt {
  label: string;
  query: string;
}

interface Messages {
  headerTitle: string;
  headerSummary: string;
  navHome: string;
  navLogin: string;
  navSearch: string;
  navAdmin: string;
  navLogout: string;
  headerCaption: string;
  searchPlaceholder: string;
  searchButton: string;
  searchClear: string;
  searchSuggestionsLabel: string;
  landingSearchEyebrow: string;
  landingSearchTitle: string;
  landingSearchLead: string;
  landingSearchPlaceholder: string;
  landingSearchContextHint: string;
  landingQuickPrompts: LandingQuickPrompt[];
  landingEyebrow: string;
  landingTitle: string;
  landingLead: string;
  landingPrimaryCta: string;
  landingSecondaryCta: string;
  landingDesignIntent: string;
  landingPrinciples: string[];
  landingHighlights: Array<{ value: string; label: string }>;
  landingPortfolioSystem: string;
  landingShowcaseTitle: string;
  landingShowcaseCopy: string;
  landingQuoteEyebrow: string;
  landingQuote: string;
  catalogEyebrow: string;
  catalogTitle: string;
  catalogIntro: string;
  catalogSearchLabel: string;
  catalogSearchPlaceholder: string;
  catalogCategoryPlaceholder: string;
  catalogClearFilters: string;
  catalogViewCaseStudy: string;
  catalogOpenProject: string;
  catalogNoImage: string;
  catalogNoSummary: string;
  catalogLoading: string;
  catalogError: string;
  catalogNoMatches: string;
  catalogNoFilterMatches: string;
  detailBack: string;
  detailLoading: string;
  detailNotFound: string;
  detailProjectOverview: string;
  detailCategory: string;
  detailClient: string;
  detailUpdated: string;
  detailTechnologies: string;
  detailIndependent: string;
  detailRecentlyCurated: string;
  detailNotSpecified: string;
  detailStrategyLayer: string;
  detailExecutionLayer: string;
  detailTechnicalLayer: string;
  detailBusinessGoal: string;
  detailProblem: string;
  detailSolution: string;
  detailDeliveryScope: string;
  detailResponsibilityScope: string;
  detailArchitecture: string;
  detailAIUsage: string;
  detailIntegrations: string;
  detailTechnicalDecisions: string;
  detailChallenges: string;
  detailResults: string;
  detailMetrics: string;
  detailTimeline: string;
  detailClientContext: string;
  detailVisualUnavailable: string;
  authPublicEyebrow: string;
  authPublicLoginTitle: string;
  authPublicLoginDescription: string;
  authPublicSignupTitle: string;
  authPublicSignupDescription: string;
  authPublicLocalRestriction: string;
  authAdminEyebrow: string;
  authAdminTitle: string;
  authAdminDescription: string;
  authBackToPortfolio: string;
  authBackToPublicLogin: string;
  detailAssistantLocalRestriction: string;
}

type AuthMessageKey =
  | 'authPublicEyebrow'
  | 'authPublicLoginTitle'
  | 'authPublicLoginDescription'
  | 'authPublicSignupTitle'
  | 'authPublicSignupDescription'
  | 'authPublicLocalRestriction'
  | 'authAdminEyebrow'
  | 'authAdminTitle'
  | 'authAdminDescription'
  | 'authBackToPortfolio'
  | 'authBackToPublicLogin'
  | 'detailAssistantLocalRestriction';

type CoreMessages = Omit<Messages, AuthMessageKey>;

const authMessages: Record<PublicLocale, Pick<Messages, AuthMessageKey>> = {
  es: {
    authPublicEyebrow: 'Acceso público',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
  },
  ca: {
    authPublicEyebrow: 'Accés públic',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
  },
  en: {
    authPublicEyebrow: 'Public access',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
  },
  de: {
    authPublicEyebrow: 'Öffentlicher Zugang',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
  },
};

const baseMessages: Record<PublicLocale, CoreMessages> = {
  es: {
    headerTitle: 'Portfolio de proyectos',
    headerSummary: 'Estrategia, ejecución y criterio técnico.',
    navHome: 'Inicio',
    navLogin: 'Login',
    navSearch: 'Buscar',
    navAdmin: 'Admin',
    navLogout: 'Salir',
    headerCaption: 'Ing. Marlon Ly Bellido',
    searchPlaceholder: 'Busca proyectos por tecnología, cliente o concepto…',
    searchButton: 'Buscar',
    searchClear: 'Limpiar búsqueda',
    searchSuggestionsLabel: 'Sugerencias de búsqueda',
    landingSearchEyebrow: 'BÚSQUEDA GUIADA',
    landingSearchTitle: '',
    landingSearchLead: 'Busca proyectos, casos y experiencias reales.',
    landingSearchPlaceholder: 'Busca un proyecto, tecnología o tema...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Muéstrame la migración PLC de Printer 05', query: 'Printer 05' },
      { label: 'Quiero casos con Allen-Bradley y CompactLogix', query: 'CompactLogix' },
      { label: 'Enséñame automatización industrial con Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Busca motion control con SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Trabajo digital seleccionado',
    landingTitle: 'Una portada más editorial para presentar producto, arquitectura y ejecución.',
    landingLead:
      'PortfolioForge reúne proyectos públicos con mejor jerarquía, bloques más definidos y una lectura más clara del valor detrás de cada entrega.',
    landingPrimaryCta: 'Explorar proyectos',
    landingSecondaryCta: 'Ver catálogo',
    landingDesignIntent: 'Intención visual',
    landingPrinciples: [
      'Casos de estudio con narrativa clara y foco en decisiones.',
      'Composición modular pensada para desktop, tablet y móvil.',
      'Paleta dark mantenida con contraste editorial y ritmo visual.',
    ],
    landingHighlights: [
      { value: '01', label: 'Portfolio público con historias de proyecto estructuradas' },
      { value: '02', label: 'Catálogo modular y responsive para trabajo seleccionado' },
      { value: '03', label: 'Búsqueda y detalle alineados con el mismo sistema visual' },
    ],
    landingPortfolioSystem: 'Sistema de portfolio',
    landingShowcaseTitle: 'Una composición modular para leer el portfolio como publicación, no como listado.',
    landingShowcaseCopy:
      'La landing ahora separa mensaje, exploración y catálogo en bloques más amplios para que el contenido respire mejor en desktop sin perder claridad en móvil.',
    landingQuoteEyebrow: 'Ritmo editorial',
    landingQuote: '“Bloques sólidos, mejor uso del ancho y una estructura visual que prioriza contexto, narrativa y exploración.”',
    catalogEyebrow: 'Índice de proyectos',
    catalogTitle: 'Casos de estudio seleccionados',
    catalogIntro: '',
    catalogSearchLabel: 'Buscar proyectos',
    catalogSearchPlaceholder: 'Buscar por nombre del proyecto…',
    catalogCategoryPlaceholder: 'Categoría',
    catalogClearFilters: 'Limpiar filtros',
    catalogViewCaseStudy: 'Ver caso de estudio',
    catalogOpenProject: 'Abrir proyecto',
    catalogNoImage: 'Sin imagen',
    catalogNoSummary: 'Sin resumen disponible.',
    catalogLoading: 'Cargando proyectos…',
    catalogError: 'Error',
    catalogNoMatches: 'No hay proyectos que coincidan con tu búsqueda.',
    catalogNoFilterMatches: 'No hay proyectos para los filtros actuales.',
    detailBack: '← Volver a proyectos',
    detailLoading: 'Cargando proyecto…',
    detailNotFound: 'Proyecto no encontrado.',
    detailProjectOverview: 'Resumen del proyecto',
    detailCategory: 'Categoría',
    detailClient: 'Cliente',
    detailUpdated: 'Actualizado',
    detailTechnologies: 'Tecnologías',
    detailIndependent: 'Independiente / Interno',
    detailRecentlyCurated: 'Curado recientemente',
    detailNotSpecified: 'No especificado',
    detailStrategyLayer: 'Estrategia',
    detailExecutionLayer: 'Ejecución',
    detailTechnicalLayer: 'Técnica',
    detailBusinessGoal: 'Contexto de negocio / Objetivo',
    detailProblem: 'Problema',
    detailSolution: 'Solución',
    detailDeliveryScope: 'Delivery Scope',
    detailResponsibilityScope: 'Responsibility Scope',
    detailArchitecture: 'Arquitectura',
    detailAIUsage: 'Uso de IA',
    detailIntegrations: 'Integraciones',
    detailTechnicalDecisions: 'Decisiones técnicas',
    detailChallenges: 'Desafíos',
    detailResults: 'Resultados',
    detailMetrics: 'Métricas',
    detailTimeline: 'Timeline',
    detailClientContext: 'Contexto cliente',
    detailVisualUnavailable: 'Visual del caso de estudio no disponible',
  },
  ca: {
    headerTitle: 'Portfoli de projectes',
    headerSummary: 'Estratègia, execució i criteri tècnic.',
    navHome: 'Inici',
    navLogin: 'Login',
    navSearch: 'Cercar',
    navAdmin: 'Admin',
    navLogout: 'Sortir',
    headerCaption: 'Eng. Marlon Ly Bellido',
    searchPlaceholder: 'Cerca projectes per tecnologia, client o concepte…',
    searchButton: 'Cercar',
    searchClear: 'Netejar cerca',
    searchSuggestionsLabel: 'Suggeriments de cerca',
    landingSearchEyebrow: 'CERCA GUIADA',
    landingSearchTitle: '',
    landingSearchLead: 'Cerca projectes, casos i experiències reals.',
    landingSearchPlaceholder: 'Cerca un projecte, tecnologia o tema...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Mostra’m la migració PLC de Printer 05', query: 'Printer 05' },
      { label: 'Vull casos amb Allen-Bradley i CompactLogix', query: 'CompactLogix' },
      { label: 'Ensenya’m automatització industrial amb Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Busca motion control amb SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Treball digital seleccionat',
    landingTitle: 'Una portada més editorial per presentar producte, arquitectura i execució.',
    landingLead:
      'PortfolioForge reuneix projectes públics amb millor jerarquia, blocs més definits i una lectura més clara del valor de cada entrega.',
    landingPrimaryCta: 'Explorar projectes',
    landingSecondaryCta: 'Veure catàleg',
    landingDesignIntent: 'Intenció visual',
    landingPrinciples: [
      'Casos d’estudi amb narrativa clara i focus en decisions.',
      'Composició modular pensada per a escriptori, tauleta i mòbil.',
      'Paleta dark mantinguda amb contrast editorial i ritme visual.',
    ],
    landingHighlights: [
      { value: '01', label: 'Portfoli públic amb històries de projecte estructurades' },
      { value: '02', label: 'Catàleg modular i responsive per al treball seleccionat' },
      { value: '03', label: 'Cerca i detall alineats amb el mateix sistema visual' },
    ],
    landingPortfolioSystem: 'Sistema de portfoli',
    landingShowcaseTitle: 'Una composició modular per llegir el portfoli com una publicació, no com un llistat.',
    landingShowcaseCopy:
      'La landing ara separa missatge, exploració i catàleg en blocs més amplis perquè el contingut respiri millor en desktop sense perdre claredat en mòbil.',
    landingQuoteEyebrow: 'Ritme editorial',
    landingQuote: '“Blocs sòlids, millor ús de l’amplada i una estructura visual que prioritza context, narrativa i exploració.”',
    catalogEyebrow: 'Índex de projectes',
    catalogTitle: 'Casos d’estudi seleccionats',
    catalogIntro: '',
    catalogSearchLabel: 'Cercar projectes',
    catalogSearchPlaceholder: 'Cercar pel nom del projecte…',
    catalogCategoryPlaceholder: 'Categoria',
    catalogClearFilters: 'Netejar filtres',
    catalogViewCaseStudy: 'Veure cas d’estudi',
    catalogOpenProject: 'Obrir projecte',
    catalogNoImage: 'Sense imatge',
    catalogNoSummary: 'Sense resum disponible.',
    catalogLoading: 'Carregant projectes…',
    catalogError: 'Error',
    catalogNoMatches: 'No hi ha projectes que coincideixin amb la teva cerca.',
    catalogNoFilterMatches: 'No hi ha projectes per als filtres actuals.',
    detailBack: '← Tornar als projectes',
    detailLoading: 'Carregant projecte…',
    detailNotFound: 'Projecte no trobat.',
    detailProjectOverview: 'Resum del projecte',
    detailCategory: 'Categoria',
    detailClient: 'Client',
    detailUpdated: 'Actualitzat',
    detailTechnologies: 'Tecnologies',
    detailIndependent: 'Independent / Intern',
    detailRecentlyCurated: 'Curat recentment',
    detailNotSpecified: 'No especificat',
    detailStrategyLayer: 'Estratègia',
    detailExecutionLayer: 'Execució',
    detailTechnicalLayer: 'Tècnica',
    detailBusinessGoal: 'Context de negoci / Objectiu',
    detailProblem: 'Problema',
    detailSolution: 'Solució',
    detailDeliveryScope: 'Abast del delivery',
    detailResponsibilityScope: 'Abast de responsabilitat',
    detailArchitecture: 'Arquitectura',
    detailAIUsage: 'Ús d’IA',
    detailIntegrations: 'Integracions',
    detailTechnicalDecisions: 'Decisions tècniques',
    detailChallenges: 'Reptes',
    detailResults: 'Resultats',
    detailMetrics: 'Mètriques',
    detailTimeline: 'Timeline',
    detailClientContext: 'Context client',
    detailVisualUnavailable: 'Visual del cas d’estudi no disponible',
  },
  en: {
    headerTitle: 'Project portfolio',
    headerSummary: 'Strategy, execution, and technical judgment.',
    navHome: 'Home',
    navLogin: 'Login',
    navSearch: 'Search',
    navAdmin: 'Admin',
    navLogout: 'Logout',
    headerCaption: 'Marlon Ly Bellido · Engineer',
    searchPlaceholder: 'Search projects by technology, client, or concept…',
    searchButton: 'Search',
    searchClear: 'Clear search',
    searchSuggestionsLabel: 'Search suggestions',
    landingSearchEyebrow: 'GUIDED SEARCH',
    landingSearchTitle: '',
    landingSearchLead: 'Search projects, cases, and real-world work.',
    landingSearchPlaceholder: 'Search a project, technology, or topic...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Show me the Printer 05 PLC migration', query: 'Printer 05' },
      { label: 'I want Allen-Bradley and CompactLogix work', query: 'CompactLogix' },
      { label: 'Find industrial automation projects with Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Search motion-control work with SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Selected digital work',
    landingTitle: 'A more editorial cover to present product, architecture, and execution.',
    landingLead:
      'PortfolioForge brings public projects together with stronger hierarchy, clearer blocks, and a more legible reading of the value behind each delivery.',
    landingPrimaryCta: 'Explore projects',
    landingSecondaryCta: 'Browse catalog',
    landingDesignIntent: 'Design intent',
    landingPrinciples: [
      'Case studies with clear narrative and decision-focused structure.',
      'Modular composition designed for desktop, tablet, and mobile.',
      'Dark palette preserved with editorial contrast and visual rhythm.',
    ],
    landingHighlights: [
      { value: '01', label: 'Public portfolio with structured project stories' },
      { value: '02', label: 'Responsive modular catalog for selected work' },
      { value: '03', label: 'Search and detail aligned under the same visual system' },
    ],
    landingPortfolioSystem: 'Portfolio system',
    landingShowcaseTitle: 'A modular composition to read the portfolio like a publication, not a list.',
    landingShowcaseCopy:
      'The landing page now separates message, exploration, and catalog into broader blocks so content can breathe on desktop without losing clarity on mobile.',
    landingQuoteEyebrow: 'Editorial rhythm',
    landingQuote: '“Solid blocks, better use of width, and a visual structure that prioritizes context, narrative, and exploration.”',
    catalogEyebrow: 'Project index',
    catalogTitle: 'Selected case studies',
    catalogIntro: '',
    catalogSearchLabel: 'Search projects',
    catalogSearchPlaceholder: 'Search by project name…',
    catalogCategoryPlaceholder: 'Category',
    catalogClearFilters: 'Clear filters',
    catalogViewCaseStudy: 'View case study',
    catalogOpenProject: 'Open project',
    catalogNoImage: 'No image',
    catalogNoSummary: 'No summary available.',
    catalogLoading: 'Loading projects…',
    catalogError: 'Error',
    catalogNoMatches: 'No projects match your search.',
    catalogNoFilterMatches: 'No projects match the current filters.',
    detailBack: '← Back to projects',
    detailLoading: 'Loading project…',
    detailNotFound: 'Project not found.',
    detailProjectOverview: 'Project overview',
    detailCategory: 'Category',
    detailClient: 'Client',
    detailUpdated: 'Updated',
    detailTechnologies: 'Technologies',
    detailIndependent: 'Independent / Internal',
    detailRecentlyCurated: 'Recently curated',
    detailNotSpecified: 'Not specified',
    detailStrategyLayer: 'Strategy',
    detailExecutionLayer: 'Execution',
    detailTechnicalLayer: 'Technical',
    detailBusinessGoal: 'Business context / goal',
    detailProblem: 'Problem',
    detailSolution: 'Solution',
    detailDeliveryScope: 'Delivery scope',
    detailResponsibilityScope: 'Responsibility scope',
    detailArchitecture: 'Architecture',
    detailAIUsage: 'AI usage',
    detailIntegrations: 'Integrations',
    detailTechnicalDecisions: 'Technical decisions',
    detailChallenges: 'Challenges',
    detailResults: 'Results',
    detailMetrics: 'Metrics',
    detailTimeline: 'Timeline',
    detailClientContext: 'Client context',
    detailVisualUnavailable: 'Case study visual unavailable',
  },
  de: {
    headerTitle: 'Projektportfolio',
    headerSummary: 'Strategie, Umsetzung und technisches Urteilsvermögen.',
    navHome: 'Start',
    navLogin: 'Login',
    navSearch: 'Suche',
    navAdmin: 'Admin',
    navLogout: 'Abmelden',
    headerCaption: 'Ing. Marlon Ly Bellido',
    searchPlaceholder: 'Projekte nach Technologie, Kunde oder Konzept suchen…',
    searchButton: 'Suchen',
    searchClear: 'Suche löschen',
    searchSuggestionsLabel: 'Suchvorschläge',
    landingSearchEyebrow: 'GEFÜHRTE SUCHE',
    landingSearchTitle: '',
    landingSearchLead: 'Suche Projekte, Fälle und reale Arbeit.',
    landingSearchPlaceholder: 'Suche ein Projekt, eine Technologie oder ein Thema...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Zeig mir die PLC-Migration von Printer 05', query: 'Printer 05' },
      { label: 'Ich suche Arbeiten mit Allen-Bradley und CompactLogix', query: 'CompactLogix' },
      { label: 'Finde Industrial-Automation-Projekte mit Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Suche Motion-Control mit SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Ausgewählte digitale Arbeit',
    landingTitle: 'Eine editorischere Startseite für Produkt, Architektur und Umsetzung.',
    landingLead:
      'PortfolioForge bündelt öffentliche Projekte mit besserer Hierarchie, klareren Blöcken und einer besser lesbaren Darstellung des Werts hinter jeder Lieferung.',
    landingPrimaryCta: 'Projekte erkunden',
    landingSecondaryCta: 'Katalog ansehen',
    landingDesignIntent: 'Designabsicht',
    landingPrinciples: [
      'Case Studies mit klarer Narration und Fokus auf Entscheidungen.',
      'Modulare Komposition für Desktop, Tablet und Mobile.',
      'Dark-Palette mit editorischem Kontrast und visuellem Rhythmus.',
    ],
    landingHighlights: [
      { value: '01', label: 'Öffentliches Portfolio mit strukturierten Projektgeschichten' },
      { value: '02', label: 'Responsiver modularer Katalog für ausgewählte Arbeiten' },
      { value: '03', label: 'Suche und Detailseiten im selben visuellen System' },
    ],
    landingPortfolioSystem: 'Portfolio-System',
    landingShowcaseTitle: 'Eine modulare Komposition, damit sich das Portfolio wie eine Publikation liest, nicht wie eine Liste.',
    landingShowcaseCopy:
      'Die Landing trennt jetzt Botschaft, Exploration und Katalog in großzügigere Blöcke, damit der Inhalt auf Desktop besser atmet, ohne auf Mobile an Klarheit zu verlieren.',
    landingQuoteEyebrow: 'Editorialer Rhythmus',
    landingQuote: '„Solide Blöcke, bessere Nutzung der Breite und eine visuelle Struktur, die Kontext, Narration und Exploration priorisiert.“',
    catalogEyebrow: 'Projektindex',
    catalogTitle: 'Ausgewählte Fallstudien',
    catalogIntro: '',
    catalogSearchLabel: 'Projekte suchen',
    catalogSearchPlaceholder: 'Nach Projektnamen suchen…',
    catalogCategoryPlaceholder: 'Kategorie',
    catalogClearFilters: 'Filter löschen',
    catalogViewCaseStudy: 'Fallstudie ansehen',
    catalogOpenProject: 'Projekt öffnen',
    catalogNoImage: 'Kein Bild',
    catalogNoSummary: 'Keine Zusammenfassung verfügbar.',
    catalogLoading: 'Projekte werden geladen…',
    catalogError: 'Fehler',
    catalogNoMatches: 'Keine Projekte passen zu deiner Suche.',
    catalogNoFilterMatches: 'Keine Projekte für die aktuellen Filter.',
    detailBack: '← Zurück zu Projekten',
    detailLoading: 'Projekt wird geladen…',
    detailNotFound: 'Projekt nicht gefunden.',
    detailProjectOverview: 'Projektüberblick',
    detailCategory: 'Kategorie',
    detailClient: 'Kunde',
    detailUpdated: 'Aktualisiert',
    detailTechnologies: 'Technologien',
    detailIndependent: 'Unabhängig / Intern',
    detailRecentlyCurated: 'Kürzlich kuratiert',
    detailNotSpecified: 'Nicht angegeben',
    detailStrategyLayer: 'Strategie',
    detailExecutionLayer: 'Umsetzung',
    detailTechnicalLayer: 'Technik',
    detailBusinessGoal: 'Geschäftskontext / Ziel',
    detailProblem: 'Problem',
    detailSolution: 'Lösung',
    detailDeliveryScope: 'Lieferumfang',
    detailResponsibilityScope: 'Verantwortungsbereich',
    detailArchitecture: 'Architektur',
    detailAIUsage: 'KI-Einsatz',
    detailIntegrations: 'Integrationen',
    detailTechnicalDecisions: 'Technische Entscheidungen',
    detailChallenges: 'Herausforderungen',
    detailResults: 'Ergebnisse',
    detailMetrics: 'Metriken',
    detailTimeline: 'Timeline',
    detailClientContext: 'Kundenkontext',
    detailVisualUnavailable: 'Visual der Fallstudie nicht verfügbar',
  },
};

export function getMessages(locale: PublicLocale): Messages {
  return {
    ...(baseMessages[locale] ?? baseMessages.es),
    ...(authMessages[locale] ?? authMessages.es),
  };
}
