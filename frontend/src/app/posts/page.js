'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { API_BASE } from '@/utils/api';
import StoryBar from '@/app/components/StoryBar';
import Avatar from '@/app/components/Avatar';
import Icon from '@/app/components/Icon';
import { useReactions } from './useReactions';
import formatDate from '@/utils/format';

const PAGE_SIZE = 10;

const PRIVACY_META = {
  public: { label: 'Public', icon: 'eye' },
  almost_private: { label: 'Followers', icon: 'users' },
  private: { label: 'Private', icon: 'lock' },
};

function PrivacyBadge({ post }) {
  if (post.group_id) return null;

  const meta = PRIVACY_META[post.privacy];
  if (!meta) return null;

  return (
    <span className={`glass-post-privacy glass-post-privacy-${post.privacy}`}>
      <Icon name={meta.icon} size={12} />
      {meta.label}
    </span>
  );
}

function PostSkeleton() {
  return (
    <div className="glass-post-card glass-post-skeleton" aria-hidden="true">
      <div className="glass-post-header">
        <div className="glass-skeleton-avatar" />
        <div className="glass-skeleton-lines">
          <span className="glass-skeleton-line glass-skeleton-line--short" />
          <span className="glass-skeleton-line glass-skeleton-line--tiny" />
        </div>
      </div>
      <span className="glass-skeleton-line" />
      <span className="glass-skeleton-line glass-skeleton-line--wide" />
    </div>
  );
}

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
      <header className="glass-feed-header">
        <div>
          <h1 className="glass-feed-title">Feed</h1>
          <p className="glass-feed-subtitle">See what everyone is sharing</p>
        </div>
        <button
          type="button"
          className="glass-button glass-create-post-button"
          onClick={() => router.push('/posts/create')}
        >
          <Icon name="share" size={16} />
          Create Post
        </button>
      </header>

      <StoryBar />

      {postsLoading && (
        <div className="glass-feed-list">
          {[0, 1, 2].map((i) => (
            <PostSkeleton key={i} />
          ))}
        </div>
      )}

      {!postsLoading && posts.length === 0 && (
        <div className="glass-empty-state glass-feed-empty">
          <Icon name="message" size={26} />
          <p>No posts yet. Be the first to share something.</p>
          <button
            type="button"
            className="glass-button glass-button-inline"
            onClick={() => router.push('/posts/create')}
          >
            Create your first post
          </button>
        </div>
      )}

      {!postsLoading && posts.length > 0 && (
        <div className="glass-feed-list">
          {posts.map((post) => {
            const liked = post.user_reaction === 'like';
            const disliked = post.user_reaction === 'dislike';

            return (
              <article key={post.id} className="glass-post-card">
                <header className="glass-post-header">
                  <div className="glass-post-author-row">
                    <Avatar
                      src={post.avatar}
                      name={post.full_name}
                      size={40}
                    />
                    <div className="glass-post-author-meta">
                      <span className="glass-post-author">
                        {post.full_name}
                      </span>
                      <span className="glass-post-time">
                        {formatDate(post.created_at)}
                      </span>
                    </div>
                  </div>
                  <PrivacyBadge post={post} />
                </header>

                <h2 className="glass-post-title">{post.title}</h2>

                {post.image && (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    className="glass-post-image"
                    src={`${API_BASE}${post.image}`}
                    alt=""
                  />
                )}

                <p className="glass-post-content">{post.content}</p>

                <footer className="glass-reactions">
                  <button
                    type="button"
                    className={`glass-react-btn ${liked ? 'glass-react-like' : ''}`}
                    aria-pressed={liked}
                    onClick={() => react(post.id, 'like')}
                  >
                    <Icon name="thumbUp" filled={liked} />
                    <span className="glass-react-count">{post.like_count || 0}</span>
                  </button>
                  <button
                    type="button"
                    className={`glass-react-btn ${disliked ? 'glass-react-dislike' : ''}`}
                    aria-pressed={disliked}
                    onClick={() => react(post.id, 'dislike')}
                  >
                    <Icon name="thumbDown" filled={disliked} />
                    <span className="glass-react-count">
                      {post.dislike_count || 0}
                    </span>
                  </button>
                  <button
                    type="button"
                    className="glass-react-btn glass-react-comment"
                    onClick={() => router.push(`/posts/comment/${post.id}`)}
                  >
                    <Icon name="message" />
                    <span>Comment</span>
                  </button>
                </footer>
              </article>
            );
          })}
        </div>
      )}

      {!postsLoading && hasMore && posts.length > 0 && (
        <button
          type="button"
          className="glass-button glass-load-more"
          onClick={loadMore}
          disabled={loadingMore}
        >
          {loadingMore ? 'Loading...' : 'Load more'}
        </button>
      )}
    </div>
  );
}
