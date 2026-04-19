# Superficies de frontend y aplicación

## Principio

El frontend se divide en experiencia pública y consola admin. Ambas consumen el mismo backend, pero con contratos y objetivos distintos.

## Superficie pública

### Landing `/`

- entrada principal al portfolio;
- presenta branding global y acceso a búsqueda/catálogo;
- debe dirigir al detalle de proyectos.

### Detalle `/projects/:slug`

- vista principal del proyecto;
- consume detalle público enriquecido;
- muestra capas estrategia, ejecución y técnica;
- renderiza tecnologías, media y assistant si está habilitado.

### Search `/search`

- búsqueda híbrida con barra, filtros y listado de resultados;
- permite explorar por evidencia real, no por tags decorativos;
- muestra explicación corta y tecnologías relevantes.

### Auth pública

- `/login`
- `/signup`
- `/verify-email`
- `/complete-profile`

Estas pantallas habilitan la elegibilidad para assistant y persistencia de sesión.

## Superficie admin

### Layout `/admin`

Debe agrupar módulos protegidos por rol admin.

### Proyectos

- `/admin/projects`
- `/admin/projects/new`
- `/admin/projects/:id`

Capacidades:

- CRUD base;
- edición de enrichment;
- asociación de tecnologías;
- carga/edición de media;
- edición de `source_markdown_url`;
- control de publicación/estado.

### Tecnologías

- `/admin/technologies`
- `/admin/technologies/new`
- `/admin/technologies/:id`

### Usuarios

- `/admin/users`
- `/admin/users/:id`

### Ajustes

- `/admin/settings/case-studies`
- `/admin/settings` → redirección a `/admin/settings/case-studies`

## Módulos de UI recomendados

- `app/` para router, layouts y providers;
- `features/catalog` para catálogo/detalle/assistant;
- `features/search` para resultados y filtros;
- `features/auth` para login/signup/OTP/perfil;
- `features/admin-projects` para CRUD canónico de proyectos;
- `features/admin-technologies`, `admin-users`, `admin-settings` para backoffice;
- `shared/i18n`, `shared/types`, `shared/api`, `shared/lib` para infraestructura común.

## Reglas de frontend

- usar naming `project` aunque existan componentes legacy `product` durante migración;
- modelar `media` como contrato principal y `images` como fallback;
- la UI pública debe ser más breve que el markdown canónico;
- assistant solo aparece si backend devuelve `assistant_available=true` y la sesión lo permite;
- locale público soportado: `es`, `ca`, `en`, `de`.

## Estado mínimo que debe manejar el cliente

- sesión autenticada y claims de elegibilidad;
- locale actual;
- filtros de búsqueda;
- estado de formularios admin;
- errores de API y feedback de operaciones.
