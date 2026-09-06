'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/context/ToastContext';
import { useRouter } from 'next/navigation';
import { API_BASE } from '@/utils/api';
import Avatar from '@/components/Avatar';
import styles from '../styles/profile.module.css';

export default function NotificationsPage() {
  const { auth, loading: authLoading } = useAuth();
  const router = useRouter();
  const showToast = useToast();
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!authLoading && !auth.authenticated) {
      router.push('/login');
    }
  }, [authLoading, auth.authenticated, router]);

  useEffect(() => {
    if (!auth.authenticated) return;

    async function fetchRequests() {
      try {
        const res = await fetch(`${API_BASE}/follow/requests`, {
          credentials: 'include',
        });
        const data = await res.json();
        setRequests(data || []);
      } catch {
        setRequests([]);
      } finally {
        setLoading(false);
      }
    }

    fetchRequests();
  }, [auth.authenticated]);

  const respond = async (followerId, action) => {
    try {
      const res = await fetch(`${API_BASE}/follow/${followerId}/${action}`, {
        method: 'POST',
        credentials: 'include',
      });
      if (res.ok) {
        setRequests((prev) => prev.filter((r) => r.follower_id !== followerId));
      } else {
        showToast('Failed to respond');
      }
    } catch {
      showToast('Could not reach the server');
    }
  };

  if (authLoading || !auth.authenticated) return null;

  return (
    <div className={styles.container}>
      <h1 className={styles.name}>Follow Requests</h1>

      {loading && <p className={styles.status}>Loading...</p>}
      {!loading && requests.length === 0 && (
        <p className={styles.status}>No pending requests</p>
      )}

      {requests.map((req) => (
        <div key={req.follower_id} className={styles.requestCard}>
          <Avatar src={req.avatar} name={req.first_name} size={40} />
          <div className={styles.requestInfo}>
            <p className={styles.requestName}>{req.first_name}</p>
            {req.email && <p className={styles.requestEmail}>{req.email}</p>}
          </div>
          <div className={styles.requestActions}>
            <button
              className={`${styles.requestButton} ${styles.acceptButton}`}
              onClick={() => respond(req.follower_id, 'accept')}
            >
              Accept
            </button>
            <button
              className={`${styles.requestButton} ${styles.declineButton}`}
              onClick={() => respond(req.follower_id, 'decline')}
            >
              Decline
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
