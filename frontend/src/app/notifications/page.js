'use client';

import { useEffect, useState } from 'react';

import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';

import {
  handleAcceptInvitation,
  handleRejectInvitation,
  handleAcceptRequest,
  handleRejectRequest,
} from './utils';

import styles from '../styles/notifications.module.css';

export default function NotificationsPage() {
  const { auth, loading } = useRequireAuth();

  const [notifications, setNotifications] = useState([]);
  const [followRequests, setFollowRequests] = useState([]);

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

    const fetchFollowRequests = async () => {
      try {
        const res = await apiRequest('/follow/requests');
        setFollowRequests(res || []);
      } catch (error) {
        console.error('Failed to fetch follow requests:', error);
        setFollowRequests([]);
      }
    };

    fetchNotifications();
    fetchFollowRequests();
  }, [auth.authenticated]);

  // Remove notification from UI
  const removeNotification = (notificationId) => {
    setNotifications((prev) =>
      prev.filter((notification) => notification.id !== notificationId)
    );
  };

  const respondToFollowRequest = async (followerId, action) => {
    try {
      await apiRequest(`/follow/${followerId}/${action}`, 'POST');
      setFollowRequests((prev) =>
        prev.filter((req) => req.follower_id !== followerId)
      );
    } catch (error) {
      console.error('Failed to respond to follow request:', error);
    }
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

      {followRequests.length > 0 && (
        <>
          <h1 style={{ marginTop: 35 }}>Follow Requests</h1>

          <div className={styles.notificationList}>
            {followRequests.map((req) => (
              <div key={req.follower_id} className={styles.notification}>
                <div className={styles.content}>
                  <p>
                    {req.first_name}
                    {req.last_name ? ` ${req.last_name}` : ''} wants to
                    follow you
                  </p>
                </div>

                <div className={styles.actions}>
                  <button
                    onClick={() =>
                      respondToFollowRequest(req.follower_id, 'accept')
                    }
                  >
                    Accept
                  </button>

                  <button
                    onClick={() =>
                      respondToFollowRequest(req.follower_id, 'decline')
                    }
                  >
                    Decline
                  </button>
                </div>
              </div>
            ))}
          </div>
        </>
      )}

      {notifications.length === 0 ? (
        <p className={styles.empty}>No notifications</p>
      ) : (
        <div className={styles.notificationList}>
          {notifications.map((notification) => (
            <div
              key={notification.id}
              className={`${styles.notification} ${
                notification.is_read ? styles.read : styles.unread
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

              {/* Group Invitation */}
              {notification.type === 'group_invtation' && (
                <div className={styles.actions}>
                  <button
                    onClick={async () => {
                      const success =
                        await handleAcceptInvitation(notification);

                      if (success) {
                        removeNotification(notification.id);
                      }
                    }}
                  >
                    Accept
                  </button>

                  <button
                    onClick={async () => {
                      const success =
                        await handleRejectInvitation(notification);

                      if (success) {
                        removeNotification(notification.id);
                      }
                    }}
                  >
                    Reject
                  </button>
                </div>
              )}

              {/* Join Request */}
              {notification.type === 'join_request' && (
                <div className={styles.actions}>
                  <button
                    onClick={async () => {
                      const success =
                        await handleAcceptRequest(notification);

                      if (success) {
                        removeNotification(notification.id);
                      }
                    }}
                  >
                    Accept
                  </button>

                  <button
                    onClick={async () => {
                      const success =
                        await handleRejectRequest(notification);

                      if (success) {
                        removeNotification(notification.id);
                      }
                    }}
                  >
                    Reject
                  </button>
                </div>
              )}

              {notification.type === 'message' && (
                <span className={styles.messageType}>Message</span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}