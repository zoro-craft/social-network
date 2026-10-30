import { useCallback, useEffect, useState } from 'react';
import { Image } from 'expo-image';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import AppButton from '@/components/AppButton';
import Card from '@/components/Card';
import HeaderBar from '@/components/HeaderBar';
import PostCard from '@/components/PostCard';
import SegmentedControl from '@/components/SegmentedControl';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { apiRequest, mediaSource } from '@/lib/api';
import formatDate, { formatEventDateTime } from '@/lib/format';
import { Event, Group, Post } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

export default function GroupDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth } = useAuth();

  const groupId = Number(id);
  const [group, setGroup] = useState<Group | null>(null);
  const [tab, setTab] = useState<'posts' | 'events'>('posts');
  const [loading, setLoading] = useState(true);

  const loadInfo = useCallback(async () => {
    setLoading(true);
    try {
      const data = await apiRequest<Group>(`/groups/${groupId}`, {
        method: 'POST',
        data: { group_id: groupId },
      });
      setGroup(data);
    } catch {
      setGroup(null);
    } finally {
      setLoading(false);
    }
  }, [groupId]);

  useEffect(() => {
    if (auth.authenticated) loadInfo();
  }, [auth.authenticated, loadInfo]);

  if (loading) {
    return (
      <View style={styles.screen}>
        <HeaderBar title="Group" onBack={() => router.back()} />
        <Loading />
      </View>
    );
  }

  if (!group) {
    return (
      <View style={styles.screen}>
        <HeaderBar title="Group" onBack={() => router.back()} />
        <EmptyState message="Group unavailable" />
      </View>
    );
  }

  return (
    <View style={styles.screen}>
      <HeaderBar title={group.titel} onBack={() => router.back()} />

      {/* Posts and events are plain views, so the whole page scrolls as one column. */}
      <ScrollView contentContainerStyle={styles.content}>
        <Card style={styles.headerCard}>
          <View style={styles.headerRow}>
            <GroupHeaderImage group={group} />
            <View style={styles.headerText}>
              <Text style={styles.groupTitle}>{group.titel}</Text>
              <Text style={styles.groupDesc}>{group.description}</Text>
              {typeof group.member_count === 'number' && (
                <Text style={styles.memberCount}>
                  {group.member_count} member{group.member_count === 1 ? '' : 's'}
                </Text>
              )}
            </View>
          </View>

          <View style={styles.actionGrid}>
            <AppButton title="Invite Members" variant="ghost" onPress={() => router.push(`/groups/${groupId}/invite`)} />
            <AppButton title="Create Post" variant="ghost" onPress={() => router.push(`/posts/create?group_id=${groupId}`)} />
            <AppButton title="Create Event" variant="ghost" onPress={() => router.push(`/groups/${groupId}/create-event`)} />
            <AppButton title="Group Members" variant="ghost" onPress={() => router.push(`/groups/${groupId}/members`)} />
          </View>
        </Card>

        <SegmentedControl
          options={[
            { label: 'Posts', value: 'posts' },
            { label: 'Events', value: 'events' },
          ]}
          value={tab}
          onChange={(v) => setTab(v as 'posts' | 'events')}
        />

        {tab === 'posts' ? (
          <GroupPosts groupId={groupId} />
        ) : (
          <GroupEvents groupId={groupId} />
        )}
      </ScrollView>
    </View>
  );
}

function GroupHeaderImage({ group }: { group: Group }) {
  const source = mediaSource(group.image);
  const [failed, setFailed] = useState(false);

  if (!source || failed) {
    return (
      <View style={styles.icon}>
        <Text style={styles.iconText}>
          {(group.titel || '?').charAt(0).toUpperCase()}
        </Text>
      </View>
    );
  }

  return (
    <Image
      source={source}
      style={styles.groupImage}
      contentFit="cover"
      onError={() => setFailed(true)}
    />
  );
}

function GroupPosts({ groupId }: { groupId: number }) {
  const router = useRouter();
  const { auth } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!auth.authenticated) return;
    apiRequest<Post[]>(`/groups/${groupId}/posts`)
      .then((data) => setPosts(data || []))
      .catch(() => setPosts([]))
      .finally(() => setLoading(false));
  }, [auth.authenticated, groupId]);

  if (loading) return <Loading label="Loading posts..." />;
  if (posts.length === 0) return <EmptyState message="No posts in this group yet." />;

  return (
    <View style={styles.list}>
      {posts.map((post) => (
        <PostCard
          key={post.id}
          post={post}
          onComment={(p) => router.push(`/posts/${p.id}`)}
        />
      ))}
    </View>
  );
}

function GroupEvents({ groupId }: { groupId: number }) {
  const { auth } = useAuth();
  const showToast = useToast();
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<Event[]>(`/events?group_id=${groupId}`);
      setEvents(data || []);
    } catch {
      showToast('Failed to fetch events');
    } finally {
      setLoading(false);
    }
  }, [groupId, showToast]);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  const vote = async (event: Event, option: string) => {
    try {
      const res = await apiRequest<{ error?: string }>('/events/vote', {
        method: 'POST',
        data: { event_id: event.id, vote: option },
      });
      if (res && (res as any).error) {
        showToast((res as any).error);
      } else {
        await load();
      }
    } catch (err: any) {
      showToast(err?.message || 'Failed to save vote');
    }
  };

  if (loading) return <Loading label="Loading events..." />;
  if (events.length === 0) return <EmptyState message="No events yet. Create the first one!" />;

  return (
    <View style={styles.list}>
      {events.map((event) => {
        const when = formatEventDateTime(event.event_date, event.event_time);
        return (
          <Card key={event.id}>
            <View style={styles.eventHeader}>
              <Text style={styles.eventTitle}>{event.title}</Text>
              <Text style={styles.eventTime}>{formatDate(event.created_at)}</Text>
            </View>
            {event.description ? <Text style={styles.eventDesc}>{event.description}</Text> : null}
            {when ? (
              <View style={styles.whenPill}>
                <Ionicons name="calendar-outline" size={14} color={colors.accentHover} />
                <Text style={styles.whenText}>{when}</Text>
              </View>
            ) : null}
            <View style={styles.voteRow}>
              {(event.options || []).map((option) => {
                const count = (event.vote_counts || []).find((vc) => vc.option === option)?.count || 0;
                const selected = event.user_response === option;
                return (
                  <Pressable
                    key={option}
                    onPress={() => vote(event, option)}
                    style={[styles.voteOption, selected && styles.voteSelected]}
                  >
                    <Text style={[styles.voteLabel, selected && styles.voteLabelSelected]}>
                      {option}
                    </Text>
                    <Text style={[styles.voteMeta, selected && styles.voteLabelSelected]}>
                      {count} vote{count === 1 ? '' : 's'}
                      {selected ? ' · Your vote' : ''}
                    </Text>
                  </Pressable>
                );
              })}
            </View>
          </Card>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  content: {
    padding: spacing.lg,
    gap: spacing.lg,
    paddingBottom: spacing.xl,
  },
  headerCard: {
    gap: spacing.md,
  },
  headerRow: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  icon: {
    width: 56,
    height: 56,
    borderRadius: radii.md,
    backgroundColor: colors.bgElevated,
    borderWidth: 1,
    borderColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  iconText: {
    color: colors.textSecondary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 24,
  },
  groupImage: {
    width: 56,
    height: 56,
    borderRadius: radii.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.bgElevated,
  },
  headerText: {
    flex: 1,
  },
  groupTitle: {
    color: colors.textPrimary,
    fontFamily: fonts.sansBold,
    fontSize: 18,
  },
  groupDesc: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 14,
    marginTop: 4,
  },
  memberCount: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 12,
    marginTop: spacing.xs,
  },
  actionGrid: {
    gap: spacing.sm,
  },
  list: {
    gap: spacing.lg,
  },
  eventHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  eventTitle: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 16,
  },
  eventTime: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  eventDesc: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 14,
    marginTop: spacing.xs,
  },
  whenPill: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
    marginTop: spacing.sm,
  },
  whenText: {
    color: colors.accentHover,
    fontFamily: fonts.sansMedium,
    fontSize: 14,
  },
  voteRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
    marginTop: spacing.md,
  },
  voteOption: {
    flexDirection: 'column',
    gap: 2,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    borderRadius: radii.pill,
    borderWidth: 1,
    borderColor: colors.borderStrong,
    backgroundColor: colors.bgElevated,
  },
  voteSelected: {
    backgroundColor: colors.accentSoft,
    borderColor: colors.accentBorder,
  },
  voteLabel: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 14,
  },
  voteMeta: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  voteLabelSelected: {
    color: colors.accent,
  },
});