'use client';

import { API_BASE } from '../../lib/api';

export function useVisibility(setProfile) {
  const toggleVisibility = async () => {
    const res = await fetch(`${API_BASE}/profile/toggle-visibility`, {
      method: 'POST',
      credentials: 'include',
    });

    if (res.ok) {
      setProfile((prev) => ({ ...prev, is_public: !prev.is_public }));
    }

    return res;
  };

  return { toggleVisibility };
}