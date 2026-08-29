export async function apiRequest(endpoint, method = "GET", data = null) {
  const url = `http://localhost:8080${endpoint}`;

  const options = {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  };

  if (data) {
    options.body = JSON.stringify(data);
  }
  const res = await fetch(url, options);
  const result = await res.json();

  if (!res.ok) {
    throw new Error(result.message || "Something went wrong!");
  }
  return result;
}

export async function logoutUser() {
  return await apiRequest("/api/logout", "POST");
}