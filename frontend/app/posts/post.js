'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../context/AuthContext';
import { API_BASE } from '../lib/api';
import { useReactions } from './useReactions';
import styles from '../styles/home.module.css';
import formatDate from '../lib/format';

export default function Posts() {
  const router = useRouter();
  const { auth, loading } = useAuth();
  const [posts, setPosts] = useState([]);
  const [postsLoading, setPostsLoading] = useState(true);
  const { react } = useReactions(setPosts);

  useEffect(() => {
    if (!loading && !auth.authenticated) {
      router.push('/login');
    }
  }, [loading, auth.authenticated, router]);

  useEffect(() => {
    if (!auth.authenticated) return;

    async function fetchPosts() {
      try {
        const res = await fetch(`${API_BASE}/posts`, {
          credentials: 'include',
        });
        const data = await res.json();
        setPosts(data || []);
      } catch (err) {
        setPosts([]);
      } finally {
        setPostsLoading(false);
      }
    }

    fetchPosts();
  }, [auth.authenticated]);

  if (loading || !auth.authenticated) {
    return null;
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.heading}>Feed</h1>
        <button
          className={styles.createButton}
          onClick={() => router.push('/posts/create')}
        >
          + Create Post
        </button>
      </div>

      {postsLoading && <p className={styles.status}>Loading posts...</p>}

      {!postsLoading && posts.length === 0 && (
        <p className={styles.status}>
          No posts yet. Be the first to share something.
        </p>
      )}

      {posts.map((post) => (
        <div key={post.id} className={styles.postCard}>
          <div className={styles.postHeader}>
            <p className={styles.postAuthor}>Posted by {post.nickname}</p>
            <p className={styles.postDate}>{formatDate(post.created_at)}</p>
          </div>
          <p className={styles.postTitle}>{post.title}</p>

          {post.image && (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              className={styles.postImage}
              src={`${API_BASE}${post.image}`}
              alt=""
            />
          )}

          <p className={styles.postContent}>{post.content}</p>

          <div className={styles.reactions}>
            <button
              className={`${styles.reactButton} ${
                post.user_reaction === 'like' ? styles.likeActive : ''
              }`}
              onClick={() => react(post.id, 'like')}
            >
              👍 {post.like_count || 0}
            </button>
            <button
              className={`${styles.reactButton} ${
                post.user_reaction === 'dislike' ? styles.dislikeActive : ''
              }`}
              onClick={() => react(post.id, 'dislike')}
            >
              👎 {post.dislike_count || 0}
            </button>
            <button
              className={styles.reactButton}
              onClick={() => router.push(`/posts/comment/${post.id}`)}
            >
              💬 Comment
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
