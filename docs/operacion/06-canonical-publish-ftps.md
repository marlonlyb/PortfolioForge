# Runbook de publicación canónica por FTPS

## Uso

Publica por FTPS toda la carpeta del slug canonical y verifica la URL pública final.

## Comando

Desde la raíz del repo:

```bash
go run ./cmd canonical-publish --case-dir "<ruta-del-caso>"
```

Opciones soportadas:

- `--slug <slug>`
- `--dry-run`

Ayuda implícita del comando:

- resuelve un canonical dentro de `90. dev_portfolioforge/<slug>/`;
- publica la carpeta completa del slug, no solo el `.md`.

## Variables requeridas

- `PF_FTP_HOST`
- `PF_FTP_USER`
- `PF_FTP_PASSWORD`
- `PF_PUBLIC_BASE`

Opcionales con default:

- `PF_FTP_PORT` → `21`
- `PF_FTP_REMOTE_BASE` → `/`

## Casos de entrada válidos

`--case-dir` puede apuntar a:

- el directorio del slug;
- la carpeta `90. dev_portfolioforge/`;
- una raíz superior que contenga `90. dev_portfolioforge/`.

Si hay múltiples slugs válidos, debes pasar `--slug`.

## Dry run recomendado

Antes de subir:

```bash
go run ./cmd canonical-publish --case-dir "<ruta-del-caso>" --slug "<slug>" --dry-run
```

Úsalo para validar:

- slug resuelto;
- directorio local;
- directorio remoto;
- URL pública final.

## Criterio de éxito

El comando debe imprimir:

- `slug=...`
- `local=...`
- `remote=...`
- `url=...`
- `publicado OK ...`

Además la URL final debe responder por HTTPS.

## Relación con el workflow admin

- este comando es el equivalente CLI del step `publish_canonical`;
- si ya estás usando el workflow persistido de `/admin/settings/case-studies`, no necesitas publicar aparte salvo diagnóstico o ejecución manual controlada.

## Fallos esperables

- faltan variables `PF_*`;
- el slug no existe o hay múltiples slugs sin `--slug`;
- el markdown `<slug>.md` no existe;
- error de login FTPS;
- error TLS/timeout;
- la verificación HTTPS final falla.

## Regla editorial asociada

- publicar solo cuando el canonical local en `90. dev_portfolioforge/<slug>/<slug>.md` ya es la versión correcta en castellano.
