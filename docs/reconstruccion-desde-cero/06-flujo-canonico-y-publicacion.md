# Flujo canónico, ingestión y publicación

## Tesis central

PortfolioForge no se alimenta directamente del repositorio fuente. Primero existe un markdown canónico en castellano; luego ese markdown se publica en una URL estable; después, desde esa fuente publicada y alineada, se deriva el runtime y se alimenta el assistant.

## Artefacto editorial obligatorio

- ruta local fuente: `90. dev_portfolioforge/<slug>/<slug>.md`
- publicación remota esperada: `https://mlbautomation.com/dev/portfolioforge/<slug>/<slug>.md`

## Flujo de reconstrucción editorial

### 1. Generación del canonical

Un agente o proceso analiza el repositorio/carpeta fuente y produce un único `.md` canónico con estructura estable:

- metadata
- summary
- strategy
- execution
- technical
- media
- validation notes

### 2. Publicación remota

El canonical se publica en una URL HTTPS estable por `slug`. Esa copia remota debe seguir en castellano y coincidir con la local.

### 3. Ingestión a runtime

Desde admin se crea o actualiza el proyecto usando como referencia editorial la misma fuente ya publicada y alineada. La UI runtime debe comprimir y estructurar, no copiar literal.

### 4. Enrichment derivado

Se poblan:

- campos base del proyecto;
- `project_profiles`;
- relaciones con tecnologías;
- media optimizada;
- `source_markdown_url`;
- localizaciones derivadas si aplica.

### 5. Indexación y assistant

- se recompone `project_search_documents`;
- se regeneran embeddings/explanations si cambió contenido indexable;
- assistant queda habilitado solo cuando el markdown remoto es accesible y el proyecto está elegible.

## Reglas editoriales obligatorias

- el markdown responde a fondo; la UI responde rápido;
- cada sección runtime debe condensarse en 1 a 3 ideas cortas cuando sea lista o resumen;
- no inventar hechos no sostenidos por la evidencia fuente;
- sanitizar datos sensibles antes del markdown, no después;
- corregir primero la fuente canónica y luego sincronizar runtime.

## Modelo de publicación

### Fuente local

Es la verdad editorial primaria.

### Fuente remota (`source_markdown_url`)

Es la fuente publicada estable del proyecto. El assistant la consume completa y el runtime debe mantenerse alineado con ella.

### Runtime DB/UI

Es una proyección ejecutiva y estructurada derivada de esa misma fuente publicada, para catálogo, detalle, localización y búsqueda.

## Condiciones de cierre de un caso

Un caso no está completo si ocurre cualquiera de estas situaciones:

- el markdown local está en castellano pero la URL publicada no;
- la UI muestra contenido incompatible con el canonical;
- el proyecto es buscable pero no se reindexó tras cambios sustanciales;
- `assistant_available` se expone sin `source_markdown_url` válido.
