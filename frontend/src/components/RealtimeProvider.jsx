'use client';

import { createContext, useContext, useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';
import { registerSessionLost } from '@/utils/session';
import { API_BASE } from '@/utils/api';

const RealtimeContext = createContext({
  subscribe: () => () => {},
  isConnected: false,
});

export function useRealtime() {
  return useContext(RealtimeContext);
}

export function RealtimeProvider({ children }) {
  const { auth, setAuth } = useAuth();
  const router = useRouter();
  const socketRef = useRef(null);
  const handlersRef = useRef({});
  const [isConnected, setIsConnected] = useState(false);

  // single shared logout routine triggered by the /me revalidation loop;
  // the socket teardown below runs when auth flips to logged-out
  const handleLogout = () => {
    setAuth({ authenticated: false, nickname: '', id: 0 });
    router.push('/login');
  };

  useEffect(() => {
    // register once; handleLogout is stable (setAuth + router)
    registerSessionLost(handleLogout);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const subscribe = (type, handler) => {
    const handlers = handlersRef.current;
    handlers[type] = handlers[type] || [];
    handlers[type].push(handler);

    return () => {
      handlers[type] = (handlers[type] || []).filter(
        (h) => h !== handler
      );
    };
  };

  useEffect(() => {
    if (!auth.authenticated) {
      if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
      }
      return;
    }

    const ws = new WebSocket(
      `ws://${API_BASE.replace(/^https?:\/\//, '')}/ws`
    );
    socketRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      let env;
      try {
        env = JSON.parse(event.data);
      } catch {
        return;
      }

      const typeHandlers = handlersRef.current[env.type];
      if (typeHandlers) {
        typeHandlers.forEach((handler) => handler(env.payload));
      }
    };

    ws.onclose = () => {
      setIsConnected(false);
      socketRef.current = null;
    };

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
      }
    };
  }, [auth.authenticated]);

  return (
    <RealtimeContext.Provider value={{ subscribe, isConnected }}>
      {children}
    </RealtimeContext.Provider>
  );
}