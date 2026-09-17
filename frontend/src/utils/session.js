// session.js holds a single shared session-loss routine.
// It is registered once by RealtimeProvider and triggered by a failed /me on the 10s revalidation loop.

let sessionLostHandler = () => {};

export function registerSessionLost(handler) {
  sessionLostHandler = handler;
}

export function handleSessionLost() {
  sessionLostHandler();
}