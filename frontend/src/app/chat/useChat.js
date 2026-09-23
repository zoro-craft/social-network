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

