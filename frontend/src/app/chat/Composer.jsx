'use client';

import { useRef, useState } from 'react';
import EmojiPicker from './EmojiPicker';

export default function Composer({ chat }) {
  const [text, setText] = useState('');
  const [pickerOpen, setPickerOpen] = useState(false);
  const inputRef = useRef(null);
  const { selected, canChat, send, convKey } = chat;

  const key = selected ? convKey(selected.type, selected.id) : null;
  const disabled = !key || canChat[key] === false;

  const handleSend = async () => {
    const content = text;
    if (!content.trim()) return;
    const res = await send(content);
    if (res.ok) {
      setText('');
      setPickerOpen(false);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const insertEmoji = (emoji) => {
    const el = inputRef.current;
    const start = el.selectionStart ?? text.length;
    const end = el.selectionEnd ?? text.length;
    setText(text.slice(0, start) + emoji + text.slice(end));
    setPickerOpen(false);
    requestAnimationFrame(() => {
      el.focus();
      const pos = start + emoji.length;
      el.setSelectionRange(pos, pos);
    });
  };

  return (
    <div className="glass-chat-composer">
      {pickerOpen && <EmojiPicker onPick={insertEmoji} />}
      <div className="glass-chat-composer-row">
        <button
          className="glass-chat-emoji-btn"
          onClick={() => setPickerOpen((v) => !v)}
          disabled={disabled}
          aria-label="Toggle emoji picker"
        >
          😊
        </button>
        <textarea
          ref={inputRef}
          className="glass-chat-input"
          value={text}
          maxLength={5000}
          disabled={disabled}
          placeholder={disabled ? 'You can no longer send here' : 'Write a message...'}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <button className="glass-chat-send" onClick={handleSend} disabled={disabled || !text.trim()}>
          Send
        </button>
      </div>
    </div>
  );
}