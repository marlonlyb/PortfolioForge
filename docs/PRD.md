# PRD — PortfolioForge

## 1. Propósito del producto

PortfolioForge es una plataforma para convertir experiencia profesional real en proyectos públicos navegables, buscables y explicables.

Su objetivo no es solo “mostrar trabajos”, sino operar un sistema editorial donde:

- un **markdown canónico** concentra la narrativa principal del proyecto;
- la **UI + DB** persisten una versión estructurada, resumida y operativa de ese contenido;
- la **búsqueda** recupera proyectos por evidencia real;
- el **assistant** responde preguntas grounded exclusivamente en el markdown publicado vía `source_markdown_url`.

El producto actual ya funciona como portfolio público, consola admin de curación y runtime de búsqueda/assistant. La dirección vigente consolida al `.md` canónico como fuente editorial principal para altas y actualizaciones.

## 1.1 Documentos normativos vigentes

Para el flujo editorial y operativo de proyectos, los únicos documentos normativos vigentes dentro del repo son:

- `docs/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md` — cómo generar el markdown canónico;
- `docs/PROJECT-RUNTIME-INGESTION-GUIDE.md` — cómo poblar o actualizar la UI/DB runtime desde ese markdown;
- `docs/PRD.md` — marco de producto, arquitectura y reglas de relación entre markdown, UI/DB, búsqueda y assistant.

Regla explícita:

- cualquier documento similar ubicado fuera de `docs/`, especialmente los workflows en `/home/marlon_ly/Workspace/`, debe tratarse como **legacy informativo** y **no** como fuente de verdad operativa.

---

## 2. Principios del producto real

1. **El repositorio fuente no se publica directamente**: primero se transforma en un `.md` canónico.
2. **El `.md` canónico es la fuente editorial principal**: resume, organiza y estabiliza el proyecto.
3. **La UI admin no debe reinventar el proyecto**: debe persistir y ajustar lo ya definido en el markdown.
4. **La DB y la UI son la capa runtime estructurada**: exponen catálogo, detalle, búsqueda, readiness, localización y media optimizada.
5. **El assistant no usa chat libre**: responde sobre un proyecto concreto usando solo el markdown remoto configurado.
6. **El dominio funcional es `project`** aunque persistan nombres legacy como `product` y `brand`.

---

## 3. Usuarios principales

### 3.1 Visitante público

Quiere:

- descubrir experiencia por búsqueda real (`React`, `SCADA`, `commissioning`, etc.);
- explorar catálogo y detalle de proyectos;
- entender problema, ejecución y solución técnica;
- consultar el assistant cuando el proyecto lo permite y su sesión es elegible.

### 3.2 Administrador / autor del portfolio

Quiere:

- crear y actualizar proyectos consistentes;
- reutilizar un markdown canónico como fuente de verdad editorial;
- gestionar tecnologías, media, traducciones y readiness;
- publicar proyectos con assistant grounded y búsqueda útil.

---

## 4. Arquitectura real del sistema

## 4.1 Backend

Stack actual:

- Go 1.20
- Echo v4
- PostgreSQL 16
- pgvector + FTS + pg_trgm
- OpenAI para embeddings, explanations y assistant

Estructura real:

- `cmd/` — wiring y rutas
- `domain/services/` — reglas de búsqueda, proyectos, assistant, auth
- `domain/ports/` — contratos
- `infrastructure/postgres/` — repositorios y composición de documentos de búsqueda
- `infrastructure/handlers/` — APIs públicas, privadas y admin
- `infrastructure/localization/` — traducción persistida y aplicación por locale
- `model/` — contratos `Project`, `ProjectProfile`, `AdminProject`, `User`, etc.
- `sqlmigrations/` — evolución real del esquema

## 4.2 Frontend

Stack actual:

- React 19
- TypeScript 5.9
- Vite 7

Features reales:

- `landing/` — entrada pública
- `search/` — búsqueda híbrida y resultados
- `catalog/` — catálogo y detalle público
- `admin-products/` y `admin-projects/` — administración de proyectos durante transición legacy
- `admin-technologies/` — CRUD de tecnologías
- `auth/` — login, signup, OTP, perfil completo
- `shared/i18n/` — locale público `es`, `ca`, `en`, `de`

## 4.3 Capas de información del producto

PortfolioForge ya opera con cuatro capas distintas:

1. **Repositorio fuente / evidencia bruta**
2. **Markdown canónico del proyecto**
3. **UI + DB runtime resumida y estructurada**
4. **Assistant grounded en markdown remoto**

Relación entre capas:

- el repo fuente sirve para generar el `.md` canónico;
- el `.md` canónico sirve para poblar el proyecto en UI/DB;
- la UI/DB alimenta catálogo, detalle, búsqueda y localización;
- `source_markdown_url` publica ese mismo `.md` para el assistant.

---

## 5. Estado funcional real

El repositorio actual implementa o soporta:

- portfolio público con catálogo y detalle por slug;
- búsqueda híbrida por evidencia real;
- explicaciones breves por resultado;
- assistant/chat por proyecto grounded en markdown remoto;
- `source_markdown_url` privado en admin y `assistant_available` derivado en público;
- enriquecimiento estructurado mediante `project_profiles`;
- media optimizada por variantes `low_url` / `medium_url` / `high_url` con `fallback_url` como respaldo final;
- localización persistida por proyecto y campo;
- auto-traducción desde español base y override manual;
- readiness de búsqueda;
- recomposición del documento de búsqueda y re-embedding;
- transición progresiva de legacy `product` hacia dominio `project`.

---

## 6. Modelo de dominio actual

## 6.1 Dominio funcional vs storage heredado

La realidad actual es híbrida:

| Dominio funcional | Storage / compat actual | Nota |
|---|---|---|
| `project` | tabla `products` | transición activa |
| `client_name` | columna `brand` | público usa `client_name`, admin/storage aún aceptan `brand` |
| `published` | `active` | `active=true` equivale a proyecto publicado |
| `project_profiles` | `project_profiles` | enriquecimiento real del detalle |
| `media` | `project_media` + `images` legacy | `images` queda como compat/fallback |

Regla: la documentación canónica debe hablar en lenguaje **project**, pero sin ocultar que la persistencia todavía depende de nombres legacy.

## 6.2 Entidad pública `Project`

Campos runtime relevantes:

- `id`
- `name`
- `slug`
- `description`
- `category`
- `client_name`
- `status`
- `featured`
- `active`
- `assistant_available`
- `images`
- `media[]`
- `profile`
- `technologies[]`

## 6.3 Perfil enriquecido `ProjectProfile`

Campos reales actuales:

- `business_goal`
- `problem_statement`
- `solution_summary`
- `delivery_scope`
- `responsibility_scope`
- `architecture`
- `ai_usage`
- `integrations`
- `technical_decisions`
- `challenges`
- `results`
- `metrics`
- `timeline`

Los campos mínimos nuevos ya presentes en el producto son:

- `delivery_scope`
- `responsibility_scope`

## 6.4 Tecnologías

Las tecnologías existen como entidad separada y deben crearse antes de asociarse a un proyecto.

Campos reales:

- `name`
- `slug`
- `category`
- `icon`
- `color`

## 6.5 Media

Contrato vigente por ítem:

- `low_url`
- `medium_url`
- `high_url`
- `fallback_url`
- `caption`
- `alt_text`
- `featured`
- `sort_order`
- `media_type`

La UI usa media optimizada como contrato principal y reconstruye `images` solo para compatibilidad/fallback.

## 6.6 Localización persistida

Locales públicos soportados:

- `es` (base)
- `ca`
- `en`
- `de`

Los campos traducibles persistidos incluyen:

- `name`, `description`, `category`
- `business_goal`, `problem_statement`, `solution_summary`
- `delivery_scope`, `responsibility_scope`
- `architecture`, `ai_usage`
- `integrations`, `technical_decisions`, `challenges`, `results`, `metrics`, `timeline`

Cada traducción se guarda por `project_id + locale + field_key`, con modo `auto` o `manual`.

---

## 7. Rol del markdown canónico

El `.md` canónico es la **fuente editorial principal** del proyecto.

Debe cumplir estos roles:

- condensar la evidencia real del repositorio/carpeta;
- estabilizar naming, narrativa y estructura;
- servir como insumo para poblar los campos del proyecto en la UI;
- servir como documento que luego puede publicarse para el assistant;
- separar la autoría editorial de la representación runtime en DB.

### Regla operativa

La generación del markdown canónico debe seguir `docs/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`.

Si ya existe un markdown canónico confiable para un proyecto, ese archivo debe leerse primero y tratarse como fuente de verdad editorial. Re-analizar el repositorio completo solo se justifica cuando:

- el `.md` no existe;
- el `.md` está desactualizado;
- aparece evidencia nueva relevante.

---

## 8. Rol de `source_markdown_url`

`source_markdown_url` es la URL HTTPS pública del markdown canónico usado por el assistant.

Reglas reales del sistema:

- se configura solo en admin;
- no debe exponerse en la API pública;
- debe ser HTTPS válida;
- si existe y no está vacía, el proyecto público expone `assistant_available=true`;
- el assistant además exige sesión autenticada elegible;
- el backend descarga el markdown remoto, lo fragmenta y selecciona secciones relevantes antes de consultar el modelo.

En resumen:

- **markdown canónico** = fuente editorial principal;
- **`source_markdown_url`** = publicación remota de esa fuente para el assistant.

---

## 9. Relación entre markdown, UI/DB y assistant

### 9.1 Markdown → UI/DB

La UI/DB persiste una versión resumida y estructurada del proyecto:

- campos base del proyecto;
- `project_profiles`;
- tecnologías relacionadas;
- media optimizada;
- localizaciones;
- readiness y search document.

La compresión editorial y el mapping campo por campo deben seguir `docs/PROJECT-RUNTIME-INGESTION-GUIDE.md`.

### 9.2 UI/DB → Búsqueda

La búsqueda no indexa el markdown remoto directamente. Indexa la composición runtime del proyecto persistido.

Pesos reales del documento de búsqueda:

- **A**: `solution_summary`, `name`, `brand/client_name`
- **B**: `architecture`, `description`, tecnologías
- **C**: `business_goal`, `problem_statement`, `ai_usage`

El texto de embedding también usa principalmente:

- `name`
- `brand/client_name`
- `description`
- `solution_summary`
- `architecture`
- `business_goal`
- `problem_statement`
- `ai_usage`
- tecnologías

### 9.3 Markdown remoto → Assistant

El assistant:

- resuelve el proyecto por slug;
- verifica que el proyecto esté activo;
- verifica que exista `source_markdown_url`;
- descarga el markdown remoto;
- selecciona chunks relevantes por términos de la pregunta e historial;
- responde grounded en esas secciones.

Por eso el markdown canónico debe ser suficientemente completo y bien estructurado para soportar preguntas reales.

---

## 10. Modelo editorial actual del detalle

El detalle público ya está organizado por tres capas de información:

## 10.1 Estrategia

Campos actuales:

- `business_goal`
- `problem_statement`
- `solution_summary`

Pregunta que responde: **por qué existió el proyecto y qué solución entregó**.

## 10.2 Ejecución

Campos actuales:

- `delivery_scope`
- `responsibility_scope`
- `challenges`
- `results`
- `timeline`

Pregunta que responde: **cómo se ejecutó, qué se entregó, qué retos hubo y qué resultados produjo**.

## 10.3 Técnica

Campos actuales:

- `architecture`
- `ai_usage`
- `integrations`
- `technical_decisions`
- `metrics`

Pregunta que responde: **cómo estaba construido, con qué integraciones, qué decisiones técnicas importaron y con qué métricas puede leerse**.

Regla editorial: cualquier nuevo markdown canónico debe redactarse de manera que estas tres capas puedan mapearse sin ambigüedad.

---

## 11. Flujos importantes del producto

## 11.1 Flujo público

1. el visitante entra a landing, search o catálogo;
2. consulta proyectos publicados;
3. abre `/projects/:slug`;
4. consume overview, media, detalle por capas y contexto de búsqueda;
5. si `assistant_available=true`, el detalle muestra acceso condicional al assistant según sesión.

## 11.2 Flujo admin de alta/edición

1. crear o editar el proyecto base;
2. guardar descripción, categoría, client/context y `source_markdown_url` si aplica;
3. cargar media optimizada;
4. asociar tecnologías existentes;
5. completar `project_profiles`;
6. revisar traducciones;
7. verificar readiness;
8. re-embed cuando cambie contenido indexable.

## 11.3 Flujo de markdown canónico

1. analizar repo/carpeta fuente;
2. producir un único `.md` canónico del proyecto siguiendo `docs/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`;
3. usar ese `.md` para poblar el proyecto en PortfolioForge siguiendo `docs/PROJECT-RUNTIME-INGESTION-GUIDE.md`;
4. si se quiere assistant, publicar ese `.md` en una URL HTTPS;
5. guardar esa URL en `source_markdown_url`.

## 11.4 Flujo del assistant

1. el proyecto debe estar activo;
2. debe existir `source_markdown_url` válido;
3. el usuario debe estar autenticado;
4. el usuario debe tener `can_use_project_assistant=true`;
5. el backend responde usando solo el markdown remoto del proyecto.

Condiciones típicas de elegibilidad del usuario:

- cuenta autenticada;
- email verificado donde corresponda;
- perfil completo (`full_name`, `company`).

## 11.5 Flujo de media

1. el admin carga `media[]` con variantes optimizadas;
2. define `featured` y `sort_order`;
3. el frontend usa `low_url` para cards, `medium_url` para hero/galería y `high_url` para lightbox, con `fallback_url` como último respaldo;
4. `images` queda como derivado de compatibilidad.

## 11.6 Flujo de localización

1. el contenido base se mantiene en español;
2. al cambiar campos traducibles, el sistema sincroniza traducciones automáticas para `ca`, `en`, `de`;
3. el admin puede sobrescribir manualmente cualquier campo;
4. la API pública aplica localización cuando recibe `?lang=`.

## 11.7 Flujo de búsqueda y enriquecimiento

1. el admin actualiza proyecto, profile, tecnologías o media;
2. el backend recompone `project_search_documents`;
3. si hay embeddings activos, genera nuevo embedding;
4. la búsqueda híbrida fusiona FTS + fuzzy + semántica + filtros estructurados;
5. cada resultado muestra explicación acotada por evidencia.

---

## 12. Reglas operativas para alta y actualización de proyectos

1. **Si existe markdown canónico, se usa primero.**
2. **No se inventa contenido en la UI** si el `.md` ya define el proyecto.
3. **Las tecnologías se crean antes** de asociarlas al proyecto.
4. **`source_markdown_url` solo se guarda si hay URL HTTPS pública real.**
5. **`source_markdown_url` nunca debe filtrarse públicamente.**
6. **`Published` editorial mapea a `active` runtime.**
7. **La verificación posterior es obligatoria**: guardar en DB no basta.
8. **Cambios en campos indexables exigen refresh de búsqueda y normalmente re-embed.**
9. **Las traducciones manuales no deben ser sobrescritas por auto-sync.**
10. **Media y narrativa deben pertenecer al proyecto correcto**; mezclar assets ajenos invalida la carga.

### Verificación mínima post-alta/actualización

- `title/name`
- `description`
- `category`
- `client/context`
- `active/published`
- tecnologías asociadas
- profile completo según el markdown
- media correcta y ordenada
- `assistant_available` si corresponde
- ausencia pública de `source_markdown_url`
- localizaciones persistidas si se editaron
- readiness consistente

---

## 13. Convenciones útiles para contenido canónico

## 13.1 Convenciones de naming

- título corto, técnico y descriptivo;
- no meter el cliente en el título;
- `client/context` vive aparte.

## 13.2 Convenciones de serialización editorial

Para máxima compatibilidad con la UI actual:

- `integrations`, `technical_decisions`, `challenges`, `results`, `timeline` → una línea por ítem;
- `metrics` → `clave: valor` por línea.

Esto importa porque la UI serializa/deserializa estos bloques desde texto plano hacia arrays/objetos.

## 13.3 Convenciones de media

Contrato editorial recomendado por imagen:

- `low_url`
- `medium_url`
- `high_url`
- `fallback_url`
- `caption`
- `alt_text`
- `featured`
- `sort_order`

Uso recomendado:

- `low_url` para catálogo/miniatura;
- `medium_url` para detalle y galería principal;
- `high_url` para ampliación;
- `fallback_url` como respaldo final cuando falte la variante preferida.

Convención pública vigente de ejemplo:

- `low_url` → `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen0N_low.webp`
- `medium_url` → `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen0N_medium.webp`
- `high_url` → `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen0N_high.webp`
- `fallback_url` → asset base/original cuando aplique.

## 13.4 Convenciones de authoring del markdown

El markdown canónico debe:

- distinguir claramente Estrategia / Ejecución / Técnica;
- evitar claims sin evidencia;
- priorizar hechos, decisiones, resultados y tradeoffs;
- escribir pensando a la vez en lectura humana, mapeo a UI y uso posterior por assistant.

---

## 14. Restricciones y realidades actuales

- el sistema todavía persiste sobre tablas legacy `products`;
- el dominio público y editorial correcto es `project`;
- la búsqueda runtime consume sobre todo campos resumidos de DB, no todo el richness posible del markdown;
- el assistant sí depende del markdown remoto real;
- localización, media optimizada y capas de detalle ya son parte del producto vigente;
- la transición `product` → `project` sigue activa y debe documentarse, no ignorarse.

---

## 15. Definición canónica para trabajo futuro

Para crear proyectos equivalentes en PortfolioForge, la referencia correcta es:

1. generar o actualizar un **único markdown canónico** por proyecto, usando `docs/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`;
2. usar ese markdown para poblar el runtime estructurado, usando `docs/PROJECT-RUNTIME-INGESTION-GUIDE.md`;
3. mantener el detalle editorial organizado por **Estrategia / Ejecución / Técnica**;
4. publicar el markdown en una URL HTTPS si el proyecto usará assistant;
5. verificar siempre coherencia entre markdown, UI/DB, búsqueda y assistant.

Ese es el contrato real del producto hoy.

Regla final de gobierno documental:

- para estas funciones, no deben usarse como referencia operativa principal documentos legacy fuera de `docs/`, aunque se conserven por contexto histórico.
