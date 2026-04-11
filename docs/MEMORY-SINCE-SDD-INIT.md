# Memory since `/sdd-init`

## Propósito

Este documento resume, en orden práctico, el trabajo realizado desde la inicialización SDD del proyecto hasta el estado actual. No reemplaza el `PRD`; sirve como memoria operativa e histórica del avance.

## 1. Inicialización SDD

- Se ejecutó `/sdd-init` en modo `engram`.
- Se detectó el stack real del proyecto:
  - Go 1.20 + Echo v4
  - PostgreSQL con pgx/v5
  - React 19 + TypeScript 5.9 + Vite 7
- Se persistió el contexto arquitectónico y el skill registry.

## 2. Limpieza de documentación

- Se eliminó la documentación heredada de tienda / e-commerce.
- Se consolidó la documentación activa en:
  - `docs/PRD.md`
- Se definió el criterio documental:
  - `README.md` como guía del repo
  - `docs/` para documentación del proyecto
  - sin duplicación innecesaria

## 3. Definición del producto

- El portfolio quedó definido como una plataforma con:
  - búsqueda por evidencia real
  - explicaciones acotadas por resultado
  - catálogo público
  - panel admin
- Se descartó explícitamente el chatbot libre.
- La home quedó definida como el punto de entrada principal con barra de búsqueda.
- Se formalizó el uso de OpenAI como proveedor concreto para embeddings y explanations.

## 4. Dirección visual

- Se aplicó una paleta oscura tipo developer a la landing.
- La referencia visual se tomó solo como guía de color y atmósfera, no como layout literal.
- Se dejó como dirección visual oficial una estética técnica, oscura y de alto contraste.

## 5. Change SDD: `evidence-based-project-search`

### Planeamiento

- Se abrió el change `evidence-based-project-search`.
- Se produjeron y persistieron:
  - proposal
  - design
  - spec
  - tasks

### Implementación

Se implementó:

- búsqueda híbrida con:
  - filtros estructurados
  - PostgreSQL FTS
  - `pg_trgm`
  - `pgvector`
- servicio de búsqueda y ranking
- explicaciones por resultado basadas en evidencia
- endpoint público de búsqueda
- frontend de búsqueda:
  - SearchBar
  - SearchResultsPage
  - SearchFilters
  - SearchResultCard
- readiness y re-embed admin
- limpieza fuerte de remanentes e-commerce

### Verificación

- Se corrigieron warnings de diseño/implementación.
- Se validó:
  - `go build ./...`
  - `go test ./...`
  - `npx tsc --noEmit`
- El change fue archivado como completado.

## 6. Base de datos e infraestructura local

- Se creó la base dedicada `portfolioforge`.
- Se evitó tocar la base `proyectoemlb` del proyecto original.
- Se instalaron / habilitaron extensiones:
  - `unaccent`
  - `pg_trgm`
  - `vector`
- Se aplicaron las migraciones SQL del proyecto.
- Se creó `.env` local con la conexión correcta.
- Se confirmó que el backend arranca correctamente en `:8080`.

## 7. README del proyecto

- `README.md` fue reescrito para reflejar el estado real del sistema.
- Ahora documenta:
  - stack actual
  - variables de entorno correctas
  - uso de OpenAI
  - rutas públicas y admin
  - setup local completo
  - flujo de prueba manual

## 8. Change SDD: `portfolio-domain-enrichment`

### Planeamiento

- Se abrió el change `portfolio-domain-enrichment`.
- Se produjeron:
  - proposal
  - design
  - spec
  - tasks

### Objetivo

Enriquecer el dominio del portfolio con:

- integración real con OpenAI
- CRUD admin de tecnologías
- enriquecimiento de proyectos (`project_profiles`)
- re-embedding al actualizar contenido relevante

### Implementación

Se implementó:

- `OpenAIEmbeddingProvider` con `text-embedding-3-small`
- `OpenAIExplanationProvider` con `gpt-4o-mini`
- fallbacks seguros a providers no-op/template
- CRUD admin de tecnologías
- endpoint `PUT /api/v1/admin/projects/:id/enrichment`
- enriquecimiento del formulario admin de proyectos

### Correcciones posteriores

Durante verify se detectaron gaps y se corrigieron:

- regeneración del spec faltante
- `GET /api/v1/admin/technologies/:id`
- atomicidad del enrichment + re-embedding
- fallo de embedding hace fallar la request
- embeddings desde texto compuesto crudo, no desde `tsvector`
- inclusión de `problem_statement` en composición
- prompt endurecido a exactamente una oración
- tests adicionales del change

### Verificación y archivo

- Se alcanzó verificación final con estado archive-ready.
- Quedó deuda de verificación manual en flujos UI frontend por falta de runner frontend.
- El change fue archivado con warnings únicamente.

## 9. OpenAI y configuración local

- Se definió usar OpenAI como proveedor real.
- Variables relevantes:
  - `OPENAI_API_KEY`
  - `ENABLE_SEMANTIC_SEARCH=true`
- La app ya puede correr con embeddings reales si la API key está presente.

## 10. Estado funcional actual

Hoy se puede probar:

- búsqueda pública por evidencia real
- resultados con explicación
- login admin
- CRUD de tecnologías
- enriquecimiento admin de proyectos
- recomposición y re-embed

## 11. Datos locales preparados para prueba

- Se creó un proyecto de prueba orientado a `SIEMENS`.
- Se creó un usuario admin local para desarrollo.
- Se confirmó que el login admin funciona contra `POST /api/v1/public/login`.

## 12. Deuda y siguientes pasos sugeridos

Pendientes razonables para próximos changes:

1. tests frontend / E2E para flujos admin
2. renaming interno `product -> project`
3. carga de proyectos reales
4. ajuste fino de prompts y scoring semántico
5. eliminación final de tablas heredadas no usadas (`orders`, `order_items`, `product_variants`)

## 13. Resumen ejecutivo

Desde `/sdd-init`, PortfolioForge pasó de ser una base transicional heredada a una plataforma funcional con:

- búsqueda híbrida real
- semántica con OpenAI
- explicación acotada por evidencia
- panel admin para tecnologías y enrichment
- infraestructura local operativa
- documentación consolidada
- dos changes SDD completos y archivados
