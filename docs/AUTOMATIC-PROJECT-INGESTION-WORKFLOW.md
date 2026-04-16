# Flujo automático de análisis de repositorio/carpeta a markdown fuente

## Objetivo

Definir el workflow deseado para tomar una ruta de repositorio o carpeta real, analizar su evidencia, producir un markdown fuente alineado con los campos reales de PortfolioForge y usar ese mismo archivo como fuente de verdad para crear un nuevo proyecto en PortfolioForge.

Este documento no describe un importador ya implementado en código. Describe el proceso recomendado para generar una fuente consistente dentro de cada repositorio/carpeta y luego consumir esa fuente para cargar un proyecto nuevo con el menor nivel posible de reinterpretación manual.

Importante: el workflow automático y la carga manual están alineados a nivel lógico, pero todavía **no** son un contrato 1:1 exacto entre fuente markdown, UI y storage.

## 0. Mapping canónico actual

| Nombre editorial / workflow | UI o storage actual | Estado |
|---|---|---|
| Client / Context | `brand` | Legacy, pero sigue siendo el canal operativo principal en admin/storage. |
| Published | `active` | Flag real de publicación. |
| Technologies | `technology_ids` | La fuente puede listar nombres; la persistencia final usa IDs existentes. |
| Main images | `media` + `images` legacy derivado | El markdown puede listar imágenes para compatibilidad, pero el contrato editorial canónico ya debe pensar en media items. |

## 1. Resultado esperado

Dado un path local de repositorio o carpeta, el proceso debe producir un archivo markdown que contenga, como mínimo, estos campos:

- Title
- Summary / Description
- Category
- Client / Context
- Main images
- Published
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

Y además debe producir, aunque sea como secciones de análisis fuente todavía no persistidas 1:1 en la UI:

- Role / Scope
- Stakeholders
- Constraints
- Delivery Strategy
- Risks / Mitigations
- Evidence Sources

Ese archivo debe quedar dentro del repositorio/carpeta analizada y convertirse en la fuente de verdad para la creación del proyecto en PortfolioForge.

Si ese markdown va a alimentar el assistant del proyecto, además debe existir una URL pública HTTPS estable para servirlo en runtime. Las páginas públicas siguen siendo públicas, pero el chat solo se habilita para sesiones elegibles.

Ruta recomendada dentro del proyecto analizado:

- `90. dev_portfolioforge/<Project_Name>.md`

Ejemplo:

- `90. dev_portfolioforge/CAN_Bus_Crane_Monitoring.md`

## 2. Principio rector

La fuente markdown debe ser una **traducción editorial de evidencia real**, no una invención promocional.

Fuentes válidas de evidencia:

- estructura del repositorio;
- README y docs internas;
- nombres de módulos, carpetas y paquetes;
- archivos de configuración y dependencias;
- código fuente;
- capturas, assets y diagramas;
- notas técnicas o documentación adjunta;
- contexto externo confirmado por quien carga el proyecto.

## 3. Flujo recomendado

### Paso 1 — Recibir la ruta y delimitar alcance

Entrada:

- path del repositorio o carpeta principal;
- opcionalmente, contexto adicional confirmado por el autor.

Preguntas mínimas que el proceso debe resolver:

- ¿qué producto/sistema es?
- ¿qué parte implementó realmente el autor?
- ¿qué evidencia existe para sostener problema, solución y resultados?

### Paso 2 — Levantar inventario técnico

Analizar:

- README, docs, diagrams;
- `package.json`, `go.mod`, `requirements`, `Dockerfile`, CI, infra files, etc.;
- estructura de carpetas;
- endpoints, integraciones, protocolos, jobs, workers, servicios;
- assets visuales reutilizables para gallery.

Objetivo:

- identificar stack, arquitectura, integraciones y contexto operativo.

### Paso 3 — Mapear tecnologías

Derivar una lista inicial de tecnologías detectadas y clasificar cada una como:

- tecnología existente en `/admin/technologies`;
- tecnología faltante que debe crearse primero.

Recordatorio del contrato actual de tecnología:

- `name`
- `category`
- `icon`
- `brand color`

Si faltan tecnologías, el flujo correcto es:

1. crear esas tecnologías en `/admin/technologies`;
2. luego cargar el proyecto.

### Paso 4 — Derivar naming editorial del proyecto

Aplicar estas reglas:

- usar nombre corto y técnico;
- no incluir cliente/empresa en el título;
- basar el título en la tecnología, dominio o concepto central.

Ejemplos:

- correcto: `CAN Bus Crane Telemetry`
- correcto: `Industrial Commissioning Support Platform`
- incorrecto: `Proyecto PRECOR`

El cliente va en **Client / Context**.

### Paso 5 — Seleccionar y normalizar imágenes

Preparar assets con esta convención:

- `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_low.webp`
- `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_medium.webp`
- `https://mlbautomation.com/dev/portfolioforge/<project-slug>/imagen01_high.webp`

Reglas:

- mínimo recomendado de **5 imágenes**;
- usar `_medium` como referencia normal de **Main images**;
- `_low` para catálogo;
- `_high` para ampliación;
- mantener naming consistente por orden.

Importante: esta URL base se considera el **default público/canónico de assets** para PortfolioForge. Se documenta separada del contrato editorial para poder reemplazar host/path base más adelante sin cambiar la estructura `imagen0N_<variant>.webp` ni el modelo de media items.

Contrato editorial canónico por ítem de media:

- `low`
- `medium`
- `high`
- `caption`
- `alt_text`
- `featured`
- `sort_order`

Compatibilidad actual:

- el markdown fuente puede mantener un bloque `Main images` con las variantes `_medium` para no romper el flujo actual;
- pero la fuente editorial V2 debe pensar siempre en una colección de media items, no solo en una lista plana de URLs.

### Paso 6 — Mapear evidencia a campos PortfolioForge

#### Campos base

- **Title** → nombre técnico corto
- **Summary / Description** → resumen editorial sustentado en README/docs/código
- **Category** → clasificación reusable
- **Client / Context** → cliente o entorno operativo
- **Main images** → variantes `_medium` de la galería
- **Published** → por defecto `true` en la fuente generada; usar `false` solo si existe una decisión explícita de mantener el proyecto interno hasta nueva revisión

#### Campos enriquecidos

- **Technologies** → stack y plataformas relevantes
- **Business Goal** → por qué existía el proyecto
- **Problem Statement** → dolor o restricción concreta
- **Solution Summary** → qué se construyó y cómo resolvió el problema
- **Architecture** → componentes, datos, protocolos, servicios
- **AI Usage** → solo si hubo uso real de IA
- **Integrations** → sistemas externos, buses, APIs, brokers, vendors
- **Technical Decisions** → tradeoffs y elecciones clave
- **Challenges** → dificultades reales
- **Results** → outcomes observables
- **Metrics** → indicadores con `key: value`
- **Timeline** → hitos o fases

Observación importante: varias de estas secciones hoy se guardan en UI como texto/listas simples. Aun así, la redacción fuente ya debe ser semi-estructurada para mejorar reutilización futura en retrieval/chatbot.

### Paso 7 — Incrustar management/lifecycle donde corresponde

Como la UI actual no tiene campos dedicados para PM/lifecycle, el proceso debe ubicar esa información dentro de los campos disponibles:

- management goals / scope → **Business Goal**
- execution constraints → **Problem Statement**
- delivery tradeoffs → **Technical Decisions**
- milestones / rollout → **Timeline**
- delivery outcomes → **Results** + **Metrics**

Mapping fuente recomendado más explícito:

- role/scope → **Business Goal** + **Solution Summary**
- stakeholders → **Business Goal** o **Problem Statement**
- constraints → **Problem Statement**
- delivery strategy → **Technical Decisions** + **Timeline**
- risks/mitigations → **Challenges**
- evidence sources → **Results**, **Metrics** y referencias explicitas en **Summary / Description**

Estas secciones deben analizarse siempre, aunque todavía no existan como columnas dedicadas.

### Paso 8 — Generar markdown fuente

La salida debe ser un archivo markdown con bloques claramente editables.

Regla operativa:

- ese archivo debe guardarse dentro de la carpeta analizada, idealmente en `90. dev_portfolioforge/`;
- el archivo generado pasa a ser la fuente editorial canónica del proyecto;
- cuando el proyecto deba exponer assistant, también debe resolverse cómo publicar ese markdown en una URL HTTPS alcanzable desde el backend;
- cualquier carga futura a PortfolioForge debe partir de ese `.md`, no de volver a improvisar el contenido desde cero.

### Paso 9 — Consumir el markdown fuente para crear el proyecto en PortfolioForge

Una vez generado el archivo fuente, el siguiente paso del workflow es usarlo como entrada para crear el proyecto en PortfolioForge.

Orden recomendado:

1. leer `90. dev_portfolioforge/<Project_Name>.md`;
2. resolver tecnologías por nombre y mapearlas a `technology_ids` existentes en `/admin/technologies`;
3. mapear los campos editoriales al contrato real actual del admin/storage;
4. si el markdown ya está publicado, escribir esa URL en `source_markdown_url` al crear o actualizar el proyecto;
5. crear o completar el proyecto en PortfolioForge usando ese archivo como fuente de verdad;
6. verificar el resultado real en payload admin/público o DB comparándolo contra el markdown fuente;
7. revisar readiness y publicar cuando corresponda.

Mapping canónico mínimo al momento de cargar:

- `Title` → `name`
- `Summary / Description` → `description`
- `Category` → `category`
- `Client / Context` → `brand` (legacy) / significado funcional `client/context`
- `Published` → `active`
- `Markdown Source URL` → `source_markdown_url` (solo admin/privado; deriva `assistant_available` en público)
- `Technologies` → resolución de nombres → `technology_ids`
- `Main images` + `Media items` → `media` + `images` legacy derivado
- `Business Goal` → `business_goal`
- `Problem Statement` → `problem_statement`
- `Solution Summary` → `solution_summary`
- `Architecture` → `architecture`
- `AI Usage` → `ai_usage`
- `Integrations` → `integrations`
- `Technical Decisions` → `technical_decisions`
- `Challenges` → `challenges`
- `Results` → `results`
- `Metrics` → `metrics`
- `Timeline` → `timeline`

Regla crítica:

- si el `.md` ya existe en `90. dev_portfolioforge/`, el flujo debe preferir consumir ese archivo antes que reanalizar todo el repositorio;
- solo se debe volver a analizar la carpeta si el archivo fuente no existe o está desactualizado respecto a la evidencia nueva.

Regla crítica adicional para assistant:

- si existe markdown fuente utilizable para el assistant, el flujo debe establecer una URL pública HTTPS para servirlo;
- esa URL debe persistirse en `source_markdown_url` durante create/update del proyecto;
- la ausencia de esa URL implica que el proyecto puede quedar bien importado editorialmente, pero sin assistant habilitable.

Regla crítica adicional:

- si durante la importación hay fallback de parseo, omisión de tecnologías, bloques ricos vacíos, sustitución de media o mezcla con assets ajenos, el proceso debe marcarse como **fallido** aunque haya escrituras exitosas en DB.

## 4. Template recomendado de markdown fuente (V2)

```md
# Project Source

## Title
CAN Bus Crane Telemetry

## Summary / Description
Sistema de telemetría y diagnóstico para grúa industrial basado en evidencia de señales CAN Bus, visualización operativa y soporte de troubleshooting en campo.

## Category
industrial automation

## Client / Context
PRECOR

## Main images
https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen01_medium.webp
https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen02_medium.webp
https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen03_medium.webp
https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen04_medium.webp
https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen05_medium.webp

## Media items
- low: https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen01_low.webp | medium: https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen01_medium.webp | high: https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen01_high.webp | caption: Vista general del tablero | alt_text: Tablero principal de telemetría de grúa | featured: true | sort_order: 0
- low: https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen02_low.webp | medium: https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen02_medium.webp | high: https://mlbautomation.com/dev/portfolioforge/can-bus-crane-telemetry/imagen02_high.webp | caption: Señales CAN en diagnóstico | alt_text: Vista de señales CAN normalizadas | featured: false | sort_order: 1

## Published
true

## Technologies
CAN Bus
React
Go
PostgreSQL

## Business Goal
Reducir el tiempo de diagnóstico y mejorar la visibilidad operativa del sistema de grúa.

## Problem Statement
La información de señales y eventos estaba dispersa, dificultando troubleshooting y validación en campo.

## Solution Summary
Se construyó una solución de observabilidad operativa que centraliza señales, eventos y visualización para soporte de diagnóstico técnico.

## Architecture
Frontend de monitoreo + capa backend/API + captura/procesamiento de eventos + persistencia histórica.

## AI Usage


## Integrations
name: CAN Bus gateway | type: protocol gateway | direction: inbound | note: captura de señales de campo
name: MQTT broker | type: messaging | direction: bidirectional | note: distribución de eventos operativos
name: Maintenance API | type: external API | direction: outbound | note: consulta de tickets y soporte

## Technical Decisions
decision: separar captura de señales de visualización operativa | why: aislar carga de adquisición | tradeoff: más coordinación entre componentes
decision: persistir eventos críticos para análisis posterior | why: soportar troubleshooting histórico | tradeoff: mayor costo de almacenamiento

## Challenges
challenge: normalización de señales propietarias | mitigation: diccionario de señales + validación en campo | status: resuelto parcialmente
challenge: ruido y frecuencia variable de eventos | mitigation: filtros y ventanas de agregación | status: mitigado

## Results
result: diagnóstico más rápido | impact: menor tiempo de troubleshooting | evidence: validación de campo y feedback operativo
result: mayor visibilidad de estado operativo | impact: mejor coordinación de soporte | evidence: uso sostenido del tablero

## Metrics
diagnosis_time_pct: -35% | period: commissioning baseline vs rollout
signals_observed_count: 120+ | unit: signals

## Timeline
phase: discovery | objective: relevar señales y contexto operativo | outcome: mapa inicial de componentes y actores
phase: signal mapping | objective: normalizar tramas CAN | outcome: catálogo usable para visualización
phase: implementation | objective: construir interfaz y backend de soporte | outcome: flujo funcional integrado
phase: field validation | objective: probar en entorno real | outcome: ajustes de ruido/frecuencia
phase: rollout | objective: habilitar uso operativo | outcome: adopción inicial y feedback

## Role / Scope
Rol principal del autor, alcance implementado y límites explícitos de responsabilidad.

## Stakeholders
Equipos o actores impactados por el proyecto.

## Constraints
Restricciones técnicas, operativas o regulatorias.

## Delivery Strategy
Estrategia de rollout, validación, coordinación o despliegue.

## Risks / Mitigations
Riesgos observados y cómo se mitigaron.

## Evidence Sources
README, código, assets, documentación técnica, tickets, validación de campo, métricas observadas.
```

## 5. Reglas de serialización obligatorias

Para que el markdown fuente sea compatible con la UI actual:

- `Integrations` → una línea por ítem
- `Technical Decisions` → una línea por ítem
- `Challenges` → una línea por ítem
- `Results` → una línea por ítem
- `Timeline` → una línea por ítem
- `Metrics` → `key: value` una línea por métrica

Pero la convención editorial recomendada dentro de cada línea es:

- `Technical Decisions` → `decision: ... | why: ... | tradeoff: ...`
- `Integrations` → `name: ... | type: ... | direction: ... | note: ...`
- `Results` → `result: ... | impact: ... | evidence: ...`
- `Timeline` → `phase: ... | objective: ... | outcome: ...`
- `Challenges` → `challenge: ... | mitigation: ... | status: ...`
- `Metrics` → `key: value`, con keys normalizadas y unidad/período explícitos cuando sea posible

## 6. Criterios de calidad del análisis

El proceso debe evitar estos errores:

- inventar cliente, resultados o métricas no confirmadas;
- poner el cliente en el título;
- declarar IA si no hubo uso real;
- usar imágenes sin convención consistente;
- pensar `Main images` solo como lista plana y perder metadata de media;
- mezclar gestión/lifecycle en campos inexistentes;
- omitir tecnologías que luego sí aparecen en código o docs.

## 6.1 Verificación obligatoria después de importar

La ingestión automática no termina cuando el insert/update devuelve éxito. Debe existir una comparación explícita entre el markdown fuente y el estado resultante.

Validación mínima obligatoria:

- `title`;
- `Published` del markdown contra `active` real;
- `client/context`;
- `assistant_available=true` cuando exista `source_markdown_url`;
- ausencia de `source_markdown_url` en el payload público;
- tecnologías por cantidad y nombre;
- campos narrativos principales del profile;
- listas enriquecidas no vacías cuando la fuente trae contenido;
- `metrics` presentes;
- `media` e `images` legacy apuntando al proyecto correcto, sin contaminación de assets ajenos.

Chequeos operativos esperados:

- revisar `GET /api/v1/admin/products/:id` o DB para el estado completo;
- revisar `GET /api/v1/public/projects/:slug?lang=es` solo si el resultado quedó activo;
- si existe markdown URL, validar primero que el payload público mantenga `assistant_available=true` sin exponer `source_markdown_url`, luego probar `POST /api/v1/private/projects/:slug/assistant/messages` con una sesión elegible y una pregunta básica;
- si el markdown dice `Published=false`, el chequeo correcto es `404` en el detalle público y ausencia del slug en el listado público.

Reglas duras:

- `Published=false` nunca puede terminar como proyecto activo/público;
- si alguna URL de media no existe o el importador intenta sustituirla por un asset default de otro proyecto, eso es fallo de importación, no fallback aceptable;
- si el parser consume solo parte del markdown, el resultado debe descartarse hasta completar la importación correctamente.

Además, debe marcar explícitamente el nivel alcanzado de calidad:

- **search readiness**: suficiente para ser encontrado en el sistema actual;
- **case study readiness**: suficiente para sostener una narrativa técnica/delivery sólida;
- **assistant readiness**: suficiente para extracción y reuso por retrieval/chatbot futuro.

Hoy el producto implementa principalmente el primer nivel. El assistant ya existe, pero los niveles superiores siguen siendo checkpoints editoriales del workflow.

Caveat operativo importante:

- la calidad runtime del assistant depende de la alcanzabilidad del markdown remoto desde el backend;
- el sistema actual usa cache local y stale fallback para degradación controlada, pero el source of truth sigue siendo la URL remota configurada;
- por eso la reachability del host markdown no es opcional: afecta disponibilidad y frescura de respuestas.

## 7. Ejemplo de referencia metodológica

El caso de grúa/CAN Bus analizado para PRECOR es una buena referencia de proceso:

- el título editorial correcto debe centrarse en el dominio técnico (`CAN Bus`, `telemetry`, `crane operations`), no en el cliente;
- el cliente queda en `Client / Context`;
- arquitectura, integraciones, desafíos y timeline deben salir de evidencia concreta del proyecto.

Ese caso sirve para validar la metodología, pero el producto no debe quedar acoplado a ese cliente ni a un único vertical industrial.

## 8. Uso práctico hoy

Hoy el workflow recomendado es:

1. analizar repo/carpeta;
2. generar markdown fuente en `90. dev_portfolioforge/<Project_Name>.md`;
3. usar ese archivo como fuente de verdad del proyecto;
4. crear tecnologías faltantes en `/admin/technologies`;
5. cargar el proyecto en PortfolioForge a partir del `.md`, siguiendo `docs/MANUAL-PROJECT-INGESTION-WORKFLOW.md` para el mapeo al admin actual;
6. si el markdown debe habilitar assistant, publicar el archivo en una URL HTTPS y persistirla en `source_markdown_url`;
7. verificar campo por campo el resultado real de la importación, incluyendo assistant y su gating autenticado;
8. revisar readiness y publicar cuando esté listo.

## 9. Prompt operativo recomendado

Cuando se quiera trabajar sobre un nuevo repositorio/carpeta, el prompt debe pedir explícitamente dos cosas en secuencia:

1. analizar la evidencia y generar el archivo markdown fuente dentro del proyecto;
2. usar ese mismo archivo para crear el proyecto en PortfolioForge.

Plantilla recomendada:

```text
Analiza la ruta:
"<RUTA_WSL_O_LOCAL_DEL_PROYECTO>"

Debes seguir el workflow de `docs/AUTOMATIC-PROJECT-INGESTION-WORKFLOW.md`.

Objetivos:
1. analizar el repositorio/carpeta y generar el markdown fuente del proyecto dentro de:
   `<RUTA_PROYECTO>/90. dev_portfolioforge/<Project_Name>.md`
2. usar ese archivo `.md` como fuente de verdad para crear un nuevo proyecto en PortfolioForge

Reglas:
- no improvisar contenido fuera de la evidencia real;
- si el archivo `.md` ya existe, usarlo como entrada principal y solo reanalizar la carpeta si hace falta actualizarlo;
- mapear los campos editoriales al contrato real actual de PortfolioForge;
- publicar el markdown fuente en una URL HTTPS y persistirla en `source_markdown_url` cuando se quiera assistant;
- resolver tecnologías por nombre hacia `technology_ids`;
- usar `Main images` como compatibilidad, pero respetar `Media items` como contrato editorial canónico;
- tratar `Client / Context` como mapping a `brand` legacy;
- tratar `Published` como mapping a `active`;
- verificar después de importar el payload admin/público o DB contra el `.md`; además comprobar `assistant_available`, ausencia pública de `source_markdown_url`, invisibilidad del chat en sesión anónima y respuesta básica del endpoint assistant privado con una sesión elegible cuando aplique; si hay fallback parcial o media contaminada, declararlo fallo;
- mantener el criterio técnico + gestión + assistant-readiness definido en la documentación.
```

## 10. Evolución futura sugerida

Más adelante, este markdown fuente puede convertirse en contrato para:

- importador semiautomático;
- validación previa a publicación;
- tooling de extracción asistida desde repositorios reales.

Schema objetivo recomendado para esa evolución:

- **base project**: identidad, resumen, categoría, client/context, publicación
- **media**: variantes, metadata editorial, featured, orden
- **technical profile**: problema, solución, arquitectura, IA, integraciones, decisiones
- **delivery profile**: role/scope, stakeholders, constraints, strategy, risks, timeline, results, metrics, evidence sources
- **assistant/retrieval metadata**: readiness tiers, provenance, confidence, keywords normalizados

Esto debe entenderse como evolución futura, no como implementación actual.
