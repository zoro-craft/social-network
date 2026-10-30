import { useCallback, useEffect, useState } from 'react';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { useFocusEffect, useRouter } from 'expo-router';
import Avatar from '@/components/Avatar';
import SegmentedControl from '@/components/SegmentedControl';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useRealtime } from '@/context/RealtimeContext';
import { apiRequest } from '@/lib/api';
import { ChatMessage, ConversationsResponse } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

export default function ChatScreen() {
  const router = useRouter();
  const { auth } = useAuth();
  const { subscribe } = useRealtime();

  const [tab, setTab] = useState<'users' | 'groups'>('users');
  const [conversations, setConversations] = useState<ConversationsResponse>({ users: [], groups: [] });
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<ConversationsResponse>('/conversations');
      setConversations({ users: data.users || [], groups: data.groups || [] });
    } catch {
      setConversations({ users: [], groups: [] });
    } finally {
      setLoading(false);
    }
  }, []);

  // Refetch on focus so a thread started elsewhere shows up in the sidebar.
  useFocusEffect(
    useCallback(() => {
      if (auth.authenticated) load();
    }, [auth.authenticated, load]),
  );

  // Live messages refresh the sidebar row preview + order.
  useEffect(() => {
    if (!auth.authenticated) return;
    const onMessage = (payload: ChatMessage) => {
      if (payload.type === 'private') {
        const otherID = payload.sender_id === auth.user?.id ? payload.ref_id : payload.sender_id;
        setConversations((prev) => {
          const entry = prev.users.find((u) => u.user_id === otherID);
          if (!entry) return prev;
          const updated = { ...entry, last_message: payload.content, last_message_at: payload.created_at };
          return { ...prev, users: [updated, ...prev.users.filter((u) => u.user_id !== otherID)] };
        });
      } else {
        setConversations((prev) => {
          const entry = prev.groups.find((g) => g.group_id === payload.ref_id);
          if (!entry) return prev;
          const updated = { ...entry, last_message: payload.content, last_message_at: payload.created_at };
          return { ...prev, groups: [updated, ...prev.groups.filter((g) => g.group_id !== payload.ref_id)] };
        });
      }
    };
    const unsub1 = subscribe('chat_new_message', onMessage);
    const unsub2 = subscribe('chat_new_group_message', onMessage);
    return () => {
      unsub1();
      unsub2();
    };
  }, [auth.authenticated, auth.user?.id, subscribe]);

  // Presence toggles the online dots live.
  useEffect(() => {
    if (!auth.authenticated) return;
    return subscribe('presence', (payload: { user_id: number; online: boolean }) => {
      setConversations((prev) => ({
        ...prev,
        users: prev.users.map((u) => (u.user_id === payload.user_id ? { ...u, online: payload.online } : u)),
      }));
    });
  }, [auth.authenticated, subscribe]);

  const openThread = (type: 'private' | 'group', id: number, name: string) => {
    router.push(`/chat/thread?type=${type}&id=${id}&name=${encodeURIComponent(name)}`);
  };

  if (loading) {
    return <Loading />;
  }

  const isUsers = tab === 'users';
  const rows: {
    key: string;
    id: number;
    name: string;
    avatar?: string;
    online?: boolean;
    last?: string | null;
    time?: number | null;
  }[] = isUsers
    ? conversations.users.map((u) => ({
        key: `private:${u.user_id}`,
        id: u.user_id,
        name: u.name,
        avatar: u.avatar,
        online: u.online,
        last: u.last_message,
        time: u.last_message_at,
      }))
    : conversations.groups.map((g) => ({
        key: `group:${g.group_id}`,
        id: g.group_id,
        name: g.title,
        last: g.last_message,
        time: g.last_message_at,
      }));

  return (
    <View style={styles.screen}>
      <View style={styles.header}>
        <Text style={styles.title}>Chat</Text>
      </View>
      <View style={styles.tabWrap}>
        <SegmentedControl
          options={[
            { label: 'Users', value: 'users' },
            { label: 'Groups', value: 'groups' },
          ]}
          value={tab}
          onChange={(v) => setTab(v as 'users' | 'groups')}
        />
      </View>
      <FlatList
        data={rows}
        keyExtractor={(item) => item.key}
        contentContainerStyle={styles.list}
        ListEmptyComponent={<EmptyState message={isUsers ? 'No conversations yet.' : 'No group chats yet.'} />}
        renderItem={({ item }) => (
          <Pressable
            style={styles.row}
            onPress={() => openThread((isUsers ? 'private' : 'group') as any, item.id, item.name)}
          >
            {isUsers ? (
              <View>
                <Avatar src={item.avatar} name={item.name} size={40} />
                {item.online ? <View style={styles.onlineDot} /> : null}
              </View>
            ) : (
              <View style={styles.groupIcon}>
                <Text style={styles.groupIconText}>{(item.name || '?').charAt(0).toUpperCase()}</Text>
              </View>
            )}
            <View style={styles.rowBody}>
              <View style={styles.rowTop}>
                <Text style={styles.rowName} numberOfLines={1}>
                  {item.name}
                </Text>
                {item.time ? <Text style={styles.rowTime}>{formatChatTime(item.time)}</Text> : null}
              </View>
              <Text style={styles.rowLast} numberOfLines={1}>
                {item.last || 'No messages yet'}
              </Text>
            </View>
          </Pressable>
        )}
      />
    </View>
  );
}

function formatChatTime(ts: number) {
  const d = new Date(ts);
  const now = new Date();
  const sameDay = d.toDateString() === now.toDateString();
  if (sameDay) {
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
  return d.toLocaleDateString();
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
  tabWrap: {
    padding: spacing.lg,
  },
  list: {
    paddingHorizontal: spacing.lg,
    gap: spacing.xs,
    paddingBottom: spacing.xl,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  onlineDot: {
    position: 'absolute',
    right: 0,
    bottom: 0,
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: colors.success,
    borderWidth: 2,
    borderColor: colors.bg,
  },
  groupIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: colors.bgHover,
    borderWidth: 1,
    borderColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  groupIconText: {
    color: colors.textSecondary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 16,
  },
  rowBody: {
    flex: 1,
  },
  rowTop: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  rowName: {
    flex: 1,
    color: colors.textPrimary,
    fontFamily: fonts.sansMedium,
    fontSize: 15,
  },
  rowTime: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 12,
    marginLeft: spacing.sm,
  },
  rowLast: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 13,
    marginTop: 2,
  },
});