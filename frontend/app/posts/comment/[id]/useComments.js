'use client';

import { useState } from 'react';
import { API_BASE } from '../../../lib/api';

export function useComments(postId) {
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchPostAndComments = async () => {
    setLoading(true);
    try {
      const [postRes, commentsRes] = await Promise.all([
        fetch(`${API_BASE}/posts/${postId}`, { credentials: 'include' }),
        fetch(`${API_BASE}/posts/${postId}/comment`, {
          credentials: 'include',
        }),
      ]);

      setPost(postRes.ok ? await postRes.json() : null);
      setComments(commentsRes.ok ? (await commentsRes.json()) || [] : []);
    } catch {
      setPost(null);
      setComments([]);
    } finally {
      setLoading(false);
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

  return { post, comments, loading, fetchPostAndComments, addComment };
}
