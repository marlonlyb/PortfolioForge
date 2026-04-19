# Despliegue en Hostinger VPS

## Objetivo

Dejar PortfolioForge publicado en un solo VPS Ubuntu 24.04 con Nginx + Let's Encrypt al frente, frontend Vite servido estáticamente, API Go bajo `systemd` y PostgreSQL 16 local.

## Topología fija

- dominio público: `https://portfolioforge.mlbautomation.com`
- IP del VPS: `187.127.73.230`
- frontend: `/srv/portfolioforge/frontend/current`
- backend: `/srv/portfolioforge/api/current/portfolioforge-api`
- `IMAGES_DIR`: `/var/lib/portfolioforge/images`
- env del backend: `/etc/portfolioforge/backend.env`

## 1. DNS

1. crear/confirmar `A record` para `portfolioforge.mlbautomation.com -> 187.127.73.230`
2. esperar propagación antes de emitir el certificado TLS

## 2. Paquetes base del VPS

```bash
sudo apt-get update
sudo apt-get install -y nginx certbot python3-certbot-nginx curl git unzip
sudo apt-get install -y golang-go nodejs npm
sudo apt-get install -y postgresql-16 postgresql-client-16 postgresql-16-pgvector
```

Si el VPS no trae Node 20+ o Go suficientemente nuevo desde `apt`, instala la versión soportada por el repositorio antes del primer build.

## 3. Usuario, rutas y ownership

```bash
sudo groupadd --system portfolioforge || true
sudo useradd --system --gid portfolioforge --home /srv/portfolioforge --shell /usr/sbin/nologin portfolioforge || true

sudo install -d -m 755 -o portfolioforge -g portfolioforge /srv/portfolioforge
sudo install -d -m 755 -o portfolioforge -g portfolioforge /srv/portfolioforge/source
sudo install -d -m 755 -o portfolioforge -g portfolioforge /srv/portfolioforge/frontend/releases
sudo install -d -m 755 -o portfolioforge -g portfolioforge /srv/portfolioforge/api/releases
sudo install -d -m 750 -o portfolioforge -g portfolioforge /var/lib/portfolioforge/images
sudo install -d -m 750 -o root -g portfolioforge /etc/portfolioforge
sudo install -d -m 755 -o root -g root /var/www/certbot
```

## 4. Publicar el código fuente en el VPS

```bash
sudo -u portfolioforge git clone <repo-url> /srv/portfolioforge/source/current
```

En releases siguientes puedes actualizar el checkout con `git fetch --all --tags && git checkout <ref>` dentro de `/srv/portfolioforge/source/current`.

## 5. Variables de entorno del backend

Crear `/etc/portfolioforge/backend.env` con modo `640`, owner `root` y group `portfolioforge`.

Valores mínimos para este despliegue:

```dotenv
SERVER_PORT=8080
ALLOWED_ORIGINS=https://portfolioforge.mlbautomation.com
ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
IMAGES_DIR=/var/lib/portfolioforge/images
JWT_SECRET_KEY=<secret>
DB_USER=portfolioforge
DB_PASSWORD=<secret>
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=portfolioforge
DB_SSL_MODE=disable
VITE_API_BASE_URL=
```

`VITE_API_BASE_URL=` debe quedar vacío al construir el frontend para que el bundle use rutas relativas del mismo origen (`/api/v1/...`).

## 6. Bootstrap de PostgreSQL y migraciones

Seguir el runbook detallado de [`03-db-bootstrap-migraciones.md`](./03-db-bootstrap-migraciones.md). Resumen mínimo:

```bash
sudo -u postgres createuser --login --pwprompt portfolioforge
sudo -u postgres createdb --owner=portfolioforge portfolioforge
sudo -u postgres psql -d portfolioforge -c "CREATE EXTENSION IF NOT EXISTS unaccent;"
sudo -u postgres psql -d portfolioforge -c "CREATE EXTENSION IF NOT EXISTS pg_trgm;"
sudo -u postgres psql -d portfolioforge -c "CREATE EXTENSION IF NOT EXISTS vector;"

sudo -u portfolioforge bash -lc 'cd /srv/portfolioforge/source/current && for f in sqlmigrations/*.sql; do psql -v ON_ERROR_STOP=1 -h 127.0.0.1 -U portfolioforge -d portfolioforge -f "$f"; done'
```

Antes de poner tráfico público o aplicar nuevas migraciones:

```bash
sudo install -d -m 750 -o portfolioforge -g portfolioforge /srv/portfolioforge/backups
sudo -u postgres pg_dump -Fc portfolioforge > "/srv/portfolioforge/backups/pre-cutover-$(date +%Y%m%d%H%M%S).dump"
```

## 7. Build y publicación de releases

Desde `/srv/portfolioforge/source/current`:

```bash
timestamp="$(date +%Y%m%d%H%M%S)"
frontend_release="/srv/portfolioforge/frontend/releases/$timestamp"
api_release="/srv/portfolioforge/api/releases/$timestamp"

cd /srv/portfolioforge/source/current/client
npm ci
VITE_API_BASE_URL= npm run build

mkdir -p "$frontend_release"
cp -R dist/. "$frontend_release/"

cd /srv/portfolioforge/source/current
mkdir -p "$api_release"
go build -o "$api_release/portfolioforge-api" ./cmd

ln -sfn "$frontend_release" /srv/portfolioforge/frontend/current
ln -sfn "$api_release" /srv/portfolioforge/api/current
```

## 8. Instalar `systemd` y Nginx

```bash
sudo cp deploy/hostinger/systemd/portfolioforge-api.service.example /etc/systemd/system/portfolioforge-api.service
sudo systemctl daemon-reload
sudo systemctl enable --now portfolioforge-api
sudo certbot certonly --standalone -d portfolioforge.mlbautomation.com
sudo cp deploy/hostinger/nginx/portfolioforge.conf.example /etc/nginx/sites-available/portfolioforge.conf
sudo ln -sfn /etc/nginx/sites-available/portfolioforge.conf /etc/nginx/sites-enabled/portfolioforge.conf
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

Validaciones inmediatas:

```bash
sudo systemd-analyze verify /etc/systemd/system/portfolioforge-api.service
sudo systemctl status portfolioforge-api --no-pager
sudo journalctl -u portfolioforge-api -n 100 --no-pager
sudo tail -n 100 /var/log/nginx/error.log
```

## 9. Smoke checks obligatorios

1. `curl -I https://portfolioforge.mlbautomation.com/` → `200` y certificado válido.
2. abrir `https://portfolioforge.mlbautomation.com/admin/projects` y recargar → debe volver `index.html` sin 404.
3. `curl -fsS https://portfolioforge.mlbautomation.com/health`
4. `curl -fsS https://portfolioforge.mlbautomation.com/api/v1/projects >/dev/null`
5. reiniciar el VPS o al menos ejecutar `sudo systemctl restart portfolioforge-api && sudo systemctl reload nginx` y repetir checks.

## 10. Logs y diagnóstico rápido

- `sudo journalctl -u portfolioforge-api -f`
- `sudo systemctl status portfolioforge-api --no-pager`
- `sudo tail -f /var/log/nginx/access.log /var/log/nginx/error.log`
- `sudo nginx -t`

## 11. Rollback exacto

1. generar backup antes de tocar symlinks si el incidente ocurre después de migraciones o writes recientes:

   ```bash
   sudo -u postgres pg_dump -Fc portfolioforge > "/srv/portfolioforge/backups/rollback-$(date +%Y%m%d%H%M%S).dump"
   ```

2. identificar los releases previos válidos:

   ```bash
   ls -1 /srv/portfolioforge/frontend/releases
   ls -1 /srv/portfolioforge/api/releases
   ```

3. restaurar symlinks:

   ```bash
   sudo ln -sfn /srv/portfolioforge/frontend/releases/<previous-frontend-ts> /srv/portfolioforge/frontend/current
   sudo ln -sfn /srv/portfolioforge/api/releases/<previous-api-ts> /srv/portfolioforge/api/current
   ```

4. restaurar snapshot previo de `/etc/portfolioforge/backend.env` si el problema fue de configuración:

   ```bash
   sudo install -m 640 -o root -g portfolioforge <saved-backend.env> /etc/portfolioforge/backend.env
   ```

5. reiniciar en este orden:

   ```bash
   sudo systemctl restart portfolioforge-api
   sudo nginx -t
   sudo systemctl reload nginx
   ```

6. si la falla provino de migraciones incompatibles o corrupción operativa, restaurar el dump más reciente de `pg_dump` antes de reabrir tráfico.

## 12. Criterio de aceptación

- Nginx válido (`nginx -t`)
- unit file válido (`systemd-analyze verify`)
- API viva tras restart/reboot
- `/`, `/admin/projects`, `/health` y un endpoint DB-backed de `/api/v1/...` responden correctamente
- `IMAGES_DIR`, DB y releases persisten después del reinicio
