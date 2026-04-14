# Guía de creación manual de proyectos en la UI actual

## Objetivo

Esta guía explica cómo cargar un proyecto manualmente **hoy** en la UI admin de PortfolioForge, usando los campos reales del formulario y respetando sus contratos actuales.

Importante: esta guía describe la **verdad operativa actual de la UI**. Está alineada con el workflow automático a nivel lógico, pero todavía no existe un contrato 1:1 exacto entre markdown fuente, UI y storage.

Regla vigente para nuevas altas: si el proyecto ya tiene `90. dev_portfolioforge/<Project_Name>.md`, ese archivo debe tomarse como fuente editorial canónica. La UI se usa para persistir/ejecutar ese contenido en el sistema actual, no para inventar desde cero una versión alternativa del proyecto.

## 0. Mapping canónico vs nombres reales actuales

| Nombre editorial / workflow | UI o storage actual | Nota |
|---|---|---|
| Client / Context | `brand` | Legacy en admin/storage; en consumo público el concepto correcto es `client_name`. |
| Published | `active` | Flag real de publicación. |
| Technologies | `technology_ids` | La UI muestra nombres, pero persiste IDs. |
| Main images | `media` + `images` legacy derivado | El contrato real ya es más rico que una lista plana. |

Regla: documentar y redactar con el vocabulario editorial correcto, pero cargar respetando los nombres/campos reales actuales.

## 1. Flujo correcto hoy

### Paso 0 — Verificar si ya existe markdown fuente

Antes de tocar la UI, revisar si el repositorio/carpeta del proyecto ya contiene:

- `90. dev_portfolioforge/<Project_Name>.md`

Si existe:

- leer primero ese archivo completo;
- usarlo como referencia principal para título, summary, category, client/context, technologies, media y rich profile;
- no inventar ni reescribir en UI contenido que ya esté correctamente capturado allí;
- considerar el re-análisis del repo como secundario, solo si el markdown está ausente o claramente desactualizado frente a evidencia nueva.

Si no existe, recién ahí corresponde apoyarse en el workflow de análisis/manual según la evidencia disponible.

### Paso 1 — Crear primero las tecnologías

Antes de crear el proyecto, ve a:

- `/admin/technologies`

Cada tecnología debe existir previamente para poder seleccionarla luego en el proyecto.

### Campos actuales de tecnología

- **name**
- **category**
- **icon**
- **brand color**

Si una tecnología no existe todavía, créala primero. Esto es obligatorio para mantener filtros, embeddings y búsqueda por stack consistentes.

---

### Paso 2 — Crear el proyecto

Ve a:

- `/admin/projects/new`

La pantalla actual tiene tres bloques relevantes:

1. **Project profile**
2. **Rich Profile (Search Enrichment)**
3. **Traducciones persistidas** (solo disponibles después de guardar)

Para la carga inicial, enfócate en los dos primeros bloques.

Regla editorial durante toda la carga:

- la UI no debe convertirse en una segunda fuente de verdad separada del markdown fuente;
- cuando haya diferencias, primero corregir/confirmar la fuente markdown y luego persistir en UI;
- solo agregar en UI información que no estuviera en el markdown cuando exista evidencia nueva y verificable.

## 2. Campos reales del formulario de proyecto

## A. Project profile

### 1. Title

Campo real: `name`

**Qué poner**
- un nombre corto y técnico;
- sin nombre del cliente o empresa;
- basado en la tecnología o concepto central.

**Regla recomendada**
- bien: `CAN Bus Crane Telemetry`
- bien: `SCADA Commissioning Dashboard`
- evitar: `PRECOR Crane Project`

El cliente o contexto va en **Client / Context**, no en el título.

---

### 2. Summary / Description

Campo real: `description`

**Qué poner**
- resumen claro del proyecto;
- contexto técnico y funcional;
- qué se construyó y para qué.

**Buenas prácticas**
- 2 a 6 párrafos breves o un bloque denso pero legible;
- mencionar dominio, actor principal, sistema y valor generado;
- evitar descripciones vacías tipo “plataforma innovadora”.

Este campo impacta directamente la búsqueda y el readiness.

---

### 3. Category

Campo real: `category`

**Qué poner**
- una categoría consistente y reutilizable.

Ejemplos:

- `industrial automation`
- `embedded systems`
- `web platform`
- `ai tooling`
- `data platform`

---

### 4. Client / Context

Campo real actual: `brand` (heredado)  
Dominio funcional correcto: `client/context`

**Qué poner**
- nombre del cliente, marca, unidad o contexto operativo;
- si no conviene publicar cliente, usar contexto neutral pero preciso.

Ejemplos:

- `PRECOR`
- `Industrial automation lab`
- `Internal R&D`

---

### 5. Main images

En la UI actual no hay un textarea “Main images”; hoy se cargan como **Optimized media**.

Por cada imagen, la UI pide:

- `Low / catálogo URL` → `thumbnail_url`
- `Medium / galería URL` → `medium_url`
- `High / ampliada URL` → `full_url`
- `Legacy / fallback URL (optional)` → `url`
- `Caption`
- `Alt text`
- `Sort order`
- `Featured image`

### Convención recomendada de imágenes

- usar como convención pública/canónica por defecto URLs completas con patrón `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_<low|medium|high>.webp`;
- cargar **mínimo 5 imágenes** cuando el proyecto tenga suficiente material;
- marcar una imagen como **Featured image**;
- usar normalmente las variantes `_medium` como base visual principal de `Main images`.

Nota: esta base `https://mlbautomation.com/dev/portfolioforge/` debe tratarse como el **default público de assets**. Si más adelante cambia el host o path base, debe actualizarse esta convención sin alterar el contrato editorial `low/medium/high` ni la preferencia por `_medium` en `Main images`.

Modelo editorial canónico por imagen:

- `low`
- `medium`
- `high`
- `caption`
- `alt_text`
- `featured`
- `sort_order`

Esto está alineado con el comportamiento actual del frontend, que deriva la lista legacy de imágenes priorizando `medium_url`, luego `full_url`, luego `thumbnail_url`.

---

### 6. Markdown source URL (private)

Campo real: `source_markdown_url`

**Qué hace**
- se guarda solo en admin/privado;
- no debe exponerse en el payload público del proyecto;
- cuando existe y no está vacío, el frontend público puede mostrar el entrypoint del project assistant mediante `assistant_available=true`.

**Cuándo completarlo**
- complétalo cuando el proyecto tenga markdown fuente público y ese proyecto deba exponer assistant en el detalle público;
- déjalo vacío si no quieres assistant público o si el markdown aún no está listo.

**Validación obligatoria**
- debe ser URL `https://` válida;
- debe ser pública y alcanzable desde el backend;
- debe apuntar al markdown fuente real del proyecto, no a una landing o HTML genérico;
- debe contener el contenido esperado del case study para grounding del assistant.

**Regla operativa**
- la visibilidad del assistant se deriva de este campo; no existe un toggle público separado.

---

### 7. Published

Campo real: `active`

**Qué hace**
- si está activo, el proyecto puede entrar al sitio público;
- si está desactivado, permanece interno.

**Recomendación**
- para nuevos markdowns fuente, usar `Published=true` por defecto y persistir `active=true` salvo que exista una decisión explícita de mantener el proyecto interno.
- publicar solo cuando ya tenga al menos descripción, categoría, tecnologías y solution summary.
- si el markdown fuente dice `Published=false`, dejar `active=false` sin excepciones; no usar readiness ni "se ve bien" como justificación para publicarlo.

## B. Rich Profile (Search Enrichment)

Estos campos viven en `project_profiles` + relación con tecnologías.

### 8. Technologies

Campo real: `technology_ids`

**Regla obligatoria**
- las tecnologías deben existir primero en `/admin/technologies`.

**Qué poner**
- stack real, protocolos, frameworks, plataformas o dominios técnicos relevantes.

Ejemplos:

- `CAN Bus`
- `React`
- `Go`
- `PostgreSQL`
- `OpenAI`
- `Siemens`

---

### 9. Business Goal

Campo real: `business_goal`

**Qué poner**
- el objetivo de negocio, operativo o de entrega;
- por qué el proyecto existía.

Ejemplo:

> Reducir el tiempo de diagnóstico y mejorar la visibilidad operativa de un sistema de telemetría de grúa basado en CAN Bus.

---

### 10. Problem Statement

Campo real: `problem_statement`

**Qué poner**
- el problema concreto antes de la solución;
- restricciones, dispersión de información, fallas de trazabilidad o dolor operativo.

---

### 11. Solution Summary

Campo real: `solution_summary`

**Qué poner**
- resumen de la solución implementada;
- arquitectura funcional principal;
- cómo se resolvió el problema.

Este es uno de los campos más importantes para búsqueda y readiness.

---

### 12. Architecture

Campo real: `architecture`

**Qué poner**
- componentes principales;
- flujos de datos;
- capas del sistema;
- servicios, protocolos o integraciones críticas.

---

### 13. AI Usage

Campo real: `ai_usage`

**Qué poner**
- cómo se usó IA si aplica;
- dejar vacío si el proyecto no utilizó IA.

No fuerces este campo si el proyecto no tuvo uso real de IA.

---

### 14. Integrations

Campo real: `integrations`

**Regla de serialización actual**
- **una línea por integración**

**Convención editorial recomendada**
- `name: ... | type: ... | direction: ... | note: ...`

Ejemplo:

```text
name: CAN Bus gateway | type: protocol gateway | direction: inbound | note: captura de señales de campo
name: MQTT broker | type: messaging | direction: bidirectional | note: distribución de eventos operativos
name: REST maintenance API | type: external API | direction: outbound | note: consulta de tickets/diagnóstico
```

---

### 15. Technical Decisions

Campo real: `technical_decisions`

**Regla de serialización actual**
- **una línea por decisión**

**Convención editorial recomendada**
- `decision: ... | why: ... | tradeoff: ...`

Ejemplo:

```text
decision: usar parser dedicado para tramas CAN de alta frecuencia | why: reducir pérdida de eventos | tradeoff: mayor complejidad de mantenimiento
decision: separar captura de telemetría y visualización operativa | why: aislar carga de adquisición | tradeoff: más coordinación entre componentes
decision: persistir eventos críticos para auditoría técnica | why: soportar troubleshooting histórico | tradeoff: mayor costo de almacenamiento
```

---

### 16. Challenges

Campo real: `challenges`

**Regla de serialización actual**
- **una línea por desafío**

**Convención editorial recomendada**
- `challenge: ... | mitigation: ... | status: ...`

---

### 17. Results

Campo real: `results`

**Regla de serialización actual**
- **una línea por resultado**

**Convención editorial recomendada**
- `result: ... | impact: ... | evidence: ...`

---

### 18. Metrics

Campo real: `metrics`

**Regla de serialización actual**
- formato **`key: value`**;
- una línea por métrica.

Ejemplo:

```text
diagnosis_time_pct: -35% | period: commissioning baseline vs rollout
signals_observed_count: 120+ | unit: signals
commissioning_errors_pct: -18% | period: first release window
```

Si una línea no usa `key: value`, el formulario lanza error al parsear.

Recomendación editorial:

- usar keys normalizadas (`snake_case`);
- explicitar unidad y/o período cuando sea posible;
- evitar métricas vagas sin evidencia.

---

### 19. Timeline

Campo real: `timeline`

**Regla de serialización actual**
- **una línea por hito, etapa o fase**

**Convención editorial recomendada**
- `phase: ... | objective: ... | outcome: ...`

Ejemplo:

```text
phase: discovery | objective: relevar señales y contexto operativo | outcome: mapa inicial de componentes y actores
phase: signal mapping | objective: normalizar tramas CAN | outcome: catálogo usable para visualización
phase: implementation | objective: construir interfaz y backend de soporte | outcome: flujo funcional integrado
phase: field validation | objective: probar en entorno real | outcome: ajustes de ruido/frecuencia
phase: rollout | objective: habilitar uso operativo | outcome: adopción inicial y feedback
```

## 3. Dónde poner información de management/lifecycle hoy

La UI actual **no tiene campos dedicados para**:

- role/scope
- stakeholders
- constraints
- delivery strategy
- risks/mitigations
- evidence sources

Por eso, esa información debe incrustarse en los campos disponibles:

- rol y alcance real → **Business Goal** + **Solution Summary**
- stakeholders → **Business Goal** o **Problem Statement**
- restricciones, contexto y dolor operativo → **Problem Statement**
- estrategia de entrega y tradeoffs → **Technical Decisions**
- riesgos y mitigaciones → **Challenges**
- hitos, rollout y fases → **Timeline**
- outcomes de delivery → **Results** y **Metrics**
- fuentes de evidencia → mencionar explícitamente en **Results**, **Metrics** o **Summary / Description**

No inventes campos nuevos fuera del contrato actual.

## 4. Readiness: qué significa hoy y qué no

- **Search readiness**: nivel realmente soportado hoy por código y admin.
- **Case study readiness**: calidad editorial para que el proyecto funcione como caso técnico completo.
- **Assistant readiness**: calidad de estructuración para el assistant ya implementado hoy y para futura evolución de retrieval.

Hoy la UI valida y expone sobre todo el primer nivel. El assistant ya existe, pero su calidad sigue dependiendo de disciplina editorial durante la carga manual y de que el markdown remoto sea alcanzable.

## 5. Orden recomendado de carga manual

1. crear tecnologías si faltan;
2. completar **Title**;
3. completar **Summary / Description**;
4. definir **Category**;
5. completar **Client / Context**;
6. si el proyecto debe tener assistant, completar **Markdown source URL (private)** con una URL HTTPS pública y alcanzable;
7. cargar al menos 5 imágenes con variantes `_low`, `_medium`, `_high`;
8. seleccionar **Technologies**;
9. redactar **Business Goal** y **Problem Statement**;
10. redactar muy bien **Solution Summary**;
11. completar **Architecture**, **AI Usage** si aplica, listas, métricas y timeline;
12. guardar;
13. verificar payload admin/DB y, si quedó activo, payload público;
14. revisar readiness;
15. re-embed si el proyecto ya está listo para búsqueda.

## 6. Verificación obligatoria después de guardar

No alcanza con que la UI acepte el submit. Después de guardar, hay que comparar el resultado real contra `90. dev_portfolioforge/<Project_Name>.md`.

Validación mínima obligatoria:

- `GET /api/v1/admin/products/:id` o revisión equivalente en DB para confirmar campos base y media;
- `GET /api/v1/public/projects/:slug?lang=es` si `active=true`;
- si el proyecto debe tener assistant, verificar que admin/DB conserve `source_markdown_url`, que el payload público exponga `assistant_available=true` y que no filtre `source_markdown_url`;
- si el proyecto debe tener assistant, probar `POST /api/v1/public/projects/:slug/assistant/messages` con una pregunta simple;
- si el markdown dice `Published=false`, confirmar lo contrario: endpoint público `404` y ausencia del slug en `GET /api/v1/public/projects`.

Checklist campo por campo:

- [ ] `Title` coincide con la fuente o con la normalización editorial aprobada
- [ ] `active` respeta `Published`
- [ ] `Client / Context` no quedó truncado o reemplazado por otro valor incorrecto
- [ ] `source_markdown_url` quedó cargado correctamente en admin cuando aplica
- [ ] `Technologies` coincide por cantidad y nombre con el markdown
- [ ] `Business Goal`, `Problem Statement`, `Solution Summary`, `Architecture` y `AI Usage` no quedaron vacíos si la fuente tiene contenido
- [ ] `Integrations`, `Technical Decisions`, `Challenges`, `Results`, `Timeline` conservan todos los ítems esperados
- [ ] `Metrics` conserva sus claves y valores
- [ ] `images` legacy y `media` apuntan al proyecto correcto
- [ ] la imagen principal y las primeras imágenes no están contaminadas con assets ajenos o placeholders
- [ ] el payload público no expone `source_markdown_url`
- [ ] `assistant_available` coincide con la presencia real de `source_markdown_url`
- [ ] el assistant responde al menos una pregunta básica cuando está habilitado

Regla de fallo:

- si el importador/manual carga solo una parte;
- si faltan tecnologías o bloques ricos que sí están en el markdown;
- si aparece fallback de media o mezcla con assets de otro proyecto;

entonces la carga debe tratarse como **fallida**, aunque la escritura en DB haya devuelto éxito.

## 7. Checklist mínimo antes de publicar

- [ ] Title corto y técnico
- [ ] Summary / Description con contexto real
- [ ] Category consistente
- [ ] Client / Context correcto
- [ ] `active` alineado con `Published` del markdown fuente
- [ ] `source_markdown_url` configurado si el proyecto debe tener assistant
- [ ] mínimo 5 imágenes
- [ ] una imagen marcada como principal
- [ ] Technologies seleccionadas
- [ ] Solution Summary completo
- [ ] Problem Statement y Architecture razonables
- [ ] Integrations / Technical Decisions / Challenges / Results / Timeline con una línea por ítem
- [ ] Líneas semi-estructuradas en listas enriquecidas (`decision: ... | why: ...`, etc.)
- [ ] Metrics con formato `key: value`
- [ ] Información de delivery/PM distribuida en campos existentes sin perder contexto
- [ ] payload admin/DB verificado contra el markdown fuente
- [ ] payload público verificado si el proyecto queda activo
- [ ] `assistant_available` verificado si existe markdown fuente
- [ ] `source_markdown_url` sigue siendo admin-only
- [ ] assistant probado con una pregunta simple si está habilitado
- [ ] media principal y galería sin assets ajenos ni contaminación de otro proyecto

## 8. Ejemplo de criterio editorial

Un caso como el analizado recientemente de grúa/CAN Bus para PRECOR no debería titularse con el nombre del cliente. La carga correcta sería algo más cercano a:

- **Title**: `CAN Bus Crane Telemetry`
- **Client / Context**: `PRECOR`

Ese criterio hace que el portfolio sea más reusable, más técnico y menos dependiente de naming comercial.
