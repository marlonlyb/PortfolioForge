# Runbook de publicación canónica por FTPS

> Estado: **legacy / opcional / compatibilidad**.

Este documento no describe el flujo estándar principal de PortfolioForge. El flujo estándar actual es:

1. crear o actualizar el canonical en la carpeta del caso;
2. publicar manualmente la URL final fuera de PortfolioForge;
3. una vez exista esa URL, continuar con runtime/UI usando `source_markdown_url`.

Usa `canonical-publish` solo si necesitas mantener una ruta legacy, hacer diagnóstico o aprovechar compatibilidad con entornos existentes.

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
- ese step y este comando siguen existiendo, pero no forman parte del camino estándar recomendado para crear o cerrar un case study;
- si ya tienes la URL remota publicada manualmente por el host externo del usuario, no necesitas este comando.

## Fallos esperables

- faltan variables `PF_*`;
- el slug no existe o hay múltiples slugs sin `--slug`;
- el markdown `<slug>.md` no existe;
- error de login FTPS;
- error TLS/timeout;
- la verificación HTTPS final falla.

## Regla editorial asociada

- publicar solo cuando el canonical local en `90. dev_portfolioforge/<slug>/<slug>.md` ya es la versión correcta en castellano;
- si se usa esta ruta legacy, la salida esperada sigue siendo la misma: obtener la URL final que luego se guardará como `source_markdown_url` para runtime/UI.
