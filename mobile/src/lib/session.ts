type SessionLostHandler = () => void;

let handlers: SessionLostHandler[] = [];

export function registerSessionLost(handler: SessionLostHandler): () => void {
  handlers.push(handler);
  return () => {
    handlers = handlers.filter((h) => h !== handler);
  };
}

export function handleSessionLost(): void {
  handlers.forEach((h) => {
    try {
      h();
    } catch {
      // ignore
    }
  });
}