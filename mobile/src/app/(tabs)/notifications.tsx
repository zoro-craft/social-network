import { useCallback, useEffect, useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import AppButton from '@/components/AppButton';
import Card from '@/components/Card';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useRealtime } from '@/context/RealtimeContext';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import formatDate from '@/lib/format';
import {
  handleAcceptFollow,
  handleAcceptInvitation,
  handleAcceptRequest,
  handleDeclineFollow,
  handleMarkRead,
  handleRejectInvitation,
  handleRejectRequest,
} from '@/lib/notifications';
import { Notification } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

export default function NotificationsScreen() {
  const { auth } = useAuth();
  const { subscribe } = useRealtime();
  const showToast = useToast();

  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [busyId, setBusyId] = useState<number | null>(null);

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<{ items: Notification[]; has_more: boolean }>('/notifications');
      setNotifications(data.items || []);
      setHasMore(data.has_more);
    } catch {
      setNotifications([]);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  useEffect(() => {
    if (!auth.authenticated) return;
    const unsub1 = subscribe('notification', (payload: Notification) => {
      setNotifications((prev) => (prev.some((n) => n.id === payload.id) ? prev : [payload, ...prev]));
    });
    const unsub2 = subscribe('notification_removed', (payload: Notification) => {
      setNotifications((prev) => prev.filter((n) => n.id !== payload.id));
    });
    return () => {
      unsub1();
      unsub2();
    };
  }, [auth.authenticated, subscribe]);

  const removeNotification = (id: number) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
  };

  const loadMore = async () => {
    if (loadingMore) return;
    const oldest = notifications[notifications.length - 1];
    if (!oldest) return;
    setLoadingMore(true);
    try {
      const data = await apiRequest<{ items: Notification[]; has_more: boolean }>(
        `/notifications?before_id=${oldest.id}`,
      );
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

  const respond = async (handler: (n: Notification) => Promise<{ ok: boolean; status?: number }>, notification: Notification) => {
    setBusyId(notification.id);
    const result = await handler(notification);
    setBusyId(null);
    if (result.ok) {
      removeNotification(notification.id);
    } else if (result.status === 400 || result.status === 404) {
      showToast('This notification is no longer available');
    } else {
      showToast('Something went wrong. Please try again.');
    }
  };

  const renderActions = (notification: Notification) => {
    // Both buttons share the row's busy flag so a request can't be double-fired.
    const busy = busyId === notification.id;

    const pair = (
      accept: (n: Notification) => Promise<{ ok: boolean; status?: number }>,
      decline: (n: Notification) => Promise<{ ok: boolean; status?: number }>,
      labels: [string, string] = ['Accept', 'Decline'],
    ) => (
      <View style={styles.actions}>
        <AppButton
          title={labels[0]}
          variant="success"
          onPress={() => respond(accept, notification)}
          loading={busy}
          disabled={busy}
          style={styles.smallBtn}
        />
        <AppButton
          title={labels[1]}
          variant="danger"
          onPress={() => respond(decline, notification)}
          disabled={busy}
          style={styles.smallBtn}
        />
      </View>
    );

    switch (notification.type) {
      case 'group_invitation':
        return pair(handleAcceptInvitation, handleRejectInvitation);
      case 'join_request':
        return pair(handleAcceptRequest, handleRejectRequest);
      case 'follow_request':
        return pair(handleAcceptFollow, handleDeclineFollow);
      case 'event_created':
        return (
          <View style={styles.actions}>
            <AppButton
              title="Mark as read"
              variant="ghost"
              onPress={() => respond(handleMarkRead, notification)}
              loading={busy}
              disabled={busy}
              style={styles.smallBtn}
            />
          </View>
        );
      default:
        return null;
    }
  };

  if (loading) {
    return (
      <View style={styles.screen}>
        <View style={styles.header}>
          <Text style={styles.title}>Notifications</Text>
        </View>
        <Loading />
      </View>
    );
  }

  return (
    <View style={styles.screen}>
      <View style={styles.header}>
        <Text style={styles.title}>Notifications</Text>
      </View>
      {notifications.length === 0 ? (
        <EmptyState message="No notifications" />
      ) : (
        <ScrollView contentContainerStyle={styles.list}>
          {notifications.map((notification) => (
            <Card
              key={notification.id}
              style={[
                styles.card,
                !notification.is_read && styles.cardUnread,
              ]}
            >
              <View style={styles.content}>
                <Text style={styles.message}>{notification.message}</Text>
                <Text style={styles.time}>{formatDate(notification.created_at)}</Text>
              </View>
              {renderActions(notification)}
            </Card>
          ))}
          {hasMore && (
            <AppButton
              title={loadingMore ? 'Loading...' : 'Load more'}
              variant="ghost"
              loading={loadingMore}
              onPress={loadMore}
            />
          )}
        </ScrollView>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  header: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  title: {
    color: colors.textPrimary,
    fontFamily: fonts.sansBold,
    fontSize: 22,
    letterSpacing: -0.02,
  },
  list: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  card: {
    padding: spacing.lg,
    gap: spacing.md,
    borderLeftWidth: 3,
    borderLeftColor: colors.border,
  },
  cardUnread: {
    borderLeftColor: colors.accent,
  },
  content: {
    gap: spacing.xs,
  },
  message: {
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
  time: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  actions: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  smallBtn: {
    flex: 1,
    paddingVertical: spacing.sm,
  },
});