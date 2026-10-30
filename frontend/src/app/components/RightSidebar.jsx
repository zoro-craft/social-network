'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Avatar from './Avatar';
import Icon from './Icon';

// Presentation-only placeholders: there is no suggestions/trending endpoint
// yet, so these are static mock rows. Following a user is local UI state and
// intentionally does not touch the API.
const SUGGESTED = [
  {
    id: 1,
    name: 'Maya Chen',
    handle: '@mayachen',
    reason: 'Followed by Amina',
  },
  {
    id: 2,
    name: 'Youssef B.',
    handle: '@youssef',
    reason: 'New to your network',
  },
  { id: 3, name: 'Lina Park', handle: '@linapark', reason: 'Popular this week' },
];

const TRENDING = [
  { tag: '#designsystems', meta: '2.4k posts' },
  { tag: '#buildinpublic', meta: '1.8k posts' },
  { tag: '#uianimation', meta: '940 posts' },
];

export default function RightSidebar() {
  const router = useRouter();
  const [following, setFollowing] = useState([]);

  const toggle = (id) =>
    setFollowing((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id],
    );

  return (
    <aside
      className="glass-shell-side glass-rail-right"
      aria-label="Suggested people and trending topics"
    >
      <div className="glass-side-card">
        <div className="glass-side-head">
          <h2 className="glass-side-title">Suggested for you</h2>
        </div>

        <ul className="glass-side-list">
          {SUGGESTED.map((person) => {
            const isFollowing = following.includes(person.id);

            return (
              <li key={person.id} className="glass-suggest-row">
                <button
                  type="button"
                  className="glass-side-row glass-suggest-person"
                  onClick={() => router.push('/search')}
                >
                  <Avatar name={person.name} size={40} />
                  <span className="glass-side-row-body">
                    <span className="glass-side-row-title">{person.name}</span>
                    <span className="glass-side-row-meta">
                      {person.handle} · {person.reason}
                    </span>
                  </span>
                </button>
                <button
                  type="button"
                  className={`glass-follow-btn ${isFollowing ? 'glass-follow-btn-on' : ''}`}
                  onClick={() => toggle(person.id)}
                  aria-pressed={isFollowing}
                >
                  {isFollowing ? 'Following' : 'Follow'}
                </button>
              </li>
            );
          })}
        </ul>
      </div>

      <div className="glass-side-card">
        <div className="glass-side-head">
          <h2 className="glass-side-title">Trending</h2>
          <Icon name="share" size={14} />
        </div>

        <ul className="glass-side-list">
          {TRENDING.map((topic, index) => (
            <li key={topic.tag}>
              <button type="button" className="glass-trend-row">
                <span className="glass-trend-rank">{index + 1}</span>
                <span className="glass-side-row-body">
                  <span className="glass-side-row-title">{topic.tag}</span>
                  <span className="glass-side-row-meta">{topic.meta}</span>
                </span>
              </button>
            </li>
          ))}
        </ul>
      </div>
    </aside>
  );
}
