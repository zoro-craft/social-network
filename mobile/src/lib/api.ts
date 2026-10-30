import { Platform } from 'react-native';
import Constants from 'expo-constants';
import * as SecureStore from 'expo-secure-store';

// The Go server speaks /api directly on :8080 (see tolerateAPIPrefix), so no
// proxy is required for local development. When you do run Caddy instead, set
// EXPO_PUBLIC_API_URL to its :8000 origin.
//
// Native picks the LAN IP from hostUri so a physical phone reaches your machine;
// 10.0.2.2 is the Android emulator's alias for the host.
function getDefaultHost(): string {
  if (Platform.OS === 'web') return 'http://localhost:8080';

  const expoHost = Constants.expoConfig?.hostUri
    ?.replace(/^https?:\/\//, '')
    .split(':')[0];
  if (expoHost && expoHost !== 'localhost' && expoHost !== '127.0.0.1') {
    return `http://${expoHost}:8080`;
  }

  return Platform.OS === 'android' ? 'http://10.0.2.2:8080' : 'http://localhost:8080';
}

const DEFAULT_HOST = getDefaultHost();

export const API_ORIGIN = (process.env.EXPO_PUBLIC_API_URL || DEFAULT_HOST).replace(/\/+$/, '');
export const API_BASE = `${API_ORIGIN}/api`;

// The backend authenticates with an HttpOnly "session_token" cookie.
// React Native has no same-origin cookie jar, so we capture the token from
// Set-Cookie on /login (and /register) and send it back on every request.
const SESSION_KEY = 'session_token';

let cachedToken: string | null = null;

export async function getSessionToken(): Promise<string | null> {
  if (cachedToken !== null) return cachedToken;
  try {
    cachedToken = await SecureStore.getItemAsync(SESSION_KEY);
  } catch {
    cachedToken = null;
  }
  return cachedToken;
}

export async function setSessionToken(token: string | null) {
  cachedToken = token;
  try {
    if (token) {
      await SecureStore.setItemAsync(SESSION_KEY, token);
    } else {
      await SecureStore.deleteItemAsync(SESSION_KEY);
    }
  } catch {
    // best effort; the in-memory token still works for the session
  }
}

// Server image paths are relative, e.g. /uploads/avatars/x. The web frontend
// builds them as `${API_BASE}${src}`; we need a fully-qualified URI for <Image>.
export function imageUrl(src?: string | null): string | null {
  if (!src) return null;
  const clean = src.startsWith('/') ? src : `/${src}`;
  return `${API_ORIGIN}/api${clean}`;
}

export type MediaSource = {
  uri: string;
  headers?: Record<string, string>;
};

function mediaHeaders(): Record<string, string> | undefined {
  if (Platform.OS === 'web' || !cachedToken) return undefined;
  return { Cookie: `session_token=${cachedToken}` };
}

export function mediaSource(src?: string | null): MediaSource | null {
  const uri = imageUrl(src);
  if (!uri) return null;
  const headers = mediaHeaders();
  return headers ? { uri, headers } : { uri };
}

// Reel videos are served from the same /uploads root, so they need the identical
// absolute-URI treatment; this alias just keeps the call sites readable.
export const mediaUrl = imageUrl;

function parseSessionToken(setCookie: string | string[] | null | undefined): string | null {
  if (!setCookie) return null;
  const values = Array.isArray(setCookie) ? setCookie : [setCookie];
  for (const value of values) {
    const match = /session_token=([^;,]+)/.exec(value);
    if (match) return match[1];
  }
  return null;
}

export type ApiOptions = {
  method?: string;
  data?: unknown;
  formData?: FormData;
  headers?: Record<string, string>;
};

export type ApiResponse<T = unknown> = {
  ok: boolean;
  status: number;
  data: T;
  headers: Record<string, string>;
};

async function request<T = unknown>(endpoint: string, options: ApiOptions = {}): Promise<ApiResponse<T>> {
  const { method = 'GET', data, formData, headers = {} } = options;
  const url = `${API_BASE}${endpoint}`;

  const requestHeaders: Record<string, string> = { ...headers };
  if (data !== undefined) requestHeaders['Content-Type'] = 'application/json';

  // Only hand over the cookie header manually on native. Browsers forbid setting
  // it and rely on the same-origin proxy + CORS instead.
  if (Platform.OS !== 'web') {
    const token = await getSessionToken();
    if (token) requestHeaders['Cookie'] = `session_token=${token}`;
  }

  const init: RequestInit = {
    method,
    headers: requestHeaders,
    credentials: 'include',
  };
  if (formData) {
    init.body = formData as unknown as BodyInit;
  } else if (data !== undefined) {
    init.body = JSON.stringify(data);
  }

  const res = await fetch(url, init);

  // Capture the session cookie whenever the server sets one (login / register).
  const setCookie = (res.headers as unknown as {
    get?: (name: string) => string | null;
  }).get?.('set-cookie');
  if (setCookie) {
    const token = parseSessionToken(setCookie);
    if (token) await setSessionToken(token);
  }

  const headersMap: Record<string, string> = {};
  res.headers.forEach?.((value, key) => {
    headersMap[key.toLowerCase()] = value;
  });

  const text = await res.text();
  let result: any = {};
  if (text) {
    try {
      result = JSON.parse(text);
    } catch {
      result = { error: text };
    }
  }

  return { ok: res.ok, status: res.status, data: result as T, headers: headersMap };
}

// User-friendly helper that throws on non-2xx, mirroring frontend utils/api.js.
export async function apiRequest<T = unknown>(
  endpoint: string,
  options: ApiOptions = {},
): Promise<T> {
  const res = await request<T>(endpoint, options);
  if (!res.ok) {
    const message =
      (res.data as any)?.error || (res.data as any)?.message || 'Something went wrong!';
    const error: any = new Error(message);
    error.status = res.status;
    throw error;
  }
  return res.data;
}

// Access to raw response (status + headers) for cursor-based pagination
// that reads X-Max-* headers the way the web frontend does.
export async function apiRaw<T = unknown>(
  endpoint: string,
  options: ApiOptions = {},
): Promise<ApiResponse<T>> {
  return request<T>(endpoint, options);
}

export function wsUrl(): string {
  const host = API_ORIGIN.replace(/^https?:\/\//, '');
  const proto = API_ORIGIN.startsWith('https:') ? 'wss' : 'ws';
  return `${proto}://${host}/api/ws`;
}