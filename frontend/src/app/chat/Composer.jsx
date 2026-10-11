'use client';

import { useRef, useState } from 'react';
import EmojiPicker from './EmojiPicker';
import Icon from '@/app/components/Icon';

// Composer panel: message input, Enter-to-send, and the emoji picker toggle.
export default function Composer({ chat }) {
  const [text, setText] = useState(''); // draft message content.
  const [pickerOpen, setPickerOpen] = useState(false); // whether the emoji picker is open.
  const inputRef = useRef(null);
  const { selected, canChat, send, convKey } = chat;

  const key = selected ? convKey(selected.type, selected.id) : null;
  const disabled = !key || canChat[key] === false;

  // Send the current draft and clear it on success.
  const handleSend = async () => {
    const content = text;
    if (!content.trim()) return;
    const res = await send(content);
    if (res.ok) {
      setText('');
      setPickerOpen(false);
    }
  };

  // Enter sends; Shift+Enter inserts a new line instead.
  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // Insert a picked emoji at the cursor and restore focus to the input.
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
          <Icon name="smile" />
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