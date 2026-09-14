"use client";

import { useEffect, useState } from "react";
import { useRequireAuth } from "@/hooks/useRequireAuth";
import { apiRequest } from "@/utils/api";
import { useToast } from "@/app/context/ToastContext";

import {
  handleAcceptInvitation,
  handleRejectInvitation,
  handleAcceptRequest,
  handleRejectRequest,
} from "./utils";

import styles from "../styles/notifications.module.css";

export default function NotificationsPage() {
  const { auth, loading } = useRequireAuth();
  const showToast = useToast();

  const [notifications, setNotifications] = useState([]);
  const [followRequests, setFollowRequests] = useState([]);

  useEffect(() => {
    if (!auth.authenticated) return;

    const loadList = async (endpoint, setter) => {
      try {
        setter((await apiRequest(endpoint)) || []);
      } catch {
        setter([]);
      }
    };

    loadList("/notifications", setNotifications);
    loadList("/follow/requests", setFollowRequests);
  }, [auth.authenticated]);

  const removeNotification = (notificationId) =>
    setNotifications((prev) => prev.filter((n) => n.id !== notificationId));

  const respondToFollowRequest = async (followerId, action) => {
    try {
      await apiRequest(`/follow/${followerId}/${action}`, "POST");
      setFollowRequests((prev) => prev.filter((req) => req.follower_id !== followerId));
    } catch (err) {
      showToast(err.message);
    }
  };

  const respondToNotification = async (handler, notification) => {
    if (await handler(notification)) removeNotification(notification.id);
  };

  if (loading) return <div>Loading...</div>;
  if (!auth.authenticated) return null;

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
                    {req.last_name ? ` ${req.last_name}` : ""} wants to follow you
                  </p>
                </div>
                <div className={styles.actions}>
                  <button onClick={() => respondToFollowRequest(req.follower_id, "accept")}>
                    Accept
                  </button>
                  <button onClick={() => respondToFollowRequest(req.follower_id, "decline")}>
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
                <small>{new Date(notification.created_at).toLocaleString()}</small>
              </div>

              {notification.type === "group_invtation" && (
                <div className={styles.actions}>
                  <button
                    onClick={() => respondToNotification(handleAcceptInvitation, notification)}
                  >
                    Accept
                  </button>
                  <button
                    onClick={() => respondToNotification(handleRejectInvitation, notification)}
                  >
                    Reject
                  </button>
                </div>
              )}

              {notification.type === "join_request" && (
                <div className={styles.actions}>
                  <button
                    onClick={() => respondToNotification(handleAcceptRequest, notification)}
                  >
                    Accept
                  </button>
                  <button
                    onClick={() => respondToNotification(handleRejectRequest, notification)}
                  >
                    Reject
                  </button>
                </div>
              )}

              {notification.type === "message" && (
                <span className={styles.messageType}>Message</span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
