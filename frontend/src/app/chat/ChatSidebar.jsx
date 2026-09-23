'use client';

import { useState } from 'react';

function formatTime(ts) {
  return new Date(ts).toLocaleString();
}

export default function ChatSidebar({ chat }) {
  const [tab, setTab] = useState('users');
  const { users, groups, selected, badges, selectConversation, convKey } = chat;

  const isSelected = (type, id) => selected && selected.type === type && selected.id === id;

  const renderRow = (entry, type, id, name) => {
    const key = convKey(type, id);
    const badgeCount = badges.counts[key] || 0;
    return (
      <button
        key={key}
        className={`glass-chat-conv-row ${isSelected(type, id) ? 'glass-chat-conv-active' : ''}`}
        onClick={() => selectConversation(type, id, name)}
      >
        <div className="glass-chat-conv-info">
          <div className="glass-chat-conv-title">
            {name}
            {type === 'user' && (
              <span className={`glass-chat-online ${entry.online ? 'glass-chat-online-true' : ''}`} />
            )}
          </div>
          <div className={`glass-chat-conv-preview ${badgeCount > 0 ? 'glass-chat-conv-preview-unread' : ''}`}>
            {entry.last_message || 'No messages yet'}
          </div>
        </div>
        <div className="glass-chat-conv-meta">
          {entry.last_message_at ? <small>{formatTime(entry.last_message_at)}</small> : null}
          {badgeCount > 0 && <span className="glass-chat-row-badge">{badgeCount}</span>}
        </div>
      </button>
    );
  };

  const list =
    tab === 'users'
      ? users.map((u) => renderRow(u, 'user', u.user_id, u.name))
      : groups.map((g) => renderRow(g, 'group', g.group_id, g.title));

  return (
    <aside className="glass-chat-sidebar">
      <div className="glass-chat-tabs">
        <button
          className={tab === 'users' ? 'glass-chat-tab glass-chat-tab-active' : 'glass-chat-tab'}
          onClick={() => setTab('users')}
        >
          Users
          {badges.users > 0 && <span className="glass-chat-tab-badge">{badges.users}</span>}
        </button>
        <button
          className={tab === 'groups' ? 'glass-chat-tab glass-chat-tab-active' : 'glass-chat-tab'}
          onClick={() => setTab('groups')}
        >
          Groups
          {badges.groups > 0 && <span className="glass-chat-tab-badge">{badges.groups}</span>}
        </button>
      </div>

      <div className="glass-chat-conv-list">
        {list.length === 0 ? (
          <p className="glass-chat-empty-list">
            {tab === 'users' ? 'No conversations yet' : 'No groups yet'}
          </p>
        ) : (
          list
        )}
      </div>
    </aside>
  );
}