export const API_BASE = '/api';

export async function apiRequest(
  endpoint,
  method = 'GET',
  data = null,
  isFormData = false,
) {
  const url = `${API_BASE}${endpoint}`;

  const options = {
    method,
    headers: {},
    credentials: 'include',
  };

  if (data) {
    if (isFormData) {
      options.body = data;
    } else {
      options.headers['Content-Type'] = 'application/json';
      options.body = JSON.stringify(data);
    }
  }

  const res = await fetch(url, options);

  const text = await res.text();

  let result = {};

  if (text) {
    try {
      result = JSON.parse(text);
    } catch {
      result = { error: text };
    }
  }

  if (!res.ok) {
    throw new Error(result.message  ||  "Something went wrong!");
    }

  return result;
}
