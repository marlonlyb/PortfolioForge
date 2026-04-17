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

Y además debe ser **seguro para publicación**:

- no debe exponer información sensible, contractual, personal u operativa innecesaria;
- debe redactar o generalizar detalles que sirvan para comprender el caso, pero no para revelar datos privados del cliente o del proyecto;
- debe tratar el repositorio fuente como evidencia de trabajo, no como material a volcar literalmente en el markdown final.

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
- sanitiza datos sensibles o confidenciales antes de escribir;
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

Notas obligatorias sobre privacidad:

- `Client / Context` puede usar el nombre del cliente solo si ya es público o está explícitamente autorizado; si no, debe describirse el contexto operativo de forma genérica.
- `Repository Source` no debe exponer rutas internas ni nombres de carpetas privadas. Si el repositorio no es público, usa una fórmula genérica como `private repository analyzed`.
- No agregues secciones como `Evidence Sources`, anexos de trazabilidad documental, listados de archivos internos ni `Stakeholders` con nombres propios.

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

Regla de sanitización del mapeo:

- `Client / Context` debe priorizar contexto operativo sobre razón social si el nombre del cliente no es público.
- Cualquier referencia a personas debe escribirse por rol o función, nunca por nombre propio, correo o teléfono.

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

Usa solo formatos seguros para publicación:

- porcentajes;
- ratios;
- variaciones comparativas;
- unidades operativas no sensibles;
- comparaciones antes/después.

No incluyas:

- montos de inversión;
- costos, tarifas o presupuestos exactos;
- números de cotización, propuesta, factura u orden de compra;
- identificadores documentales usados solo internamente.

Ejemplos:

- `users_supported: 120/day`
- `latency: <200ms`
- `commissioning_time_reduction: 30%`

Mejores ejemplos cuando el dato puede ser sensible:

- `manual_intervention_reduction: 40%`
- `commissioning_time_reduction: 30%`
- `validation_cycles: 2x faster than previous setup`
- `simultaneous_crane_contexts_supported: 2 operating modes`

## 9.5 Sanitización obligatoria

Antes de escribir el markdown final, transforma o elimina cualquier dato sensible.

No mostrar literalmente:

- nombres y apellidos de stakeholders, operadores, supervisores o contrapartes;
- correos, teléfonos, firmas, usuarios o identificadores personales;
- números de cotización, presupuestos, propuestas, informes, órdenes de compra o tickets;
- montos económicos absolutos o condiciones comerciales;
- rutas internas, nombres exactos de archivos privados, backups, carpetas del cliente o referencias documentales crudas;
- `Evidence Sources`, listados de archivos usados, ni trazabilidad documental detallada;
- IPs internas, hostnames, credenciales, tokens, serial numbers o identificadores de activos;
- logs, tramas, dumps o payloads crudos cuando revelen datos operativos sensibles;
- volúmenes exactos de operación, capacidad instalada o fechas internas si no son necesarias para entender el caso.

Reemplazar por:

- roles o funciones en vez de nombres propios (`contraparte técnica`, `supervisión de planta`, `proveedor de integración`);
- tipos documentales genéricos (`propuesta técnica inicial`, `informe de validación`, `manual del sistema`);
- porcentajes, ratios o comparativas en vez de montos absolutos;
- descripciones editoriales de evidencia en vez de rutas o nombres de archivo;
- contexto operativo genérico cuando el nombre del cliente o de terceros no deba publicarse.

## 9.6 Tono

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

Semántica real en PortfolioForge:

- `low`: card / miniatura del catálogo;
- `medium`: imagen principal del detalle / galería;
- `high`: vista ampliada / lightbox;
- `featured`: define qué asset aparece primero;
- `sort_order`: orden restante después de `featured`.

Reglas:

1. prioriza imágenes realmente útiles para explicar el proyecto;
2. evita placeholders o assets genéricos;
3. evita media que no pertenezca claramente al proyecto;
4. si no hay media suficiente, deja explícito el gap;
5. si hay suficiente material, intenta dejar definido un set de al menos 5 assets relevantes;
6. no uses etiquetas ambiguas como `Main images` si solo enumeran una variante; documenta cada asset con `low | medium | high`.

---

## 11. Restricciones

El agente no debe:

- escribir sobre temas no verificables;
- inventar resultados o métricas;
- mezclar proyectos distintos;
- convertir el markdown en documentación interna del repo;
- copiar código extenso innecesario;
- producir un documento genérico sin mapeo claro a PortfolioForge;
- incluir nombres propios de stakeholders o terceros cuando baste el rol;
- incluir números de documentos, cotizaciones, presupuestos, propuestas o tickets internos;
- incluir montos económicos exactos o condiciones comerciales;
- incluir secciones de `Evidence Sources`, rutas internas, listados de archivos o anexos de trazabilidad sensible;
- exponer secretos, infraestructura interna o datos operativos no necesarios para comprender el caso.

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
8. el documento serviría tanto para assistant como para carga manual en UI;
9. no aparecen nombres propios, IDs documentales, montos exactos ni rutas internas innecesarias;
10. no existe una sección `Evidence Sources` ni trazabilidad documental sensible;
11. las métricas están expresadas como porcentajes, ratios, comparativas o unidades operativas seguras;
12. si `Repository Source` aparece, está sanitizado y no expone rutas privadas.

---

## 13. Condiciones de fallo

La tarea debe considerarse fallida si ocurre cualquiera de estas condiciones:

- no se puede producir un único `.md` coherente;
- faltan las capas Strategy / Execution / Technical;
- el documento no puede mapearse razonablemente a PortfolioForge;
- el contenido depende de especulación fuerte;
- la media propuesta está contaminada con assets ajenos;
- el texto final es solo un README reescrito y no un caso editorial canónico;
- el resultado está pensado para documentación del repo, no para PortfolioForge;
- el markdown expone nombres personales, identificadores documentales, montos exactos o rutas internas sensibles;
- el documento incluye `Evidence Sources` o trazabilidad documental cruda como si fuera un anexo interno.

---

## 14. Plantilla de salida recomendada

```md
# <Project Title>

## Metadata
- Slug: <project-slug>
- Published: true
- Category: <category>
- Client / Context: <public client name if already public, otherwise generic operating context>
- Repository Source: <public repo URL if already public, otherwise "private repository analyzed">

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
- privacy_sanitized: yes/no
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
