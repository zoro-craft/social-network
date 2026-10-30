import { useCallback, useEffect, useState } from 'react';
import { FlatList, StyleSheet, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import Avatar from '@/components/Avatar';
import AppButton from '@/components/AppButton';
import HeaderBar from '@/components/HeaderBar';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { apiRequest } from '@/lib/api';
import { UserSearchResult } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

export default function InviteMembersScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();

  const groupId = Number(id);
  const [users, setUsers] = useState<UserSearchResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [invitingId, setInvitingId] = useState<number | null>(null);

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<UserSearchResult[]>(`/groups/${groupId}/users`);
      setUsers(data || []);
    } catch {
      setUsers([]);
    } finally {
      setLoading(false);
    }
  }, [groupId]);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  const invite = async (user: UserSearchResult) => {
    setInvitingId(user.id);
    try {
      await apiRequest(`/groups/${groupId}/invitations`, {
        method: 'POST',
        data: { group_id: groupId, invited_user_id: user.id },
      });
      showToast(`Invitation sent to ${user.first_name || user.full_name}`, false);
    } catch (err: any) {
      showToast(err?.message || 'Could not send invitation');
    } finally {
      setInvitingId(null);
    }
  };

  return (
    <View style={styles.screen}>
      <HeaderBar title="Invite Members" onBack={() => router.back()} />
      {loading ? (
        <Loading />
      ) : (
        <FlatList
          data={users}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          renderItem={({ item }) => (
            <View style={styles.row}>
              <Avatar src={item.avatar} name={item.full_name || item.first_name} size={48} />
              <Text style={styles.name} numberOfLines={1}>
                {item.first_name} {item.last_name}
              </Text>
              <AppButton
                title={invitingId === item.id ? 'Inviting...' : 'Invite'}
                variant="ghost"
                loading={invitingId === item.id}
                onPress={() => invite(item)}
                style={styles.inviteBtn}
              />
            </View>
          )}
          ListEmptyComponent={<EmptyState message="No users available to invite." />}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  list: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    backgroundColor: colors.bgCard,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 14,
    padding: spacing.md,
  },
  name: {
    flex: 1,
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
  inviteBtn: {
    paddingHorizontal: spacing.lg,
  },
});