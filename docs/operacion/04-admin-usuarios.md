# Primer admin y operaciones de usuarios

## Objetivo

Bootstrapear el primer administrador y operar el módulo `/admin/users` sin tocar código.

## Hechos del sistema

- no existe endpoint público para crear un admin directamente;
- el login admin usa `POST /api/v1/admin/login`;
- una vez existe un admin, la UI/API admin puede promover o degradar usuarios estándar usando solo `is_admin`;
- los usuarios que ya son admin están protegidos contra edición desde este flujo;
- un admin no puede borrarse a sí mismo y los admins no se eliminan desde la UI.

## Runbook A — bootstrap del primer admin

### Paso 1: crear un usuario local normal

Con backend levantado:

```bash
curl -X POST http://localhost:8080/api/v1/public/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"ChangeMe123!"}'
```

Esto crea un usuario `local` estándar.

### Paso 2: promoverlo en PostgreSQL

Como todavía no existe admin, la promoción inicial debe hacerse por SQL:

```bash
psql -d portfolioforge -c "UPDATE users SET is_admin = TRUE, updated_at = EXTRACT(EPOCH FROM now())::int WHERE LOWER(email) = LOWER('admin@example.com') AND deleted_at IS NULL;"
```

Verificación:

```bash
psql -d portfolioforge -c "SELECT email, is_admin, auth_provider, email_verified, deleted_at FROM users WHERE LOWER(email) = LOWER('admin@example.com');"
```

### Paso 3: iniciar sesión como admin

UI:

- `http://localhost:5173/admin/login`

API:

```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"ChangeMe123!"}'
```

## Runbook B — promoción/degradación posterior por UI/API

Una vez existe al menos un admin:

- UI: `/admin/users`
- listar: `GET /api/v1/admin/users`
- detalle: `GET /api/v1/admin/users/:id`
- cambiar rol: `PATCH /api/v1/admin/users/:id`
- soft delete: `DELETE /api/v1/admin/users/:id`

Payload admitido para update:

```json
{ "is_admin": true }
```

No se admite editar email, perfil ni identidad desde este flujo.

## Runbook C — operaciones de usuario público relacionadas

### Perfil obligatorio para assistant público

Ruta:

- `PUT /api/v1/private/me/profile`

Payload:

```json
{ "full_name": "Nombre Apellido", "company": "Mi empresa" }
```

### OTP y verificación

Rutas:

- `POST /api/v1/public/email-verification/request`
- `POST /api/v1/public/email-verification/resend`
- `POST /api/v1/public/email-verification/verify`

## Reglas operativas importantes

- el primer admin se bootstrapea fuera de la UI, por SQL;
- después de eso, usa `/admin/users` como flujo canónico;
- no intentes editar admins existentes desde `/admin/users/:id`: el backend lo bloquea por diseño;
- si solo necesitas un admin inicial para operar el sistema, no hace falta marcar `email_verified=true` para login admin local.

## Follow-up/riesgo documentado

- el repo no ofrece todavía un comando o migración dedicada para bootstrapear el primer admin; hoy la promoción inicial depende de SQL manual.
