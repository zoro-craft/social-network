import { useEffect } from 'react';
import { useRouter } from 'expo-router';
import { useAuth } from '@/context/AuthContext';

// Guards a screen: redirects to /login once loading finishes if not authenticated.
// Return null while loading / not authenticated.
export function useRequireAuth() {
  const { auth, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !auth.authenticated) {
      router.replace('/login');
    }
  }, [loading, auth.authenticated, router]);

  return { auth, loading };
}