# Guía única para poblar la UI desde el markdown canónico

## 1. Objetivo

Esta guía define **cómo llenar y actualizar la UI admin de PortfolioForge** tomando como punto de partida operativo la **URL canónica publicada** del proyecto (`source_markdown_url`).

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

- `docs/operacion/CANONICAL-PROJECT-MARKDOWN-AGENT-GUIDE.md`
- `docs/PRD.md`
- `docs/README.md`

Para trabajo operativo dentro de este repo, esta guía debe tratarse como la **referencia única** del flujo:

- **markdown canónico → UI runtime**

---

## 3. Principio rector

### 3.1 Fuente de verdad

Para esta tarea operativa de runtime, la fuente que se lee primero es la **URL canónica publicada** del proyecto.

Ese documento remoto se trata como:

- fuente editorial principal;
- fuente base en castellano para el contenido runtime;
- fuente remota completa que también consume el assistant;
- referencia principal para poblar o corregir la UI.

Punto de partida obligatorio en este flujo:

- si el proyecto ya existe, se parte de la `source_markdown_url` ya guardada en admin;
- si el proyecto todavía no existe, hay que pedir primero la URL canónica publicada antes de leer o crear nada;
- esa URL no la publica PortfolioForge como flujo estándar: la publica manualmente personal externo del usuario en su host en la nube;
- para esta tarea no se debe buscar el canonical en la PC, discos locales ni rutas auxiliares dentro de `PortfolioForge`.

Regla de coherencia obligatoria:

- la URL publicada en `source_markdown_url` debe estar en castellano y representar la versión editorial vigente del proyecto;
- si la URL publicada sigue en inglés u otro idioma, la fuente todavía está desalineada y no debe darse por cerrada;
- cualquier verificación contra una fuente local fuera de la UI/runtime pertenece al flujo editorial/publish, no al primer paso de esta tarea.
- la convención de publicación remota debe ser estable y basada en slug: `https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md`.

Diferencia importante:

- el assistant consume la copia remota completa publicada en `source_markdown_url`;
- la UI/DB no debe copiar ese markdown literal, sino derivar una representación resumida, estructurada y operativa del mismo contenido.

No corresponde re-analizar el repositorio completo ni buscar archivos canónicos locales salvo que:

- el markdown no exista;
- el markdown esté desactualizado;
- haya evidencia nueva relevante.

### 3.2 La UI no duplica el markdown

La UI **no** debe copiar el markdown completo ni intentar competir con él.

Tampoco debe reformatear el markdown casi literal cambiando solo el HTML o metiendo prefijos repetitivos.

La UI debe persistir:

- la **versión breve**;
- la **versión estructurada**;
- la **versión visualmente clara**;
- la **versión útil para catálogo, detalle, búsqueda, localización y readiness**.
- la **versión sintetizada**, con solo los puntos clave por sección.

Contrato explícito:

- el markdown canónico es la fuente de verdad completa;
- la UI runtime es una vista ejecutiva derivada;
- cada sección de UI debe condensarse en **1 a 3 ideas breves** o bullets cortos;
- queda prohibido volcar una línea larga del canónico casi textual dentro de la UI;
- queda prohibido repetir prefijos redundantes como `decision:`, `challenge:`, `result:`, `phase:`, `why:`, `impact:`, `evidence:` o equivalentes si el contenedor ya aporta el contexto.

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

Si un bullet de UI se parece demasiado a una línea completa del markdown canónico, también está mal comprimido.

La UI debe mostrar:

- **lo suficiente para decidir**;
- **no todo lo disponible para documentar**.

---

## 5. Flujo correcto hoy

### Resumen operativo corto

1. obtener o confirmar `source_markdown_url`;
2. leer el markdown remoto completo desde esa URL;
3. confirmar que la URL publicada esté accesible, en castellano y siga siendo la fuente correcta;
4. resolver tecnologías antes de guardar el proyecto;
5. poblar campos base del proyecto usando una proyección resumida y estructurada de esa fuente remota;
6. poblar `project_profiles` con versión resumida y estructurada;
7. cargar o revisar media solo si el markdown la define con claridad suficiente;
8. guardar o mantener `source_markdown_url` si el proyecto debe exponer assistant;
9. si cambió la base en castellano, regenerar localizaciones derivadas (`ca`, `en`, `de`) desde esa base;
10. verificar el resultado real en admin/público o DB;
11. refrescar búsqueda y re-embed cuando cambie contenido indexable.

Precondición explícita:

- si el proyecto todavía no existe y la URL final aún no fue publicada en el host externo del usuario, el flujo runtime/UI debe detenerse hasta que esa publicación manual exista.

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
- `Industry Type` → `industry_type`
- `Final Product` → `final_product`
- `Published` → `active`
- `Technologies` → `technology_ids`
- `Media` → `media[]` + `images` legacy derivado
- `Markdown source URL` → `source_markdown_url`

Regla de localización asociada:

- el contenido base siempre se escribe y corrige primero en castellano (`es`);
- las versiones públicas `ca`, `en` y `de` se derivan desde ese castellano;
- `client_name` participa en esa misma localización pública, aunque storage/admin legacy aún acepten `brand`.
- identificadores técnicos, marcas o nombres de plataforma pueden mantenerse en idioma original si traducirlos degrada precisión.

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

### 7.2.1 `industry_type`

Fuente: frontmatter `industry_type` + bloque visible `## Metadata`

Reglas:

- normalizar whitespace y persistir el valor base en castellano como texto editorial corto;
- aceptar keys legacy solo como compatibilidad transicional y convertirlas a copy editorial ES antes de guardar;
- límite máximo: 160 caracteres;
- participar en localización pública/admin igual que `final_product`, con fallback a `es` cuando no exista override.

### 7.2.2 `final_product`

Fuente: frontmatter `final_product` + bloque visible `## Metadata`

Reglas:

- debe quedar como una frase corta y específica del entregable final;
- se persiste top-level en `products.final_product`;
- participa en localización pública derivada (`ca`, `en`, `de`) igual que el resto del copy editorial corto.
- ruido histórico innecesario.

Objetivo:

- este campo debe servir a la vez para catálogo, overview público y búsqueda.

### 7.3 `category`

Fuente: `Category`

Qué debe quedar:

- una sola categoría consistente y reusable.

Regla:

- preferir taxonomías estables como `industrial automation`, `embedded systems`, `web platform`, etc.

### 7.4 `client_name` / `brand` / Client-Context

Fuente: `Client / Context`

Qué debe quedar:

- contexto breve;
- legible;
- útil para comprender el proyecto sin saturar la UI.

Regla:

- no usar una frase excesivamente larga si puede comprimirse sin perder sentido;
- si el nombre del cliente no debe publicarse, usar contexto operativo sanitizado.
- ese valor en castellano es la fuente de verdad; otras locales se regeneran desde aquí y cualquier ajuste no automático debe guardarse como override manual.

Excepción:

- si el campo contiene un identificador técnico, un nombre de producto o una marca cuyo cambio degrade precisión, se puede mantener igual entre idiomas.

### 7.5 `source_markdown_url`

Fuente: URL pública del markdown canónico

Qué debe quedar:

- URL HTTPS exacta del markdown fuente;
- solo en admin/privado.

Reglas:

- no se expone en la API pública;
- si existe y no está vacía, el proyecto puede declarar `assistant_available=true` en público;
- el assistant sigue requiriendo sesión elegible.
- la convención por defecto debe ser: `https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md`.
- `<slug>` debe coincidir con el `Slug` del markdown canónico y con el `slug` runtime del proyecto.
- si se va a crear un proyecto nuevo y todavía no existe esta URL, el flujo debe detenerse para pedirla antes de continuar.

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
- cuando el proyecto se crea o actualiza automáticamente desde el canonical y todavía no tiene media persistida suficiente, el sistema debe sembrar por defecto **7 imágenes** usando la convención `https://mlbautomation.com/dev/portfolioforge/<slug>/imagenNN_{low|medium|high}.webp` y el fallback global `https://mlbautomation.com/dev/portfolioforge/imagen_fallback/Logo_500_500.png`;
- ese sembrado automático debe respetar media manual existente y solo completar faltantes hasta llegar al mínimo esperado;
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
- texto directo y breve, sin prefijos redundantes.

Ejemplo correcto:

- `CAN Bus como backbone entre la medición existente y la estación de monitoreo`

Regla de compresión:

- preferir **3 a 5 integraciones** importantes;
- no enumerar integraciones irrelevantes o incidentales.

### 7.17 `technical_decisions`

Fuente: `Technical Decisions`

Formato recomendado:

- una línea por decisión;
- texto directo y breve, sin `decision:`, `why:` ni `tradeoff:`.

Ejemplo correcto:

- `Exponer datos por Ethernet/UDP para consumo externo futuro sin afirmar una integración ERP ya completada`

Regla de compresión:

- preferir **3 a 5 decisiones** realmente estructurales.

### 7.18 `challenges`

Fuente: `Challenges`

Formato recomendado:

- una línea por desafío;
- texto directo y breve, sin `challenge:`, `mitigation:` ni `status:`.

Ejemplo correcto:

- `La comunicación USB sobre UTP no era confiable para la distancia requerida y se rediseñó alrededor de CAN directa`

Regla de compresión:

- preferir **3 a 5 retos** significativos.

### 7.19 `results`

Fuente: `Results`

Formato recomendado:

- una línea por resultado;
- texto directo y breve, sin `result:`, `impact:` ni `evidence:`.

Ejemplo correcto:

- `Se instaló visualización en dos pantallas para ambos contextos de grúa y mejoró la supervisión desde piso`

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
- texto directo y breve, sin `phase:`, `objective:` ni `outcome:`.

Ejemplo correcto:

- `A inicios de 2024 quedó documentado el retrofit instalado y la arquitectura posterior a la modificación`

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
- no dejar filas auto vacías que degraden el fallback de arrays, métricas o listas;
- si el markdown canónico estaba en otro idioma por error, primero se corrige la base en castellano y recién después se regeneran `ca`, `en` y `de`.
- la URL publicada en `source_markdown_url` debe reflejar también esa base en castellano; no debe quedar una copia remota en otro idioma.

### 9.3 Assistant

Reglas:

- el assistant responde sobre el markdown remoto en `source_markdown_url`;
- la UI solo controla si el proyecto está preparado para declararlo disponible;
- `source_markdown_url` nunca debe filtrarse públicamente.
- si el assistant debe responder en coherencia con la base española, la copia remota publicada también debe estar en castellano.
- para este flujo operativo, esa URL remota es también el punto de partida de lectura para poblar o corregir la UI.
- FTPS o `canonical-publish` no forman parte del camino principal documentado para esta tarea; si existen, pertenecen a compatibilidad legacy.

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
- se intenta poblar o corregir la UI sin `source_markdown_url` cuando el proyecto todavía no existe y no se pidió la URL al usuario;
- cambia contenido indexable y no se refresca búsqueda;
- localizaciones vacías o mal sincronizadas pisan el contenido base.

---

## 12. Regla final de trabajo

Cuando exista markdown canónico confiable y publicado:

1. **se obtiene o confirma primero `source_markdown_url`**;
2. **se lee primero esa fuente remota completa**;
3. **se resume para UI**;
4. **no se reescribe libremente**;
5. **no se copia entero**;
6. **se verifica el resultado real**;
7. **el assistant sigue usando el mismo markdown remoto completo**;
8. **no se buscan canonicals en PC, discos locales ni copias auxiliares dentro de `PortfolioForge` como primer paso de esta tarea**.

En una frase:

- **`source_markdown_url` = punto de partida operativo para esta tarea**
- **UI runtime = claridad resumida derivada de esa fuente**
- **assistant = grounding sobre esa fuente remota completa**
