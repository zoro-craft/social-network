import { useCallback, useEffect, useState } from 'react';
import { Image } from 'expo-image';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import AppButton from '@/components/AppButton';
import Card from '@/components/Card';
import SegmentedControl from '@/components/SegmentedControl';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { apiRequest, mediaSource } from '@/lib/api';
import { Group } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

export default function GroupsScreen() {
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();

  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<'all' | 'my'>('all');
  const [actionId, setActionId] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await apiRequest<Group[]>('/groups');
      setGroups(data || []);
    } catch {
      setGroups([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  const requestJoin = async (group: Group) => {
    setActionId(group.id);
    try {
      const data = await apiRequest<{ message?: string }>(`/groups/${group.id}/join`, {
        method: 'POST',
      });
      setGroups((current) =>
        current.map((item) => (item.id === group.id ? { ...item, is_pending: true } : item))
      );
      showToast(data?.message || 'Request sent', false);
      await load();
    } catch (err: any) {
      if (err?.status === 409) {
        setGroups((current) =>
          current.map((item) => (item.id === group.id ? { ...item, is_pending: true } : item))
        );
      }
      showToast(err?.message || 'Could not send join request');
    } finally {
      setActionId(null);
    }
  };

  const cancelRequest = async (group: Group) => {
    setActionId(group.id);
    try {
      const data = await apiRequest<{ message?: string }>(`/groups/${group.id}/join`, {
        method: 'DELETE',
      });
      setGroups((current) =>
        current.map((item) => (item.id === group.id ? { ...item, is_pending: false } : item))
      );
      showToast(data?.message || 'Join request cancelled', false);
    } catch (err: any) {
      showToast(err?.message || 'Could not cancel join request');
    } finally {
      setActionId(null);
    }
  };

  const visible = tab === 'all' ? groups.filter((g) => !g.is_member) : groups.filter((g) => g.is_member);

  return (
    <View style={styles.screen}>
      <View style={styles.header}>
        <Text style={styles.title}>Groups</Text>
        <Pressable onPress={() => router.push('/groups/create')} hitSlop={8} style={styles.createBtn}>
          <Ionicons name="add" size={24} color={colors.textPrimary} />
          <Text style={styles.createText}>Create Group</Text>
        </Pressable>
      </View>

      <View style={styles.tabWrap}>
        <SegmentedControl
          options={[
            { label: 'All Groups', value: 'all' },
            { label: 'My Groups', value: 'my' },
          ]}
          value={tab}
          onChange={(v) => setTab(v as 'all' | 'my')}
        />
      </View>

      {loading ? (
        <Loading />
      ) : (
        <FlatList
          data={visible}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          renderItem={({ item }) => {
            const busy = actionId === item.id;
            return (
              <Card>
                <View style={styles.groupRow}>
                  <GroupAvatar group={item} />
                  <View style={styles.groupInfo}>
                    <Text style={styles.groupName}>{item.titel}</Text>
                    <Text style={styles.groupDesc} numberOfLines={2}>
                      {item.description}
                    </Text>
                    <View style={styles.groupMeta}>
                      {typeof item.member_count === 'number' && (
                        <Text style={styles.metaText}>
                          {item.member_count} member{item.member_count === 1 ? '' : 's'}
                        </Text>
                      )}
                      {item.is_pending && <Text style={styles.pendingText}>Request pending</Text>}
                    </View>
                  </View>
                </View>
                {item.is_member ? (
                  <AppButton
                    title="Enter Group"
                    variant="ghost"
                    onPress={() => router.push(`/groups/${item.id}`)}
                  />
                ) : item.is_pending ? (
                  <AppButton
                    title={busy ? 'Cancelling...' : 'Cancel request'}
                    variant="danger"
                    loading={busy}
                    onPress={() => cancelRequest(item)}
                  />
                ) : (
                  <AppButton
                    title={busy ? 'Requesting...' : 'Request to join'}
                    variant="primary"
                    loading={busy}
                    onPress={() => requestJoin(item)}
                  />
                )}
              </Card>
            );
          }}
          ListEmptyComponent={

            <EmptyState
              message={tab === 'all' ? 'No groups yet. Create the first one!' : 'You have not joined any groups yet.'}
            />
          }
        />
      )}
    </View>
  );
}

function GroupAvatar({ group }: { group: Group }) {
  const source = mediaSource(group.image);
  const [failed, setFailed] = useState(false);
  const initial = (group.titel || '?').charAt(0).toUpperCase();

  if (!source || failed) {
    return (
      <View style={styles.groupIcon}>
        <Text style={styles.groupIconText}>{initial}</Text>
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

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
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
  createBtn: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.xs,
  },
  createText: {
    color: colors.accentHover,
    fontFamily: fonts.sansSemiBold,
    fontSize: 14,
  },
  tabWrap: {
    padding: spacing.lg,
  },
  list: {
    paddingHorizontal: spacing.lg,
    gap: spacing.lg,
    paddingBottom: spacing.xl,
  },
  groupRow: {
    flexDirection: 'row',
    gap: spacing.md,
    marginBottom: spacing.md,
  },
  groupIcon: {
    width: 48,
    height: 48,
    borderRadius: radii.md,
    backgroundColor: colors.bgElevated,
    borderWidth: 1,
    borderColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  groupImage: {
    width: 48,
    height: 48,
    borderRadius: radii.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.bgElevated,
  },
  groupIconText: {
    color: colors.textSecondary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 20,
  },
  groupInfo: {
    flex: 1,
  },
  groupName: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 16,
  },
  groupDesc: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 14,
    marginTop: 2,
  },
  groupMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    flexWrap: 'wrap',
    gap: spacing.sm,
    marginTop: spacing.xs,
  },
  metaText: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 12,
  },
  pendingText: {
    color: colors.pending,
    fontFamily: fonts.sansSemiBold,
    fontSize: 12,
  },
});
