# Frontend (`client/`)

Frontend actual de PortfolioForge usando React 19 + Vite + TypeScript.

## Requisitos

- Node.js 20+
- Backend corriendo en `http://localhost:8080`
- `ALLOWED_ORIGINS=http://localhost:5173` en el backend para desarrollo Vite

## Variables de entorno

Copia `.env.example` a `.env` y define:

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

- `VITE_API_BASE_URL`: base path del backend actual.

## Scripts

```bash
npm install
npm run dev
```

La app Vite abre por defecto en `http://localhost:5173`.

## Verificación mínima

- `/` renderiza el shell público.
- `/login` renders the only public auth entry point: Google plus passwordless email OTP.
- `/admin/login` sigue disponible por URL directa pero no se expone en la navegación pública.
- `/admin` sigue accesible bajo `RequireAdmin` con su layout actual.
