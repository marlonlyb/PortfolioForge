# Documentación de PortfolioForge

## Propósito

Este directorio separa tres niveles de documentación para no mezclar producto, operación diaria y reconstrucción técnica completa.

## Mapa documental

### 1. Documento rector

- `docs/PRD.md`
  - define qué producto es PortfolioForge;
  - fija sus principios, capacidades y límites;
  - establece la relación entre markdown canónico, runtime, búsqueda y assistant.

### 2. Guías operativas

- `docs/operacion/README.md`
  - índice operativo oficial;
  - define el orden de uso de los runbooks vigentes;
  - deja explícito que README fuera de `docs/` puede estar desactualizado.
- `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
  - cómo generar o actualizar el markdown canónico en la carpeta estudiada original.
- `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md`
  - cómo convertir ese markdown en runtime UI/DB resumido, estructurado y verificable.
- `docs/operacion/01-setup-local.md`
  - setup local real, arranque y smoke test.
- `docs/operacion/02-entorno-servicios.md`
  - variables de entorno, servicios requeridos y dependencias operativas.
- `docs/operacion/03-db-bootstrap-migraciones.md`
  - bootstrap de base de datos y orden ejecutable de migraciones.
- `docs/operacion/04-admin-usuarios.md`
  - bootstrap del primer admin y operaciones de usuarios admin.
- `docs/operacion/05-case-study-workflow.md`
  - runbook del workflow admin persistido para case studies.
- `docs/operacion/06-canonical-publish-ftps.md`
  - runbook del publish canónico por FTPS.
- `docs/operacion/07-localization-backfill.md`
  - runbook de backfill de localización.
- `docs/operacion/08-search-readiness-reembed.md`
  - readiness de búsqueda, refresh y re-embed.
- `docs/operacion/09-troubleshooting.md`
  - fallos operativos comunes y diagnóstico inicial.

### 3. Blueprint de reconstrucción desde cero

- `docs/reconstruccion-desde-cero/README.md`
  - índice del pack técnico para recrear PortfolioForge desde cero.

## Matriz rápida de responsabilidades

| Documento | Responde qué | Cuándo usarlo |
|---|---|---|
| `docs/PRD.md` | Qué es el producto, qué resuelve y cómo se relacionan sus capas | Cuando necesitas entender el sistema a nivel producto y arquitectura conceptual |
| `docs/operacion/README.md` | Cuál es el mapa operativo vigente y qué runbook manda | Cuando vas a instalar, operar, reparar o publicar el sistema actual |
| `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md` | Cómo crear o actualizar el canonical correcto | Cuando estás trabajando sobre la fuente editorial del proyecto |
| `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md` | Cómo convertir esa fuente en UI/DB resumida y verificable | Cuando estás poblando o corrigiendo runtime |
| `docs/reconstruccion-desde-cero/README.md` | Cómo reconstruir todo el sistema desde cero | Cuando el objetivo es greenfield o replicación completa |

## Orden recomendado de lectura

### Para entender el producto

1. `docs/PRD.md`
2. `docs/README.md`

### Para operar el flujo vigente

1. `docs/PRD.md`
2. `docs/operacion/README.md`
3. el runbook específico que corresponda
4. `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md` y/o `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md` cuando aplique

### Para reconstruir el sistema desde cero

1. `docs/PRD.md`
2. `docs/reconstruccion-desde-cero/README.md`
3. resto del pack en el orden indicado por su índice

## Regla de mantenimiento

- si cambia el producto, primero se corrige `docs/PRD.md`;
- si cambia el flujo operativo, se corrigen las guías en `docs/operacion/`;
- si cambia la arquitectura objetivo o el orden de construcción greenfield, se corrige `docs/reconstruccion-desde-cero/`.

## Autoridad documental operativa

- dentro de este repositorio, la fuente de verdad para instalación y operación es `docs/operacion/`;
- `README.md` del repo y `client/README.md` pueden servir como contexto histórico, pero no deben prevalecer sobre los runbooks de `docs/operacion/` si existe contradicción;
- si detectas una contradicción, corrige `docs/operacion/` primero y deja el resto como follow-up documental fuera de alcance si no puede tocarse en esa tarea.
