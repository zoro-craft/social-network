'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { API_BASE, apiRequest } from '@/utils/api';
import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import ReelItem from './ReelItem';
import ReelComments from './ReelComments';

const PAGE_SIZE = 10;

// Reels: a snap-scrolling column of full-height videos where only the visible reel
// plays. Likes, saves, shares, views and comments all patch the single affected
// reel in place, using the counter bundle those endpoints return.
export default function ReelFeed() {
  const router = useRouter();
  const [mode, setMode] = useState('feed');

  return (
    <>
      <div className="glass-reels-head">
        <h1 className="glass-reels-title">Reels</h1>

        <div className="glass-reels-tabs">
          <button
            type="button"
            className={`glass-reels-tab ${mode === 'feed' ? 'glass-reels-tab-active' : ''}`}
            onClick={() => setMode('feed')}
          >
            For you
          </button>
          <button
            type="button"
            className={`glass-reels-tab ${mode === 'saved' ? 'glass-reels-tab-active' : ''}`}
            onClick={() => setMode('saved')}
          >
            Saved
          </button>
        </div>

        <button
          className="glass-button glass-create-post-button"
          onClick={() => router.push('/reels/create')}
        >
          + Create Reel
        </button>
      </div>

      {/* Remounting on `mode` is what resets the column: a fresh feed starts from a
          clean slate, so no effect has to clear state after the fact. */}
      <ReelColumn key={mode} mode={mode} />
    </>
  );
}

function ReelColumn({ mode }) {
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();

  const [reels, setReels] = useState([]);
  const [loading, setLoading] = useState(true);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [activeId, setActiveId] = useState(null);
  const [muted, setMuted] = useState(true);
  const [commentsFor, setCommentsFor] = useState(null);

  const scrollRef = useRef(null);
  const itemRefs = useRef(new Map());
  const viewedRef = useRef(new Set());
  const snapshotMaxId = useRef(null);
  const loadingRef = useRef(false);

  const endpoint = mode === 'saved' ? '/reels/saved' : '/reels';

  useEffect(() => {
    if (!auth.authenticated) return;

    let cancelled = false;

    const load = async () => {
      try {
        const res = await fetch(`${API_BASE}${endpoint}?limit=${PAGE_SIZE}&offset=0`, {
          credentials: 'include',
        });
        const data = await res.json();

        if (cancelled) return;

        if (!res.ok || !Array.isArray(data)) {
          setReels([]);
          return;
        }

        // Pin the first page's ceiling so later pages ignore reels posted mid-scroll.
        snapshotMaxId.current = res.headers.get('X-Max-Reel-Id');
        setReels(data);
        setHasMore(data.length === PAGE_SIZE);
      } catch {
        if (!cancelled) {
          setReels([]);
          setHasMore(false);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    load();

    return () => {
      cancelled = true;
    };
  }, [auth.authenticated, endpoint]);

  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMore) return;

    loadingRef.current = true;
    setLoadingMore(true);

    try {
      const params = new URLSearchParams({
        limit: String(PAGE_SIZE),
        offset: String(reels.length),
      });
      if (mode !== 'saved' && snapshotMaxId.current) {
        params.set('before', snapshotMaxId.current);
      }

      const res = await fetch(`${API_BASE}${endpoint}?${params.toString()}`, {
        credentials: 'include',
      });
      const data = await res.json();

      if (!res.ok || !Array.isArray(data)) {
        setHasMore(false);
        return;
      }

      setReels((prev) => {
        const seen = new Set(prev.map((reel) => reel.id));
        return [...prev, ...data.filter((reel) => !seen.has(reel.id))];
      });
      setHasMore(data.length === PAGE_SIZE);
    } catch {
      setHasMore(false);
    } finally {
      loadingRef.current = false;
      setLoadingMore(false);
    }
  }, [endpoint, hasMore, mode, reels.length]);

  // Watch the slide column: whichever reel is at least 60% visible becomes the
  // active one (which starts its playback) and, near the end, prefetches more.
  useEffect(() => {
    const root = scrollRef.current;
    if (!root) return;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;

          const id = Number(entry.target.dataset.reelId);
          if (!Number.isFinite(id)) return;

          setActiveId(id);

          const index = reels.findIndex((reel) => reel.id === id);
          if (index >= 0 && index >= reels.length - 2) loadMore();
        });
      },
      { root, threshold: 0.6 },
    );

    // Drop slides a refetch removed, then observe whatever is on screen now.
    const ids = new Set(reels.map((reel) => reel.id));
    itemRefs.current.forEach((_, id) => {
      if (!ids.has(id)) itemRefs.current.delete(id);
    });
    itemRefs.current.forEach((node) => observer.observe(node));

    return () => observer.disconnect();
  }, [reels, loadMore]);

  const applyStats = useCallback((reelId, stats) => {
    setReels((prev) =>
      prev.map((reel) =>
        reel.id === reelId
          ? {
              ...reel,
              like_count: stats.like_count,
              comment_count: stats.comment_count,
              view_count: stats.view_count,
              save_count: stats.save_count,
              share_count: stats.share_count,
              liked: stats.liked,
              saved: stats.saved,
            }
          : reel,
      ),
    );
  }, []);

  // A view is counted once per reel per session, when it first comes into view.
  // The apiRequest promise makes the setState land after an await, not synchronously
  // inside the effect.
  useEffect(() => {
    if (!activeId || viewedRef.current.has(activeId)) return;

    viewedRef.current.add(activeId);

    apiRequest(`/reels/${activeId}/view`, 'POST')
      .then((stats) => applyStats(activeId, stats))
      .catch(() => {
        // A failed view count must never interrupt playback.
      });
  }, [activeId, applyStats]);

  const act = async (reel, action) => {
    try {
      const stats = await apiRequest(`/reels/${reel.id}/${action}`, 'POST');
      applyStats(reel.id, stats);
      return stats;
    } catch (err) {
      if (err?.status === 401) {
        router.push('/login');
        return null;
      }
      showToast(err?.message || 'Could not update reel');
      return null;
    }
  };

  const share = async (reel) => {
    const stats = await act(reel, 'share');
    if (!stats) return;

    const url = `${window.location.origin}/reels?reel=${reel.id}`;
    try {
      if (navigator.share) {
        await navigator.share({ title: 'Reel', text: reel.caption || '', url });
        return;
      }
      await navigator.clipboard.writeText(url);
      showToast('Link copied', false);
    } catch {
      // Dismissing the native share sheet is not an error worth reporting.
    }
  };

  const remove = async (reel) => {
    if (!window.confirm('Delete this reel?')) return;

    try {
      await apiRequest(`/reels/${reel.id}`, 'DELETE');
      setReels((prev) => prev.filter((item) => item.id !== reel.id));
      showToast('Reel deleted', false);
    } catch (err) {
      showToast(err?.message || 'Could not delete reel');
    }
  };

  // One stable callback for every slide: the id rides on the node's data attribute,
  // so refs are never pointlessly detached and re-attached between renders.
  const registerItem = useCallback((node) => {
    if (!node) return;
    const id = Number(node.dataset.reelId);
    if (Number.isFinite(id)) itemRefs.current.set(id, node);
  }, []);

  const commentsReel = reels.find((reel) => reel.id === commentsFor) || null;

  return (
    <>
      <div className="glass-reels-scroll" ref={scrollRef}>
        {loading && (
          <div className="glass-reels-empty">
            <p className="glass-muted-text">Loading reels...</p>
          </div>
        )}

        {!loading && reels.length === 0 && (
          <div className="glass-reels-empty">
            <p className="glass-muted-text">
              {mode === 'saved'
                ? 'No saved reels yet. Tap the bookmark on a reel to keep it here.'
                : 'No reels yet. Be the first to post one!'}
            </p>
          </div>
        )}

        {reels.map((reel) => (
          <ReelItem
            key={reel.id}
            reel={reel}
            active={reel.id === activeId}
            muted={muted}
            itemRef={registerItem}
            isOwn={reel.user_id === auth.id}
            onToggleMute={() => setMuted((prev) => !prev)}
            onLike={(item) => act(item, 'like')}
            onSave={(item) => act(item, 'save')}
            onShare={share}
            onOpenComments={(item) => setCommentsFor(item.id)}
            onDelete={remove}
          />
        ))}

        {loadingMore && <p className="glass-muted-text glass-empty-state">Loading more...</p>}
      </div>

      {commentsReel && (
        <ReelComments
          key={commentsReel.id}
          reel={commentsReel}
          onClose={() => setCommentsFor(null)}
          onCommentAdded={(reelId) =>
            setReels((prev) =>
              prev.map((reel) =>
                reel.id === reelId ? { ...reel, comment_count: reel.comment_count + 1 } : reel,
              ),
            )
          }
        />
      )}
    </>
  );
}
