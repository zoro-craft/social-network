"use client";

import { useEffect, useState } from "react";
import { useRequireAuth } from "@/hooks/useRequireAuth";
import { apiRequest } from "@/utils/api";
import { useToast } from "@/app/context/ToastContext";
import { useRealtime } from "@/components/RealtimeProvider";

import {
  handleAcceptInvitation,
  handleRejectInvitation,
  handleAcceptRequest,
  handleRejectRequest,
  handleAcceptFollow,
  handleDeclineFollow,
  handleMarkRead,
} from "./utils";

export default function NotificationsPage() {
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();
  const { subscribe } = useRealtime();

  const [notifications, setNotifications] = useState([]);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);

  const removeNotification = (id) => {
    setNotifications((prev) => {
      const next = prev.filter((n) => n.id !== id);
      // Re-anchor the cursor on the oldest row still visible (keyset paging).
      return next;
    });
  };

  // Fetch the first page of unread notifications, newest first, 10 rows + exact unread count.
  useEffect(() => {
    if (!auth.authenticated) return;
    let cancelled = false;
    apiRequest("/notifications")
      .then((data) => {
        if (cancelled) return;
        const items = data.items || [];
        setNotifications(items);
        setHasMore(data.has_more);
        // This fetch carried the exact unread total; reconcile the mounted bell.
        window.dispatchEvent(new Event('notifications:changed'));
      })
      .catch(() => {
        if (cancelled) return;
        setNotifications([]);
        setHasMore(false);
      });
    return () => {
      cancelled = true;
    };
  }, [auth.authenticated]);

  // Live arrivals prepend; dedupe by id so a row is never shown twice.
  useEffect(() => {
    if (!auth.authenticated) return;
    return subscribe("notification", (payload) => {
      setNotifications((prev) => {
        if (prev.some((n) => n.id === payload.id)) return prev;
        return [payload, ...prev];
      });
    });
  }, [auth.authenticated, subscribe]);

  // Removal signal: a notification was deleted server-side (e.g. a cancelled follow
  // request); drop the row so every open tab removes it without a refresh.
  useEffect(() => {
    if (!auth.authenticated) return;
    return subscribe("notification_removed", (payload) => {
      setNotifications((prev) => prev.filter((n) => n.id !== payload.id));
    });
  }, [auth.authenticated, subscribe]);

  // Load the next 10 rows older than the oldest row already shown.
  const loadMore = async () => {
    if (loadingMore) return;
    const oldest = notifications[notifications.length - 1];
    if (!oldest) return;
    setLoadingMore(true);
    try {
      const data = await apiRequest(`/notifications?before_id=${oldest.id}`);
      const items = data.items || [];
      setNotifications((prev) => {
        const seen = new Set(prev.map((n) => n.id));
        return [...prev, ...items.filter((n) => !seen.has(n.id))];
      });
      setHasMore(data.has_more);
    } catch {
      setHasMore(false);
    } finally {
      setLoadingMore(false);
    }
  };

  const respond = async (handler, notification) => {
    const result = await handler(notification);
    if (result.ok) {
      removeNotification(notification.id);
      window.dispatchEvent(new Event('notifications:changed'));
    } else if (result.status === 400 || result.status === 404) {
      showToast("This notification is no longer available");
    } else {
      showToast("Something went wrong. Please try again.");
    }
  };

  const markRead = async (notification, index) => {
    if (notification.is_read) return;
    removeNotification(notification.id);
    const result = await handleMarkRead(notification);
    if (result.ok) {
      window.dispatchEvent(new Event('notifications:changed'));
      return;
    }
    setNotifications((prev) => {
      if (prev.some((n) => n.id === notification.id)) return prev;
      const next = [...prev];
      next.splice(Math.min(index, next.length), 0, notification);
      return next;
    });
    if (result.status === 400 || result.status === 404) {
      showToast("This notification is no longer available");
    } else {
      showToast("Something went wrong. Please try again.");
    }
  };

  const renderActions = (notification) => {
    switch (notification.type) {
      case "group_invitation":
        return (
          <div className="glass-notification-actions">
            <button onClick={() => respond(handleAcceptInvitation, notification)}>
              Accept
            </button>
            <button onClick={() => respond(handleRejectInvitation, notification)}>
              Decline
            </button>
          </div>
        );
      case "join_request":
        return (
          <div className="glass-notification-actions">
            <button onClick={() => respond(handleAcceptRequest, notification)}>
              Accept
            </button>
            <button onClick={() => respond(handleRejectRequest, notification)}>
              Decline
            </button>
          </div>
        );
      case "follow_request":
        return (
          <div className="glass-notification-actions">
            <button onClick={() => respond(handleAcceptFollow, notification)}>
              Accept
            </button>
            <button onClick={() => respond(handleDeclineFollow, notification)}>
              Decline
            </button>
          </div>
        );
      case "event_created":
        return (
          <div className="glass-notification-actions">
            <button onClick={() => respond(handleMarkRead, notification)}>
              Mark as read
            </button>
          </div>
        );
      default:
        return null;
    }
  };

  if (loading) return <div>Loading...</div>;
  if (!auth.authenticated) return null;

  return (
    <div className="glass-notifications-page">
      <h1>Notifications</h1>

      {notifications.length === 0 ? (
        <p className="glass-notification-empty">No notifications</p>
      ) : (
        <>
          <div className="glass-notifications-list">
            {notifications.map((notification, index) => (
              <div
                key={notification.id}
                className={`glass-notification-card ${
                  notification.is_read ? "glass-notification-read" : "glass-notification-unread"
                }`}
                onClick={(e) => {
                  if (e.target.closest("button, a")) return;
                  markRead(notification, index);
                }}
              >
                <div className="glass-notification-content">
                  <p>{notification.message}</p>
                  <small>{new Date(notification.created_at).toLocaleString()}</small>
                </div>
                {renderActions(notification)}
              </div>
            ))}
          </div>
          {hasMore && (
            <button
              className="glass-notifications-load-more"
              onClick={loadMore}
              disabled={loadingMore}
            >
              {loadingMore ? "Loading..." : "Load more"}
            </button>
          )}
        </>
      )}
    </div>
  );
}