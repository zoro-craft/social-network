import { Image } from 'expo-image';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import formatDate from '@/lib/format';
import { mediaSource } from '@/lib/api';
import { Post } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';
import Card from '@/components/Card';

interface Props {
  post: Post;
  onComment?: (post: Post) => void;
  onReact?: (post: Post, type: 'like' | 'dislike') => void;
  showName?: boolean;
}

export default function PostCard({ post, onComment, onReact, showName = true }: Props) {
  const image = mediaSource(post.image);

  const hasLiked = post.user_reaction === 'like';
  const hasDisliked = post.user_reaction === 'dislike';

  return (
    <Card>
      <View style={styles.header}>
        <Text style={styles.author}>{showName ? post.full_name : ''}</Text>
        <Text style={styles.time}>{formatDate(post.created_at)}</Text>
      </View>
      <Text style={styles.title}>{post.title}</Text>
      {post.content ? <Text style={styles.content}>{post.content}</Text> : null}
      {image ? (
        <Image source={image} style={styles.image} contentFit="cover" />
      ) : null}
      <View style={styles.actions}>
        <Pressable
          onPress={() => onReact?.(post, 'like')}
          style={[styles.pill, hasLiked && styles.pillLike]}
        >
          <Ionicons
            name={hasLiked ? 'thumbs-up' : 'thumbs-up-outline'}
            size={16}
            color={hasLiked ? colors.success : colors.textSecondary}
          />
          <Text style={[styles.pillText, hasLiked && styles.pillTextLike]}>
            {post.like_count}
          </Text>
        </Pressable>
        <Pressable
          onPress={() => onReact?.(post, 'dislike')}
          style={[styles.pill, hasDisliked && styles.pillDislike]}
        >
          <Ionicons
            name={hasDisliked ? 'thumbs-down' : 'thumbs-down-outline'}
            size={16}
            color={hasDisliked ? colors.danger : colors.textSecondary}
          />
          <Text style={[styles.pillText, hasDisliked && styles.pillTextDislike]}>
            {post.dislike_count}
          </Text>
        </Pressable>
        {onComment ? (
          <Pressable
            onPress={() => onComment(post)}
            style={styles.pill}
          >
            <Ionicons name="chatbubble-outline" size={16} color={colors.textSecondary} />
            <Text style={styles.pillText}>Comment</Text>
          </Pressable>
        ) : null}
      </View>
    </Card>
  );
}

// 'none' is not a valid backend reaction; the server toggles on repeat clicks.

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  author: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 15,
  },
  time: {
    color: colors.textMuted,
    fontFamily: fonts.mono,
    fontSize: 12,
  },
  title: {
    color: colors.textPrimary,
    fontFamily: fonts.sansSemiBold,
    fontSize: 17,
    marginTop: spacing.sm,
  },
  content: {
    color: colors.textBody,
    fontFamily: fonts.sans,
    fontSize: 15,
    lineHeight: 24,
    marginTop: spacing.xs,
  },
  image: {
    width: '100%',
    aspectRatio: 4 / 3,
    borderRadius: radii.md,
    marginTop: spacing.md,
    backgroundColor: colors.bgElevated,
  },
  actions: {
    flexDirection: 'row',
    gap: spacing.sm,
    marginTop: spacing.md,
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
  pillLike: {
    backgroundColor: colors.successSoft,
    borderColor: colors.successBorder,
  },
  pillDislike: {
    backgroundColor: colors.dangerSoft,
    borderColor: colors.dangerBorder,
  },
  pillText: {
    color: colors.textSecondary,
    fontFamily: fonts.sansMedium,
    fontSize: 13,
  },
  pillTextLike: {
    color: colors.success,
  },
  pillTextDislike: {
    color: colors.danger,
  },
});