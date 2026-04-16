import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { httpGet } from '../../shared/api/http';
import { getToken, setToken, removeToken } from '../../shared/storage';

// ─── Types ────────────────────────────────────────────────────────────

export interface SessionUser {
  id: string;
  email: string;
  is_admin: boolean;
  auth_provider: 'local' | 'google';
  email_verified: boolean;
  full_name?: string;
  company?: string;
  profile_completed: boolean;
  assistant_eligible: boolean;
  can_use_project_assistant: boolean;
  created_at: string;
  last_login_at?: string;
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
  refreshSession: () => Promise<SessionUser | null>;
  setUser: (user: SessionUser | null) => void;
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
  async function refreshSession(): Promise<SessionUser | null> {
    const activeToken = state.token ?? getToken();
    if (!activeToken) {
      setState((current) => ({ ...current, user: null, loading: false }));
      return null;
    }

    try {
      const user = await httpGet<SessionUser>('/api/v1/private/me');
      setState((current) => ({ ...current, user, loading: false }));
      return user;
    } catch {
      removeToken();
      setState({ user: null, token: null, loading: false });
      return null;
    }
  }

  useEffect(() => {
    if (!state.token) return;

    void refreshSession();
  }, [state.token]);

  const login = (response: LoginResponse) => {
    setToken(response.token);
    setState({ user: response.user, token: response.token, loading: false });
  };

  const setUser = (user: SessionUser | null) => {
    setState((current) => ({ ...current, user }));
  };

  const logout = () => {
    removeToken();
    setState({ user: null, token: null, loading: false });
  };

  return (
    <SessionContext.Provider value={{ ...state, login, refreshSession, setUser, logout }}>
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
