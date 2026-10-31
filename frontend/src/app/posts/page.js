'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { API_BASE } from '@/utils/api';
import PostMedia from '@/components/PostMedia';
import { useReactions } from './useReactions';
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
<div className="glass-feed-header">
				<h1 className="glass-feed-title">Feed</h1>
				<button
					className={`glass-button glass-create-post-button`}
					onClick={() => router.push('/posts/create')}
				>
					+ Create Post
				</button>
			</div>

			{postsLoading && <p className="glass-muted-text">Loading posts...</p>}

			{!postsLoading && posts.length === 0 && (
				<p className="glass-muted-text">
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

          <PostMedia src={post.image} />

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
          className={`glass-button glass-load-more`}
          onClick={loadMore}
          disabled={loadingMore}
        >
          {loadingMore ? 'Loading...' : 'Load more'}
        </button>
      )}
    </div>
  );
}
