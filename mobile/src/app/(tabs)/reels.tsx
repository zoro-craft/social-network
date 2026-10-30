import { useCallback, useEffect, useRef, useState } from 'react';
import {
  Alert,
  FlatList,
  Pressable,
  Share,
  StyleSheet,
  Text,
  View,
  type ViewToken,
} from 'react-native';
import { useFocusEffect, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import ReelCard from '@/components/ReelCard';
import { EmptyState, Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { API_ORIGIN, apiRaw } from '@/lib/api';
import { Reel, ReelStats } from '@/lib/types';
import { colors, fonts, spacing } from '@/lib/theme';

const PAGE_SIZE = 10;

export default function ReelsScreen() {
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();

  const [reels, setReels] = useState<Reel[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);
  const [muted, setMuted] = useState(true);
  const [height, setHeight] = useState(0);

  const viewedRef = useRef(new Set<number>());
  const snapshotMaxId = useRef<number | null>(null);
  const loadingMoreRef = useRef(false);
  const loadedOnceRef = useRef(false);

  const load = useCallback(
    async (showSpinner: boolean) => {
      if (showSpinner) setLoading(true);

      try {
        const res = await apiRaw<Reel[]>(`/reels?limit=${PAGE_SIZE}&offset=0`);

        if (!res.ok) {
          if (res.status === 401) {
            router.replace('/login');
            return;
          }
          setReels([]);
          setHasMore(false);
          return;
        }

        const maxId = Number(res.headers['x-max-reel-id'] || 0);
        snapshotMaxId.current = maxId || null;
        // A fresh page is a new viewing session: views are counted again per reel.
        viewedRef.current = new Set();
        setActiveIndex(0);
        setReels(res.data || []);
        setHasMore((res.data || []).length === PAGE_SIZE);
      } catch {
        setReels([]);
        setHasMore(false);
      } finally {
        setLoading(false);
      }
    },
    [router],
  );

  // Refetch on tab focus so a reel created on the create screen shows up on return.
  useFocusEffect(
    useCallback(() => {
      if (!auth.authenticated) return;
      load(!loadedOnceRef.current);
      loadedOnceRef.current = true;
    }, [auth.authenticated, load]),
  );

  const loadMore = async () => {
    if (loadingMoreRef.current || !hasMore || snapshotMaxId.current === null) return;

    loadingMoreRef.current = true;
    setLoadingMore(true);

    try {
      const res = await apiRaw<Reel[]>(
        `/reels?before=${snapshotMaxId.current}&limit=${PAGE_SIZE}&offset=${reels.length}`,
      );

      if (!res.ok) {
        setHasMore(false);
        return;
      }

      const batch = res.data || [];
      setReels((prev) => {
        const seen = new Set(prev.map((reel) => reel.id));
        return [...prev, ...batch.filter((reel) => !seen.has(reel.id))];
      });
      setHasMore(batch.length === PAGE_SIZE);
    } catch {
      setHasMore(false);
    } finally {
      loadingMoreRef.current = false;
      setLoadingMore(false);
    }
  };

  const applyStats = useCallback((reelId: number, stats: ReelStats) => {
    setReels((prev) =>
      prev.map((reel) =>
        reel.id === reelId
          ? {
              ...reel,
              like_count: stats.like_count,
              comment_count: stats.comment_count,
              view_count: stats.view_count,
              save_count: stats.save_count,
              share_count: stats.share_count,
              liked: stats.liked,
              saved: stats.saved,
            }
          : reel,
      ),
    );
  }, []);

  // Count a view once per reel per session, when it first scrolls into view.
  useEffect(() => {
    const reel = reels[activeIndex];
    if (!reel || viewedRef.current.has(reel.id)) return;

    viewedRef.current.add(reel.id);

    apiRaw<ReelStats>(`/reels/${reel.id}/view`, { method: 'POST' })
      .then((res) => {
        if (res.ok) applyStats(reel.id, res.data);
      })
      .catch(() => {
        // A failed view count must never interrupt playback.
      });
  }, [activeIndex, reels, applyStats]);

  const act = async (reel: Reel, action: 'like' | 'save' | 'share') => {
    try {
      const res = await apiRaw<ReelStats>(`/reels/${reel.id}/${action}`, { method: 'POST' });

      if (!res.ok) {
        if (res.status === 401) {
          router.replace('/login');
          return null;
        }
        showToast('Could not update reel');
        return null;
      }

      applyStats(reel.id, res.data);
      return res.data;
    } catch {
      showToast('Could not reach the server');
      return null;
    }
  };

  const share = async (reel: Reel) => {
    const stats = await act(reel, 'share');
    if (!stats) return;

    try {
      await Share.share({
        message: `Check out this reel: ${API_ORIGIN}/reels?reel=${reel.id}`,
      });
    } catch {
      // Dismissing the share sheet is not an error worth reporting.
    }
  };

  const remove = async (reel: Reel) => {
    try {
      const res = await apiRaw(`/reels/${reel.id}`, { method: 'DELETE' });

      if (!res.ok) {
        if (res.status === 401) {
          router.replace('/login');
          return;
        }
        showToast('Could not delete reel');
        return;
      }

      setReels((prev) => prev.filter((item) => item.id !== reel.id));
      setActiveIndex(0);
      showToast('Reel deleted', false);
    } catch {
      showToast('Could not reach the server');
    }
  };

  const confirmDelete = (reel: Reel) => {
    Alert.alert('Delete reel', 'This cannot be undone.', [
      { text: 'Cancel', style: 'cancel' },
      { text: 'Delete', style: 'destructive', onPress: () => remove(reel) },
    ]);
  };

  // Ref-held so FlatList never sees these props change (it throws if they do).
  const viewabilityConfig = useRef({ itemVisiblePercentThreshold: 60 }).current;
  const onViewableItemsChanged = useRef(
    ({ viewableItems }: { viewableItems: Array<ViewToken<Reel>> }) => {
      const first = viewableItems[0];
      if (first?.index != null) setActiveIndex(first.index);
    },
  ).current;

  return (
    <View style={styles.screen}>
      <View style={styles.header}>
        <Text style={styles.title}>Reels</Text>
        <Pressable
          onPress={() => router.push('/reels/create')}
          hitSlop={8}
          style={styles.iconButton}
        >
          <Ionicons name="add-circle-outline" size={26} color={colors.textPrimary} />
        </Pressable>
      </View>

      {loading ? (
        <Loading />
      ) : (
        // Measuring the viewport lets each slide be exactly one page, which is what
        // pagingEnabled needs to land on a single reel at a time.
        <View
          style={styles.viewport}
          onLayout={(event) => setHeight(event.nativeEvent.layout.height)}
        >
          {height > 0 ? (
            <FlatList
              data={reels}
              keyExtractor={(item) => String(item.id)}
              pagingEnabled
              showsVerticalScrollIndicator={false}
              viewabilityConfig={viewabilityConfig}
              onViewableItemsChanged={onViewableItemsChanged}
              onEndReached={loadMore}
              onEndReachedThreshold={0.5}
              initialNumToRender={2}
              maxToRenderPerBatch={2}
              windowSize={3}
              removeClippedSubviews
              renderItem={({ item, index }) => (
                <ReelCard
                  reel={item}
                  height={height}
                  active={index === activeIndex}
                  muted={muted}
                  isOwn={item.user_id === auth.user?.id}
                  onToggleMute={() => setMuted((prev) => !prev)}
                  onLike={(reel) => act(reel, 'like')}
                  onSave={(reel) => act(reel, 'save')}
                  onShare={share}
                  onComments={(reel) => router.push(`/reels/${reel.id}`)}
                  onDelete={confirmDelete}
                />
              )}
              ListEmptyComponent={
                <View style={{ height }}>
                  <EmptyState message="No reels yet. Be the first to post one!" />
                </View>
              }
              ListFooterComponent={
                loadingMore ? <Text style={styles.loadingMore}>Loading more...</Text> : null
              }
            />
          ) : null}
        </View>
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
  iconButton: {
    padding: spacing.xs,
  },
  viewport: {
    flex: 1,
  },
  loadingMore: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 13,
    textAlign: 'center',
    paddingVertical: spacing.md,
  },
});
