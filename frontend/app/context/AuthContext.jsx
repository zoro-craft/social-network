'use client';

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [auth, setAuth] = useState({ authenticated: false, nickname: '' });
  const [loading, setLoading] = useState(true);

  const refreshAuth = useCallback(async () => {
    try {
      const res = await fetch('http://localhost:8080/me', {
        credentials: 'include',
      });
      if (!res.ok) throw new Error();
      const data = await res.json();
      setAuth(data);
    } catch {
      setAuth({ authenticated: false, nickname: '' });
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const res = await fetch('http://localhost:8080/me', {
          credentials: 'include',
        });
        if (!res.ok) throw new Error();
        const data = await res.json();
        if (!cancelled) setAuth(data);
      } catch {
        if (!cancelled) setAuth({ authenticated: false, nickname: '' });
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <AuthContext.Provider value={{ auth, setAuth, loading, refreshAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
