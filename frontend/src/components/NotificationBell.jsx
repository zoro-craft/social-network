'use client';

import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';
import { useRealtime } from '@/components/RealtimeProvider';
import { apiRequest } from '@/utils/api';

// Navbar bell showing the server-supplied exact unread total.
export default function NotificationBell() {
  const router = useRouter();
  const { auth } = useAuth();
  const { subscribe } = useRealtime();
  const [unread, setUnread] = useState(0); // exact unread notification total.

  // Fetch the server-supplied exact unread total for the badge.
  const refreshUnread = useCallback(async () => {
    try {
      const data = await apiRequest('/notifications');
      setUnread(data.unread_count || 0);
    } catch {
      // keep the previous count on transient errors
    }
  }, []);

  // Initial and sign-in catch-up: show the server-supplied exact unread total.
  useEffect(() => {
    if (!auth.authenticated) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setUnread(0);
      return;
    }
    refreshUnread();
  }, [auth.authenticated, refreshUnread]);

  // Live updates: refetch the exact total on each new notification.
  useEffect(() => {
    if (!auth.authenticated) return;
    return subscribe('notification', refreshUnread);
  }, [auth.authenticated, subscribe, refreshUnread]);

  // Removal signal: refetch so the badge mimics the server-side deletion.
  useEffect(() => {
    if (!auth.authenticated) return;
    return subscribe('notification_removed', refreshUnread);
  }, [auth.authenticated, subscribe, refreshUnread]);

  // After acting on a notification, refresh the badge to the exact unread total.
  useEffect(() => {
    if (!auth.authenticated) return;
    window.addEventListener('notifications:changed', refreshUnread);
    return () => window.removeEventListener('notifications:changed', refreshUnread);
  }, [auth.authenticated, refreshUnread]);

  return (
    <button
      className="glass-nav-search glass-nav-bell"
      onClick={() => router.push('/notifications')}
      aria-label="Notifications"
    >
      🔔
      {unread > 0 && <span className="glass-bell-badge">{unread}</span>}
    </button>
  );
}