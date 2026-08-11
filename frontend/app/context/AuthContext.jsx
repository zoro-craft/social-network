'use client';

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from 'react';

import { API_BASE } from '../lib/api';

const AuthContext = createContext(null);

const defaultAuth = {
  authenticated: false,
  nickname: '',
};

export function AuthProvider({ children }) {
  const [auth, setAuth] = useState(defaultAuth);
  const [loading, setLoading] = useState(true);

  const refreshAuth = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/me`, {
        credentials: 'include',
      });

      if (!res.ok) {
        throw new Error();
      }

      const data = await res.json();
      setAuth(data);
    } catch {
      setAuth(defaultAuth);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function checkAuth() {
      try {
        const res = await fetch(`${API_BASE}/me`, {
          credentials: 'include',
        });

        if (!res.ok) {
          throw new Error();
        }

        const data = await res.json();

        if (!cancelled) {
          setAuth(data);
        }
      } catch {
        if (!cancelled) {
          setAuth(defaultAuth);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    checkAuth();

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <AuthContext.Provider
      value={{
        auth,
        setAuth,
        loading,
        refreshAuth,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
