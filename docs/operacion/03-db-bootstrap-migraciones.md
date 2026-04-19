# Base de datos, bootstrap y migraciones

## Objetivo

Dejar PostgreSQL listo para que el backend pueda arrancar y operar el esquema real del repo.

## 1. Requisitos previos

- PostgreSQL 16 levantado.
- `psql` disponible.
- paquete `pgvector` instalado en el host PostgreSQL.

En Ubuntu/Debian, el propio repo ya documenta:

```bash
sudo apt-get install -y postgresql-16-pgvector
```

## 2. Crear la base local

La base esperada por defecto es `portfolioforge`.

Ejemplo:

```bash
createdb portfolioforge
```

## 3. Aplicar todas las migraciones en orden lexical

La forma más segura y mantenible es ejecutar **todos** los archivos de `sqlmigrations/` en orden lexical:

```bash
for f in sqlmigrations/*.sql; do psql -v ON_ERROR_STOP=1 -d portfolioforge -f "$f"; done
```

Orden actual del repo:

1. `20240617_2206_create_user.sql`
2. `20240624_1609_create_products.sql`
3. `20240625_2312_create_purchase_order.sql`
4. `20240627_1503_create_invoice.sql`
5. `20240627_1505_create_invoice_details.sql`
6. `20260410_1200_store_mvp_base.sql`
7. `20260411_0900_portfolioforge_extension.sql`
8. `20260412_0100_search_extensions.sql`
9. `20260412_0200_search_schema.sql`
10. `20260412_0300_search_backfill.sql`
11. `20260412_0400_search_embedding_text.sql`
12. `20260413_0900_project_localizations.sql`
13. `20260413_1200_project_media_variants.sql`
14. `20260414_1000_site_settings.sql`
15. `20260414_1100_project_assistant_chat.sql`
16. `20260415_1000_authenticated_project_assistant.sql`
17. `20260415_1100_local_email_verification.sql`
18. `20260416_0900_admin_user_management.sql`
19. `20260416_0900_layered_project_detail_ui.sql`
20. `20260416_1400_standard_public_auth.sql`
21. `20260418_0900_case_study_workflow_runs.sql`
22. `20260418_2200_compact_project_profile_lists.sql`
23. `20260418_2300_search_client_name_alignment.sql`

## 4. Validaciones mínimas post-migración

### Extensiones

```bash
psql -d portfolioforge -c "SELECT extname FROM pg_extension WHERE extname IN ('unaccent','pg_trgm','vector') ORDER BY extname;"
```

### Tablas operativas clave

```bash
psql -d portfolioforge -c "SELECT to_regclass('public.users'), to_regclass('public.products'), to_regclass('public.project_profiles'), to_regclass('public.project_search_documents'), to_regclass('public.project_localizations');"
```

## 5. Bootstrap inicial recomendado

Después de migrar:

1. arranca backend y frontend;
2. crea o promueve el primer admin siguiendo [`04-admin-usuarios.md`](./04-admin-usuarios.md);
3. crea tecnologías antes de intentar enrichment completo;
4. usa [`08-search-readiness-reembed.md`](./08-search-readiness-reembed.md) para validar indexación.

## 6. Criterio operativo

- no des por válida una instalación solo porque `go run ./cmd` arranca;
- si faltan migraciones, rutas públicas o admin pueden fallar con `500` aunque el proceso esté vivo.

## Riesgos conocidos

- el directorio `sqlmigrations/` mezcla piezas portfolio y legado e-commerce; hoy la forma respaldada por el repo es aplicar todo el set en orden lexical;
- si en el futuro se separa ese historial, deberá documentarse como follow-up fuera de este runbook.
