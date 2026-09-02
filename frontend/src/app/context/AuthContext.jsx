'use client';

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from 'react';
import { usePathname } from 'next/navigation';
import { API_BASE } from '@/utils/api';

const AuthContext = createContext({
  auth: { authenticated: false, nickname: '' },
  loading: true,
  refreshAuth: async () => {},
});

async function fetchAuth() {
  try {
    const res = await fetch(`${API_BASE}/me`, { credentials: 'include' });
    if (!res.ok) throw new Error();
    const data = await res.json();
    if (data.status === 'success') {
      return { authenticated: true, nickname: data.user.nickname, ...data.user };
    }
    return { authenticated: false, nickname: '' };
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
