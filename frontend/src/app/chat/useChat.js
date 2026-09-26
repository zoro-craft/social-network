'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { apiRequest } from '@/utils/api';
import { conversationKeyFromPayload, otherParticipant } from '@/utils/chat';

const BATCH = 10;
const PAGING_COOLDOWN = 2000;

// useChat owns all chat state and realtime wiring for the /chat page.
export function useChat(auth, subscribe) {
  const [users, setUsers] = useState([]); // sidebar users tab.
  const [groups, setGroups] = useState([]); // sidebar groups tab.
  const [selected, setSelected] = useState(null); // currently open conversation.
  const [messages, setMessages] = useState({}); // per-key thread messages.
  const [hasMore, setHasMore] = useState({}); // per-key older-history availability.
  const [canChat, setCanChat] = useState({}); // per-key send permission.
  const [badges, setBadges] = useState({ users: 0, groups: 0, counts: {} }); // unread totals from the bell.

  const selectedRef = useRef(selected); // latest selected read by socket callbacks.
  // eslint-disable-next-line react-hooks/refs -- latest-value ref read by socket callbacks
  selectedRef.current = selected;
  const loadingMoreRef = useRef({}); // per-key in-flight older-batch guard.
  const lastFetchRef = useRef({}); // per-key last fetch time for the cooldown.

  // Build a stable key for a conversation in state maps and window events.
  const convKey = (type, id) => (type === 'group' ? `group:${id}` : `private:${id}`);

  // Snapshot semantics: fetch the sidebar once on page mount only.
  useEffect(() => {
    if (!auth.authenticated) return;
    let cancelled = false;
    apiRequest('/conversations')
      .then((data) => {
        if (cancelled) return;
        const userList = data.users || [];
        const groupList = data.groups || [];
        setUsers(userList);
        setGroups(groupList);
        const cc = {};
        userList.forEach((u) => {
          cc[`private:${u.user_id}`] = u.can_chat;
        });
        groupList.forEach((g) => {
          cc[`group:${g.group_id}`] = true;
        });
        setCanChat(cc);
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [auth.authenticated]);

  const loadHistory = useCallback(async (type, id, beforeId, isMore) => {
    const key = convKey(type, id);
    const path = type === 'group' ? `/groups/${id}/messages` : `/messages/${id}`;
    const query = `?limit=${BATCH}${beforeId ? `&before=${beforeId}` : ''}`;
    try {
      const data = await apiRequest(`${path}${query}`);
      const batch = (data.messages || []).slice().reverse();
      setMessages((prev) => {
        // Fresh load: the server batch is authoritative and replaces the cached thread.
        if (!isMore) {
          return { ...prev, [key]: batch };
        }
        const seen = new Set((prev[key] || []).map((m) => m.id));
        const fresh = batch.filter((m) => !seen.has(m.id));
        const merged = [...fresh, ...(prev[key] || [])];
        return { ...prev, [key]: merged };
      });
      setHasMore((prev) => ({ ...prev, [key]: !!data.has_more }));
    } catch (e) {
      if (e.status === 400) {
        setCanChat((prev) => ({ ...prev, [key]: false }));
      }
      setHasMore((prev) => ({ ...prev, [key]: false }));
    }
  }, []);

  // Always re-fetch the newest batch on select so messages that arrived while closed show up.
  const selectConversation = useCallback(
    (type, id, name) => {
      const key = convKey(type, id);
      setSelected({ type, id, name });
      window.dispatchEvent(new CustomEvent('chat:open', { detail: { key } }));
      loadHistory(type, id, null, false);
    },
    [loadHistory],
  );

  // fetchOlder prepends the batch older than the oldest loaded message.
  const fetchOlder = useCallback((type, id) => {
    const key = convKey(type, id);
    const oldest = (messages[key] || [])[0];
    if (!oldest) return;
    loadingMoreRef.current[key] = true;
    lastFetchRef.current[key] = Date.now();
    loadHistory(type, id, oldest.id, true).finally(() => {
      loadingMoreRef.current[key] = false;
    });
  }, [messages, loadHistory]);

  // Polled by Thread while pinned at the top; the 2s cooldown gates older-batch loads.
  const loadOlder = useCallback(() => {
    const sel = selectedRef.current;
    if (!sel) return;
    const { type, id } = sel;
    const key = convKey(type, id);
    if (loadingMoreRef.current[key] || !hasMore[key]) return;
    if (PAGING_COOLDOWN - (Date.now() - (lastFetchRef.current[key] || 0)) > 0) return;
    fetchOlder(type, id);
  }, [hasMore, fetchOlder]);

  const send = useCallback(
    async (content) => {
      const sel = selectedRef.current;
      if (!sel) return { ok: false };
      const { type, id } = sel;
      const key = convKey(type, id);
      const body = type === 'group' ? { content } : { recipient_id: id, content };
      const path = type === 'group' ? `/groups/${id}/messages` : `/messages`;
      try {
        await apiRequest(path, 'POST', body);
        return { ok: true };
      } catch (e) {
        if (e.status === 400) {
          setCanChat((prev) => ({ ...prev, [key]: false }));
        }
        return { ok: false };
      }
    },
    [],
  );

  // Live arrivals: refresh sidebar rows always, and append to the thread when it is open.
  useEffect(() => {
    if (!auth.authenticated) return;
    const onMessage = (payload) => {
      const key = conversationKeyFromPayload(payload, auth.id);
      const otherID = otherParticipant(payload, auth.id);
      const sel = selectedRef.current;
      const isOpen = sel && convKey(sel.type, sel.id) === key;

      if (isOpen) {
        setMessages((prev) => {
          const list = prev[key] || [];
          if (list.some((m) => m.id === payload.id)) return prev;
          return { ...prev, [key]: [...list, payload] };
        });
      }

      if (payload.type === 'private') {
        setUsers((prev) => {
          const entry = prev.find((u) => u.user_id === otherID);
          if (!entry) return prev;
          const updated = { ...entry, last_message: payload.content, last_message_at: payload.created_at };
          return [updated, ...prev.filter((u) => u.user_id !== otherID)];
        });
      } else {
        setGroups((prev) => {
          const entry = prev.find((g) => g.group_id === otherID);
          if (!entry) return prev;
          const updated = { ...entry, last_message: payload.content, last_message_at: payload.created_at };
          return [updated, ...prev.filter((g) => g.group_id !== otherID)];
        });
      }
    };
    const unsub1 = subscribe('chat_new_message', onMessage);
    const unsub2 = subscribe('chat_new_group_message', onMessage);
    return () => {
      unsub1();
      unsub2();
    };
  }, [auth.authenticated, auth.id, subscribe]);

  // Presence envelopes update the online dots live.
  useEffect(() => {
    if (!auth.authenticated) return;
    return subscribe('presence', (payload) => {
      setUsers((prev) =>
        prev.map((u) => (u.user_id === payload.user_id ? { ...u, online: payload.online } : u)),
      );
    });
  }, [auth.authenticated, subscribe]);

  // The bell owns the counts; the page renders them from the chat:unread broadcast.
  useEffect(() => {
    if (!auth.authenticated) return;
    const onUnread = (e) =>
      setBadges({ users: e.detail.users, groups: e.detail.groups, counts: e.detail.counts });
    window.addEventListener('chat:unread', onUnread);
    // Ask the bell to broadcast its current counts so pre-existing unread reaches rows/tabs.
    window.dispatchEvent(new CustomEvent('chat:sync'));
    return () => window.removeEventListener('chat:unread', onUnread);
  }, [auth.authenticated]);

  return {
    auth,
    users,
    groups,
    selected,
    messages,
    hasMore,
    canChat,
    badges,
    convKey,
    selectConversation,
    loadOlder,
    send,
  };
}