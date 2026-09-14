export const API_BASE = "http://localhost:8080";

export async function apiRequest(endpoint, method = "GET", data = null, isFormData = false) {
  const url = `http://localhost:8080${endpoint}`;

  const options = {
    method: method,
    headers: {},
    credentials: "include",
  };

  if (data) {
    if (isFormData) {
      options.body = data;
    } else {
      options.headers["Content-Type"] = "application/json";
      options.body = JSON.stringify(data);
    }
  }

  const res = await fetch(url, options);
  const result = await res.json();

  if (!res.ok) {
    throw new Error(result.error  || "Something went wrong!");
  }

  return result;
}