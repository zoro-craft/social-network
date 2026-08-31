'use client';

import { useState } from 'react';
import { API_BASE } from '../../lib/api';

export function useFollow(userId) {
  const [status, setStatus] = useState('none'); // none | pending | accepted

  const fetchStatus = async () => {
    try {
      const res = await fetch(`${API_BASE}/follow/${userId}/status`, {
        credentials: 'include',
      });
      const data = await res.json();
      setStatus(data.status || 'none');
    } catch {
      setStatus('none');
    }
  };

  const follow = async () => {
    const res = await fetch(`${API_BASE}/follow/${userId}`, {
      method: 'POST',
      credentials: 'include',
    });
    if (res.ok) {
      const data = await res.json();
      setStatus(data.status);
    }
    return res;
  };

  const unfollow = async () => {
    const res = await fetch(`${API_BASE}/follow/${userId}/unfollow`, {
      method: 'POST',
      credentials: 'include',
    });
    if (res.ok) setStatus('none');
    return res;
  };

  return { status, fetchStatus, follow, unfollow };
}
