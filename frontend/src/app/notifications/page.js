'use client';

import { useEffect, useState } from 'react';

import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';

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


  const handleAcceptInvitation = async (notification) => {
    try {
      await apiRequest(
        `/invitations/${notification.invitation_id}/accept`,
        {
          method: 'POST',
        }
      );

      removeNotification(notification.id);
    } catch (error) {
      console.error('Failed to accept invitation:', error);
    }
  };

  const handleRejectInvitation = async (notification) => {
    try {
      await apiRequest(
        `/invitations/${notification.invitation_id}/reject`,
        {
          method: 'POST',
        }
      );

      removeNotification(notification.id);
    } catch (error) {
      console.error('Failed to reject invitation:', error);
    }
  };

  const handleAcceptRequest = async (notification) => {
    try {
      await apiRequest(
        `/join_requests/${notification.request_id}/accept`,
        {
          method: 'POST',
        }
      );

      removeNotification(notification.id);
    } catch (error) {
      console.error('Failed to accept join request:', error);
    }
  };

  const handleRejectRequest = async (notification) => {
    try {
      await apiRequest(
        `/join_requests/${notification.request_id}/reject`,
        {
          method: 'POST',
        }
      );

      removeNotification(notification.id);
    } catch (error) {
      console.error('Failed to reject join request:', error);
    }
  };

  // Remove notification from UI
  const removeNotification = (notificationId) => {
    setNotifications((prev) =>
      prev.filter((notification) => notification.id !== notificationId)
    );
  };

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
            <div
              key={notification.id}
              className={`${styles.notification} ${
                notification.is_read
                  ? styles.read
                  : styles.unread
              }`}
            >
              <div className={styles.content}>
                <p>{notification.message}</p>

                <small>
                  {new Date(
                    notification.created_at
                  ).toLocaleString()}
                </small>
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

              {notification.type === 'message' && (
                <span className={styles.messageType}>
                  Message
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}