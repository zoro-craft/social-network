import { Image } from 'expo-image';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { mediaSource } from '@/lib/api';
import { Post } from '@/lib/types';
import { colors, fonts, radii, spacing } from '@/lib/theme';

interface Props {
  posts?: Post[] | null;
  isPrivateHidden?: boolean;
  onOpenPost: (post: Post) => void;
}

export default function PostGrid({ posts, isPrivateHidden = false, onOpenPost }: Props) {
  if (isPrivateHidden) {
    return <Text style={styles.empty}>🔒 Private account.</Text>;
  }
  if (!posts || posts.length === 0) {
    return <Text style={styles.empty}>No posts yet.</Text>;
  }

  return (
    <View style={styles.grid}>
      {posts.map((post) => {
        const image = mediaSource(post.image);
        return (
          <Pressable
            key={post.id}
            onPress={() => onOpenPost(post)}
            style={styles.tile}
          >
            {image ? (
              <Image source={image} style={styles.tileImage} contentFit="cover" />
            ) : (
              <View style={styles.tileTextWrap}>
                <Text style={styles.tileText} numberOfLines={2}>
                  {post.title}
                </Text>
              </View>
            )}
          </Pressable>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  grid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  tile: {
    // 3 columns with room to spare for the 2 gaps, so a tile never wraps
    // onto a row of its own on narrow phones.
    width: '30%',
    aspectRatio: 1,
    borderRadius: radii.md,
    overflow: 'hidden',
    backgroundColor: colors.bgElevated,
    borderWidth: 1,
    borderColor: colors.border,
  },
  tileImage: {
    width: '100%',
    height: '100%',
  },
  tileTextWrap: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: spacing.xs,
  },
  tileText: {
    color: colors.textSecondary,
    fontFamily: fonts.sans,
    fontSize: 13,
    textAlign: 'center',
  },
  empty: {
    color: colors.textMuted,
    fontFamily: fonts.sans,
    textAlign: 'center',
    paddingVertical: spacing.xl,
  },
});