'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useToast } from '@/app/context/ToastContext';
import { apiRequest } from '@/utils/api';
import formatDate from '@/utils/format';
import Avatar from './Avatar';
import Icon from './Icon';

const PAGE_SIZE = 20;

// Comment drawer for one reel. It owns its own pagination and posting state so the
// feed only has to say which reel is open and react to the new comment count.
export default function ReelComments({ reel, onClose, onCommentAdded }) {
  const router = useRouter();
  const showToast = useToast();

  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const [content, setContent] = useState('');
  const [submitting, setSubmitting] = useState(false);

  // The parent mounts this per reel (via `key`), so the first page is only ever
  // fetched once per open and `loading` can start at true.
  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      try {
        const data = await apiRequest(`/reels/${reel.id}/comments?limit=${PAGE_SIZE}&offset=0`);
        if (cancelled) return;

        const list = Array.isArray(data) ? data : [];
        setComments(list);
        setHasMore(list.length === PAGE_SIZE);
      } catch {
        if (!cancelled) {
          setComments([]);
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
  }, [reel]);

  useEffect(() => {
    const onKeyDown = (event) => {
      if (event.key === 'Escape') onClose();
    };

    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [onClose]);

  if (!reel) return null;

  // Comments come back newest first, so "older" means a smaller id than the oldest
  // one currently on screen.
  const loadOlder = async () => {
    if (loadingMore || comments.length === 0) return;

    setLoadingMore(true);
    try {
      const oldest = comments[comments.length - 1].id;
      const data = await apiRequest(
        `/reels/${reel.id}/comments?limit=${PAGE_SIZE}&offset=0&before=${oldest}`,
      );
      const list = Array.isArray(data) ? data : [];
      setComments((prev) => [...prev, ...list]);
      setHasMore(list.length === PAGE_SIZE);
    } catch {
      setHasMore(false);
    } finally {
      setLoadingMore(false);
    }
  };

  const submit = async (event) => {
    event.preventDefault();

    const text = content.trim();
    if (!text) return;

    setSubmitting(true);
    try {
      const form = new FormData();
      form.append('content', text);

      const created = await apiRequest(`/reels/${reel.id}/comments`, 'POST', form, true);
      setComments((prev) => [created, ...prev]);
      setContent('');
      // The feed holds the authoritative count, so just tell it to tick up.
      onCommentAdded?.(reel.id);
    } catch (err) {
      if (err?.status === 401) {
        router.push('/login');
        return;
      }
      showToast(err?.message || 'Could not add comment');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="glass-reel-comments-backdrop" onClick={onClose}>
      <div className="glass-reel-comments" onClick={(event) => event.stopPropagation()}>
        <div className="glass-reel-comments-head">
          <h3>Comments ({reel.comment_count})</h3>
          <button
            type="button"
            className="glass-story-icon-btn"
            onClick={onClose}
            aria-label="Close comments"
          >
            <Icon name="close" />
          </button>
        </div>

        <div className="glass-reel-comments-list">
          {loading && <p className="glass-muted-text">Loading comments...</p>}

          {!loading && comments.length === 0 && (
            <p className="glass-muted-text">No comments yet. Be the first!</p>
          )}

          {comments.map((comment) => (
            <div key={comment.id} className="glass-reel-comment">
              <Avatar src={comment.avatar} name={comment.first_name} size={32} />
              <div className="glass-reel-comment-body">
                <span className="glass-reel-comment-author">
                  {comment.first_name || 'Someone'}
                </span>
                <p className="glass-reel-comment-text">{comment.content}</p>
                <span className="glass-reel-comment-time">
                  {formatDate(comment.created_at)}
                </span>
              </div>
            </div>
          ))}

          {!loading && hasMore && (
            <button
              type="button"
              className="glass-button glass-load-more"
              onClick={loadOlder}
              disabled={loadingMore}
            >
              {loadingMore ? 'Loading...' : 'Load older comments'}
            </button>
          )}
        </div>

        <form className="glass-reel-comment-form" onSubmit={submit}>
          <input
            className="glass-input"
            type="text"
            placeholder="Write a comment..."
            value={content}
            onChange={(event) => setContent(event.target.value)}
            maxLength={1000}
            disabled={submitting}
          />
          <button className="glass-button" type="submit" disabled={submitting || !content.trim()}>
            {submitting ? '...' : 'Post'}
          </button>
        </form>
      </div>
    </div>
  );
}
