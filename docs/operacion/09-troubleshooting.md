# Troubleshooting operativo

## Regla base

Antes de asumir bug de código, valida en este orden:

1. entorno;
2. migraciones;
3. variables;
4. disponibilidad de servicios externos;
5. datos del proyecto;
6. recién después, código.

## Backend no arranca

### Síntoma

`go run ./cmd` corta al inicio con error de env.

### Revisar

- faltan `SERVER_PORT`, `ALLOWED_ORIGINS`, `ALLOWED_METHODS`, `IMAGES_DIR`, `JWT_SECRET_KEY` o DB vars;
- `IMAGES_DIR` apunta a una ruta no creada;
- `.env` no está presente o no coincide con el entorno actual.

## `/health` responde pero rutas fallan con `500`

### Revisar

- migraciones incompletas;
- especialmente tablas/search/localization/assistant faltantes;
- ver [`03-db-bootstrap-migraciones.md`](./03-db-bootstrap-migraciones.md).

## Frontend no conecta al backend

### Revisar

- backend en `http://localhost:8080`;
- `VITE_API_BASE_URL` correcta si se overrideó;
- `ALLOWED_ORIGINS` incluye `http://localhost:5173`;
- consola de red del navegador.

## Login admin rechaza un usuario válido

### Revisar

- el usuario realmente tiene `is_admin = true`;
- estás usando `/admin/login` o `POST /api/v1/admin/login`;
- no estás intentando convertir el primer admin desde una UI sin admin previo.

## OTP no entrega correos

### Revisar

- `SMTP_HOST`, `SMTP_PORT`, `EMAIL_FROM_ADDRESS` configurados;
- credenciales SMTP válidas;
- si SMTP no está configurado, el sistema no puede entregar OTP reales.

## Workflow de case study aparece deshabilitado

### Revisar

- `PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS`
- `PF_FTP_HOST`
- `PF_FTP_USER`
- `PF_FTP_PASSWORD`
- `PF_PUBLIC_BASE`

Sin eso, `/api/v1/admin/settings/case-study-workflow` informa no configurado.

## `canonical-publish` falla

### Revisar

- falta alguna `PF_*` requerida;
- `--case-dir` no contiene un slug canonical válido;
- hay múltiples slugs y faltó `--slug`;
- FTPS/TLS/credenciales;
- la URL pública final no responde por HTTPS.

## El import del workflow creó proyecto pero quedó no visible

### Causa esperada

El importador crea proyectos nuevos con `active=false`.

### Acción

- revisar `/admin/projects/:id`;
- completar enrichment/tecnologías si falta;
- publicar/activar explícitamente desde admin cuando corresponda.

## Localización falla

### Revisar

- `OPENAI_API_KEY`;
- locale soportada (`ca`, `en`, `de`);
- proyecto base en español válido;
- no esperar cambios de ranking multilingüe: la localización aplica display fields, no rediseña el retrieval.

## Search/reembed no refleja cambios

### Revisar

- readiness del proyecto;
- si el contenido realmente cambió;
- `ENABLE_SEMANTIC_SEARCH` y `OPENAI_API_KEY` si esperas embedding nuevo;
- usar `POST /api/v1/admin/projects/:id/reembed` después de enrichment/import/localización.

## Assistant no aparece

### Revisar

- el proyecto tiene `source_markdown_url` cargada;
- el payload público da `assistant_available=true`;
- el usuario está autenticado y es elegible;
- en usuarios públicos por email, el perfil debe quedar completo (`full_name`, `company`).

## Assistant aparece pero responde error

### Revisar

- `OPENAI_API_KEY`;
- accesibilidad HTTPS del markdown remoto;
- endpoint correcto: `POST /api/v1/private/projects/:slug/assistant/messages`.

## Contradicción entre docs

### Regla

- si hay conflicto entre `docs/operacion/` y README fuera de `docs/`, toma `docs/operacion/` como referencia vigente;
- deja el resto como follow-up documental si está fuera del alcance permitido.
