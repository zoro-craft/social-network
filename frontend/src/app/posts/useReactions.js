'use client';

import { useRouter } from 'next/navigation';
import { API_BASE } from '@/utils/api';

export function useReactions(setPosts) {
  const router = useRouter();

  const react = async (postId, type) => {
    try {
      const res = await fetch(`${API_BASE}/posts/react`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ post_id: postId, type }),
      });

      if (res.status === 401) {
        router.push('/login');
        return;
      }

      if (!res.ok) {
        console.error('Failed to react to post');
        return;
      }

      const data = await res.json();

      setPosts((prev) =>
        prev.map((p) =>
          p.id === postId
            ? {
                ...p,
                like_count: data.like_count || 0,
                dislike_count: data.dislike_count || 0,
                user_reaction: data.user_reaction || null,
              }
            : p,
        ),
      );
    } catch (err) {
      console.error('Reaction error:', err);
    }
  };

  return { react };
}
