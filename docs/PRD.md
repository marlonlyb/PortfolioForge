# PRD - PortfolioForge

## 1. Resumen ejecutivo

PortfolioForge es una plataforma para convertir experiencia profesional real en case studies navegables, buscables y explicables. El producto combina un sitio público orientado a discovery con una consola admin donde cada proyecto se modela como una entrada rica de portfolio, no como una ficha comercial genérica.

La dirección actual del producto ya no se limita a “publicar proyectos manualmente”. El objetivo es establecer un flujo completo de ingestión que permita partir de evidencia real de trabajo —repositorios, carpetas de proyecto, documentación técnica, imágenes y notas operativas— para producir una fuente markdown estructurada y, desde allí, cargar proyectos consistentes dentro del admin.

## 2. Problema

Los portfolios tradicionales fallan en tres niveles:

- muestran resultados visuales, pero no el razonamiento técnico ni el impacto de negocio;
- no responden bien a búsquedas por evidencia real (`CAN Bus`, `SCADA`, `commissioning`, `React`, `Go`, `OpenAI`, etc.);
- dependen de carga manual inconsistente, lo que dificulta escalar un portfolio amplio y riguroso.

Además, gran parte de la experiencia profesional valiosa nace fuera de una “ficha” preparada para marketing: vive en repositorios, carpetas de entregables, diagramas, capturas, documentos técnicos y contexto operativo disperso. PortfolioForge debe capturar ese material y transformarlo en entradas publicables, comparables y buscables.

## 3. Visión del producto

PortfolioForge debe ser un portfolio técnico de alta credibilidad con dos capacidades principales:

1. **Descubrimiento público por evidencia real**  
   El visitante describe una necesidad, tecnología o contexto y el sistema recupera proyectos relevantes con explicaciones breves, limitadas a evidencia del proyecto.

2. **Producción estructurada de case studies**  
   El administrador puede transformar proyectos reales en entradas consistentes del portfolio usando un modelo de contenido claro, primero manual y progresivamente asistido por análisis de repositorios/carpetas.

## 4. Dirección actual

### 4.1 Lo que ya existe

El sistema ya implementa:

- sitio público con landing, búsqueda, catálogo y detalle de proyecto;
- búsqueda híbrida por evidencia real (filtros estructurados + FTS + fuzzy + semántica);
- explicaciones resumidas por resultado;
- project assistant por proyecto expuesto en el detalle público cuando existe markdown fuente configurado;
- admin de proyectos;
- admin de tecnologías;
- enriquecimiento de proyectos con `project_profiles`;
- recomposición del documento de búsqueda y re-embedding;
- revisión visual local basada en Playwright.

### 4.2 Lo que ahora se está consolidando

La dirección de producto actual incorpora un flujo editorial más riguroso:

- analizar un repositorio o carpeta de proyecto real;
- extraer contexto técnico, operativo y de negocio;
- mapear evidencia a los campos reales de PortfolioForge;
- generar una fuente markdown estable y revisable;
- tratar esa fuente markdown como origen editorial canónico del proyecto;
- usar esa fuente para completar el admin con menor fricción y mayor consistencia.

Esto convierte al markdown estructurado en la **fuente editorial canónica** entre la evidencia bruta y la carga en UI. La UI actual debe entenderse como la capa operativa de persistencia/ejecución de ese contenido, no como el lugar primario donde se inventa o redefine el proyecto.

## 5. Objetivos

### 5.1 Objetivos de producto

- demostrar experiencia real con evidencia técnica y resultados;
- permitir que un potencial cliente encuentre proyectos por problema, stack, arquitectura o contexto;
- ofrecer entradas de portfolio suficientemente ricas para funcionar como case studies;
- reducir la carga manual repetitiva en la creación de proyectos;
- preparar el sistema para ingestión semiautomática y, después, automática asistida.

### 5.2 Objetivos operativos del flujo de contenido

- normalizar cómo se nombran los proyectos;
- normalizar tecnologías, imágenes y campos enriquecidos;
- evitar portfolios “bonitos pero vacíos”;
- permitir trabajar desde material técnico real, no desde memoria improvisada.

## 6. No objetivos actuales

Fuera del alcance actual:

- chatbot libre conversacional sin grounding ni restricción por proyecto;
- CMS genérico multi-contenido;
- gestión completa de proyectos/PM dentro del producto;
- ingestión totalmente autónoma sin revisión humana;
- DAM/media management avanzada;
- blog o publishing editorial generalista;
- e-commerce, checkout, órdenes o pagos.

## 7. Usuarios principales

### 7.1 Visitante

Quiere:

- buscar experiencia por términos reales;
- entender rápido problema, solución, arquitectura y resultados;
- validar capacidad técnica sin hablar primero con el autor.

### 7.2 Administrador / autor del portfolio

Quiere:

- crear proyectos publicables con estructura consistente;
- cargar tecnologías reutilizables;
- enriquecer proyectos con contexto de negocio y ejecución;
- reutilizar documentación técnica real para poblar el portfolio.

## 8. Propuesta de valor

PortfolioForge no es solo un portfolio visual. Es un sistema para:

- **presentar experiencia real** como evidencia estructurada;
- **buscar** esa experiencia con modelos híbridos de recuperación;
- **explicar** por qué un proyecto coincide con una búsqueda;
- **operacionalizar** la creación de nuevos case studies desde material técnico existente.

## 9. Modelo de contenido vigente en admin

El modelo vigente debe alinearse con los campos reales del código y de la UI admin.

### 9.0 Mapping canónico, nombres vigentes y gap actual

La documentación debe distinguir entre:

- **nombre editorial / workflow**: el lenguaje con el que se analiza y redacta el proyecto;
- **nombre real actual de UI/storage**: el contrato que hoy existe en código y base.

| Dominio editorial / público | UI o storage actual | Observación |
|---|---|---|
| `Client / Context` | `brand` (legacy) | En consumo público el concepto correcto es `client/context` / `client_name`, pero la carga admin y la composición de búsqueda todavía dependen fuertemente de `brand`. |
| `Published` | `active` | `active = true` equivale a proyecto publicable; el naming de producto correcto sigue siendo Published/Unpublished. |
| `Technologies` | relación por `technology_ids` | La fuente editorial puede listar nombres, pero la persistencia real ocurre por IDs de tecnologías ya creadas. |
| `Main images` | `media` + `images` legacy derivado | La UI actual ya maneja variantes y metadata; la lista plana existe solo como compatibilidad y fallback. |

Regla importante: **el flujo manual y el flujo automático ya están alineados conceptualmente, pero todavía no son 1:1 exactos**. La documentación debe reflejar esa alineación lógica sin prometer equivalencia literal entre markdown fuente y storage actual.

### 9.1 Campos base de proyecto

- **Title**
- **Summary / Description**
- **Category**
- **Client / Context**
- **Markdown Source URL** (`source_markdown_url`) — solo admin/privado; habilita el assistant cuando existe
- **Main images**
- **Published**

### 9.2 Campos de enriquecimiento

- **Technologies**
- **Business Goal**
- **Problem Statement**
- **Solution Summary**
- **Architecture**
- **AI Usage**
- **Integrations**
- **Technical Decisions**
- **Challenges**
- **Results**
- **Metrics**
- **Timeline**

### 9.3 Reglas operativas clave

- Las tecnologías deben existir primero en `/admin/technologies`.
- Los campos actuales de tecnología son: `name`, `category`, `icon`, `brand color`.
- En el admin actual, `Client / Context` todavía viaja por el campo heredado `brand`, pero el dominio funcional correcto es `client/context`.
- `source_markdown_url` es privado/admin-only; la API pública no debe exponer esa URL.
- La API pública sí expone `assistant_available`, derivado de si `source_markdown_url` está presente y no vacío.
- `Main images` hoy se representa con media optimizada (`thumbnail_url`, `medium_url`, `full_url`) y una lista legacy derivada para compatibilidad.
- Para nuevas fuentes markdown generadas o normalizadas, la recomendación editorial por defecto es `Published=true`; `Published=false` debe usarse solo cuando se quiera mantener el proyecto fuera de la API pública.
- La implementación real de búsqueda/readiness actual está optimizada sobre todo para `title`, `description`, `client/context`, `technologies`, `solution_summary`, `architecture`, `business_goal`, `problem_statement` y `ai_usage`; otros bloques ricos siguen siendo valiosos editorialmente, pero hoy no tienen el mismo peso operativo en código.

## 10. Convenciones editoriales actuales

### 10.1 Naming de proyectos

El nombre fuente de un proyecto debe seguir esta convención:

- nombre corto y técnico;
- sin incluir cliente o empresa en el título;
- basado en la tecnología, dominio o concepto central.

Ejemplos válidos:

- `CAN Bus Crane Telemetry`
- `SCADA Commissioning Dashboard`
- `React Search Explainability Platform`

Ejemplos no recomendados:

- `Proyecto PRECOR`
- `Sistema de Empresa X`
- `Cliente Y - Dashboard`

El cliente o contexto pertenece al campo **Client / Context**, no al título.

### 10.2 Convención de imágenes

Para una fuente markdown de proyecto:

- usar como convención pública/canónica por defecto URLs completas con el patrón `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_<low|medium|high>.webp`;
- mantener al menos **5 imágenes** por proyecto cuando el material lo permita;
- en **Main images** usar normalmente las variantes `_medium` como referencia principal;
- reservar `_low` para catálogo / miniatura y `_high` para ampliación.

Importante: esta convención es el **default público de assets** para PortfolioForge. Se documenta explícitamente así para que el origen/base pueda cambiarse más adelante sin rediseñar el modelo editorial de media.

Pero el contrato editorial canónico ya no debe pensarse como “lista plana de imágenes”. Cada imagen debe entenderse como un **media item** con, idealmente, estos atributos:

- `low`
- `medium`
- `high`
- `caption`
- `alt_text`
- `featured`
- `sort_order`

Esta convención está alineada con la UI admin actual, que captura tres variantes por imagen: catálogo, galería y ampliada.

### 10.3 Serialización de campos enriquecidos

Las reglas actuales del sistema son de almacenamiento, no de calidad editorial. La UI sigue serializando/deserializando estos campos desde texto plano, pero el estándar fuente debe anticipar consumo futuro por búsqueda, case study y asistentes.

Convención fuente recomendada:

- **Technical Decisions** → `decision: ... | why: ... | tradeoff: ...`
- **Integrations** → `name: ... | type: ... | direction: ... | note: ...`
- **Results** → `result: ... | impact: ... | evidence: ...`
- **Timeline** → `phase: ... | objective: ... | outcome: ...`
- **Challenges** → `challenge: ... | mitigation: ... | status: ...`
- **Metrics** → `key: value`, con keys normalizadas y, cuando aplique, unidad/período explícitos

Compatibilidad con la UI actual:

- **una línea por ítem** para:
  - `integrations`
  - `technical_decisions`
  - `challenges`
  - `results`
  - `timeline`
- **`key: value` por línea** para `metrics`

Estas reglas existen porque la UI serializa/deserializa esos campos desde texto plano hacia arrays u objetos JSON. La mejora buscada no cambia DB/UI hoy; cambia cómo se redacta la fuente para que luego sea más recuperable y explicable.

### 10.4 Gaps actuales para narrativa de delivery / PM

La UI actual **no tiene campos dedicados** para:

- role/scope
- stakeholders
- constraints
- delivery strategy
- risks/mitigations
- evidence sources

Por lo tanto, hoy deben mapearse así:

- **role/scope** → `Business Goal` + `Solution Summary`
- **stakeholders** → `Business Goal` o `Problem Statement`
- **constraints** → `Problem Statement`
- **delivery strategy** → `Technical Decisions` + `Timeline`
- **risks/mitigations** → `Challenges`
- **evidence sources** → `Results`, `Metrics` y referencias explícitas dentro de `Summary / Description` cuando haga falta

Este mapping es transitorio, pero obligatorio mientras no exista un perfil de delivery dedicado.

## 11. Flujo manual actual

El flujo operativo disponible hoy es:

1. si existe `90. dev_portfolioforge/<Project_Name>.md`, leerlo primero y usarlo como fuente de verdad editorial;
2. crear tecnologías faltantes en `/admin/technologies`;
3. crear el proyecto base en `/admin/projects/new`;
4. cargar media optimizada y definir la imagen principal;
5. completar el rich profile usando el markdown fuente como referencia principal;
6. si el proyecto debe tener assistant, guardar `source_markdown_url` con una URL HTTPS pública y alcanzable del markdown fuente;
7. guardar;
8. verificar el resultado real en payload admin/público o DB comparando campo por campo contra el markdown fuente;
9. revisar readiness;
10. ejecutar re-embed cuando corresponda.

Regla operativa: si el markdown fuente ya existe, la carga manual no debe reinventar contenido ya definido allí. Reanalizar el repositorio completo debe ser una acción secundaria, reservada para cuando falte ese archivo o haya evidencia nueva que lo vuelva incompleto/desactualizado.

Contrato mínimo de verificación post-import:

- confirmar `title`;
- confirmar que `Published` del markdown coincide con `active` en storage/API;
- confirmar tecnologías por cantidad y nombre;
- confirmar que `business_goal`, `problem_statement`, `solution_summary`, `architecture`, `ai_usage` y demás bloques ricos no quedaron vacíos si el markdown sí trae contenido;
- confirmar que `integrations`, `technical_decisions`, `challenges`, `results`, `timeline` y `metrics` preservaron su contenido;
- confirmar que la galería y `images` legacy usan assets del proyecto correcto y no quedaron contaminadas con placeholders o assets de otro proyecto;
- si el proyecto debe tener assistant, confirmar `assistant_available=true` en el payload público y confirmar que `source_markdown_url` solo se vea en admin;
- si el markdown fuente dice `Published=false`, el resultado correcto es proyecto inactivo y ausencia del slug en la API pública.

## 12. Flujo objetivo de ingestión desde repositorio/carpeta

El flujo hacia el que apunta el producto es:

1. recibir una ruta de repositorio o carpeta;
2. analizar estructura, README, código, integraciones, assets y evidencia documental;
3. detectar tecnologías y contexto técnico;
4. producir un markdown fuente alineado con los campos reales del admin;
5. revisar y ajustar editorialmente el markdown hasta consolidarlo como fuente canónica del proyecto;
6. cargar el contenido en la UI actual a partir de ese markdown;
7. verificar el payload resultante contra el markdown fuente antes de considerar exitosa la ingestión;
8. publicar y reindexar cuando corresponda.

Este flujo debe permitir tomar proyectos del mundo real —por ejemplo, un caso industrial como el de grúa con CAN Bus analizado recientemente para PRECOR— sin acoplar el producto a un único cliente ni a un único vertical.

Además, el análisis automático debe producir y revisar aunque sea fuera de la UI actual estas secciones fuente obligatorias:

- role/scope
- stakeholders
- constraints
- delivery strategy
- risks/mitigations
- evidence sources

Aunque hoy no se persistan como columnas dedicadas, forman parte del análisis mínimo para lograr case studies y futura preparación para asistentes.

## 12.1 Assistant de proyecto implementado hoy

Estado actual implementado:

- el detalle público puede exponer un botón de assistant por proyecto;
- la API pública responde por `POST /api/v1/public/projects/:slug/assistant/messages`;
- el backend resuelve el proyecto por slug, verifica `source_markdown_url`, descarga el markdown remoto y selecciona secciones relevantes antes de consultar OpenAI;
- la URL markdown se mantiene privada/admin-only y la señal pública es `assistant_available`.

Comportamiento operativo actual:

- el assistant usa markdown remoto como source of truth;
- el backend mantiene cache de chunks de markdown con persistencia temporal local;
- si el fetch remoto falla o el host está inestable, el backend puede responder usando cache vigente o stale fallback persistido;
- esto mejora disponibilidad, pero la calidad de respuesta sigue dependiendo de la alcanzabilidad y frescura del markdown remoto.

Importante: esto **sí está implementado hoy**. Lo que todavía queda como evolución futura no es el assistant en sí, sino una capa de retrieval más sofisticada.

## 13. Requerimientos funcionales

### RF1. Gestión de proyectos

El sistema debe permitir crear, editar, publicar y despublicar proyectos con los campos base y enriquecidos vigentes.

### RF2. Gestión de tecnologías

El sistema debe permitir crear tecnologías reutilizables con `name`, `category`, `icon` y `brand color`, para luego asociarlas a proyectos.

### RF3. Búsqueda por evidencia real

El sistema debe recuperar proyectos por términos presentes o semánticamente relacionados con:

- título;
- summary / description;
- client / context;
- tecnologías;
- solution summary;
- architecture;
- business goal;
- problem statement;
- ai usage;

Importante: en la implementación real actual, la cobertura más fuerte está en `title`, `description`, `client/context`, `technologies`, `solution_summary`, `architecture`, `business_goal`, `problem_statement` y `ai_usage`. Secciones como `integrations`, `technical decisions`, `results`, `timeline` y `challenges` son muy valiosas para narrativa y futuro retrieval, pero hoy no están igualmente ponderadas por la composición principal de búsqueda.

### RF4. Explicación acotada

El sistema debe explicar brevemente por qué un resultado coincide con la búsqueda sin inventar información fuera de la evidencia del proyecto.

### RF5. Flujo editorial reproducible

El sistema debe tener una documentación operativa suficiente para que una persona pueda:

- crear proyectos manualmente hoy;
- preparar fuentes markdown consistentes;
- mapear evidencia técnica a los campos reales del producto.

### RF6. Preparación para ingestión asistida

El producto debe quedar preparado para un flujo futuro donde un análisis de repositorio/carpeta genere la fuente markdown inicial de un proyecto.

### RF7. Assistant markdown-grounded por proyecto

El sistema debe permitir responder preguntas sobre un proyecto específico usando exclusivamente su markdown fuente remoto configurado en `source_markdown_url`, sin exponer esa URL en la API pública.

### RF8. Degradación controlada del assistant

El assistant debe seguir siendo utilizable ante fallas temporales del fetch remoto usando cache local y stale fallback cuando exista contenido previamente resuelto para ese proyecto.

## 14. Requerimientos no funcionales

- documentación rigurosa y alineada al contrato real del código;
- consistencia de naming entre docs y UI actual;
- separación clara entre evidencia, fuente editorial y publicación;
- capacidad de crecer sin rehacer el modelo de contenido;
- mínima dependencia de interpretación subjetiva durante la carga;
- validación post-import obligatoria sobre payload/API y no solo sobre éxito de escritura en DB;
- fallbacks parciales o contaminaciones de media tratados como fallo de ingestión, no como éxito parcial.

## 15. Restricción actual importante

La UI todavía no tiene campos dedicados para gestión/lifecycle/PM, por lo que esa información debe incrustarse en los campos existentes.

Regla editorial actual:

- objetivos, alcance y stakeholders → **Business Goal** / **Problem Statement**;
- decisiones de ejecución, tradeoffs y rollout → **Technical Decisions** / **Timeline**;
- coordinación, operación y resultados de entrega → **Results** / **Metrics** / **Timeline**.

Hasta que exista un modelo específico para PM, esta es la forma correcta de conservar ese contexto sin perder rigor.

## 16. Tiers de calidad del contenido

Para evitar falsas expectativas, PortfolioForge debe distinguir tres niveles conceptuales de calidad:

### 16.1 Search readiness

Nivel mínimo realmente soportado hoy en código.

Busca asegurar que el proyecto sea encontrable y explicable en búsqueda con campos suficientes para FTS, fuzzy y embeddings.

### 16.2 Case study readiness

Nivel editorial superior.

Exige que el proyecto tenga narrativa técnica y de delivery suficientemente clara: problema, solución, decisiones, tradeoffs, resultados, timeline y media coherente. Hoy este nivel depende sobre todo de disciplina editorial, no de validaciones completas del sistema.

### 16.3 Assistant readiness

Nivel mixto entre capacidad implementada y evolución futura.

Hoy ya existe un assistant por proyecto grounded en markdown remoto, pero la calidad del resultado sigue dependiendo de que la fuente esté bien redactada y sea estable. El nivel completo exige además redacción semi-estructurada, evidencia explícita, naming consistente, metadata de delivery y secciones que permitan extracción/reuso por asistentes sin re-interpretación excesiva.

Situación actual: el producto ya cubre **search readiness** y una primera implementación real de assistant readiness operativa. Aun así, **case study readiness** completa y una assistant readiness más profunda siguen dependiendo de calidad editorial y evolución del retrieval.

## 17. Esquema objetivo futuro recomendado

Sin presentarlo como implementación vigente, la evolución natural del repo apunta a separar el contenido en capas más explícitas:

### 17.1 Base project

- identidad editorial (`title`, `slug`, `category`)
- resumen público (`description`, `client/context`, estado de publicación)
- taxonomía (`technology_ids`)

### 17.2 Media

- colección de media items
- variantes `low`, `medium`, `high`
- metadata editorial (`caption`, `alt_text`)
- orden y `featured`

### 17.3 Technical profile

- `business_goal`
- `problem_statement`
- `solution_summary`
- `architecture`
- `ai_usage`
- `integrations`
- `technical_decisions`

### 17.4 Delivery profile

- `role_scope`
- `stakeholders`
- `constraints`
- `delivery_strategy`
- `risks_mitigations`
- `timeline`
- `results`
- `metrics`
- `evidence_sources`

### 17.5 Assistant / retrieval metadata

- readiness tiers
- evidence confidence / verification flags
- normalized keywords / synonyms
- source provenance
- extraction quality notes
- persistencia dedicada de chunks / embeddings de markdown para retrieval más robusto y menos dependiente del fetch remoto en tiempo real

Este target schema es una recomendación de evolución para desacoplar mejor búsqueda, narrativa de case study y consumo por asistentes. No describe la implementación actual, que hoy usa fetch remoto del markdown con cache local y stale fallback.

## 18. Estado actual y próximos pasos

### Estado actual

- el producto ya soporta búsqueda híbrida y proyectos enriquecidos;
- el admin ya soporta tecnologías, media optimizada, enrichment y traducciones persistidas;
- el producto ya soporta assistant por proyecto basado en markdown fuente remoto configurado desde admin;
- la documentación se está ajustando para reflejar el flujo real de trabajo.

### Próximos pasos recomendados

- formalizar un template markdown estable para nuevas entradas;
- definir un pipeline repetible de análisis de repositorio/carpeta → markdown fuente;
- evaluar posterior importador automático desde markdown hacia el admin/API;
- evaluar persistencia dedicada de chunks/embeddings del markdown fuente para reducir dependencia del fetch remoto en runtime;
- separar con más claridad perfiles técnicos, de delivery y de retrieval en futuras iteraciones de schema;
- seguir reduciendo naming heredado (`product`, `brand`) en código interno.

## 19. Documentos complementarios

- `docs/MANUAL-PROJECT-INGESTION-WORKFLOW.md` — cómo cargar manualmente un proyecto hoy
- `docs/AUTOMATIC-PROJECT-INGESTION-WORKFLOW.md` — cómo analizar una carpeta/repositorio y producir el markdown fuente
- `docs/MEMORY-SDD.md` — memoria histórica del proceso SDD y del trabajo de agentes
