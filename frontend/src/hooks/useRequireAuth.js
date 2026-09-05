'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/context/AuthContext';

export function useRequireAuth() {
  const router = useRouter();
  const { auth, loading } = useAuth();

  useEffect(() => {
    if (!loading && !auth.authenticated) {
      router.push('/login');
    }
  }, [loading, auth.authenticated, router]);

  return { auth, loading };
}

