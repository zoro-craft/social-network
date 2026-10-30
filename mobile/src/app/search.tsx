import { useCallback, useEffect, useState } from 'react';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import AppInput from '@/components/AppInput';
import Avatar from '@/components/Avatar';
import HeaderBar from '@/components/HeaderBar';
import { EmptyState } from '@/components/StateViews';
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { apiRequest } from '@/lib/api';
import { UserSearchResult } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

const DEBOUNCE_MS = 300;

export default function SearchScreen() {
  const router = useRouter();
  const { auth, loading } = useRequireAuth();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<UserSearchResult[]>([]);
  const [searching, setSearching] = useState(false);
  const [done, setDone] = useState(false);

  const search = useCallback(async (q: string) => {
    if (!q.trim()) {
      setResults([]);
      setSearching(false);
      setDone(false);
      return;
    }
    setSearching(true);
    try {
      const data = await apiRequest<UserSearchResult[]>(`/users/search?q=${encodeURIComponent(q.trim())}`);
      setResults(data || []);
    } catch {
      setResults([]);
    } finally {
      setSearching(false);
      setDone(true);
    }
  }, []);

  useEffect(() => {
    const handle = setTimeout(() => search(query), DEBOUNCE_MS);
    return () => clearTimeout(handle);
  }, [query, search]);

  if (loading || !auth.authenticated) return null;

  return (
    <View style={styles.screen}>
      <HeaderBar title="Search" onBack={() => router.back()} />
      <View style={styles.searchWrap}>
        <AppInput
          value={query}
          onChangeText={setQuery}
          placeholder="Search users..."
          autoCapitalize="none"
        />
      </View>

      {searching ? (
        <EmptyState message="Searching..." />
      ) : !done ? null : results.length === 0 ? (
        <EmptyState message="No users found" />
      ) : (
        <FlatList
          data={results}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          renderItem={({ item }) => (
            <Pressable
              style={styles.row}
              onPress={() => router.push(`/profile/${item.id}`)}
            >
              <Avatar src={item.avatar} name={item.full_name} size={40} />
              <Text style={styles.name}>{item.full_name}</Text>
            </Pressable>
          )}
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
  searchWrap: {
    padding: spacing.lg,
  },
  list: {
    paddingHorizontal: spacing.lg,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  name: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
});