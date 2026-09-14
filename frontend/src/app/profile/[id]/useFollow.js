"use client";

import { useState } from "react";
import { apiPost } from "@/utils/api";
import { useToast } from "@/app/context/ToastContext";

export function useFollow(id, reload) {
  const showToast = useToast();
  const [pending, setPending] = useState(false);

  const toggleFollow = async () => {
    setPending(true);
    try {
      await apiPost(`/follow/${id}`);
      await reload();
    } catch (err) {
      showToast(err.message);
    } finally {
      setPending(false);
    }
  };

  return { toggleFollow, pending };
}
