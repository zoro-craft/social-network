'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from './context/AuthContext';

export default function Home() {
  const router = useRouter();
  const { auth, loading } = useAuth();

  useEffect(() => {
    if (!loading && !auth.authenticated) {
      router.push('/login');
    }
  }, [loading, auth.authenticated, router]);

  if (loading || !auth.authenticated) {
    return null; // or a loading spinner
  }

  return (
    <div
      style={{
        maxWidth: 600,
        color: 'white',
        margin: '60px auto',
        textAlign: 'center',
      }}
    >
      <h1>Welcome to 01Forum</h1>
      <p>Connect, share, and chat with people around the world.</p>
    </div>
  );
}
