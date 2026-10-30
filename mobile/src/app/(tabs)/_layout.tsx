import { useCallback, useEffect, useState } from 'react';
import { Redirect } from 'expo-router';
import { StyleSheet, Text, View } from 'react-native';
import { Tabs } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useRealtime } from '@/context/RealtimeContext';
import { apiRequest } from '@/lib/api';
import { colors, fonts } from '@/lib/theme';

function NotificationTabBadge() {
  const { auth } = useAuth();
  const { subscribe } = useRealtime();
  const [unread, setUnread] = useState(0);

  const refresh = useCallback(async () => {
    if (!auth.authenticated) return;
    try {
      const data = await apiRequest<{ unread_count: number }>('/notifications');
      setUnread(data.unread_count || 0);
    } catch {
      // ignore
    }
  }, [auth.authenticated]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  useEffect(() => {
    if (!auth.authenticated) return;
    const unsub1 = subscribe('notification', () => setUnread((u) => u + 1));
    const unsub2 = subscribe('notification_removed', () => setUnread((u) => Math.max(0, u - 1)));
    return () => {
      unsub1();
      unsub2();
    };
  }, [auth.authenticated, subscribe]);

  if (unread <= 0) return null;
  return (
    <View style={styles.badge}>
      <Text style={styles.badgeText}>{unread > 99 ? '99+' : unread}</Text>
    </View>
  );
}

export default function TabLayout() {
  const { auth, loading } = useAuth();

  if (loading) return <Loading />;
  if (!auth.authenticated) return <Redirect href="/login" />;

  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarStyle: styles.tabBar,
        tabBarActiveTintColor: colors.textPrimary,
        tabBarInactiveTintColor: colors.textMuted,
        tabBarLabelStyle: styles.tabLabel,
      }}
    >
      <Tabs.Screen
        name="feed"
        options={{
          title: 'Feed',
          tabBarIcon: ({ color, size }) => <Ionicons name="home" size={size} color={color} />,
        }}
      />
      <Tabs.Screen
        name="reels"
        options={{
          title: 'Reels',
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="play-circle" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="groups"
        options={{
          title: 'Groups',
          tabBarIcon: ({ color, size }) => <Ionicons name="people" size={size} color={color} />,
        }}
      />
      <Tabs.Screen
        name="chat"
        options={{
          title: 'Chat',
          tabBarIcon: ({ color, size }) => <Ionicons name="chatbubble" size={size} color={color} />,
        }}
      />
      <Tabs.Screen
        name="notifications"
        options={{
          title: 'Notifications',
          tabBarIcon: ({ color, size }) => (
            <View>
              <Ionicons name="notifications" size={size} color={color} />
              <NotificationTabBadge />
            </View>
          ),
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Profile',
          tabBarIcon: ({ color, size }) => <Ionicons name="person-circle" size={size} color={color} />,
        }}
      />
    </Tabs>
  );
}

const styles = StyleSheet.create({
  tabBar: {
    backgroundColor: colors.bg,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  tabLabel: {
    fontFamily: fonts.sansMedium,
    fontSize: 11,
  },
  badge: {
    position: 'absolute',
    top: -4,
    right: -10,
    minWidth: 18,
    height: 18,
    borderRadius: 9,
    backgroundColor: colors.accent,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 4,
  },
  badgeText: {
    color: colors.onAccent,
    fontFamily: fonts.sansBold,
    fontSize: 10,
  },
});