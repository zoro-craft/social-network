'use client';

import { useEffect, useState } from 'react';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api'; 
import { handleAcceptInvitation , handleRejectInvitation , handleAcceptRequest , handleRejectRequest } from './utils';
import styles from '../styles/notifications.module.css';


export default function NotificationsPage() {
  const { auth, loading } = useRequireAuth();

  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    if (!auth.authenticated) return;

    const fetchNotifications = async () => {
      try {
        const res = await apiRequest('/notifications');
        setNotifications(res || []);
      } catch (error) {
        console.error('Failed to fetch notifications:', error);
        setNotifications([]);
      }
    };

    fetchNotifications();
  }, [auth.authenticated]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!auth.authenticated) {
    return null;
  }
  return (
    <div className={styles.container}>
      <h1>Notifications</h1>

      {notifications.length === 0 ? (
        <p className={styles.empty}>
          No notifications
        </p>
      ) : (
        <div className={styles.notificationList}>
          {notifications.map((notification) => (
            <div key={notification.id} className={styles.notification} >
              <div className={styles.content}>
                <p>{notification.message}</p>
              </div>
              {notification.type === 'group_invtation' && (
                <div className={styles.actions}>
                  <button
                    onClick={() =>
                      handleAcceptInvitation(notification)
                    }
                  >
                    Accept
                  </button>

                  <button
                    onClick={() =>
                      handleRejectInvitation(notification)
                    }
                  >
                    Reject
                  </button>
                </div>
              )}

              {notification.type === 'join_request' && (
                <div className={styles.actions}>
                  <button
                    onClick={() =>
                      handleAcceptRequest(notification)
                    }
                  >
                    Accept
                  </button>

                  <button
                    onClick={() =>
                      handleRejectRequest(notification)
                    }
                  >
                    Reject
                  </button>
                </div>
              )}

            </div>
          ))}
        </div>
      )}
    </div>
  );
}