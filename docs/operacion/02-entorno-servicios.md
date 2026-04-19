# Entorno y servicios requeridos

## Regla de lectura

Este documento concentra en un solo lugar las dependencias operativas del repo actual.

## Servicios externos

### Obligatorios para levantar el backend base

- PostgreSQL 16

### Obligatorios por feature

- OpenAI:
  - requerido para embeddings reales;
  - requerido para translations/backfill reales;
  - requerido para assistant por proyecto.
- Google Sign-In:
  - requerido si se quiere login público con Google.
- SMTP:
  - requerido para entrega real del OTP por email.
- FTPS:
  - requerido solo para `canonical-publish` o para la ruta legacy del workflow admin persistido cuando todavía publica canónicos.

## Variables de entorno del backend

### Base obligatoria

| Variable | Requerida | Uso real |
|---|---|---|
| `SERVER_PORT` | Sí | Puerto donde Echo expone la API |
| `ALLOWED_ORIGINS` | Sí | CORS general del backend |
| `ALLOWED_METHODS` | Sí | Métodos permitidos por CORS |
| `IMAGES_DIR` | Sí | Ruta local obligatoria validada al arrancar |
| `JWT_SECRET_KEY` | Sí | Firma de JWT |
| `DB_USER` | Sí | PostgreSQL |
| `DB_PASSWORD` | Sí | PostgreSQL |
| `DB_HOST` | Sí | PostgreSQL |
| `DB_PORT` | Sí | PostgreSQL |
| `DB_NAME` | Sí | PostgreSQL |
| `DB_SSL_MODE` | Sí | PostgreSQL |

Valores de desarrollo observados en ejemplos del repo:

- `SERVER_PORT=8080`
- `ALLOWED_ORIGINS=http://localhost:5173`
- `ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS`
- `IMAGES_DIR=./images`
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_NAME=portfolioforge`
- `DB_SSL_MODE=disable`

Valores de producción acordados para el despliegue en Hostinger VPS:

- `SERVER_PORT=8080`
- `ALLOWED_ORIGINS=https://portfolioforge.mlbautomation.com`
- `ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS`
- `IMAGES_DIR=/var/lib/portfolioforge/images`
- `DB_HOST=127.0.0.1`
- `DB_PORT=5432`
- `DB_NAME=portfolioforge`
- `DB_SSL_MODE=disable`

### Ajustes opcionales de pool DB

| Variable | Requerida | Uso real |
|---|---|---|
| `DB_MIN_CONN` | No | Override del pool mínimo |
| `DB_MAX_CONN` | No | Override del pool máximo |

Si no se definen, el código usa `min=3` y `max=100`.

### Search / AI

| Variable | Requerida | Cuándo |
|---|---|---|
| `ENABLE_SEMANTIC_SEARCH` | No | Activa refresh de embeddings y peso semántico |
| `OPENAI_API_KEY` | Condicional | Necesaria si `ENABLE_SEMANTIC_SEARCH=true`, si se usa assistant, o si se ejecuta localización automática |

Notas operativas:

- si `ENABLE_SEMANTIC_SEARCH=true` pero falta `OPENAI_API_KEY`, el backend arranca con warning y degrada búsqueda semántica a proveedores no-op/template;
- el assistant **no** degrada útilmente: sin `OPENAI_API_KEY` queda no disponible;
- `localization-backfill` falla si el traductor OpenAI no está configurado.

### Auth público y admin

| Variable | Requerida | Cuándo |
|---|---|---|
| `GOOGLE_CLIENT_ID` | Condicional | Login público con Google |

### OTP por email

| Variable | Requerida | Cuándo |
|---|---|---|
| `SMTP_HOST` | Condicional | SMTP real |
| `SMTP_PORT` | Condicional | SMTP real |
| `SMTP_USERNAME` | No, pero normalmente sí | SMTP autenticado |
| `SMTP_PASSWORD` | No, pero normalmente sí | SMTP autenticado |
| `EMAIL_FROM_ADDRESS` | Condicional | SMTP real |
| `EMAIL_FROM_NAME` | No | remitente legible |

Regla: si empiezas a configurar SMTP, el backend exige al menos `SMTP_HOST`, `SMTP_PORT` y `EMAIL_FROM_ADDRESS`.

### Workflow de case study y publish FTPS legacy

| Variable | Requerida | Uso |
|---|---|---|
| `PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS` | Sí solo para workflow legacy | allowlist de rutas fuente aceptadas |
| `PF_FTP_HOST` | Sí solo para publish legacy | host FTPS |
| `PF_FTP_PORT` | No | puerto FTPS, default `21` |
| `PF_FTP_USER` | Sí solo para publish legacy | usuario FTPS |
| `PF_FTP_PASSWORD` | Sí solo para publish legacy | password FTPS |
| `PF_PUBLIC_BASE` | Sí solo para publish legacy | base pública para construir la URL canónica |
| `PF_FTP_REMOTE_BASE` | No | raíz remota; default `/` |

Notas operativas:

- para el flujo estándar real de case studies, PortfolioForge necesita la URL remota ya publicada; esa publicación la hace manualmente personal externo del usuario en su host;
- FTPS no es requisito del flujo estándar de runtime/UI;
- si se usa la ruta legacy con `publish_canonical`, entonces sí aplican las variables `PF_FTP_*` y `PF_PUBLIC_BASE`.

## Variables del frontend

| Variable | Requerida | Uso real |
|---|---|---|
| `VITE_API_BASE_URL` | No | si falta, el frontend usa `http://localhost:8080`; en producción same-origin debe definirse como cadena vacía (`VITE_API_BASE_URL=`) |
| `VITE_GOOGLE_CLIENT_ID` | Condicional | render del botón de login público con Google |

Valores recomendados para Hostinger VPS:

- `VITE_API_BASE_URL=`
- `VITE_GOOGLE_CLIENT_ID` solo si login Google estará activo en producción

## Matriz rápida por operación

| Operación | DB | OpenAI | Google | SMTP | FTPS |
|---|---:|---:|---:|---:|---:|
| Backend base | Sí | No | No | No | No |
| Search semántico | Sí | Sí | No | No | No |
| Assistant | Sí | Sí | No | No | No |
| Signup/login local | Sí | No | No | Parcial |
| Login Google | Sí | No | Sí | No | No |
| OTP real | Sí | No | No | Sí | No |
| Localization backfill | Sí | Sí | No | No | No |
| Canonical publish legacy | No | No | No | No | Sí |
| Runtime/UI con URL ya publicada | Sí | Sí si hay localización | No | No | No |
| Case-study workflow persistido legacy completo | Sí | Sí si hay localización | No | No | Sí |

## Notas de autoridad

- esta tabla es la referencia operativa vigente;
- si otra documentación fuera de `docs/` enumera menos variables o defaults distintos, úsala solo como contexto legacy.
