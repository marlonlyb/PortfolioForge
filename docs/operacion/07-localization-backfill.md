# Runbook de localization backfill

## Qué hace

Regenera localizaciones derivadas de un proyecto desde la base en español (`es`).

Aplica a campos públicos y de profile, incluyendo `client_name`.

## Comando

```bash
go run ./cmd localization-backfill [--project-id <uuid>] [--locale ca,en]
```

## Requisitos

- DB accesible (`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SSL_MODE`);
- `OPENAI_API_KEY` válida;
- proyecto base en español coherente antes de regenerar.

## Reglas de locales

Locales soportadas para backfill:

- `ca`
- `en`
- `de`

Si no pasas `--locale`, el backfill usa todas las locales derivadas soportadas.

## Modo recomendado

### 1. Smoke con un solo proyecto

```bash
go run ./cmd localization-backfill --project-id "<uuid-del-proyecto>" --locale ca,en
```

### 2. Verificación posterior

Validar:

- detalle público localizado;
- `client_name` localizado;
- admin translations consistente;
- sin pérdida de overrides manuales.

### 3. Batch completo

```bash
go run ./cmd localization-backfill
```

## Qué preserva

- las traducciones manuales no deben sobrescribirse;
- el sistema rellena solo campos auto gestionados.

## Cuándo ejecutarlo

- después de cambiar contenido base en español;
- después de importar/actualizar un canonical;
- antes de cerrar una corrección editorial que afecta payload público multilenguaje.

## Cuándo no ejecutarlo todavía

- si el contenido base en `es` sigue dudoso;
- si el canonical remoto y el runtime base todavía están desalineados;
- si la operación solo toca un override manual puntual que no quieres regenerar en lote.

## Riesgo conocido

- el comando valida DB obligatoria de forma explícita, pero la necesidad de `OPENAI_API_KEY` aparece al ejecutar el traductor; trátala como requisito efectivo del runbook.
