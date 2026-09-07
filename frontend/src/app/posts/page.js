'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { API_BASE } from '@/utils/api';
import { useReactions } from './useReactions';
import styles from '../styles/posts.module.css';
import formatDate from '@/utils/format';

const PAGE_SIZE = 10;

export default function Posts() {
  const router = useRouter();
  const { auth, loading } = useRequireAuth();
  const [posts, setPosts] = useState([]);
  const [postsLoading, setPostsLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [snapshotMaxId, setSnapshotMaxId] = useState(null);
  const [hasMore, setHasMore] = useState(true);
  const { react } = useReactions(setPosts);

  useEffect(() => {
    if (!auth.authenticated) return;

    async function fetchFirstPage() {
      try {
        const res = await fetch(
          `${API_BASE}/posts?limit=${PAGE_SIZE}&offset=0`,
          {
            credentials: 'include',
          },
        );
        const data = await res.json();
        
        if (!res.ok || !Array.isArray(data)) {
          setPosts([]);
          setHasMore(false);
          setPostsLoading(false);
          return;
        }

        setSnapshotMaxId(res.headers.get('X-Max-Post-Id'));
        setPosts(data);
        setHasMore(data.length === PAGE_SIZE);
      } catch {
        setPosts([]);
        setHasMore(false);
      } finally {
        setPostsLoading(false);
      }
    }

    fetchFirstPage();
  }, [auth.authenticated]);

  const loadMore = async () => {
    if (!snapshotMaxId || loadingMore) return;

    setLoadingMore(true);
    try {
      const res = await fetch(
        `${API_BASE}/posts?before=${snapshotMaxId}&limit=${PAGE_SIZE}&offset=${posts.length}`,
        { credentials: 'include' },
      );
      const data = await res.json();
      
      if (!res.ok || !Array.isArray(data)) {
        setHasMore(false);
        setLoadingMore(false);
        return;
      }

      setPosts((prev) => [...prev, ...data]);
      setHasMore(data.length === PAGE_SIZE);
    } catch {
      setHasMore(false);
    } finally {
      setLoadingMore(false);
    }
  };

  if (loading || !auth.authenticated) {
    return null;
  }

  return (
    <div className="glass-feed-page">
      <div className={styles.feedHeader}>
        <h1 className={styles.feedTitle}>Feed</h1>
        <button
          className={`glass-button ${styles.createPostButton}`}
          onClick={() => router.push('/posts/create')}
        >
          + Create Post
        </button>
      </div>

      {postsLoading && <p className={styles.mutedText}>Loading posts...</p>}

      {!postsLoading && posts.length === 0 && (
        <p className={styles.mutedText}>
          No posts yet. Be the first to share something.
        </p>
      )}

      {posts.map((post) => (
        <div key={post.id} className="glass-post-card">
          <div className="glass-post-header">
            <span className="glass-post-author">Posted by {post.full_name}</span>
            <span>{formatDate(post.created_at)}</span>
          </div>
          <div className="glass-post-title">{post.title}</div>

          {post.image && (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              className="glass-post-image"
              src={`${API_BASE}${post.image}`}
              alt=""
            />
          )}

          <div className="glass-post-content">{post.content}</div>

          <div className="glass-reactions">
            <button
              className={`glass-react-btn ${post.user_reaction === 'like' ? 'glass-react-like' : ''}`}
              onClick={() => react(post.id, 'like')}
            >
              👍 {post.like_count || 0}
            </button>
            <button
              className={`glass-react-btn ${post.user_reaction === 'dislike' ? 'glass-react-dislike' : ''}`}
              onClick={() => react(post.id, 'dislike')}
            >
              👎 {post.dislike_count || 0}
            </button>
            <button
              className="glass-react-btn"
              onClick={() => router.push(`/posts/comment/${post.id}`)}
            >
              💬 Comment
            </button>
          </div>
        </div>
      ))}

      {!postsLoading && hasMore && (
        <button
          className={`glass-button ${styles.loadMoreButton}`}
          onClick={loadMore}
          disabled={loadingMore}
        >
          {loadingMore ? 'Loading...' : 'Load more'}
        </button>
      )}
    </div>
  );
}
