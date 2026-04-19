# PortfolioForge

PortfolioForge es una plataforma editorial para convertir experiencia profesional real en **proyectos públicos navegables, buscables y explicables**.

Combina cuatro piezas en un solo sistema:

- **markdown canónico** como fuente editorial del caso;
- **UI + DB runtime** como proyección resumida, estructurada y operativa;
- **búsqueda híbrida** para recuperar proyectos por evidencia real;
- **assistant por proyecto** grounded en el markdown remoto publicado.

---

## Resumen ejecutivo

Este repositorio implementa el sistema completo de PortfolioForge:

- **backend** en Go + Echo;
- **frontend** en React + TypeScript + Vite;
- **persistencia** en PostgreSQL 16;
- **search stack** con FTS + `pg_trgm` + `pgvector`;
- **consola admin** para proyectos, tecnologías, localización, búsqueda y workflows editoriales;
- **assistant autenticado** por proyecto usando `source_markdown_url`.

### Qué resuelve

PortfolioForge no trata los proyectos como simples fichas visuales. Los trata como un sistema editorial con trazabilidad:

1. la evidencia vive en un repositorio o carpeta fuente;
2. se consolida en un canonical en `90. dev_portfolioforge/<slug>/<slug>.md`;
3. se publica en una URL estable guardada en `source_markdown_url`;
4. desde esa fuente publicada se derivan dos consumos:
   - **assistant** -> usa el markdown remoto completo;
   - **UI/DB** -> usa una representación resumida y estructurada del mismo contenido.

### Regla central

La UI **no debe copiar** el markdown literal. Debe mostrar solo los puntos clave, con compresión editorial y estructura operativa.

---

## Capacidades del producto

### Público

- landing con búsqueda guiada;
- catálogo de proyectos;
- detalle público tipo case study;
- búsqueda híbrida con filtros por categoría, cliente y tecnologías;
- explicaciones breves por resultado;
- assistant visible solo para sesiones elegibles cuando el proyecto lo permite.

### Admin

- CRUD de proyectos;
- CRUD de tecnologías;
- enriquecimiento estructurado mediante `project_profiles`;
- configuración privada de `source_markdown_url`;
- readiness de búsqueda y reembed;
- workflow persistido para case studies.

### Plataforma

- localización pública (`es`, `ca`, `en`, `de`);
- media optimizada por variantes y fallback;
- canonical publish por FTPS;
- localization backfill;
- recomposición del documento de búsqueda y embeddings.

---

## Stack real

- **Backend**: Go **1.25.x**, Echo v4, pgx/v5
- **Frontend**: React 19, TypeScript 5.9, Vite 7
- **Base de datos**: PostgreSQL 16
- **Extensiones DB**: `unaccent`, `pg_trgm`, `vector (pgvector)`
- **IA / embeddings / assistant**: OpenAI
- **Auth**: JWT + login local + Google opcional + verificación email + elegibilidad para assistant
- **Testing frontend**: Vitest + Playwright

---

## Cómo está organizado el repo

```text
cmd/                              entrypoints y subcomandos operativos
domain/                           lógica de dominio y puertos
infrastructure/                   adapters HTTP, PostgreSQL, assistant, localización, etc.
model/                            contratos de proyecto, usuario, media, workflow
sqlmigrations/                    evolución real del esquema
client/                           frontend React/Vite
docs/                             documentación autoritativa del sistema
```

---

## Quick start local

> Este bloque es una entrada rápida. Para operación real y completa, manda `docs/operacion/`.

### Prerrequisitos

- Go 1.25.x
- Node.js 20+
- PostgreSQL 16
- `pgvector` instalado en PostgreSQL

### 1. Preparar entorno

```bash
cp .env.example .env
mkdir -p images
```

Si también vas a usar overrides del frontend:

```bash
cp client/.env.example client/.env
```

### 2. Preparar base de datos

Sigue estos runbooks:

- `docs/operacion/02-entorno-servicios.md`
- `docs/operacion/03-db-bootstrap-migraciones.md`

### 3. Arrancar backend

```bash
go run ./cmd
```

Health check esperado:

```bash
curl http://localhost:8080/health
```

### 4. Arrancar frontend

```bash
cd client
npm install
npm run dev -- --host 0.0.0.0 --port 5173
```

### 5. Bootstrap inicial

Después del arranque:

- si no existe admin, sigue `docs/operacion/04-admin-usuarios.md`
- si vas a operar canonicals, sigue `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
- si vas a poblar runtime, sigue `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md`

---

## Documentación recomendada

### Punto de entrada general

- `docs/README.md`

### Documento rector del producto

- `docs/PRD.md`

### Operación real del repositorio

- `docs/operacion/README.md`

Ese directorio es la **fuente de verdad operativa** para:

- setup local;
- variables de entorno;
- migraciones;
- bootstrap del primer admin;
- workflow de case studies;
- publish FTPS;
- localization backfill;
- search readiness / refresh / reembed;
- troubleshooting.

### Reconstrucción desde cero

- `docs/reconstruccion-desde-cero/README.md`

---

## Validación y testing

### Backend

```bash
go test ./...
```

### Frontend

```bash
cd client
npm test
npm run build
```

### E2E visual

```bash
cd client
npm run test:e2e
```

---

## Para quién es útil este repo

- **Producto / arquitectura** -> `docs/PRD.md`
- **Onboarding técnico** -> este `README.md` + `docs/README.md`
- **Operación real** -> `docs/operacion/README.md`
- **Rebuild greenfield** -> `docs/reconstruccion-desde-cero/README.md`

---

## Regla documental importante

Si encuentras contradicciones entre este `README.md` y la documentación dentro de `docs/operacion/`, **manda `docs/operacion/`**.

Este archivo debe funcionar como:

- primera impresión del repositorio;
- resumen ejecutivo-técnico;
- mapa de entrada rápido.

La operación detallada vive en `/docs`.

---

## Qué leer primero si llegas nuevo

1. `docs/README.md`
2. `docs/PRD.md`
3. `docs/operacion/README.md`

Con eso ya puedes entender:

- qué es el sistema;
- cómo se organiza;
- cómo levantarlo;
- cómo operar el flujo editorial y runtime sin adivinar.
