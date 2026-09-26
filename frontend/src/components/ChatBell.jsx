'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';
import { useRealtime } from '@/components/RealtimeProvider';
import { conversationKeyFromPayload } from '@/utils/chat';

// Navbar chat icon that owns the browser-only unread counts and broadcasts them.
export default function ChatBell() {
  const router = useRouter();
  const { auth } = useAuth();
  const { subscribe } = useRealtime();
  const [total, setTotal] = useState(0); // app-wide unread total.
  const countsRef = useRef({}); // per-conversation unread counts.
  const openRef = useRef(new Set()); // conversation keys currently open (not counted).

  // Broadcast the current counts so the chat page can render tab and row badges.
  const sync = () => {
    const counts = countsRef.current;
    let usersTab = 0;
    let groupsTab = 0;
    for (const key in counts) {
      if (key.startsWith('private:')) usersTab += counts[key];
      else groupsTab += counts[key];
    }
    setTotal(usersTab + groupsTab);
    window.dispatchEvent(
      new CustomEvent('chat:unread', {
        detail: { users: usersTab, groups: groupsTab, counts: { ...counts } },
      }),
    );
  };

  // +1 only for messages from someone else whose conversation is not open.
  const bump = (payload) => {
    if (payload.sender_id === auth.id) return;
    const key = conversationKeyFromPayload(payload, auth.id);
    if (openRef.current.has(key)) return;
    countsRef.current[key] = (countsRef.current[key] || 0) + 1;
    sync();
  };

  useEffect(() => {
    if (!auth.authenticated) return;
    const unsub1 = subscribe('chat_new_message', bump);
    const unsub2 = subscribe('chat_new_group_message', bump);
    return () => {
      unsub1();
      unsub2();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [auth.authenticated, subscribe, auth.id]);

  // The chat page tells us which conversation is open; opening one resets its count.
  useEffect(() => {
    if (!auth.authenticated) return;
    const onOpen = (e) => {
      const key = e.detail && e.detail.key;
      openRef.current = key ? new Set([key]) : new Set();
      if (key && countsRef.current[key]) {
        countsRef.current[key] = 0;
        sync();
      }
    };
    window.addEventListener('chat:open', onOpen);
    return () => window.removeEventListener('chat:open', onOpen);
  }, [auth.authenticated]);

  // The chat page requests the current counts on mount so rows/tabs show pre-existing unread.
  useEffect(() => {
    if (!auth.authenticated) return;
    const onSync = () => sync();
    window.addEventListener('chat:sync', onSync);
    return () => window.removeEventListener('chat:sync', onSync);
  }, [auth.authenticated]);

  // Counts are browser-only and start at zero on every sign-in or refresh.
  useEffect(() => {
    if (auth.authenticated) return;
    countsRef.current = {};
    openRef.current = new Set();
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setTotal(0);
  }, [auth.authenticated]);

  return (
    <button
      className="glass-nav-search glass-nav-bell glass-chat-bell"
      onClick={() => router.push('/chat')}
      aria-label="Chat"
    >
      💬
      {total > 0 && <span className="glass-bell-badge">{total}</span>}
    </button>
  );
}