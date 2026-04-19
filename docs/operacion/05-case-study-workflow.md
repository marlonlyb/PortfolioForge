# Runbook del case-study workflow persistido

## Qué resuelve

Este workflow admin persiste ejecuciones de publish/import/localización/reembed para un case study que **ya** tiene canonical existente bajo `90. dev_portfolioforge/<slug>/`.

No genera canonical desde una carpeta cruda.

## Precondiciones

- backend levantado;
- admin autenticado;
- DB migrada incluyendo `20260418_0900_case_study_workflow_runs.sql`;
- variables `PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS`, `PF_FTP_HOST`, `PF_FTP_USER`, `PF_FTP_PASSWORD`, `PF_PUBLIC_BASE` configuradas;
- `OPENAI_API_KEY` si se va a ejecutar localización automática.

## Disponibilidad

### UI

- `/admin/settings`
- `/admin/settings/case-studies`

### API

- `GET /api/v1/admin/settings/case-study-workflow`

Si la configuración falta, el sistema responde workflow no disponible.

## Input del workflow

Payload de creación:

```json
{
  "source_path": "/ruta/permitida/90. dev_portfolioforge/mi-slug",
  "slug": "mi-slug",
  "run_localization_backfill": true,
  "run_reembed": true,
  "locales": ["ca", "en", "de"]
}
```

Reglas:

- `source_path` es obligatorio;
- la ruta debe quedar dentro de `PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS`;
- `slug` es opcional si la ruta ya resuelve un único slug válido;
- `locales` admite subset; por defecto la lógica soporta `ca`, `en`, `de`.

## Steps reales

1. `resolve_source`
2. `publish_canonical`
3. `import_or_update_project`
4. `localization_backfill`
5. `reembed`

Estados posibles:

- `pending`
- `blocked`
- `awaiting_confirmation`
- `running`
- `succeeded`
- `failed`
- `skipped`

## Semántica operativa

### 1. Resolve source

- normaliza la ruta;
- valida allowlist;
- resuelve slug, carpeta canonical y markdown canonical.

### 2. Publish canonical

- requiere confirmación manual;
- publica toda la carpeta del slug por FTPS;
- verifica luego la URL pública final.

### 3. Import or update project

- requiere confirmación manual;
- si el slug ya existe, actualiza nombre, descripción, categoría, `client_name`, `source_markdown_url` y media derivada;
- si no existe, crea el proyecto con `active=false`.

### 4. Localization backfill

- opcional;
- corre después del import;
- regenera locales desde español base;
- preserva overrides manuales.

### 5. Reembed

- opcional;
- corre después del import y después de localización si esa etapa está activa.

## Secuencia recomendada por UI

1. abrir `/admin/settings/case-studies`;
2. cargar `source_path` al canonical existente;
3. decidir si corre localización y reembed;
4. iniciar run;
5. confirmar `publish_canonical`;
6. ejecutar `publish_canonical`;
7. confirmar `import_or_update_project`;
8. ejecutar `import_or_update_project`;
9. continuar con `resume` o ejecutar el siguiente step pendiente;
10. revisar logs y estado final.

## Endpoints útiles

- crear run: `POST /api/v1/admin/settings/case-study-runs`
- ver run: `GET /api/v1/admin/settings/case-study-runs/:id`
- logs: `GET /api/v1/admin/settings/case-study-runs/:id/logs`
- confirmar: `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/confirm`
- iniciar: `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/start`
- retry: `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/retry`
- continuar: `POST /api/v1/admin/settings/case-study-runs/:id/resume`

## Criterio de éxito

- `canonical_url` poblada;
- `project_id` poblado;
- pasos terminales en `succeeded` o `skipped` según opciones;
- proyecto visible en `/admin/projects/:id`;
- si hubo reembed, readiness consistente y búsqueda refrescada.

## Límites actuales

- no genera canonical desde raw folder;
- depende de FTPS y allowlist de filesystem externo al repo;
- si falla un step, se reintenta ese step, no se reinventa manualmente el estado del run.
