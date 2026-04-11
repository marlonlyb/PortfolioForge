import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { httpGet } from '../../shared/api/http';
import { getToken, setToken, removeToken } from '../../shared/storage';

// ─── Types ────────────────────────────────────────────────────────────

export interface SessionUser {
  id: string;
  email: string;
  is_admin: boolean;
  created_at: string;
}

interface LoginResponse {
  user: SessionUser;
  token: string;
  expires_in: number;
}

interface SessionState {
  user: SessionUser | null;
  token: string | null;
  loading: boolean;
}

interface SessionContextValue extends SessionState {
  login: (response: LoginResponse) => void;
  logout: () => void;
}

// ─── Context ──────────────────────────────────────────────────────────

const SessionContext = createContext<SessionContextValue | null>(null);

// ─── Provider ─────────────────────────────────────────────────────────

export function SessionProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<SessionState>(() => {
    const existingToken = getToken();
    return {
      user: null,
      token: existingToken,
      loading: existingToken !== null,
    };
  });

  // Restore session from /private/me when a token exists at boot
  useEffect(() => {
    if (!state.token) return;

    let cancelled = false;

    httpGet<SessionUser>('/api/v1/private/me')
      .then((user) => {
        if (!cancelled) {
          setState({ user, token: state.token, loading: false });
        }
      })
      .catch(() => {
        // Token is invalid or expired — clear it
        if (!cancelled) {
          removeToken();
          setState({ user: null, token: null, loading: false });
        }
      });

    return () => {
      cancelled = true;
    };
  }, [state.token]);

  const login = (response: LoginResponse) => {
    setToken(response.token);
    setState({ user: response.user, token: response.token, loading: false });
  };

  const logout = () => {
    removeToken();
    setState({ user: null, token: null, loading: false });
  };

  return (
    <SessionContext.Provider value={{ ...state, login, logout }}>
      {children}
    </SessionContext.Provider>
  );
}

// ─── Hook ─────────────────────────────────────────────────────────────

export function useSession(): SessionContextValue {
  const ctx = useContext(SessionContext);
  if (!ctx) {
    throw new Error('useSession must be used within a SessionProvider');
  }
  return ctx;
}
