'use client';

import { useEffect, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';
import { API_BASE, apiRequest } from '@/utils/api';
import Avatar from './Avatar';

// Left gutter of the app shell: the groups you belong to, and who is online now.
export default function LeftSidebar() {
  const router = useRouter();
  const pathname = usePathname();
  const { auth } = useAuth();

  const [groups, setGroups] = useState([]);
  const [online, setOnline] = useState([]);

  useEffect(() => {
    if (!auth.authenticated) return;

    let cancelled = false;

    // Refetched on navigation so a group joined elsewhere, or a fresh sign-in,
    // shows up right away.
    apiRequest('/groups')
      .then((data) => {
        if (cancelled) return;
        const mine = (Array.isArray(data) ? data : []).filter((group) => group.is_member);
        setGroups(mine);
      })
      .catch(() => {
        if (!cancelled) setGroups([]);
      });

    apiRequest('/users/online')
      .then((data) => {
        if (!cancelled) setOnline(Array.isArray(data) ? data : []);
      })
      .catch(() => {
        if (!cancelled) setOnline([]);
      });

    return () => {
      cancelled = true;
    };
  }, [auth.authenticated, pathname]);

  if (!auth.authenticated) return null;

  return (
    <aside
      className="glass-shell-side"
      aria-label="Your groups and people online"
    >
      <div className="glass-side-card">
        <div className="glass-side-head">
          <h2 className="glass-side-title">Your groups</h2>
          <button
            type="button"
            className="glass-side-action"
            onClick={() => router.push('/groups')}
          >
            See all
          </button>
        </div>

        {groups.length === 0 ? (
          <p className="glass-side-empty">You have not joined any groups yet.</p>
        ) : (
          <ul className="glass-side-list">
            {groups.map((group) => (
              <li key={group.id}>
                <button
                  type="button"
                  className="glass-side-row"
                  onClick={() => router.push(`/groups/${group.id}`)}
                >
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    className="glass-side-thumb"
                    src={group.image ? `${API_BASE}${group.image}` : '/group-avatar.jpeg'}
                    alt=""
                    onError={(event) => {
                      event.currentTarget.onerror = null;
                      event.currentTarget.src = '/group-avatar.jpeg';
                    }}
                  />
                  <span className="glass-side-row-body">
                    <span className="glass-side-row-title">{group.titel}</span>
                    {group.description ? (
                      <span className="glass-side-row-meta">{group.description}</span>
                    ) : null}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        )}

        <button
          type="button"
          className="glass-side-cta"
          onClick={() => router.push('/groups/create')}
        >
          + Create a group
        </button>
      </div>

      <div className="glass-side-card">
        <div className="glass-side-head">
          <h2 className="glass-side-title">Online now</h2>
          <span className="glass-side-count">{online.length}</span>
        </div>

        {online.length === 0 ? (
          <p className="glass-side-empty">Nobody else is online right now.</p>
        ) : (
          <ul className="glass-side-list">
            {online.map((person) => (
              <li key={person.id}>
                <button
                  type="button"
                  className="glass-side-row"
                  onClick={() => router.push(`/profile/${person.id}`)}
                >
                  <span className="glass-side-avatar">
                    <Avatar src={person.avatar} name={person.full_name} size={34} />
                    <span className="glass-side-dot glass-side-dot-online" />
                  </span>
                  <span className="glass-side-row-body">
                    <span className="glass-side-row-title">{person.full_name}</span>
                    <span className="glass-side-row-meta">Online</span>
                  </span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </aside>
  );
}
