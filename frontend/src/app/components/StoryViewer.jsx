'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { API_BASE, apiRequest } from '@/utils/api';
import formatDate from '@/utils/format';
import Avatar from './Avatar';
import Icon from './Icon';

// How long a single story stays on screen before auto-advancing.
const STORY_DURATION = 5000;

// The rail remounts this component (via `key`) whenever a different story is
// opened, so the position only ever has to be derived from `startIndex` once.
export default function StoryViewer({ groups, startIndex, viewerId, onClose, onChanged }) {
  const [position, setPosition] = useState({ group: startIndex ?? 0, story: 0 });
  // Progress is tagged with the story it belongs to, so switching stories reads
  // as 0% immediately without a state reset inside an effect.
  const [tick, setTick] = useState({ storyId: null, value: 0 });

  const open = startIndex !== null;
  const group = groups[position.group];
  const story = group?.stories[position.story];

  const next = useCallback(() => {
    const current = groups[position.group];

    if (!current) {
      onClose();
      return;
    }

    if (position.story + 1 < current.stories.length) {
      setPosition({ group: position.group, story: position.story + 1 });
      return;
    }

    if (position.group + 1 < groups.length) {
      setPosition({ group: position.group + 1, story: 0 });
      return;
    }

    onClose();
  }, [groups, position, onClose]);

  const previous = useCallback(() => {
    if (position.story > 0) {
      setPosition({ group: position.group, story: position.story - 1 });
      return;
    }

    if (position.group > 0) {
      const earlier = groups[position.group - 1];
      setPosition({
        group: position.group - 1,
        story: Math.max(0, (earlier?.stories.length || 1) - 1),
      });
    }
  }, [groups, position]);

  // The timer and key handlers need the latest callbacks without restarting.
  const nextRef = useRef(next);
  const previousRef = useRef(previous);
  useEffect(() => {
    nextRef.current = next;
    previousRef.current = previous;
  }, [next, previous]);

  const storyId = story?.id;

  useEffect(() => {
    if (!open || storyId === undefined) return;

    const startedAt = Date.now();
    const timer = setInterval(() => {
      const value = Math.min(1, (Date.now() - startedAt) / STORY_DURATION);
      setTick({ storyId, value });
      if (value >= 1) {
        clearInterval(timer);
        nextRef.current();
      }
    }, 50);

    // Let the author know who saw their story; a failure must not block playback.
    apiRequest(`/stories/${storyId}/view`, 'POST').catch(() => {});

    return () => clearInterval(timer);
  }, [open, storyId]);

  useEffect(() => {
    if (!open) return;

    const onKeyDown = (event) => {
      if (event.key === 'Escape') onClose();
      if (event.key === 'ArrowRight') nextRef.current();
      if (event.key === 'ArrowLeft') previousRef.current();
    };

    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [open, onClose]);

  if (!open || !story) return null;

  const progress = tick.storyId === storyId ? tick.value : 0;
  const isOwn = story.user_id === viewerId;

  const remove = async () => {
    if (!isOwn) return;

    const remaining = group.stories.length - 1;

    try {
      await apiRequest(`/stories/${story.id}`, 'DELETE');
    } catch {
      return;
    }

    if (onChanged) onChanged();

    if (remaining > 0) {
      // The refetch shifts the next story into this index.
      setPosition((prev) => ({
        group: prev.group,
        story: Math.min(prev.story, remaining - 1),
      }));
    } else {
      onClose();
    }
  };

  return (
    <div className="glass-story-viewer" onClick={onClose}>
      <div className="glass-story-viewer-inner" onClick={(event) => event.stopPropagation()}>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          className="glass-story-media"
          src={`${API_BASE}${story.image}`}
          alt={group.full_name}
        />

        <div className="glass-story-bars">
          {group.stories.map((item, index) => (
            <div key={item.id} className="glass-story-bar">
              <span
                className="glass-story-bar-fill"
                style={{
                  width:
                    index < position.story
                      ? '100%'
                      : index === position.story
                        ? `${progress * 100}%`
                        : '0%',
                }}
              />
            </div>
          ))}
        </div>

        <div className="glass-story-head">
          <Avatar src={group.avatar} name={group.full_name} size={34} />
          <div className="glass-story-head-text">
            <span className="glass-story-name">{group.full_name}</span>
            <span className="glass-story-time">{formatDate(story.created_at)}</span>
          </div>
          {isOwn && (
            <button
              type="button"
              className="glass-story-icon-btn"
              onClick={remove}
              aria-label="Delete story"
            >
              <Icon name="trash" />
            </button>
          )}
          <button
            type="button"
            className="glass-story-icon-btn"
            onClick={onClose}
            aria-label="Close story"
          >
            <Icon name="close" />
          </button>
        </div>

        {/* Tap zones: the left third steps back, the rest advances. */}
        <div className="glass-story-nav">
          <button
            type="button"
            className="glass-story-nav-prev"
            onClick={previous}
            aria-label="Previous story"
          />
          <button
            type="button"
            className="glass-story-nav-next"
            onClick={() => nextRef.current()}
            aria-label="Next story"
          />
        </div>

        {story.caption && <p className="glass-story-caption">{story.caption}</p>}
      </div>
    </div>
  );
}
