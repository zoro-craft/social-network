'use client';

import { useEffect, useRef } from 'react';

// Stable empty array so an unloaded conversation's items reference never changes per render.
const EMPTY = [];

// Thread panel: message bubbles, open-at-bottom scroll, and top paging.
export default function Thread({ chat }) {
  const listRef = useRef(null); // the scrollable message list element.
  const prevKeyRef = useRef(null); // the last rendered conversation key.
  const lastStateRef = useRef(null); // last scroll snapshot for prepend preservation.
  const { auth, selected, messages, hasMore, loadOlder, convKey } = chat;

  const key = selected ? convKey(selected.type, selected.id) : null;
  const items = key ? messages[key] || EMPTY : EMPTY;
  const more = key ? !!hasMore[key] : false;

  // Scroll the message list to the newest message.
  const scrollToBottom = () => {
    const el = listRef.current;
    if (!el) return;
    // Defer to the next tick so the browser paints the new content before the scroll lands.
    setTimeout(() => {
      el.scrollTop = el.scrollHeight;
    }, 0);
  };

  useEffect(() => {
    const el = listRef.current;
    const switched = prevKeyRef.current !== key;
    prevKeyRef.current = key;
    if (!el || items.length === 0) return;

    const prev = lastStateRef.current;
    const lastId = items[items.length - 1].id;
    const height = el.scrollHeight;
    lastStateRef.current = { key, lastId, count: items.length, height };

    // Open or a new message: land at the newest message, even while still loading.
    if (switched || !prev || prev.key !== key || (prev.key === key && prev.lastId !== lastId)) {
      scrollToBottom();
      return;
    }

    // Older batch prepended in place: shift the viewport down so it never jumps to the top.
    if (prev.key === key && prev.lastId === lastId && items.length > prev.count) {
      const added = height - prev.height;
      if (added > 0) {
        el.scrollTop += added;
      }
    }
  }, [key, items]);

  // Load older history when pinned at the top; poll so a throttled load fires on cooldown.
  const handleScroll = () => {
    const el = listRef.current;
    if (el && el.scrollTop <= 0 && more) {
      loadOlder();
    }
  };

  useEffect(() => {
    const id = setInterval(() => {
      const el = listRef.current;
      if (el && el.scrollTop <= 0 && more) {
        loadOlder();
      }
    }, 250);
    return () => clearInterval(id);
  }, [more, loadOlder]);

  return (
    <div className="glass-chat-thread">
      <div className="glass-chat-thread-header">{selected ? selected.name : ''}</div>
      <div className="glass-chat-messages" ref={listRef} onScroll={handleScroll}>
        {items.map((m) => (
          <div
            key={m.id}
            className={`glass-chat-bubble ${m.sender_id === auth.id ? 'glass-chat-bubble-own' : ''}`}
          >
            <span className="glass-chat-bubble-name">{m.sender_name}</span>
            <span className="glass-chat-bubble-content">{m.content}</span>
          </div>
        ))}
      </div>
    </div>
  );
}