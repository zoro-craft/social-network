"use client";

import { useState, useEffect, useCallback } from "react";
import { apiRequest } from "@/utils/api";

export function useProfile(id) {
  const [profile, setProfile] = useState(null);
  const [error, setError] = useState(null);

  const loadProfile = useCallback(async () => {
    try {
      const data = await apiRequest(`/profile/${id}`);
      
      if (!data.user) {
        throw new Error(data.error || "User not found");
      }

      const validRelations = ["self", "accepted", "pending", "none"];
      const relation = validRelations.includes(data.relation) ? data.relation : "none";

      setProfile({ ...data, relation });
      setError(null);
    } catch (err) {
      setProfile(null);
      setError(err.message);
    }
  }, [id]);

  useEffect(() => {
    loadProfile();
  }, [loadProfile]);

  return { profile, error, reload: loadProfile };
}