# Hostinger VPS deployment assets

## Runtime layout

```text
/srv/portfolioforge/
├── source/current                     # checkout used to build releases
├── frontend/
│   ├── releases/<timestamp>/          # published Vite dist/
│   └── current -> releases/<timestamp>
└── api/
    ├── releases/<timestamp>/
    │   └── portfolioforge-api         # Go binary built from ./cmd
    └── current -> releases/<timestamp>

/etc/portfolioforge/backend.env        # backend secrets/runtime env
/var/lib/portfolioforge/images         # persistent IMAGES_DIR
```

## Repo templates

- `nginx/portfolioforge.conf.example` → `/etc/nginx/sites-available/portfolioforge.conf`
- `systemd/portfolioforge-api.service.example` → `/etc/systemd/system/portfolioforge-api.service`

## Command map

### Deploy

```bash
timestamp="$(date +%Y%m%d%H%M%S)"
frontend_release="/srv/portfolioforge/frontend/releases/$timestamp"
api_release="/srv/portfolioforge/api/releases/$timestamp"

cd /srv/portfolioforge/source/current/client
npm ci
VITE_API_BASE_URL= npm run build

cd /srv/portfolioforge/source/current
mkdir -p "$frontend_release" "$api_release"
cp -R client/dist/. "$frontend_release/"
go build -o "$api_release/portfolioforge-api" ./cmd
ln -sfn "$frontend_release" /srv/portfolioforge/frontend/current
ln -sfn "$api_release" /srv/portfolioforge/api/current
sudo systemctl restart portfolioforge-api
sudo systemctl reload nginx
```

### Smoke

```bash
sudo nginx -t
sudo systemd-analyze verify /etc/systemd/system/portfolioforge-api.service
sudo systemctl status portfolioforge-api --no-pager
curl -fsS https://portfolioforge.mlbautomation.com/ >/dev/null
curl -fsS https://portfolioforge.mlbautomation.com/health
curl -fsS https://portfolioforge.mlbautomation.com/api/v1/projects >/dev/null
```

### Rollback

```bash
sudo -u postgres pg_dump -Fc portfolioforge > "/srv/portfolioforge/backups/pre-rollback-$(date +%Y%m%d%H%M%S).dump"
ln -sfn /srv/portfolioforge/frontend/releases/<previous> /srv/portfolioforge/frontend/current
ln -sfn /srv/portfolioforge/api/releases/<previous> /srv/portfolioforge/api/current
sudo install -m 640 -o root -g portfolioforge <previous-backend.env> /etc/portfolioforge/backend.env
sudo systemctl restart portfolioforge-api
sudo systemctl reload nginx
```

Consulta `docs/operacion/10-hostinger-vps-deployment.md` para el runbook completo, bootstrap del VPS y checklist de validación.
