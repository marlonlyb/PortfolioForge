# Frontend MVP (`client/`)

Base SPA for the store MVP using React 19 + Vite + TypeScript.

## Requisitos

- Node.js 20+
- Backend corriendo en `http://localhost:8080`
- `ALLOWED_ORIGINS=http://localhost:5173` en el backend para desarrollo Vite

## Variables de entorno

Copia `.env.example` a `.env` y define:

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_PAYPAL_CLIENT_ID=your_paypal_sandbox_client_id
```

- `VITE_API_BASE_URL`: base path del backend store.
- `VITE_PAYPAL_CLIENT_ID`: client id público de PayPal Sandbox para el futuro checkout.

## Scripts

```bash
npm install
npm run dev
```

La app Vite abre por defecto en `http://localhost:5173`.

## Verificación mínima de Phase 1

- `/` renderiza el shell público.
- `/login` renderiza placeholder de autenticación.
- `/checkout` existe como ruta protegida base y redirige a `/login` hasta conectar sesión real en Phase 2.
