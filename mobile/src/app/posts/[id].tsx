import { useCallback, useEffect, useRef, useState } from 'react';
import {
  FlatList,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { Image } from 'expo-image';
import { useLocalSearchParams, useRouter } from 'expo-router';
import AppInput from '@/components/AppInput';
import AppButton from '@/components/AppButton';
import Avatar from '@/components/Avatar';
import HeaderBar from '@/components/HeaderBar';
import PostCard from '@/components/PostCard';
import PhotoPicker from '@/components/PhotoPicker';
import { Loading } from '@/components/StateViews';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/context/ToastContext';
import { apiRaw, mediaSource } from '@/lib/api';
import { appendAsset, PickedImage } from '@/hooks/useImagePicker';
import formatDate from '@/lib/format';
import { Comment, Post } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

const PAGE_SIZE = 10;

export default function PostCommentsScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { auth } = useAuth();
  const showToast = useToast();

  const postId = Number(id);
  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [snapshotMaxId, setSnapshotMaxId] = useState<number | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);

  const [content, setContent] = useState('');
  const [image, setImage] = useState<PickedImage | null>(null);
  const [sending, setSending] = useState(false);
  const loadingMoreRef = useRef(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [postRes, commentsRes] = await Promise.all([
        apiRaw<Post>(`/posts/${postId}`),
        apiRaw<Comment[]>(`/posts/${postId}/comment?limit=${PAGE_SIZE}&offset=0`),
      ]);
      setPost(postRes.ok ? postRes.data : null);
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
      setPost(null);
      setComments([]);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  }, [postId]);

  useEffect(() => {
    if (auth.authenticated) load();
  }, [auth.authenticated, load]);

  const loadMore = async () => {
    if (loadingMoreRef.current || !hasMore || snapshotMaxId === null) return;
    loadingMoreRef.current = true;
    setLoadingMore(true);
    try {
      const res = await apiRaw<Comment[]>(
        `/posts/${postId}/comment?before=${snapshotMaxId}&limit=${PAGE_SIZE}&offset=${comments.length}`,
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

  const addComment = async () => {
    if (!content.trim() && !image) {
      showToast('Write a comment first');
      return;
    }
    setSending(true);
    try {
      const formData = new FormData();
      formData.append('content', content.trim());
      if (image) appendAsset(formData, 'image', image);
      const res = await apiRaw<Comment>(`/posts/${postId}/comment`, { method: 'POST', formData });
      if (res.ok) {
        setComments((prev) => [res.data, ...prev]);
        setContent('');
        setImage(null);
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
    <View>
      {post ? <PostCard post={post} /> : null}
      <View style={styles.compose}>
        <AppInput
          value={content}
          onChangeText={setContent}
          placeholder="Write a comment..."
          multiline
          maxLength={2000}
        />
        <PhotoPicker label="Attach image" value={image} onChange={setImage} />
        <AppButton
          title={sending ? 'Posting...' : 'Post'}
          onPress={addComment}
          loading={sending}
        />
      </View>
      {comments.length > 0 && <Text style={styles.commentsLabel}>Comments</Text>}
    </View>
  );

  return (
    <KeyboardAvoidingView
      style={styles.flex}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <HeaderBar title="Post & Comments" onBack={() => router.back()} />
      {loading ? (
        <Loading />
      ) : (
        <FlatList
          data={comments}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.list}
          ListHeaderComponent={header}
          renderItem={({ item }) => <CommentRow comment={item} onOpenUser={() => router.push(`/profile/${item.user_id}`)} />}
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
    </KeyboardAvoidingView>
  );
}

function CommentRow({ comment, onOpenUser }: { comment: Comment; onOpenUser: () => void }) {
  const image = mediaSource(comment.image);
  return (
    <View style={styles.comment}>
      <Pressable onPress={onOpenUser}>
        <Avatar src={comment.avatar} name={comment.first_name} size={32} />
      </Pressable>
      <View style={styles.commentBody}>
        <View style={styles.commentMeta}>
          <Pressable onPress={onOpenUser}>
            <Text style={styles.commentName}>{comment.first_name}</Text>
          </Pressable>
          {comment.created_at ? (
            <Text style={styles.commentTime}>{formatDate(comment.created_at)}</Text>
          ) : null}
        </View>
        <Text style={styles.commentContent}>{comment.content}</Text>
        {image ? (
          <Image source={image} style={styles.commentImage} contentFit="cover" />
        ) : null}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  flex: {
    flex: 1,
    backgroundColor: colors.bg,
  },
  list: {
    padding: spacing.lg,
    gap: spacing.lg,
  },
  compose: {
    gap: spacing.md,
    backgroundColor: colors.bgCard,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radii.lg,
    padding: spacing.lg,
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
  commentImage: {
    width: '100%',
    height: 180,
    borderRadius: radii.md,
    marginTop: spacing.sm,
    backgroundColor: colors.bgElevated,
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