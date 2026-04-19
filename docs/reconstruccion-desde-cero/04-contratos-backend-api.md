# Contratos de backend y API

## Principio

Separar tres superficies HTTP: pública, privada autenticada y admin. El rebuild no necesita clonar exactamente todos los handlers actuales, pero sí debe preservar estas capacidades contractuales.

## API pública

### Proyectos

- `GET /api/v1/public/projects`
  - lista proyectos publicados/activos;
  - soporta localización pública por `?lang=`;
  - debe devolver `assistant_available` derivado.

- `GET /api/v1/public/projects/:slug`
  - devuelve detalle público enriquecido;
  - incluye `profile`, `technologies`, `media` e `images` de compatibilidad.

### Búsqueda

- `GET /api/v1/public/search`
  - query principal + filtros por categoría, cliente y tecnologías;
  - respuesta con `data[]` y `meta`;
  - resultado incluye `explanation` y `evidence[]`;
  - debe tener CORS público y rate limiting.

### Auth pública

- `POST /api/v1/public/signup`
- `POST /api/v1/public/login`
- `POST /api/v1/public/login/google`
- `POST /api/v1/public/email-verification/request`
- `POST /api/v1/public/email-verification/resend`
- `POST /api/v1/public/email-verification/verify`

### Site settings pública

- `GET /api/v1/public/site-settings`

## API privada autenticada

### Sesión y perfil

- `GET /api/v1/private/me`
- `PUT /api/v1/private/me/profile`

### Assistant por proyecto

- `POST /api/v1/private/projects/:slug/assistant/messages`
  - recibe `question`, `history[]`, `lang`;
  - valida proyecto activo + `source_markdown_url` + elegibilidad del usuario;
  - responde `answer` grounded.

## API admin

### Login admin

- `POST /api/v1/admin/login`

### Gestión de proyectos

- `GET /api/v1/admin/projects`
- `POST /api/v1/admin/projects`
- `GET /api/v1/admin/projects/:id`
- `PUT/PATCH /api/v1/admin/projects/:id`
- `DELETE /api/v1/admin/projects/:id`
- `PATCH /api/v1/admin/projects/:id/status`

### Enrichment y localizaciones

- `PUT /api/v1/admin/projects/:id/enrichment`
- `GET /api/v1/admin/projects/:id/localizations`
- `PUT /api/v1/admin/projects/:id/localizations/:locale`

### Search admin

- readiness por proyecto o global;
- reembed individual;
- reembed batch de stale content.

### Tecnologías

- CRUD completo de `technologies`.

### Usuarios

- listado admin;
- detalle;
- update;
- delete.

### Ajustes globales

- `GET /api/v1/admin/site-settings`
- `PUT /api/v1/admin/site-settings`

### Workflow de case study

- `GET /api/v1/admin/settings/case-study-workflow`
- `POST /api/v1/admin/settings/case-study-runs`
- `GET /api/v1/admin/settings/case-study-runs/:id`
- `GET /api/v1/admin/settings/case-study-runs/:id/logs`
- `POST /api/v1/admin/settings/case-study-runs/:id/resume`
- `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/confirm`
- `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/start`
- `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/retry`

## Contratos de datos relevantes

### Project público

Debe exponer al menos:

- `id`, `name`, `slug`, `description`, `category`, `client_name`
- `status`, `featured`, `active`, `assistant_available`
- `images`, `media`, `profile`, `technologies`

### AdminProject write

Debe aceptar:

- `name`, `description`, `category`, `client_name`
- `source_markdown_url`
- `active`
- `media[]`
- aliases legacy `product_name`, `brand` solo como compatibilidad temporal

### SearchResponse

- `data[]` con `id`, `slug`, `title`, `category`, `client_name`, `summary`, `technologies`, `hero_image`, `score`, `explanation`, `evidence[]`
- `meta` con total, page_size, cursor, query y filtros aplicados

## Reglas backend no negociables

- nunca exponer `source_markdown_url` en payload público;
- aplicar locale público sin mutar la fuente base;
- usar contratos explícitos para compatibilidad legacy;
- desacoplar repositorio/search/assistant del HTTP layer.
