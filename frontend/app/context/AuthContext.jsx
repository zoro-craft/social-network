'use client';

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from 'react';
import { usePathname } from 'next/navigation';
import { API_BASE } from '../lib/api';

const AuthContext = createContext(null);

async function fetchAuth() {
  try {
    const res = await fetch(`${API_BASE}/me`, { credentials: 'include' });
    if (!res.ok) throw new Error();
    return await res.json();
  } catch {
    return { authenticated: false, nickname: '' };
  }
}

export function AuthProvider({ children }) {
  const [auth, setAuth] = useState({ authenticated: false, nickname: '' });
  const [loading, setLoading] = useState(true);
  const pathname = usePathname();

  const refreshAuth = useCallback(async () => {
    const data = await fetchAuth();
    setAuth(data);
    setLoading(false);
  }, []);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    refreshAuth();
  }, [pathname, refreshAuth]);

  return (
    <AuthContext.Provider value={{ auth, setAuth, loading, refreshAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
