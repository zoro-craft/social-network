'use client';

const EMOJIS = [
  '😀', '😄', '😁', '😂', '🤣', '😊', '😍', '😘',
  '😎', '🤩', '😢', '😭', '😡', '🤔', '🙄', '😴',
  '👍', '👎', '👏', '🙌', '🤝', '🙏', '💪', '🔥',
  '✨', '❤️', '💔', '🎉', '🎂', '🌈', '⭐', '☀️',
  '🌙', '⚽', '🎮', '🚀', '💯', '✅', '❌', '🍕',
];

export default function EmojiPicker({ onPick }) {
  return (
    <div className="glass-chat-emoji-picker">
      {EMOJIS.map((e) => (
        <button key={e} type="button" className="glass-chat-emoji" onClick={() => onPick(e)}>
          {e}
        </button>
      ))}
    </div>
  );
}