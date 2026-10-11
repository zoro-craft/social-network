'use client';

import { useEffect } from 'react';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { useRealtime } from '@/app/components/RealtimeProvider';
import ChatSidebar from './ChatSidebar';
import Thread from './Thread';
import Composer from './Composer';
import { useChat } from './useChat';

// Chat page: sidebar, thread, and composer, guarded by auth.
export default function ChatPage() {
  const { auth, loading } = useRequireAuth();
  const { subscribe } = useRealtime();
  const chat = useChat(auth, subscribe);

  // Tell the bell nothing is open on mount, and clear it again on leave.
  useEffect(() => {
    if (!auth.authenticated) return;
    window.dispatchEvent(new CustomEvent('chat:open', { detail: { key: null } }));
    return () => window.dispatchEvent(new CustomEvent('chat:open', { detail: { key: null } }));
  }, [auth.authenticated]);

  if (loading) return <div>Loading...</div>;
  if (!auth.authenticated) return null;

  return (
    <div className="glass-chat-page">
      <ChatSidebar chat={chat} />
      <div className="glass-chat-main">
        {chat.selected ? (
          <>
            <Thread chat={chat} />
            <Composer chat={chat} />
          </>
        ) : (
          <div className="glass-chat-empty">
            <p>Select a conversation to start chatting</p>
          </div>
        )}
      </div>
    </div>
  );
}