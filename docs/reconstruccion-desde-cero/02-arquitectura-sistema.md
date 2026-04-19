# Arquitectura del sistema

## Vista general

PortfolioForge se reconstruye como un sistema web full-stack con backend Go, frontend React y PostgreSQL como fuente de persistencia. La arquitectura separa claramente capa editorial, capa runtime y capacidades de IA.

## Componentes principales

### 1. Frontend web

- React 19 + TypeScript + Vite.
- Dos superficies: pública y admin.
- Consume APIs HTTP del backend.
- Aplica localización pública por `lang` y controla sesión cliente.

### 2. Backend API

- Go + Echo.
- Expone rutas públicas, privadas autenticadas y admin.
- Orquesta dominio `project`, búsqueda, assistant, auth, localización y workflow operativo.
- Encapsula integraciones con OpenAI, SMTP, Google identity y tooling legacy de publicación/compatibilidad (`publish_canonical`, FTPS).

### 3. PostgreSQL

- Persistencia operacional principal.
- Usa FTS, `pg_trgm` y `pgvector` para búsqueda híbrida.
- Mantiene tablas de proyectos, perfiles, tecnologías, media, localizaciones, users y workflow runs.

### 4. Capa editorial externa al runtime

- El markdown canónico no nace dentro del runtime.
- Se produce en la carpeta estudiada original bajo `90. dev_portfolioforge/<slug>/<slug>.md`.
- Su publicación remota alimenta `source_markdown_url`.

## Bounded contexts sugeridos

### Catalog / Project Runtime

Responsable de proyectos públicos, detalle, enrichment admin y media.

### Search

Responsable de componer documentos indexables, embeddings, explicaciones y readiness de indexación.

### Assistant

Responsable de recuperar markdown remoto, cachearlo y responder preguntas grounded por proyecto.

### Localization

Responsable de derivar y aplicar locales públicas desde la base `es`.

### Identity & Access

Responsable de signup/login, login Google, verificación de email, perfil completo y elegibilidad para assistant.

### Editorial Workflow

Responsable de alinear el runtime con la fuente canónica ya publicada, importar/actualizar proyecto, backfill de localización y reembed; conserva `publish_canonical`/FTPS solo como compatibilidad legacy.

## Flujo macro de datos

1. Se genera o corrige markdown canónico fuera del runtime.
2. Se publica en URL remota estable por `slug`.
3. Solo cuando existe `source_markdown_url`, admin importa o actualiza el proyecto en runtime.
4. Backend persiste proyecto base, enrichment, tecnologías, media y localizaciones como proyección resumida/estructurada de la fuente remota.
5. Search recompone documento indexable y embeddings.
6. Público consulta catálogo, detalle o búsqueda.
7. Assistant usa el markdown remoto del proyecto autenticado.

## Decisiones estructurales obligatorias

- mantener separación entre handlers HTTP, servicios de dominio y repositorios;
- modelar contratos públicos/admin de forma explícita;
- tratar localización, búsqueda y assistant como capacidades dependientes del contenido editorial base;
- permitir transición controlada de storage legacy (`products`, `brand`) sin contaminar el lenguaje del dominio nuevo.

## Stack de referencia del rebuild

- Backend: Go 1.25+, Echo v4, pgx.
- Frontend: React 19, TypeScript 5.9, Vite 7.
- DB: PostgreSQL 16 con `unaccent`, `pg_trgm`, `vector`.
- IA: OpenAI para embeddings, explanations y assistant.

## Riesgos arquitectónicos a controlar

- mezclar markdown fuente con contenido UI resumido;
- romper la coherencia entre `client_name` público y `brand` legacy;
- acoplar assistant al contenido runtime en vez del markdown remoto;
- tratar FTPS o `publish_canonical` como flujo principal en vez de compatibilidad legacy;
- tratar búsqueda semántica como opcional de producto cuando en realidad forma parte del discoverability core.
