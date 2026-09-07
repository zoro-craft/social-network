'use client';

import { useState } from 'react';
import { API_BASE } from '@/utils/api';

const PAGE_SIZE = 10;

export function useComments(postId) {
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [snapshotMaxId, setSnapshotMaxId] = useState(null);
  const [hasMore, setHasMore] = useState(true);

  const fetchPostAndComments = async () => {
    setLoading(true);
    try {
      const [postRes, commentsRes] = await Promise.all([
        fetch(`${API_BASE}/posts/${postId}`, { credentials: 'include' }),
        fetch(
          `${API_BASE}/posts/${postId}/comment?limit=${PAGE_SIZE}&offset=0`,
          {
            credentials: 'include',
          },
        ),
      ]);

      setPost(postRes.ok ? await postRes.json() : null);

      if (commentsRes.ok) {
        const data = await commentsRes.json();
        setSnapshotMaxId(commentsRes.headers.get('X-Max-Comment-Id'));
        setComments(data || []);
        setHasMore((data || []).length === PAGE_SIZE);
      } else {
        setComments([]);
        setHasMore(false);
      }
    } catch {
      setPost(null);
      setComments([]);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  };

  const loadMoreComments = async () => {
    if (!snapshotMaxId || loadingMore) return;

    setLoadingMore(true);
    try {
      const res = await fetch(
        `${API_BASE}/posts/${postId}/comment?before=${snapshotMaxId}&limit=${PAGE_SIZE}&offset=${comments.length}`,
        { credentials: 'include' },
      );

      if (res.ok) {
        const data = await res.json();
        setComments((prev) => [...prev, ...(data || [])]);
        setHasMore((data || []).length === PAGE_SIZE);
      } else {
        setHasMore(false);
      }
    } catch {
      setHasMore(false);
    } finally {
      setLoadingMore(false);
    }
  };

  const addComment = async (content, imageFile) => {
    const form = new FormData();
    form.append('content', content);
    if (imageFile) {
      form.append('image', imageFile);
    }

    const res = await fetch(`${API_BASE}/posts/${postId}/comment`, {
      method: 'POST',
      credentials: 'include',
      body: form,
    });

    if (res.ok) {
      const newComment = await res.json();
      setComments((prev) => [newComment, ...prev]);
    }

    return res;
  };

  // ADDED: Function to handle reactions on this single post
  const reactToPost = async (type) => {
    try {
      const res = await fetch(`${API_BASE}/posts/react`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ post_id: Number(postId), type }),
      });

      if (res.status === 401) {
        window.location.href = '/login';
        return;
      }

      if (!res.ok) return;

      const data = await res.json();
      setPost((prev) => ({
        ...prev,
        like_count: data.like_count || 0,
        dislike_count: data.dislike_count || 0,
        user_reaction: data.user_reaction || null,
      }));
    } catch (err) {
      console.error('Failed to react:', err);
    }
  };

  // ADDED: Export reactToPost
  return {
    post,
    comments,
    loading,
    loadingMore,
    hasMore,
    fetchPostAndComments,
    loadMoreComments,
    addComment,
    reactToPost,
  };
}
