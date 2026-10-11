'use client';

import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@/utils/api';
import { useAuth } from '@/app/context/AuthContext';
import Avatar from './Avatar';
import StoryViewer from './StoryViewer';

async function fetchStoryGroups() {
  const data = await apiRequest('/stories');
  return Array.isArray(data) ? data : [];
}

// The horizontal "stories" rail above the feed, Instagram-style: your own bubble
// first, then everyone you follow who has an active story.
export default function StoryBar() {
  const router = useRouter();
  const { auth } = useAuth();

  const [groups, setGroups] = useState([]);
  const [activeIndex, setActiveIndex] = useState(null);

  // Used after a story is deleted, when we need to refetch on demand.
  const reload = useCallback(async () => {
    try {
      setGroups(await fetchStoryGroups());
    } catch {
      setGroups([]);
    }
  }, []);

  useEffect(() => {
    if (!auth.authenticated) return;

    let cancelled = false;

    const load = async () => {
      try {
        const data = await fetchStoryGroups();
        if (!cancelled) setGroups(data);
      } catch {
        if (!cancelled) setGroups([]);
      }
    };

    load();
    return () => {
      cancelled = true;
    };
  }, [auth.authenticated]);

  const ownGroup = groups.find((group) => group.user_id === auth.id);
  const others = groups.filter((group) => group.user_id !== auth.id);
  const rail = ownGroup ? [ownGroup, ...others] : others;

  const openGroup = (group) => {
    const index = rail.findIndex((item) => item.user_id === group.user_id);
    if (index >= 0) setActiveIndex(index);
  };

  const openOwn = () => {
    if (ownGroup) {
      openGroup(ownGroup);
    } else {
      router.push('/stories/create');
    }
  };

  return (
    <div className="glass-stories-bar">
      <div className="glass-stories-scroll">
        <div className="glass-story-tile-wrap">
          <button type="button" className="glass-story-tile" onClick={openOwn}>
            <span
              className={`glass-story-ring ${
                ownGroup && !ownGroup.seen ? 'glass-story-ring-unseen' : ''
              }`}
            >
              {/* The rail sits above the fold on the feed, so it is the page's
                  LCP image — load it eagerly instead of lazily. */}
              <Avatar src={auth.avatar} name={auth.first_name} size={56} eager />
            </span>
            <span className="glass-story-label">Your story</span>
          </button>
          {/* Sibling of the tile button, so the two clicks never overlap. */}
          <button
            type="button"
            className="glass-story-add"
            onClick={() => router.push('/stories/create')}
            aria-label="Add to story"
          >
            +
          </button>
        </div>

        {others.map((group) => (
          <button
            key={group.user_id}
            type="button"
            className="glass-story-tile"
            onClick={() => openGroup(group)}
          >
            <span
              className={`glass-story-ring ${
                group.seen ? '' : 'glass-story-ring-unseen'
              }`}
            >
              <Avatar src={group.avatar} name={group.full_name} size={56} />
            </span>
            <span className="glass-story-label">{group.full_name.split(' ')[0]}</span>
          </button>
        ))}
      </div>

      <StoryViewer
        key={activeIndex ?? 'closed'}
        groups={rail}
        startIndex={activeIndex}
        viewerId={auth.id}
        onClose={() => setActiveIndex(null)}
        onChanged={reload}
      />
    </div>
  );
}
