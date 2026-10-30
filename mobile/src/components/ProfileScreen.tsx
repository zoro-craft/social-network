import { useCallback, useEffect, useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import Avatar from '@/components/Avatar';
import AppButton from '@/components/AppButton';
import Card from '@/components/Card';
import FollowButton from '@/components/FollowButton';
import PostGrid from '@/components/PostGrid';
import UsersModal from '@/components/UsersModal';
import EditProfileModal from '@/components/EditProfileModal';
import { Loading } from '@/components/StateViews';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { ProfileData, UserSearchResult } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

interface Props {
  userId: number;
}

export default function ProfileScreen({ userId }: Props) {
  const router = useRouter();
  const showToast = useToast();
  const [profile, setProfile] = useState<ProfileData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [following, setFollowing] = useState(false);
  const [modalType, setModalType] = useState<'followers' | 'following' | null>(null);
  const [usersList, setUsersList] = useState<UserSearchResult[]>([]);
  const [isEditing, setIsEditing] = useState(false);

  const loadProfile = useCallback(async () => {
    try {
      const data = await apiRequest<ProfileData>(`/profile/${userId}`);
      if (!data.user) throw new Error('User not found');
      const relation = ['self', 'accepted', 'pending', 'none'].includes(data.relation)
        ? data.relation
        : 'none';
      setProfile({ ...data, relation });
      setError(null);
    } catch (err: any) {
      setProfile(null);
      setError(err?.message || 'User loading failed');
    }
  }, [userId]);

  useEffect(() => {
    loadProfile();
  }, [loadProfile]);

  const toggleFollow = async () => {
    setFollowing(true);
    try {
      await apiRequest(`/follow/${userId}`, { method: 'POST' });
      await loadProfile();
    } catch (err: any) {
      showToast(err?.message || 'Could not update follow state');
    } finally {
      setFollowing(false);
    }
  };

  const openModal = async (type: 'followers' | 'following') => {
    if (!profile) return;
    const hidden = profile.user.is_private && profile.relation !== 'self' && profile.relation !== 'accepted';
    if (hidden) return;
    setModalType(type);
    try {
      const data = await apiRequest<UserSearchResult[]>(`/profile/${userId}/${type}`);
      setUsersList(data || []);
    } catch (err: any) {
      setModalType(null);
      showToast(err?.message || 'Failed to load list');
    }
  };

  if (error) {
    return (
      <View style={styles.errorWrap}>
        <Card>
          <Text style={styles.heading}>Profile unavailable</Text>
          <Text style={styles.muted}>{error}</Text>
          <AppButton title="Back to feed" variant="ghost" onPress={() => router.replace('/(tabs)/feed')} />
        </Card>
      </View>
    );
  }

  if (!profile) {
    return <Loading />;
  }

  const { user, relation, followers_count, following_count, posts } = profile;
  const isPrivateHidden = user.is_private && relation !== 'self' && relation !== 'accepted';

  const editConfig =
    relation === 'self'
      ? { text: 'Edit Profile', onPress: () => setIsEditing(true) }
      : null;
  const followProps = editConfig ? null : {
    relation,
    isPrivate: user.is_private,
    onToggle: toggleFollow,
    disabled: following,
  };

  // A profile is taller than the viewport (avatar + stats + post grid), so the
  // whole thing scrolls as one column.
  return (
    <ScrollView contentContainerStyle={styles.wrap}>
      <Card style={styles.sidebar}>
        <Avatar src={user.avatar} name={user.first_name} size={112} />
        <Text style={styles.name}>
          {user.first_name} {user.last_name}
        </Text>
        <Text style={styles.handle}>@{user.nickname || user.first_name}</Text>
        <Text style={styles.badge}>
          {user.is_private ? 'Private Profile' : 'Public Profile'}
        </Text>
        {!isPrivateHidden && (
          <>
            {user.date_of_birth ? (
              <Text style={styles.muted}>Born: {user.date_of_birth.split('T')[0]}</Text>
            ) : null}
            {user.about_me ? <Text style={styles.muted}>{user.about_me}</Text> : null}
          </>
        )}
        {editConfig ? (
          <AppButton
            title={editConfig.text}
            variant="ghost"
            onPress={editConfig.onPress}
          />
        ) : (
          <FollowButton {...followProps!} />
        )}
      </Card>

      <View style={styles.right}>
        <View style={styles.stats}>
          <StatCard label="Followers" count={followers_count || 0} onPress={() => openModal('followers')} />
          <StatCard label="Following" count={following_count || 0} onPress={() => openModal('following')} />
          <StatCard
            label="Joined"
            count={user.created_at ? user.created_at.split('-')[0] : 'Recently'}
          />
        </View>

        <Card style={styles.postsCard}>
          <Text style={styles.postsHeading}>Posts</Text>
          <PostGrid
            posts={posts}
            isPrivateHidden={isPrivateHidden}
            onOpenPost={(post) => router.push(`/posts/${post.id}`)}
          />
        </Card>
      </View>

      {isEditing && (
        <EditProfileModal
          user={user}
          onClose={() => setIsEditing(false)}
          onShowError={showToast}
          onSaved={async () => {
            setIsEditing(false);
            await loadProfile();
            showToast('Profile updated', false);
          }}
        />
      )}

      {modalType && (
        <UsersModal
          type={modalType}
          users={usersList}
          onClose={() => setModalType(null)}
          onOpenUser={(id) => router.push(`/profile/${id}`)}
        />
      )}
    </ScrollView>
  );
}

function StatCard({ label, count, onPress }: { label: string; count: string | number; onPress?: () => void }) {
  return (
    <Card style={styles.statCard} onPress={onPress}>
      <Text style={styles.statLabel}>{label}</Text>
      <Text style={styles.statCount}>{count}</Text>
    </Card>
  );
}

const styles = StyleSheet.create({
  wrap: {
    padding: spacing.lg,
    gap: spacing.lg,
    paddingBottom: spacing.xl,
  },
  errorWrap: {
    flex: 1,
    justifyContent: 'center',
    padding: spacing.lg,
  },
  sidebar: {
    alignItems: 'center',
    gap: spacing.sm,
  },
  name: {
    color: colors.textPrimary,
    fontFamily: fonts.sansBold,
    fontSize: 20,
    letterSpacing: -0.02,
  },
  handle: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
  badge: {
    color: colors.textMuted,
    fontFamily: fonts.sansSemiBold,
    fontSize: 11,
    letterSpacing: 0.08,
    textTransform: 'uppercase',
  },
  muted: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 14,
    textAlign: 'center',
  },
  right: {
    gap: spacing.lg,
  },
  stats: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  statCard: {
    flex: 1,
    alignItems: 'center',
    paddingVertical: spacing.lg,
  },
  statLabel: {
    color: colors.textMuted,
    fontFamily: fonts.sansSemiBold,
    fontSize: 11,
    letterSpacing: 0.08,
    textTransform: 'uppercase',
  },
  statCount: {
    color: colors.textPrimary,
    fontFamily: fonts.mono,
    fontSize: 20,
  },
  postsCard: {
    padding: spacing.lg,
  },
  postsHeading: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 17,
    marginBottom: spacing.md,
  },
  heading: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 18,
    marginBottom: spacing.sm,
  },
});