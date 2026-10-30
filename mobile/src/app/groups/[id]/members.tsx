import { useCallback, useEffect, useState } from 'react';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import Avatar from '@/components/Avatar';
import HeaderBar from '@/components/HeaderBar';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { apiRequest } from '@/lib/api';
import { UserSearchResult } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

export default function GroupMembersScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth } = useAuth();

  const groupId = Number(id);
  const [members, setMembers] = useState<UserSearchResult[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<UserSearchResult[]>(`/groups/${groupId}/members`);
      setMembers(data || []);
    } catch {
      setMembers([]);
    } finally {
      setLoading(false);
    }
  }, [groupId]);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  return (
    <View style={styles.screen}>
      <HeaderBar title="Group Members" onBack={() => router.back()} />
      {loading ? (
        <Loading />
      ) : (
        <FlatList
          data={members}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          renderItem={({ item }) => {
            // full_name already contains both parts, so only join them when we
            // are falling back to the split fields.
            const name =
              item.full_name ||
              [item.first_name, item.last_name].filter(Boolean).join(' ');
            return (
              <Pressable
                style={styles.row}
                onPress={() => router.push(`/profile/${item.id}`)}
              >
                <Avatar src={item.avatar} name={name} size={48} />
                <Text style={styles.name} numberOfLines={1}>
                  {name}
                </Text>
              </Pressable>
            );
          }}
          ListEmptyComponent={<EmptyState message="No members yet." />}
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
    color: colors.textPrimary,
    fontFamily: fonts.sans,
    fontSize: 15,
  },
});