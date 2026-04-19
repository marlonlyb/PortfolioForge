# Modelo de dominio y datos

## Regla de naming

El dominio canónico es `project`. La persistencia puede reutilizar `products` por compatibilidad, pero el rebuild debe diseñarse pensando en entidades de portfolio, no en e-commerce.

## Entidades núcleo

### Project

Campos mínimos:

- `id`
- `name`
- `slug`
- `description`
- `category`
- `client_name` (fallback legacy: `brand`)
- `status` (`draft|published|archived`)
- `featured`
- `active`
- `assistant_available` derivado de `source_markdown_url`
- `source_markdown_url` privado a nivel admin/backend
- `images` como compatibilidad/fallback
- timestamps

Persistencia real observada:

- base en `products`
- extensiones en migraciones de abril 2026 para `slug`, `category`, `active`, `client_name`, `status`, `featured`, `source_markdown_url`

### ProjectProfile

Describe detalle enriquecido del proyecto. Campos:

- `business_goal`
- `problem_statement`
- `solution_summary`
- `delivery_scope`
- `responsibility_scope`
- `architecture`
- `ai_usage`
- `integrations` (lista estructurada)
- `technical_decisions` (lista estructurada)
- `challenges` (lista estructurada)
- `results` (lista estructurada)
- `metrics` (objeto)
- `timeline` (lista estructurada)

### Technology

- `id`
- `name`
- `slug`
- `category`
- `icon`
- `color`

Relación N:N con proyectos mediante `project_technologies`.

### ProjectMedia

Contrato objetivo:

- `fallback_url`
- `low_url`
- `medium_url`
- `high_url`
- `caption`
- `alt_text`
- `featured`
- `sort_order`
- `media_type`

Notas:

- el storage real pasó por una transición `url/thumbnail_url/full_url` hacia variantes explícitas;
- `images[]` debe reconstruirse como compatibilidad de lectura.

### ProjectLocalization

- clave compuesta por `project_id + locale + field_key`
- `locale`: `ca|en|de`
- `value` JSONB
- `mode`: `auto|manual`
- `source_hash`

La base `es` vive en `products`/`project_profiles`; no se persiste como fila de localización.

### User

Campos relevantes:

- `email`
- `password_hash` / login local
- `auth_provider`
- `provider_subject`
- `email_verified`
- `full_name`
- `company`
- `profile_completed`
- `assistant_eligible`
- `can_use_project_assistant`
- `is_admin`

### SiteSettings

- `public_hero_logo_url`
- `public_hero_logo_alt`

### CaseStudyWorkflowRun

Controla la operación editorial asistida desde admin:

- fuente resuelta
- opciones de localización/reembed
- URL canónica publicada
- proyecto asociado
- pasos y logs persistidos

## Tablas mínimas del rebuild

- `users`
- `products` o `projects`
- `project_profiles`
- `technologies`
- `project_technologies`
- `project_media`
- `project_localizations`
- `project_search_documents`
- `site_settings`
- `case_study_workflow_runs`
- `case_study_workflow_steps`
- `case_study_workflow_logs`

## Índices/capacidades obligatorias

- índice por `slug`, `category`, `active`;
- FTS sobre documento compuesto;
- GIN `pg_trgm` sobre texto de búsqueda;
- HNSW/pgvector para `search_embedding` cuando extensión esté disponible;
- unicidad por `technology.slug` y correo normalizado de user.

## Reglas de consistencia

- `slug` estable y reutilizable en URL pública y canonical publish;
- `client_name` es el campo público/localizable; `brand` solo compat;
- si cambia contenido indexable, debe recomponerse `project_search_documents`;
- si cambia base `es`, las localizaciones derivadas deben invalidarse o regenerarse;
- un proyecto no debe anunciar assistant si `source_markdown_url` está vacío o roto.
