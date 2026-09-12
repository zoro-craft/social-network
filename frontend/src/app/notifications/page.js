'use client';

import { useEffect, useState } from 'react';

import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/utils/api';


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
    <div className="glass-notifications-page">
      <h1>Notifications</h1>

      {followRequests.length > 0 && (
        <>
          <h1 style={{ marginTop: 35 }}>Follow Requests</h1>

          <div className="glass-notifications-list">
            {followRequests.map((req) => (
              <div key={req.follower_id} className="glass-notification-card">
                <div className="glass-notification-content">
                  <p>
                    {req.first_name}
                    {req.last_name ? ` ${req.last_name}` : ''} wants to
                    follow you
                  </p>
                </div>

                <div className="glass-notification-actions">
                  <button
                    onClick={() =>
                      respondToFollowRequest(
                        req.follower_id,
                        'accept'
                      )
                    }
                  >
                    Accept
                  </button>

                  <button
                    onClick={() =>
                      respondToFollowRequest(
                        req.follower_id,
                        'decline'
                      )
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

      {notifications.length === 0 && followRequests.length === 0 ? (
        <p className="glass-notification-empty">
          No notifications
        </p>
      ) : (
        <div className="glass-notifications-list">
          {notifications.map((notification) => (
            <div
              key={notification.id}
              className={`glass-notification-card ${
                notification.is_read
                  ? 'glass-notification-read'
                  : 'glass-notification-unread'
              }`}
            >
              <div className="glass-notification-content">
                <p>{notification.message}</p>

                <small>
                  {new Date(
                    notification.created_at
                  ).toLocaleString()}
                </small>
              </div>

              {notification.type === 'group_invtation' && (
                <div className="glass-notification-actions">
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
                <div className="glass-notification-actions">
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
                <span className="glass-message-type">
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