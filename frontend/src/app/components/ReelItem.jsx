'use client';

import { useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { API_BASE } from '@/utils/api';
import formatDate from '@/utils/format';
import Avatar from './Avatar';
import Icon from './Icon';

// One full-height reel slide. Playback is driven by the parent's `active` flag
// (from an IntersectionObserver) so only the reel in view is ever playing.
export default function ReelItem({
  reel,
  active,
  muted,
  itemRef,
  isOwn,
  onToggleMute,
  onLike,
  onSave,
  onShare,
  onOpenComments,
  onDelete,
}) {
  const router = useRouter();
  const videoRef = useRef(null);

  // Muted and playback live in one effect on purpose: browsers refuse unmuted
  // autoplay, so `muted` has to be applied to the element *before* the first
  // play() rather than in a separate effect that may run afterwards.
  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    video.muted = muted;

    if (active) {
      // Swallow the rejection: a blocked autoplay must not break scrolling.
      const started = video.play();
      if (started && typeof started.catch === 'function') started.catch(() => {});
    } else {
      video.pause();
    }
  }, [active, muted]);

  return (
    <div className="glass-reel" ref={itemRef} data-reel-id={reel.id}>
      <video
        ref={videoRef}
        className="glass-reel-video"
        src={`${API_BASE}${reel.video}`}
        playsInline
        loop
        preload="metadata"
      />

      <button
        type="button"
        className="glass-reel-mute"
        onClick={onToggleMute}
        aria-label={muted ? 'Unmute reel' : 'Mute reel'}
      >
        <Icon name={muted ? 'volumeOff' : 'volume'} />
      </button>

      <div className="glass-reel-overlay">
        <div className="glass-reel-info">
          <div className="glass-reel-author">
            <button
              type="button"
              className="glass-reel-author-name"
              onClick={() => router.push(`/profile/${reel.user_id}`)}
            >
              <Avatar src={reel.avatar} name={reel.full_name} size={32} />
              <span>{reel.full_name}</span>
            </button>
            {isOwn && (
              <button
                type="button"
                className="glass-story-icon-btn"
                onClick={onDelete}
                aria-label="Delete reel"
              >
                <Icon name="trash" />
              </button>
            )}
          </div>

          {reel.caption && <p className="glass-reel-caption">{reel.caption}</p>}

          <span className="glass-reel-time">
            {formatDate(reel.created_at)} · <Icon name="eye" /> {reel.view_count}
          </span>
        </div>

        <div className="glass-reel-actions">
          <button
            type="button"
            className={`glass-reel-action ${reel.liked ? 'glass-reel-action-liked' : ''}`}
            onClick={() => onLike(reel)}
            aria-label="Like reel"
          >
            <span>
              <Icon name="heart" filled={reel.liked} />
            </span>
            <span className="glass-reel-action-count">{reel.like_count}</span>
          </button>

          <button
            type="button"
            className="glass-reel-action"
            onClick={() => onOpenComments(reel)}
            aria-label="Open comments"
          >
            <span>
              <Icon name="message" />
            </span>
            <span className="glass-reel-action-count">{reel.comment_count}</span>
          </button>

          <button
            type="button"
            className="glass-reel-action"
            onClick={() => onShare(reel)}
            aria-label="Share reel"
          >
            <span>
              <Icon name="share" />
            </span>
            <span className="glass-reel-action-count">{reel.share_count}</span>
          </button>

          <button
            type="button"
            className={`glass-reel-action ${reel.saved ? 'glass-reel-action-saved' : ''}`}
            onClick={() => onSave(reel)}
            aria-label={reel.saved ? 'Remove from saved' : 'Save reel'}
          >
            <span>
              <Icon name="bookmark" filled={reel.saved} />
            </span>
            <span className="glass-reel-action-count">{reel.save_count}</span>
          </button>
        </div>
      </div>
    </div>
  );
}
