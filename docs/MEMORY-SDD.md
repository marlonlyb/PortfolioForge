# MEMORY-SDD

## Propósito

Este documento resume el trabajo acumulado desde la inicialización SDD del proyecto hasta el estado actual. No reemplaza al `PRD`; funciona como memoria operativa del proceso, de los changes realizados y de las decisiones recientes sobre documentación, estandarización y flujo de ingestión.

## 1. Inicialización SDD y contexto detectado

- Se ejecutó `/sdd-init` en modo `engram`.
- Se detectó el stack real del proyecto:
  - Go 1.20 + Echo v4
  - PostgreSQL + pgx/v5
  - React 19 + TypeScript 5.9 + Vite 7
- Se registró el contexto arquitectónico base del repo.
- Se estableció el uso de `docs/` como ubicación principal de documentación del producto, manteniendo `README.md` como guía general del repositorio.

## 2. Limpieza de herencia e-commerce y redefinición del dominio

Al inicio, el proyecto aún arrastraba nombres y estructuras heredadas de un catálogo/product store.

Trabajo realizado:

- eliminación progresiva de documentación heredada de tienda/e-commerce;
- redefinición del dominio funcional hacia `project`, `client_name`, `technologies`, `project_profiles`;
- aceptación explícita de deuda transicional en naming interno (`product`, `brand`) para no bloquear el avance del MVP.

Resultado:

- el producto quedó reposicionado como portfolio técnico;
- la documentación dejó de hablar de catálogo comercial y pasó a hablar de projects, case studies y search evidence.

## 3. Consolidación del producto PortfolioForge

Se formalizó que PortfolioForge debía ser:

- una plataforma de portfolio interactivo;
- con búsqueda por evidencia real como interacción principal;
- con explicaciones breves por resultado, limitadas por evidencia;
- sin chatbot libre.

También se fijó que:

- la landing es la entrada principal del producto;
- la barra de búsqueda es la pieza central de la experiencia pública;
- la propuesta visual debe ser oscura, técnica y orientada a developer.

## 4. Change SDD: `evidence-based-project-search`

### Planeamiento

- Se abrió el change `evidence-based-project-search`.
- Se generaron proposal, design, spec y tasks.

### Implementación principal

Se implementó:

- búsqueda híbrida con filtros estructurados;
- PostgreSQL FTS;
- fuzzy matching con `pg_trgm`;
- capa semántica con `pgvector`;
- servicio de ranking;
- explicaciones resumidas por resultado basadas en evidencia del proyecto;
- endpoint público de búsqueda;
- UI frontend de búsqueda con barra, filtros, resultados y cards.

### Verificación y cierre

- Se corrigieron gaps detectados en verify.
- Se validó backend y typecheck frontend.
- El change fue archivado como completado.

## 5. Base local, infraestructura y entorno real

Durante el trabajo inicial se normalizó el entorno de desarrollo:

- se creó y fijó la base `portfolioforge` como DB local correcta;
- se evitó reutilizar la base del proyecto original (`proyectoemlb`);
- se instalaron / habilitaron extensiones `unaccent`, `pg_trgm` y `vector`;
- se aplicaron migraciones del proyecto;
- se dejó el backend funcionando localmente en `:8080`.

Esto fue clave porque varias decisiones posteriores —especialmente búsqueda semántica y re-embedding— dependían de una base local coherente.

## 6. README y documentación de repositorio

Se reescribió `README.md` para reflejar:

- estado real del producto;
- stack y variables de entorno correctas;
- rutas públicas y admin;
- dependencias de PostgreSQL;
- flujo manual básico de prueba.

La documentación dejó de ser un reflejo histórico del proyecto anterior y pasó a describir el producto actual con precisión operativa.

## 7. Change SDD: `portfolio-domain-enrichment`

### Objetivo

Enriquecer el dominio del portfolio con:

- OpenAI como proveedor real;
- CRUD admin de tecnologías;
- `project_profiles` para enriquecer cada proyecto;
- re-composición y re-embedding del documento de búsqueda.

### Implementación principal

Se implementó:

- `OpenAIEmbeddingProvider` con `text-embedding-3-small`;
- `OpenAIExplanationProvider` con `gpt-4o-mini`;
- fallbacks seguros a providers no-op/template;
- CRUD admin de tecnologías;
- endpoint `PUT /api/v1/admin/projects/:id/enrichment`;
- enriquecimiento del formulario admin de proyectos.

### Ajustes posteriores detectados en verify

Se corrigieron varios gaps importantes:

- regeneración de spec faltante;
- `GET /api/v1/admin/technologies/:id`;
- atomicidad de enrichment + re-embedding;
- error de embedding pasa a fallar la request;
- embeddings calculados desde texto compuesto real y no desde `tsvector`;
- inclusión de `problem_statement` en la composición;
- endurecimiento del prompt de explicación a exactamente una oración;
- ampliación de tests.

### Resultado

- el change quedó archive-ready y luego archivado;
- la app ya soporta proyectos enriquecidos reales y búsqueda semántica funcional.

## 8. Estandarización del modelo de proyecto en admin

Con el enriquecimiento ya funcionando, se consolidó qué campos existen realmente hoy en la UI/admin:

### Campos base

- Title
- Summary / Description
- Category
- Client / Context
- Main images
- Published

### Campos enriquecidos

- Technologies
- Business Goal
- Problem Statement
- Solution Summary
- Architecture
- AI Usage
- Integrations
- Technical Decisions
- Challenges
- Results
- Metrics
- Timeline

Esta aclaración fue importante porque parte del vocabulario del código seguía heredado (`name`, `description`, `brand`, `images`, `active`), mientras que la documentación necesitaba hablar en términos de producto y UI actual.

## 9. Reglas de serialización y contratos reales del formulario

Se verificó y documentó que el admin actual serializa así:

- `integrations`, `technical_decisions`, `challenges`, `results`, `timeline` → una línea por ítem;
- `metrics` → formato `key: value`, una línea por métrica.

Esto no es solo una convención editorial: está alineado con los serializadores reales del frontend y con el contrato JSON que el backend normaliza al guardar `project_profiles`.

## 10. Gestión de tecnologías y mapping operativo

Se consolidó que las tecnologías deben existir primero en `/admin/technologies`.

Campos actuales de tecnología:

- `name`
- `category`
- `icon`
- `brand color`

Aprendizaje relevante:

- la calidad del portfolio depende de un catálogo de tecnologías consistente;
- si el catálogo es pobre o ambiguo, también se degradan filtros, embeddings y recuperación por stack.

## 11. Evolución de estándares documentales

En trabajo reciente se reforzaron criterios para que la documentación no sea aspiracional sino contractual:

- los docs deben reflejar campos y flujos reales del código;
- las guías deben explicar la UI disponible hoy, no una UI futura imaginada;
- la terminología oficial debe ser de portfolio, aunque el código aún tenga nombres heredados;
- se deben evitar duplicaciones entre `README.md` y `docs/`, manteniendo continuidad entre ambos.

Esto llevó a una actualización explícita de estándares de documentación dentro del proyecto y a una revisión de qué documentos necesitaban ser refrescados.

## 12. Trabajo reciente de análisis externo de proyectos

Se empezó a formalizar un flujo para trabajar con proyectos del mundo real desde evidencia externa al sistema.

Líneas de trabajo recientes:

- análisis de proyectos reales desde repositorio/carpeta;
- extracción de señales técnicas, de dominio e integración;
- mapping de tecnologías detectadas hacia el catálogo del admin;
- definición de una fuente markdown intermedia para describir proyectos antes de cargarlos en UI.

Caso usado como referencia metodológica:

- análisis de un caso PRECOR vinculado a grúa / CAN Bus para validar cómo derivar título técnico, client/context, arquitectura, integraciones, desafíos y resultados a partir de evidencia real.

Importante:

- este caso sirve como ejemplo del proceso;
- no debe convertir la documentación del producto en documentación específica de un cliente.

## 13. Planificación de ingestión basada en markdown fuente

Se definió con más claridad una dirección editorial/técnica:

- antes de cargar en admin, conviene generar un markdown fuente estable;
- ese markdown debe mapear 1:1 con los campos reales de PortfolioForge;
- debe incluir convenciones de naming, imágenes y serialización;
- debe permitir revisión humana antes de publicación.

La consecuencia práctica es que PortfolioForge ya no se entiende solo como “admin + búsqueda”, sino también como un sistema que necesita un **pipeline de producción de contenido técnico**.

## 14. Convenciones recientes de naming e imágenes

Se reforzó la convención para nuevos proyectos fuente:

- nombres cortos y técnicos;
- sin cliente/empresa dentro del título;
- título basado en la tecnología o concepto central;
- imágenes con convención pública/canónica por defecto `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_<low|medium|high>.webp`;
- mínimo recomendado de 5 imágenes;
- `Main images` usando normalmente variantes `_medium`.

Esta convención está pensada para alinear la preparación del material con el formulario admin existente y con la experiencia visual del portfolio.

## 15. Gestión/lifecycle dentro de campos disponibles

También se aclaró una restricción importante:

- la UI actual no tiene campos dedicados para PM, delivery management o lifecycle.

Por eso, la información de ejecución debe incrustarse en los campos ya existentes, por ejemplo:

- alcance/objetivo → `Business Goal`;
- problema real y restricciones → `Problem Statement`;
- hitos y rollout → `Timeline`;
- tradeoffs y decisiones de ejecución → `Technical Decisions`;
- impacto y outcomes → `Results` + `Metrics`.

## 16. Estado funcional actual

Hoy el sistema permite:

- búsqueda pública por evidencia real;
- resultados con explicación breve;
- login admin;
- CRUD admin de tecnologías;
- carga/edición de proyectos;
- enriched profile por proyecto;
- media optimizada con variantes low/medium/high;
- readiness de búsqueda;
- re-embed individual y batch;
- traducciones persistidas para campos públicos seleccionados.

## 17. Deuda y siguientes pasos sugeridos

Pendientes razonables:

1. seguir reduciendo naming heredado `product -> project` y `brand -> client/context`;
2. formalizar template/source markdown definitivo para nuevas cargas;
3. evaluar importación semi-automática desde markdown hacia admin/API;
4. incorporar más proyectos reales usando el flujo de análisis externo;
5. sumar verificación frontend/E2E para flujos admin.

## 17.1 Regla documental consolidada para nuevas altas

Se formalizó una decisión adicional para mantener consistencia entre workflows y documentación:

- cuando un proyecto analizado ya tiene `90. dev_portfolioforge/<Project_Name>.md`, ese archivo pasa a ser la fuente editorial canónica para crear el proyecto en PortfolioForge;
- la UI actual debe entenderse como capa de persistencia/ejecución del contenido de ese markdown, no como fuente primaria de autoría;
- la carga manual debe comenzar leyendo ese archivo si existe y no debe reinventar contenido ya capturado allí;
- reanalizar el repositorio completo queda como acción secundaria cuando el markdown fuente no exista o esté desactualizado frente a evidencia nueva.

## 18. Resumen ejecutivo

Desde `/sdd-init`, PortfolioForge pasó de una base heredada y ambigua a una plataforma funcional con búsqueda híbrida, enrichment, tecnologías y documentación mucho más cercana al contrato real del producto.

La etapa reciente añadió algo decisivo: una visión clara de cómo transformar evidencia técnica externa —repositorios, carpetas y documentación real— en una fuente markdown estructurada y, desde ahí, en entradas de portfolio ricas, publicables y buscables.

En la actualización documental más reciente, además, se consolidó que esa fuente markdown ya no es solo una ayuda intermedia: si existe dentro de `90. dev_portfolioforge/`, debe operar como fuente de verdad editorial para futuras altas del proyecto en PortfolioForge.

## 19. Realineación reciente del modelo documental de proyectos

En la revisión documental más reciente se hicieron explícitas varias decisiones arquitectónicas/editoriales que antes estaban implícitas o incompletas.

### 19.1 Mapping canónico y mismatch actual

Se dejó documentado que:

- `Client / Context` sigue entrando por `brand` en admin/storage legacy;
- el significado público/canónico correcto es `client/context` / `client_name`;
- `Published` en realidad mapea a `active`;
- `Technologies` se persiste realmente vía `technology_ids`;
- `Main images` no debe pensarse ya como lista plana, sino como `media` con variantes y metadata, aunque siga existiendo una lista legacy derivada.

También se aclaró expresamente que el flujo manual y el automático están alineados lógicamente, pero todavía no son 1:1 exactos.

### 19.2 Fortalecimiento del contrato editorial fuente

Se decidió que el estándar fuente debe anticipar consumo futuro por retrieval/chatbot aun sin cambiar DB/UI hoy.

Nueva dirección documentada:

- `Technical Decisions` → `decision: ... | why: ... | tradeoff: ...`
- `Integrations` → `name: ... | type: ... | direction: ... | note: ...`
- `Results` → `result: ... | impact: ... | evidence: ...`
- `Timeline` → `phase: ... | objective: ... | outcome: ...`
- `Challenges` → `challenge: ... | mitigation: ... | status: ...`
- `Metrics` se mantiene como `key: value`, pero con keys normalizadas y unidad/período explícitos cuando sea posible.

La conclusión fue que el storage actual sigue siendo texto/listas simples, pero la fuente editorial ya no debe escribirse como listas débiles si se quiere assistant-readiness futura.

### 19.3 Media contract ampliado

Se reforzó que la convención de imágenes debe incluir:

- slug corto de proyecto;
- naming público por defecto `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_<low|medium|high>.webp`;
- mínimo 5 imágenes cuando el material lo permita;
- `Main images` normalmente apuntando a `_medium`;
- modelo canónico de media item con `low`, `medium`, `high`, `caption`, `alt_text`, `featured`, `sort_order`.

Además, se dejó explícito que ese host/path base es el **default público de assets** y no una restricción rígida del modelo editorial; si cambia en el futuro, debe ajustarse la convención documental manteniendo el contrato de variantes y extensión `.webp`.

### 19.4 PM / delivery aún sin schema dedicado

Se documentó explícitamente que la UI actual no tiene campos propios para:

- role/scope;
- stakeholders;
- constraints;
- delivery strategy;
- risks/mitigations;
- evidence sources.

Decisión vigente:

- mapear esa información dentro de `Business Goal`, `Problem Statement`, `Technical Decisions`, `Challenges`, `Timeline`, `Results` y `Metrics`;
- exigir que el workflow automático igualmente analice esas secciones aunque aún no se persistan 1:1.

### 19.5 Readiness tiers y realidad de búsqueda actual

Se formalizó una distinción clave:

- **search readiness**;
- **case study readiness**;
- **assistant readiness**.

Hallazgo importante:

- el producto actual implementa principalmente search readiness;
- todavía no garantiza por sí solo case study quality ni assistant-readiness completa.

También se dejó explícito que la implementación de búsqueda actual es más fuerte sobre:

- `title` / `name`;
- `description`;
- `client/context` (`brand`/`client_name`);
- `technologies`;
- `solution_summary`;
- `architecture`;
- `business_goal`;
- `problem_statement`;
- `ai_usage`.

Mientras tanto, secciones como `integrations`, `technical decisions`, `results`, `timeline` y `challenges` siguen siendo críticas editorialmente, pero hoy no tienen el mismo nivel de consumo directo en search/explanation.

### 19.6 Target schema futuro recomendado

Finalmente, se propuso como dirección futura separar mejor el contenido en:

- base project;
- media;
- technical profile;
- delivery profile;
- assistant/retrieval metadata.

Esta recomendación quedó documentada como evolución futura grounded en el repo actual, no como implementación ya existente.

## 20. Lección operativa reciente: no confiar en éxito parcial de importación

En la corrección reciente del proyecto `can-bus-crane-monitoring` quedó confirmado un failure mode concreto del flujo actual:

- el markdown fuente podía estar bien redactado y aun así el resultado persistido quedar incompleto;
- se detectó una combinación de `active=true` pese a `Published=false`, `project_profiles` vacíos, cero tecnologías asociadas y media contaminada con assets ajenos;
- además, un fallback de media sustituyó imágenes del proyecto por `Foto_perfil_2026_Cuadrado...`, lo que demostró que el sistema estaba tratando un parseo/import parcial como si fuera éxito.

Decisión documental consolidada a partir de ese caso:

- después de cada importación manual o automática hay que verificar payload admin/DB y payload público cuando aplique;
- la verificación debe ser campo por campo contra `90. dev_portfolioforge/<Project_Name>.md`;
- si el markdown dice `Published=false`, el proyecto debe quedar inactivo y fuera de la API pública;
- fallbacks, omisiones o contaminación de media deben clasificarse como fallo de ingestión.

Actualización posterior de estándar documental:

- para nuevos markdowns fuente y templates recomendados, `Published=true` pasa a ser el default editorial;
- `Published=false` sigue siendo autoritativo cuando se use de forma explícita y debe reflejarse como `active=false` sin excepciones.

## 21. Changes SDD archivados: `project-assistant-chat` + `authenticated-project-assistant`

El change `project-assistant-chat` introdujo la capacidad de assistant por proyecto y luego `authenticated-project-assistant` redefinió su contrato de acceso final.

Resultado consolidado vigente:

- el detalle del proyecto sigue siendo público;
- la API pública expone `assistant_available` como capacidad del proyecto derivada de `source_markdown_url` y no debe filtrar la URL privada;
- el backend expone el envío de mensajes del assistant solo en la ruta privada autenticada `POST /api/v1/private/projects/:slug/assistant/messages`;
- el assistant usa OpenAI solo desde backend y responde grounded en markdown fuente remoto del proyecto;
- usuarios públicos entran por Google o por OTP de email passwordless desde `/login`;
- admins locales conservan elegibilidad por el flujo oculto `/admin/login` con email/password;
- el assistant queda oculto para usuarios signed out y solo se habilita para sesiones elegibles.

## 22. Revisión visual local con Playwright

Se incorporó revisión visual/local basada en Playwright para cubrir flujos reales del frontend, incluyendo catálogo, detalle de proyecto y visibilidad del assistant solo para sesiones autenticadas y elegibles cuando corresponde.

Aprendizaje práctico:

- Playwright ya forma parte del stack operativo local para revisión visual y smoke coverage del flujo público.

## 23. Incidentes operativos reales posteriores al assistant

### 23.1 Migración faltante rompía proyectos públicos

Se confirmó un incidente real después de traer el change del assistant:

- los proyectos públicos devolvían `500` hasta aplicar la migración `sqlmigrations/20260414_1100_project_assistant_chat.sql` sobre la base activa;
- la causa fue dependencia del código nuevo sobre la columna `source_markdown_url` antes de que existiera en esa DB.

Lección consolidada:

- si el código nuevo toca payloads públicos y queries que derivan `assistant_available`, una DB sin esa migración puede romper endpoints públicos existentes.

### 23.2 `assistant_available` dependía de dato real, no solo del código

También se verificó un segundo problema real:

- aunque el feature estuviera desplegado, `assistant_available` seguía en `false` hasta configurar `source_markdown_url` en la DB/admin del proyecto correspondiente;
- no existe toggle independiente: la disponibilidad pública del assistant se deriva de ese campo.

Lección consolidada:

- parte del rollout del assistant es operacional/editorial, no solo de código; hay que poblar `source_markdown_url` proyecto por proyecto.

### 23.3 Host remoto de markdown inestable desde el backend

Se detectó además una condición real de runtime:

- el host remoto del markdown fue inestable o intermitente desde el entorno backend;
- para evitar degradación total del assistant, la implementación quedó con cache persistida de chunks y stale fallback local;
- el source of truth sigue siendo la URL remota, pero el assistant puede seguir respondiendo con contenido previamente resuelto cuando el fetch actual falla.

Decisión operativa/documental:

- la disponibilidad del assistant mejora con cache + stale fallback, pero la calidad y frescura siguen dependiendo de la reachability del markdown remoto.
