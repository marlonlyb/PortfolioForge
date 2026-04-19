# Alcance y metas del producto

## Qué es PortfolioForge

PortfolioForge es una plataforma editorial y runtime para convertir experiencia profesional real en proyectos públicos navegables, buscables y explicables. No es un repositorio de código expuesto al usuario final ni un CMS genérico.

## Meta principal

Construir un portfolio técnico donde cada proyecto tenga cuatro capas coherentes:

1. evidencia de origen en repositorio o carpeta estudiada;
2. markdown canónico en castellano;
3. representación runtime estructurada en UI/DB;
4. assistant grounded en la publicación remota de ese markdown.

## Usuarios objetivo

### Visitante público

Necesita:

- descubrir proyectos por búsqueda textual, fuzzy y semántica;
- leer catálogo y detalle sin ruido editorial;
- consultar assistant cuando el proyecto lo permita y su sesión sea elegible.

### Administrador/editor

Necesita:

- crear y actualizar proyectos desde una base editorial canónica;
- gestionar tecnologías, media, estado de publicación y localizaciones;
- ejecutar el workflow de importación/publicación del caso;
- mantener searchable content y assistant alineados.

## Capacidades incluidas en el MVP objetivo

- catálogo público de proyectos publicados;
- detalle público por `slug`;
- búsqueda híbrida con explicación breve por resultado;
- assistant por proyecto con acceso autenticado;
- backoffice para proyectos, tecnologías, usuarios y ajustes globales;
- persistencia de detalle enriquecido (`project_profiles`);
- media optimizada por variantes;
- localización persistida derivada desde `es`;
- workflow admin para publicar/importar casos de estudio;
- reindexado/re-embedding cuando cambie contenido indexable.

## Fuera de alcance inicial

- chat libre no grounded;
- edición WYSIWYG del markdown dentro del producto;
- generación automática completa del markdown a partir de carpetas crudas como capacidad cerrada de MVP;
- multi-tenant;
- analítica avanzada o monetización del portfolio.

## Restricciones de producto

- idioma editorial base: castellano;
- assistant usa solo `source_markdown_url`, no la UI resumida ni el repo fuente;
- la UI no debe duplicar el markdown completo;
- el dominio público se expresa como `project`, aunque existan nombres legacy en storage;
- cualquier rebuild debe preservar la distinción entre fuente editorial y representación runtime.

## Criterios de aceptación del rebuild

El rebuild está funcionalmente correcto si permite:

- dar de alta un proyecto con contenido enriquecido y tecnologías;
- publicar `source_markdown_url` coherente con el slug;
- ver el proyecto en catálogo/detalle público;
- buscarlo por nombre, cliente/contexto, tecnologías y contenido clave;
- responder preguntas con assistant usando el markdown remoto;
- exponer versiones públicas en `es`, `ca`, `en`, `de`.
