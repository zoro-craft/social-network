import { useCallback, useEffect, useRef, useState } from 'react';
import {
  FlatList,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  Share,
  StyleSheet,
  Text,
  useWindowDimensions,
  View,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { VideoView, useVideoPlayer } from 'expo-video';
import AppButton from '@/components/AppButton';
import AppInput from '@/components/AppInput';
import Avatar from '@/components/Avatar';
import HeaderBar from '@/components/HeaderBar';
import { Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { API_ORIGIN, apiRaw, mediaSource } from '@/lib/api';
import formatDate from '@/lib/format';
import { Reel, ReelComment, ReelStats } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

const PAGE_SIZE = 20;

export default function ReelDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();
  const { width } = useWindowDimensions();

  const reelId = Number(id);

  const [reel, setReel] = useState<Reel | null>(null);
  const [comments, setComments] = useState<ReelComment[]>([]);
  const [loading, setLoading] = useState(true);
  const [snapshotMaxId, setSnapshotMaxId] = useState<number | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [content, setContent] = useState('');
  const [sending, setSending] = useState(false);

  const loadingMoreRef = useRef(false);

  // Recreated whenever the loaded reel changes, which is exactly when the source
  // becomes available; it starts playing because this is the focused video.
  const source = reel ? mediaSource(reel.video) : null;
  const player = useVideoPlayer(source, (instance) => {
    instance.loop = true;
    instance.play();
  });

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [reelRes, commentsRes] = await Promise.all([
        apiRaw<Reel>(`/reels/${reelId}`),
        apiRaw<ReelComment[]>(`/reels/${reelId}/comments?limit=${PAGE_SIZE}&offset=0`),
      ]);

      setReel(reelRes.ok ? reelRes.data : null);

      if (commentsRes.ok) {
        const maxId = Number(commentsRes.headers['x-max-comment-id'] || 0);
        setSnapshotMaxId(maxId || null);
        setComments(commentsRes.data || []);
        setHasMore((commentsRes.data || []).length === PAGE_SIZE);
      } else {
        setComments([]);
        setHasMore(false);
      }
    } catch {
      setReel(null);
      setComments([]);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  }, [reelId]);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  const loadMore = async () => {
    if (loadingMoreRef.current || !hasMore || snapshotMaxId === null) return;

    loadingMoreRef.current = true;
    setLoadingMore(true);

    try {
      const res = await apiRaw<ReelComment[]>(
        `/reels/${reelId}/comments?before=${snapshotMaxId}&limit=${PAGE_SIZE}&offset=${comments.length}`,
      );

      if (res.ok) {
        setComments((prev) => [...prev, ...(res.data || [])]);
        setHasMore((res.data || []).length === PAGE_SIZE);
      } else {
        setHasMore(false);
      }
    } catch {
      setHasMore(false);
    } finally {
      loadingMoreRef.current = false;
      setLoadingMore(false);
    }
  };

  const act = async (action: 'like' | 'save' | 'share') => {
    if (!reel) return null;

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

      setReel((prev) =>
        prev
          ? {
              ...prev,
              like_count: res.data.like_count,
              comment_count: res.data.comment_count,
              view_count: res.data.view_count,
              save_count: res.data.save_count,
              share_count: res.data.share_count,
              liked: res.data.liked,
              saved: res.data.saved,
            }
          : prev,
      );
      return res.data;
    } catch {
      showToast('Could not reach the server');
      return null;
    }
  };

  const share = async () => {
    if (!reel) return;

    const stats = await act('share');
    if (!stats) return;

    try {
      await Share.share({
        message: `Check out this reel: ${API_ORIGIN}/reels?reel=${reel.id}`,
      });
    } catch {
      // Dismissing the share sheet is not an error worth reporting.
    }
  };

  const addComment = async () => {
    const text = content.trim();
    if (!text) return;

    setSending(true);
    try {
      const form = new FormData();
      form.append('content', text);

      const res = await apiRaw<ReelComment>(`/reels/${reelId}/comments`, {
        method: 'POST',
        formData: form,
      });

      if (res.ok) {
        setComments((prev) => [res.data, ...prev]);
        setReel((prev) =>
          prev ? { ...prev, comment_count: prev.comment_count + 1 } : prev,
        );
        setContent('');
      } else if (res.status === 401) {
        router.replace('/login');
      } else {
        showToast((res.data as any)?.error || 'Could not add comment');
      }
    } catch {
      showToast('Could not reach the server');
    } finally {
      setSending(false);
    }
  };

  const header = (
    <View style={styles.header}>
      <View style={[styles.videoWrap, { height: Math.min(width * 1.2, 380) }]}>
        <VideoView
          player={player}
          style={styles.video}
          contentFit="contain"
          nativeControls={false}
        />
      </View>

      {reel ? (
        <View style={styles.meta}>
          <Pressable
            style={styles.author}
            onPress={() => router.push(`/profile/${reel.user_id}`)}
          >
            <Avatar src={reel.avatar} name={reel.full_name} size={36} />
            <View>
              <Text style={styles.authorName}>{reel.full_name}</Text>
              <Text style={styles.time}>{formatDate(reel.created_at)}</Text>
            </View>
          </Pressable>

          <View style={styles.pillRow}>
            <Pressable
              onPress={() => act('like')}
              style={[styles.pill, reel.liked && styles.pillActive]}
            >
              <Ionicons
                name={reel.liked ? 'heart' : 'heart-outline'}
                size={16}
                color={reel.liked ? colors.danger : colors.textSecondary}
              />
              <Text style={styles.pillText}>{reel.like_count}</Text>
            </Pressable>

            <Pressable onPress={share} style={styles.pill}>
              <Ionicons name="paper-plane-outline" size={16} color={colors.textSecondary} />
              <Text style={styles.pillText}>{reel.share_count}</Text>
            </Pressable>

            <Pressable
              onPress={() => act('save')}
              style={[styles.pill, reel.saved && styles.pillActive]}
            >
              <Ionicons
                name={reel.saved ? 'bookmark' : 'bookmark-outline'}
                size={16}
                color={reel.saved ? colors.textPrimary : colors.textSecondary}
              />
              <Text style={styles.pillText}>{reel.save_count}</Text>
            </Pressable>

            <View style={styles.pill}>
              <Ionicons name="eye-outline" size={16} color={colors.textSecondary} />
              <Text style={styles.pillText}>{reel.view_count}</Text>
            </View>
          </View>

          {reel.caption ? <Text style={styles.caption}>{reel.caption}</Text> : null}
        </View>
      ) : null}

      <Text style={styles.commentsLabel}>Comments</Text>
    </View>
  );

  return (
    <KeyboardAvoidingView
      style={styles.flex}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <HeaderBar title="Reel" onBack={() => router.back()} />

      {loading ? (
        <Loading />
      ) : (
        <FlatList
          data={comments}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          ListHeaderComponent={header}
          renderItem={({ item }) => (
            <View style={styles.comment}>
              <Pressable onPress={() => router.push(`/profile/${item.user_id}`)}>
                <Avatar src={item.avatar} name={item.first_name} size={32} />
              </Pressable>
              <View style={styles.commentBody}>
                <View style={styles.commentMeta}>
                  <Text style={styles.commentName}>{item.first_name}</Text>
                  <Text style={styles.commentTime}>{formatDate(item.created_at)}</Text>
                </View>
                <Text style={styles.commentContent}>{item.content}</Text>
              </View>
            </View>
          )}
          ListEmptyComponent={
            <Text style={styles.empty}>No comments yet. Be the first!</Text>
          }
          ListFooterComponent={
            hasMore ? (
              <Pressable onPress={loadMore} disabled={loadingMore} style={styles.loadMore}>
                <Text style={styles.loadMoreText}>
                  {loadingMore ? 'Loading...' : 'Load older comments'}
                </Text>
              </Pressable>
            ) : null
          }
        />
      )}

      <View style={styles.composer}>
        <View style={styles.composerInput}>
          <AppInput
            value={content}
            onChangeText={setContent}
            placeholder="Write a comment..."
            maxLength={1000}
          />
        </View>
        <AppButton
          title={sending ? '...' : 'Post'}
          onPress={addComment}
          loading={sending}
          disabled={!content.trim()}
        />
      </View>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  flex: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  list: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  header: {
    gap: spacing.md,
  },
  videoWrap: {
    width: '100%',
    borderRadius: radii.lg,
    overflow: 'hidden',
    backgroundColor: '#000',
    borderWidth: 1,
    borderColor: colors.border,
  },
  // Spelled out rather than Spread from StyleSheet.absoluteFillObject, which this
  // React Native version no longer exposes in its types.
  video: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  meta: {
    gap: spacing.sm,
    backgroundColor: colors.bgCard,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radii.lg,
    padding: spacing.lg,
  },
  author: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  authorName: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
  time: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  pillRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  pill: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    borderRadius: radii.pill,
    backgroundColor: colors.bgElevated,
    borderWidth: 1,
    borderColor: colors.border,
  },
  pillActive: {
    borderColor: colors.borderStrong,
  },
  pillText: {
    color: colors.textSecondary,
    fontFamily: fonts.sansMedium,
    fontSize: 13,
  },
  caption: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 15,
    lineHeight: 21,
  },
  commentsLabel: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
  comment: {
    flexDirection: 'row',
    gap: spacing.md,
    backgroundColor: colors.bgCard,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radii.lg,
    padding: spacing.md,
  },
  commentBody: {
    flex: 1,
  },
  commentMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  commentName: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 14,
  },
  commentTime: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 11,
  },
  commentContent: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 15,
    marginTop: spacing.xs,
  },
  empty: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    fontSize: 14,
    textAlign: 'center',
    paddingVertical: spacing.xl,
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
  composer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    padding: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  composerInput: {
    flex: 1,
  },
});
