'use client';

import { useEffect, useMemo, useRef } from 'react';

export default function Thread({ chat }) {
  const listRef = useRef(null);
  const lastSeenRef = useRef(null);
  const prevHeightRef = useRef(0);
  const { auth, selected, messages, hasMore, loadOlder, convKey } = chat;

  const key = selected ? convKey(selected.type, selected.id) : null;
  const items = useMemo(() => (key ? messages[key] || [] : []), [key, messages]);
  const more = key ? !!hasMore[key] : false;

  // Keep the view in place after the list changes: scroll to the newest on open / new
  // message, and shift down by the added height when an older batch is prepended so the
  // view never jumps to the top.
  useEffect(() => {
    const el = listRef.current;
    if (!el || items.length === 0) return;
    const heightBefore = prevHeightRef.current;
    prevHeightRef.current = el.scrollHeight;
    const firstId = items[0].id;
    const lastId = items[items.length - 1].id;
    const seen = lastSeenRef.current;
    const isNewConversation = !seen || seen.key !== key;
    const hasNewMessage = seen && seen.key === key && seen.lastId !== lastId;
    const hasOlderBatch = seen && seen.key === key && seen.firstId !== firstId;
    lastSeenRef.current = { key, firstId, lastId };

    if (isNewConversation || hasNewMessage) {
      const scrollToBottom = () => {
        el.scrollTop = el.scrollHeight;
      };
      requestAnimationFrame(scrollToBottom);
      requestAnimationFrame(scrollToBottom);
      return;
    }

    if (hasOlderBatch) {
      const added = el.scrollHeight - heightBefore;
      if (added > 0) {
        el.scrollTop += added;
      }
    }
  }, [key, items]);

  // Older history loads when scrolled to the very top. Poll the scroll position so a
  // throttled top-load fires as soon as the cooldown elapses even while the user just
  // hangs at the top, without needing to scroll away and back.
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