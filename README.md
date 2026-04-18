# PortfolioForge

PortfolioForge es un portfolio interactivo construido con React, Go y PostgreSQL. Su núcleo combina búsqueda híbrida por evidencia real, explicaciones acotadas y un project assistant por proyecto, grounded en markdown fuente del propio caso de estudio, sin chatbot libre generalista.

## Estado actual

El proyecto ya tiene implementado y archivado en SDD:

- búsqueda híbrida por evidencia real
- explicaciones acotadas por resultado
- integración con OpenAI para embeddings y explanations
- project assistant por proyecto basado en markdown fuente
- CRUD admin de tecnologías
- enriquecimiento de proyectos desde admin
- re-composición y re-embedding del documento de búsqueda
- automatización local de revisión visual con Playwright

## Stack

- **Backend**: Go 1.20 + Echo v4 + pgx/v5
- **Frontend**: React 19 + TypeScript 5.9 + Vite 7
- **Base de datos**: PostgreSQL 16
- **Extensiones DB**: `unaccent`, `pg_trgm`, `vector (pgvector)`
- **LLM / Embeddings**: OpenAI (`text-embedding-3-small`, `gpt-4o-mini`)
- **Auth**: JWT Bearer sessions, canonical `/login` as the only visible public nav entry, Google or local email/password for all local users, `/signup` kept for direct access or auth-flow links, `/admin/login` kept as a compatibility alias to the same login behavior, OTP reserved for email verification/resend, and assistant enabled only for eligible sessions
- **Revisión visual local**: Playwright

## Funcionalidad principal

### Búsqueda por evidencia real
La interacción principal del producto. El visitante escribe uno o varios conceptos, por ejemplo `SIEMENS COMMISSIONING`, y el sistema:

1. recupera proyectos usando:
   - filtros estructurados (`category`, `client`, `technologies`)
   - búsqueda léxica con PostgreSQL FTS
   - búsqueda fuzzy con `pg_trgm`
   - búsqueda semántica con `pgvector` + OpenAI embeddings
2. fusiona el ranking por relevancia
3. muestra una explicación corta basada solo en evidencia del proyecto

### Sitio público
- landing con buscador principal
- resultados de búsqueda
- catálogo de proyectos
- detalle de proyecto tipo case study con payload público más liviano por presentación
- páginas públicas de proyecto siempre accesibles; assistant visible solo para sesiones elegibles cuando `assistant_available=true`

### Consola admin
- login con JWT
- CRUD admin de proyectos
- CRUD admin de tecnologías
- enriquecimiento de proyectos (`project_profiles`)
- configuración privada de `source_markdown_url` por proyecto
- readiness de búsqueda
- re-embed manual y recomposición del documento de búsqueda

## Requisitos de PostgreSQL

PortfolioForge requiere estas extensiones:

| Extensión | Propósito |
|-----------|-----------|
| `unaccent` | búsqueda insensible a acentos |
| `pg_trgm` | tolerancia a typos y búsqueda fuzzy |
| `vector` | búsqueda semántica por embeddings |

Instalación en Ubuntu/Debian:

```bash
sudo apt-get install -y postgresql-16-pgvector
```

## Base de datos local

La base de datos usada por este proyecto es:

- **DB_NAME**: `portfolioforge`

No uses `proyectoemlb`; pertenece al proyecto original y debe permanecer intacta.

## Setup local

### 1. Backend

1. Copia `.env.example` a `.env`
2. Configura PostgreSQL, JWT y OpenAI
3. Asegúrate de que la base `portfolioforge` exista
4. Aplica las migraciones en orden:

```bash
psql -d portfolioforge -f sqlmigrations/20240617_2206_create_user.sql
psql -d portfolioforge -f sqlmigrations/20240624_1609_create_products.sql
psql -d portfolioforge -f sqlmigrations/20260410_1200_store_mvp_base.sql
psql -d portfolioforge -f sqlmigrations/20260411_0900_portfolioforge_extension.sql
psql -d portfolioforge -f sqlmigrations/20260412_0100_search_extensions.sql
psql -d portfolioforge -f sqlmigrations/20260412_0200_search_schema.sql
psql -d portfolioforge -f sqlmigrations/20260412_0300_search_backfill.sql
psql -d portfolioforge -f sqlmigrations/20260412_0400_search_embedding_text.sql
psql -d portfolioforge -f sqlmigrations/20260414_1100_project_assistant_chat.sql
psql -d portfolioforge -f sqlmigrations/20260415_1000_authenticated_project_assistant.sql
psql -d portfolioforge -f sqlmigrations/20260415_1100_local_email_verification.sql
```

5. Inicia el backend:

```bash
go run cmd/*.go
```

Backend por defecto:

- `http://localhost:8080`

### 2. Frontend

1. Copia `client/.env.example` a `client/.env` si necesitas sobreescribir la API
2. Instala dependencias:

```bash
cd client
npm install
```

3. Inicia Vite:

```bash
npm run dev -- --host 0.0.0.0 --port 5173
```

Frontend por defecto:

- `http://localhost:5173`

## Variables de entorno

### Backend

| Variable | Descripción | Requerido |
|----------|-------------|-----------|
| `SERVER_PORT` | Puerto del backend | Sí |
| `JWT_SECRET_KEY` | Secreto para firmar JWT | Sí |
| `DB_HOST` | Host de PostgreSQL | Sí |
| `DB_PORT` | Puerto de PostgreSQL | Sí |
| `DB_USER` | Usuario de PostgreSQL | Sí |
| `DB_PASSWORD` | Password de PostgreSQL | Sí |
| `DB_NAME` | Nombre de la base de datos | Sí |
| `DB_SSL_MODE` | SSL mode de PostgreSQL | Sí |
| `ENABLE_SEMANTIC_SEARCH` | Activa capa semántica | No |
| `OPENAI_API_KEY` | API key de OpenAI | Sí, para embeddings reales, explanations y project assistant |
| `GOOGLE_CLIENT_ID` | Client ID para verificar Google ID tokens en backend | Sí, para sign-in público con Google |
| `SMTP_HOST` | Host SMTP para OTP por email | No, pero requerido para entrega real |
| `SMTP_PORT` | Puerto SMTP para OTP por email | No, pero requerido para entrega real |
| `SMTP_USERNAME` | Usuario SMTP | No |
| `SMTP_PASSWORD` | Password SMTP | No |
| `EMAIL_FROM_ADDRESS` | Remitente usado por los OTP | No, pero requerido para entrega real |
| `EMAIL_FROM_NAME` | Nombre de remitente para OTP | No |

### Frontend

| Variable | Descripción | Requerido |
|----------|-------------|-----------|
| `VITE_API_BASE_URL` | URL base del backend | No |
| `VITE_GOOGLE_CLIENT_ID` | Google client ID para renderizar el botón público de sign-in | Sí, para sign-in público con Google |

Si `VITE_API_BASE_URL` no existe, el frontend usa `http://localhost:8080` por defecto.

## Rutas principales

### Frontend público
- `/` — landing con buscador principal
- `/search` — resultados de búsqueda
- `/projects` — catálogo
- `/projects/:slug` — detalle de proyecto
- `/login` — canonical login with Google or local email/password for any registered user
- `/signup` — direct-access signup route kept inside the auth flow (not shown in the public header)
- `/admin/login` — compatibility alias to the shared login screen, still reachable by direct URL
- `/verify-email` — 6-digit OTP flow for email verification / re-verification only
- `/complete-profile` — completar `full_name` y `company` antes de usar el assistant

### Frontend admin
- `/admin/projects`
- `/admin/projects/new`
- `/admin/projects/:id`
- `/admin/technologies`
- `/admin/technologies/new`
- `/admin/technologies/:id`

### API pública
- `GET /api/v1/public/search?q=...&technologies=...&category=...&client=...`
- `POST /api/v1/public/signup`
- `POST /api/v1/public/login`
- `POST /api/v1/public/login/google`
- `POST /api/v1/public/email-verification/request`
- `POST /api/v1/public/email-verification/resend`
- `POST /api/v1/public/email-verification/verify`
- `POST /api/v1/admin/login`
- `POST /api/v1/public/email-verification/request`
- `POST /api/v1/public/email-verification/resend`
- `POST /api/v1/public/email-verification/verify`
- `GET /api/v1/public/projects`
- `GET /api/v1/public/projects/:slug`

### API privada
- `GET /api/v1/private/me`
- `PUT /api/v1/private/me/profile`
- `POST /api/v1/private/projects/:slug/assistant/messages`

### API admin
- `GET /api/v1/admin/projects/:id/readiness`
- `POST /api/v1/admin/projects/:id/reembed`
- `POST /api/v1/admin/projects/reembed-stale`
- `PUT /api/v1/admin/projects/:id/enrichment`
- `GET /api/v1/admin/technologies`
- `GET /api/v1/admin/technologies/:id`
- `POST /api/v1/admin/technologies`
- `PUT /api/v1/admin/technologies/:id`
- `DELETE /api/v1/admin/technologies/:id`

## Cómo probar rápido

1. Levanta backend y frontend
2. Entra a `http://localhost:5173/login`
3. Usa Google o inicia sesión con email/password en `/login`; desde ahí puedes abrir `Sign up`, o entrar directo a `/signup`, para crear una cuenta local pública nueva
4. Si quieres validar compatibilidad con bookmarks antiguos, entra manualmente a `http://localhost:5173/admin/login`; el formulario aplica las mismas reglas que `/login`
5. Crea tecnologías en `/admin/technologies`
6. Enriquece un proyecto desde `/admin/projects/:id`
7. Si el proyecto debe exponer assistant, configura `source_markdown_url` con una URL HTTPS pública que apunte al markdown fuente
8. Ejecuta re-embed
9. Prueba búsquedas como:
   - `SIEMENS`
   - `commissioning`
   - `react scada`
10. Si `assistant_available=true` en el payload público, abre `/projects/:slug`, confirma que el assistant no aparece en sesión anónima y luego entra con Google o con una cuenta local verificada
11. Si el usuario público entra por email, valida que reciba un OTP de 6 dígitos, que `/verify-email` complete la sesión y que el assistant siga bloqueado hasta completar perfil
12. Si el usuario público no tiene `full_name` o `company`, completa `/complete-profile` y valida que recién entonces se habilita el chat

## Arquitectura

### Backend

Arquitectura hexagonal (ports & adapters):

```text
cmd/                    entry points, wiring, routes
domain/ports/           interfaces
domain/services/        lógica de dominio
infrastructure/postgres repos PostgreSQL
infrastructure/handlers handlers Echo
infrastructure/embedding providers de embeddings
infrastructure/explanation providers de explanations
model/                  modelos de dominio
sqlmigrations/          migraciones SQL
```

### Frontend

Estructura orientada por features:

```text
client/src/
  app/
  features/
    landing/
    search/
    catalog/
    admin-products/
    admin-technologies/
    auth/
  shared/
```

## Notas de transición

- internamente todavía existen nombres heredados como `product` y `brand`
- el dominio funcional oficial del producto es `project`, `client_name`, `technologies`, `project_profiles`
- el mapping canónico vigente es: `Client / Context -> brand` (storage/admin legacy, consumo público `client_name`), `Published -> active`, `Technologies -> technology_ids`, `Main images -> media` + lista legacy derivada
- la carga manual en UI y el workflow automático ya están alineados a nivel lógico, pero todavía no son un contrato 1:1 exacto
- para proyectos analizados desde repositorio/carpeta, `90. dev_portfolioforge/<Project_Name>.md` pasa a ser la fuente editorial canónica para crear el proyecto en PortfolioForge
- cuando ese markdown ya existe, la UI debe entenderse como la capa actual de persistencia/ejecución del contenido definido allí, no como la fuente primaria de autoría
- si el markdown fuente ya existe, volver a analizar el repo completo debe ser secundario y solo justificarse cuando haya evidencia nueva o desactualización real
- el markdown canónico local y su copia publicada en `source_markdown_url` deben mantenerse en castellano y alineados; no debe quedar la URL remota en otro idioma si la base editorial oficial es `es`
- la convención oficial de publicación del markdown canónico es `https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md`
- si se corrige la base en castellano de un proyecto, hay que regenerar después las locales derivadas (`ca`, `en`, `de`)
- `source_markdown_url` solo existe en admin/privado; la API pública expone `assistant_available` derivado de ese campo
- las tablas `orders`, `order_items` y `product_variants` siguen en la base heredada pero no forman parte del portfolio activo

## Documentación

Documentos normativos vigentes dentro del repo:

- `docs/PRD.md` — fuente de verdad del producto
- `docs/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md` — cómo generar el markdown canónico
- `docs/PROJECT-RUNTIME-INGESTION-GUIDE.md` — cómo poblar o actualizar la UI/DB desde ese markdown

Regla:

- cualquier workflow similar fuera de `docs/` debe tratarse como material legacy/informativo, no como especificación operativa vigente.

Regla documental vigente para nuevas altas de proyectos:

- si existe `90. dev_portfolioforge/<Project_Name>.md` dentro del repo/carpeta analizada, ese archivo debe ser el punto de partida preferido para crear el proyecto en PortfolioForge
- ese markdown y su copia publicada para `source_markdown_url` deben estar en castellano como fuente editorial primaria
- la URL remota publicada debe construirse siempre con el mismo slug del proyecto: `https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md`
- la carga manual en UI debe limitarse a persistir y revisar lo que el markdown fuente ya define, evitando re-redactar o reinventar contenido sin evidencia nueva
- para nuevos markdowns fuente, la recomendación por defecto ahora es `Published=true`; usar `Published=false` solo cuando exista una decisión explícita de mantener el proyecto interno/no público
- después de importar o editar un proyecto a partir de ese markdown, hay que verificar el resultado real en payload admin/público o en DB; escribir en la base no alcanza como criterio de éxito
- la verificación mínima obligatoria es campo por campo contra el `.md`: `title`, `active/published`, `client/context`, tecnologías, narrativa del `profile`, listas enriquecidas, métricas y media principal
- si el proyecto debe tener assistant, además hay que persistir `source_markdown_url`, verificar `assistant_available=true` en el payload público, confirmar que la URL markdown no se expone públicamente y validar el flujo de sign-in elegible (Google verificado + perfil completo, o admin local)
- si el canonical o el runtime base estaban en otro idioma por error, primero se corrige el castellano y luego se regeneran `ca`, `en` y `de`
- si el markdown fuente dice `Published=false`, el proyecto no puede quedar activo ni visible públicamente; si el importador cae en fallback, parsea parcialmente o mezcla assets ajenos, eso debe tratarse como fallo y no como import exitoso

### Publicación del markdown canónico

El repo ya incluye un subcomando para publicar por FTPS el canonical generado dentro de `90. dev_portfolioforge`:

```bash
go run ./cmd canonical-publish --case-dir "<ruta-del-caso>"
```

Comportamiento:

- detecta el directorio `90. dev_portfolioforge/<slug>/`
- publica toda la carpeta del slug, no solo el `.md`
- valida luego la URL pública final

Variables de entorno esperadas:

- `PF_FTP_HOST`
- `PF_FTP_PORT`
- `PF_FTP_USER`
- `PF_FTP_PASSWORD`
- `PF_PUBLIC_BASE`
- `PF_FTP_REMOTE_BASE` (opcional; por defecto `/` para cuentas FTP ya restringidas a la carpeta remota correcta)

Uso recomendado:

- si hay más de un slug dentro de `90. dev_portfolioforge`, usar `--slug <slug>`
- usar `--dry-run` para verificar destino remoto y URL pública antes de subir

### Dirección documental actual

- el estándar editorial ya debe escribir contenido pensando en tres niveles: **search readiness**, **case study readiness** y **assistant readiness**
- el código actual implementa **search readiness** y un project assistant markdown-grounded por proyecto; **case study readiness** y una assistant-readiness más profunda siguen dependiendo de disciplina editorial y futura evolución del retrieval
- listas como `Technical Decisions`, `Integrations`, `Results`, `Timeline` y `Challenges` deben redactarse en líneas semi-estructuradas para reforzar recuperación futura y consumo por asistentes
- `Main images` sigue existiendo por compatibilidad editorial, pero el contrato canónico ya debe pensarse como ítems de media con variantes `low`, `medium`, `high`, más `caption`, `alt_text`, `featured` y `sort_order`
- la convención pública/canónica por defecto para assets de imagen es `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen0<numero>_<low|medium|high>.webp`; `Main images` debe seguir priorizando `_medium` y el mínimo recomendado se mantiene en 5 imágenes cuando haya material suficiente

## Troubleshooting

- Si los proyectos públicos devuelven `500` después de traer el último código, revisa migraciones pendientes, especialmente `sqlmigrations/20260414_1100_project_assistant_chat.sql`.
- Si no aparece el assistant en un proyecto habilitado, verifica `assistant_available=true`, que exista una sesión elegible (`can_use_project_assistant=true`) y que el usuario público haya iniciado sesión con Google con email verificado y perfil completo.
- Si el chat falla después de aparecer, verifica `source_markdown_url`, que el markdown sea alcanzable por HTTPS desde el backend, que `OPENAI_API_KEY` esté configurada y que el request vaya a `POST /api/v1/private/projects/:slug/assistant/messages`.
