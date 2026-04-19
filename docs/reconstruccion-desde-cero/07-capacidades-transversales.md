# Capacidades transversales

## Búsqueda

### Objetivo

Permitir discovery por evidencia real, no solo por navegación manual.

### Estrategia mínima

- FTS ponderado sobre nombre, cliente/contexto, summary técnico y tecnologías;
- fuzzy matching con `pg_trgm`;
- embeddings semánticos con `pgvector`;
- explicación breve por resultado para justificar la coincidencia.

### Datos indexables mínimos

- `name`
- `client_name`/`brand`
- `description`
- `solution_summary`
- `architecture`
- `business_goal`
- `problem_statement`
- `ai_usage`
- nombres de tecnologías

## Assistant

### Regla de grounding

Solo debe responder con contexto del markdown remoto del proyecto seleccionado.

### Requisitos de acceso

- proyecto activo;
- `source_markdown_url` presente;
- usuario autenticado;
- usuario con `can_use_project_assistant=true`.

### Capacidades técnicas necesarias

- fetch/cache de markdown remoto;
- validación de slug y acceso;
- prompt grounded por proyecto;
- historial corto de conversación por request.

## Localización

### Idiomas públicos

- `es` base
- `ca`
- `en`
- `de`

### Regla de derivación

Todo contenido parte de `es`. Las demás locales se generan como derivadas persistidas y pueden admitir override manual por campo.

### Campos traducibles mínimos

- `name`, `description`, `category`, `client_name`
- todos los campos de `ProjectProfile`

## Auth y sesión

### Capacidades mínimas

- signup/login local;
- login Google;
- verificación de email por OTP;
- perfil completo;
- separación entre usuario público y admin.

### Objetivo funcional

El auth no existe solo para proteger admin; también habilita el assistant por proyecto.

## Admin y operación

### Módulos mínimos

- proyectos;
- tecnologías;
- usuarios;
- site settings;
- workflow de case study.

### Workflow operativo mínimo

Pasos persistidos:

1. `resolve_source`
2. `publish_canonical`
3. `import_or_update_project`
4. `localization_backfill`
5. `reembed`

Los pasos 2 y 3 requieren confirmación explícita según la implementación actual.

## Seguridad y sanitización

- no exponer `source_markdown_url` en API pública;
- no publicar secretos, datos personales ni material contractual en canonical/runtime;
- normalizar correo de user para unicidad;
- proteger rutas admin y privadas por middleware de autenticación/autorización.
