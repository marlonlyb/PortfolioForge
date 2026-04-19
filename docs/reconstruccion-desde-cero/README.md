# Pack de reconstrucción desde cero — PortfolioForge

## Objetivo

Este directorio define el mínimo documental necesario para reconstruir PortfolioForge desde cero sin depender de conocimiento tácito del repo actual. Toma como fuente normativa el PRD vigente, las guías operativas de markdown/runtime y el esquema real expresado en `sqlmigrations/`.

## Cómo leer este pack

1. `01-alcance-producto.md` fija qué producto se debe construir y qué no.
2. `02-arquitectura-sistema.md` describe la forma general del sistema y sus bounded contexts.
3. `03-modelo-dominio-datos.md` define entidades, relaciones y reglas de persistencia.
4. `04-contratos-backend-api.md` fija las capacidades y contratos HTTP/backend mínimos.
5. `05-superficies-frontend.md` traduce el producto a rutas, pantallas y módulos UI.
6. `06-flujo-canonico-y-publicacion.md` documenta la cadena editorial obligatoria desde el markdown canónico hasta assistant/search.
7. `07-capacidades-transversales.md` concentra búsqueda, assistant, localización, auth/admin y workflow operativo.
8. `08-roadmap-implementacion.md` propone orden de construcción para un greenfield realista.

## Principios de uso

- Este pack habla en dominio `project`, aunque reconoce persistencia legacy `products`/`brand` cuando afecta compatibilidad.
- El idioma base del sistema es castellano (`es`).
- El markdown canónico es fuente editorial primaria; la UI/DB es representación runtime resumida.
- El assistant solo responde sobre markdown remoto publicado en `source_markdown_url`.

## Fuentes de verdad utilizadas

- `docs/PRD.md`
- `docs/operacion/PROJECT-RUNTIME-INGESTION-GUIDE.md`
- `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
- `sqlmigrations/20260410_1200_store_mvp_base.sql` → `sqlmigrations/20260418_2300_search_client_name_alignment.sql`

## Resultado esperado

Si se sigue este pack, el sistema reconstruido debe poder:

- publicar un portfolio navegable por catálogo, detalle y búsqueda;
- administrar proyectos, tecnologías, traducciones, media y settings;
- ingerir proyectos desde markdown canónico;
- ofrecer assistant autenticado y grounded por proyecto;
- mantener coherencia entre markdown, runtime, localización y búsqueda.
