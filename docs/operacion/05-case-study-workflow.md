# Flujo principal de case studies

## Propósito

Este es el **documento principal** para operar un case study nuevo o actualizado de principio a fin.

Su objetivo es dejar el happy path claro, lineal y con pocos saltos:

1. crear o corregir el canonical en la carpeta del caso;
2. publicarlo manualmente fuera de PortfolioForge hasta obtener `source_markdown_url`;
3. solo entonces crear o actualizar el runtime en la UI usando esa URL remota;
4. después ejecutar localización, readiness/reembed y verificación final.

## Clasificación documental para este flujo

### Documento principal

- `docs/operacion/05-case-study-workflow.md` *(este documento)*
  - explica el happy path completo;
  - define dónde hay que detenerse;
  - indica cuándo entrar a los documentos auxiliares.

### Documentos auxiliares

- `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
  - guía de autoría/corrección del canonical en la carpeta del caso.
- `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md`
  - referencia detallada para mapear el markdown publicado a campos runtime/UI.
- `docs/operacion/07-localization-backfill.md`
  - runbook auxiliar para regenerar locales derivadas.
- `docs/operacion/08-search-readiness-reembed.md`
  - runbook auxiliar para readiness, refresh y embeddings.

### Documento legacy / compatibilidad

- `docs/operacion/06-canonical-publish-ftps.md`
  - no es el camino principal;
  - se conserva solo como compatibilidad, diagnóstico u operación legacy.

## Regla central del flujo

PortfolioForge **no** arranca el runtime desde un archivo local.

Para UI/runtime, la fuente operativa es la **URL remota ya publicada** en `source_markdown_url`.

Secuencia obligatoria:

1. preparar `90. dev_portfolioforge/<slug>/<slug>.md` en la carpeta original del caso;
2. publicar manualmente ese canonical fuera de PortfolioForge;
3. verificar que la URL final exista y responda por HTTPS;
4. recién después usar esa URL en la UI/runtime.

Si esa URL todavía no existe, **el flujo debe detenerse**.

---

## Happy path operativo recomendado

## Paso 1 — Crear o corregir el canonical local

Objetivo:

- dejar correcto el archivo `90. dev_portfolioforge/<slug>/<slug>.md` en la carpeta real del caso de estudio.

Usa:

- `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`

Debe quedar claro que:

- la fuente editorial local vive en la carpeta del caso, no en una copia ad hoc dentro de PortfolioForge;
- el canonical base debe estar en castellano;
- ese archivo es el insumo que luego se publicará fuera del sistema.

Checklist mínimo al cerrar este paso:

- existe `90. dev_portfolioforge/<slug>/`;
- existe `90. dev_portfolioforge/<slug>/<slug>.md`;
- el slug del archivo es el que se usará también en runtime;
- el contenido ya es la versión editorial correcta para publicar.

## Paso 2 — Publicar manualmente fuera de PortfolioForge

Objetivo:

- obtener la URL remota final que se guardará como `source_markdown_url`.

Reglas:

- esta publicación manual ocurre fuera de PortfolioForge;
- la UI no resuelve este paso como flujo estándar;
- la ruta FTPS interna del repo no es obligatoria para el flujo principal.

Resultado esperado:

- existe una URL remota estable, por ejemplo bajo la convención `https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md`.

## Paso 3 — Verificar la URL remota antes de tocar runtime

Este es el punto de control más importante del flujo.

Antes de abrir la UI para crear o corregir runtime, verifica que:

- la URL final ya existe;
- responde correctamente por HTTPS;
- apunta al markdown canónico correcto del mismo slug;
- sigue en castellano y representa la versión editorial vigente.

### Regla de detención obligatoria

Si todavía **no** existe `source_markdown_url`, o si la URL existe pero aún no es la versión correcta, hay que detenerse aquí.

No corresponde todavía:

- crear el proyecto nuevo en la UI partiendo de un archivo local;
- actualizar runtime desde una copia local del canonical;
- dar por listo el case study;
- lanzar localización o reembed.

En otras palabras: **sin `source_markdown_url` válida no arranca el flujo runtime/UI**.

## Paso 4 — Crear o actualizar runtime en la UI usando la URL remota

Objetivo:

- poblar o corregir el proyecto runtime a partir de la URL ya publicada.

Usa:

- `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md`

Reglas operativas:

- si el proyecto ya existe, se parte de su `source_markdown_url` guardada;
- si el proyecto no existe, primero debe existir la URL y recién después se crea el runtime;
- la UI no debe leer el archivo local como fuente primaria de esta etapa;
- la UI/DB guarda una proyección resumida y estructurada, no una copia literal del markdown.

Qué se hace aquí:

1. confirmar o cargar `source_markdown_url`;
2. leer el markdown remoto publicado;
3. mapear nombre, descripción, categoría, client/context, perfil enriquecido, tecnologías y media según la guía runtime;
4. guardar o actualizar el proyecto;
5. validar que el runtime quedó coherente con la fuente remota.

## Paso 5 — Ejecutar localización derivada si aplica

Objetivo:

- regenerar `ca`, `en` y `de` desde la base en castellano ya corregida en runtime.

Usa:

- `docs/operacion/07-localization-backfill.md`

Solo debe hacerse cuando:

- el contenido base en `es` ya quedó correcto;
- `source_markdown_url` y runtime ya están alineados.

No debe ejecutarse todavía si el proyecto base sigue dudoso.

## Paso 6 — Verificar readiness, refresh y reembed

Objetivo:

- dejar consistente la capa de búsqueda y embeddings respecto al contenido ya actualizado.

Usa:

- `docs/operacion/08-search-readiness-reembed.md`

Secuencia mínima:

1. comprobar readiness;
2. corregir faltantes si los hubiera;
3. ejecutar reembed/refresh del proyecto;
4. validar búsqueda real con queries representativas.

## Paso 7 — Verificación final de “listo y terminado”

Un case study se considera listo cuando se cumplen **todas** estas condiciones:

1. el canonical local correcto existe en `90. dev_portfolioforge/<slug>/<slug>.md`;
2. la versión publicada manualmente ya existe en la URL final remota;
3. esa URL remota es la `source_markdown_url` real del proyecto;
4. el runtime en UI quedó creado o actualizado usando esa URL remota, no una fuente local;
5. nombre, descripción, categoría, `client_name`/contexto, perfil y tecnologías quedaron coherentes con el canonical publicado;
6. si aplicaba, la localización derivada quedó regenerada sin romper overrides manuales;
7. si aplicaba, readiness y reembed quedaron consistentes;
8. el proyecto puede verificarse correctamente en admin y en su consumo público/búsqueda según corresponda.

Si alguno de esos puntos falla, el case study todavía no debe darse por cerrado.

---

## Resumen lineal en una sola lista

1. crear o corregir `90. dev_portfolioforge/<slug>/<slug>.md`;
2. publicar manualmente ese canonical fuera de PortfolioForge;
3. obtener la URL remota final;
4. verificar que esa URL ya responde y será `source_markdown_url`;
5. si la URL no existe o no es correcta, detenerse;
6. crear o actualizar runtime/UI usando esa URL remota;
7. regenerar localización si aplica;
8. ejecutar readiness/reembed si aplica;
9. verificar admin, consumo público y búsqueda;
10. solo entonces marcar el case study como terminado.

---

## Workflow persistido disponible en el producto actual

El producto actual sigue teniendo un workflow admin persistido para steps como publish/import/localización/reembed sobre un case study con canonical ya existente.

Eso significa que siguen existiendo:

- UI: `/admin/settings` y `/admin/settings/case-studies`
- API: `GET /api/v1/admin/settings/case-study-workflow`

Payload de creación de run:

```json
{
  "source_path": "/ruta/permitida/90. dev_portfolioforge/mi-slug",
  "slug": "mi-slug",
  "run_localization_backfill": true,
  "run_reembed": true,
  "locales": ["ca", "en", "de"]
}
```

Steps reales implementados:

1. `resolve_source`
2. `publish_canonical` *(legacy / compatibilidad)*
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

### Cómo interpretar ese workflow hoy

- sirve como tooling persistido de la aplicación;
- **no** redefine el happy path principal descrito arriba;
- `publish_canonical` debe seguir tratándose como ruta legacy;
- `import_or_update_project` en el flujo operativo correcto solo debe considerarse después de que exista la URL remota publicada.

### Endpoints útiles

- crear run: `POST /api/v1/admin/settings/case-study-runs`
- ver run: `GET /api/v1/admin/settings/case-study-runs/:id`
- logs: `GET /api/v1/admin/settings/case-study-runs/:id/logs`
- confirmar: `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/confirm`
- iniciar: `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/start`
- retry: `POST /api/v1/admin/settings/case-study-runs/:id/steps/:step/retry`
- continuar: `POST /api/v1/admin/settings/case-study-runs/:id/resume`

## Precondiciones técnicas si se usa el workflow persistido

- backend levantado;
- admin autenticado;
- DB migrada incluyendo `20260418_0900_case_study_workflow_runs.sql`;
- variables `PF_CASE_STUDY_ALLOWED_SOURCE_ROOTS` configuradas;
- si se usará la ruta legacy FTPS, además `PF_FTP_HOST`, `PF_FTP_USER`, `PF_FTP_PASSWORD`, `PF_PUBLIC_BASE`;
- `OPENAI_API_KEY` si se va a ejecutar localización automática.

## Límites actuales

- no genera canonical desde una carpeta raw;
- el flujo principal real depende de publicación manual externa previa;
- la UI/runtime no debe tratar un archivo local como punto de partida principal;
- el workflow persistido legacy depende de FTPS y de una allowlist de filesystem externo al repo;
- si falla un step del workflow persistido, se reintenta ese step, no se reinventa manualmente el estado del run.
