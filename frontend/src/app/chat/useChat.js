'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { apiRequest } from '@/utils/api';
import { conversationKeyFromPayload, otherParticipant } from '@/utils/chat';

const BATCH = 10;
const PAGING_COOLDOWN = 2000;

export function useChat(auth, subscribe) {
  const [users, setUsers] = useState([]);
  const [groups, setGroups] = useState([]);
  const [selected, setSelected] = useState(null);
  const [messages, setMessages] = useState({});
  const [hasMore, setHasMore] = useState({});
  const [canChat, setCanChat] = useState({});
  const [badges, setBadges] = useState({ users: 0, groups: 0, counts: {} });

  const selectedRef = useRef(selected);
  // eslint-disable-next-line react-hooks/refs -- latest-value ref read by socket callbacks
  selectedRef.current = selected;
  const loadingMoreRef = useRef({});
  const lastFetchRef = useRef({});

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
        const seen = new Set((prev[key] || []).map((m) => m.id));
        const fresh = batch.filter((m) => !seen.has(m.id));
        const merged = isMore ? [...fresh, ...(prev[key] || [])] : [...fresh];
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

  const selectConversation = useCallback(
    (type, id, name) => {
      const key = convKey(type, id);
      setSelected({ type, id, name });
      window.dispatchEvent(new CustomEvent('chat:open', { detail: { key } }));
      if (!messages[key]) {
        loadHistory(type, id, null, false);
      }
    },
    [messages, loadHistory],
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

  // loadOlder is polled by Thread while the user stays at the very top; the 2s cooldown
  // naturally gates how often the next older batch can load.
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
