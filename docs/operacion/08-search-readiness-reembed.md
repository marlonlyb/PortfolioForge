# Search readiness, refresh y re-embed

## Objetivo

Mantener consistente el documento de búsqueda, su readiness y, cuando aplica, el embedding semántico.

## 1. Qué mide readiness

El backend clasifica un proyecto según estos campos:

- `name`
- `description` (mínimo útil de 20 caracteres)
- `category`
- al menos una tecnología asociada
- `profile.solution_summary`

Niveles:

- `incomplete`: falta `name` o `description`
- `basic`: tiene base pero le faltan perfil/tecnologías clave
- `complete`: tiene base + `solution_summary` + tecnologías

## 2. Endpoint de readiness

```http
GET /api/v1/admin/projects/:id/readiness
```

Úsalo antes de re-embed para distinguir:

- problema de contenido faltante;
- problema de refresh/indexación.

## 3. Refresh de un proyecto

### API

```http
POST /api/v1/admin/projects/:id/reembed
```

### Efecto real

1. recompone el documento de búsqueda (`project_search_documents`);
2. si el contenido cambió y `ENABLE_SEMANTIC_SEARCH=true`, vuelve a generar embedding;
3. si no hubo cambio de contenido, evita trabajo innecesario.

## 4. Refresh masivo

### API

```http
POST /api/v1/admin/projects/reembed-stale
```

Efecto:

- recorre proyectos `active = TRUE`;
- refresca documento/embedding uno por uno.

## 5. Cuándo usar cada uno

### Reembed de un proyecto

Después de:

- editar enrichment;
- cambiar tecnologías;
- importar un case study concreto;
- corregir texto base de un proyecto puntual.

### Reembed masivo

Después de:

- una migración que altere composición de search;
- una corrección global de datos;
- un lote de imports o localizaciones.

## 6. Dependencias semánticas

Para embeddings reales necesitas:

- `ENABLE_SEMANTIC_SEARCH=true`
- `OPENAI_API_KEY`

Sin eso:

- el refresh documental sigue existiendo;
- la parte semántica queda degradada o sin actualizar.

## 7. Secuencia recomendada

1. verificar readiness;
2. corregir contenido faltante si el nivel no es suficiente;
3. ejecutar reembed del proyecto;
4. validar búsqueda pública con queries reales;
5. si hubo cambios amplios, ejecutar refresh masivo.

## 8. Verificación mínima

- `GET /api/v1/admin/projects/:id/readiness`
- `GET /api/v1/public/search?q=<query>`
- revisar que el proyecto aparezca y que la explicación/resultados sean coherentes.
