# Guía maestra para generar `<nombre_del_proyecto>.md`

## 1. Propósito

Este documento es un **prompt operativo** para un agente que debe analizar un repositorio fuente y producir el archivo canónico `<nombre_del_proyecto>.md` que luego será incorporado a PortfolioForge.

El agente **no debe modificar PortfolioForge** ni poblar la UI. Su único entregable es:

- **un único archivo `.md` canónico del proyecto**.

---

## 2. Regla principal

Usa el repositorio fuente **solo** para generar el markdown canónico del proyecto.

No hagas ninguno de estos pasos como parte de esta tarea:

- no crear registros en la UI admin;
- no escribir en la base de datos;
- no generar múltiples documentos intermedios como salida final;
- no producir JSON, YAML o estructuras alternativas como entregable principal;
- no inventar features que no estén respaldadas por el repositorio fuente.

El resultado final debe ser exactamente:

- **un solo `.md` canónico**.

---

## 3. Para qué se usará ese `.md`

Ese markdown tendrá dos usos posteriores:

1. **publicarse en una URL** para el chatbot / assistant mediante `source_markdown_url`;
2. **servir como fuente resumida** para poblar los campos del proyecto en la UI de PortfolioForge.

Por eso el documento debe ser:

- editorialmente claro para humanos;
- fácilmente mapeable a campos estructurados;
- suficientemente rico para responder preguntas en el assistant.

---

## 4. Premisas editoriales obligatorias

1. **La fuente es el repositorio real**: README, código, estructura, configuración, tests, assets y documentación interna.
2. **No inferir sin evidencia**: si algo no puede sostenerse, no debe afirmarse como hecho.
3. **Preferir precisión sobre marketing**.
4. **Separar estrategia, ejecución y técnica**.
5. **El markdown debe ser estable**: no escribirlo como una nota improvisada.
6. **El documento debe poder resumirse luego en UI sin perder sentido**.

---

## 5. Proceso obligatorio del agente

## Paso 1 — Inspección del repositorio fuente

Analiza como mínimo:

- `README` principal y secundarios;
- estructura de carpetas;
- entry points y wiring;
- modelos y contratos;
- integraciones externas;
- tests significativos;
- assets y media relevantes;
- archivos de configuración y despliegue;
- documentación adicional disponible.

## Paso 2 — Extracción de evidencia

Identifica y separa:

- propósito del proyecto;
- problema que resuelve;
- usuarios o contexto operativo;
- alcance de entrega;
- responsabilidades reales asumidas;
- arquitectura y stack;
- integraciones;
- decisiones técnicas;
- riesgos, retos y tradeoffs;
- resultados, métricas o señales de impacto;
- media disponible y reutilizable.

## Paso 3 — Normalización editorial

Convierte la evidencia en narrativa canónica:

- elimina ruido o detalles accidentales;
- consolida naming;
- distingue hechos de interpretación;
- convierte listas dispersas en bloques reutilizables;
- organiza el contenido según las tres capas requeridas.

## Paso 4 — Producción del `.md`

Entrega un solo archivo:

- `<nombre_del_proyecto>.md`

No entregues un resumen alternativo fuera del propio markdown.

---

## 6. Estructura obligatoria del markdown

El markdown debe incluir, como mínimo, estas secciones y este orden lógico.

```md
# <Project Title>

## Metadata
- Slug:
- Published:
- Category:
- Client / Context:
- Repository Source:

## Summary

## Strategy
### Business Goal
### Problem Statement
### Solution Summary

## Execution
### Delivery Scope
### Responsibility Scope
### Challenges
### Results
### Timeline

## Technical
### Architecture
### AI Usage
### Integrations
### Technical Decisions
### Metrics
### Technologies

## Media

## Validation Notes
```

Si falta evidencia para alguna sección, no la inventes: marca la ausencia de forma explícita y breve.

---

## 7. Mapeo obligatorio hacia PortfolioForge

El agente debe escribir el markdown pensando en este mapeo posterior:

| Sección del markdown | Campo runtime esperado en PortfolioForge |
|---|---|
| Title | `name` |
| Summary | `description` |
| Category | `category` |
| Client / Context | `client_name` (storage legacy: `brand`) |
| Business Goal | `business_goal` |
| Problem Statement | `problem_statement` |
| Solution Summary | `solution_summary` |
| Delivery Scope | `delivery_scope` |
| Responsibility Scope | `responsibility_scope` |
| Architecture | `architecture` |
| AI Usage | `ai_usage` |
| Integrations | `integrations` |
| Technical Decisions | `technical_decisions` |
| Challenges | `challenges` |
| Results | `results` |
| Metrics | `metrics` |
| Timeline | `timeline` |
| Technologies | relación con `technology_ids` |
| Media | `media[]` + fallback `images[]` |

---

## 8. Reglas específicas por capa

## 8.1 Strategy

Debe responder:

- ¿qué necesidad existía?
- ¿qué objetivo de negocio u operación motivó el proyecto?
- ¿qué solución se entregó a nivel ejecutivo?

Incluye solo contenido que ayude a entender el **por qué**.

## 8.2 Execution

Debe responder:

- ¿qué se entregó realmente?
- ¿qué parte del trabajo estuvo bajo responsabilidad directa?
- ¿qué retos aparecieron y cómo se manejaron?
- ¿qué resultados o hitos hubo?

Incluye el **cómo se ejecutó y qué impacto tuvo**.

## 8.3 Technical

Debe responder:

- ¿cómo estaba construido?
- ¿qué decisiones técnicas fueron importantes?
- ¿qué integraciones existían?
- ¿qué métricas o señales técnicas ayudan a entenderlo?

Incluye el **cómo funciona**.

---

## 9. Reglas de estilo y serialización

## 9.1 Título

- corto;
- técnico;
- sin meter el cliente dentro del título salvo que sea indispensable para identificar el caso.

## 9.2 Summary

- 1 a 3 párrafos;
- debe poder reutilizarse como descripción pública;
- no repetir mecánicamente todo el documento.

## 9.3 Listas estructuradas

Escribe estas secciones como listas con un ítem por línea:

- Integrations
- Technical Decisions
- Challenges
- Results
- Timeline
- Technologies

Para máxima compatibilidad con PortfolioForge, favorece formatos semi-estructurados como:

- `name: ... | type: ... | note: ...`
- `decision: ... | why: ... | tradeoff: ...`
- `result: ... | impact: ... | evidence: ...`
- `phase: ... | objective: ... | outcome: ...`

## 9.4 Metrics

Escribe métricas como:

- `key: value`

Ejemplos:

- `users_supported: 120/day`
- `latency: <200ms`
- `commissioning_time_reduction: 30%`

## 9.5 Tono

- factual;
- técnico;
- sin claims grandilocuentes;
- con tradeoffs cuando existan.

---

## 10. Convenciones para media

Si el repositorio fuente tiene media reutilizable o si la evidencia permite definirla, documenta cada asset pensado para PortfolioForge con esta intención:

- `low`
- `medium`
- `high`
- `caption`
- `alt_text`
- `featured`
- `sort_order`

Reglas:

1. prioriza imágenes realmente útiles para explicar el proyecto;
2. evita placeholders o assets genéricos;
3. evita media que no pertenezca claramente al proyecto;
4. si no hay media suficiente, deja explícito el gap;
5. si hay suficiente material, intenta dejar definido un set de al menos 5 assets relevantes.

---

## 11. Restricciones

El agente no debe:

- escribir sobre temas no verificables;
- inventar resultados o métricas;
- mezclar proyectos distintos;
- convertir el markdown en documentación interna del repo;
- copiar código extenso innecesario;
- producir un documento genérico sin mapeo claro a PortfolioForge.

---

## 12. Validaciones obligatorias antes de entregar

Antes de considerar terminado el `.md`, verifica que:

1. existe **un solo archivo final**;
2. el documento tiene las tres capas: **Strategy / Execution / Technical**;
3. cada capa contiene evidencia útil y no relleno;
4. el contenido puede mapearse a los campos reales de PortfolioForge;
5. `Delivery Scope` y `Responsibility Scope` están diferenciados;
6. media y tecnologías están normalizadas si existe evidencia;
7. no se exponen suposiciones como hechos;
8. el documento serviría tanto para assistant como para carga manual en UI.

---

## 13. Condiciones de fallo

La tarea debe considerarse fallida si ocurre cualquiera de estas condiciones:

- no se puede producir un único `.md` coherente;
- faltan las capas Strategy / Execution / Technical;
- el documento no puede mapearse razonablemente a PortfolioForge;
- el contenido depende de especulación fuerte;
- la media propuesta está contaminada con assets ajenos;
- el texto final es solo un README reescrito y no un caso editorial canónico;
- el resultado está pensado para documentación del repo, no para PortfolioForge.

---

## 14. Plantilla de salida recomendada

```md
# <Project Title>

## Metadata
- Slug: <project-slug>
- Published: true
- Category: <category>
- Client / Context: <client or context>
- Repository Source: <repo path or URL>

## Summary
<2-4 párrafos con síntesis editorial>

## Strategy

### Business Goal
...

### Problem Statement
...

### Solution Summary
...

## Execution

### Delivery Scope
...

### Responsibility Scope
...

### Challenges
- challenge: ... | mitigation: ... | status: ...

### Results
- result: ... | impact: ... | evidence: ...

### Timeline
- phase: ... | objective: ... | outcome: ...

## Technical

### Architecture
...

### AI Usage
...

### Integrations
- name: ... | type: ... | direction: ... | note: ...

### Technical Decisions
- decision: ... | why: ... | tradeoff: ...

### Metrics
- key: value

### Technologies
- <technology>

## Media
- asset: ... | low: ... | medium: ... | high: ... | caption: ... | alt_text: ... | featured: true | sort_order: 0

## Validation Notes
- confirmed_from_readme: yes/no
- confirmed_from_code: yes/no
- gaps_or_uncertainties: ...
```

---

## 15. Instrucción final para el agente

Analiza el repositorio fuente como evidencia.

Genera **solo** el archivo canónico `<nombre_del_proyecto>.md`.

Ese archivo debe quedar listo para:

1. publicarse luego en una URL para `source_markdown_url`;
2. servir como base editorial para poblar el proyecto en PortfolioForge.

Si la evidencia no alcanza para producir un markdown canónico confiable, debes declararlo explícitamente en `Validation Notes` en lugar de inventar contenido.
