# PRD - PortfolioForge

## 1. Resumen

PortfolioForge es una plataforma de portfolio interactivo construida con React, Go y PostgreSQL para presentar proyectos como case studies técnicos y de negocio, con una experiencia pública centrada en búsqueda inteligente y una consola privada para administración de contenido.

La interacción principal del producto será una búsqueda en la página principal que recupere proyectos por evidencia real y explique su relevancia con un resumen breve y acotado, sin depender de un chatbot conversacional.

## 2. Problema

Los portfolios tradicionales suelen ser demasiado estáticos y no comunican bien:

- el contexto del problema
- la solución implementada
- las decisiones técnicas
- el uso de AI
- los resultados obtenidos
- la capacidad real de ejecución técnica

Además, cuando un potencial cliente busca experiencia en un dominio concreto como `SIEMENS`, `commissioning`, `SCADA` o `React`, el portfolio tradicional no responde bien a esa intención porque no conecta términos de búsqueda con evidencia concreta del trabajo realizado.

## 3. Objetivo del Producto

Construir un portfolio profesional e interactivo que:

- permita encontrar proyectos relevantes a partir de conceptos reales del cliente
- explique por qué cada proyecto es relevante sin inventar información
- muestre proyectos como case studies completos
- demuestre capacidad técnica real
- permita administración propia del contenido
- convierta visitas en oportunidades de contacto

## 4. Objetivos Específicos

- Publicar proyectos desde una consola admin propia.
- Permitir búsqueda de proyectos por evidencia real.
- Explicar la relevancia de cada resultado con resumen breve y contextual.
- Mostrar un catálogo navegable y filtrable.
- Mostrar detalle enriquecido por proyecto.
- Asociar tecnologías a cada proyecto.
- Gestionar media por proyecto.
- Capturar leads desde el sitio público.
- Mantener una base técnica escalable y propia.

## 5. No Objetivos

Fuera del alcance inicial:

- e-commerce
- checkout
- carrito
- pagos
- órdenes
- facturación
- blog complejo
- CMS genérico
- multirol avanzado
- chatbot conversacional libre

## 6. Usuarios

### Visitante

Quiere:

- entender rápidamente el perfil profesional
- escribir uno o varios conceptos y encontrar proyectos relevantes
- entender por qué cada proyecto coincide con su búsqueda
- ver detalle técnico y resultados
- contactar al autor

### Administrador

Quiere:

- crear y editar proyectos
- publicar o despublicar contenido
- asociar tecnologías y media
- revisar leads de contacto

## 7. Propuesta de Valor

PortfolioForge no será solo una vitrina visual. Será un portfolio con estructura de producto:

- búsqueda principal por evidencia real
- explicación resumida por resultado
- catálogo público
- detalle tipo case study
- administración propia
- datos estructurados
- capacidad de crecimiento futuro

## 8. Estado Actual del Proyecto

Ya existe una base funcional con:

- landing pública
- catálogo de proyectos
- detalle individual
- login admin
- listado admin de proyectos
- creación y edición
- publish/draft

Sin embargo, el modelo actual todavía arrastra naming transicional heredado de un catálogo de productos.

También existe una limitación funcional importante: hoy la búsqueda del catálogo solo filtra por nombre y categoría, lo cual no cubre la interacción principal deseada para el portfolio.

## 9. Criterio Documental

El PRD es el documento principal de referencia del proyecto.

Reglas de documentación:

- Toda la documentación del proyecto, excepto `README.md`, debe vivir dentro de `docs/`.
- No se deben crear documentos nuevos si repiten información ya cubierta por el PRD.
- Solo se separarán documentos cuando aporten claridad real y cubran un tema distinto.
- Los historiales de sesión solo se crearán cuando el usuario los pida explícitamente.

Reglas de lenguaje:

- usar `project` en lugar de `product`
- usar `lead` o `contact lead` en lugar de `order`
- usar `project media` en lugar de conceptos de catálogo comercial
- usar `case study` en lugar de ficha comercial

## 10. Visión Funcional

### Sitio Público

- Landing con buscador principal
- Catálogo de proyectos
- Resumen contextual por resultado
- Detalle de proyecto
- Contacto

### Dirección visual

- La landing debe usar una paleta oscura con estética orientada a developer.
- La referencia visual aportada por el usuario debe tomarse como guía de color y atmósfera, no como copia literal de layout.
- La prioridad visual está en contraste, legibilidad, foco en la barra de búsqueda y sensación técnica/profesional.

### Consola Admin

- Login
- CRUD de proyectos
- Gestión de publicación
- Gestión futura de tecnologías
- Gestión futura de media
- Gestión futura de leads

## 11. Requerimientos Funcionales

### RF1. Landing

El sistema debe mostrar una página principal con propuesta de valor, buscador principal de proyectos y accesos al catálogo y contacto.

### RF2. Búsqueda Principal

El sistema debe permitir que el visitante escriba uno o varios conceptos y recupere proyectos relevantes a partir de evidencia real del contenido del proyecto.

La búsqueda debe considerar al menos:

- `client_name`
- `technologies`
- `title`
- `summary`
- `description`
- `solution_summary`
- `architecture`
- `technical_decisions`
- `results`

La búsqueda debe tolerar variaciones razonables de escritura y priorizar coincidencias fuertes en cliente, tecnologías y contenido técnico.

La capa semántica debe diseñarse desde el inicio, pero sin acoplar el producto a un proveedor específico de embeddings en esta etapa.

### RF3. Explicación de Relevancia

El sistema debe mostrar debajo de cada resultado una frase breve que explique por qué el proyecto es relevante para la búsqueda actual.

Esta explicación debe:

- generarse solo a partir de evidencia del proyecto recuperado
- estar limitada a una frase breve
- evitar inferencias no sustentadas
- reflejar el concepto buscado por el usuario

### RF4. Catálogo

El sistema debe listar proyectos publicados y permitir filtrado complementario por criterios estructurados.

Los filtros estructurados iniciales deben incluir:

- categoría
- cliente
- tecnologías

### RF5. Detalle de Proyecto

El sistema debe mostrar para cada proyecto:

- título
- resumen
- descripción
- categoría
- contexto o cliente
- media principal

En fases siguientes también debe incluir:

- problema
- solución
- arquitectura
- integraciones
- uso de AI
- decisiones técnicas
- desafíos
- resultados
- métricas
- timeline
- links externos

### RF6. Admin de Proyectos

El administrador debe poder:

- crear proyectos
- editar proyectos
- publicar
- despublicar

Además, debe poder completar la información necesaria para que la búsqueda y la explicación de relevancia funcionen correctamente.

### RF7. Tecnologías

El sistema debe permitir asociar múltiples tecnologías a un proyecto.

### RF8. Media

El sistema debe permitir asociar múltiples assets a un proyecto:

- imágenes
- video
- recursos embebibles

### RF9. Leads

El sitio público debe permitir capturar contactos con:

- nombre
- email
- empresa
- mensaje
- proyecto de interés

### RF10. Gestión de Leads

El administrador debe poder listar y actualizar el estado de los leads.

## 12. Requerimientos No Funcionales

- Frontend SPA responsivo.
- Backend API claro y mantenible.
- Persistencia en PostgreSQL.
- Autenticación JWT para admin.
- Arquitectura preparada para crecer sin rehacer el proyecto.
- Separación clara entre frontend y backend.
- Las respuestas resumidas deben estar acotadas a evidencia real del proyecto.
- La capa AI no debe decidir resultados ni inventar contenido.
- La capa semántica utilizará OpenAI (text-embedding-3-small) como proveedor de embeddings, manteniendo una interfaz reemplazable en el código.

## 13. Modelo de Datos Objetivo

### `projects`

- `id`
- `title`
- `slug`
- `summary`
- `description`
- `category`
- `client_name`
- `hero_image`
- `featured`
- `status`
- `repo_url`
- `demo_url`
- `video_url`
- `started_at`
- `finished_at`
- `created_at`
- `updated_at`

### `technologies`

- `id`
- `name`
- `slug`
- `category`
- `icon`
- `color`

### `project_technologies`

- `project_id`
- `technology_id`

### `project_profiles`

- `project_id`
- `business_goal`
- `problem_statement`
- `solution_summary`
- `architecture`
- `integrations`
- `ai_usage`
- `technical_decisions`
- `challenges`
- `results`
- `metrics`
- `timeline`

### Campos derivados para búsqueda

- `search_document`
- `search_terms`
- `search_summary`
- `search_embedding`

Estos campos pueden materializarse en base de datos o generarse en capa de aplicación, pero deben existir conceptualmente para soportar recuperación y explicación de relevancia.

### `project_media`

- `id`
- `project_id`
- `type`
- `url`
- `caption`
- `sort_order`

### `contact_leads`

- `id`
- `name`
- `email`
- `company`
- `message`
- `project_interest`
- `status`
- `created_at`

## 14. Nota de Transición Técnica

La implementación actual todavía conserva nombres internos heredados en algunas capas del código. Esa deuda técnica se acepta temporalmente solo para acelerar el MVP, pero el dominio oficial del proyecto es el descrito en este PRD.

## 15. MVP Recomendado

El MVP debe incluir:

- landing con buscador principal
- recuperación de proyectos por evidencia real
- explicación resumida por resultado
- catálogo público con filtros estructurados mínimos
- detalle básico de proyecto
- login admin
- CRUD admin de proyectos
- publish/draft
- estructura inicial para technologies
- formulario de contacto
- base preparada para evolución a case study completo

Además, el MVP debe implementar la integración con OpenAI para embeddings y configurar la estrategia de fallback.

No incluye chatbot conversacional.

## 16. Fases

### Fase 1

Priorizar la interacción principal:

- limpiar terminología
- enfocar el producto como portfolio
- mover la búsqueda a la landing como interacción principal
- implementar recuperación por evidencia real
- implementar explicación resumida por resultado
- estabilizar catálogo y admin para soportar esta búsqueda

### Fase 2

Enriquecer el dominio:

- technologies
- project_profiles
- project_media
- mejorar ranking y cobertura de evidencia

### Fase 3

Agregar conversión:

- contact_leads
- gestión admin de leads

### Fase 4

Cerrar deuda técnica:

- renombrar completamente `product` -> `project` en backend, frontend y documentación

## 17. Métricas de Éxito

- cantidad de proyectos publicados
- porcentaje de proyectos con perfil completo
- porcentaje de búsquedas que devuelven resultados relevantes
- porcentaje de resultados con explicación válida basada en evidencia
- cantidad de leads recibidos
- tasa de conversión visita -> contacto
- tiempo de publicación de un nuevo proyecto

## 18. Riesgos

- mantener demasiado tiempo conceptos heredados del catálogo de productos
- crecer el alcance antes de cerrar el MVP
- mezclar features del proyecto original con el portfolio nuevo
- introducir una capa AI que invente o resuma sin sustento

## 19. Decisión de Producto

El proyecto debe evolucionar desde la base actual sin rehacerse desde cero, pero eliminando progresivamente toda terminología y documentación heredada que no pertenezca al dominio del portfolio.

La prioridad de implementación es la búsqueda de proyectos por evidencia real y la explicación breve de relevancia por resultado. No habrá chatbot conversacional; la AI se usará solo como capa acotada de resumen sobre resultados ya recuperados por el sistema.
