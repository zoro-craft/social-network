import { apiRequest } from "@/utils/api";

export async function fetchGroups() {
    try {
        const data = await apiRequest("/groups");
        return data;
    } catch (error) {
        console.error("fetchGroups error:", error);
        throw error;
    }
}