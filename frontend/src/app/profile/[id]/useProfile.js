"use client";

import { useCallback, useEffect, useState } from "react";
import { apiGet } from "@/utils/api";

const VALID_RELATIONS = ["self", "accepted", "pending", "none"];

const normalizeRelation = (relation) =>
  VALID_RELATIONS.includes(relation) ? relation : "none";

export function useProfile(id) {
  const [profile, setProfile] = useState(null);
  const [error, setError] = useState(null);

  const loadProfile = useCallback(
    () =>
      apiGet(`/profile/${id}`).then((data) => {
        if (!data.user) throw new Error(data.error || "User not found");
        return { ...data, relation: normalizeRelation(data.relation) };
      }),
    [id],
  );

  const applyProfile = useCallback((data) => {
    setProfile(data);
    setError(null);
  }, []);

  const applyError = useCallback((err) => {
    setProfile(null);
    setError(err.message);
  }, []);

  const reload = useCallback(() => {
    loadProfile().then(applyProfile, applyError);
  }, [loadProfile, applyProfile, applyError]);

  useEffect(() => {
    let cancelled = false;
    const ifActive = (fn) => (value) => {
      if (!cancelled) fn(value);
    };

    loadProfile().then(ifActive(applyProfile), ifActive(applyError));

    return () => {
      cancelled = true;
    };
  }, [loadProfile, applyProfile, applyError]);

  return { profile, error, reload };
}
