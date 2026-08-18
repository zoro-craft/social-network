'use client';

import { useRouter } from 'next/navigation';
import { API_BASE } from '../lib/api';

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

      const data = await res.json();
      if (!res.ok) {
        router.push('/login')
        return;
      }

      setPosts((prev) =>
        prev.map((p) =>
          p.id === postId
            ? {
                ...p,
                like_count: data.like_count,
                dislike_count: data.dislike_count,
                user_reaction: data.user_reaction,
              }
            : p,
        ),
      );
    } catch {
      // ignore for now
    }
  };

  return { react };
}
