# Operación de PortfolioForge

## Propósito

Este directorio es la **fuente de verdad operativa** del repositorio actual.

Si `README.md` del repo, `client/README.md` o notas fuera de `docs/` contradicen algo de aquí, **manda `docs/operacion/`**.

Nota para contribuidores nuevos:

- la documentación operativa usa el dominio `project`, aunque todavía existan nombres legacy como `products` o `brand` en partes del sistema;
- archivos fuera de `docs/`, incluyendo `client/README.md` y `cmd/.env.example`, pueden existir por contexto o compatibilidad, pero no mandan sobre `docs/operacion/`.

## Qué cubre

- instalación local real;
- variables de entorno y servicios requeridos;
- bootstrap de base de datos y migraciones;
- bootstrap del primer admin y gestión de usuarios;
- workflow persistido de case studies;
- publicación canónica manual externa como flujo estándar;
- tooling FTPS legacy/opcional de compatibilidad;
- backfill de localización;
- readiness, refresh y re-embed de búsqueda;
- troubleshooting operativo.

## Tooling operador de OpenCode

Además de los runbooks, este repo puede operarse con comandos globales de OpenCode ejecutados **desde PortfolioForge** para que reutilicen `docs/` como fuente de verdad:

- `/pf-canonical-create <directory>`
- `/pf-ui-create <source_markdown_url>`
- `/pf-ui-update <source_markdown_url>`

Interpretación correcta:

- no son capacidades del producto ni endpoints de la aplicación;
- son wrappers de conveniencia sobre el workflow canonical-first ya documentado;
- `/pf-canonical-create` cubre la etapa canonical/editorial y actualiza el inventario operador local `.atl/case-study-local-index.md`;
- `/pf-ui-create` y `/pf-ui-update` arrancan siempre desde la URL remota publicada, nunca desde un archivo local;
- `publish_canonical` / FTPS sigue siendo tooling legacy/opcional de compatibilidad.

Inventario operador local:

- `.atl/case-study-local-index.md` vive fuera del producto/runtime, está ignorado por git y guarda solo cinco campos por caso: `slug`, `project name`, `source repo local path`, `canonical local path`, `source_markdown_url`.
- dentro de ese archivo, `source_markdown_url` puede aparecer primero como **URL objetivo prevista por convención**; solo se considera URL operativa real después de publicación manual externa y verificación HTTPS exitosa.

## Índice de runbooks

1. [`01-setup-local.md`](./01-setup-local.md)
   - arranque local completo y smoke test.
2. [`02-entorno-servicios.md`](./02-entorno-servicios.md)
   - variables de entorno, dependencias y modos requeridos/opcionales.
3. [`03-db-bootstrap-migraciones.md`](./03-db-bootstrap-migraciones.md)
   - creación de DB, extensiones y orden ejecutable de migraciones.
4. [`04-admin-usuarios.md`](./04-admin-usuarios.md)
   - primer admin y operaciones de usuarios.
5. [`05-case-study-workflow.md`](./05-case-study-workflow.md)
   - **documento principal** del happy path de case studies: canonical local → publicación manual externa → runtime/UI con `source_markdown_url` → localización → readiness/reembed → cierre.
6. [`06-canonical-publish-ftps.md`](./06-canonical-publish-ftps.md)
   - tooling FTPS **legacy/opcional** sin pasar por la UI.
7. [`07-localization-backfill.md`](./07-localization-backfill.md)
   - runbook **auxiliar** para regeneración de locales derivadas desde español.
8. [`08-search-readiness-reembed.md`](./08-search-readiness-reembed.md)
   - runbook **auxiliar** para readiness, refresh documental y embeddings.
9. [`09-troubleshooting.md`](./09-troubleshooting.md)
   - fallos frecuentes y diagnóstico inicial.
10. [`10-hostinger-vps-deployment.md`](./10-hostinger-vps-deployment.md)
   - runbook autoritativo para despliegue productivo mínimo en Hostinger VPS.

## Relación con otras piezas de `docs/`

- [`../PRD.md`](../PRD.md): contrato de producto y límites del sistema.
- [`./05-case-study-workflow.md`](./05-case-study-workflow.md): documento principal para operar un case study de punta a punta.
- [`./CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`](./CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md): guía auxiliar para producir o corregir el markdown canónico.
- [`./PROJECT-RUNTIME-INGESTION-GUIDE.md`](./PROJECT-RUNTIME-INGESTION-GUIDE.md): referencia auxiliar para convertir ese canonical remoto en runtime UI/DB.
- [`./06-canonical-publish-ftps.md`](./06-canonical-publish-ftps.md): ruta legacy/opcional de compatibilidad.

## Orden recomendado

### Para levantar el repo

1. `02-entorno-servicios.md`
2. `03-db-bootstrap-migraciones.md`
3. `01-setup-local.md`
4. `04-admin-usuarios.md`

### Para desplegar en producción (Hostinger VPS)

1. `02-entorno-servicios.md`
2. `03-db-bootstrap-migraciones.md`
3. `10-hostinger-vps-deployment.md`
4. `04-admin-usuarios.md` si también necesitas bootstrap del primer admin

### Para operar ingestion/publicación

1. `05-case-study-workflow.md`
2. `CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md` para preparar el canonical local
3. publicación manual externa al host final del usuario para obtener `source_markdown_url`
4. `PROJECT-RUNTIME-INGESTION-GUIDE.md` para cargar/actualizar runtime desde esa URL remota
5. `07-localization-backfill.md` solo cuando ya corresponda regenerar locales
6. `08-search-readiness-reembed.md` solo cuando ya corresponda cerrar readiness/reembed
7. `06-canonical-publish-ftps.md` solo si necesitas compatibilidad/diagnóstico legacy

## Lectura recomendada para case studies

### Principal

1. `05-case-study-workflow.md`

### Auxiliares

2. `CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
3. `PROJECT-RUNTIME-INGESTION-GUIDE.md`
4. `07-localization-backfill.md`
5. `08-search-readiness-reembed.md`

### Legacy

6. `06-canonical-publish-ftps.md`

## Regla de detención para este flujo

- si todavía no existe la URL remota final que será `source_markdown_url`, el operador debe detenerse antes de crear o actualizar runtime/UI;
- la UI/runtime no parte de un archivo local como fuente primaria;
- un case study no está cerrado hasta verificar canonical local, URL remota publicada, runtime actualizado, localización (si aplica) y readiness/reembed (si aplica).

### Para resolver incidentes

1. `09-troubleshooting.md`
2. el runbook específico afectado

## Regla explícita

- estos runbooks son autoritativos para implementar y operar el repo actual;
- si algo fuera de `docs/` sigue diciendo otra cosa, trátalo como material legacy o pendiente de sincronización.
