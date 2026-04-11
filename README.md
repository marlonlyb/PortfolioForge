# PortfolioForge

PortfolioForge es un portfolio interactivo construido con React, Go y PostgreSQL. Su núcleo es una búsqueda híbrida por evidencia real que combina filtros estructurados, FTS, búsqueda fuzzy y búsqueda semántica con embeddings de OpenAI para recuperar proyectos relevantes y explicarlos en lenguaje natural, sin chatbot libre.

## Estado actual

El proyecto ya tiene implementado y archivado en SDD:

- búsqueda híbrida por evidencia real
- explicaciones acotadas por resultado
- integración con OpenAI para embeddings y explanations
- CRUD admin de tecnologías
- enriquecimiento de proyectos desde admin
- re-composición y re-embedding del documento de búsqueda

## Stack

- **Backend**: Go 1.20 + Echo v4 + pgx/v5
- **Frontend**: React 19 + TypeScript 5.9 + Vite 7
- **Base de datos**: PostgreSQL 16
- **Extensiones DB**: `unaccent`, `pg_trgm`, `vector (pgvector)`
- **LLM / Embeddings**: OpenAI (`text-embedding-3-small`, `gpt-4o-mini`)
- **Auth**: JWT Bearer para admin

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
- detalle de proyecto tipo case study

### Consola admin
- login con JWT
- CRUD admin de proyectos
- CRUD admin de tecnologías
- enriquecimiento de proyectos (`project_profiles`)
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
| `OPENAI_API_KEY` | API key de OpenAI | Sí, para embeddings reales |

### Frontend

| Variable | Descripción | Requerido |
|----------|-------------|-----------|
| `VITE_API_BASE_URL` | URL base del backend | No |

Si `VITE_API_BASE_URL` no existe, el frontend usa `http://localhost:8080` por defecto.

## Rutas principales

### Frontend público
- `/` — landing con buscador principal
- `/search` — resultados de búsqueda
- `/projects` — catálogo
- `/projects/:slug` — detalle de proyecto
- `/login` — acceso admin

### Frontend admin
- `/admin/projects`
- `/admin/projects/new`
- `/admin/projects/:id`
- `/admin/technologies`
- `/admin/technologies/new`
- `/admin/technologies/:id`

### API pública
- `GET /api/v1/public/search?q=...&technologies=...&category=...&client=...`
- `POST /api/v1/public/login`
- `GET /api/v1/public/projects`
- `GET /api/v1/public/projects/:slug`

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
3. Inicia sesión con un usuario admin local
4. Crea tecnologías en `/admin/technologies`
5. Enriquece un proyecto desde `/admin/projects/:id`
6. Ejecuta re-embed
7. Prueba búsquedas como:
   - `SIEMENS`
   - `commissioning`
   - `react scada`

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
- las tablas `orders`, `order_items` y `product_variants` siguen en la base heredada pero no forman parte del portfolio activo

## Documentación

- `docs/PRD.md` — fuente de verdad del producto
- `docs/MEMORY-SINCE-SDD-INIT.md` — memoria resumida del trabajo realizado desde `/sdd-init`
