import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
} from 'react';
import { useAuth } from '@/context/AuthContext';
import { getSessionToken, wsUrl } from '@/lib/api';
import { handleSessionLost } from '@/lib/session';

type Handler = (payload: any) => void;

interface RealtimeContextValue {
  subscribe: (type: string, handler: Handler) => () => void;
}

const RealtimeContext = createContext<RealtimeContextValue>({
  subscribe: () => () => {},
});

export function RealtimeProvider({ children }: { children: React.ReactNode }) {
  const { auth } = useAuth();
  const handlersRef = useRef<Map<string, Set<Handler>>>(new Map());
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const cancelledRef = useRef(false);

  const subscribe = useCallback((type: string, handler: Handler) => {
    if (!handlersRef.current.has(type)) handlersRef.current.set(type, new Set());
    handlersRef.current.get(type)!.add(handler);
    return () => {
      handlersRef.current.get(type)?.delete(handler);
    };
  }, []);

  useEffect(() => {
    cancelledRef.current = false;

    if (!auth.authenticated) {
      if (wsRef.current) {
        try {
          wsRef.current.close();
        } catch {
          // ignore
        }
        wsRef.current = null;
      }
      return;
    }

    let cancelled = false;
    let socket: WebSocket | null = null;

    const connect = async () => {
      if (cancelled || cancelledRef.current) return;
      const token = await getSessionToken();
      try {
        type WSConstructor = new (
          url: string,
          protocols?: string | string[] | null,
          options?: { headers?: Record<string, string> },
        ) => WebSocket;
        socket = new ((WebSocket as unknown) as WSConstructor)(wsUrl(), null, {
          headers: token ? { Cookie: `session_token=${token}` } : {},
        });
        wsRef.current = socket;

        socket.onopen = () => {
          // connected
        };

        socket.onclose = () => {
          if (!cancelled && !cancelledRef.current) {
            reconnectRef.current = setTimeout(connect, 3000);
          }
        };

        socket.onerror = () => {
          // handled by onclose
        };

        socket.onmessage = (event) => {
          try {
            const envelope = JSON.parse((event as any).data ?? event.data);
            const handlers = handlersRef.current.get(envelope.type);
            if (handlers) {
              handlers.forEach((h) => {
                try {
                  h(envelope.payload);
                } catch {
                  // ignore handler errors
                }
              });
            }
          } catch {
            // ignore malformed frames
          }
        };
      } catch {
        reconnectRef.current = setTimeout(connect, 3000);
      }
    };

    connect();

    return () => {
      cancelled = true;
      cancelledRef.current = true;
      if (reconnectRef.current) clearTimeout(reconnectRef.current);
      if (socket) {
        try {
          socket.close();
        } catch {
          // ignore
        }
      }
      socket = null;
      wsRef.current = null;
    };
  }, [auth.authenticated]);

  const value = useMemo(() => ({ subscribe }), [subscribe]);

  return <RealtimeContext.Provider value={value}>{children}</RealtimeContext.Provider>;
}

export function useRealtime(): RealtimeContextValue {
  return useContext(RealtimeContext);
}

export function sessionLostFromSocket() {
  handleSessionLost();
}