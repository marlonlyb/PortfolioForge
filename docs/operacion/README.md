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
- publicación canónica por FTPS;
- backfill de localización;
- readiness, refresh y re-embed de búsqueda;
- troubleshooting operativo.

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
   - workflow admin persistido para publish/import/localización/reembed.
6. [`06-canonical-publish-ftps.md`](./06-canonical-publish-ftps.md)
   - publicación canónica por FTPS sin pasar por la UI.
7. [`07-localization-backfill.md`](./07-localization-backfill.md)
   - regeneración de locales derivadas desde español.
8. [`08-search-readiness-reembed.md`](./08-search-readiness-reembed.md)
   - readiness, refresh documental y embeddings.
9. [`09-troubleshooting.md`](./09-troubleshooting.md)
   - fallos frecuentes y diagnóstico inicial.

## Relación con otras piezas de `docs/`

- [`../PRD.md`](../PRD.md): contrato de producto y límites del sistema.
- [`./CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`](./CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md): cómo producir el markdown canónico.
- [`./PROJECT-RUNTIME-INGESTION-GUIDE.md`](./PROJECT-RUNTIME-INGESTION-GUIDE.md): cómo convertir ese canonical en runtime UI/DB.

## Orden recomendado

### Para levantar el repo

1. `02-entorno-servicios.md`
2. `03-db-bootstrap-migraciones.md`
3. `01-setup-local.md`
4. `04-admin-usuarios.md`

### Para operar ingestion/publicación

1. `CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
2. `06-canonical-publish-ftps.md`
3. `PROJECT-RUNTIME-INGESTION-GUIDE.md`
4. `05-case-study-workflow.md` si usas el flujo admin persistido
5. `07-localization-backfill.md`
6. `08-search-readiness-reembed.md`

### Para resolver incidentes

1. `09-troubleshooting.md`
2. el runbook específico afectado

## Regla explícita

- estos runbooks son autoritativos para implementar y operar el repo actual;
- si algo fuera de `docs/` sigue diciendo otra cosa, trátalo como material legacy o pendiente de sincronización.
