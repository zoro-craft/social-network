import { headers } from "next/headers";
import { API_BASE } from "@/utils/api";

async function serverApiRequest(endpoint) {
  const requestHeaders = await headers();

  const host = requestHeaders.get("x-forwarded-host") || requestHeaders.get("host");
  const proto = requestHeaders.get("x-forwarded-proto") || "http";

  const cookie = requestHeaders.get("cookie") || "";

  const res = await fetch(`${proto}://${host}${API_BASE}${endpoint}`, {
    method: "GET",
    headers: { Cookie: cookie },
    cache: "no-store",
  });

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
    const error = new Error(result.error);
    error.status = res.status;
    throw error;
  }

  return result;
}

export async function getProfile(id) {
  try {
    const data = await serverApiRequest(`/profile/${id}`);

    if (!data.user) {
      throw new Error(data.error || "User not found");
    }

    const validRelations = ["self", "accepted", "pending", "none"];
    const relation = validRelations.includes(data.relation) ? data.relation : "none";

    return { profile: { ...data, relation }, error: null };
  } catch (err) {
    return { profile: null, error: err.message };
  }
}
