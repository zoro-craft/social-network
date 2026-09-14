export const API_BASE = "http://localhost:8080";

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function parseResponse(res) {
  const text = await res.text();
  if (!text) return {};

  try {
    return JSON.parse(text);
  } catch {
    return { error: text };
  }
}

export async function apiRequest(
  endpoint,
  { method = "GET", body = null, isForm = false } = {},
) {
  const options = { method, credentials: "include", headers: {} };

  if (body) {
    if (isForm) {
      options.body = body;
    } else {
      options.headers["Content-Type"] = "application/json";
      options.body = JSON.stringify(body);
    }
  }

  const res = await fetch(`${API_BASE}${endpoint}`, options);
  const result = await parseResponse(res);

  if (!res.ok) {
    throw new ApiError(
      result.message || result.error || "Something went wrong!",
      res.status,
    );
  }

  return result;
}

export const apiGet = (endpoint) => apiRequest(endpoint);

export const apiPost = (endpoint, body = null) =>
  apiRequest(endpoint, { method: "POST", body });

export const apiPostForm = (endpoint, formData) =>
  apiRequest(endpoint, { method: "POST", body: formData, isForm: true });
