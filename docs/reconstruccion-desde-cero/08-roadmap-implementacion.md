# Roadmap de implementación y orden de construcción

## Objetivo

Definir un orden de build que reduzca retrabajo y preserve la relación correcta entre dominio, runtime y editorial.

## Fase 1 — Fundaciones

1. crear esquema base de usuarios y proyectos;
2. introducir `slug`, `category`, `active`, `client_name`, `status`, `featured`;
3. montar backend Go con health/config/routing base;
4. montar frontend React con router y layouts público/admin;
5. habilitar auth admin mínimo.

## Fase 2 — Runtime de proyectos

1. implementar `projects/products` CRUD admin;
2. implementar detalle/listado público por `slug`;
3. agregar `project_profiles`;
4. agregar tecnologías y relación N:N;
5. agregar `project_media` y serialización `media + images`.

## Fase 3 — Editorial coherente

1. fijar estructura del markdown canónico;
2. implementar importación/actualización desde canonical a runtime;
3. persistir `source_markdown_url`;
4. validar compresión editorial en la UI admin.

## Fase 4 — Discoverability

1. crear `project_search_documents`;
2. componer documento FTS + trigram;
3. integrar embeddings;
4. exponer búsqueda pública y search admin readiness;
5. agregar explanations por resultado.

## Fase 5 — Assistant y auth pública

1. signup/login local;
2. verificación email;
3. login Google opcional según prioridad;
4. elegibilidad para assistant;
5. endpoint privado de assistant grounded en markdown remoto.

## Fase 6 — Localización

1. crear `project_localizations`;
2. definir campos traducibles;
3. backfill automático desde `es`;
4. overrides manuales por locale/campo;
5. aplicación de locale en respuestas públicas.

## Fase 7 — Operación completa

1. site settings públicos/admin;
2. workflow persistido de case study runs;
3. herramientas de publicación canónica y resume/retry;
4. endurecimiento de permisos, rate limits y observabilidad.

## Dependencias críticas

- assistant depende de auth + `source_markdown_url` + canonical publicado;
- búsqueda buena depende de enrichment consistente y tecnologías resueltas;
- localización depende de base `es` estabilizada;
- workflow admin depende de tener importación, publicación y reembed ya implementados.

## Definición de listo del rebuild

Se considera reconstruido cuando un editor puede tomar un canonical existente, publicarlo, importarlo al runtime, traducirlo, reindexarlo y mostrarlo en público con assistant operativo y búsqueda útil.
