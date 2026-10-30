import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from 'react';
import { AppState } from 'react-native';
import { apiRequest, setSessionToken } from '@/lib/api';
import { handleSessionLost } from '@/lib/session';
import { User } from '@/lib/types';

interface AuthState {
  authenticated: boolean;
  user: User | null;
}

interface AuthContextValue {
  auth: AuthState;
  loading: boolean;
  refreshAuth: () => Promise<boolean>;
  setAuth: (auth: AuthState) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue>({
  auth: { authenticated: false, user: null },
  loading: true,
  refreshAuth: async () => false,
  setAuth: () => {},
  logout: () => {},
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [auth, setAuthState] = useState<AuthState>({ authenticated: false, user: null });
  const [loading, setLoading] = useState(true);
  const authRef = useRef(auth);

  const setAuth = useCallback((next: AuthState) => {
    authRef.current = next;
    setAuthState(next);
  }, []);

  const refreshAuth = useCallback(async (): Promise<boolean> => {
    try {
      const data = await apiRequest<{ status: string; user: User }>('/me');
      if (data.status !== 'success' || !data.user) {
        const error = new Error('unauthorized') as Error & { status?: number };
        error.status = 401;
        throw error;
      }
      setAuth({ authenticated: true, user: data.user });
      return true;
    } catch (error: unknown) {
      const status = (error as { status?: number } | null)?.status;
      if (status === 401 || status === 403) {
        setAuth({ authenticated: false, user: null });
        await setSessionToken(null);
        return false;
      }
      return authRef.current.authenticated;
    }
  }, [setAuth]);

  // Restore the session once on mount.
  useEffect(() => {
    let cancelled = false;
    refreshAuth().finally(() => {
      if (!cancelled) setLoading(false);
    });
    return () => {
      cancelled = true;
    };
  }, [refreshAuth]);

  // 10s session revalidation loop, like the web frontend.
  useEffect(() => {
    if (!auth.authenticated) return;
    const interval = setInterval(async () => {
      if (AppState.currentState !== 'active') return;
      const ok = await refreshAuth();
      if (!ok) handleSessionLost();
    }, 10000);
    return () => clearInterval(interval);
  }, [auth.authenticated, refreshAuth]);

  const logout = useCallback(() => {
    setAuth({ authenticated: false, user: null });
    void setSessionToken(null);
  }, [setAuth]);

  const value = useMemo(
    () => ({ auth, loading, refreshAuth, setAuth, logout }),
    [auth, loading, refreshAuth, setAuth, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  return useContext(AuthContext);
}