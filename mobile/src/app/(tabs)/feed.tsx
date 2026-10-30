import { useCallback, useEffect, useRef, useState } from 'react';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import PostCard from '@/components/PostCard';
import StoryBar from '@/components/StoryBar';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { apiRaw } from '@/lib/api';
import { Post } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

const PAGE_SIZE = 10;

export default function FeedScreen() {
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();

  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [snapshotMaxId, setSnapshotMaxId] = useState<number | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const loadingRef = useRef(false);

  const fetchFirstPage = useCallback(async () => {
    setLoading(true);
    try {
      const res = await apiRaw<Post[]>(`/posts?limit=${PAGE_SIZE}&offset=0`);
      if (!res.ok) {
        if (res.status === 401) {
          router.replace('/login');
        }
        setPosts([]);
        return;
      }
      const maxId = Number(res.headers['x-max-post-id'] || 0);
      setSnapshotMaxId(maxId || null);
      setPosts(res.data || []);
      setHasMore((res.data || []).length === PAGE_SIZE);
    } catch {
      setPosts([]);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  }, [router]);

  useEffect(() => {
    if (auth.authenticated) fetchFirstPage();
  }, [auth.authenticated, fetchFirstPage]);

  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMore || snapshotMaxId === null) return;
    loadingRef.current = true;
    setLoadingMore(true);
    try {
      const res = await apiRaw<Post[]>(
        `/posts?before=${snapshotMaxId}&limit=${PAGE_SIZE}&offset=${posts.length}`,
      );
      if (res.ok) {
        const batch = res.data || [];
        setPosts((prev) => {
          const seen = new Set(prev.map((p) => p.id));
          return [...prev, ...batch.filter((p) => !seen.has(p.id))];
        });
        setHasMore(batch.length === PAGE_SIZE);
      } else {
        setHasMore(false);
      }
    } catch {
      setHasMore(false);
    } finally {
      loadingRef.current = false;
      setLoadingMore(false);
    }
  }, [hasMore, snapshotMaxId, posts.length]);

  const react = async (post: Post, type: 'like' | 'dislike') => {
    try {
      const res = await apiRaw<{ like_count: number; dislike_count: number; user_reaction: string | null }>(
        '/posts/react',
        { method: 'POST', data: { post_id: post.id, type } },
      );
      if (res.ok) {
        setPosts((prev) =>
          prev.map((p) =>
            p.id === post.id
              ? {
                  ...p,
                  like_count: res.data.like_count || 0,
                  dislike_count: res.data.dislike_count || 0,
                  user_reaction: res.data.user_reaction || null,
                }
              : p,
          ),
        );
      }
    } catch (err: any) {
      if (err?.status === 401) {
        router.replace('/login');
        return;
      }
      showToast('Could not register reaction');
    }
  };

  const openComments = (post: Post) => router.push(`/posts/${post.id}`);
  const openCreate = () => router.push('/posts/create');
  const openSearch = () => router.push('/search');

  return (
    <View style={styles.screen}>
      <View style={styles.header}>
        <Text style={styles.title}>Feed</Text>
        <View style={styles.headerActions}>
          <Pressable onPress={openSearch} hitSlop={8} style={styles.iconButton}>
            <Ionicons name="search-outline" size={22} color={colors.textPrimary} />
          </Pressable>
          <Pressable onPress={openCreate} hitSlop={8} style={styles.iconButton}>
            <Ionicons name="create-outline" size={24} color={colors.textPrimary} />
          </Pressable>
        </View>
      </View>

      {loading ? (
        <Loading />
      ) : (
        <FlatList
          data={posts}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          renderItem={({ item }) => (
            <PostCard post={item} onComment={openComments} onReact={react} />
          )}
          ListHeaderComponent={<StoryBar />}
          ListEmptyComponent={<EmptyState message="No posts yet. Create the first one!" />}
          ListFooterComponent={
            hasMore ? (
              <Pressable
                onPress={loadMore}
                disabled={loadingMore}
                style={styles.loadMore}
              >
                <Text style={styles.loadMoreText}>
                  {loadingMore ? 'Loading...' : 'Load more'}
                </Text>
              </Pressable>
            ) : null
          }
          onEndReached={loadMore}
          onEndReachedThreshold={0.4}
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
  headerActions: {
    flexDirection: 'row',
    gap: spacing.lg,
  },
  iconButton: {
    padding: spacing.xs,
  },
  list: {
    padding: spacing.lg,
    gap: spacing.lg,
  },
  loadMore: {
    paddingVertical: spacing.md,
    alignItems: 'center',
  },
  loadMoreText: {
    color: colors.accentHover,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
});