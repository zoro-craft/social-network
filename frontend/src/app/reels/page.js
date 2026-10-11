'use client';

import { useRequireAuth } from '@/hooks/useRequireAuth';
import ReelFeed from '@/app/components/ReelFeed';

export default function ReelsPage() {
  const { auth, loading } = useRequireAuth();

  if (loading || !auth.authenticated) {
    return null;
  }

  return (
    <div className="glass-reels-page">
      <ReelFeed />
    </div>
  );
}
