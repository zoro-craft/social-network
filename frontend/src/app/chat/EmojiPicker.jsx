'use client';

// Fixed set of common emojis shown in the picker.
const EMOJIS = [
  '😀', '😄', '😁', '😂', '🤣', '😊', '😍', '😘',
  '😎', '🤩', '😢', '😭', '😡', '🤔', '🙄', '😴',
  '👍', '👎', '👏', '🙌', '🤝', '🙏', '💪', '🔥',
  '✨', '❤️', '💔', '🎉', '🎂', '🌈', '⭐', '☀️',
  '🌙', '⚽', '🎮', '🚀', '💯', '✅', '❌', '🍕',
];

// Grid of emoji buttons; clicking one calls onPick with the emoji.
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