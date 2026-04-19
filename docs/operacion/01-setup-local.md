# Setup local y arranque real

## Objetivo

Levantar backend, frontend y base de datos local con el contrato real del repo.

Nota de onboarding:

- este setup local vigente es manual/local sobre tu host; no depende de Docker ni de Docker Compose dentro de este repo;
- no uses `cmd/.env.example` como referencia de runtime: para variables vigentes manda la raíz `.env.example`, `client/.env.example` y [`02-entorno-servicios.md`](./02-entorno-servicios.md).

## Prerrequisitos

- Go con versión compatible con `go.mod` (`go 1.25.x`).
- Node.js 20+.
- PostgreSQL 16.
- Extensión `pgvector` instalada en el host PostgreSQL.
- `npm` disponible para `client/`.

## 1. Preparar variables de entorno

Usa como base los ejemplos del repo:

- raíz: `.env.example`
- frontend: `client/.env.example`

`cmd/.env.example` existe solo como compatibilidad documental legacy y no debe usarse para bootstrap local.

Variables mínimas del backend que deben existir antes de arrancar:

- `SERVER_PORT`
- `ALLOWED_ORIGINS`
- `ALLOWED_METHODS`
- `IMAGES_DIR`
- `JWT_SECRET_KEY`
- `DB_USER`
- `DB_PASSWORD`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_SSL_MODE`

`ENABLE_SEMANTIC_SEARCH`, `OPENAI_API_KEY`, `GOOGLE_CLIENT_ID`, SMTP y variables `PF_*` son condicionales según el flujo; ver [`02-entorno-servicios.md`](./02-entorno-servicios.md).

## 2. Preparar directorio de imágenes local

El backend valida `IMAGES_DIR`, y el ejemplo apunta a `./images`.

Si no existe, créalo desde la raíz del repo:

```bash
mkdir -p images
```

## 3. Preparar base de datos

Sigue [`03-db-bootstrap-migraciones.md`](./03-db-bootstrap-migraciones.md) y vuelve aquí cuando la DB ya esté migrada.

## 4. Arrancar backend

Desde la raíz del repo:

```bash
go run ./cmd
```

Checks mínimos:

```bash
curl http://localhost:8080/health
```

Respuesta esperada: HTTP 200 con `PortfolioForge is running`.

## 5. Arrancar frontend

Desde `client/`:

```bash
npm install
npm run dev -- --host 0.0.0.0 --port 5173
```

Notas:

- el frontend cae por defecto a `http://localhost:8080` si `VITE_API_BASE_URL` no está definido;
- para desarrollo local normal, `ALLOWED_ORIGINS` debe incluir `http://localhost:5173`.

## 6. Smoke test mínimo

1. Abrir `http://localhost:5173/`.
2. Verificar que la landing cargue.
3. Abrir `http://localhost:5173/login`.
4. Si ya existe admin, abrir también `http://localhost:5173/admin/login` y validar acceso.
5. Si todavía no existe admin, seguir [`04-admin-usuarios.md`](./04-admin-usuarios.md).

## 7. Arranque mínimo por escenario

### Solo catálogo/search léxico-fuzzy

Requiere:

- backend válido;
- frontend válido;
- DB migrada;
- `ENABLE_SEMANTIC_SEARCH=false` o vacío.

### Search semántico real

Además requiere:

- `ENABLE_SEMANTIC_SEARCH=true`
- `OPENAI_API_KEY`

### Assistant por proyecto

Además requiere:

- `OPENAI_API_KEY`
- proyecto con `source_markdown_url` accesible por HTTPS
- usuario autenticado elegible

### OTP por email real

Además requiere SMTP válido.

## Cierre operativo

Después del arranque inicial, los siguientes runbooks suelen ser:

- [`04-admin-usuarios.md`](./04-admin-usuarios.md)
- [`05-case-study-workflow.md`](./05-case-study-workflow.md)
- [`08-search-readiness-reembed.md`](./08-search-readiness-reembed.md)
