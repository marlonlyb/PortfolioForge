# Guía única para poblar la UI desde el markdown canónico

## 1. Objetivo

Esta guía define **cómo llenar y actualizar la UI admin de PortfolioForge** tomando como fuente de verdad un **markdown canónico ya generado** y, cuando aplique, publicado en `source_markdown_url` para el assistant.

Su propósito no es explicar cómo analizar un repositorio ni cómo escribir el markdown canónico. Su propósito es uno solo:

- **convertir el markdown canónico en una versión runtime resumida, estructurada y visualmente limpia dentro de la UI/DB**.

Regla central:

- **el markdown canónico responde a fondo; la UI responde rápido**.

---

## 2. Alcance

Esta guía cubre:

- creación manual de proyectos en la UI;
- actualización de proyectos existentes a partir del markdown canónico;
- criterios de compresión editorial para evitar ruido visual;
- reglas de mapping hacia campos reales del sistema;
- verificación posterior en payload admin/público o DB.

Esta guía **no** cubre:

- cómo generar el markdown canónico desde un repositorio/carpeta;
- cómo redactar el `.md` fuente en detalle;
- cómo diseñar el producto completo.

Para eso ya existen:

- `docs/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
- `docs/PRD.md`

Para trabajo operativo dentro de este repo, esta guía debe tratarse como la **referencia única** del flujo:

- **markdown canónico → UI runtime**

---

## 3. Principio rector

### 3.1 Fuente de verdad

Si un proyecto ya tiene markdown canónico confiable, ese archivo se usa primero y se trata como:

- fuente editorial principal;
- fuente de `source_markdown_url` para el assistant;
- referencia principal para poblar o corregir la UI.

No corresponde re-analizar el repositorio completo salvo que:

- el markdown no exista;
- el markdown esté desactualizado;
- haya evidencia nueva relevante.

### 3.2 La UI no duplica el markdown

La UI **no** debe copiar el markdown completo ni intentar competir con él.

La UI debe persistir:

- la **versión breve**;
- la **versión estructurada**;
- la **versión visualmente clara**;
- la **versión útil para catálogo, detalle, búsqueda, localización y readiness**.

### 3.3 El assistant no lee la UI resumida

El assistant usa el markdown remoto publicado en `source_markdown_url`.

Por eso:

- el markdown puede y debe conservar más contexto;
- la UI debe quedarse con la información clave y operativa.

---

## 4. Regla editorial de compresión

### 4.1 Qué sí pasa a la UI

Debe pasar a la UI solo lo que permita entender rápido:

- qué es el proyecto;
- qué problema resolvió;
- qué solución se entregó;
- qué se hizo realmente;
- cómo estaba construido a nivel útil;
- qué decisiones, retos y resultados son más importantes;
- qué tecnologías son realmente nucleares.

### 4.2 Qué debe quedarse solo en el markdown canónico

Debe quedarse solo en el markdown:

- contexto largo;
- matices históricos completos;
- explicaciones extensas;
- trazabilidad detallada;
- evidencia narrativa rica;
- aclaraciones largas de alcance;
- justificaciones documentales o comparativas extensas.

### 4.3 Regla práctica

Si un texto en UI obliga al usuario a “leer demasiado para entender lo principal”, está mal comprimido.

La UI debe mostrar:

- **lo suficiente para decidir**;
- **no todo lo disponible para documentar**.

---

## 5. Flujo correcto hoy

1. leer el markdown canónico completo;
2. confirmar si ese markdown sigue siendo la fuente correcta;
3. resolver tecnologías antes de guardar el proyecto;
4. poblar campos base del proyecto;
5. poblar `project_profiles` con versión resumida y estructurada;
6. cargar o revisar media solo si el markdown la define con claridad suficiente;
7. guardar `source_markdown_url` si el proyecto debe exponer assistant;
8. verificar el resultado real en admin/público o DB;
9. refrescar búsqueda y re-embed cuando cambie contenido indexable.

---

## 6. Reglas transversales obligatorias

### 6.1 No inventar

Si el markdown no sostiene algo como hecho:

- no se agrega a la UI.

### 6.2 No expandir innecesariamente

Si la UI ya tiene una versión correcta y breve de un bloque:

- no se reemplaza por una versión más larga solo porque el markdown tenga más detalle.

### 6.3 No crear una segunda autoría

Si hay diferencia entre UI y markdown:

- primero se corrige o confirma la fuente;
- luego se sincroniza la UI.

### 6.4 Mantener privacidad y sanitización

La UI runtime debe respetar las mismas restricciones del markdown canónico:

- sin nombres personales cuando baste el rol;
- sin identificadores documentales internos;
- sin montos exactos;
- sin anexos de trazabilidad;
- sin secretos o infraestructura sensible.

### 6.5 Respetar el contrato real del sistema

Mapping operativo actual:

- `Client / Context` → `brand` (legacy admin/storage; consumo público `client_name`)
- `Published` → `active`
- `Technologies` → `technology_ids`
- `Media` → `media[]` + `images` legacy derivado
- `Markdown source URL` → `source_markdown_url`

---

## 7. Guía campo por campo

## A. Project profile

### 7.1 `name`

Fuente: `Title`

Qué debe quedar:

- corto;
- técnico;
- claro;
- sin cliente en el título salvo necesidad real.

Regla:

- el nombre en UI debe ser **más identificable que comercial**.

### 7.2 `description`

Fuente: `Summary`

Qué debe quedar:

- **1 a 2 párrafos**;
- contexto funcional;
- solución entregada;
- valor generado.

Qué evitar:

- pegar todo el summary completo si es muy largo;
- repetir detalles ya explicados mejor en otras secciones;
- ruido histórico innecesario.

Objetivo:

- este campo debe servir a la vez para catálogo, overview público y búsqueda.

### 7.3 `category`

Fuente: `Category`

Qué debe quedar:

- una sola categoría consistente y reusable.

Regla:

- preferir taxonomías estables como `industrial automation`, `embedded systems`, `web platform`, etc.

### 7.4 `brand` / Client-Context

Fuente: `Client / Context`

Qué debe quedar:

- contexto breve;
- legible;
- útil para comprender el proyecto sin saturar la UI.

Regla:

- no usar una frase excesivamente larga si puede comprimirse sin perder sentido;
- si el nombre del cliente no debe publicarse, usar contexto operativo sanitizado.

### 7.5 `source_markdown_url`

Fuente: URL pública del markdown canónico

Qué debe quedar:

- URL HTTPS exacta del markdown fuente;
- solo en admin/privado.

Reglas:

- no se expone en la API pública;
- si existe y no está vacía, el proyecto puede declarar `assistant_available=true` en público;
- el assistant sigue requiriendo sesión elegible.

### 7.6 `active`

Fuente: `Published`

Regla:

- `Published=true` → `active=true`
- `Published=false` → `active=false`

No reinterpretar esto por “se ve bien” o readiness parcial.

### 7.7 `media`

Fuente: `Media`

Contrato editorial esperado por ítem:

- `low_url`
- `medium_url`
- `high_url`
- `fallback_url`
- `caption`
- `alt_text`
- `featured`
- `sort_order`

Uso real en frontend:

- `low_url` → card/catálogo
- `medium_url` → hero/galería del detalle
- `high_url` → lightbox/ampliación
- `fallback_url` → respaldo final

Reglas:

- si el markdown define media estructurada, se sincroniza;
- si el markdown **no** define media estructurada de forma suficiente, **no se borra media existente por defecto**;
- limpiar o reemplazar media debe ser una decisión explícita, no un efecto colateral de una sincronización de texto.

---

## B. Rich Profile / Search Enrichment

### 7.8 `technology_ids`

Fuente: `Technologies`

Qué debe quedar:

- solo tecnologías nucleares del caso;
- protocolos, plataformas, stacks y componentes realmente importantes.

Reglas:

- si una tecnología no existe, se crea primero;
- evitar inflar el proyecto con tecnologías marginales o incidentales.

### 7.9 `business_goal`

Fuente: `Business Goal`

Qué debe quedar:

- **1 párrafo corto**;
- objetivo de negocio, operación o entrega;
- el “por qué” del proyecto.

### 7.10 `problem_statement`

Fuente: `Problem Statement`

Qué debe quedar:

- **1 párrafo corto**;
- problema previo;
- dolor operativo o restricción principal.

### 7.11 `solution_summary`

Fuente: `Solution Summary`

Qué debe quedar:

- **1 párrafo corto**;
- solución principal;
- enfoque de alto nivel;
- cómo resolvió el problema.

Regla:

- este es uno de los campos más importantes para búsqueda y readiness.

### 7.12 `delivery_scope`

Fuente: `Delivery Scope`

Qué debe quedar:

- lista breve de lo entregado realmente;
- preferir **3 a 6 bullets**.

Regla:

- no listar todo lo posible;
- priorizar alcance visible y verificable.

### 7.13 `responsibility_scope`

Fuente: `Responsibility Scope`

Qué debe quedar:

- lista breve de responsabilidad real;
- idealmente **3 a 6 bullets**;
- puede incluir una subsección muy corta de “no evidenciado” si es importante para no sobredeclarar.

### 7.14 `architecture`

Fuente: `Architecture`

Qué debe quedar:

- **1 a 2 párrafos**;
- componentes principales;
- flujo general de datos;
- capas y protocolos esenciales.

Qué evitar:

- convertirlo en especificación exhaustiva.

### 7.15 `ai_usage`

Fuente: `AI Usage`

Regla:

- si no hubo uso real de IA, dejarlo vacío o claramente negativo y breve;
- no forzar este campo.

### 7.16 `integrations`

Fuente: `Integrations`

Formato recomendado:

- una línea por integración;
- convención: `name: ... | type: ... | direction: ... | note: ...`

Regla de compresión:

- preferir **3 a 5 integraciones** importantes;
- no enumerar integraciones irrelevantes o incidentales.

### 7.17 `technical_decisions`

Fuente: `Technical Decisions`

Formato recomendado:

- una línea por decisión;
- convención: `decision: ... | why: ... | tradeoff: ...`

Regla de compresión:

- preferir **3 a 5 decisiones** realmente estructurales.

### 7.18 `challenges`

Fuente: `Challenges`

Formato recomendado:

- una línea por desafío;
- convención: `challenge: ... | mitigation: ... | status: ...`

Regla de compresión:

- preferir **3 a 5 retos** significativos.

### 7.19 `results`

Fuente: `Results`

Formato recomendado:

- una línea por resultado;
- convención: `result: ... | impact: ... | evidence: ...`

Regla de compresión:

- preferir **3 a 5 resultados** de alto valor.

### 7.20 `metrics`

Fuente: `Metrics`

Formato obligatorio:

- `key: value`

Reglas:

- usar pocas métricas y de alto valor;
- priorizar porcentajes, ratios, comparativas o unidades operativas seguras;
- evitar claves vagas o métricas sin evidencia.

Recomendación:

- normalmente **3 a 6 métricas** bastan.

### 7.21 `timeline`

Fuente: `Timeline`

Formato recomendado:

- una línea por fase;
- convención: `phase: ... | objective: ... | outcome: ...`

Regla de compresión:

- preferir **3 a 5 hitos**;
- no convertir la UI en cronología exhaustiva.

---

## 8. Dónde poner la información que no tiene campo propio

La UI actual no tiene columnas dedicadas para varios conceptos del markdown. Deben redistribuirse así:

- rol y alcance real → `business_goal` + `solution_summary`
- stakeholders → `business_goal` o `problem_statement`, por rol, nunca por nombre
- restricciones → `problem_statement`
- estrategia de entrega → `technical_decisions` + `timeline`
- riesgos y mitigaciones → `challenges`
- outcomes de delivery → `results` + `metrics`

No crear campos nuevos fuera del contrato runtime actual.

---

## 9. Reglas de búsqueda, localización y assistant

### 9.1 Búsqueda

Cambios en estos campos son indexables y exigen refresh de búsqueda:

- `name`
- `description`
- `brand` / client-context
- `solution_summary`
- `architecture`
- `business_goal`
- `problem_statement`
- `ai_usage`
- tecnologías

Además:

- normalmente conviene re-embed después de cambios relevantes en contenido indexable;
- la búsqueda runtime usa DB, **no** el markdown remoto directamente.

### 9.2 Localización

Reglas:

- el contenido base vive en español;
- los cambios deben sincronizar traducciones automáticas cuando aplique;
- las traducciones manuales no deben sobrescribirse;
- no dejar filas auto vacías que degraden el fallback de arrays, métricas o listas.

### 9.3 Assistant

Reglas:

- el assistant responde sobre el markdown remoto en `source_markdown_url`;
- la UI solo controla si el proyecto está preparado para declararlo disponible;
- `source_markdown_url` nunca debe filtrarse públicamente.

---

## 10. Verificación obligatoria después de guardar

Guardar en DB no basta.

Después de crear o actualizar un proyecto, verificar al menos:

- `name`
- `description`
- `category`
- `brand/client-context`
- `active`
- `source_markdown_url` si aplica
- tecnologías asociadas
- `project_profiles` completos según el markdown
- media correcta y no contaminada
- `assistant_available` si corresponde
- ausencia pública de `source_markdown_url`
- búsqueda refrescada
- localizaciones coherentes

Checklist mínimo:

- [ ] la UI refleja la **síntesis correcta** del markdown y no una copia literal excesiva
- [ ] no faltan campos clave que sí están en el markdown
- [ ] no se agregaron claims nuevos sin evidencia
- [ ] las listas enriquecidas conservaron los ítems importantes
- [ ] `metrics` conserva sus claves y formato
- [ ] la media pertenece al proyecto correcto
- [ ] si el markdown no define media suficiente, no se borró media existente por accidente
- [ ] `assistant_available` coincide con la presencia real de `source_markdown_url`

---

## 11. Condiciones de fallo

La carga o actualización debe tratarse como fallida si ocurre cualquiera de estas:

- se copia el markdown casi completo a la UI y se vuelve visualmente ruidosa;
- faltan tecnologías, bloques enriquecidos o narrativa clave que sí están en el markdown;
- se inventa contenido en la UI sin respaldo en el markdown o en evidencia nueva verificada;
- se mezclan assets ajenos o media contaminada;
- se borra media útil sin decisión explícita;
- `source_markdown_url` queda mal configurado o expuesto públicamente;
- cambia contenido indexable y no se refresca búsqueda;
- localizaciones vacías o mal sincronizadas pisan el contenido base.

---

## 12. Regla final de trabajo

Cuando exista markdown canónico confiable:

1. **se lee primero**;
2. **se resume para UI**;
3. **no se reescribe libremente**;
4. **no se copia entero**;
5. **se verifica el resultado real**;
6. **el assistant sigue usando el markdown remoto completo**.

En una frase:

- **markdown canónico = profundidad**
- **UI runtime = claridad resumida**
- **assistant = grounding sobre la fuente remota**
